package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

func CreateWorkspaceLifecycleScenarioV1(scenario string, params *WorkspaceParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	url := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	resource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)

	response := newWorkspaceResponse(params.Workspace.Name, secalib.WorkspaceProviderV1, resource, secalib.ApiVersion1,
		params.Tenant, params.Region,
		params.Workspace.InitialLabels)

	// Create a workspace
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = newWorkspaceStatus(secalib.CreatingStatusState)
	response.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: url, params: params, response: response, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get created workspace
	setWorkspaceStatusState(response.Status, secalib.ActiveStatusState)
	response.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: url, params: params, response: response, currentState: "GetCreatedWorkspace", nextState: "UpdateWorkspace"}); err != nil {
		return nil, err
	}

	// Update the workspace
	setModifiedRegionalResourceMetadata(response.Metadata)
	setWorkspaceStatusState(response.Status, secalib.UpdatingStatusState)
	response.Labels = params.Workspace.UpdatedLabels
	response.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: url, params: params, response: response, currentState: "UpdateWorkspace", nextState: "GetUpdatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get updated workspace
	setWorkspaceStatusState(response.Status, secalib.ActiveStatusState)
	response.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: url, params: params, response: response, currentState: "GetUpdatedWorkspace", nextState: "DeleteWorkspace"}); err != nil {
		return nil, err
	}

	// Delete the workspace
	response.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: url, params: params, response: response, currentState: "DeleteWorkspace", nextState: "GetDeletedWorkspace"}); err != nil {
		return nil, err
	}

	// Get deleted workspace
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: url, params: params, response: response, currentState: "GetDeletedWorkspace", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}
