package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
)

func ConfigureSubnetLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.SubnetLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := *params.Workspace
	network := *params.Network
	routeTable := *params.RouteTable
	internetGateway := *params.InternetGateway
	subnetInitial := *params.SubnetInitial
	subnetUpdated := *params.SubnetUpdated

	workspaceURL := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	networkURL := generators.GenerateNetworkURL(constants.NetworkProviderV1, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)
	subnetURL := generators.GenerateSubnetURL(constants.NetworkProviderV1, subnetInitial.Metadata.Tenant, subnetInitial.Metadata.Workspace, subnetInitial.Metadata.Network, subnetInitial.Metadata.Name)
	gatewayURL := generators.GenerateInternetGatewayURL(constants.NetworkProviderV1, internetGateway.Metadata.Tenant, internetGateway.Metadata.Workspace, internetGateway.Metadata.Name)
	routeUrl := generators.GenerateRouteTableURL(constants.NetworkProviderV1, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(workspace.Metadata.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(workspace.Metadata.Tenant).Region(workspace.Metadata.Region).
		Labels(workspace.Labels).
		Build()
	if err != nil {
		return err
	}

	// Create workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspaceResponse, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Network
	networkResponse, err := builders.NewNetworkBuilder().
		Name(network.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(network.Metadata.Tenant).Workspace(network.Metadata.Workspace).Region(network.Metadata.Region).
		Spec(&network.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create network
	if err := configurator.ConfigureCreateNetworkStub(networkResponse, networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Internet Gateway
	gatewayInitialResponse, err := builders.NewInternetGatewayBuilder().
		Name(internetGateway.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(internetGateway.Metadata.Tenant).Workspace(internetGateway.Metadata.Workspace).Region(internetGateway.Metadata.Region).
		Spec(&internetGateway.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create internet gateway
	if err := configurator.ConfigureCreateInternetGatewayStub(gatewayInitialResponse, gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created internet gateway
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(gatewayInitialResponse, gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayInitialResponse, gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Route table
	routeResponse, err := builders.NewRouteTableBuilder().
		Name(routeTable.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(routeTable.Metadata.Tenant).Workspace(routeTable.Metadata.Workspace).Network(routeTable.Metadata.Network).Region(routeTable.Metadata.Region).
		Spec(&routeTable.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a route table
	if err := configurator.ConfigureCreateRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created route table
	if err := configurator.ConfigureGetCreatingRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created network
	if err := configurator.ConfigureGetCreatingNetworkStub(networkResponse, networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(networkResponse, networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Subnet
	subnetInitialResponse, err := builders.NewSubnetBuilder().
		Name(subnetInitial.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(subnetInitial.Metadata.Tenant).Workspace(subnetInitial.Metadata.Workspace).Network(subnetInitial.Metadata.Network).Region(subnetInitial.Metadata.Region).
		Spec(&subnetInitial.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create subnet
	if err := configurator.ConfigureCreateSubnetStub(subnetInitialResponse, subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created subnet
	if err := configurator.ConfigureGetCreatingSubnetStub(subnetInitialResponse, subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(subnetInitialResponse, subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Update subnet (change zone)
	subnetUpdatedResponse, err := builders.NewSubnetBuilder().
		Name(subnetUpdated.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(subnetUpdated.Metadata.Tenant).Workspace(subnetUpdated.Metadata.Workspace).Network(subnetUpdated.Metadata.Network).Region(subnetUpdated.Metadata.Region).
		Spec(&subnetUpdated.Spec).
		Build()
	if err != nil {
		return err
	}

	if err := configurator.ConfigureUpdateSubnetStub(subnetUpdatedResponse, subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated subnet
	if err := configurator.ConfigureGetUpdatingSubnetStub(subnetUpdatedResponse, subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(subnetUpdatedResponse, subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete Subnet
	if err := configurator.ConfigureDeleteStub(subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the internet gateway
	if err := configurator.ConfigureDeleteStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the route table
	if err := configurator.ConfigureDeleteStub(routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the network
	if err := configurator.ConfigureDeleteStub(networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(networkURL, scenario.MockParams); err != nil {
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
