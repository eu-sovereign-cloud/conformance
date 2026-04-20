package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureSecurityGroupRuleConstraintsViolationsV1(scenario *mockscenarios.Scenario, p params.SecurityGroupRuleConstraintsViolationsV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)

	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	for _, url := range []string{
		generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, p.OverLengthNameSecurityGroupRule.Metadata.Tenant, p.OverLengthNameSecurityGroupRule.Metadata.Workspace, p.OverLengthNameSecurityGroupRule.Metadata.Name),
		generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, p.InvalidPatternNameSecurityGroupRule.Metadata.Tenant, p.InvalidPatternNameSecurityGroupRule.Metadata.Workspace, p.InvalidPatternNameSecurityGroupRule.Metadata.Name),
		generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, p.OverLengthLabelValueSecurityGroupRule.Metadata.Tenant, p.OverLengthLabelValueSecurityGroupRule.Metadata.Workspace, p.OverLengthLabelValueSecurityGroupRule.Metadata.Name),
		generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, p.OverLengthAnnotationSecurityGroupRule.Metadata.Tenant, p.OverLengthAnnotationSecurityGroupRule.Metadata.Workspace, p.OverLengthAnnotationSecurityGroupRule.Metadata.Name),
	} {
		if err := configurator.ConfigurePutUnprocessableEntityStub(url, scenario.MockParams); err != nil {
			return err
		}
	}

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
