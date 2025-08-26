package mock

import (
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

func CreateAuthorizationLifecycleScenarioV1(scenario string, params AuthorizationParamsV1) (*wiremock.Client, error) {
	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	roleUrl := secalib.GenerateRoleURL(params.Tenant, params.Role.Name)
	roleAssignmentUrl := secalib.GenerateRoleAssignmentURL(params.Tenant, params.RoleAssignment.Name)

	roleResource := secalib.GenerateRoleResource(params.Tenant, params.Role.Name)
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(params.Tenant, params.RoleAssignment.Name)

	// Role
	roleResponse := roleResponseV1{
		Metadata: metadataResponse{
			Name:       params.Role.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleKind,
			Tenant:     params.Tenant,
		},
	}
	for _, perm := range params.Role.Permissions {
		roleResponse.Permissions = append(roleResponse.Permissions, rolePermissionResponseV1{
			Provider:  perm.Provider,
			Resources: append([]string{}, perm.Resources...),
			Verb:      append([]string{}, perm.VerbInitial...),
		})
	}

	// Create a role
	roleResponse.Metadata.Verb = http.MethodPut
	roleResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	roleResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	roleResponse.Metadata.ResourceVersion = 1
	roleResponse.Status.State = secalib.CreatingStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     roleResponse,
		template:     roleResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    "GetCreatedRole",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created role
	roleResponse.Metadata.Verb = http.MethodGet
	roleResponse.Status.State = secalib.ActiveStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     roleResponse,
		template:     roleResponseTemplateV1,
		currentState: "GetCreatedRole",
		nextState:    "UpdateRole",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the role
	roleResponse.Metadata.Verb = http.MethodPut
	roleResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	roleResponse.Metadata.ResourceVersion = roleResponse.Metadata.ResourceVersion + 1
	for i, perm := range params.Role.Permissions {
		roleResponse.Permissions[i].Verb = append([]string{}, perm.VerbUpdated...)
	}
	roleResponse.Status.State = secalib.UpdatingStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     roleResponse,
		template:     roleResponseTemplateV1,
		currentState: "UpdateRole",
		nextState:    "GetUpdatedRole",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get updated role
	roleResponse.Metadata.Verb = http.MethodGet
	roleResponse.Status.State = secalib.ActiveStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     roleResponse,
		template:     roleResponseTemplateV1,
		currentState: "GetUpdatedRole",
		nextState:    "CreateRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignmentResponse := roleAssignmentResponseV1{
		Metadata: metadataResponse{
			Name:       params.RoleAssignment.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleAssignmentResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleAssignmentKind,
			Tenant:     params.Tenant,
		},
		Subs:  params.RoleAssignment.SubsInitial,
		Roles: params.RoleAssignment.Roles,
	}
	for _, scope := range params.RoleAssignment.Scopes {
		roleAssignmentResponse.Scopes = append(roleAssignmentResponse.Scopes, roleAssignmentScopeResponseV1{
			Tenants:    scope.Tenants,
			Regions:    scope.Regions,
			Workspaces: scope.Workspaces,
		})
	}

	// Create a role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodPut
	roleAssignmentResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	roleAssignmentResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	roleAssignmentResponse.Metadata.ResourceVersion = 1
	roleAssignmentResponse.Status.State = secalib.CreatingStatusState
	roleAssignmentResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     roleAssignmentResponse,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "CreateRoleAssignment",
		nextState:    "GetCreatedRoleAssignment",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodGet
	roleAssignmentResponse.Status.State = secalib.ActiveStatusState
	roleAssignmentResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     roleAssignmentResponse,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "GetCreatedRoleAssignment",
		nextState:    "UpdateRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodPut
	roleAssignmentResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	roleAssignmentResponse.Metadata.ResourceVersion = roleAssignmentResponse.Metadata.ResourceVersion + 1
	roleAssignmentResponse.Subs = params.RoleAssignment.SubsUpdated
	roleAssignmentResponse.Status.State = secalib.UpdatingStatusState
	roleAssignmentResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     roleAssignmentResponse,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "UpdateRoleAssignment",
		nextState:    "GetUpdatedRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get updated role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodGet
	roleAssignmentResponse.Status.State = secalib.ActiveStatusState
	roleAssignmentResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     roleAssignmentResponse,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "GetUpdatedRoleAssignment",
		nextState:    "DeleteRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Delete the role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     roleAssignmentResponse,
		currentState: "DeleteRoleAssignment",
		nextState:    "GetDeletedRoleAssignment",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     roleAssignmentResponse,
		currentState: "GetDeletedRoleAssignment",
		nextState:    "DeleteRole",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the role
	roleResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     roleResponse,
		currentState: "DeleteRole",
		nextState:    "GetDeletedRole",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted role
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     roleResponse,
		currentState: "GetDeletedRole",
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
