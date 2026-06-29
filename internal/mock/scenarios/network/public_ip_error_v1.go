package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigurePublicIpErrorV1 sets up mock stubs for the public ip error scenarios suite.
// Creates a valid workspace environment before testing error scenarios,
// all invalid public ip requests returning 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create public ip with invalid region
//   - Create public ip with non-existent workspace
//   - Create public ip with invalid IP version
func ConfigurePublicIpErrorV1(scenario *mockscenarios.Scenario, p params.PublicIpErrorV1Params) error {
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
	invalidRegionURL := generators.GeneratePublicIpURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidRegionPublicIp.Metadata.Tenant,
		p.InvalidRegionPublicIp.Metadata.Workspace,
		p.InvalidRegionPublicIp.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	nonExistentWorkspaceURL := generators.GeneratePublicIpURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentWorkspacePublicIp.Metadata.Tenant,
		p.NonExistentWorkspacePublicIp.Metadata.Workspace,
		p.NonExistentWorkspacePublicIp.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentWorkspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid IP version violation
	invalidVersionURL := generators.GeneratePublicIpURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidVersionPublicIp.Metadata.Tenant,
		p.InvalidVersionPublicIp.Metadata.Workspace,
		p.InvalidVersionPublicIp.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidVersionURL, scenario.MockParams); err != nil {
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
