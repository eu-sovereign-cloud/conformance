package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureSubnetErrorV1 sets up mock stubs for the subnet error scenarios suite.
// Creates a valid workspace + network + internet gateway + route table environment
// before testing error scenarios, all invalid subnet requests returning 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create subnet with invalid region
//   - Create subnet with invalid zone
//   - Create subnet with non-existent workspace
//   - Create subnet with non-existent network
//   - Create subnet with non-existent routeTableRef
//   - Create subnet with CIDR outside network CIDR
func ConfigureSubnetErrorV1(scenario *mockscenarios.Scenario, p params.SubnetErrorV1Params) error {
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

	// Create network
	network := p.Network
	networkURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)

	if err := configurator.ConfigureCreateNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Create internet gateway
	internetGateway := p.InternetGateway
	internetGatewayURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, internetGateway.Metadata.Tenant, internetGateway.Metadata.Workspace, internetGateway.Metadata.Name)

	if err := configurator.ConfigureCreateInternetGatewayStub(internetGateway, internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(internetGateway, internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(internetGateway, internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Create route table
	routeTable := p.RouteTable
	routeTableURL := generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)

	if err := configurator.ConfigureCreateRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid region violation
	invalidRegionURL := generators.GenerateSubnetURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidRegionSubnet.Metadata.Tenant,
		p.InvalidRegionSubnet.Metadata.Workspace,
		p.InvalidRegionSubnet.Metadata.Network,
		p.InvalidRegionSubnet.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid zone violation
	invalidZoneURL := generators.GenerateSubnetURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidZoneSubnet.Metadata.Tenant,
		p.InvalidZoneSubnet.Metadata.Workspace,
		p.InvalidZoneSubnet.Metadata.Network,
		p.InvalidZoneSubnet.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidZoneURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	nonExistentWorkspaceURL := generators.GenerateSubnetURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentWorkspaceSubnet.Metadata.Tenant,
		p.NonExistentWorkspaceSubnet.Metadata.Workspace,
		p.NonExistentWorkspaceSubnet.Metadata.Network,
		p.NonExistentWorkspaceSubnet.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentWorkspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent network violation
	nonExistentNetworkURL := generators.GenerateSubnetURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentNetworkSubnet.Metadata.Tenant,
		p.NonExistentNetworkSubnet.Metadata.Workspace,
		p.NonExistentNetworkSubnet.Metadata.Network,
		p.NonExistentNetworkSubnet.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentNetworkURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent routeTableRef violation
	nonExistentRouteTableRefURL := generators.GenerateSubnetURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentRouteTableRefSubnet.Metadata.Tenant,
		p.NonExistentRouteTableRefSubnet.Metadata.Workspace,
		p.NonExistentRouteTableRefSubnet.Metadata.Network,
		p.NonExistentRouteTableRefSubnet.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentRouteTableRefURL, scenario.MockParams); err != nil {
		return err
	}

	// CIDR outside network CIDR violation
	outsideCidrURL := generators.GenerateSubnetURL(
		sdkconsts.NetworkProviderV1Name,
		p.OutsideCidrSubnet.Metadata.Tenant,
		p.OutsideCidrSubnet.Metadata.Workspace,
		p.OutsideCidrSubnet.Metadata.Network,
		p.OutsideCidrSubnet.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(outsideCidrURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete route table teardown
	if err := configurator.ConfigureDeleteStub(routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete internet gateway teardown
	if err := configurator.ConfigureDeleteStub(internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingInternetGatewayStub(internetGateway, internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete network teardown
	if err := configurator.ConfigureDeleteStub(networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingNetworkStub(network, networkURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete workspace teardown
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
