package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigureAuthorizationLifecycleScenarioV1(scenario string, params *AuthorizationLifeCycleParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	roleUrl := generators.GenerateRoleURL(authorizationProviderV1, params.Tenant, params.Role.Name)
	roleAssignUrl := generators.GenerateRoleAssignmentURL(authorizationProviderV1, params.Tenant, params.RoleAssignment.Name)

	// Role
	roleResponse, err := builders.NewRoleBuilder().
		Name(params.Role.Name).
		Provider(authorizationProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).
		Spec(params.Role.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a role
	if err := configurator.configureCreateRoleStub(roleResponse, roleUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created role
	if err := configurator.configureGetActiveRoleStub(roleResponse, roleUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Update the role
	roleResponse.Spec = *params.Role.UpdatedSpec
	if err := configurator.configureUpdateRoleStub(roleResponse, roleUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated role
	if err := configurator.configureGetActiveRoleStub(roleResponse, roleUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name(params.RoleAssignment.Name).
		Provider(authorizationProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).
		Spec(params.RoleAssignment.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a role assignment
	if err := configurator.configureCreateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created role assignment
	if err := configurator.configureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Update the role assignment
	roleAssignResponse.Spec = *params.RoleAssignment.UpdatedSpec
	if err := configurator.configureUpdateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated role assignment
	if err := configurator.configureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the role assignment
	if err := configurator.configureDeleteStub(roleAssignUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted role assignment
	if err := configurator.configureGetNotFoundStub(roleAssignUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the role
	if err := configurator.configureDeleteStub(roleUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted role
	if err := configurator.configureGetNotFoundStub(roleUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	return configurator.client, nil
}

func ConfigureAuthorizationListScenarioV1(scenario string, params *AuthorizationListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	roleUrl := generators.GenerateRoleListURL(authorizationProviderV1, params.Tenant)
	roleAssignUrl := generators.GenerateRoleAssignmentListURL(authorizationProviderV1, params.Tenant)

	// Create roles
	rolesList, err := bulkCreateRolesStubV1(configurator, params.getBaseParams(), params.Roles)
	if err != nil {
		return nil, err
	}
	rolesResponse, err := builders.NewRoleIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).
		Items(rolesList).
		Build()
	if err != nil {
		return nil, err
	}

	// List Roles
	if err := configurator.configureGetListRoleStub(rolesResponse, roleUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List Roles with limit 1
	rolesResponse.Items = rolesList[:1]
	if err := configurator.configureGetListRoleStub(rolesResponse, roleUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListRoleStub(rolesResponse, roleUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List roles with limit and label
	rolesResponse.Items = rolesListWithLabel(rolesList)[:1]
	if err := configurator.configureGetListRoleStub(rolesResponse, roleUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create role assignments
	roleAssignmentsList, err := bulkCreateRoleAssignmentsStubV1(configurator, params.getBaseParams(), params.RoleAssignments)
	if err != nil {
		return nil, err
	}
	roleAssignResponse, err := builders.NewRoleAssignmentIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).
		Items(roleAssignmentsList).
		Build()
	if err != nil {
		return nil, err
	}

	// List RoleAssignments
	if err := configurator.configureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List Roles with limit 1
	roleAssignResponse.Items = roleAssignmentsList[:1]
	if err := configurator.configureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List roles with limit and label
	roleAssignResponse.Items = rolesAssignWithLabel(roleAssignmentsList)[:1]
	if err := configurator.configureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Delete role assignments
	for _, roleAssignment := range roleAssignmentsList {
		roleAssignUrl := generators.GenerateRoleAssignmentURL(authorizationProviderV1, params.Tenant, roleAssignment.Metadata.Name)

		// Delete the role assignment
		if err := configurator.configureDeleteStub(roleAssignUrl, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.configureGetNotFoundStub(roleAssignUrl, params.getBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete roles
	for _, role := range rolesList {
		roleUrl := generators.GenerateRoleURL(authorizationProviderV1, params.Tenant, role.Metadata.Name)

		// Delete the role assignment
		if err := configurator.configureDeleteStub(roleUrl, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.configureGetNotFoundStub(roleUrl, params.getBaseParams()); err != nil {
			return nil, err
		}
	}

	return configurator.client, nil
}
