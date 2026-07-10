package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureRouteTableErrorV1 sets up mock stubs for the route table error scenarios suite.
// Creates a valid workspace + network + internet gateway environment before testing
// error scenarios, all invalid route table requests returning 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create route table with invalid region
//   - Create route table with non-existent workspace
//   - Create route table with non-existent network
//   - Create route table with non-existent targetRef
func ConfigureRouteTableErrorV1(scenario *mockscenarios.Scenario, p params.RouteTableErrorV1Params) error {
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

	// Invalid region violation
	invalidRegionURL := generators.GenerateRouteTableURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidRegionRouteTable.Metadata.Tenant,
		p.InvalidRegionRouteTable.Metadata.Workspace,
		p.InvalidRegionRouteTable.Metadata.Network,
		p.InvalidRegionRouteTable.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace violation
	nonExistentWorkspaceURL := generators.GenerateRouteTableURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentWorkspaceRouteTable.Metadata.Tenant,
		p.NonExistentWorkspaceRouteTable.Metadata.Workspace,
		p.NonExistentWorkspaceRouteTable.Metadata.Network,
		p.NonExistentWorkspaceRouteTable.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentWorkspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent network violation
	nonExistentNetworkURL := generators.GenerateRouteTableURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentNetworkRouteTable.Metadata.Tenant,
		p.NonExistentNetworkRouteTable.Metadata.Workspace,
		p.NonExistentNetworkRouteTable.Metadata.Network,
		p.NonExistentNetworkRouteTable.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentNetworkURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent targetRef violation
	nonExistentTargetRefURL := generators.GenerateRouteTableURL(
		sdkconsts.NetworkProviderV1Name,
		p.NonExistentTargetRefRouteTable.Metadata.Tenant,
		p.NonExistentTargetRefRouteTable.Metadata.Workspace,
		p.NonExistentTargetRefRouteTable.Metadata.Network,
		p.NonExistentTargetRefRouteTable.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentTargetRefURL, scenario.MockParams); err != nil {
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
