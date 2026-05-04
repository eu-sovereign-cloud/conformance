package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureInternetGatewayConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.InternetGatewayConstraintsValidationV1Params) error {
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

	// Over-length name Validation
	overLengthNameInternetGateway := p.OverLengthNameInternetGateway
	overLengthNameURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, overLengthNameInternetGateway.Metadata.Tenant, overLengthNameInternetGateway.Metadata.Workspace, overLengthNameInternetGateway.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name Validation
	invalidPatternNameInternetGateway := p.InvalidPatternNameInternetGateway
	invalidPatternNameURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, invalidPatternNameInternetGateway.Metadata.Tenant, invalidPatternNameInternetGateway.Metadata.Workspace, invalidPatternNameInternetGateway.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value Validation
	overLengthLabelInternetGateway := p.OverLengthLabelValueInternetGateway
	overLengthLabelURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, overLengthLabelInternetGateway.Metadata.Tenant, overLengthLabelInternetGateway.Metadata.Workspace, overLengthLabelInternetGateway.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value Validation
	overLengthAnnotationInternetGateway := p.OverLengthAnnotationInternetGateway
	overLengthAnnotationURL := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, overLengthAnnotationInternetGateway.Metadata.Tenant, overLengthAnnotationInternetGateway.Metadata.Workspace, overLengthAnnotationInternetGateway.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
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
