package workspace

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigureListScenarioV1(scenario string, params *mock.WorkspaceListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := stubs.NewScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	url := generators.GenerateWorkspaceListURL(mock.WorkspaceProviderV1, params.Tenant)

	// Create workspaces
	workspaceList, err := stubs.BulkCreateWorkspacesStubV1(configurator, params.GetBaseParams(), params.Workspaces)
	if err != nil {
		return nil, err
	}
	workspaceListResponse, err := builders.NewWorkspaceIteratorBuilder().
		Provider(mock.StorageProviderV1).
		Tenant(params.Tenant).
		Items(workspaceList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with limit
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
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
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, params.GetBaseParams(), mock.PathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with limit & labels
	workspaceListResponse.Items = workspaceWithLabel(workspaceList[:1])
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	for _, workspace := range workspaceList {
		url := generators.GenerateWorkspaceURL(mock.WorkspaceProviderV1, workspace.Metadata.Name, workspace.Metadata.Name)

		// Delete the workspace
		if err := configurator.ConfigureDeleteStub(url, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(url, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}
	return configurator.Client, nil
}
