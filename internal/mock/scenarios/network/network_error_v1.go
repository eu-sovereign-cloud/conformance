package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureNetworkErrorV1 sets up mock stubs for the network error scenarios suite.
// Creates a valid workspace environment before testing error scenarios,
// all invalid network requests returning 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create network with invalid region
//   - Create network with invalid SKU
//   - Create network with non-existent workspace
func ConfigureNetworkErrorV1(scenario *mockscenarios.Scenario, p params.NetworkErrorV1Params) error {
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
	invalidRegionURL := generators.GenerateNetworkURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidRegionNetwork.Metadata.Tenant,
		p.InvalidRegionNetwork.Metadata.Workspace,
		p.InvalidRegionNetwork.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid SKU violation
	invalidSkuURL := generators.GenerateNetworkURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidSkuNetwork.Metadata.Tenant,
		p.InvalidSkuNetwork.Metadata.Workspace,
		p.InvalidSkuNetwork.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidSkuURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	nonExistentWorkspaceURL := generators.GenerateNetworkURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentWorkspaceNetwork.Metadata.Tenant,
		p.NonExistentWorkspaceNetwork.Metadata.Workspace,
		p.NonExistentWorkspaceNetwork.Metadata.Name,
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

	return scenario.FinishConfiguration(configurator)
}
