package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"

	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, mockParams *mock.MockParams, suiteParams *params.AuthorizationLifeCycleV1Params) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	role := *suiteParams.RoleInitial
	roleAssignment := *suiteParams.RoleAssignmentInitial

	configurator, err := stubs.NewStubConfigurator(scenario, mockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	roleUrl := generators.GenerateRoleURL(constants.AuthorizationProviderV1, role.Metadata.Tenant, role.Metadata.Name)
	roleAssignUrl := generators.GenerateRoleAssignmentURL(constants.AuthorizationProviderV1, roleAssignment.Metadata.Tenant, roleAssignment.Metadata.Name)
	// Role
	roleResponse, err := builders.NewRoleBuilder().
		Name(role.Metadata.Name).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(role.Metadata.Tenant).
		Spec(&role.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a role
	if err := configurator.ConfigureCreateRoleStub(roleResponse, roleUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created role
	if err := configurator.ConfigureGetActiveRoleStub(roleResponse, roleUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the role
	roleResponse.Spec = suiteParams.RoleUpdated.Spec
	if err := configurator.ConfigureUpdateRoleStub(roleResponse, roleUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated role
	if err := configurator.ConfigureGetActiveRoleStub(roleResponse, roleUrl, mockParams); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignment.Metadata.Name).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(roleAssignment.Metadata.Tenant).
		Spec(&roleAssignment.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a role assignment
	if err := configurator.ConfigureCreateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created role assignment
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the role assignment
	roleAssignResponse.Spec = suiteParams.RoleAssignmentUpdated.Spec
	if err := configurator.ConfigureUpdateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated role assignment
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the role assignment
	if err := configurator.ConfigureDeleteStub(roleAssignUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted role assignment
	if err := configurator.ConfigureGetNotFoundStub(roleAssignUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the role
	if err := configurator.ConfigureDeleteStub(roleUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted role
	if err := configurator.ConfigureGetNotFoundStub(roleUrl, mockParams); err != nil {
		return nil, err
	}

	// Finish the stubs configuration
	if client, err := configurator.Finish(); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}
