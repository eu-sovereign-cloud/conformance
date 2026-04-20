package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureRoleAssignmentConstraintsViolationsV1 sets up mock stubs for the role assignment
// constraints violations suite. Each role assignment in the params targets a different
// constraint violation, all returning 422 Unprocessable Entity.
func ConfigureRoleAssignmentConstraintsViolationsV1(scenario *mockscenarios.Scenario, p params.RoleAssignmentConstraintsViolationsV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	// Over-length name violation
	overLengthNameRoleAssignment := p.OverLengthNameRoleAssignment
	overLengthNameURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overLengthNameRoleAssignment.Metadata.Tenant, overLengthNameRoleAssignment.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name violation
	invalidPatternNameRoleAssignment := p.InvalidPatternNameRoleAssignment
	invalidPatternNameURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, invalidPatternNameRoleAssignment.Metadata.Tenant, invalidPatternNameRoleAssignment.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value violation
	overLengthLabelRoleAssignment := p.OverLengthLabelValueRoleAssignment
	overLengthLabelURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overLengthLabelRoleAssignment.Metadata.Tenant, overLengthLabelRoleAssignment.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value violation
	overLengthAnnotationRoleAssignment := p.OverLengthAnnotationRoleAssignment
	overLengthAnnotationURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overLengthAnnotationRoleAssignment.Metadata.Tenant, overLengthAnnotationRoleAssignment.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
