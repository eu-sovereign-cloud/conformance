package mockworkspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureWorkspaceConstraintsViolationsV1 sets up mock stubs for the workspace constraints
// violations suite. Each workspace in the params targets a different constraint violation,
// all returning 422 Unprocessable Entity.
func ConfigureWorkspaceConstraintsViolationsV1(scenario *mockscenarios.Scenario, p params.WorkspaceConstraintsValidationV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	// Over-length name violation
	overLengthNameURL := generators.GenerateWorkspaceURL(
		sdkconsts.WorkspaceProviderV1Name,
		p.OverLengthNameWorkspace.Metadata.Tenant,
		p.OverLengthNameWorkspace.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name violation
	invalidPatternNameURL := generators.GenerateWorkspaceURL(
		sdkconsts.WorkspaceProviderV1Name,
		p.InvalidPatternNameWorkspace.Metadata.Tenant,
		p.InvalidPatternNameWorkspace.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value violation
	overLengthLabelURL := generators.GenerateWorkspaceURL(
		sdkconsts.WorkspaceProviderV1Name,
		p.OverLengthLabelValueWorkspace.Metadata.Tenant,
		p.OverLengthLabelValueWorkspace.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value violation
	overLengthAnnotationURL := generators.GenerateWorkspaceURL(
		sdkconsts.WorkspaceProviderV1Name,
		p.OverLengthAnnotationWorkspace.Metadata.Tenant,
		p.OverLengthAnnotationWorkspace.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
