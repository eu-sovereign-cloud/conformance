package mockstorage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
)

func ConfigureProviderLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.StorageProviderLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := *params.Workspace
	blockStorageInitial := *params.BlockStorageInitial
	blockStorageUpdated := *params.BlockStorageUpdated
	imageInitial := *params.ImageInitial
	imageUpdated := *params.ImageUpdated

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
		return err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(blockStorageInitial.Metadata.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(blockStorageInitial.Metadata.Tenant).Workspace(blockStorageInitial.Metadata.Workspace).Region(blockStorageInitial.Metadata.Region).
		Spec(&blockStorageInitial.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the block storage
	blockResponse.Spec = blockStorageUpdated.Spec
	if err := configurator.ConfigureUpdateBlockStorageStub(blockResponse, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Image
	imageResponse, err := builders.NewImageBuilder().
		Name(imageInitial.Metadata.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(imageInitial.Metadata.Tenant).Region(imageInitial.Metadata.Region).
		Spec(&imageInitial.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create an image
	if err := configurator.ConfigureCreateImageStub(imageResponse, imageUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created image
	if err := configurator.ConfigureGetActiveImageStub(imageResponse, imageUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the image
	imageResponse.Spec = imageUpdated.Spec
	if err := configurator.ConfigureUpdateImageStub(imageResponse, imageUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated image
	if err := configurator.ConfigureGetActiveImageStub(imageResponse, imageUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the image
	if err := configurator.ConfigureDeleteStub(imageUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted image
	if err := configurator.ConfigureGetNotFoundStub(imageUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
