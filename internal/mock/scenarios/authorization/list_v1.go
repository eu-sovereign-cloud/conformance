package authorization

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/wiremock/go-wiremock"
)

func ConfigureListScenarioV1(scenario string, params *mock.AuthorizationListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := stubs.NewScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	roleUrl := generators.GenerateRoleListURL(mock.AuthorizationProviderV1, params.Tenant)
	roleAssignUrl := generators.GenerateRoleAssignmentListURL(mock.AuthorizationProviderV1, params.Tenant)

	// Create roles
	rolesList, err := stubs.BulkCreateRolesStubV1(configurator, params.GetBaseParams(), params.Roles)
	if err != nil {
		return nil, err
	}
	rolesResponse, err := builders.NewRoleIteratorBuilder().
		Provider(mock.StorageProviderV1).
		Tenant(params.Tenant).
		Items(rolesList).
		Build()
	if err != nil {
		return nil, err
	}

	// List Roles
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List Roles with limit 1
	rolesResponse.Items = rolesList[:1]
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List roles with label
	rolesListWithLabel := func(rolesList []schema.Role) []schema.Role {
		var filteredRoles []schema.Role
		for _, role := range rolesList {
			if val, ok := role.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredRoles = append(filteredRoles, role)
			}
		}
		return filteredRoles
	}
	rolesResponse.Items = rolesListWithLabel(rolesList)
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, params.GetBaseParams(), mock.PathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List roles with limit and label
	rolesResponse.Items = rolesListWithLabel(rolesList)[:1]
	if err := configurator.ConfigureGetListRoleStub(rolesResponse, roleUrl, params.GetBaseParams(), mock.PathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create role assignments
	roleAssignmentsList, err := stubs.BulkCreateRoleAssignmentsStubV1(configurator, params.GetBaseParams(), params.RoleAssignments)
	if err != nil {
		return nil, err
	}
	roleAssignResponse, err := builders.NewRoleAssignmentIteratorBuilder().
		Provider(mock.StorageProviderV1).
		Tenant(params.Tenant).
		Items(roleAssignmentsList).
		Build()
	if err != nil {
		return nil, err
	}

	// List RoleAssignments
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List Roles with limit 1
	roleAssignResponse.Items = roleAssignmentsList[:1]
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List roles with label
	rolesAssignWithLabel := func(rolesAssignList []schema.RoleAssignment) []schema.RoleAssignment {
		var filteredRoles []schema.RoleAssignment
		for _, roleAssign := range rolesAssignList {
			if val, ok := roleAssign.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredRoles = append(filteredRoles, roleAssign)
			}
		}
		return filteredRoles
	}
	roleAssignResponse.Items = rolesAssignWithLabel(roleAssignmentsList)
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.GetBaseParams(), mock.PathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List roles with limit and label
	roleAssignResponse.Items = rolesAssignWithLabel(roleAssignmentsList)[:1]
	if err := configurator.ConfigureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Delete role assignments
	for _, roleAssignment := range roleAssignmentsList {
		roleAssignUrl := generators.GenerateRoleAssignmentURL(mock.AuthorizationProviderV1, params.Tenant, roleAssignment.Metadata.Name)

		// Delete the role assignment
		if err := configurator.ConfigureDeleteStub(roleAssignUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(roleAssignUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete roles
	for _, role := range rolesList {
		roleUrl := generators.GenerateRoleURL(mock.AuthorizationProviderV1, params.Tenant, role.Metadata.Name)

		// Delete the role assignment
		if err := configurator.ConfigureDeleteStub(roleUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(roleUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}

	return configurator.Client, nil
}
