package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigurePublicIpLifecycleScenarioV1(scenario *mockscenarios.Scenario, params params.PublicIpLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := params.Workspace
	publicIp := params.PublicIpInitial

	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	publicIpURL := generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, publicIp.Metadata.Tenant, publicIp.Metadata.Workspace, publicIp.Metadata.Name)

	// Workspace

	// Create workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Public ip

	// Create public ip
	if err := configurator.ConfigureCreatePublicIpStub(publicIp, publicIpURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created public ip
	if err := configurator.ConfigureGetCreatingPublicIpStub(publicIp, publicIpURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActivePublicIpStub(publicIp, publicIpURL, scenario.MockParams); err != nil {
		return err
	}

	// Update public ip (change address)
	publicIp = params.PublicIpUpdated
	if err := configurator.ConfigureUpdatePublicIpStub(publicIp, publicIpURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated public ip
	if err := configurator.ConfigureGetUpdatingPublicIpStub(publicIp, publicIpURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActivePublicIpStub(publicIp, publicIpURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the public ip
	if err := configurator.ConfigureDeleteStub(publicIpURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingPublicIpStub(publicIp, publicIpURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(publicIpURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
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
