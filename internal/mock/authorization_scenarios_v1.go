package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	authorization "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.authorization.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
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
		Build()
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
		Build()
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
	if err := configurator.configureDeleteStub(roleAssignUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted role assignment
	if err := configurator.configureGetNotFoundStub(roleAssignUrl, params); err != nil {
		return nil, err
	}

	// Delete the role
	if err := configurator.configureDeleteStub(roleUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted role
	if err := configurator.configureGetNotFoundStub(roleUrl, params); err != nil {
		return nil, err
	}

	return configurator.client, nil
}

func CreateAuthorizationListLifecycleScenarioV1(scenario string, params *AuthorizationListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	var rolesList []schema.Role

	// Role
	for _, role := range *params.Role {
		roleUrl := generators.GenerateRoleURL(authorizationProviderV1, params.Tenant, role.Name)

		roleResponse, err := builders.NewRoleBuilder().
			Name(role.Name).
			Provider(authorizationProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).
			Spec(role.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a role
		if err := configurator.configureCreateRoleStub(roleResponse, roleUrl, params); err != nil {
			return nil, err
		}
		rolesList = append(rolesList, *roleResponse)
	}

	roleResource := generators.GenerateRoleResource(authorizationProviderV1, params.Tenant)
	rolesResponse := &authorization.RoleIterator{
		Metadata: schema.ResponseMetadata{
			Provider: authorizationProviderV1,
			Resource: roleResource,
			Verb:     http.MethodGet,
		},
	}

	rolesResponse.Items = rolesList
	// List Roles
	if err := configurator.configureGetListRoleStub(rolesResponse, generators.GenerateRoleListURL(authorizationProviderV1, params.Tenant), params, nil); err != nil {
		return nil, err
	}

	// List Roles with limit 1

	rolesResponse.Items = rolesList[:1]
	if err := configurator.configureGetListRoleStub(rolesResponse, generators.GenerateRoleListURL(authorizationProviderV1, params.Tenant), params, pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListRoleStub(rolesResponse, generators.GenerateRoleListURL(authorizationProviderV1, params.Tenant), params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}
	// List roles with limit and label

	rolesResponse.Items = rolesListWithLabel(rolesList)[:1]
	if err := configurator.configureGetListRoleStub(rolesResponse, generators.GenerateRoleListURL(authorizationProviderV1, params.Tenant), params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// RoleAssignment

	var rolesAssignmentList []schema.RoleAssignment

	for _, roleAssignment := range *params.RoleAssignment {
		roleAssignmentUrl := generators.GenerateRoleAssignmentURL(authorizationProviderV1, params.Tenant, roleAssignment.Name)
		roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
			Name(roleAssignment.Name).
			Provider(authorizationProviderV1).
			ApiVersion(apiVersion1).
			Tenant(params.Tenant).
			Labels(roleAssignment.InitialLabels).
			Spec(roleAssignment.InitialSpec).
			Build()

		if err != nil {
			return nil, err
		}

		// Create a role assignment
		if err := configurator.configureCreateRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, params); err != nil {
			return nil, err
		}

		rolesAssignmentList = append(rolesAssignmentList, *roleAssignResponse)
	}

	roleAssignUrl := generators.GenerateRoleAssignmentListURL(authorizationProviderV1, params.Tenant)
	roleAssignResource := generators.GenerateRoleAssignmentListResource(params.Tenant)
	roleAssignResponse := &authorization.RoleAssignmentIterator{
		Metadata: schema.ResponseMetadata{
			Provider: authorizationProviderV1,
			Resource: roleAssignResource,
			Verb:     http.MethodGet,
		},
		Items: rolesAssignmentList,
	}

	// List RoleAssignments
	if err := configurator.configureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params, nil); err != nil {
		return nil, err
	}
	// List Roles with limit 1

	roleAssignResponse.Items = rolesAssignmentList[:1]
	if err := configurator.configureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params, pathParamsLimit("1")); err != nil {
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
	roleAssignResponse.Items = rolesAssignWithLabel(rolesAssignmentList)
	if err := configurator.configureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}
	// List roles with limit and label

	roleAssignResponse.Items = rolesAssignWithLabel(rolesAssignmentList)[:1]
	if err := configurator.configureGetListRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Delete RoleAssignments
	for _, roleAssignment := range rolesAssignmentList {
		roleAssignUrl := generators.GenerateRoleAssignmentURL(authorizationProviderV1, params.Tenant, roleAssignment.Metadata.Name)

		// Delete the role assignment
		if err := configurator.configureDeleteStub(roleAssignUrl, params, false); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.configureGetNotFoundStub(roleAssignUrl, params, true); err != nil {
			return nil, err
		}
	}

	// Delete Roles
	for _, role := range rolesList {
		roleUrl := generators.GenerateRoleURL(authorizationProviderV1, params.Tenant, role.Metadata.Name)

		// Delete the role assignment
		if err := configurator.configureDeleteStub(roleUrl, params, false); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.configureGetNotFoundStub(roleUrl, params, true); err != nil {
			return nil, err
		}
	}

	return configurator.client, nil
}
