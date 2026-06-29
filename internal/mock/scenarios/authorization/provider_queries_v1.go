package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/conformance/pkg/mock"
	"github.com/eu-sovereign-cloud/conformance/pkg/mock/stubs"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func ConfigureProviderQueriesV1(scenario *mockscenarios.Scenario, params params.AuthorizationProviderQueriesV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}
	roles := params.Roles
	roleAssignments := params.RoleAssignments

	// Generate URLs
	roleUrl := generators.GenerateRoleListURL(sdkconsts.AuthorizationProviderV1Name, roles.Items[0].Metadata.Tenant)
	roleAssignmentUrl := generators.GenerateRoleAssignmentListURL(sdkconsts.AuthorizationProviderV1Name, roleAssignments.Items[0].Metadata.Tenant)

	// Roles
	err = stubs.BulkCreateRolesStubV1(configurator, scenario.MockParams, roles.Items)
	if err != nil {
		return err
	}
	rolesResponse := &params.Roles

	// List roles
	if err := configurator.ConfigureListRoleStub(rolesResponse, roleUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List roles with limit 1
	rolesResponse.Items = roles.Items[:1]
	if err := configurator.ConfigureListRoleStub(rolesResponse, roleUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
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
	rolesResponse.Items = rolesListWithLabel(roles.Items)
	if err := configurator.ConfigureListRoleStub(rolesResponse, roleUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List roles with limit and label
	rolesResponse.Items = rolesListWithLabel(roles.Items)[:1]
	if err := configurator.ConfigureListRoleStub(rolesResponse, roleUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Role assignments
	err = stubs.BulkCreateRoleAssignmentsStubV1(configurator, scenario.MockParams, roleAssignments.Items)
	if err != nil {
		return err
	}
	roleAssignResponse := &params.RoleAssignments

	// List role assignments
	if err := configurator.ConfigureListRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List role assignments with limit 1
	roleAssignResponse.Items = roleAssignments.Items[:1]
	if err := configurator.ConfigureListRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
	}

	// List role assignments with label
	rolesAssignWithLabel := func(rolesAssignList []schema.RoleAssignment) []schema.RoleAssignment {
		var filteredRoles []schema.RoleAssignment
		for _, roleAssign := range rolesAssignList {
			if val, ok := roleAssign.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredRoles = append(filteredRoles, roleAssign)
			}
		}
		return filteredRoles
	}
	roleAssignResponse.Items = rolesAssignWithLabel(roleAssignments.Items)
	if err := configurator.ConfigureListRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List role assignments with limit and label
	roleAssignResponse.Items = rolesAssignWithLabel(roleAssignments.Items)[:1]
	if err := configurator.ConfigureListRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Delete role assignments
	for _, roleAssignment := range roleAssignments.Items {
		roleAssignUrl := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, roleAssignment.Metadata.Tenant, roleAssignment.Metadata.Name)

		// Delete the role assignment
		if err := configurator.ConfigureDeleteStub(roleAssignUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted role assignment
		if err := configurator.ConfigureGetDeletingRoleAssignmentStub(&roleAssignment, roleAssignUrl, scenario.MockParams); err != nil {
			return err
		}
		if err := configurator.ConfigureGetNotFoundStub(roleAssignUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete roles
	for _, role := range roles.Items {
		roleUrl := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, role.Metadata.Tenant, role.Metadata.Name)

		// Delete the role assignment
		if err := configurator.ConfigureDeleteStub(roleUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted role assignment
		if err := configurator.ConfigureGetDeletingRoleStub(&role, roleUrl, scenario.MockParams); err != nil {
			return err
		}
		if err := configurator.ConfigureGetNotFoundStub(roleUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
