package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureNetworkConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.NetworkConstraintsValidationV1Params) error {
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

	// Over-length name Validation
	overLengthNameNetwork := p.OverLengthNameNetwork
	overLengthNameURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, overLengthNameNetwork.Metadata.Tenant, overLengthNameNetwork.Metadata.Workspace, overLengthNameNetwork.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name Validation
	invalidPatternNameNetwork := p.InvalidPatternNameNetwork
	invalidPatternNameURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, invalidPatternNameNetwork.Metadata.Tenant, invalidPatternNameNetwork.Metadata.Workspace, invalidPatternNameNetwork.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value Validation
	overLengthLabelNetwork := p.OverLengthLabelValueNetwork
	overLengthLabelURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, overLengthLabelNetwork.Metadata.Tenant, overLengthLabelNetwork.Metadata.Workspace, overLengthLabelNetwork.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value Validation
	overLengthAnnotationNetwork := p.OverLengthAnnotationNetwork
	overLengthAnnotationURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, overLengthAnnotationNetwork.Metadata.Tenant, overLengthAnnotationNetwork.Metadata.Workspace, overLengthAnnotationNetwork.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length cidr ipv4 Validation
	overLengthCidrIpv4Network := p.OverLengthCidrIpv4Network
	overLengthCidrIpv4URL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, overLengthCidrIpv4Network.Metadata.Tenant, overLengthCidrIpv4Network.Metadata.Workspace, overLengthCidrIpv4Network.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthCidrIpv4URL, scenario.MockParams); err != nil {
		return err
	}

	// Under-length cidr ipv4 Validation
	underLengthCidrIpv4Network := p.UnderLengthCidrIpv4Network
	underLengthCidrIpv4URL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, underLengthCidrIpv4Network.Metadata.Tenant, underLengthCidrIpv4Network.Metadata.Workspace, underLengthCidrIpv4Network.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(underLengthCidrIpv4URL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length cidr ipv6 Validation
	overLengthCidrIpv6Network := p.OverLengthCidrIpv6Network
	overLengthCidrIpv6URL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, overLengthCidrIpv6Network.Metadata.Tenant, overLengthCidrIpv6Network.Metadata.Workspace, overLengthCidrIpv6Network.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthCidrIpv6URL, scenario.MockParams); err != nil {
		return err
	}

	// Under-length cidr ipv6 Validation
	underLengthCidrIpv6Network := p.UnderLengthCidrIpv6Network
	underLengthCidrIpv6URL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, underLengthCidrIpv6Network.Metadata.Tenant, underLengthCidrIpv6Network.Metadata.Workspace, underLengthCidrIpv6Network.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(underLengthCidrIpv6URL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length additional cidr ipv4 Validation
	overLengthAdditionalCidrIpv4Network := p.OverLengthAdditionalCidrIpv4Network
	overLengthAdditionalCidrIpv4URL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, overLengthAdditionalCidrIpv4Network.Metadata.Tenant, overLengthAdditionalCidrIpv4Network.Metadata.Workspace, overLengthAdditionalCidrIpv4Network.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAdditionalCidrIpv4URL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length additional cidr ipv6 Validation
	overLengthAdditionalCidrIpv6Network := p.OverLengthAdditionalCidrIpv6Network
	overLengthAdditionalCidrIpv6URL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, overLengthAdditionalCidrIpv6Network.Metadata.Tenant, overLengthAdditionalCidrIpv6Network.Metadata.Workspace, overLengthAdditionalCidrIpv6Network.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAdditionalCidrIpv6URL, scenario.MockParams); err != nil {
		return err
	}

	// Under-length additional cidr ipv4 Validation
	underLengthAdditionalCidrIpv4Network := p.UnderLengthAdditionalCidrIpv4Network
	underLengthAdditionalCidrIpv4URL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, underLengthAdditionalCidrIpv4Network.Metadata.Tenant, underLengthAdditionalCidrIpv4Network.Metadata.Workspace, underLengthAdditionalCidrIpv4Network.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(underLengthAdditionalCidrIpv4URL, scenario.MockParams); err != nil {
		return err
	}

	// Under-length additional cidr ipv6 Validation
	underLengthAdditionalCidrIpv6Network := p.UnderLengthAdditionalCidrIpv6Network
	underLengthAdditionalCidrIpv6URL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, underLengthAdditionalCidrIpv6Network.Metadata.Tenant, underLengthAdditionalCidrIpv6Network.Metadata.Workspace, underLengthAdditionalCidrIpv6Network.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(underLengthAdditionalCidrIpv6URL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length additional cidr ipv4 Validation
	overMaxItemsAdditionalCidrsNetwork := p.OverMaxItemsAdditionalCidrsNetwork
	overMaxItemsAdditionalCidrsURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, overMaxItemsAdditionalCidrsNetwork.Metadata.Tenant, overMaxItemsAdditionalCidrsNetwork.Metadata.Workspace, overMaxItemsAdditionalCidrsNetwork.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxItemsAdditionalCidrsURL, scenario.MockParams); err != nil {
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
