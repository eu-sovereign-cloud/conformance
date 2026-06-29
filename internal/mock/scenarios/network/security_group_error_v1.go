package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureSecurityGroupErrorV1 sets up mock stubs for the security group
// error scenarios suite. Creates a valid workspace environment before testing
// error scenarios, all invalid security group requests returning 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create security group with invalid region
//   - Create security group with non-existent workspace
//   - Create security group with non-existent rule ref
func ConfigureSecurityGroupErrorV1(scenario *mockscenarios.Scenario, p params.SecurityGroupErrorV1Params) error {
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
	invalidRegionURL := generators.GenerateSecurityGroupURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidRegionSecurityGroup.Metadata.Tenant,
		p.InvalidRegionSecurityGroup.Metadata.Workspace,
		p.InvalidRegionSecurityGroup.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	nonExistentWorkspaceURL := generators.GenerateSecurityGroupURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentWorkspaceSecurityGroup.Metadata.Tenant,
		p.NonExistentWorkspaceSecurityGroup.Metadata.Workspace,
		p.NonExistentWorkspaceSecurityGroup.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentWorkspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent rule ref violation
	nonExistentRuleRefURL := generators.GenerateSecurityGroupURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentRuleRefSecurityGroup.Metadata.Tenant,
		p.NonExistentRuleRefSecurityGroup.Metadata.Workspace,
		p.NonExistentRuleRefSecurityGroup.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentRuleRefURL, scenario.MockParams); err != nil {
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
