package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureRoleAssignmentErrorV1 sets up mock stubs for the role assignment error
// scenarios suite. All invalid role assignment requests return 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create role assignment with non-existent role ref
//   - Create role assignment with non-existent tenant in scope
//   - Create role assignment with non-existent region in scope
//   - Create role assignment with non-existent workspace in scope
//   - Create role assignment with non-existent sub
func ConfigureRoleAssignmentErrorV1(scenario *mockscenarios.Scenario, p params.RoleAssignmentErrorV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	// Non-existent role ref
	nonExistentRoleRefURL := generators.GenerateRoleAssignmentURL(
		sdkconsts.AuthorizationProviderV1Name,
		p.NonExistentRoleRefRoleAssignment.Metadata.Tenant,
		p.NonExistentRoleRefRoleAssignment.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentRoleRefURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent tenant in scope
	nonExistentScopeTenantURL := generators.GenerateRoleAssignmentURL(
		sdkconsts.AuthorizationProviderV1Name,
		p.NonExistentScopeTenantRoleAssignment.Metadata.Tenant,
		p.NonExistentScopeTenantRoleAssignment.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentScopeTenantURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent region in scope
	nonExistentScopeRegionURL := generators.GenerateRoleAssignmentURL(
		sdkconsts.AuthorizationProviderV1Name,
		p.NonExistentScopeRegionRoleAssignment.Metadata.Tenant,
		p.NonExistentScopeRegionRoleAssignment.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentScopeRegionURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent workspace in scope
	nonExistentScopeWorkspaceURL := generators.GenerateRoleAssignmentURL(
		sdkconsts.AuthorizationProviderV1Name,
		p.NonExistentScopeWorkspaceRoleAssignment.Metadata.Tenant,
		p.NonExistentScopeWorkspaceRoleAssignment.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentScopeWorkspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Non-existent sub
	nonExistentSubURL := generators.GenerateRoleAssignmentURL(
		sdkconsts.AuthorizationProviderV1Name,
		p.NonExistentSubRoleAssignment.Metadata.Tenant,
		p.NonExistentSubRoleAssignment.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentSubURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
