package mock

import (
	"bytes"
	"log"
	"net/http"
	"text/template"

	responsesTemplate "github.com/eu-sovereign-cloud/conformance/internal/mock/responses"
	"github.com/wiremock/go-wiremock"
)

const (
	// Version
	Version1 = "v1"

	// Base Url for WireMock server
	WorkspaceProviderV1     = "/providers/seca.workspace/v1/"
	ComputeProviderV1       = "/providers/seca.compute/v1/"
	NetworkProviderV1       = "/providers/seca.network/v1/"
	StorageProviderV1       = "/providers/seca.storage/v1/"
	AuthorizationProviderV1 = "/providers/seca.authorization/v1/"

	// Resource
	WorkspaceResourceURL     = "seca.workspace/workspaces/"
	ComputeResourceURL       = "seca.compute/workspaces/"
	NetworkResourceURL       = "seca.network/workspaces/"
	StorageResourceURL       = "seca.storage/workspaces/"
	AuthorizationResourceURL = "seca.authorization/workspaces/"

	// ScenarioName
	ScenarioName = "Use Case Lifecycle"

	// State for the scenario
	StartedState = "Started"
	CreatedState = "UsecaseCreated"
	UpdatedState = "UsecaseUpdated"
	DeletedState = "UsecaseDeleted"

	// State for the workspace
	CreatingState = "creating"
	UpdatingState = "updating"

	// Kind
	WorkspaceKind = "workspace"
)

func CreateWorkspaceScenario(workspaceMock MockParams) error {
	wm := wiremock.NewClient(workspaceMock.WireMockURL)

	defer wm.ResetAllScenarios()

	workspaceMock.WireMockURL = WorkspaceProviderV1 + "tenants/" + workspaceMock.TenantName + "/workspaces/" + workspaceMock.WorkspaceName

	workspaceMetadata := UsecaseMetadata{
		Name:     workspaceMock.WorkspaceName,
		Tenant:   workspaceMock.TenantName,
		Region:   workspaceMock.Region,
		Version:  Version1,
		Kind:     WorkspaceKind,
		Resource: WorkspaceResourceURL,
		State:    CreatingState,
	}

	putStub(wm, UsecaseStubMetadata{
		Params:             workspaceMock,
		Metadata:           workspaceMetadata,
		Template:           responsesTemplate.WorkspacePutTemplateResponse,
		ScenarioState:      StartedState,
		NextScenarioState:  CreatedState,
		ScenarioHttpStatus: http.StatusCreated, // 201 Created
		ScenarioPriority:   1,
	})

	// Update Workspace
	workspaceMetadata.State = UpdatingState

	putStub(wm, UsecaseStubMetadata{
		Params:             workspaceMock,
		Metadata:           workspaceMetadata,
		Template:           responsesTemplate.WorkspacePutTemplateResponse,
		ScenarioState:      CreatedState,
		NextScenarioState:  UpdatedState,
		ScenarioHttpStatus: http.StatusOK, // 200 OK
		ScenarioPriority:   1,
	})

	// Delete Workspace
	// First delete the workspace
	deleteStub(wm, UsecaseStubMetadata{
		Params:             workspaceMock,
		Metadata:           workspaceMetadata,
		ScenarioState:      UpdatedState,
		NextScenarioState:  DeletedState,
		ScenarioHttpStatus: http.StatusAccepted, // 202 Accepted
		ScenarioPriority:   1,
	})
	// Second delete the workspace
	deleteStub(wm, UsecaseStubMetadata{
		Params:             workspaceMock,
		Metadata:           workspaceMetadata,
		ScenarioState:      DeletedState,
		NextScenarioState:  StartedState,
		ScenarioHttpStatus: http.StatusAccepted, // 202 Accepted
		ScenarioPriority:   1,
	})
	return nil
}

func CreateComputeScenario(computeParams MockParams) error {
	// Work in progress

	return nil
}

func CreateNetworkScenario(networkMock MockParams) error {
	// Work in progress

	return nil
}

func putStub(wm *wiremock.Client, stubMetadata UsecaseStubMetadata) {
	processTemplate, err := processTemplate(stubMetadata.Template, stubMetadata.Metadata)
	if err != nil {
		log.Printf("Error processing template: %v\n", err)
	}

	wm.StubFor(wiremock.Put(wiremock.URLPathMatching(stubMetadata.Params.WireMockURL)).
		WithHeader("Authorization", wiremock.Matching("Bearer "+stubMetadata.Params.Token)).
		InScenario(ScenarioName).
		WhenScenarioStateIs(stubMetadata.ScenarioState).
		WillSetStateTo(stubMetadata.NextScenarioState).
		WillReturnResponse(
			wiremock.NewResponse().
				WithStatus(int64(stubMetadata.ScenarioHttpStatus)).
				WithHeader("Content-Type", "application/json").
				WithJSONBody(processTemplate),
		).
		AtPriority(int64(stubMetadata.ScenarioPriority)))
}

func getStub(wm *wiremock.Client, stubMetadata UsecaseStubMetadata) {
	wm.StubFor(wiremock.Get(wiremock.URLPathMatching(stubMetadata.Params.WireMockURL)).
		WithHeader("Authorization", wiremock.Matching("Bearer "+stubMetadata.Params.Token)).
		InScenario(ScenarioName).
		WhenScenarioStateIs(stubMetadata.ScenarioState).
		WillSetStateTo(stubMetadata.NextScenarioState).
		WillReturnResponse(
			wiremock.NewResponse().
				WithStatus(int64(stubMetadata.ScenarioHttpStatus)),
		).
		AtPriority(int64(stubMetadata.ScenarioPriority)))
}

func deleteStub(wm *wiremock.Client, stubMetadata UsecaseStubMetadata) {
	wm.StubFor(wiremock.Delete(wiremock.URLPathMatching(stubMetadata.Params.WireMockURL)).
		WithHeader("Authorization", wiremock.Matching("Bearer "+stubMetadata.Params.Token)).
		InScenario(ScenarioName).
		WhenScenarioStateIs(stubMetadata.ScenarioState).
		WillSetStateTo(stubMetadata.NextScenarioState).
		WillReturnResponse(
			wiremock.NewResponse().
				WithStatus(int64(stubMetadata.ScenarioHttpStatus)),
		).
		AtPriority(int64(stubMetadata.ScenarioPriority)))
}

func processTemplate(templ string, data any) (string, error) {
	tmpl := template.Must(template.New("response").Parse(templ))

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
