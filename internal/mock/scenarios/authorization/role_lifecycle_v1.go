package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureRoleLifecycleScenarioV1(scenario *mockscenarios.Scenario, params params.RoleLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	role := params.RoleInitial

	// Generate URLs
	roleUrl := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, role.Metadata.Tenant, role.Metadata.Name)

	// Create a role
	if err := configurator.ConfigureCreateRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created role
	if err := configurator.ConfigureGetCreatingRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the role
	params.RoleUpdated.Status = role.Status
	role = params.RoleUpdated
	if err := configurator.ConfigureUpdateRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated role
	if err := configurator.ConfigureGetUpdatingRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the role
	if err := configurator.ConfigureDeleteStub(roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted role
	if err := configurator.ConfigureGetDeletingRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(roleUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
