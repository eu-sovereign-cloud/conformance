package network

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/conformance/pkg/mock/scenarios"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureSecurityGroupRuleConstraintsValidationV1(scenario *scenarios.Scenario, p params.SecurityGroupRuleConstraintsValidationV1Params) error {
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

	// Invalid direction Validation
	invalidDirectionSecurityGroupRule := p.InvalidDirectionSecurityGroupRule
	invalidDirectionURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, invalidDirectionSecurityGroupRule.Metadata.Tenant, invalidDirectionSecurityGroupRule.Metadata.Workspace, invalidDirectionSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidDirectionURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid version Validation
	invalidVersionSecurityGroupRule := p.InvalidVersionSecurityGroupRule
	invalidVersionURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, invalidVersionSecurityGroupRule.Metadata.Tenant, invalidVersionSecurityGroupRule.Metadata.Workspace, invalidVersionSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidVersionURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid protocol Validation
	invalidProtocolSecurityGroupRule := p.InvalidProtocolSecurityGroupRule
	invalidProtocolURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, invalidProtocolSecurityGroupRule.Metadata.Tenant, invalidProtocolSecurityGroupRule.Metadata.Workspace, invalidProtocolSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidProtocolURL, scenario.MockParams); err != nil {
		return err
	}

	// Over max port from Validation
	overMaxPortFromSecurityGroupRule := p.OverMaxPortFromSecurityGroupRule
	overMaxPortFromURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, overMaxPortFromSecurityGroupRule.Metadata.Tenant, overMaxPortFromSecurityGroupRule.Metadata.Workspace, overMaxPortFromSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxPortFromURL, scenario.MockParams); err != nil {
		return err
	}

	// Under min port from Validation
	underMinPortFromSecurityGroupRule := p.UnderMinPortFromSecurityGroupRule
	underMinPortFromURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, underMinPortFromSecurityGroupRule.Metadata.Tenant, underMinPortFromSecurityGroupRule.Metadata.Workspace, underMinPortFromSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(underMinPortFromURL, scenario.MockParams); err != nil {
		return err
	}

	// Over max port to Validation
	overMaxPortToSecurityGroupRule := p.OverMaxPortToSecurityGroupRule
	overMaxPortToURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, overMaxPortToSecurityGroupRule.Metadata.Tenant, overMaxPortToSecurityGroupRule.Metadata.Workspace, overMaxPortToSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxPortToURL, scenario.MockParams); err != nil {
		return err
	}

	// Under min port to Validation
	underMinPortToSecurityGroupRule := p.UnderMinPortToSecurityGroupRule
	underMinPortToURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, underMinPortToSecurityGroupRule.Metadata.Tenant, underMinPortToSecurityGroupRule.Metadata.Workspace, underMinPortToSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(underMinPortToURL, scenario.MockParams); err != nil {
		return err
	}

	// Over max port list Validation
	overMaxPortListSecurityGroupRule := p.OverMaxPortListSecurityGroupRule
	overMaxPortListURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, overMaxPortListSecurityGroupRule.Metadata.Tenant, overMaxPortListSecurityGroupRule.Metadata.Workspace, overMaxPortListSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxPortListURL, scenario.MockParams); err != nil {
		return err
	}

	// Under min port list Validation
	underMinPortListSecurityGroupRule := p.UnderMinPortListSecurityGroupRule
	underMinPortListURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, underMinPortListSecurityGroupRule.Metadata.Tenant, underMinPortListSecurityGroupRule.Metadata.Workspace, underMinPortListSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(underMinPortListURL, scenario.MockParams); err != nil {
		return err
	}

	// Over max icmp type Validation
	overMaxIcmpTypeSecurityGroupRule := p.OverMaxIcmpTypeSecurityGroupRule
	overMaxIcmpTypeURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, overMaxIcmpTypeSecurityGroupRule.Metadata.Tenant, overMaxIcmpTypeSecurityGroupRule.Metadata.Workspace, overMaxIcmpTypeSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxIcmpTypeURL, scenario.MockParams); err != nil {
		return err
	}

	// Over max icmp code Validation
	overMaxIcmpCodeSecurityGroupRule := p.OverMaxIcmpCodeSecurityGroupRule
	overMaxIcmpCodeURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, overMaxIcmpCodeSecurityGroupRule.Metadata.Tenant, overMaxIcmpCodeSecurityGroupRule.Metadata.Workspace, overMaxIcmpCodeSecurityGroupRule.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxIcmpCodeURL, scenario.MockParams); err != nil {
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
