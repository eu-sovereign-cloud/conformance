package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureSubnetLifecycleScenarioV1(scenario *mockscenarios.Scenario, params params.SubnetLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := params.Workspace
	network := params.Network
	routeTable := params.RouteTable
	internetGateway := params.InternetGateway
	subnet := params.SubnetInitial

	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	networkURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)
	subnetURL := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, subnet.Metadata.Tenant, subnet.Metadata.Workspace, subnet.Metadata.Network, subnet.Metadata.Name)
	gatewayURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, internetGateway.Metadata.Tenant, internetGateway.Metadata.Workspace, internetGateway.Metadata.Name)
	routeUrl := generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)

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

	// Network

	// Create network
	if err := configurator.ConfigureCreateNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Internet Gateway

	// Create internet gateway
	if err := configurator.ConfigureCreateInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created internet gateway
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Route table

	// Create a route table
	if err := configurator.ConfigureCreateRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created route table
	if err := configurator.ConfigureGetCreatingRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created network
	if err := configurator.ConfigureGetCreatingNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Subnet

	// Create subnet
	if err := configurator.ConfigureCreateSubnetStub(subnet, subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created subnet
	if err := configurator.ConfigureGetCreatingSubnetStub(subnet, subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(subnet, subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Update subnet (change zone)
	subnet = params.SubnetUpdated
	if err := configurator.ConfigureUpdateSubnetStub(subnet, subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated subnet
	if err := configurator.ConfigureGetUpdatingSubnetStub(subnet, subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(subnet, subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete Subnet
	if err := configurator.ConfigureDeleteStub(subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingSubnetStub(subnet, subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the internet gateway
	if err := configurator.ConfigureDeleteStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the route table
	if err := configurator.ConfigureDeleteStub(routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the network
	if err := configurator.ConfigureDeleteStub(networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
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
