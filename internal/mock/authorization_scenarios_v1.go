package mock

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

type AuthorizationScenariosV1 struct {
	Scenarios
}

func NewAuthorizationScenariosV1(authToken string, tenant string, mockURL string) *AuthorizationScenariosV1 {
	return &AuthorizationScenariosV1{
		Scenarios: Scenarios{
			params: secalib.GeneralParams{
				AuthToken: authToken,
				Tenant:    tenant,
			},
			mockURL: mockURL,
		},
	}
}

func (scenarios *AuthorizationScenariosV1) ConfigureLifecycleScenario(id string, params *secalib.AuthorizationLifeCycleParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to Authorization Lifecycle Scenario")

	name := "AuthorizationLifecycleV1_" + id

	wm, err := scenarios.newClient()
	if err != nil {
		return nil, err
	}

	roleUrl := secalib.GenerateRoleURL(scenarios.params.Tenant, params.Role.Name)
	roleAssignmentUrl := secalib.GenerateRoleAssignmentURL(scenarios.params.Tenant, params.RoleAssignment.Name)

	roleResource := secalib.GenerateRoleResource(scenarios.params.Tenant, params.Role.Name)
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(scenarios.params.Tenant, params.RoleAssignment.Name)

	// Role
	roleResponse := &resourceResponse[secalib.RoleSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Role.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleKind,
			Tenant:     scenarios.params.Tenant,
		},
		Status: &secalib.Status{},
		Spec:   &secalib.RoleSpecV1{},
	}
	for _, perm := range params.Role.InitialSpec.Permissions {
		roleResponse.Spec.Permissions = append(roleResponse.Spec.Permissions, &secalib.RoleSpecPermissionV1{
			Provider:  perm.Provider,
			Resources: append([]string{}, perm.Resources...),
			Verb:      append([]string{}, perm.Verb...),
		})
	}

	// Create a role
	roleResponse.Metadata.Verb = http.MethodPut
	roleResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	roleResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	roleResponse.Metadata.ResourceVersion = 1
	roleResponse.Status.State = secalib.CreatingStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          roleUrl,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          roleUrl,
		params:       scenarios.params,
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
	roleResponse.Spec = params.Role.UpdatedSpec
	roleResponse.Status.State = secalib.UpdatingStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          roleUrl,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          roleUrl,
		params:       scenarios.params,
		response:     roleResponse,
		template:     roleResponseTemplateV1,
		currentState: "GetUpdatedRole",
		nextState:    "CreateRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignmentResponse := &resourceResponse[secalib.RoleAssignmentSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.RoleAssignment.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleAssignmentResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleAssignmentKind,
			Tenant:     scenarios.params.Tenant,
		},
		Status: &secalib.Status{},
		Spec: &secalib.RoleAssignmentSpecV1{
			Subs:  params.RoleAssignment.InitialSpec.Subs,
			Roles: params.RoleAssignment.InitialSpec.Roles,
		},
	}
	for _, scope := range params.RoleAssignment.InitialSpec.Scopes {
		roleAssignmentResponse.Spec.Scopes = append(roleAssignmentResponse.Spec.Scopes, &secalib.RoleAssignmentSpecScopeV1{
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
	if err := configurePutStub(wm, name, stubConfig{
		url:          roleAssignmentUrl,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          roleAssignmentUrl,
		params:       scenarios.params,
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
	roleAssignmentResponse.Spec = params.RoleAssignment.UpdatedSpec
	roleAssignmentResponse.Status.State = secalib.UpdatingStatusState
	roleAssignmentResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          roleAssignmentUrl,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          roleAssignmentUrl,
		params:       scenarios.params,
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
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          roleAssignmentUrl,
		params:       scenarios.params,
		response:     roleAssignmentResponse,
		currentState: "DeleteRoleAssignment",
		nextState:    "GetDeletedRoleAssignment",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          roleAssignmentUrl,
		params:       scenarios.params,
		response:     roleAssignmentResponse,
		currentState: "GetDeletedRoleAssignment",
		nextState:    "DeleteRole",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the role
	roleResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          roleUrl,
		params:       scenarios.params,
		response:     roleResponse,
		currentState: "DeleteRole",
		nextState:    "GetDeletedRole",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted role
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          roleUrl,
		params:       scenarios.params,
		response:     roleResponse,
		currentState: "GetDeletedRole",
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
