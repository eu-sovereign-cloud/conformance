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
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type CreateBlockStorageV1TestSuite struct {
	suites.RegionalTestSuite

	StorageSkus []string

	params *params.CreateBlockStorageV1Params
}

func NewCreateBlockStorageV1TestSuite(regionalTestSuite suites.RegionalTestSuite, storageSkus []string) *CreateBlockStorageV1TestSuite {
	suite := &CreateBlockStorageV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       storageSkus,
	}
	suite.ScenarioName = constants.CreateBlockStorageV1SuiteName.String()
	return suite
}

func (suite *CreateBlockStorageV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("storage")

	// Select sku
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj, err := generators.GenerateSkuRefObject(storageSkuName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	blockStorageName := generators.GenerateBlockStorageName()

	storageSize := generators.GenerateBlockStorageSize()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: storageSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	params := &params.CreateBlockStorageV1Params{
		Workspace:    workspace,
		BlockStorage: blockStorage,
	}

	suite.params = params

	err = suites.SetupMockIfEnabled(suite.TestSuite, mockstorage.ConfigureCreateBlockStorageScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *CreateBlockStorageV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.StorageProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels

	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Block storage

	// Create a block storage
	block := suite.params.BlockStorage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created block storage
	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(block.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
		Name:      block.Metadata.Name,
	}
	block = stepsBuilder.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1, blockWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion

	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", suite.Client.StorageV1, block)
	stepsBuilder.GetBlockStorageWithErrorV1Step("Get the deleted block storage", suite.Client.StorageV1, blockWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *CreateBlockStorageV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
