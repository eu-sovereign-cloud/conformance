package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureSecurityGroupLifecycleScenarioV1(scenario *mockscenarios.Scenario, params params.SecurityGroupLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := params.Workspace
	securityGroup := params.SecurityGroupInitial

	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	securityGroupURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, securityGroup.Metadata.Tenant, securityGroup.Metadata.Workspace, securityGroup.Metadata.Name)

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

	// Security group

	// Create security group
	if err := configurator.ConfigureCreateSecurityGroupStub(securityGroup, securityGroupURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created security group
	if err := configurator.ConfigureGetCreatingSecurityGroupStub(securityGroup, securityGroupURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupStub(securityGroup, securityGroupURL, scenario.MockParams); err != nil {
		return err
	}

	// Update security group (change rules)
	securityGroup = params.SecurityGroupUpdated
	if err := configurator.ConfigureUpdateSecurityGroupStub(securityGroup, securityGroupURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated security group
	if err := configurator.ConfigureGetUpdatingSecurityGroupStub(securityGroup, securityGroupURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupStub(securityGroup, securityGroupURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the security group
	if err := configurator.ConfigureDeleteStub(securityGroupURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingSecurityGroupStub(securityGroup, securityGroupURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(securityGroupURL, scenario.MockParams); err != nil {
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
