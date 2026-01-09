package mockworkspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/wiremock/go-wiremock"
)

func ConfigureListScenarioV1(scenario string, params *params.WorkspaceListParamsV1) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	configurator, err := stubs.NewStubConfigurator(scenario, params.MockParams)
	if err != nil {
		return nil, err
	}

	url := generators.GenerateWorkspaceListURL(constants.WorkspaceProviderV1, params.Tenant)

	// Create workspaces
	workspaceList, err := stubs.BulkCreateWorkspacesStubV1(configurator, params.BaseParams, params.Workspaces)
	if err != nil {
		return nil, err
	}
	workspaceListResponse, err := builders.NewWorkspaceIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(params.Tenant).
		Items(workspaceList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, params.MockParams, nil); err != nil {
		return nil, err
	}

	// List with limit
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, params.MockParams, mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with labels
	workspaceWithLabel := func(workspaceList []schema.Workspace) []schema.Workspace {
		var filteredWorkspaces []schema.Workspace
		for _, ws := range workspaceList {
			if val, ok := ws.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredWorkspaces = append(filteredWorkspaces, ws)
			}
		}
		return filteredWorkspaces
	}
	workspaceListResponse.Items = workspaceWithLabel(workspaceList)
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, params.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with limit & labels
	workspaceListResponse.Items = workspaceWithLabel(workspaceList[:1])
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, params.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	for _, workspace := range workspaceList {
		url := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Name, workspace.Metadata.Name)

		// Delete the workspace
		if err := configurator.ConfigureDeleteStub(url, params.MockParams); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(url, params.MockParams); err != nil {
			return nil, err
		}
	}
	return configurator.Client, nil
}
