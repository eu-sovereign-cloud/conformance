package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/wiremock/go-wiremock"
)

func CreateAuthorizationLifecycleScenarioV1(scenario string, params *AuthorizationParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	roleUrl := generators.GenerateRoleURL(authorizationProviderV1, params.Tenant, params.Role.Name)
	roleAssignUrl := generators.GenerateRoleAssignmentURL(authorizationProviderV1, params.Tenant, params.RoleAssignment.Name)

	// Role
	roleResponse, err := builders.NewRoleBuilder().
		Name(params.Role.Name).
		Provider(authorizationProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).
		Spec(params.Role.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a role
	if err := configurator.configureCreateRoleStub(roleResponse, roleUrl, params); err != nil {
		return nil, err
	}

	// Get the created role
	if err := configurator.configureGetActiveRoleStub(roleResponse, roleUrl, params); err != nil {
		return nil, err
	}

	// Update the role
	roleResponse.Spec = *params.Role.UpdatedSpec
	if err := configurator.configureUpdateRoleStub(roleResponse, roleUrl, params); err != nil {
		return nil, err
	}

	// Get the updated role
	if err := configurator.configureGetActiveRoleStub(roleResponse, roleUrl, params); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name(params.RoleAssignment.Name).
		Provider(authorizationProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).
		Spec(params.RoleAssignment.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a role assignment
	if err := configurator.configureCreateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params); err != nil {
		return nil, err
	}

	// Get the created role assignment
	if err := configurator.configureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params); err != nil {
		return nil, err
	}

	// Update the role assignment
	roleAssignResponse.Spec = *params.RoleAssignment.UpdatedSpec
	if err := configurator.configureUpdateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params); err != nil {
		return nil, err
	}

	// Get the updated role assignment
	if err := configurator.configureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params); err != nil {
		return nil, err
	}

	// Delete the role assignment
	if err := configurator.configureDeleteStub(roleAssignUrl, params, false); err != nil {
		return nil, err
	}

	// Get the deleted role assignment
	if err := configurator.configureGetNotFoundStub(roleAssignUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the role
	if err := configurator.configureDeleteStub(roleUrl, params, false); err != nil {
		return nil, err
	}

	// Get the deleted role
	if err := configurator.configureGetNotFoundStub(roleUrl, params, true); err != nil {
		return nil, err
	}

	return configurator.client, nil
}
