package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	workspaceV1 "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigWorkspaceLifecycleScenarioV1(scenario string, params *WorkspaceParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	url := secalib.GenerateWorkspaceURL(params.Tenant, (*params.Workspace)[0].Name)
	resource := secalib.GenerateWorkspaceResource(params.Tenant, (*params.Workspace)[0].Name)

	response := newWorkspaceResponse((*params.Workspace)[0].Name, secalib.WorkspaceProviderV1, resource, secalib.ApiVersion1,
		params.Tenant, params.Region,
		(*params.Workspace)[0].InitialLabels)

	// Create a workspace
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = secalib.NewWorkspaceStatus(secalib.CreatingResourceState)
	response.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: url, params: params, headers: headerParamsGeneric(params.AuthToken), responseBody: response, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	secalib.SetWorkspaceStatusState(response.Status, secalib.ActiveResourceState)
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: url, params: params, headers: headerParamsGeneric(params.AuthToken), responseBody: response, currentState: "GetCreatedWorkspace", nextState: "UpdateWorkspace"}); err != nil {
		return nil, err
	}

	// Update the workspace
	setModifiedRegionalResourceMetadata(response.Metadata)
	secalib.SetWorkspaceStatusState(response.Status, secalib.UpdatingResourceState)
	response.Labels = (*params.Workspace)[0].UpdatedLabels
	response.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: url, params: params, headers: headerParamsGeneric(params.AuthToken), responseBody: response, currentState: "UpdateWorkspace", nextState: "GetUpdatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the updated workspace
	secalib.SetWorkspaceStatusState(response.Status, secalib.ActiveResourceState)
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: url, params: params, headers: headerParamsGeneric(params.AuthToken), responseBody: response, currentState: "GetUpdatedWorkspace", nextState: "DeleteWorkspace"}); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: url, params: params, currentState: "DeleteWorkspace", nextState: "GetDeletedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: url, params: params, currentState: "GetDeletedWorkspace", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}

func ConfigWorkspaceListLifecycleScenarioV1(scenario string, params *WorkspaceParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	// Workspace
	var workspaceList []schema.Workspace
	for i := range *params.Workspace {
		resource := secalib.GenerateWorkspaceResource(params.Tenant, (*params.Workspace)[i].Name)
		response := newWorkspaceResponse((*params.Workspace)[i].Name, secalib.WorkspaceProviderV1, resource, secalib.ApiVersion1,
			params.Tenant, params.Region,
			(*params.Workspace)[i].InitialLabels)
		var nextState string
		if i < len(*params.Workspace)-1 {
			nextState = (*params.Workspace)[i+1].Name
		} else {
			nextState = "GetWorkspaceList"
		}
		// Create a workspace
		setCreatedRegionalResourceMetadata(response.Metadata)
		response.Status = secalib.NewWorkspaceStatus(secalib.CreatingResourceState)
		response.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateWorkspaceURL(params.Tenant, (*params.Workspace)[i].Name), params: params, headers: headerParamsGeneric(params.AuthToken), responseBody: response, currentState: startedScenarioState, nextState: nextState}); err != nil {
			return nil, err
		}
		workspaceList = append(workspaceList, *response)
	}
	// List workspaces

	workspaceResource := secalib.GenerateWorkspaceListResource(params.Tenant)
	workspaceListResponse := &workspaceV1.WorkspaceIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.WorkspaceProviderV1,
			Resource: workspaceResource,
			Verb:     http.MethodGet,
		},
		Items: workspaceList,
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateWorkspaceListURL(params.Tenant), params: params, headers: headerParamsGeneric(params.AuthToken), responseBody: workspaceListResponse, currentState: (*params.Workspace)[len(*params.Workspace)-1].Name, nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	// List with limit
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateWorkspaceListURL(params.Tenant), params: params, headers: headerParamsLimit(params.AuthToken, "1"), responseBody: workspaceListResponse, currentState: startedScenarioState, nextState: "ListWithLimit"}); err != nil {
		return nil, err
	}
	// List with labels

	workspaceWithLabelResponse := &workspaceV1.WorkspaceIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.WorkspaceProviderV1,
			Resource: workspaceResource,
			Verb:     http.MethodGet,
		},
	}
	workspaceWithLabel := func(workspaceList []schema.Workspace) []schema.Workspace {
		var filteredWorkspaces []schema.Workspace
		for _, ws := range workspaceList {
			if val, ok := ws.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformance {
				filteredWorkspaces = append(filteredWorkspaces, ws)
			}
		}
		return filteredWorkspaces
	}
	workspaceWithLabelResponse.Items = workspaceWithLabel(workspaceList)
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateWorkspaceListURL(params.Tenant), params: params, headers: headerParamsLabel(params.AuthToken, secalib.EnvLabel, secalib.EnvConformance), responseBody: workspaceWithLabelResponse, currentState: "ListWithLimit", nextState: "ListWithLabels"}); err != nil {
		return nil, err
	}
	// List with limit & labels

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateWorkspaceListURL(params.Tenant), params: params, headers: headerParamsLimitAndLabel(params.AuthToken, "1", secalib.EnvLabel, secalib.EnvConformance), responseBody: workspaceWithLabelResponse, currentState: "ListWithLabels", nextState: startedScenarioState}); err != nil {
		return nil, err
	}
	return wm, nil
}
