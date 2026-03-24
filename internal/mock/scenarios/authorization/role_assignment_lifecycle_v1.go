package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureRoleAssignmentLifecycleScenarioV1(scenario *mockscenarios.Scenario, params params.RoleAssignmentLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	roleAssignment := params.RoleAssignmentInitial

	// Generate URLs
	roleAssignmentUrl := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, roleAssignment.Metadata.Tenant, roleAssignment.Metadata.Name)

	// Create a role assignment
	if err := configurator.ConfigureCreateRoleAssignmentStub(roleAssignment, roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created role assignment
	if err := configurator.ConfigureGetCreatingRoleAssignmentStub(roleAssignment, roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignment, roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the role assignment
	params.RoleAssignmentUpdated.Status = roleAssignment.Status
	roleAssignment = params.RoleAssignmentUpdated
	if err := configurator.ConfigureUpdateRoleAssignmentStub(roleAssignment, roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated role assignment
	if err := configurator.ConfigureGetUpdatingRoleAssignmentStub(roleAssignment, roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignment, roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the role assignment
	if err := configurator.ConfigureDeleteStub(roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted role assignment
	if err := configurator.ConfigureGetDeletingRoleAssignmentStub(roleAssignment, roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(roleAssignmentUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
