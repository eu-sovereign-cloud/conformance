package mock

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wiremock/go-wiremock"
)

func CreateWorkspaceScenario(scenario string, mockParams MockParams) (*wiremock.Client, error) {
	wm := wiremock.NewClient(mockParams.MockURL)
	if err := wm.ResetAllScenarios(); err != nil {
		return nil, err
	}

	mockURL := fmt.Sprintf(workspaceV1BaseURL, mockParams.Tenant, mockParams.Workspace)
	mockParams.MockURL = mockURL

	resource := fmt.Sprintf(workspaceResource, mockParams.Tenant, mockParams.Workspace)

	response := workspaceResponseV1{
		Metadata: regionalMetadataResponseV1{
			Name:       mockParams.Workspace,
			Provider:   workspaceV1Provider,
			Resource:   resource,
			ApiVersion: version1,
			Kind:       workspaceKind,
			Tenant:     mockParams.Tenant,
			Region:     mockParams.Region,
		},
		Status: statusResponseV1{},
	}

	// Create a workspace
	response.Metadata.Verb = http.MethodPut
	response.Status.State = creatingStatusState
	response.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = 1
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       mockParams,
		response:     response,
		template:     workspaceTemplateResponse,
		currentState: startedScenarioState,
		nextState:    creatingScenarioState,
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get created workspace
	response.Metadata.Verb = http.MethodGet
	response.Status.State = activeStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       mockParams,
		response:     response,
		template:     workspaceTemplateResponse,
		currentState: creatingScenarioState,
		nextState:    createdScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Update the workspace
	response.Metadata.Verb = http.MethodPut
	response.Status.State = updatingStatusState
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = response.Metadata.ResourceVersion + 1
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       mockParams,
		response:     response,
		template:     workspaceTemplateResponse,
		currentState: createdScenarioState,
		nextState:    updatingScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get updated workspace
	response.Metadata.Verb = http.MethodGet
	response.Status.State = activeStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       mockParams,
		response:     response,
		template:     workspaceTemplateResponse,
		currentState: updatingScenarioState,
		nextState:    updatedScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Delete the workspace
	response.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       mockParams,
		response:     response,
		currentState: updatedScenarioState,
		nextState:    deletingScenarioState,
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Re-delete the workspace
	response.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       mockParams,
		response:     response,
		currentState: deletingScenarioState,
		nextState:    redeletingScenarioState,
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get deleted workspace
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       mockParams,
		response:     response,
		currentState: redeletingScenarioState,
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
