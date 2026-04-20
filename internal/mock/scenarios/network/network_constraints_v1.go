package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureNetworkConstraintsViolationsV1(scenario *mockscenarios.Scenario, p params.NetworkConstraintsViolationsV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	internetGateway := p.InternetGateway

	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	gatewayURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, internetGateway.Metadata.Tenant, internetGateway.Metadata.Workspace, internetGateway.Metadata.Name)

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

	// Create internet gateway (needed as dependency for network spec)
	if err := configurator.ConfigureCreateInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	for _, url := range []string{
		generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, p.OverLengthNameNetwork.Metadata.Tenant, p.OverLengthNameNetwork.Metadata.Workspace, p.OverLengthNameNetwork.Metadata.Name),
		generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, p.InvalidPatternNameNetwork.Metadata.Tenant, p.InvalidPatternNameNetwork.Metadata.Workspace, p.InvalidPatternNameNetwork.Metadata.Name),
		generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, p.OverLengthLabelValueNetwork.Metadata.Tenant, p.OverLengthLabelValueNetwork.Metadata.Workspace, p.OverLengthLabelValueNetwork.Metadata.Name),
		generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, p.OverLengthAnnotationNetwork.Metadata.Tenant, p.OverLengthAnnotationNetwork.Metadata.Workspace, p.OverLengthAnnotationNetwork.Metadata.Name),
	} {
		if err := configurator.ConfigurePutUnprocessableEntityStub(url, scenario.MockParams); err != nil {
			return err
		}
	}

	// Teardown internet gateway
	if err := configurator.ConfigureDeleteStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingInternetGatewayStub(internetGateway, gatewayURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(gatewayURL, scenario.MockParams); err != nil {
		return err
	}

	// Teardown workspace
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
