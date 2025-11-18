package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/conformance/secalib/builders"
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
	roleResponse, err := builders.NewRoleBuilder().
		Name(params.Role.Name).
		Provider(secalib.AuthorizationProviderV1).
		Resource(roleResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Spec(params.Role.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

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
	roleResponse.Spec = *params.Role.UpdatedSpec
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
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name(params.RoleAssignment.Name).
		Provider(secalib.AuthorizationProviderV1).
		Resource(roleAssignResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Spec(params.RoleAssignment.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

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
	roleAssignResponse.Spec = *params.RoleAssignment.UpdatedSpec
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
