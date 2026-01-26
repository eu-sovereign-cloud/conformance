package mockstorage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, mockParams *mock.MockParams, suiteParams *params.StorageLifeCycleV1Params) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	workspace := *suiteParams.Workspace
	blockStorageInitial := *suiteParams.BlockStorageInitial
	blockStorageUpdated := *suiteParams.BlockStorageUpdated
	imageInitial := *suiteParams.ImageInitial
	imageUpdated := *suiteParams.ImageUpdated

	configurator, err := stubs.NewStubConfigurator(scenario, mockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, blockStorageInitial.Metadata.Tenant, blockStorageInitial.Metadata.Workspace, blockStorageInitial.Metadata.Name)
	imageUrl := generators.GenerateImageURL(constants.StorageProviderV1, imageInitial.Metadata.Tenant, imageInitial.Metadata.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(workspace.Metadata.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(workspace.Metadata.Tenant).Region(workspace.Metadata.Region).
		Labels(workspace.Labels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(blockStorageInitial.Metadata.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(blockStorageInitial.Metadata.Tenant).Workspace(blockStorageInitial.Metadata.Workspace).Region(blockStorageInitial.Metadata.Region).
		Spec(&blockStorageInitial.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the block storage
	blockResponse.Spec = blockStorageUpdated.Spec
	if err := configurator.ConfigureUpdateBlockStorageStub(blockResponse, blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Image
	imageResponse, err := builders.NewImageBuilder().
		Name(imageInitial.Metadata.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(imageInitial.Metadata.Tenant).Region(imageInitial.Metadata.Region).
		Spec(&imageInitial.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an image
	if err := configurator.ConfigureCreateImageStub(imageResponse, imageUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created image
	if err := configurator.ConfigureGetActiveImageStub(imageResponse, imageUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the image
	imageResponse.Spec = imageUpdated.Spec
	if err := configurator.ConfigureUpdateImageStub(imageResponse, imageUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated image
	if err := configurator.ConfigureGetActiveImageStub(imageResponse, imageUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the image
	if err := configurator.ConfigureDeleteStub(imageUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted image
	if err := configurator.ConfigureGetNotFoundStub(imageUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
