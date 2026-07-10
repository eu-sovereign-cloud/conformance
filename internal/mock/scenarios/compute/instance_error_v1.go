package mockcompute

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureInstanceErrorV1 sets up mock stubs for the instance error scenarios suite.
// Creates a valid workspace + block storage environment before testing
// error scenarios, all invalid instance requests returning 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create instance with non-existent SKU
//   - Create instance with non-existent workspace
//   - Create instance with non-existent boot volume ref
//   - Create instance with invalid zone
func ConfigureInstanceErrorV1(scenario *mockscenarios.Scenario, p params.InstanceErrorV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	blockStorage := p.BlockStorage

	// Generate URLs
	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockURL := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockStorage.Metadata.Tenant, blockStorage.Metadata.Workspace, blockStorage.Metadata.Name)

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
	if err := configurator.ConfigureCreateBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid SKU violation
	invalidSkuURL := generators.GenerateInstanceURL(
		sdkconsts.ComputeProviderV1Name,
		p.InvalidSkuInstance.Metadata.Tenant,
		p.InvalidSkuInstance.Metadata.Workspace,
		p.InvalidSkuInstance.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidSkuURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	nonExistentWorkspaceURL := generators.GenerateInstanceURL(
		sdkconsts.ComputeProviderV1Name,
		p.NonExistentWorkspaceInstance.Metadata.Tenant,
		p.NonExistentWorkspaceInstance.Metadata.Workspace,
		p.NonExistentWorkspaceInstance.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentWorkspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent boot volume violation
	nonExistentBootVolumeURL := generators.GenerateInstanceURL(
		sdkconsts.ComputeProviderV1Name,
		p.NonExistentBootVolumeInstance.Metadata.Tenant,
		p.NonExistentBootVolumeInstance.Metadata.Workspace,
		p.NonExistentBootVolumeInstance.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentBootVolumeURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid zone violation
	invalidZoneURL := generators.GenerateInstanceURL(
		sdkconsts.ComputeProviderV1Name,
		p.InvalidZoneInstance.Metadata.Tenant,
		p.InvalidZoneInstance.Metadata.Workspace,
		p.InvalidZoneInstance.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidZoneURL, scenario.MockParams); err != nil {
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
