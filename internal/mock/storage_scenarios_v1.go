package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/wiremock/go-wiremock"
)

func ConfigStorageLifecycleScenarioV1(scenario string, params *StorageParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageUrl := generators.GenerateImageURL(storageProviderV1, params.Tenant, params.Image.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.configureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.configureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.configureGetActiveBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Update the block storage
	blockResponse.Spec = *params.BlockStorage.UpdatedSpec
	if err := configurator.configureUpdateBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Get the updated block storage
	if err := configurator.configureGetActiveBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Image
	imageResponse, err := builders.NewImageBuilder().
		Name(params.Image.Name).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Spec(params.Image.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an image
	if err := configurator.configureCreateImageStub(imageResponse, imageUrl, params); err != nil {
		return nil, err
	}

	// Get the created image
	if err := configurator.configureGetActiveImageStub(imageResponse, imageUrl, params); err != nil {
		return nil, err
	}

	// Update the image
	imageResponse.Spec = *params.Image.UpdatedSpec
	if err := configurator.configureUpdateImageStub(imageResponse, imageUrl, params); err != nil {
		return nil, err
	}

	// Get the updated image
	if err := configurator.configureGetActiveImageStub(imageResponse, imageUrl, params); err != nil {
		return nil, err
	}

	// Delete the image
	if err := configurator.configureDeleteStub(imageUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted image
	if err := configurator.configureGetNotFoundStub(imageUrl, params); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.configureGetNotFoundStub(blockUrl, params); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(workspaceUrl, params); err != nil {
		return nil, err
	}

	return configurator.client, nil
}
