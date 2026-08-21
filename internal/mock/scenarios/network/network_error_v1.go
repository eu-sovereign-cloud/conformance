package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureNetworkErrorV1 sets up mock stubs for the network error scenarios suite.
// Creates a valid workspace environment before testing error scenarios,
// all invalid network requests returning 422 Unprocessable Entity, and deleting a
// network that still has an active route table returns 409 Conflict.
//
// Scenarios tested:
//   - Create network with invalid region
//   - Create network with invalid SKU
//   - Create network with non-existent workspace
//   - Delete a network that still has an active route table → 409
func ConfigureNetworkErrorV1(scenario *mockscenarios.Scenario, p params.NetworkErrorV1Params) error {
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

	// Invalid region violation
	invalidRegionURL := generators.GenerateNetworkURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidRegionNetwork.Metadata.Tenant,
		p.InvalidRegionNetwork.Metadata.Workspace,
		p.InvalidRegionNetwork.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid SKU violation
	invalidSkuURL := generators.GenerateNetworkURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidSkuNetwork.Metadata.Tenant,
		p.InvalidSkuNetwork.Metadata.Workspace,
		p.InvalidSkuNetwork.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidSkuURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	nonExistentWorkspaceURL := generators.GenerateNetworkURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentWorkspaceNetwork.Metadata.Tenant,
		p.NonExistentWorkspaceNetwork.Metadata.Workspace,
		p.NonExistentWorkspaceNetwork.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentWorkspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete conflict scenario — network with an active route table
	network := p.Network
	internetGateway := p.InternetGateway
	routeTable := p.RouteTable

	networkURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)
	internetGatewayURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, internetGateway.Metadata.Tenant, internetGateway.Metadata.Workspace, internetGateway.Metadata.Name)
	routeTableURL := generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)

	// Create network
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
	if err := configurator.ConfigureCreateRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete network while route table is active → 409
	if err := configurator.ConfigureDeleteUnprocessableEntityStub(networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Teardown route table
	if err := configurator.ConfigureDeleteStub(routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingRouteTableStub(routeTable, routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Teardown internet gateway
	if err := configurator.ConfigureDeleteStub(internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingInternetGatewayStub(internetGateway, internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(internetGatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Teardown network
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
