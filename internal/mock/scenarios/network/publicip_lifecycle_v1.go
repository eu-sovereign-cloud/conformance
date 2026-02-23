package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
)

func ConfigurePublicIpLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.PublicIpLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := *params.Workspace
	publicIpInitial := *params.PublicIpInitial
	publicIpUpdated := *params.PublicIpUpdated

	workspaceURL := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	publicIpURL := generators.GeneratePublicIpURL(constants.NetworkProviderV1, publicIpInitial.Metadata.Tenant, publicIpInitial.Metadata.Workspace, publicIpInitial.Metadata.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(workspace.Metadata.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(workspace.Metadata.Tenant).Region(workspace.Metadata.Region).
		Labels(workspace.Labels).
		Build()
	if err != nil {
		return err
	}

	// Create workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Get created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Public ip
	publicIpResponse, err := builders.NewPublicIpBuilder().
		Name(publicIpInitial.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(publicIpInitial.Metadata.Tenant).Workspace(publicIpInitial.Metadata.Workspace).Region(publicIpInitial.Metadata.Region).
		Spec(&publicIpInitial.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create public ip
	if err := configurator.ConfigureCreatePublicIpStub(publicIpResponse, publicIpURL, scenario.MockParams); err != nil {
		return err
	}

	// Get created public ip
	if err := configurator.ConfigureGetActivePublicIpStub(publicIpResponse, publicIpURL, scenario.MockParams); err != nil {
		return err
	}

	// Update public ip (change address)
	publicIpResponse.Spec = publicIpUpdated.Spec
	if err := configurator.ConfigureUpdatePublicIpStub(publicIpResponse, publicIpURL, scenario.MockParams); err != nil {
		return err
	}

	// Get updated public ip
	if err := configurator.ConfigureGetActivePublicIpStub(publicIpResponse, publicIpURL, scenario.MockParams); err != nil {
		return err
	}

	// Deletions
	if err := configurator.ConfigureDeleteStub(publicIpURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(publicIpURL, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureDeleteStub(workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}

	return nil
}
