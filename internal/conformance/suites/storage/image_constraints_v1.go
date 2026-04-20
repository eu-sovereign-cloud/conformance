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

// ImageConstraintsV1TestSuite verifies that Image resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 64 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type ImageConstraintsV1TestSuite struct {
	suites.RegionalTestSuite
	StorageSkus []string

	params *params.ImageConstraintsViolationsV1Params
}

func CreateImageConstraintsV1TestSuite(regionalTestSuite suites.RegionalTestSuite, storageSkus []string) *ImageConstraintsV1TestSuite {
	suite := &ImageConstraintsV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       storageSkus,
	}
	suite.ScenarioName = constants.ImageConstraintsV1SuiteName.String()
	return suite
}

func (suite *ImageConstraintsV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	// Select sku
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]

	// Generate scenario data
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

	// Image with name exceeding maxLength: 128 (129 chars)
	overLengthName := strings.Repeat("a", 129)
	overLengthNameImage, err := builders.NewImageBuilder().
		Name(overLengthName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Image with over-length name"}).
		Spec(&schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthNameImage: %v", err)
	}

	// Image with name violating kebab-case pattern (uppercase letters)
	invalidPatternName := "Invalid-Name-With-Uppercase"
	invalidPatternNameImage, err := builders.NewImageBuilder().
		Name(invalidPatternName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Image with non-kebab-case name"}).
		Spec(&schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build invalidPatternNameImage: %v", err)
	}

	// Image with label value exceeding maxLength: 64
	overLengthLabelImage, err := builders.NewImageBuilder().
		Name(generators.GenerateImageName()).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
			"constraint-test":  strings.Repeat("x", 64),
		}).
		Annotations(schema.Annotations{"description": "Image with over-length label value"}).
		Spec(&schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthLabelImage: %v", err)
	}

	// Image with annotation value exceeding maxLength: 1024 (1025 chars)
	overLengthAnnotationImage, err := builders.NewImageBuilder().
		Name(generators.GenerateImageName()).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{
			"description":     "Image with over-length annotation value",
			"long-annotation": strings.Repeat("y", 1025),
		}).
		Spec(&schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthAnnotationImage: %v", err)
	}

	p := &params.ImageConstraintsViolationsV1Params{
		Workspace:                 workspace,
		BlockStorage:              blockStorage,
		OverLengthNameImage:       overLengthNameImage,
		InvalidPatternNameImage:   invalidPatternNameImage,
		OverLengthLabelValueImage: overLengthLabelImage,
		OverLengthAnnotationImage: overLengthAnnotationImage,
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockstorage.ConfigureImageConstraintsViolationsV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ImageConstraintsV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.StorageProviderV1Name, string(schema.GlobalTenantResourceMetadataKindResourceKindImage))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	expectWorkspaceAnnotations := workspace.Annotations
	expectWorkspaceExtensions := workspace.Extensions
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Annotations:    expectWorkspaceAnnotations,
			Extensions:     expectWorkspaceExtensions,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created Workspace
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

	// Create a block storage
	block := suite.params.BlockStorage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	expectedBlockLabels := block.Labels
	expectedBlockAnnotations := block.Annotations
	expectedBlockExtensions := block.Extensions
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Labels:         expectedBlockLabels,
			Annotations:    expectedBlockAnnotations,
			Extensions:     expectedBlockExtensions,
			Spec:           expectedBlockSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created block storage
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

	// name: maxLength 128 — must be rejected
	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with name exceeding maxLength:128 — expect rejection",
		suite.Client.StorageV1,
		suite.params.OverLengthNameImage,
	)

	// name: pattern — must be rejected
	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.StorageV1,
		suite.params.InvalidPatternNameImage,
	)

	// labels value: maxLength 64 — must be rejected
	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with label value exceeding maxLength:64 — expect rejection",
		suite.Client.StorageV1,
		suite.params.OverLengthLabelValueImage,
	)

	// annotations value: maxLength 1024 — must be rejected
	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.StorageV1,
		suite.params.OverLengthAnnotationImage,
	)

	suite.FinishScenario()
}

func (suite *ImageConstraintsV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
