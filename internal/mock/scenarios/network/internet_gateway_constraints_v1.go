package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureInternetGatewayConstraintsViolationsV1(scenario *mockscenarios.Scenario, p params.InternetGatewayConstraintsViolationsV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)

	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	violations := []*params.InternetGatewayConstraintsViolationsV1Params{}
	urls := []string{
		generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, p.OverLengthNameInternetGateway.Metadata.Tenant, p.OverLengthNameInternetGateway.Metadata.Workspace, p.OverLengthNameInternetGateway.Metadata.Name),
		generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, p.InvalidPatternNameInternetGateway.Metadata.Tenant, p.InvalidPatternNameInternetGateway.Metadata.Workspace, p.InvalidPatternNameInternetGateway.Metadata.Name),
		generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, p.OverLengthLabelValueInternetGateway.Metadata.Tenant, p.OverLengthLabelValueInternetGateway.Metadata.Workspace, p.OverLengthLabelValueInternetGateway.Metadata.Name),
		generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, p.OverLengthAnnotationInternetGateway.Metadata.Tenant, p.OverLengthAnnotationInternetGateway.Metadata.Workspace, p.OverLengthAnnotationInternetGateway.Metadata.Name),
	}
	_ = violations

	for _, url := range urls {
		if err := configurator.ConfigurePutUnprocessableEntityStub(url, scenario.MockParams); err != nil {
			return err
		}
	}

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
