package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureRouteTableConstraintsViolationsV1(scenario *mockscenarios.Scenario, p params.RouteTableConstraintsViolationsV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	network := p.Network
	internetGateway := p.InternetGateway

	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	networkURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)
	gatewayURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, internetGateway.Metadata.Tenant, internetGateway.Metadata.Workspace, internetGateway.Metadata.Name)

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

	if err := configurator.ConfigureGetCreatingNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}

	for _, url := range []string{
		generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, p.OverLengthNameRouteTable.Metadata.Tenant, p.OverLengthNameRouteTable.Metadata.Workspace, p.OverLengthNameRouteTable.Metadata.Network, p.OverLengthNameRouteTable.Metadata.Name),
		generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, p.InvalidPatternNameRouteTable.Metadata.Tenant, p.InvalidPatternNameRouteTable.Metadata.Workspace, p.InvalidPatternNameRouteTable.Metadata.Network, p.InvalidPatternNameRouteTable.Metadata.Name),
		generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, p.OverLengthLabelValueRouteTable.Metadata.Tenant, p.OverLengthLabelValueRouteTable.Metadata.Workspace, p.OverLengthLabelValueRouteTable.Metadata.Network, p.OverLengthLabelValueRouteTable.Metadata.Name),
		generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, p.OverLengthAnnotationRouteTable.Metadata.Tenant, p.OverLengthAnnotationRouteTable.Metadata.Workspace, p.OverLengthAnnotationRouteTable.Metadata.Network, p.OverLengthAnnotationRouteTable.Metadata.Name),
	} {
		if err := configurator.ConfigurePutUnprocessableEntityStub(url, scenario.MockParams); err != nil {
			return err
		}
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
