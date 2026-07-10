package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureSecurityGroupRuleErrorV1 sets up mock stubs for the security group rule
// error scenarios suite. Creates a valid workspace environment before testing
// error scenarios, all invalid security group rule requests returning 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create security group rule with invalid region
//   - Create security group rule with non-existent workspace
func ConfigureSecurityGroupRuleErrorV1(scenario *mockscenarios.Scenario, p params.SecurityGroupRuleErrorV1Params) error {
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
	invalidRegionURL := generators.GenerateSecurityGroupRuleURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidRegionSecurityGroupRule.Metadata.Tenant,
		p.InvalidRegionSecurityGroupRule.Metadata.Workspace,
		p.InvalidRegionSecurityGroupRule.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	nonExistentWorkspaceURL := generators.GenerateSecurityGroupRuleURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentWorkspaceSecurityGroupRule.Metadata.Tenant,
		p.NonExistentWorkspaceSecurityGroupRule.Metadata.Workspace,
		p.NonExistentWorkspaceSecurityGroupRule.Metadata.Name,
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
