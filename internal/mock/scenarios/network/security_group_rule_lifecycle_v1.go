package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureSecurityGroupRuleLifecycleScenarioV1(scenario *mockscenarios.Scenario, params params.SecurityGroupRuleLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := params.Workspace
	securityGroupRule := params.SecurityGroupRuleInitial

	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	securityGroupRuleURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, securityGroupRule.Metadata.Tenant, securityGroupRule.Metadata.Workspace, securityGroupRule.Metadata.Name)

	// Workspace

	// Create workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Security group rule

	// Create a security group rule
	if err := configurator.ConfigureCreateSecurityGroupRuleStub(securityGroupRule, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created security group rule
	if err := configurator.ConfigureGetCreatingSecurityGroupRuleStub(securityGroupRule, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupRuleStub(securityGroupRule, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}

	// Update the security group rule
	securityGroupRule = params.SecurityGroupRuleUpdated
	if err := configurator.ConfigureUpdateSecurityGroupRuleStub(securityGroupRule, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated security group rule
	if err := configurator.ConfigureGetUpdatingSecurityGroupRuleStub(securityGroupRule, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupRuleStub(securityGroupRule, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the security group rule
	if err := configurator.ConfigureDeleteStub(securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceURL, scenario.MockParams); err != nil {
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
