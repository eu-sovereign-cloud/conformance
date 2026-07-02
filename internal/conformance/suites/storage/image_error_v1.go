package storage

import (
	"math/rand"

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

// ImageErrorV1TestSuite verifies that Image resources with invalid references
// are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create image with invalid region
//   - Create image with invalid cpuArchitecture (not in [amd64, arm64])
//   - Create image with block storage in a different region
//   - Create image with non-existent workspace
type ImageErrorV1TestSuite struct {
	suites.RegionalTestSuite

	StorageSkus []string

	params *params.ImageErrorV1Params
}

func CreateImageErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite, storageSkus []string) *ImageErrorV1TestSuite {
	suite := &ImageErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       storageSkus,
	}
	suite.ScenarioName = constants.ImageErrorV1SuiteName.String()
	return suite
}

func (suite *ImageErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.StorageParentSuite)

	// Select sku
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]
	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)

	workspaceName := generators.GenerateWorkspaceName()
	blockStorageName := generators.GenerateBlockStorageName()
	storageSize := constants.BlockStorageInitialSize
	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	// Valid block storage ref (same region as suite)
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)

	// Cross-region block storage ref — block storage lives in a different region
	crossRegionBlockStorageName := generators.GenerateBlockStorageName()
	crossRegionBlockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, crossRegionBlockStorageName)

	// Build valid workspace
	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace for image error scenarios testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	// Build valid block storage (same region — used as real dependency)
	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "BlockStorage for image error scenarios testing"}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: storageSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	// Helper to build an image with custom overrides
	buildImage := func(name string, region string, blockStorageRef schema.Reference, cpuArch schema.ImageSpecCpuArchitecture) *schema.Image {
		image, err := builders.NewImageBuilder().
			Name(name).
			Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Region(region).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "Image for error scenario testing"}).
			Spec(&schema.ImageSpec{
				BlockStorageRef: blockStorageRef,
				CpuArchitecture: cpuArch,
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build Image: %v", err)
		}
		return image
	}

	// Invalid cpuArchitecture spec — value not in enum [amd64, arm64]
	invalidCpuArchSpec := &schema.ImageSpec{
		BlockStorageRef: *blockStorageRefObj,
		CpuArchitecture: "x86_64",
	}
	invalidCpuArchImage, err := builders.NewImageBuilder().
		Name(generators.GenerateImageName()).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Image with invalid cpuArchitecture"}).
		Spec(invalidCpuArchSpec).Build()
	if err != nil {
		t.Fatalf("Failed to build InvalidCpuArchitectureImage: %v", err)
	}

	p := &params.ImageErrorV1Params{
		Workspace:    workspace,
		BlockStorage: blockStorage,

		// Invalid region — completely random string, everything else valid
		InvalidRegionImage: buildImage(
			generators.GenerateImageName(),
			"invalid-region",
			*blockStorageRefObj,
			schema.ImageSpecCpuArchitectureAmd64,
		),

		// Invalid cpuArchitecture — "x86_64" not in enum [amd64, arm64]
		InvalidCpuArchitectureImage: invalidCpuArchImage,

		// Cross-region — image in suite.Region, blockStorageRef points to block storage in "other-region"
		// The crossRegionBlockStorage was never actually created in "other-region"
		CrossRegionBlockStorageImage: buildImage(
			generators.GenerateImageName(),
			suite.Region,
			*crossRegionBlockStorageRefObj,
			schema.ImageSpecCpuArchitectureAmd64,
		),

		// Non-existent workspace — image references a workspace that was never created
		// Note: Image is a tenant-scoped resource (no workspace in metadata),
		// but blockStorageRef points to "non-existent-workspace"
		NonExistentWorkspaceImage: buildImage(
			generators.GenerateImageName(),
			suite.Region,
			*generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, "non-existent-workspace", generators.GenerateBlockStorageName()),
			schema.ImageSpecCpuArchitectureAmd64,
		),
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockstorage.ConfigureImageErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ImageErrorV1TestSuite) TestScenario(t provider.T) {
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

	// Error scenarios — all must be rejected with 422

	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with invalid region — expect rejection",
		suite.Client.StorageV1,
		suite.params.InvalidRegionImage,
	)

	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with invalid cpuArchitecture (x86_64 not in [amd64, arm64]) — expect rejection",
		suite.Client.StorageV1,
		suite.params.InvalidCpuArchitectureImage,
	)

	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with block storage in a different region — expect rejection",
		suite.Client.StorageV1,
		suite.params.CrossRegionBlockStorageImage,
	)

	stepsBuilder.CreateOrUpdateImageExpectViolationV1Step(
		"Create an image with non-existent workspace in blockStorageRef — expect rejection",
		suite.Client.StorageV1,
		suite.params.NonExistentWorkspaceImage,
	)

	// Teardown

	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", t, suite.Client.StorageV1, block)
	stepsBuilder.WatchBlockStorageUntilDeletedV1Step("Watch the block storage deletion", t, suite.Client.StorageV1, blockWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *ImageErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
