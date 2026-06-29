package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/pkg/mock/scenarios"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureSecurityGroupConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.SecurityGroupConstraintsValidationV1Params) error {
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
	overLengthNameSecurityGroup := p.OverLengthNameSecurityGroup
	overLengthNameURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, overLengthNameSecurityGroup.Metadata.Tenant, overLengthNameSecurityGroup.Metadata.Workspace, overLengthNameSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name Validation
	invalidPatternNameSecurityGroup := p.InvalidPatternNameSecurityGroup
	invalidPatternNameURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, invalidPatternNameSecurityGroup.Metadata.Tenant, invalidPatternNameSecurityGroup.Metadata.Workspace, invalidPatternNameSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value Validation
	overLengthLabelSecurityGroup := p.OverLengthLabelValueSecurityGroup
	overLengthLabelURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, overLengthLabelSecurityGroup.Metadata.Tenant, overLengthLabelSecurityGroup.Metadata.Workspace, overLengthLabelSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value Validation
	overLengthAnnotationSecurityGroup := p.OverLengthAnnotationSecurityGroup
	overLengthAnnotationURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, overLengthAnnotationSecurityGroup.Metadata.Tenant, overLengthAnnotationSecurityGroup.Metadata.Workspace, overLengthAnnotationSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid direction Validation
	invalidDirectionSecurityGroup := p.InvalidDirectionSecurityGroup
	invalidDirectionURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, invalidDirectionSecurityGroup.Metadata.Tenant, invalidDirectionSecurityGroup.Metadata.Workspace, invalidDirectionSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidDirectionURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid version Validation
	invalidVersionSecurityGroup := p.InvalidVersionSecurityGroup
	invalidVersionURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, invalidVersionSecurityGroup.Metadata.Tenant, invalidVersionSecurityGroup.Metadata.Workspace, invalidVersionSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidVersionURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid protocol Validation
	invalidProtocolSecurityGroup := p.InvalidProtocolSecurityGroup
	invalidProtocolURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, invalidProtocolSecurityGroup.Metadata.Tenant, invalidProtocolSecurityGroup.Metadata.Workspace, invalidProtocolSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidProtocolURL, scenario.MockParams); err != nil {
		return err
	}

	// Over max port from Validation
	overMaxPortFromSecurityGroup := p.OverMaxPortFromSecurityGroup
	overMaxPortFromURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, overMaxPortFromSecurityGroup.Metadata.Tenant, overMaxPortFromSecurityGroup.Metadata.Workspace, overMaxPortFromSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxPortFromURL, scenario.MockParams); err != nil {
		return err
	}

	// Under min port from Validation
	underMinPortFromSecurityGroup := p.UnderMinPortFromSecurityGroup
	underMinPortFromURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, underMinPortFromSecurityGroup.Metadata.Tenant, underMinPortFromSecurityGroup.Metadata.Workspace, underMinPortFromSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(underMinPortFromURL, scenario.MockParams); err != nil {
		return err
	}

	// Over max port to Validation
	overMaxPortToSecurityGroup := p.OverMaxPortToSecurityGroup
	overMaxPortToURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, overMaxPortToSecurityGroup.Metadata.Tenant, overMaxPortToSecurityGroup.Metadata.Workspace, overMaxPortToSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxPortToURL, scenario.MockParams); err != nil {
		return err
	}

	// Under min port to Validation
	underMinPortToSecurityGroup := p.UnderMinPortToSecurityGroup
	underMinPortToURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, underMinPortToSecurityGroup.Metadata.Tenant, underMinPortToSecurityGroup.Metadata.Workspace, underMinPortToSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(underMinPortToURL, scenario.MockParams); err != nil {
		return err
	}

	// Over max port list Validation
	overMaxPortListSecurityGroup := p.OverMaxPortListSecurityGroup
	overMaxPortListURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, overMaxPortListSecurityGroup.Metadata.Tenant, overMaxPortListSecurityGroup.Metadata.Workspace, overMaxPortListSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxPortListURL, scenario.MockParams); err != nil {
		return err
	}

	// Under min port list Validation
	underMinPortListSecurityGroup := p.UnderMinPortListSecurityGroup
	underMinPortListURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, underMinPortListSecurityGroup.Metadata.Tenant, underMinPortListSecurityGroup.Metadata.Workspace, underMinPortListSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(underMinPortListURL, scenario.MockParams); err != nil {
		return err
	}

	// Over max icmp type Validation
	overMaxIcmpTypeSecurityGroup := p.OverMaxIcmpTypeSecurityGroup
	overMaxIcmpTypeURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, overMaxIcmpTypeSecurityGroup.Metadata.Tenant, overMaxIcmpTypeSecurityGroup.Metadata.Workspace, overMaxIcmpTypeSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxIcmpTypeURL, scenario.MockParams); err != nil {
		return err
	}

	// Over max icmp code Validation
	overMaxIcmpCodeSecurityGroup := p.OverMaxIcmpCodeSecurityGroup
	overMaxIcmpCodeURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, overMaxIcmpCodeSecurityGroup.Metadata.Tenant, overMaxIcmpCodeSecurityGroup.Metadata.Workspace, overMaxIcmpCodeSecurityGroup.Metadata.Name)
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
