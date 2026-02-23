package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
)

func ConfigureRouteTableLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.RouteTableLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := *params.Workspace
	network := *params.Network
	internetGateway := *params.InternetGateway
	routeTableInitial := *params.RouteTableInitial
	routeTableUpdated := *params.RouteTableUpdated

	workspaceURL := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	networkURL := generators.GenerateNetworkURL(constants.NetworkProviderV1, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)
	routeTableURL := generators.GenerateRouteTableURL(constants.NetworkProviderV1, routeTableInitial.Metadata.Tenant, routeTableInitial.Metadata.Workspace, routeTableInitial.Metadata.Network, routeTableInitial.Metadata.Name)
	gatewayURL := generators.GenerateInternetGatewayURL(constants.NetworkProviderV1, internetGateway.Metadata.Tenant, internetGateway.Metadata.Workspace, internetGateway.Metadata.Name)
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

	// Get created workspace
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

	// Get created network
	if err := configurator.ConfigureGetActiveNetworkStub(networkResponse, networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Route table
	routeInitialResponse, err := builders.NewRouteTableBuilder().
		Name(routeTableInitial.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(routeTableInitial.Metadata.Tenant).Workspace(routeTableInitial.Metadata.Workspace).Network(routeTableInitial.Metadata.Network).Region(routeTableInitial.Metadata.Region).
		Spec(&routeTableInitial.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create route table
	if err := configurator.ConfigureCreateRouteTableStub(routeInitialResponse, routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Get created route table
	if err := configurator.ConfigureGetActiveRouteTableStub(routeInitialResponse, routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Update route table
	routeUpdatedResponse, err := builders.NewRouteTableBuilder().
		Name(routeTableUpdated.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(routeTableUpdated.Metadata.Tenant).Workspace(routeTableUpdated.Metadata.Workspace).Network(routeTableUpdated.Metadata.Network).Region(routeTableUpdated.Metadata.Region).
		Spec(&routeTableUpdated.Spec).
		Build()
	if err != nil {
		return err
	}

	if err := configurator.ConfigureUpdateRouteTableStub(routeUpdatedResponse, routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Get updated route table
	if err := configurator.ConfigureGetActiveRouteTableStub(routeUpdatedResponse, routeTableURL, scenario.MockParams); err != nil {
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

	// Get created internet gateway
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayInitialResponse, gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete internet gateway
	if err := configurator.ConfigureDeleteStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete route table
	if err := configurator.ConfigureDeleteStub(routeTableURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(routeTableURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the network
	if err := configurator.ConfigureDeleteStub(networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted network
	if err := configurator.ConfigureGetNotFoundStub(networkURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
