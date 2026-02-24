package mockworkspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureProviderLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.WorkspaceProviderLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := *params.WorkspaceInitial

	url := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)

	response, err := builders.NewWorkspaceBuilder().
		Name(workspace.Metadata.Name).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(workspace.Metadata.Tenant).Region(workspace.Metadata.Region).
		Labels(workspace.Labels).
		Build()
	if err != nil {
		return err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(response, url, scenario.MockParams); err != nil {
		return err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetCreatingWorkspaceStub(response, url, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(response, url, scenario.MockParams); err != nil {
		return err
	}

	// Update the workspace
	if err := configurator.ConfigureUpdateWorkspaceStubWithLabels(response, url, scenario.MockParams, params.WorkspaceUpdated.Labels); err != nil {
		return err
	}

	// Get the updated workspace
	if err := configurator.ConfigureGetUpdatingWorkspaceStub(response, url, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(response, url, scenario.MockParams); err != nil {
		return err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(url, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(url, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
