package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureRoleAssignmentConstraintsValidationV1 sets up mock stubs for the role assignment
// constraints validation suite. Each role assignment in the params targets a different
// constraint Validation, all returning 422 Unprocessable Entity.
func ConfigureRoleAssignmentConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.RoleAssignmentConstraintsValidationV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	// Over-length name Validation
	overLengthNameRoleAssignment := p.OverLengthNameRoleAssignment
	overLengthNameURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overLengthNameRoleAssignment.Metadata.Tenant, overLengthNameRoleAssignment.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name Validation
	invalidPatternNameRoleAssignment := p.InvalidPatternNameRoleAssignment
	invalidPatternNameURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, invalidPatternNameRoleAssignment.Metadata.Tenant, invalidPatternNameRoleAssignment.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value Validation
	overLengthLabelRoleAssignment := p.OverLengthLabelValueRoleAssignment
	overLengthLabelURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overLengthLabelRoleAssignment.Metadata.Tenant, overLengthLabelRoleAssignment.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value Validation
	overLengthAnnotationRoleAssignment := p.OverLengthAnnotationRoleAssignment
	overLengthAnnotationURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overLengthAnnotationRoleAssignment.Metadata.Tenant, overLengthAnnotationRoleAssignment.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length sub Validation
	overLengthSubRA := p.OverLengthSubRoleAssignment
	overLengthSubURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overLengthSubRA.Metadata.Tenant, overLengthSubRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthSubURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length role name Validation
	overLengthRoleNameRA := p.OverLengthRoleNameRoleAssignment
	overLengthRoleNameURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overLengthRoleNameRA.Metadata.Tenant, overLengthRoleNameRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthRoleNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length scope tenant Validation
	overLengthScopeTenantRA := p.OverLengthScopeTenantRoleAssignment
	overLengthScopeTenantURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overLengthScopeTenantRA.Metadata.Tenant, overLengthScopeTenantRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthScopeTenantURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length scope region Validation
	overLengthScopeRegionRA := p.OverLengthScopeRegionRoleAssignment
	overLengthScopeRegionURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overLengthScopeRegionRA.Metadata.Tenant, overLengthScopeRegionRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthScopeRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length scope workspace Validation
	overLengthScopeWorkspaceRA := p.OverLengthScopeWorkspaceRoleAssignment
	overLengthScopeWorkspaceURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overLengthScopeWorkspaceRA.Metadata.Tenant, overLengthScopeWorkspaceRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthScopeWorkspaceURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
