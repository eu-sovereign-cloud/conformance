package mock

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"

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
	StartedState   = "Started"
	GettedState    = "UsecaseGetted"
	CreatedState   = "UsecaseCreated"
	UpdatedState   = "UsecaseUpdated"
	DeletedState   = "UsecaseDeleted"
	PowerOffState  = "UsecasePowerOff"
	PowerOnState   = "UsecasePowerOn"
	RestartedState = "UsecaseRestarted"

	// State for the workspace
	CreatingState = "creating"
	UpdatingState = "updating"
	ActiveState   = "active"

	// Kind
	WorkspaceKind     = "workspace"
	ComputeKind       = "compute"
	NetworkKind       = "network"
	StorageKind       = "storage"
	AuthorizationKind = "authorization"
)

func CreateWorkspaceScenario(workspaceMock MockParams) error {
	wm := wiremock.NewClient(workspaceMock.WireMockURL)
	wm.Clear()

	workspaceMock.WireMockURL = WorkspaceProviderV1 + "tenants/" + workspaceMock.TenantName + "/workspaces/" + workspaceMock.Name

	// Create Workspace
	workspaceMetadata := UsecaseMetadata{
		Name:             workspaceMock.WorkspaceName,
		CreatedAt:        time.Now().Format(time.RFC3339),
		LastModifiedAt:   time.Now().Format(time.RFC3339),
		Tenant:           workspaceMock.TenantName,
		Region:           workspaceMock.Region,
		Version:          Version1,
		Kind:             WorkspaceKind,
		Resource:         WorkspaceResourceURL,
		State:            CreatingState,
		LastTransitionAt: time.Now().Format(time.RFC3339),
	}

	err := putStub(wm, UsecaseStubMetadata{
		Params:             workspaceMock,
		Metadata:           workspaceMetadata,
		Template:           responsesTemplate.WorkspaceTemplateResponse,
		ScenarioState:      StartedState,
		NextScenarioState:  CreatedState,
		ScenarioHttpStatus: http.StatusCreated, // 201 Created
		ScenarioPriority:   1,
	})
	if err != nil {
		return err
	}

	// Update Workspace
	workspaceMetadata.State = UpdatingState

	err = putStub(wm, UsecaseStubMetadata{
		Params:             workspaceMock,
		Metadata:           workspaceMetadata,
		Template:           responsesTemplate.WorkspaceTemplateResponse,
		ScenarioState:      CreatedState,
		NextScenarioState:  UpdatedState,
		ScenarioHttpStatus: http.StatusOK, // 200 OK
		ScenarioPriority:   1,
	})
	if err != nil {
		return err
	}

	// Delete Workspace
	// First delete the workspace
	err = deleteStub(wm, UsecaseStubMetadata{
		Params:             workspaceMock,
		Metadata:           workspaceMetadata,
		ScenarioState:      UpdatedState,
		NextScenarioState:  DeletedState,
		ScenarioHttpStatus: http.StatusAccepted, // 202 Accepted
		ScenarioPriority:   1,
	})
	if err != nil {
		return err
	}

	// Second delete the workspace
	err = deleteStub(wm, UsecaseStubMetadata{
		Params:             workspaceMock,
		Metadata:           workspaceMetadata,
		ScenarioState:      DeletedState,
		NextScenarioState:  StartedState,
		ScenarioHttpStatus: http.StatusAccepted, // 202 Accepted
		ScenarioPriority:   1,
	})
	if err != nil {
		return err
	}

	workspaceMetadata.State = ActiveState
	err = getStub(wm, UsecaseStubMetadata{
		Params:             workspaceMock,
		Metadata:           workspaceMetadata,
		ScenarioHttpStatus: http.StatusOK, // 200 OK
		ScenarioPriority:   1,
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateComputeScenario(computeMock MockParams) error {

	wm := wiremock.NewClient(computeMock.WireMockURL)
	//wm.Clear()

	computeMetadata := UsecaseMetadata{
		Name:             computeMock.Name,
		Tenant:           computeMock.TenantName,
		Region:           computeMock.Region,
		CreatedAt:        time.Now().Format(time.RFC3339),
		LastModifiedAt:   time.Now().Format(time.RFC3339),
		Version:          Version1,
		Kind:             ComputeKind,
		Resource:         ComputeResourceURL,
		Workspace:        computeMock.WorkspaceName,
		State:            CreatingState,
		LastTransitionAt: time.Now().Format(time.RFC3339),
	}

	//Get sku
	computeMock.WireMockURL = ComputeProviderV1 + "tenants/" + computeMock.TenantName + "/skus/" + computeMock.SkuName
	err := getStub(wm, UsecaseStubMetadata{
		Params:             computeMock,
		Metadata:           computeMetadata,
		Template:           responsesTemplate.ComputeSkuTemplateResponse,
		ScenarioState:      StartedState,
		NextScenarioState:  CreatedState,
		ScenarioHttpStatus: http.StatusOK, // 200 OK
		ScenarioPriority:   1,
	})

	if err != nil {
		log.Printf("Error getting compute sku: %v\n", err)
		return err
	}

	// Create Instance
	computeMock.WireMockURL = ComputeProviderV1 + "tenants/" + computeMock.TenantName + "/workspaces/" + computeMock.WorkspaceName + "/instances/" + computeMock.Name
	err = putStub(wm, UsecaseStubMetadata{
		Params:             computeMock,
		Metadata:           computeMetadata,
		Template:           responsesTemplate.ComputePutTemplateResponse,
		RequestTemplate:    responsesTemplate.ComputePutRestrictionTemplate,
		ScenarioState:      CreatedState,
		NextScenarioState:  GettedState,
		ScenarioHttpStatus: http.StatusCreated, // 201 Created
		ScenarioPriority:   1,
	})

	if err != nil {
		log.Printf("Error creating compute instance: %v\n", err)
		return err
	}

	// Get Instance
	computeMetadata.State = ActiveState
	err = getStub(wm, UsecaseStubMetadata{
		Params:             computeMock,
		Metadata:           computeMetadata,
		Template:           responsesTemplate.ComputeGetTemplateResponse,
		ScenarioState:      GettedState,
		NextScenarioState:  UpdatedState,
		ScenarioHttpStatus: http.StatusOK, // 200 OK
		ScenarioPriority:   1,
	})
	if err != nil {
		log.Printf("Error getting compute instance: %v\n", err)
		return err
	}

	// Update Instance
	computeMetadata.State = UpdatingState
	err = putStub(wm, UsecaseStubMetadata{Params: computeMock,
		Metadata:           computeMetadata,
		Template:           responsesTemplate.ComputePutTemplateResponse,
		ScenarioState:      UpdatedState,
		NextScenarioState:  PowerOffState,
		ScenarioHttpStatus: http.StatusOK, // 200 OK
		ScenarioPriority:   2,
	})
	if err != nil {
		log.Printf("Error updating compute instance: %v\n", err)
		return err
	}

	// Power off Instance
	err = postStub(wm, UsecaseStubMetadata{Params: computeMock,
		Metadata:           computeMetadata,
		Template:           "",
		ScenarioState:      PowerOffState,
		NextScenarioState:  PowerOnState,
		ScenarioHttpStatus: http.StatusAccepted, // 202 Accepted
		ScenarioPriority:   1,
	})
	if err != nil {
		log.Printf("Error powering off compute instance: %v\n", err)
		return err
	}

	// Power on Instance
	err = postStub(wm, UsecaseStubMetadata{Params: computeMock,
		Metadata:           computeMetadata,
		Template:           "",
		ScenarioState:      PowerOnState,
		NextScenarioState:  RestartedState,
		ScenarioHttpStatus: http.StatusAccepted, // 202 Accepted
		ScenarioPriority:   1,
	})
	if err != nil {
		return err
	}

	// Restart Instance
	err = postStub(wm, UsecaseStubMetadata{Params: computeMock,
		Metadata:           computeMetadata,
		Template:           "",
		ScenarioState:      RestartedState,
		NextScenarioState:  DeletedState,
		ScenarioHttpStatus: http.StatusAccepted, // 202 Accepted
		ScenarioPriority:   1,
	})
	if err != nil {
		return err
	}

	// Delete Instance
	err = deleteStub(wm, UsecaseStubMetadata{Params: computeMock,
		Metadata:           computeMetadata,
		Template:           "",
		ScenarioState:      DeletedState,
		NextScenarioState:  StartedState,
		ScenarioHttpStatus: http.StatusAccepted, // 202 Accepted
		ScenarioPriority:   1,
	})
	if err != nil {
		return err
	}

	return nil
}

func CreateStorageScenario(storageMock MockParams) error {
	// Work in progress

	return nil
}

func CreateNetworkScenario(networkMock MockParams) error {
	// Work in progress

	return nil
}

func putStub(wm *wiremock.Client, stubMetadata UsecaseStubMetadata) error {
	processTemplate, err := processTemplate(stubMetadata.Template, stubMetadata.Metadata)
	if err != nil {
		log.Printf("Error processing template: %v\n", err)
		return err
	}

	err = wm.StubFor(wiremock.Put(wiremock.URLPathMatching(stubMetadata.Params.WireMockURL)).
		WithHeader("Authorization", wiremock.Matching("Bearer "+stubMetadata.Params.Token)).
		InScenario(ScenarioName).
		WhenScenarioStateIs(stubMetadata.ScenarioState).
		WillSetStateTo(stubMetadata.NextScenarioState).
		WithBodyPattern(wiremock.MatchingJsonPath(stubMetadata.RequestTemplate)).
		WillReturnResponse(
			wiremock.NewResponse().
				WithStatus(int64(stubMetadata.ScenarioHttpStatus)).
				WithHeader("Content-Type", "application/json").
				WithJSONBody(processTemplate),
		).
		AtPriority(int64(stubMetadata.ScenarioPriority)))
	if condition := err != nil; condition {
		log.Printf("Error configuring put method: %v\n", err)
		return err
	}

	return nil
}
func postStub(wm *wiremock.Client, stubMetadata UsecaseStubMetadata) error {

	err := wm.StubFor(wiremock.Post(wiremock.URLPathMatching(stubMetadata.Params.WireMockURL)).
		WithHeader("Authorization", wiremock.Matching("Bearer "+stubMetadata.Params.Token)).
		InScenario(ScenarioName).
		WhenScenarioStateIs(stubMetadata.ScenarioState).
		WillSetStateTo(stubMetadata.NextScenarioState).
		WillReturnResponse(
			wiremock.NewResponse().
				WithStatus(int64(stubMetadata.ScenarioHttpStatus)),
		).
		AtPriority(int64(stubMetadata.ScenarioPriority)))

	if err != nil {
		log.Printf("Error processing template: %v\n", err)
		return err
	}
	return nil
}

func getStub(wm *wiremock.Client, stubMetadata UsecaseStubMetadata) error {
	processTemplate, err := processTemplate(stubMetadata.Template, stubMetadata.Metadata)
	if err != nil {
		log.Printf("Error processing template: %v\n", err)
		return err
	}

	err = wm.StubFor(wiremock.Put(wiremock.URLPathMatching(stubMetadata.Params.WireMockURL)).
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
	if condition := err != nil; condition {
		log.Printf("Error configuring put method: %v\n", err)
		return err
	}

	return nil
}
func deleteStub(wm *wiremock.Client, stubMetadata UsecaseStubMetadata) error {
	err := wm.StubFor(wiremock.Delete(wiremock.URLPathMatching(stubMetadata.Params.WireMockURL)).
		WithHeader("Authorization", wiremock.Matching("Bearer "+stubMetadata.Params.Token)).
		InScenario(ScenarioName).
		WhenScenarioStateIs(stubMetadata.ScenarioState).
		WillSetStateTo(stubMetadata.NextScenarioState).
		WillReturnResponse(
			wiremock.NewResponse().
				WithStatus(int64(stubMetadata.ScenarioHttpStatus)),
		).
		AtPriority(int64(stubMetadata.ScenarioPriority)))
	if err != nil {
		log.Printf("Error configuring delete method: %v\n", err)
		return err
	}
	return nil
}

func processTemplate(templ string, data any) (map[string]interface{}, error) {
	tmpl := template.Must(template.New("response").Delims("[[", "]]").Parse(templ))

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return nil, err
	}
	var dataJsonMap map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &dataJsonMap)
	if err != nil {
		return nil, err
	}

	return dataJsonMap, nil
}
