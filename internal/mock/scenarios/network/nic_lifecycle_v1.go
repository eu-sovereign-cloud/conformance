package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureNicLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.NicLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := *params.Workspace
	network := *params.Network
	routeTable := *params.RouteTable
	subnet := *params.Subnet
	internetGateway := *params.InternetGateway
	nic := *params.NicInitial

	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	internetGatewayURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, internetGateway.Metadata.Tenant, internetGateway.Metadata.Workspace, internetGateway.Metadata.Name)
	routeTableURL := generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)
	subnetURL := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, subnet.Metadata.Tenant, subnet.Metadata.Workspace, subnet.Metadata.Network, subnet.Metadata.Name)
	nicURL := generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, nic.Metadata.Tenant, nic.Metadata.Workspace, nic.Metadata.Name)
	networkURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(workspace.Metadata.Name).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(workspace.Metadata.Tenant).Region(workspace.Metadata.Region).
		Labels(workspace.Labels).
		Build()
	if err != nil {
		return err
	}

	// Create a workspace
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
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(network.Metadata.Tenant).Workspace(network.Metadata.Workspace).Region(network.Metadata.Region).
		Spec(&network.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a network
	if err := configurator.ConfigureCreateNetworkStub(networkResponse, networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Internet Gateway
	gatewayInitialResponse, err := builders.NewInternetGatewayBuilder().
		Name(internetGateway.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(internetGateway.Metadata.Tenant).Workspace(internetGateway.Metadata.Workspace).Region(internetGateway.Metadata.Region).
		Spec(&internetGateway.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create internet gateway
	if err := configurator.ConfigureCreateInternetGatewayStub(gatewayInitialResponse, internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created internet gateway
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(gatewayInitialResponse, internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayInitialResponse, internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Route table
	routeResponse, err := builders.NewRouteTableBuilder().
		Name(routeTable.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(routeTable.Metadata.Tenant).Workspace(routeTable.Metadata.Workspace).Network(routeTable.Metadata.Network).Region(routeTable.Metadata.Region).
		Spec(&routeTable.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a route table
	if err := configurator.ConfigureCreateRouteTableStub(routeResponse, routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created route table
	if err := configurator.ConfigureGetCreatingRouteTableStub(routeResponse, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeResponse, routeTableURL, scenario.MockParams); err != nil {
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
		Name(subnet.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(subnet.Metadata.Tenant).Workspace(subnet.Metadata.Workspace).Network(subnet.Metadata.Network).Region(subnet.Metadata.Region).
		Spec(&subnet.Spec).
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

	// Nic
	nicResponse, err := builders.NewNicBuilder().
		Name(nic.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(nic.Metadata.Tenant).Workspace(nic.Metadata.Workspace).Region(nic.Metadata.Region).
		Spec(&nic.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a nic
	if err := configurator.ConfigureCreateNicStub(nicResponse, nicURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created nic
	if err := configurator.ConfigureGetCreatingNicStub(nicResponse, nicURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNicStub(nicResponse, nicURL, scenario.MockParams); err != nil {
		return err
	}

	// Update the nic
	nicResponse.Spec = params.NicUpdated.Spec
	if err := configurator.ConfigureUpdateNicStub(nicResponse, nicURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated nic
	if err := configurator.ConfigureGetUpdatingNicStub(nicResponse, nicURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNicStub(nicResponse, nicURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the internet gateway
	if err := configurator.ConfigureDeleteStub(internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the route table
	if err := configurator.ConfigureDeleteStub(routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the subnet
	if err := configurator.ConfigureDeleteStub(subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the nic
	if err := configurator.ConfigureDeleteStub(nicURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(nicURL, scenario.MockParams); err != nil {
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
