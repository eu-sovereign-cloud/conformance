package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
)

func ConfigureRoleLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.RoleLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	role := *params.RoleInitial

	// Generate URLs
	roleUrl := generators.GenerateRoleURL(constants.AuthorizationProviderV1, role.Metadata.Tenant, role.Metadata.Name)

	// Role
	roleResponse, err := builders.NewRoleBuilder().
		Name(role.Metadata.Name).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(role.Metadata.Tenant).
		Spec(&role.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a role
	if err := configurator.ConfigureCreateRoleStub(roleResponse, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created role
	if err := configurator.ConfigureGetCreatingRoleStub(roleResponse, roleUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleStub(roleResponse, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the role
	roleResponse.Spec = params.RoleUpdated.Spec
	if err := configurator.ConfigureUpdateRoleStub(roleResponse, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated role
	if err := configurator.ConfigureGetUpdatingRoleStub(roleResponse, roleUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleStub(roleResponse, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the role
	if err := configurator.ConfigureDeleteStub(roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted role
	if err := configurator.ConfigureGetNotFoundStub(roleUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
