package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	authorizationV1 "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.authorization.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func CreateAuthorizationLifecycleScenarioV1(scenario string, params *AuthorizationParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	roleUrl := secalib.GenerateRoleURL(params.Tenant, (*params.Role)[0].Name)
	roleAssignUrl := secalib.GenerateRoleAssignmentURL(params.Tenant, (*params.RoleAssignment)[0].Name)

	roleResource := secalib.GenerateRoleResource(params.Tenant, (*params.Role)[0].Name)
	roleAssignResource := secalib.GenerateRoleAssignmentResource(params.Tenant, (*params.RoleAssignment)[0].Name)
	// Role
	roleResponse, err := builders.NewRoleBuilder().
		Name((*params.Role)[0].Name).
		Provider(secalib.AuthorizationProviderV1).
		Resource(roleResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Spec((*params.Role)[0].InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a role
	setCreatedGlobalTenantResourceMetadata(roleResponse.Metadata)
	roleResponse.Status = secalib.NewResourceStatus(schema.ResourceStateCreating)
	roleResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: startedScenarioState, nextState: "GetCreatedRole"}); err != nil {
		return nil, err
	}

	// Get the created role
	secalib.SetStatusState(roleResponse.Status, schema.ResourceStateActive)
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: "GetCreatedRole", nextState: "UpdateRole"}); err != nil {
		return nil, err
	}

	// Update the role
	setModifiedGlobalTenantResourceMetadata(roleResponse.Metadata)
	secalib.SetStatusState(roleResponse.Status, schema.ResourceStateUpdating)
	roleResponse.Spec = *(*params.Role)[0].UpdatedSpec
	roleResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: "UpdateRole", nextState: "GetUpdatedRole"}); err != nil {
		return nil, err
	}

	// Get the updated role
	secalib.SetStatusState(roleResponse.Status, schema.ResourceStateActive)
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: "GetUpdatedRole", nextState: "CreateRoleAssignment"}); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name((*params.RoleAssignment)[0].Name).
		Provider(secalib.AuthorizationProviderV1).
		Resource(roleAssignResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Spec((*params.RoleAssignment)[0].InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a role assignment
	setCreatedGlobalTenantResourceMetadata(roleAssignResponse.Metadata)
	roleAssignResponse.Status = secalib.NewResourceStatus(schema.ResourceStateCreating)
	roleAssignResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: "CreateRoleAssignment", nextState: "GetCreatedRoleAssignment"}); err != nil {
		return nil, err
	}

	// Get the created role assignment
	secalib.SetStatusState(roleAssignResponse.Status, schema.ResourceStateActive)
	roleAssignResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: "GetCreatedRoleAssignment", nextState: "UpdateRoleAssignment"}); err != nil {
		return nil, err
	}

	// Update the role assignment
	setModifiedGlobalTenantResourceMetadata(roleAssignResponse.Metadata)
	secalib.SetStatusState(roleAssignResponse.Status, schema.ResourceStateUpdating)
	roleAssignResponse.Spec = *(*params.RoleAssignment)[0].UpdatedSpec
	roleAssignResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: "UpdateRoleAssignment", nextState: "GetUpdatedRoleAssignment"}); err != nil {
		return nil, err
	}

	// Get the updated role assignment
	secalib.SetStatusState(roleAssignResponse.Status, schema.ResourceStateActive)
	roleAssignResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: "GetUpdatedRoleAssignment", nextState: "DeleteRoleAssignment"}); err != nil {
		return nil, err
	}

	// Delete the role assignment
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, currentState: "DeleteRoleAssignment", nextState: "GetDeletedRoleAssignment"}); err != nil {
		return nil, err
	}

	// Get the deleted role assignment
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, currentState: "GetDeletedRoleAssignment", nextState: "DeleteRole"}); err != nil {
		return nil, err
	}

	// Delete the role
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, currentState: "DeleteRole", nextState: "GetDeletedRole"}); err != nil {
		return nil, err
	}

	// Get deleted role
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, currentState: "GetDeletedRole", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}

func CreateAuthorizationListLifecycleScenarioV1(scenario string, params *AuthorizationParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}
	var rolesList []schema.Role

	for i := range *params.Role {
		roleResource := secalib.GenerateRoleResource(params.Tenant, (*params.Role)[i].Name)
		// Role

		roleResponse, err := builders.NewRoleBuilder().
			Name((*params.Role)[i].Name).
			Provider(secalib.AuthorizationProviderV1).
			Resource(roleResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Label((*params.Role)[i].InitialLabels).
			Spec((*params.Role)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}

		// Create a role
		setCreatedGlobalTenantResourceMetadata(roleResponse.Metadata)
		roleResponse.Status = secalib.NewResourceStatus(schema.ResourceStateCreating)
		roleResponse.Metadata.Verb = http.MethodPut
		roleResponse.Labels = (*params.Role)[i].InitialLabels
		var nextState string
		if i < len(*params.Role)-1 {
			nextState = (*params.RoleAssignment)[i+1].Name
		} else {
			nextState = "GetRoleList"
		}

		if err := configurePutStub(wm, scenario,
			&stubConfig{
				url:          secalib.GenerateRoleURL(params.Tenant, (*params.Role)[i].Name),
				params:       params,
				responseBody: roleResponse,
				currentState: func() string {
					if i == 0 {
						return startedScenarioState
					}
					return (*params.RoleAssignment)[i].Name
				}(),
				nextState: nextState,
			}); err != nil {
			return nil, err
		}

		// Create Role to be listed

		rolesList = append(rolesList, *roleResponse)
	}
	roleResource := secalib.GenerateRolesResource(params.Tenant)
	rolesResponse := &authorizationV1.RoleIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.RegionProviderV1,
			Resource: roleResource,
			Verb:     http.MethodGet,
		},
	}
	rolesResponse.Items = rolesList
	// List Roles
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleListURL(params.Tenant), params: params, responseBody: rolesResponse, currentState: "GetRoleList", nextState: "GetRolesListWithLimit"}); err != nil {
		return nil, err
	}

	// List Roles with limit 1

	rolesResponse.Items = rolesList[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleListURL(params.Tenant), params: params, pathParams: pathParamsLimit("1"), responseBody: rolesResponse, currentState: "GetRolesListWithLimit", nextState: "GetRolesListWithLabel"}); err != nil {
		return nil, err
	}
	// List roles with label
	rolesListWithLabel := func(rolesList []schema.Role) []schema.Role {
		var filteredRoles []schema.Role
		for _, role := range rolesList {
			if val, ok := role.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformance {
				filteredRoles = append(filteredRoles, role)
			}
		}
		return filteredRoles
	}
	rolesResponse.Items = rolesListWithLabel(rolesList)
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleListURL(params.Tenant), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformance), responseBody: rolesResponse, currentState: "GetRolesListWithLabel", nextState: "GetRolesListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List roles with limit and label

	rolesResponse.Items = rolesListWithLabel(rolesList)[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleListURL(params.Tenant), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformance), responseBody: rolesResponse, currentState: "GetRolesListWithLimitAndLabel", nextState: (*params.RoleAssignment)[0].Name}); err != nil {
		return nil, err
	}

	// RoleAssignment

	var rolesAssignmentList []schema.RoleAssignment
	for i := range *params.RoleAssignment {
		roleAssignResource := secalib.GenerateRoleAssignmentResource(params.Tenant, (*params.RoleAssignment)[i].Name)
		// RoleAssignmentroleAssignResponse := newRoleAssignmentResponse((*params.RoleAssignment)[i].Name, secalib.AuthorizationProviderV1, roleAssignResource, secalib.ApiVersion1, params.Tenant, (*params.RoleAssignment)[i].InitialSpec)
		roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
			Name((*params.RoleAssignment)[0].Name).
			Provider(secalib.AuthorizationProviderV1).
			Resource(roleAssignResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Labels((*params.RoleAssignment)[i].InitialLabels).
			Spec((*params.RoleAssignment)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}
		// Create a role assignment
		setCreatedGlobalTenantResourceMetadata(roleAssignResponse.Metadata)
		roleAssignResponse.Status = secalib.NewResourceStatus(schema.ResourceStateCreating)
		roleAssignResponse.Metadata.Verb = http.MethodPut
		var nextState string
		if i < len(*params.RoleAssignment)-1 {
			nextState = (*params.RoleAssignment)[i+1].Name
		} else {
			nextState = "GetRoleAssignmentsList"
		}
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateRoleAssignmentURL(params.Tenant, (*params.RoleAssignment)[i].Name), params: params, responseBody: roleAssignResponse, currentState: (*params.RoleAssignment)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}

		// Create RoleAssignment to be listed

		rolesAssignmentList = append(rolesAssignmentList, *roleAssignResponse)
	}
	roleAssignResource := secalib.GenerateRoleAssignmentsResource(params.Tenant)
	roleAssignResponse := &authorizationV1.RoleAssignmentIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.RegionProviderV1,
			Resource: roleAssignResource,
			Verb:     http.MethodGet,
		},
		Items: rolesAssignmentList,
	}

	// List RoleAssignments
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleAssignmentListURL(params.Tenant), params: params, responseBody: roleAssignResponse, currentState: "GetRoleAssignmentsList", nextState: "GetRoleAssignmentsListWithLimit"}); err != nil {
		return nil, err
	}
	// List Roles with limit 1

	roleAssignResponse.Items = rolesAssignmentList[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleAssignmentListURL(params.Tenant), params: params, pathParams: pathParamsLimit("1"), responseBody: roleAssignResponse, currentState: "GetRoleAssignmentsListWithLimit", nextState: "GetRoleAssignmentsListWithLabel"}); err != nil {
		return nil, err
	}
	// List roles with label

	rolesAssignWithLabel := func(rolesAssignList []schema.RoleAssignment) []schema.RoleAssignment {
		var filteredRoles []schema.RoleAssignment
		for _, roleAssign := range rolesAssignList {
			if val, ok := roleAssign.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformance {
				filteredRoles = append(filteredRoles, roleAssign)
			}
		}
		return filteredRoles
	}
	roleAssignResponse.Items = rolesAssignWithLabel(rolesAssignmentList)
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleAssignmentListURL(params.Tenant), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformance), responseBody: roleAssignResponse, currentState: "GetRoleAssignmentsListWithLabel", nextState: "GetRoleAssignmentsListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List roles with limit and label

	roleAssignResponse.Items = rolesAssignWithLabel(rolesAssignmentList)[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleAssignmentListURL(params.Tenant), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformance), responseBody: roleAssignResponse, currentState: "GetRoleAssignmentsListWithLimitAndLabel", nextState: "DeleteRoleAssignment"}); err != nil {
		return nil, err
	}

	// Delete RoleAssignments
	for i := range *params.RoleAssignment {
		roleAssignUrl := secalib.GenerateRoleAssignmentURL(params.Tenant, (*params.RoleAssignment)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteRoleAssignment"
		} else {
			currentState = "GetDeletedRoleAssignment_" + (*params.RoleAssignment)[i-1].Name
		}

		nextState = "DeleteRoleAssignment_" + (*params.RoleAssignment)[i].Name

		// Delete the role assignment
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: roleAssignUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted role assignment (should return 404)
		nextState = func() string {
			if i < len(*params.RoleAssignment)-1 {
				return "GetDeletedRoleAssignment_" + (*params.RoleAssignment)[i].Name
			} else {
				return "DeleteRole_" + (*params.Role)[0].Name
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: roleAssignUrl, params: params, currentState: "DeleteRoleAssignment_" + (*params.RoleAssignment)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
	}

	// Delete Roles
	for i := range *params.Role {
		roleUrl := secalib.GenerateRoleURL(params.Tenant, (*params.Role)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteRole_" + (*params.Role)[i].Name
		} else {
			currentState = "GetDeletedRole_" + (*params.Role)[i-1].Name
		}

		nextState = "DeleteRole_" + (*params.Role)[i].Name

		// Delete the role
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: roleUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted role (should return 404)
		nextState = func() string {
			if i < len(*params.Role)-1 {
				return "GetDeletedRole_" + (*params.Role)[i].Name
			} else {
				return startedScenarioState
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: roleUrl, params: params, currentState: "DeleteRole_" + (*params.Role)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
	}

	return wm, nil
}
