package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigWorkspaceLifecycleScenarioV1(scenario string, params *WorkspaceParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	url := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	resource := generators.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)

	response, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).Resource(resource).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		BuildResponse()
	if err != nil {
		return nil, err
	}
	// Create a workspace
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = newWorkspaceStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: url, params: params, responseBody: response, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	setWorkspaceState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: url, params: params, responseBody: response, currentState: "GetCreatedWorkspace", nextState: "UpdateWorkspace"}); err != nil {
		return nil, err
	}

	// Update the workspace
	setModifiedRegionalResourceMetadata(response.Metadata)
	setWorkspaceState(response.Status, schema.ResourceStateUpdating)
	response.Labels = params.Workspace.UpdatedLabels
	response.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: url, params: params, responseBody: response, currentState: "UpdateWorkspace", nextState: "GetUpdatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the updated workspace
	setWorkspaceState(response.Status, schema.ResourceStateActive)
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
