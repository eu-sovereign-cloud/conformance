package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureSecurityGroupRuleConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.SecurityGroupRuleConstraintsValidationV1Params) error {
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

	// Over-length name Validation
	overLengthNameSecurityGroupRule := p.OverLengthNameSecurityGroupRule
	overLengthNameURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, overLengthNameSecurityGroupRule.Metadata.Tenant, overLengthNameSecurityGroupRule.Metadata.Workspace, overLengthNameSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name Validation
	invalidPatternNameSecurityGroupRule := p.InvalidPatternNameSecurityGroupRule
	invalidPatternNameURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, invalidPatternNameSecurityGroupRule.Metadata.Tenant, invalidPatternNameSecurityGroupRule.Metadata.Workspace, invalidPatternNameSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value Validation
	overLengthLabelSecurityGroupRule := p.OverLengthLabelValueSecurityGroupRule
	overLengthLabelURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, overLengthLabelSecurityGroupRule.Metadata.Tenant, overLengthLabelSecurityGroupRule.Metadata.Workspace, overLengthLabelSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value Validation
	overLengthAnnotationSecurityGroupRule := p.OverLengthAnnotationSecurityGroupRule
	overLengthAnnotationURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, overLengthAnnotationSecurityGroupRule.Metadata.Tenant, overLengthAnnotationSecurityGroupRule.Metadata.Workspace, overLengthAnnotationSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
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
