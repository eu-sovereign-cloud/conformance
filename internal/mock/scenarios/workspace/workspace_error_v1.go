package mockworkspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureWorkspaceErrorV1 sets up mock stubs for the workspace error scenarios suite.
// Invalid workspace requests return 422 Unprocessable Entity, and deleting a workspace
// that still has an active instance returns 409 Conflict.
//
// Scenarios tested:
//   - Create workspace with non-existent region
//   - Delete a workspace that still has an active instance → 409
func ConfigureWorkspaceErrorV1(scenario *mockscenarios.Scenario, p params.WorkspaceErrorV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	// Non-existent region workspace — expect 412
	nonExistentRegionURL := generators.GenerateWorkspaceURL(
		sdkconsts.WorkspaceProviderV1Name,
		p.NonExistentRegionWorkspace.Metadata.Tenant,
		p.NonExistentRegionWorkspace.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete conflict scenario — workspace with an active instance
	workspace := p.Workspace
	blockStorage := p.BlockStorage
	instance := p.Instance

	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockURL := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockStorage.Metadata.Tenant, workspace.Metadata.Name, blockStorage.Metadata.Name)
	instanceURL := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)

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

	// Create instance
	if err := configurator.ConfigureCreateInstanceStub(instance, instanceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingInstanceStub(instance, instanceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInstanceStub(instance, instanceURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete workspace while instance is active → 409
	if err := configurator.ConfigureDeleteUnprocessableEntityStub(workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Teardown instance
	if err := configurator.ConfigureDeleteStub(instanceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingInstanceStub(instance, instanceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(instanceURL, scenario.MockParams); err != nil {
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
