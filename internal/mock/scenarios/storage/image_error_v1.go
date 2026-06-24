package mockstorage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureImageErrorV1 sets up mock stubs for the image error scenarios suite.
// Creates a valid workspace + block storage environment before testing error scenarios,
// all invalid image requests returning 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create image with invalid region
//   - Create image with invalid cpuArchitecture
//   - Create image with block storage in a different region
//   - Create image with non-existent workspace
func ConfigureImageErrorV1(scenario *mockscenarios.Scenario, p params.ImageErrorV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)

	// Create workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Create block storage
	blockStorage := p.BlockStorage
	blockURL := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockStorage.Metadata.Tenant, blockStorage.Metadata.Workspace, blockStorage.Metadata.Name)

	if err := configurator.ConfigureCreateBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid region violation
	// Note: image region = "invalid-region", everything else valid
	invalidRegionURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		p.InvalidRegionImage.Metadata.Tenant,
		p.InvalidRegionImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid cpuArchitecture violation
	// Note: cpuArchitecture = "x86_64" (not in enum [amd64, arm64])
	invalidCpuArchURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		p.InvalidCpuArchitectureImage.Metadata.Tenant,
		p.InvalidCpuArchitectureImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidCpuArchURL, scenario.MockParams); err != nil {
		return err
	}

	// Cross-region block storage violation
	// Note: image region = suite.Region, blockStorageRef points to a block storage in a different region
	crossRegionURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		p.CrossRegionBlockStorageImage.Metadata.Tenant,
		p.CrossRegionBlockStorageImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(crossRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	// Note: image references a workspace that was never created
	nonExistentWorkspaceURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		p.NonExistentWorkspaceImage.Metadata.Tenant,
		p.NonExistentWorkspaceImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentWorkspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete block storage teardown
	if err := configurator.ConfigureDeleteStub(blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(blockURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete workspace teardown
	if err := configurator.ConfigureDeleteStub(workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
