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

func ConfigureListScenarioV1(scenario string, mockParams *mock.MockParams, suiteParams *params.AuthorizationListV1Params) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	roles := suiteParams.Roles
	roleAssignments := suiteParams.RoleAssignments

	configurator, err := stubs.NewStubConfigurator(scenario, mockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	roleUrl := generators.GenerateRoleListURL(constants.AuthorizationProviderV1, roles[0].Metadata.Tenant)
	roleAssignUrl := generators.GenerateRoleAssignmentListURL(constants.AuthorizationProviderV1, roleAssignments[0].Metadata.Tenant)

	// Create roles
	err = stubs.BulkCreateRolesStubV1(configurator, mockParams, roles)
	if err != nil {
		return nil, err
	}
	rolesResponse, err := builders.NewRoleIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(roles[0].Metadata.Tenant).
		Items(roles).
		Build()
	if err != nil {
		return nil, err
	}

	// List Roles
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, mockParams, nil); err != nil {
		return nil, err
	}

	// List Roles with limit 1
	rolesResponse.Items = roles[:1]
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, mockParams, mock.PathParamsLimit("1")); err != nil {
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
	rolesResponse.Items = rolesListWithLabel(roles)
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, mockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List roles with limit and label
	rolesResponse.Items = rolesListWithLabel(roles)[:1]
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, mockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create role assignments
	err = stubs.BulkCreateRoleAssignmentsStubV1(configurator, mockParams, roleAssignments)
	if err != nil {
		return nil, err
	}
	roleAssignResponse, err := builders.NewRoleAssignmentIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(roleAssignments[0].Metadata.Tenant).
		Items(roleAssignments).
		Build()
	if err != nil {
		return nil, err
	}

	// List RoleAssignments
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, mockParams, nil); err != nil {
		return nil, err
	}

	// List Roles with limit 1
	roleAssignResponse.Items = roleAssignments[:1]
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, mockParams, mock.PathParamsLimit("1")); err != nil {
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
	roleAssignResponse.Items = rolesAssignWithLabel(roleAssignments)
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, mockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List roles with limit and label
	roleAssignResponse.Items = rolesAssignWithLabel(roleAssignments)[:1]
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, mockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Delete role assignments
	for _, roleAssignment := range roleAssignments {
		roleAssignUrl := generators.GenerateRoleAssignmentURL(constants.AuthorizationProviderV1, roleAssignment.Metadata.Tenant, roleAssignment.Metadata.Name)

		// Delete the role assignment
		if err := configurator.ConfigureDeleteStub(roleAssignUrl, mockParams); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(roleAssignUrl, mockParams); err != nil {
			return nil, err
		}
	}

	// Delete roles
	for _, role := range roles {
		roleUrl := generators.GenerateRoleURL(constants.AuthorizationProviderV1, role.Metadata.Tenant, role.Metadata.Name)

		// Delete the role assignment
		if err := configurator.ConfigureDeleteStub(roleUrl, mockParams); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(roleUrl, mockParams); err != nil {
			return nil, err
		}
	}

	return configurator.Client, nil
}
