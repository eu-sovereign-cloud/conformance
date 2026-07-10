package mockstorage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureBlockStorageErrorV1 sets up mock stubs for the block storage
// error scenarios suite. Creates a valid workspace environment before testing
// error scenarios, all invalid block storage requests returning 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create block storage with invalid region
//   - Create block storage with invalid SKU
//   - Create block storage with non-existent workspace
func ConfigureBlockStorageErrorV1(scenario *mockscenarios.Scenario, p params.BlockStorageErrorV1Params) error {
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

	// Invalid region violation
	invalidRegionURL := generators.GenerateBlockStorageURL(
		sdkconsts.StorageProviderV1Name,
		p.InvalidRegionBlockStorage.Metadata.Tenant,
		p.InvalidRegionBlockStorage.Metadata.Workspace,
		p.InvalidRegionBlockStorage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid SKU violation
	invalidSkuURL := generators.GenerateBlockStorageURL(
		sdkconsts.StorageProviderV1Name,
		p.InvalidSkuBlockStorage.Metadata.Tenant,
		p.InvalidSkuBlockStorage.Metadata.Workspace,
		p.InvalidSkuBlockStorage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidSkuURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	// Note: uses a workspace name that was never created — URL points to "non-existent-workspace"
	nonExistentWorkspaceURL := generators.GenerateBlockStorageURL(
		sdkconsts.StorageProviderV1Name,
		p.NonExistentWorkspaceBlockStorage.Metadata.Tenant,
		p.NonExistentWorkspaceBlockStorage.Metadata.Workspace,
		p.NonExistentWorkspaceBlockStorage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentWorkspaceURL, scenario.MockParams); err != nil {
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

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
