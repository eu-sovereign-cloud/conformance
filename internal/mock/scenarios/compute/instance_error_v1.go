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
//   - Start an already-on instance → 409
//   - Stop an already-off instance → 409
//   - Restart an already-off instance → 409
func ConfigureInstanceErrorV1(scenario *mockscenarios.Scenario, p params.InstanceErrorV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	blockStorage := p.BlockStorage

	// Generate URLs
	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockURL := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockStorage.Metadata.Tenant, workspace.Metadata.Name, blockStorage.Metadata.Name)

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

	// Invalid SKU → 422
	invalidSkuURL := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, p.InvalidSkuInstance.Metadata.Tenant, workspace.Metadata.Name, p.InvalidSkuInstance.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidSkuURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace → 422
	nonExistentWorkspaceURL := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, p.NonExistentWorkspaceInstance.Metadata.Tenant, p.NonExistentWorkspaceInstance.Metadata.Workspace, p.NonExistentWorkspaceInstance.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentWorkspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent boot volume → 422
	nonExistentBootVolumeURL := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, p.NonExistentBootVolumeInstance.Metadata.Tenant, workspace.Metadata.Name, p.NonExistentBootVolumeInstance.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentBootVolumeURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid zone → 422
	invalidZoneURL := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, p.InvalidZoneInstance.Metadata.Tenant, workspace.Metadata.Name, p.InvalidZoneInstance.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidZoneURL, scenario.MockParams); err != nil {
		return err
	}

	// Power state error scenarios
	powerInstance := p.PowerInstance
	powerInstanceURL := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, powerInstance.Metadata.Tenant, workspace.Metadata.Name, powerInstance.Metadata.Name)
	powerInstanceStartURL := generators.GenerateInstanceStartURL(sdkconsts.ComputeProviderV1Name, powerInstance.Metadata.Tenant, workspace.Metadata.Name, powerInstance.Metadata.Name)
	powerInstanceStopURL := generators.GenerateInstanceStopURL(sdkconsts.ComputeProviderV1Name, powerInstance.Metadata.Tenant, workspace.Metadata.Name, powerInstance.Metadata.Name)
	powerInstanceRestartURL := generators.GenerateInstanceRestartURL(sdkconsts.ComputeProviderV1Name, powerInstance.Metadata.Tenant, workspace.Metadata.Name, powerInstance.Metadata.Name)

	// Create power instance → creating → active
	if err := configurator.ConfigureCreateInstanceStub(powerInstance, powerInstanceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingInstanceStub(powerInstance, powerInstanceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInstanceStub(powerInstance, powerInstanceURL, scenario.MockParams); err != nil {
		return err
	}

	// Start → instance transitions to powerState=on
	if err := configurator.ConfigureInstanceOperationStub(powerInstance, powerInstanceStartURL, scenario.MockParams); err != nil {
		return err
	}
	// GET after start → still active
	if err := configurator.ConfigureGetActiveInstanceStub(powerInstance, powerInstanceURL, scenario.MockParams); err != nil {
		return err
	}

	// Start on already-on instance → 412
	if err := configurator.ConfigurePostUnprocessableEntityStub(powerInstanceStartURL, scenario.MockParams); err != nil {
		return err
	}

	// Stop → instance transitions to powerState=off
	if err := configurator.ConfigureInstanceOperationStub(powerInstance, powerInstanceStopURL, scenario.MockParams); err != nil {
		return err
	}

	// Stop on already-off instance → 412
	if err := configurator.ConfigurePostUnprocessableEntityStub(powerInstanceStopURL, scenario.MockParams); err != nil {
		return err
	}

	// Restart on already-off instance → 412
	if err := configurator.ConfigurePostUnprocessableEntityStub(powerInstanceRestartURL, scenario.MockParams); err != nil {
		return err
	}

	// Teardown power instance
	if err := configurator.ConfigureDeleteStub(powerInstanceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingInstanceStub(powerInstance, powerInstanceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(powerInstanceURL, scenario.MockParams); err != nil {
		return err
	}

	// Teardown block storage
	if err := configurator.ConfigureDeleteStub(blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(blockURL, scenario.MockParams); err != nil {
		return err
	}

	// Teardown workspace
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
