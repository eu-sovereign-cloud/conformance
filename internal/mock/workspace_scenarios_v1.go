package mock

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

type WorkspaceV1Scenarios struct {
	Scenarios
}

func NewWorkspaceV1Scenarios(authToken string, tenant string, region string, mockURL string) *WorkspaceV1Scenarios {
	return &WorkspaceV1Scenarios{
		Scenarios: Scenarios{
			params: secalib.GeneralParams{
				AuthToken: authToken,
				Tenant:    tenant,
				Region:    region,
			},
			mockURL: mockURL,
		},
	}
}

func (scenarios *WorkspaceV1Scenarios) ConfigureLifecycleScenario(id string, params WorkspaceParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to Workspace Lifecycle Scenario")

	name := "WorkspaceV1Lifecycle_" + id

	wm, err := scenarios.newClient()
	if err != nil {
		return nil, err
	}

	url := secalib.GenerateWorkspaceURL(scenarios.params.Tenant, params.Workspace.Name)
	resource := secalib.GenerateWorkspaceResource(scenarios.params.Tenant, params.Workspace.Name)

	response := &resourceResponse[secalib.WorkspaceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Workspace.Name,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   resource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     scenarios.params.Tenant,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Labels: &[]secalib.Label{},
	}

	for _, labels := range *params.Workspace.InitialSpec.Labels {
		*response.Labels = append(*response.Labels, secalib.Label{
			Name:  labels.Name,
			Value: labels.Value,
		})
	}
	// Create a workspace
	response.Metadata.Verb = http.MethodPut
	response.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = 1
	response.Status.State = secalib.CreatingStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          url,
		params:       scenarios.params,
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
	response.Status.State = secalib.ActiveStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          url,
		params:       scenarios.params,
		response:     response,
		template:     workspaceResponseTemplateV1,
		currentState: "GetCreatedWorkspace",
		nextState:    "UpdateWorkspace",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the workspace

	for _, labels := range *params.Workspace.UpdatedSpec.Labels {
		*response.Labels = append(*response.Labels, secalib.Label{
			Name:  labels.Name,
			Value: labels.Value,
		})
	}
	response.Metadata.Verb = http.MethodPut
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = response.Metadata.ResourceVersion + 1
	response.Status.State = secalib.UpdatingStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          url,
		params:       scenarios.params,
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
	response.Status.State = secalib.ActiveStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          url,
		params:       scenarios.params,
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
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          url,
		params:       scenarios.params,
		response:     response,
		currentState: "DeleteWorkspace",
		nextState:    "GetDeletedWorkspace",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted workspace
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          url,
		params:       scenarios.params,
		response:     response,
		currentState: "GetDeletedWorkspace",
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
