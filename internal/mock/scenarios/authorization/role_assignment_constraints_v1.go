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

	// Empty roles Validation
	emptyRolesRA := p.EmptyRolesRoleAssignment
	emptyRolesURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, emptyRolesRA.Metadata.Tenant, emptyRolesRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyRolesURL, scenario.MockParams); err != nil {
		return err
	}

	// Over maxItems roles Validation
	overMaxItemsRolesRA := p.OverMaxItemsRolesRoleAssignment
	overMaxItemsRolesURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overMaxItemsRolesRA.Metadata.Tenant, overMaxItemsRolesRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxItemsRolesURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty role value Validation
	emptyRoleValueRA := p.EmptyRoleValueRoleAssignment
	emptyRoleValueURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, emptyRoleValueRA.Metadata.Tenant, emptyRoleValueRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyRoleValueURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty subs Validation
	emptySubsRA := p.EmptySubsRoleAssignment
	emptySubsURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, emptySubsRA.Metadata.Tenant, emptySubsRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptySubsURL, scenario.MockParams); err != nil {
		return err
	}

	// Over maxItems subs Validation
	overMaxItemsSubsRA := p.OverMaxItemsSubsRoleAssignment
	overMaxItemsSubsURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overMaxItemsSubsRA.Metadata.Tenant, overMaxItemsSubsRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxItemsSubsURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty sub value Validation
	emptySubValueRA := p.EmptySubValueRoleAssignment
	emptySubValueURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, emptySubValueRA.Metadata.Tenant, emptySubValueRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptySubValueURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty scopes Validation
	emptyScopesRA := p.EmptyScopesRoleAssignment
	emptyScopesURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, emptyScopesRA.Metadata.Tenant, emptyScopesRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyScopesURL, scenario.MockParams); err != nil {
		return err
	}

	// Over maxItems scopes Validation
	overMaxItemsScopesRA := p.OverMaxItemsScopesRoleAssignment
	overMaxItemsScopesURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overMaxItemsScopesRA.Metadata.Tenant, overMaxItemsScopesRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxItemsScopesURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty scope tenant value Validation
	emptyScopeTenantValueRA := p.EmptyScopeTenantValueRoleAssignment
	emptyScopeTenantValueURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, emptyScopeTenantValueRA.Metadata.Tenant, emptyScopeTenantValueRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyScopeTenantValueURL, scenario.MockParams); err != nil {
		return err
	}

	// Over maxItems scope tenants Validation
	overMaxItemsScopeTenantsRA := p.OverMaxItemsScopeTenantsRoleAssignment
	overMaxItemsScopeTenantsURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overMaxItemsScopeTenantsRA.Metadata.Tenant, overMaxItemsScopeTenantsRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxItemsScopeTenantsURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty scope region value Validation
	emptyScopeRegionValueRA := p.EmptyScopeRegionValueRoleAssignment
	emptyScopeRegionValueURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, emptyScopeRegionValueRA.Metadata.Tenant, emptyScopeRegionValueRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyScopeRegionValueURL, scenario.MockParams); err != nil {
		return err
	}

	// Over maxItems scope regions Validation
	overMaxItemsScopeRegionsRA := p.OverMaxItemsScopeRegionsRoleAssignment
	overMaxItemsScopeRegionsURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overMaxItemsScopeRegionsRA.Metadata.Tenant, overMaxItemsScopeRegionsRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxItemsScopeRegionsURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty scope workspace value Validation
	emptyScopeWorkspaceValueRA := p.EmptyScopeWorkspaceValueRoleAssignment
	emptyScopeWorkspaceValueURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, emptyScopeWorkspaceValueRA.Metadata.Tenant, emptyScopeWorkspaceValueRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyScopeWorkspaceValueURL, scenario.MockParams); err != nil {
		return err
	}

	// Over maxItems scope workspaces Validation
	overMaxItemsScopeWorkspacesRA := p.OverMaxItemsScopeWorkspacesRoleAssignment
	overMaxItemsScopeWorkspacesURL := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, overMaxItemsScopeWorkspacesRA.Metadata.Tenant, overMaxItemsScopeWorkspacesRA.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxItemsScopeWorkspacesURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
