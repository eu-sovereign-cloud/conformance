package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
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
	roleResponse := newRoleResponse((*params.Role)[0].Name, secalib.AuthorizationProviderV1, roleResource, secalib.ApiVersion1, params.Tenant, (*params.Role)[0].InitialSpec)

	// Create a role
	setCreatedGlobalTenantResourceMetadata(roleResponse.Metadata)
	roleResponse.Status = secalib.NewResourceStatus(secalib.CreatingResourceState)
	roleResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: startedScenarioState, nextState: "GetCreatedRole"}); err != nil {
		return nil, err
	}

	// Get the created role
	secalib.SetStatusState(roleResponse.Status, secalib.ActiveResourceState)
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: "GetCreatedRole", nextState: "UpdateRole"}); err != nil {
		return nil, err
	}

	// Update the role
	setModifiedGlobalTenantResourceMetadata(roleResponse.Metadata)
	secalib.SetStatusState(roleResponse.Status, secalib.UpdatingResourceState)
	roleResponse.Spec = *(*params.Role)[0].UpdatedSpec
	roleResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: "UpdateRole", nextState: "GetUpdatedRole"}); err != nil {
		return nil, err
	}

	// Get the updated role
	secalib.SetStatusState(roleResponse.Status, secalib.ActiveResourceState)
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: "GetUpdatedRole", nextState: "CreateRoleAssignment"}); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse := newRoleAssignmentResponse((*params.RoleAssignment)[0].Name, secalib.AuthorizationProviderV1, roleAssignResource, secalib.ApiVersion1,
		params.Tenant,
		(*params.RoleAssignment)[0].InitialSpec)

	// Create a role assignment
	setCreatedGlobalTenantResourceMetadata(roleAssignResponse.Metadata)
	roleAssignResponse.Status = secalib.NewResourceStatus(secalib.CreatingResourceState)
	roleAssignResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: "CreateRoleAssignment", nextState: "GetCreatedRoleAssignment"}); err != nil {
		return nil, err
	}

	// Get the created role assignment
	secalib.SetStatusState(roleAssignResponse.Status, secalib.ActiveResourceState)
	roleAssignResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: "GetCreatedRoleAssignment", nextState: "UpdateRoleAssignment"}); err != nil {
		return nil, err
	}

	// Update the role assignment
	setModifiedGlobalTenantResourceMetadata(roleAssignResponse.Metadata)
	secalib.SetStatusState(roleAssignResponse.Status, secalib.UpdatingResourceState)
	roleAssignResponse.Spec = *(*params.RoleAssignment)[0].UpdatedSpec
	roleAssignResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: "UpdateRoleAssignment", nextState: "GetUpdatedRoleAssignment"}); err != nil {
		return nil, err
	}

	// Get the updated role assignment
	secalib.SetStatusState(roleAssignResponse.Status, secalib.ActiveResourceState)
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
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: roleAssignUrl, params: params, currentState: "GetDeletedRoleAssignment", nextState: "DeleteRole"}); err != nil {
		return nil, err
	}

	// Delete the role
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, currentState: "DeleteRole", nextState: "GetDeletedRole"}); err != nil {
		return nil, err
	}

	// Get deleted role
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
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
		roleUrl := secalib.GenerateRoleURL(params.Tenant, (*params.Role)[i].Name)
		// roleAssignUrl := secalib.GenerateRoleAssignmentURL(params.Tenant, (*params.RoleAssignment)[i].Name)

		roleResource := secalib.GenerateRoleResource(params.Tenant, (*params.Role)[i].Name)
		// roleAssignResource := secalib.GenerateRoleAssignmentResource(params.Tenant, (*params.RoleAssignment)[i].Name)
		// Role
		roleResponse := newRoleResponse((*params.Role)[i].Name, secalib.AuthorizationProviderV1, roleResource, secalib.ApiVersion1, params.Tenant, (*params.Role)[i].InitialSpec)

		// Create a role
		setCreatedGlobalTenantResourceMetadata(roleResponse.Metadata)
		roleResponse.Status = secalib.NewResourceStatus(secalib.CreatingResourceState)
		roleResponse.Metadata.Verb = http.MethodPut
		var nextState string
		if i < len(*params.Role)-1 {
			nextState = (*params.RoleAssignment)[i+1].Name
		} else {
			nextState = "GetRoleList"
		}

		if err := configurePutStub(wm, scenario,
			&stubConfig{
				url:          roleUrl,
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
		&stubConfig{url: secalib.GenerateRolesURL(params.Tenant), params: params, responseBody: rolesResponse, currentState: "GetRoleList", nextState: "GetRolesListWithLimit"}); err != nil {
		return nil, err
	}
	// List Roles with limit 1
	rolesWithLimitResponse := &authorizationV1.RoleIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.RegionProviderV1,
			Resource: roleResource,
			Verb:     http.MethodGet,
		},
	}
	rolesWithLimitResponse.Items = rolesList[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRolesURL(params.Tenant), params: params, responseBody: rolesWithLimitResponse, currentState: "GetRolesListWithLimit", nextState: "GetRolesListWithLabel"}); err != nil {
		return nil, err
	}
	// List roles with label
	rolesWithLabelResponse := &authorizationV1.RoleIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.RegionProviderV1,
			Resource: roleResource,
			Verb:     http.MethodGet,
		},
	}
	rolesListWithLabel := func(rolesList []schema.Role) []schema.Role {
		var filteredRoles []schema.Role
		for _, role := range rolesList {
			if val, ok := role.Labels["env"]; ok && val == "conformance" {
				filteredRoles = append(filteredRoles, role)
			}
		}
		return filteredRoles
	}
	rolesWithLabelResponse.Items = rolesListWithLabel(rolesList)
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRolesURL(params.Tenant), params: params, responseBody: rolesWithLabelResponse, currentState: "GetRolesListWithLabel", nextState: "GetRolesListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List roles with limit and label
	rolesWithLimitAndLabelResponse := &authorizationV1.RoleIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.RegionProviderV1,
			Resource: roleResource,
			Verb:     http.MethodGet,
		},
	}
	rolesWithLimitAndLabelResponse.Items = rolesListWithLabel(rolesList)[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRolesURL(params.Tenant), params: params, responseBody: rolesWithLimitAndLabelResponse, currentState: "GetRolesListWithLimitAndLabel", nextState: (*params.RoleAssignment)[0].Name}); err != nil {
		return nil, err
	}

	// RoleAgreement
	var rolesAssignmentList []schema.RoleAssignment
	for i := range *params.RoleAssignment {

		roleAssignUrl := secalib.GenerateRoleAssignmentURL(params.Tenant, (*params.RoleAssignment)[i].Name)

		roleAssignResource := secalib.GenerateRoleAssignmentResource(params.Tenant, (*params.RoleAssignment)[i].Name)
		// RoleAssignment
		roleAssignResponse := newRoleAssignmentResponse((*params.RoleAssignment)[i].Name, secalib.AuthorizationProviderV1, roleAssignResource, secalib.ApiVersion1, params.Tenant, (*params.RoleAssignment)[i].InitialSpec)

		// Create a role assignment
		setCreatedGlobalTenantResourceMetadata(roleAssignResponse.Metadata)
		roleAssignResponse.Status = secalib.NewResourceStatus(secalib.CreatingResourceState)
		roleAssignResponse.Metadata.Verb = http.MethodPut
		var nextState string
		if i < len(*params.RoleAssignment)-1 {
			nextState = (*params.RoleAssignment)[i+1].Name
		} else {
			nextState = "GetRoleAssignmentsList"
		}
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: (*params.RoleAssignment)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}

		// Create Role to be listed

		rolesAssignmentList = append(rolesAssignmentList, *roleAssignResponse)
	}
	roleAssignResource := secalib.GenerateRoleAssignmentsResource(params.Tenant)
	roleAssignResponse := &authorizationV1.RoleAssignmentIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.RegionProviderV1,
			Resource: roleAssignResource,
			Verb:     http.MethodGet,
		},
	}
	roleAssignResponse.Items = rolesAssignmentList
	// List RoleAssignments
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleAssignmentsURL(params.Tenant), params: params, responseBody: roleAssignResponse, currentState: "GetRoleAssignmentsList", nextState: "GetRoleAssignmentsListWithLimit"}); err != nil {
		return nil, err
	}
	// List Roles with limit 1
	rolesAssignWithLimitResponse := &authorizationV1.RoleAssignmentIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.RegionProviderV1,
			Resource: roleAssignResource,
			Verb:     http.MethodGet,
		},
	}
	rolesAssignWithLimitResponse.Items = rolesAssignmentList[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleAssignmentsURL(params.Tenant), params: params, responseBody: rolesAssignWithLimitResponse, currentState: "GetRoleAssignmentsListWithLimit", nextState: "GetRoleAssignmentsListWithLabel"}); err != nil {
		return nil, err
	}
	// List roles with label
	rolesAssignWithLabelResponse := &authorizationV1.RoleAssignmentIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.RegionProviderV1,
			Resource: roleAssignResource,
			Verb:     http.MethodGet,
		},
	}
	rolesAssignWithLabel := func(rolesAssignList []schema.RoleAssignment) []schema.RoleAssignment {
		var filteredRoles []schema.RoleAssignment
		for _, roleAssign := range rolesAssignList {
			if val, ok := roleAssign.Labels["env"]; ok && val == "conformance" {
				filteredRoles = append(filteredRoles, roleAssign)
			}
		}
		return filteredRoles
	}
	rolesAssignWithLabelResponse.Items = rolesAssignWithLabel(rolesAssignmentList)
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleAssignmentsURL(params.Tenant), params: params, responseBody: rolesAssignWithLabelResponse, currentState: "GetRoleAssignmentsListWithLabel", nextState: "GetRoleAssignmentsListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List roles with limit and label
	rolesAssignWithLimitAndLabelResponse := &authorizationV1.RoleAssignmentIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.RegionProviderV1,
			Resource: roleAssignResource,
			Verb:     http.MethodGet,
		},
	}
	rolesAssignWithLimitAndLabelResponse.Items = rolesAssignWithLabel(rolesAssignmentList)[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateRoleAssignmentsURL(params.Tenant), params: params, responseBody: rolesAssignWithLimitAndLabelResponse, currentState: "GetRoleAssignmentsListWithLimitAndLabel", nextState: startedScenarioState}); err != nil {
		return nil, err
	}
	return wm, nil
}
