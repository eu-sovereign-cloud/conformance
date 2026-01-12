package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"

	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, params *params.AuthorizationLifeCycleParamsV1) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	configurator, err := stubs.NewStubConfigurator(scenario, params.MockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	roleUrl := generators.GenerateRoleURL(constants.AuthorizationProviderV1, params.RoleInitial.Metadata.Tenant, params.RoleInitial.Metadata.Name)
	roleAssignUrl := generators.GenerateRoleAssignmentURL(constants.AuthorizationProviderV1, params.RoleAssignmentInitial.Metadata.Tenant, params.RoleAssignmentInitial.Metadata.Name)

	// Role
	roleResponse, err := builders.NewRoleBuilder().
		Name(params.RoleInitial.Metadata.Name).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.RoleInitial.Metadata.Tenant).
		Spec(&params.RoleInitial.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a role
	if err := configurator.ConfigureCreateRoleStub(roleResponse, roleUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created role
	if err := configurator.ConfigureGetActiveRoleStub(roleResponse, roleUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the role
	roleResponse.Spec = params.RoleUpdated.Spec
	if err := configurator.ConfigureUpdateRoleStub(roleResponse, roleUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated role
	if err := configurator.ConfigureGetActiveRoleStub(roleResponse, roleUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name(params.RoleAssignmentInitial.Metadata.Name).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.RoleAssignmentInitial.Metadata.Tenant).
		Spec(&params.RoleAssignmentInitial.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a role assignment
	if err := configurator.ConfigureCreateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created role assignment
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the role assignment
	roleAssignResponse.Spec = params.RoleAssignmentUpdated.Spec
	if err := configurator.ConfigureUpdateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated role assignment
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the role assignment
	if err := configurator.ConfigureDeleteStub(roleAssignUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted role assignment
	if err := configurator.ConfigureGetNotFoundStub(roleAssignUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the role
	if err := configurator.ConfigureDeleteStub(roleUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted role
	if err := configurator.ConfigureGetNotFoundStub(roleUrl, params.MockParams); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
