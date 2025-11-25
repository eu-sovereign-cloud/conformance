package mock

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/conformance/secalib/builders"
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

	response, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(secalib.WorkspaceProviderV1).
		Resource(resource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		BuildResponse()
	if err != nil {
		return nil, err
	}
	// Create a workspace
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = secalib.NewWorkspaceStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: url, params: params, responseBody: response, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	secalib.SetWorkspaceStatusState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: url, params: params, responseBody: response, currentState: "GetCreatedWorkspace", nextState: "UpdateWorkspace"}); err != nil {
		return nil, err
	}

	// Update the workspace
	setModifiedRegionalResourceMetadata(response.Metadata)
	secalib.SetWorkspaceStatusState(response.Status, secalib.UpdatingResourceState)
	response.Labels = (*params.Workspace)[0].UpdatedLabels
	response.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: url, params: params, responseBody: response, currentState: "UpdateWorkspace", nextState: "GetUpdatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the updated workspace
	secalib.SetWorkspaceStatusState(response.Status, schema.ResourceStateActive)
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
		var currentState string
		if i < len(*params.Workspace)-1 {
			nextState = (*params.Workspace)[i+1].Name
		} else {
			nextState = "GetWorkspaceList"
		}
		if i == 0 {
			currentState = startedScenarioState
		} else {
			currentState = (*params.Workspace)[i].Name
		}
		// Create a workspace
		setCreatedRegionalResourceMetadata(response.Metadata)
		response.Status = secalib.NewWorkspaceStatus(secalib.CreatingResourceState)
		response.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateWorkspaceURL(params.Tenant, (*params.Workspace)[i].Name), params: params, responseBody: response, currentState: currentState, nextState: nextState}); err != nil {
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
		&stubConfig{
			url: secalib.GenerateWorkspaceListURL(params.Tenant), params: params, responseBody: workspaceListResponse,
			currentState: "GetWorkspaceList", nextState: "ListWorkspaceWithLimit",
		}); err != nil {
		return nil, err
	}

	// List with limit
	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateWorkspaceListURL(params.Tenant), params: params, pathParams: pathParamsLimit("1"), responseBody: workspaceListResponse,
			currentState: "ListWorkspaceWithLimit", nextState: "ListWorkspaceWithLabels",
		}); err != nil {
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
		&stubConfig{
			url: secalib.GenerateWorkspaceListURL(params.Tenant), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformance), responseBody: workspaceWithLabelResponse,
			currentState: "ListWorkspaceWithLabels", nextState: "ListWorkspaceWithLimitAndLabels",
		}); err != nil {
		return nil, err
	}
	// List with limit & labels

	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateWorkspaceListURL(params.Tenant), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformance), responseBody: workspaceWithLabelResponse,
			currentState: "ListWorkspaceWithLimitAndLabels", nextState: "DeleteWorkspace1",
		}); err != nil {
		return nil, err
	}

	// Delete all workspaces
	for i := range *params.Workspace {
		workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, (*params.Workspace)[i].Name)
		var nextState string
		var currentState string

		nextState = fmt.Sprintf("GetDeletedWorkspace%d", i+1)
		currentState = fmt.Sprintf("DeleteWorkspace%d", i+1)

		// Delete workspace
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: workspaceUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get deleted workspace
		getNextState := startedScenarioState
		if i < len(*params.Workspace)-1 {
			getNextState = fmt.Sprintf("DeleteWorkspace%d", i+2)
		}
		if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
			&stubConfig{url: workspaceUrl, params: params, currentState: nextState, nextState: getNextState}); err != nil {
			return nil, err
		}
	}

	return wm, nil
}
