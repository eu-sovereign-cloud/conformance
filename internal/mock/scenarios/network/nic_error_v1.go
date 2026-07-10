package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureNicErrorV1 sets up mock stubs for the nic error scenarios suite.
// Creates a valid workspace + network + internet gateway + route table + subnet
// environment before testing error scenarios, all invalid nic requests returning
// 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create nic with invalid region
//   - Create nic with non-existent workspace
//   - Create nic with non-existent subnetRef
//   - Create nic with non-existent publicIpRef
func ConfigureNicErrorV1(scenario *mockscenarios.Scenario, p params.NicErrorV1Params) error {
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

	// Create subnet
	subnet := p.Subnet
	subnetURL := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, subnet.Metadata.Tenant, subnet.Metadata.Workspace, subnet.Metadata.Network, subnet.Metadata.Name)

	if err := configurator.ConfigureCreateSubnetStub(subnet, subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingSubnetStub(subnet, subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(subnet, subnetURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid region violation
	invalidRegionURL := generators.GenerateNicURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidRegionNic.Metadata.Tenant,
		p.InvalidRegionNic.Metadata.Workspace,
		p.InvalidRegionNic.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	nonExistentWorkspaceURL := generators.GenerateNicURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentWorkspaceNic.Metadata.Tenant,
		p.NonExistentWorkspaceNic.Metadata.Workspace,
		p.NonExistentWorkspaceNic.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentWorkspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent subnetRef violation
	nonExistentSubnetRefURL := generators.GenerateNicURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentSubnetRefNic.Metadata.Tenant,
		p.NonExistentSubnetRefNic.Metadata.Workspace,
		p.NonExistentSubnetRefNic.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentSubnetRefURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent publicIpRef violation
	nonExistentPublicIpRefURL := generators.GenerateNicURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentPublicIpRefNic.Metadata.Tenant,
		p.NonExistentPublicIpRefNic.Metadata.Workspace,
		p.NonExistentPublicIpRefNic.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentPublicIpRefURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete subnet teardown
	if err := configurator.ConfigureDeleteStub(subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingSubnetStub(subnet, subnetURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(subnetURL, scenario.MockParams); err != nil {
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
