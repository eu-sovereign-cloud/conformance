package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

func ConfigWorkspaceLifecycleScenarioV1(scenario string, params *WorkspaceParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	url := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	resource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)

	response, err := newWorkspaceResponse(params.Workspace.Name, secalib.WorkspaceProviderV1, resource, secalib.ApiVersion1,
		params.Tenant, params.Region, params.Workspace.InitialLabels)
	if err != nil {
		return nil, err
	}

	// Create a workspace
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = secalib.NewWorkspaceStatus(secalib.CreatingResourceState)
	response.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: url, params: params, responseBody: response, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	secalib.SetWorkspaceStatusState(response.Status, secalib.ActiveResourceState)
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: url, params: params, responseBody: response, currentState: "GetCreatedWorkspace", nextState: "UpdateWorkspace"}); err != nil {
		return nil, err
	}

	// Update the workspace
	setModifiedRegionalResourceMetadata(response.Metadata)
	secalib.SetWorkspaceStatusState(response.Status, secalib.UpdatingResourceState)
	response.Labels = params.Workspace.UpdatedLabels
	response.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: url, params: params, responseBody: response, currentState: "UpdateWorkspace", nextState: "GetUpdatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the updated workspace
	secalib.SetWorkspaceStatusState(response.Status, secalib.ActiveResourceState)
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: url, params: params, responseBody: response, currentState: "GetUpdatedWorkspace", nextState: "DeleteWorkspace"}); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: url, params: params, currentState: "DeleteWorkspace", nextState: "GetDeletedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: url, params: params, currentState: "GetDeletedWorkspace", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}
