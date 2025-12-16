package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigWorkspaceLifecycleScenarioV1(scenario string, params *WorkspaceParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	url := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)

	response, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.configureCreateWorkspaceStub(response, url, params); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.configureGetActiveWorkspaceStub(response, url, params); err != nil {
		return nil, err
	}

	// Update the workspace
	if err := configurator.configureUpdateWorkspaceStubWithLabels(response, url, params, params.Workspace.UpdatedLabels); err != nil {
		return nil, err
	}

	// Get the updated workspace
	if err := configurator.configureGetActiveWorkspaceStub(response, url, params); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(url, params); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(url, params); err != nil {
		return nil, err
	}

	return configurator.client, nil
}

func ConfigWorkspaceListLifecycleScenarioV1(scenario string, params *WorkspaceListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}
	var workspaceList []schema.Workspace

	for _, workspace := range *params.Workspace {
		url := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, workspace.Name)

		response, err := builders.NewWorkspaceBuilder().
			Name(workspace.Name).
			Provider(workspaceProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).Region(params.Region).
			Labels(workspace.InitialLabels).
			Build()
		if err != nil {
			return nil, err
		}

		// Create a workspace
		if err := configurator.configureCreateWorkspaceStub(response, url, params); err != nil {
			return nil, err
		}

		workspaceList = append(workspaceList, *response)
	}

	// List workspaces
	urlList := generators.GenerateWorkspaceListURL(workspaceProviderV1, params.Tenant)
	workspaceResource := generators.GenerateWorkspaceListResource(params.Tenant)
	workspaceListResponse := &workspace.WorkspaceIterator{
		Metadata: schema.ResponseMetadata{
			Provider: workspaceProviderV1,
			Resource: workspaceResource,
			Verb:     http.MethodGet,
		},
		Items: workspaceList,
	}

	if err := configurator.configureGetListActiveWorkspaceStub(workspaceListResponse, urlList, params, nil); err != nil {
		return nil, err
	}

	// List with limit
	if err := configurator.configureGetListActiveWorkspaceStub(workspaceListResponse, urlList, params, pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with labels
	workspaceWithLabel := func(workspaceList []schema.Workspace) []schema.Workspace {
		var filteredWorkspaces []schema.Workspace
		for _, ws := range workspaceList {
			if val, ok := ws.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredWorkspaces = append(filteredWorkspaces, ws)
			}
		}
		return filteredWorkspaces
	}
	workspaceListResponse.Items = workspaceWithLabel(workspaceList)
	if err := configurator.configureGetListActiveWorkspaceStub(workspaceListResponse, urlList, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with limit & labels
	workspaceListResponse.Items = workspaceWithLabel(workspaceList[:1])
	if err := configurator.configureGetListActiveWorkspaceStub(workspaceListResponse, urlList, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	for _, workspace := range workspaceList {
		url := generators.GenerateWorkspaceURL(workspaceProviderV1, workspace.Metadata.Name, workspace.Metadata.Name)

		// Delete the workspace
		if err := configurator.configureDeleteStub(url, params); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.configureGetNotFoundStub(url, params); err != nil {
			return nil, err
		}
	}
	return configurator.client, nil
}
