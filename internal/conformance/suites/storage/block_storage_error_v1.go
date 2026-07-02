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

// BlockStorageErrorV1TestSuite verifies that BlockStorage resources with
// invalid references are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create block storage with invalid region
//   - Create block storage with invalid SKU
//   - Create block storage with non-existent workspace
type BlockStorageErrorV1TestSuite struct {
	suites.RegionalTestSuite

	StorageSkus []string

	params *params.BlockStorageErrorV1Params
}

func CreateBlockStorageErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite, storageSkus []string) *BlockStorageErrorV1TestSuite {
	suite := &BlockStorageErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       storageSkus,
	}
	suite.ScenarioName = constants.BlockStorageErrorV1SuiteName.String()
	return suite
}

func (suite *BlockStorageErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.StorageParentSuite)

	// Select sku
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]
	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)

	workspaceName := generators.GenerateWorkspaceName()
	storageSize := constants.BlockStorageInitialSize
	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	// Build valid workspace (used for setup/teardown)
	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace for block storage error scenarios testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	buildBlockStorage := func(name string, workspaceRef string, region string) *schema.BlockStorage {
		bs, err := builders.NewBlockStorageBuilder().
			Name(name).
			Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceRef).Region(region).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "BlockStorage for error scenario testing"}).
			Spec(&schema.BlockStorageSpec{
				SkuRef: *storageSkuRefObj,
				SizeGB: storageSize,
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build BlockStorage: %v", err)
		}
		return bs
	}

	// Invalid SKU — valid workspace + valid region, SKU does not exist
	invalidSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, "non-existent-sku")
	invalidSkuBlockStorage, err := builders.NewBlockStorageBuilder().
		Name(generators.GenerateBlockStorageName()).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "BlockStorage with non-existent SKU"}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *invalidSkuRefObj,
			SizeGB: storageSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build InvalidSkuBlockStorage: %v", err)
	}

	p := &params.BlockStorageErrorV1Params{
		Workspace: workspace,

		// Invalid region — completely random string, valid workspace + valid SKU
		InvalidRegionBlockStorage: buildBlockStorage(
			generators.GenerateBlockStorageName(),
			workspaceName,
			"invalid-region",
		),

		// Invalid SKU — valid workspace + valid region, SKU does not exist
		InvalidSkuBlockStorage: invalidSkuBlockStorage,

		// Non-existent workspace — workspace was never created
		NonExistentWorkspaceBlockStorage: buildBlockStorage(
			generators.GenerateBlockStorageName(),
			"non-existent-workspace",
			suite.Region,
		),
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockstorage.ConfigureBlockStorageErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *BlockStorageErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.StorageProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage))
	suite.ConfigureDepends(t, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	expectWorkspaceAnnotations := workspace.Annotations
	expectWorkspaceExtension := workspace.Extensions
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Annotations:    expectWorkspaceAnnotations,
			Extensions:     expectWorkspaceExtension,
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created workspace
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

	// invalid region — must be rejected
	stepsBuilder.CreateOrUpdateBlockStorageExpectViolationV1Step(
		"Create a block storage with invalid region — expect rejection",
		suite.Client.StorageV1,
		suite.params.InvalidRegionBlockStorage,
	)

	// invalid SKU — must be rejected
	stepsBuilder.CreateOrUpdateBlockStorageExpectViolationV1Step(
		"Create a block storage with invalid SKU — expect rejection",
		suite.Client.StorageV1,
		suite.params.InvalidSkuBlockStorage,
	)

	// non-existent workspace — must be rejected
	stepsBuilder.CreateOrUpdateBlockStorageExpectViolationV1Step(
		"Create a block storage with non-existent workspace — expect rejection",
		suite.Client.StorageV1,
		suite.params.NonExistentWorkspaceBlockStorage,
	)

	// Teardown workspace
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *BlockStorageErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
