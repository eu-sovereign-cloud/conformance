package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/wiremock/go-wiremock"
)

func ConfigureListScenarioV1(scenario string, params *params.AuthorizationListParamsV1) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	configurator, err := stubs.NewStubConfigurator(scenario, params.MockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	roleUrl := generators.GenerateRoleListURL(constants.AuthorizationProviderV1, params.Tenant)
	roleAssignUrl := generators.GenerateRoleAssignmentListURL(constants.AuthorizationProviderV1, params.Tenant)

	// Create roles
	rolesList, err := stubs.BulkCreateRolesStubV1(configurator, params.BaseParams, params.Roles)
	if err != nil {
		return nil, err
	}
	rolesResponse, err := builders.NewRoleIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(params.Tenant).
		Items(rolesList).
		Build()
	if err != nil {
		return nil, err
	}

	// List Roles
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, params.MockParams, nil); err != nil {
		return nil, err
	}

	// List Roles with limit 1
	rolesResponse.Items = rolesList[:1]
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, params.MockParams, mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List roles with label
	rolesListWithLabel := func(rolesList []schema.Role) []schema.Role {
		var filteredRoles []schema.Role
		for _, role := range rolesList {
			if val, ok := role.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredRoles = append(filteredRoles, role)
			}
		}
		return filteredRoles
	}
	rolesResponse.Items = rolesListWithLabel(rolesList)
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, params.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List roles with limit and label
	rolesResponse.Items = rolesListWithLabel(rolesList)[:1]
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, params.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create role assignments
	roleAssignmentsList, err := stubs.BulkCreateRoleAssignmentsStubV1(configurator, params.BaseParams, params.RoleAssignments)
	if err != nil {
		return nil, err
	}
	roleAssignResponse, err := builders.NewRoleAssignmentIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(params.Tenant).
		Items(roleAssignmentsList).
		Build()
	if err != nil {
		return nil, err
	}

	// List RoleAssignments
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.MockParams, nil); err != nil {
		return nil, err
	}

	// List Roles with limit 1
	roleAssignResponse.Items = roleAssignmentsList[:1]
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.MockParams, mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List roles with label
	rolesAssignWithLabel := func(rolesAssignList []schema.RoleAssignment) []schema.RoleAssignment {
		var filteredRoles []schema.RoleAssignment
		for _, roleAssign := range rolesAssignList {
			if val, ok := roleAssign.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredRoles = append(filteredRoles, roleAssign)
			}
		}
		return filteredRoles
	}
	roleAssignResponse.Items = rolesAssignWithLabel(roleAssignmentsList)
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List roles with limit and label
	roleAssignResponse.Items = rolesAssignWithLabel(roleAssignmentsList)[:1]
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Delete role assignments
	for _, roleAssignment := range roleAssignmentsList {
		roleAssignUrl := generators.GenerateRoleAssignmentURL(constants.AuthorizationProviderV1, params.Tenant, roleAssignment.Metadata.Name)

		// Delete the role assignment
		if err := configurator.ConfigureDeleteStub(roleAssignUrl, params.MockParams); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(roleAssignUrl, params.MockParams); err != nil {
			return nil, err
		}
	}

	// Delete roles
	for _, role := range rolesList {
		roleUrl := generators.GenerateRoleURL(constants.AuthorizationProviderV1, params.Tenant, role.Metadata.Name)

		// Delete the role assignment
		if err := configurator.ConfigureDeleteStub(roleUrl, params.MockParams); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(roleUrl, params.MockParams); err != nil {
			return nil, err
		}
	}

	return configurator.Client, nil
}
