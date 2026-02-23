package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
)

func ConfigureSecurityGroupLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.SecurityGroupLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := *params.Workspace
	securityGroupInitial := *params.SecurityGroupInitial
	securityGroupUpdated := *params.SecurityGroupUpdated

	workspaceURL := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	securityGroupURL := generators.GenerateSecurityGroupURL(constants.NetworkProviderV1, securityGroupInitial.Metadata.Tenant, securityGroupInitial.Metadata.Workspace, securityGroupInitial.Metadata.Name)

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

	// Security group
	securityGroupInitialResponse, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupInitial.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(securityGroupInitial.Metadata.Tenant).Workspace(securityGroupInitial.Metadata.Workspace).Region(securityGroupInitial.Metadata.Region).
		Spec(&securityGroupInitial.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create security group
	if err := configurator.ConfigureCreateSecurityGroupStub(securityGroupInitialResponse, securityGroupURL, scenario.MockParams); err != nil {
		return err
	}

	// Get created security group
	if err := configurator.ConfigureGetActiveSecurityGroupStub(securityGroupInitialResponse, securityGroupURL, scenario.MockParams); err != nil {
		return err
	}

	// Update security group (change rules)
	securityGroupUpdatedResponse, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupUpdated.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(securityGroupUpdated.Metadata.Tenant).Workspace(securityGroupUpdated.Metadata.Workspace).Region(securityGroupUpdated.Metadata.Region).
		Spec(&securityGroupUpdated.Spec).
		Build()
	if err != nil {
		return err
	}

	if err := configurator.ConfigureUpdateSecurityGroupStub(securityGroupUpdatedResponse, securityGroupURL, scenario.MockParams); err != nil {
		return err
	}

	// Get updated security group
	if err := configurator.ConfigureGetActiveSecurityGroupStub(securityGroupUpdatedResponse, securityGroupURL, scenario.MockParams); err != nil {
		return err
	}

	// Deletions
	if err := configurator.ConfigureDeleteStub(securityGroupURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(securityGroupURL, scenario.MockParams); err != nil {
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
