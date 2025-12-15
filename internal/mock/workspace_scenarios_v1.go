package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
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
