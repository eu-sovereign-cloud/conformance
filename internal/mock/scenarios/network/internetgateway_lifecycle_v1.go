package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
)

func ConfigureInternetGatewayLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.InternetGatewayLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := *params.Workspace
	internetGatewayInitial := *params.InternetGatewayInitial
	internetGatewayUpdated := *params.InternetGatewayUpdated

	workspaceURL := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	gatewayURL := generators.GenerateInternetGatewayURL(constants.NetworkProviderV1, internetGatewayInitial.Metadata.Tenant, internetGatewayInitial.Metadata.Workspace, internetGatewayInitial.Metadata.Name)

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

	// Internet gateway
	gatewayInitialResponse, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayInitial.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(internetGatewayInitial.Metadata.Tenant).Workspace(internetGatewayInitial.Metadata.Workspace).Region(internetGatewayInitial.Metadata.Region).
		Spec(&internetGatewayInitial.Spec).
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

	// Update internet gateway (change egressOnly)
	gatewayUpdatedResponse, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayUpdated.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(internetGatewayUpdated.Metadata.Tenant).Workspace(internetGatewayUpdated.Metadata.Workspace).Region(internetGatewayUpdated.Metadata.Region).
		Spec(&internetGatewayUpdated.Spec).
		Build()
	if err != nil {
		return err
	}

	if err := configurator.ConfigureUpdateInternetGatewayStub(gatewayUpdatedResponse, gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Get updated internet gateway
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayUpdatedResponse, gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete internet gateway
	if err := configurator.ConfigureDeleteStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete workspace
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
