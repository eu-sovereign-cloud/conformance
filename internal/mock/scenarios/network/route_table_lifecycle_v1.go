package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureRouteTableLifecycleScenarioV1(scenario *mockscenarios.Scenario, params params.RouteTableLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := params.Workspace
	network := params.Network
	internetGateway := params.InternetGateway
	routeTable := params.RouteTableInitial

	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	networkURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)
	routeTableURL := generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)
	gatewayURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, internetGateway.Metadata.Tenant, internetGateway.Metadata.Workspace, internetGateway.Metadata.Name)

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

	// Create route table
	if err := configurator.ConfigureCreateRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created route table
	if err := configurator.ConfigureGetCreatingRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Update route table
	routeTable = params.RouteTableUpdated
	if err := configurator.ConfigureUpdateRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated route table
	if err := configurator.ConfigureGetUpdatingRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created network
	if err := configurator.ConfigureGetCreatingNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(network, networkURL, scenario.MockParams); err != nil {
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
	if err := configurator.ConfigureDeleteStub(routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(routeTableURL, scenario.MockParams); err != nil {
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
