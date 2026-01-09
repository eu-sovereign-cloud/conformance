package mockworkspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, params *params.WorkspaceLifeCycleParamsV1) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	configurator, err := stubs.NewStubConfigurator(scenario, params.MockParams)
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
	if err := configurator.ConfigureCreateWorkspaceStub(response, url, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(response, url, params.MockParams); err != nil {
		return nil, err
	}

	// Update the workspace
	if err := configurator.ConfigureUpdateWorkspaceStubWithLabels(response, url, params.MockParams, params.Workspace.UpdatedLabels); err != nil {
		return nil, err
	}

	// Get the updated workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(response, url, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(url, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(url, params.MockParams); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
