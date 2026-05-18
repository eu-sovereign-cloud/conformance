package storage

import (
	"math/rand"
	"strings"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockstorage "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/storage"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// ImageConstraintsValidationV1TestSuite verifies that Image resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
//   - spec.cpuArchitecture: enum [amd64, arm64] (ImageSpec)
//   - spec.initializer: enum [none, cloudinit-22] (ImageSpec)
//   - spec.boot: enum [UEFI, BIOS] (ImageSpec)
type ImageConstraintsValidationV1TestSuite struct {
	suites.RegionalTestSuite
	StorageSkus []string

	params *params.ImageConstraintsValidationV1Params
}

func CreateImageConstraintsValidationV1TestSuite(regionalTestSuite suites.RegionalTestSuite, storageSkus []string) *ImageConstraintsValidationV1TestSuite {
	suite := &ImageConstraintsValidationV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       storageSkus,
	}
	suite.ScenarioName = constants.ImageConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *ImageConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.StorageParentSuite)

	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]

	workspaceName := generators.GenerateWorkspaceName()
	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)
	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)
	storageSize := constants.BlockStorageInitialSize

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Workspace for conformance testing",
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Annotations(schema.Annotations{
			"description": "BlockStorage for conformance testing",
		}).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: storageSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	buildImage := func(name string, labels schema.Labels, annotations schema.Annotations, spec *schema.ImageSpec) *schema.Image {
		image, err := builders.NewImageBuilder().
			Name(name).
			Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Region(suite.Region).
			Labels(labels).
			Annotations(annotations).
			Spec(spec).Build()
		if err != nil {
			t.Fatalf("Failed to build Image: %v", err)
		}
		return image
	}

	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}
	baseSpec := func() *schema.ImageSpec {
		return &schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}
	}

	invalidCpuArchSpec := baseSpec()
	invalidCpuArchSpec.CpuArchitecture = "x86_64"

	invalidInitializerSpec := baseSpec()
	invalidInitializerSpec.Initializer = "invalid-initializer"

	invalidBootSpec := baseSpec()
	invalidBootSpec.Boot = "PXE"

	p := &params.ImageConstraintsValidationV1Params{
		Workspace:    workspace,
		BlockStorage: blockStorage,
		OverLengthNameImage: buildImage(
			strings.Repeat("a", 129),
			baseLabels,
			schema.Annotations{"description": "Image with over-length name"},
			baseSpec(),
		),
		InvalidPatternNameImage: buildImage(
			"Invalid-Name-With-Uppercase",
			baseLabels,
			schema.Annotations{"description": "Image with non-kebab-case name"},
			baseSpec(),
		),
		OverLengthLabelValueImage: buildImage(
			generators.GenerateImageName(),
			schema.Labels{
				constants.EnvLabel: constants.EnvConformanceLabel,
				"constraint-test":  strings.Repeat("x", 64),
			},
			schema.Annotations{"description": "Image with over-length label value"},
			baseSpec(),
		),
		OverLengthAnnotationImage: buildImage(
			generators.GenerateImageName(),
			baseLabels,
			schema.Annotations{
				"description":     "Image with over-length annotation value",
				"long-annotation": strings.Repeat("y", 1025),
			},
			baseSpec(),
		),
		InvalidCpuArchitectureImage: buildImage(
			generators.GenerateImageName(),
			baseLabels,
			schema.Annotations{"description": "Image with invalid cpuArchitecture enum value"},
			invalidCpuArchSpec,
		),
		InvalidInitializerImage: buildImage(
			generators.GenerateImageName(),
			baseLabels,
			schema.Annotations{"description": "Image with invalid initializer enum value"},
			invalidInitializerSpec,
		),
		InvalidBootImage: buildImage(
			generators.GenerateImageName(),
			baseLabels,
			schema.Annotations{"description": "Image with invalid boot enum value"},
			invalidBootSpec,
		),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockstorage.ConfigureImageConstraintsViolationsV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ImageConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.StorageProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage))
	suite.ConfigureDepends(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	expectWorkspaceAnnotations := workspace.Annotations
	expectWorkspaceExtensions := workspace.Extensions
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Annotations:    expectWorkspaceAnnotations,
			Extensions:     expectWorkspaceExtensions,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			Labels:   expectWorkspaceLabels,
			Metadata: expectWorkspaceMeta,
			ResourceStatus: schema.WorkspaceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Block storage
	block := suite.params.BlockStorage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	expectedBlockLabels := block.Labels
	expectedBlockAnnotations := block.Annotations
	expectedBlockExtensions := block.Extensions
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", t, suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Labels:         expectedBlockLabels,
			Annotations:    expectedBlockAnnotations,
			Extensions:     expectedBlockExtensions,
			Spec:           expectedBlockSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(block.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
		Name:      block.Metadata.Name,
	}
	stepsBuilder.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1, blockWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			Metadata: expectedBlockMeta,
			Spec:     expectedBlockSpec,
			ResourceStatus: schema.BlockStorageStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Constraint violations
	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with name exceeding maxLength:128 — expect rejection",
		suite.Client.StorageV1,
		suite.params.OverLengthNameImage,
	)
	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.StorageV1,
		suite.params.InvalidPatternNameImage,
	)
	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with label value exceeding maxLength:63 — expect rejection",
		suite.Client.StorageV1,
		suite.params.OverLengthLabelValueImage,
	)
	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.StorageV1,
		suite.params.OverLengthAnnotationImage,
	)
	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with invalid cpuArchitecture enum value — expect rejection",
		suite.Client.StorageV1,
		suite.params.InvalidCpuArchitectureImage,
	)
	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with invalid initializer enum value — expect rejection",
		suite.Client.StorageV1,
		suite.params.InvalidInitializerImage,
	)
	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with invalid boot enum value — expect rejection",
		suite.Client.StorageV1,
		suite.params.InvalidBootImage,
	)

	suite.FinishScenario()
}

func (suite *ImageConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
