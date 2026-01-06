package workspace

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, params *mock.WorkspaceLifeCycleParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := stubs.NewScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	url := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, params.Tenant, params.Workspace.Name)

	response, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(response, url, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(response, url, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Update the workspace
	if err := configurator.ConfigureUpdateWorkspaceStubWithLabels(response, url, params.GetBaseParams(), params.Workspace.UpdatedLabels); err != nil {
		return nil, err
	}

	// Get the updated workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(response, url, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(url, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(url, params.GetBaseParams()); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
