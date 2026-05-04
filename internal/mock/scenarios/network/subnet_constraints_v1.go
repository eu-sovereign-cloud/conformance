package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureSubnetConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.SubnetConstraintsValidationV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	network := p.Network
	internetGateway := p.InternetGateway
	routeTable := p.RouteTable

	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	networkURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)
	gatewayURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, internetGateway.Metadata.Tenant, internetGateway.Metadata.Workspace, internetGateway.Metadata.Name)
	routeTableURL := generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)

	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureGetCreatingNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length name Validation
	overLengthNameSubnet := p.OverLengthNameSubnet
	overLengthNameURL := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, overLengthNameSubnet.Metadata.Tenant, overLengthNameSubnet.Metadata.Workspace, overLengthNameSubnet.Metadata.Network, overLengthNameSubnet.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name Validation
	invalidPatternNameSubnet := p.InvalidPatternNameSubnet
	invalidPatternNameURL := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, invalidPatternNameSubnet.Metadata.Tenant, invalidPatternNameSubnet.Metadata.Workspace, invalidPatternNameSubnet.Metadata.Network, invalidPatternNameSubnet.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value Validation
	overLengthLabelSubnet := p.OverLengthLabelValueSubnet
	overLengthLabelURL := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, overLengthLabelSubnet.Metadata.Tenant, overLengthLabelSubnet.Metadata.Workspace, overLengthLabelSubnet.Metadata.Network, overLengthLabelSubnet.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value Validation
	overLengthAnnotationSubnet := p.OverLengthAnnotationSubnet
	overLengthAnnotationURL := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, overLengthAnnotationSubnet.Metadata.Tenant, overLengthAnnotationSubnet.Metadata.Workspace, overLengthAnnotationSubnet.Metadata.Network, overLengthAnnotationSubnet.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureDeleteStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureDeleteStub(routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureDeleteStub(networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(networkURL, scenario.MockParams); err != nil {
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
