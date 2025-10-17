package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

func CreateAuthorizationLifecycleScenarioV1(scenario string, params *AuthorizationParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	roleUrl := secalib.GenerateRoleURL(params.Tenant, params.Role.Name)
	roleAssignUrl := secalib.GenerateRoleAssignmentURL(params.Tenant, params.RoleAssignment.Name)

	roleResource := secalib.GenerateRoleResource(params.Tenant, params.Role.Name)
	roleAssignResource := secalib.GenerateRoleAssignmentResource(params.Tenant, params.RoleAssignment.Name)

	// Role
	roleResponse := newRoleResponse(params.Role.Name, secalib.AuthorizationProviderV1, roleResource, secalib.ApiVersion1,
		params.Tenant,
		params.Role.InitialSpec)

	// Create a role
	setCreatedGlobalTenantResourceMetadata(roleResponse.Metadata)
	roleResponse.Status = newStatus(secalib.CreatingStatusState)
	roleResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, response: roleResponse, currentState: startedScenarioState, nextState: "GetCreatedRole"}); err != nil {
		return nil, err
	}

	// Get created role
	setStatusState(roleResponse.Status, secalib.ActiveStatusState)
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, response: roleResponse, currentState: "GetCreatedRole", nextState: "UpdateRole"}); err != nil {
		return nil, err
	}

	// Update the role
	setModifiedGlobalTenantResourceMetadata(roleResponse.Metadata)
	setStatusState(roleResponse.Status, secalib.UpdatingStatusState)
	roleResponse.Spec = *params.Role.UpdatedSpec
	roleResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, response: roleResponse, currentState: "UpdateRole", nextState: "GetUpdatedRole"}); err != nil {
		return nil, err
	}

	// Get updated role
	setStatusState(roleResponse.Status, secalib.ActiveStatusState)
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, response: roleResponse, currentState: "GetUpdatedRole", nextState: "CreateRoleAssignment"}); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse := newRoleAssignmentResponse(params.RoleAssignment.Name, secalib.AuthorizationProviderV1, roleAssignResource, secalib.ApiVersion1,
		params.Tenant,
		params.RoleAssignment.InitialSpec)

	// Create a role assignment
	setCreatedGlobalTenantResourceMetadata(roleAssignResponse.Metadata)
	roleAssignResponse.Status = newStatus(secalib.CreatingStatusState)
	roleAssignResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, response: roleAssignResponse, currentState: "CreateRoleAssignment", nextState: "GetCreatedRoleAssignment"}); err != nil {
		return nil, err
	}

	// Get created role assignment
	setStatusState(roleAssignResponse.Status, secalib.ActiveStatusState)
	roleAssignResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, response: roleAssignResponse, currentState: "GetCreatedRoleAssignment", nextState: "UpdateRoleAssignment"}); err != nil {
		return nil, err
	}

	// Update the role assignment
	setModifiedGlobalTenantResourceMetadata(roleAssignResponse.Metadata)
	setStatusState(roleAssignResponse.Status, secalib.UpdatingStatusState)
	roleAssignResponse.Spec = *params.RoleAssignment.UpdatedSpec
	roleAssignResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, response: roleAssignResponse, currentState: "UpdateRoleAssignment", nextState: "GetUpdatedRoleAssignment"}); err != nil {
		return nil, err
	}

	// Get updated role assignment
	setStatusState(roleAssignResponse.Status, secalib.ActiveStatusState)
	roleAssignResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, response: roleAssignResponse, currentState: "GetUpdatedRoleAssignment", nextState: "DeleteRoleAssignment"}); err != nil {
		return nil, err
	}

	// Delete the role assignment
	roleAssignResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, response: roleAssignResponse, currentState: "DeleteRoleAssignment", nextState: "GetDeletedRoleAssignment"}); err != nil {
		return nil, err
	}

	// Get deleted role assignment
	roleAssignResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: roleAssignUrl, params: params, response: roleAssignResponse, currentState: "GetDeletedRoleAssignment", nextState: "DeleteRole"}); err != nil {
		return nil, err
	}

	// Delete the role
	roleResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, response: roleResponse, currentState: "DeleteRole", nextState: "GetDeletedRole"}); err != nil {
		return nil, err
	}

	// Get deleted role
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: roleUrl, params: params, response: roleResponse, currentState: "GetDeletedRole", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}
