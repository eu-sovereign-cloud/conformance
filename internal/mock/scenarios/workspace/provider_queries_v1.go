package mockworkspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func ConfigureProviderQueriesV1(scenario *mockscenarios.Scenario, params *params.WorkspaceProviderQueriesV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspaces := params.Workspaces
	workspace := workspaces[0]

	url := generators.GenerateWorkspaceListURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant)

	// Create workspaces
	err = stubs.BulkCreateWorkspacesStubV1(configurator, scenario.MockParams, workspaces)
	if err != nil {
		return err
	}
	workspaceListResponse, err := builders.NewWorkspaceIteratorBuilder().
		Provider(constants.WorkspaceProviderV1).
		Tenant(workspace.Metadata.Tenant).
		Items(workspaces).
		Build()
	if err != nil {
		return err
	}

	// List
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with limit
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
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
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with limit & labels
	workspaceListResponse.Items = workspaceWithLabel(workspaces[:1])
	if err := configurator.ConfigureGetListActiveWorkspaceStub(workspaceListResponse, url, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	for _, workspace := range workspaces {
		url := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Name, workspace.Metadata.Name)

		// Delete the workspace
		if err := configurator.ConfigureDeleteStub(url, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetNotFoundStub(url, scenario.MockParams); err != nil {
			return err
		}
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
