package storage

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockstorage "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/storage"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type LifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	StorageSkus []string
}

func (suite *LifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.StorageProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	var err error

	// Select sku
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj, err := generators.GenerateSkuRefObject(storageSkuName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj, err := generators.GenerateBlockStorageRefObject(blockStorageName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	imageName := generators.GenerateImageName()

	initialStorageSize := generators.GenerateBlockStorageSize()
	updatedStorageSize := generators.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.MockEnabled {
		mockParams := &params.StorageLifeCycleParamsV1{
			BaseParams: &params.BaseParams{
				Tenant: suite.Tenant,
				Region: suite.Region,
				MockParams: &mock.MockParams{
					ServerURL: *suite.MockServerURL,
					AuthToken: suite.AuthToken,
				},
			},
			Workspace: &params.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					constants.EnvLabel: constants.EnvDevelopmentLabel,
				},
			},
			BlockStorage: &params.ResourceParams[schema.BlockStorageSpec]{
				Name: blockStorageName,
				InitialSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: initialStorageSize,
				},
				UpdatedSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: updatedStorageSize,
				},
			},
			Image: &params.ResourceParams[schema.ImageSpec]{
				Name: imageName,
				InitialSpec: &schema.ImageSpec{
					BlockStorageRef: *blockStorageRefObj,
					CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
				},
				UpdatedSpec: &schema.ImageSpec{
					BlockStorageRef: *blockStorageRefObj,
					CpuArchitecture: schema.ImageSpecCpuArchitectureArm64,
				},
			},
		}
		wm, err := mockstorage.ConfigureLifecycleScenarioV1(suite.ScenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.MockClient = wm
	}

	stepsBuilder := steps.NewStepsConfigurator(&suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.Tenant,
			Name:   workspaceName,
		},
	}
	expectWorkspaceMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{constants.EnvLabel: constants.EnvDevelopmentLabel}

	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   workspaceName,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, *workspaceTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Block storage

	// Create a block storage
	block := &schema.BlockStorage{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      blockStorageName,
		},
		Spec: schema.BlockStorageSpec{
			SizeGB: initialStorageSize,
			SkuRef: *storageSkuRefObj,
		},
	}
	expectedBlockMeta, err := builders.NewBlockStorageMetadataBuilder().
		Name(blockStorageName).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectedBlockSpec := &schema.BlockStorageSpec{
		SizeGB: initialStorageSize,
		SkuRef: *storageSkuRefObj,
	}
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created block storage
	blockWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	block = stepsBuilder.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1, *blockWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the block storage
	block.Spec.SizeGB = updatedStorageSize
	expectedBlockSpec.SizeGB = block.Spec.SizeGB
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Update the block storage", suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated block storage
	block = stepsBuilder.GetBlockStorageV1Step("Get the updated block storage", suite.Client.StorageV1, *blockWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Image

	// Create an image
	image := &schema.Image{
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.Tenant,
			Name:   imageName,
		},
		Spec: schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		},
	}
	expectedImageMeta, err := builders.NewImageMetadataBuilder().
		Name(imageName).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectedImageSpec := &schema.ImageSpec{
		BlockStorageRef: *blockStorageRefObj,
		CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
	}
	stepsBuilder.CreateOrUpdateImageV1Step("Create an image", suite.Client.StorageV1, image,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:      expectedImageMeta,
			Spec:          expectedImageSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created image
	imageTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   imageName,
	}
	image = stepsBuilder.GetImageV1Step("Get the created image", suite.Client.StorageV1, *imageTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:      expectedImageMeta,
			Spec:          expectedImageSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the image
	image.Spec.CpuArchitecture = schema.ImageSpecCpuArchitectureArm64
	expectedImageSpec.CpuArchitecture = image.Spec.CpuArchitecture
	stepsBuilder.CreateOrUpdateImageV1Step("Update the image", suite.Client.StorageV1, image,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:      expectedImageMeta,
			Spec:          expectedImageSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated image
	image = stepsBuilder.GetImageV1Step("Get the updated image", suite.Client.StorageV1, *imageTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:      expectedImageMeta,
			Spec:          expectedImageSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion

	stepsBuilder.DeleteImageV1Step("Delete the image", suite.Client.StorageV1, image)
	stepsBuilder.GetImageWithErrorV1Step("Get the deleted image", suite.Client.StorageV1, *imageTRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", suite.Client.StorageV1, block)
	stepsBuilder.GetBlockStorageWithErrorV1Step("Get the deleted block storage", suite.Client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *LifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
