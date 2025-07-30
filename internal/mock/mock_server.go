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
	//Base Url for WireMock server
	Provider = "/providers/seca.workspace/v1/"

	//Resource
	WorkspaceResourceURL = "seca.workspace/workspaces/"

	//ScenarioName
	ScenarioName = "Workspace Lifecycle"

	//State for the scenario
	StartedState = "Started"
	CreatedState = "WorkspaceUsecaseCreated"
	UpdatedState = "WorkspaceUsecaseUpdated"
	DeletedState = "WorkspaceUsecaseDeleted"

	//State for the workspace
	CreatingState = "creating"
	UpdatingState = "updating"

	//Kind
	WorkspaceKind = "workspace"

	//Version
	WorkspaceVersion = "v1"
)

func CreateWorkspaceScenario(workspaceMock Workspace) error {

	wm := wiremock.NewClient(workspaceMock.WireMockURL)

	defer wm.ResetAllScenarios()

	workspaceMock.WireMockURL = Provider + "tenants/" + workspaceMock.TenantName + "/workspaces/" + workspaceMock.WorkspaceName

	workspaceMetadata := WorkspaceMetadata{
		Name:     workspaceMock.WorkspaceName,
		Tenant:   workspaceMock.TenantName,
		Region:   workspaceMock.Region,
		Version:  WorkspaceVersion,
		Kind:     WorkspaceKind,
		Resource: WorkspaceResourceURL,
		State:    CreatingState,
	}

	putStub(wm, WorkspaceStubMetadata{
		WorkspaceMock:      workspaceMock,
		WorkspaceMetadata:  workspaceMetadata,
		Template:           responsesTemplate.WorkspacePutTemplateResponse,
		ScenarioState:      StartedState,
		NextScenarioState:  CreatedState,
		ScenarioHttpStatus: http.StatusCreated, // 201 Created
		ScenarioPriority:   1,
	})

	// Update Workspace
	workspaceMetadata.State = UpdatingState

	putStub(wm, WorkspaceStubMetadata{
		WorkspaceMock:      workspaceMock,
		WorkspaceMetadata:  workspaceMetadata,
		Template:           responsesTemplate.WorkspacePutTemplateResponse,
		ScenarioState:      CreatedState,
		NextScenarioState:  UpdatedState,
		ScenarioHttpStatus: http.StatusOK, // 200 OK
		ScenarioPriority:   1,
	})

	// Delete Workspace
	// First delete the workspace
	deleteStub(wm, WorkspaceStubMetadata{
		WorkspaceMock:      workspaceMock,
		WorkspaceMetadata:  workspaceMetadata,
		ScenarioState:      UpdatedState,
		NextScenarioState:  DeletedState,
		ScenarioHttpStatus: http.StatusAccepted, // 202 Accepted
		ScenarioPriority:   1,
	})
	// Second delete the workspace
	deleteStub(wm, WorkspaceStubMetadata{
		WorkspaceMock:      workspaceMock,
		WorkspaceMetadata:  workspaceMetadata,
		ScenarioState:      DeletedState,
		NextScenarioState:  StartedState,
		ScenarioHttpStatus: http.StatusAccepted, // 202 Accepted
		ScenarioPriority:   1,
	})
	return nil
}

func putStub(wm *wiremock.Client, stubMetadata WorkspaceStubMetadata) {
	processTemplate, err := processTemplate(stubMetadata.Template, stubMetadata.WorkspaceMetadata)
	if err != nil {
		log.Printf("Error processing template: %v\n", err)
	}

	wm.StubFor(wiremock.Put(wiremock.URLPathMatching(stubMetadata.WorkspaceMock.WireMockURL)).
		WithHeader("Authorization", wiremock.Matching("Bearer "+stubMetadata.WorkspaceMock.Token)).
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

func deleteStub(wm *wiremock.Client, stubMetadata WorkspaceStubMetadata) {
	wm.StubFor(wiremock.Delete(wiremock.URLPathMatching(stubMetadata.WorkspaceMock.WireMockURL)).
		WithHeader("Authorization", wiremock.Matching("Bearer "+stubMetadata.WorkspaceMock.Token)).
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
