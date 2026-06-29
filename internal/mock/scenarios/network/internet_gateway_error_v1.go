package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureInternetGatewayErrorV1 sets up mock stubs for the internet gateway
// error scenarios suite. Creates a valid workspace environment before testing
// error scenarios, all invalid internet gateway requests returning 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create internet gateway with invalid region
//   - Create internet gateway with non-existent workspace
func ConfigureInternetGatewayErrorV1(scenario *mockscenarios.Scenario, p params.InternetGatewayErrorV1Params) error {
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
	invalidRegionURL := generators.GenerateInternetGatewayURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidRegionInternetGateway.Metadata.Tenant,
		p.InvalidRegionInternetGateway.Metadata.Workspace,
		p.InvalidRegionInternetGateway.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	nonExistentWorkspaceURL := generators.GenerateInternetGatewayURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentWorkspaceInternetGateway.Metadata.Tenant,
		p.NonExistentWorkspaceInternetGateway.Metadata.Workspace,
		p.NonExistentWorkspaceInternetGateway.Metadata.Name,
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
