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

func ConfigureListScenarioV1(scenario string, mockParams *mock.MockParams, suiteParams *params.WorkspaceListV1Params) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	workspaces := suiteParams.Workspaces
	workspace := workspaces[0]
	configurator, err := stubs.NewStubConfigurator(scenario, mockParams)
	if err != nil {
		return nil, err
	}

	url := generators.GenerateWorkspaceListURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant)

	// Create workspaces
	err = stubs.BulkCreateWorkspacesStubV1(configurator, mockParams, workspaces)
	if err != nil {
		return nil, err
	}
	workspaceListResponse, err := builders.NewWorkspaceIteratorBuilder().
		Provider(constants.WorkspaceProviderV1).
		Tenant(workspace.Metadata.Tenant).
		Items(workspaces).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, mockParams, nil); err != nil {
		return nil, err
	}

	// List with limit
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, mockParams, mock.PathParamsLimit("1")); err != nil {
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
	workspaceListResponse.Items = workspaceWithLabel(workspaces)
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, mockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with limit & labels
	workspaceListResponse.Items = workspaceWithLabel(workspaces[:1])
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, mockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	for _, workspace := range workspaces {
		url := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Name, workspace.Metadata.Name)

		// Delete the workspace
		if err := configurator.ConfigureDeleteStub(url, mockParams); err != nil {
			return nil, err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(url, mockParams); err != nil {
			return nil, err
		}
	}

	// Finish the stubs configuration
	if client, err := configurator.Finish(); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}
