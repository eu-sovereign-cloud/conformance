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

	roleUrl := secalib.GenerateRoleURL(params.Tenant, params.role.Name)
	roleAssignmentUrl := secalib.GenerateRoleAssignmentURL(params.Tenant, params.roleAssignment.Name)

	roleResource := secalib.GenerateRoleResource(params.Tenant, params.role.Name)
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(params.Tenant, params.roleAssignment.Name)

	// Role
	response := roleResponseV1{
		Metadata: metadataResponse{
			Name:       params.role.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleKind,
			Tenant:     params.Tenant,
		},
	}
	for _, perm := range params.role.Permissions {
		response.Permissions = append(response.Permissions, rolePermissionResponseV1{
			Provider:  perm.Provider,
			Resources: append([]string{}, perm.Resources...),
			Verbs:     append([]string{}, perm.Verbs...),
		})
	}

	// Create a role
	response.Metadata.Verb = http.MethodPut
	response.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = 1
	response.Status.State = secalib.CreatingStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     response,
		template:     roleResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    "GetCreatedRole",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created role
	response.Metadata.Verb = http.MethodGet
	response.Status.State = secalib.ActiveStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     response,
		template:     roleResponseTemplateV1,
		currentState: "GetCreatedRole",
		nextState:    "UpdateRole",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the role
	response.Metadata.Verb = http.MethodPut
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = response.Metadata.ResourceVersion + 1
	response.Status.State = secalib.UpdatingStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     response,
		template:     roleResponseTemplateV1,
		currentState: "UpdateRole",
		nextState:    "GetUpdatedRole",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get updated role
	response.Metadata.Verb = http.MethodGet
	response.Status.State = secalib.ActiveStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     response,
		template:     roleResponseTemplateV1,
		currentState: "GetUpdatedRole",
		nextState:    "DeleteRole",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Delete the role
	response.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     response,
		currentState: "DeleteRole",
		nextState:    "GetDeletedRole",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted role
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleUrl,
		params:       params,
		response:     response,
		currentState: "GetDeletedRole",
		nextState:    "CreateRoleAssignment",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Role assignment
	responseRA := roleAssignmentResponseV1{
		Metadata: metadataResponse{
			Name:       params.roleAssignment.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleAssignmentResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleAssignmentKind,
			Tenant:     params.Tenant,
		},
		Subs:  append([]string{}, params.roleAssignment.subs...),
		Roles: append([]string{}, params.roleAssignment.roles...),
	}
	for _, scope := range params.roleAssignment.scopes {
		responseRA.Scopes = append(responseRA.Scopes, roleAssignmentScopeResponseV1{
			Tenants:    append([]string{}, scope.Tenants...),
			Regions:    append([]string{}, scope.Regions...),
			Workspaces: append([]string{}, scope.Workspaces...),
		})
	}

	// Create a role assignment
	responseRA.Metadata.Verb = http.MethodPut
	responseRA.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	responseRA.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	responseRA.Metadata.ResourceVersion = 1
	responseRA.Status.State = secalib.CreatingStatusState
	responseRA.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     responseRA,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "CreateRoleAssignment",
		nextState:    "GetCreatedRoleAssignment",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created role assignment
	responseRA.Metadata.Verb = http.MethodGet
	responseRA.Status.State = secalib.ActiveStatusState
	responseRA.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     responseRA,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "GetCreatedRoleAssignment",
		nextState:    "UpdateRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the role assignment
	responseRA.Metadata.Verb = http.MethodPut
	responseRA.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	responseRA.Metadata.ResourceVersion = responseRA.Metadata.ResourceVersion + 1
	responseRA.Status.State = secalib.UpdatingStatusState
	responseRA.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     responseRA,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "UpdateRoleAssignment",
		nextState:    "GetUpdatedRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get updated role assignment
	responseRA.Metadata.Verb = http.MethodGet
	responseRA.Status.State = secalib.ActiveStatusState
	responseRA.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     responseRA,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "GetUpdatedRoleAssignment",
		nextState:    "DeleteRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Delete the role assignment
	responseRA.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     responseRA,
		currentState: "DeleteRoleAssignment",
		nextState:    "GetDeletedRoleAssignment",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted role assignment
	responseRA.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          roleAssignmentUrl,
		params:       params,
		response:     responseRA,
		currentState: "GetDeletedRoleAssignment",
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
