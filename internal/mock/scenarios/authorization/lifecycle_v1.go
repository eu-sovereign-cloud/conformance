package authorization

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"

	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, params *mock.AuthorizationLifeCycleParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := stubs.NewStubConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	roleUrl := generators.GenerateRoleURL(constants.AuthorizationProviderV1, params.Tenant, params.Role.Name)
	roleAssignUrl := generators.GenerateRoleAssignmentURL(constants.AuthorizationProviderV1, params.Tenant, params.RoleAssignment.Name)

	// Role
	roleResponse, err := builders.NewRoleBuilder().
		Name(params.Role.Name).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).
		Spec(params.Role.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a role
	if err := configurator.ConfigureCreateRoleStub(roleResponse, roleUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created role
	if err := configurator.ConfigureGetActiveRoleStub(roleResponse, roleUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Update the role
	roleResponse.Spec = *params.Role.UpdatedSpec
	if err := configurator.ConfigureUpdateRoleStub(roleResponse, roleUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated role
	if err := configurator.ConfigureGetActiveRoleStub(roleResponse, roleUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name(params.RoleAssignment.Name).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).
		Spec(params.RoleAssignment.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a role assignment
	if err := configurator.ConfigureCreateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created role assignment
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Update the role assignment
	roleAssignResponse.Spec = *params.RoleAssignment.UpdatedSpec
	if err := configurator.ConfigureUpdateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated role assignment
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the role assignment
	if err := configurator.ConfigureDeleteStub(roleAssignUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted role assignment
	if err := configurator.ConfigureGetNotFoundStub(roleAssignUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the role
	if err := configurator.ConfigureDeleteStub(roleUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted role
	if err := configurator.ConfigureGetNotFoundStub(roleUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
