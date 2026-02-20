package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
)

func ConfigureRoleAssignmentLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.RoleAssignmentLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	roleAssignment := *params.RoleAssignmentInitial

	// Generate URLs
	roleAssignmentUrl := generators.GenerateRoleAssignmentURL(constants.AuthorizationProviderV1, roleAssignment.Metadata.Tenant, roleAssignment.Metadata.Name)

	// Role assignment
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignment.Metadata.Name).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(roleAssignment.Metadata.Tenant).
		Spec(&roleAssignment.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a role assignment
	if err := configurator.ConfigureCreateRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created role assignment
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the role assignment
	roleAssignResponse.Spec = params.RoleAssignmentUpdated.Spec
	if err := configurator.ConfigureUpdateRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated role assignment
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the role assignment
	if err := configurator.ConfigureDeleteStub(roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted role assignment
	if err := configurator.ConfigureGetNotFoundStub(roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
