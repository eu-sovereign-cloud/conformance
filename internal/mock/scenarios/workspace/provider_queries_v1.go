package workspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	pkgscenarios "github.com/eu-sovereign-cloud/conformance/pkg/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/mock/stubs"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func ConfigureProviderQueriesV1(scenario *pkgscenarios.Scenario, params params.WorkspaceProviderQueriesV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspaces := params.Workspaces

	url := generators.GenerateWorkspaceListURL(sdkconsts.WorkspaceProviderV1Name, workspaces.Items[0].Metadata.Tenant)

	// Create workspaces
	err = stubs.BulkCreateWorkspacesStubV1(configurator, scenario.MockParams, workspaces.Items)
	if err != nil {
		return err
	}
	workspaceListResponse := &params.Workspaces

	// List
	if err := configurator.ConfigureListActiveWorkspaceStub(workspaceListResponse, url, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with limit
	if err := configurator.ConfigureListActiveWorkspaceStub(workspaceListResponse, url, scenario.MockParams, scenarios.PathParamsLimit("1")); err != nil {
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
	workspaceListResponse.Items = workspaceWithLabel(workspaces.Items)
	if err := configurator.ConfigureListActiveWorkspaceStub(workspaceListResponse, url, scenario.MockParams, scenarios.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with limit & labels
	workspaceListResponse.Items = workspaceWithLabel(workspaces.Items[:1])
	if err := configurator.ConfigureListActiveWorkspaceStub(workspaceListResponse, url, scenario.MockParams, scenarios.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	for _, workspace := range workspaces.Items {
		url := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Name, workspace.Metadata.Name)

		// Delete the workspace
		if err := configurator.ConfigureDeleteStub(url, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted workspace
		if err := configurator.ConfigureGetDeletingWorkspaceStub(&workspace, url, scenario.MockParams); err != nil {
			return err
		}
		if err := configurator.ConfigureGetNotFoundStub(url, scenario.MockParams); err != nil {
			return err
		}
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
