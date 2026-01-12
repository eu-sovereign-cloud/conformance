package mockstorage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, params *params.StorageLifeCycleParamsV1) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	configurator, err := stubs.NewStubConfigurator(scenario, params.MockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, params.Workspace.Metadata.Tenant, params.Workspace.Metadata.Name)
	blockUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, params.BlockStorageInitial.Metadata.Tenant, params.BlockStorageInitial.Metadata.Workspace, params.BlockStorageInitial.Metadata.Name)
	imageUrl := generators.GenerateImageURL(constants.StorageProviderV1, params.ImageInitial.Metadata.Tenant, params.ImageInitial.Metadata.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Metadata.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Workspace.Metadata.Tenant).Region(params.Workspace.Metadata.Region).
		Labels(params.Workspace.Labels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorageInitial.Metadata.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.BlockStorageInitial.Metadata.Tenant).Workspace(params.BlockStorageInitial.Metadata.Workspace).Region(params.BlockStorageInitial.Metadata.Region).
		Spec(&params.BlockStorageInitial.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the block storage
	blockResponse.Spec = params.BlockStorageUpdated.Spec
	if err := configurator.ConfigureUpdateBlockStorageStub(blockResponse, blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Image
	imageResponse, err := builders.NewImageBuilder().
		Name(params.ImageInitial.Metadata.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.ImageInitial.Metadata.Tenant).Region(params.ImageInitial.Metadata.Region).
		Spec(&params.ImageInitial.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an image
	if err := configurator.ConfigureCreateImageStub(imageResponse, imageUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created image
	if err := configurator.ConfigureGetActiveImageStub(imageResponse, imageUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the image
	imageResponse.Spec = params.ImageUpdated.Spec
	if err := configurator.ConfigureUpdateImageStub(imageResponse, imageUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated image
	if err := configurator.ConfigureGetActiveImageStub(imageResponse, imageUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the image
	if err := configurator.ConfigureDeleteStub(imageUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted image
	if err := configurator.ConfigureGetNotFoundStub(imageUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
