package mock

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wiremock/go-wiremock"
)

func CreateWorkspaceLifecycleScenarioV1(scenario string, params WorkspaceParamsV1) (*wiremock.Client, error) {
	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(workspaceURLV1, params.Tenant, params.Name)

	response := workspaceResponseV1{
		Metadata: metadataResponse{
			Name:       params.Name,
			Provider:   workspaceProviderV1,
			Resource:   fmt.Sprintf(workspaceResource, params.Tenant, params.Name),
			ApiVersion: version1,
			Kind:       workspaceKind,
			Tenant:     params.Tenant,
			Region:     params.Region,
		},
	}

	// Create a workspace
	response.Metadata.Verb = http.MethodPut
	response.Status.State = creatingStatusState
	response.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = 1
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          url,
		params:       params,
		response:     response,
		template:     workspaceResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    "GetCreatedWorkspace",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created workspace
	response.Metadata.Verb = http.MethodGet
	response.Status.State = activeStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          url,
		params:       params,
		response:     response,
		template:     workspaceResponseTemplateV1,
		currentState: "GetCreatedWorkspace",
		nextState:    "UpdateWorkspace",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the workspace
	response.Metadata.Verb = http.MethodPut
	response.Status.State = updatingStatusState
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = response.Metadata.ResourceVersion + 1
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          url,
		params:       params,
		response:     response,
		template:     workspaceResponseTemplateV1,
		currentState: "UpdateWorkspace",
		nextState:    "GetUpdatedWorkspace",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get updated workspace
	response.Metadata.Verb = http.MethodGet
	response.Status.State = activeStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          url,
		params:       params,
		response:     response,
		template:     workspaceResponseTemplateV1,
		currentState: "GetUpdatedWorkspace",
		nextState:    "DeleteWorkspace",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Delete the workspace
	response.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          url,
		params:       params,
		response:     response,
		currentState: "DeleteWorkspace",
		nextState:    "ReDeleteWorkspace",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Re-delete the workspace
	response.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          url,
		params:       params,
		response:     response,
		currentState: "ReDeleteWorkspace",
		nextState:    "GetDeletedWorkspace",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted workspace (not found)
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          url,
		params:       params,
		response:     response,
		currentState: "GetDeletedWorkspace",
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
