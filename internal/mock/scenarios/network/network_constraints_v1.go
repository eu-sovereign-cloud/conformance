package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureNetworkConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.NetworkConstraintsValidationV1Params) error {
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

	// Over-length name Validation
	overLengthNameNetwork := p.OverLengthNameNetwork
	overLengthNameURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, overLengthNameNetwork.Metadata.Tenant, overLengthNameNetwork.Metadata.Workspace, overLengthNameNetwork.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name Validation
	invalidPatternNameNetwork := p.InvalidPatternNameNetwork
	invalidPatternNameURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, invalidPatternNameNetwork.Metadata.Tenant, invalidPatternNameNetwork.Metadata.Workspace, invalidPatternNameNetwork.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value Validation
	overLengthLabelNetwork := p.OverLengthLabelValueNetwork
	overLengthLabelURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, overLengthLabelNetwork.Metadata.Tenant, overLengthLabelNetwork.Metadata.Workspace, overLengthLabelNetwork.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value Validation
	overLengthAnnotationNetwork := p.OverLengthAnnotationNetwork
	overLengthAnnotationURL := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, overLengthAnnotationNetwork.Metadata.Tenant, overLengthAnnotationNetwork.Metadata.Workspace, overLengthAnnotationNetwork.Metadata.Name)
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
