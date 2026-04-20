package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureRoleConstraintsViolationsV1 sets up mock stubs for the role constraints violations suite.
// Each role in the params targets a different constraint violation, all returning 422 Unprocessable Entity.
func ConfigureRoleConstraintsViolationsV1(scenario *mockscenarios.Scenario, p params.RoleConstraintsViolationsV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	// Over-length name violation
	overLengthNameRole := p.OverLengthNameRole
	overLengthNameURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overLengthNameRole.Metadata.Tenant, overLengthNameRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name violation
	invalidPatternNameRole := p.InvalidPatternNameRole
	invalidPatternNameURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, invalidPatternNameRole.Metadata.Tenant, invalidPatternNameRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value violation
	overLengthLabelRole := p.OverLengthLabelValueRole
	overLengthLabelURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overLengthLabelRole.Metadata.Tenant, overLengthLabelRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value violation
	overLengthAnnotationRole := p.OverLengthAnnotationRole
	overLengthAnnotationURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overLengthAnnotationRole.Metadata.Tenant, overLengthAnnotationRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
