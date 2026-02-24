package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func ConfigureProviderQueriesV1(scenario *mockscenarios.Scenario, params *params.AuthorizationProviderQueriesV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}
	roles := params.Roles
	roleAssignments := params.RoleAssignments

	// Generate URLs
	roleUrl := generators.GenerateRoleListURL(sdkconsts.AuthorizationProviderV1Name, roles[0].Metadata.Tenant)
	roleAssignmentUrl := generators.GenerateRoleAssignmentListURL(sdkconsts.AuthorizationProviderV1Name, roleAssignments[0].Metadata.Tenant)

	// Create roles
	err = stubs.BulkCreateRolesStubV1(configurator, scenario.MockParams, roles)
	if err != nil {
		return err
	}
	rolesResponse, err := builders.NewRoleIteratorBuilder().
		Provider(sdkconsts.StorageProviderV1Name).
		Tenant(roles[0].Metadata.Tenant).
		Items(roles).
		Build()
	if err != nil {
		return err
	}

	// List Roles
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List Roles with limit 1
	rolesResponse.Items = roles[:1]
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
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
	rolesResponse.Items = rolesListWithLabel(roles)
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List roles with limit and label
	rolesResponse.Items = rolesListWithLabel(roles)[:1]
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Create role assignments
	err = stubs.BulkCreateRoleAssignmentsStubV1(configurator, scenario.MockParams, roleAssignments)
	if err != nil {
		return err
	}
	roleAssignResponse, err := builders.NewRoleAssignmentIteratorBuilder().
		Provider(sdkconsts.StorageProviderV1Name).
		Tenant(roleAssignments[0].Metadata.Tenant).
		Items(roleAssignments).
		Build()
	if err != nil {
		return err
	}

	// List RoleAssignments
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List Roles with limit 1
	roleAssignResponse.Items = roleAssignments[:1]
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
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
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List roles with limit and label
	roleAssignResponse.Items = rolesAssignWithLabel(roleAssignments)[:1]
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Delete role assignments
	for _, roleAssignment := range roleAssignments {
		roleAssignUrl := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, roleAssignment.Metadata.Tenant, roleAssignment.Metadata.Name)

		// Delete the role assignment
		if err := configurator.ConfigureDeleteStub(roleAssignUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(roleAssignUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete roles
	for _, role := range roles {
		roleUrl := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, role.Metadata.Tenant, role.Metadata.Name)

		// Delete the role assignment
		if err := configurator.ConfigureDeleteStub(roleUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(roleUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
