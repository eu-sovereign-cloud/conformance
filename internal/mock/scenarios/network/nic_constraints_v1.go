package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureNicConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.NicConstraintsValidationV1Params) error {
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
	overLengthNameNic := p.OverLengthNameNic
	overLengthNameURL := generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, overLengthNameNic.Metadata.Tenant, overLengthNameNic.Metadata.Workspace, overLengthNameNic.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name Validation
	invalidPatternNameNic := p.InvalidPatternNameNic
	invalidPatternNameURL := generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, invalidPatternNameNic.Metadata.Tenant, invalidPatternNameNic.Metadata.Workspace, invalidPatternNameNic.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value Validation
	overLengthLabelNic := p.OverLengthLabelValueNic
	overLengthLabelURL := generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, overLengthLabelNic.Metadata.Tenant, overLengthLabelNic.Metadata.Workspace, overLengthLabelNic.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value Validation
	overLengthAnnotationNic := p.OverLengthAnnotationNic
	overLengthAnnotationURL := generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, overLengthAnnotationNic.Metadata.Tenant, overLengthAnnotationNic.Metadata.Workspace, overLengthAnnotationNic.Metadata.Name)
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
