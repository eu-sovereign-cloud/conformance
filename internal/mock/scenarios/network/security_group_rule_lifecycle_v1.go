package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureSecurityGroupRuleLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.SecurityGroupRuleLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := *params.Workspace
	securityGroupRuleInitial := *params.SecurityGroupRuleInitial
	securityGroupRuleUpdated := *params.SecurityGroupRuleUpdated

	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	securityGroupRuleURL := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name,
		securityGroupRuleInitial.Metadata.Tenant,
		securityGroupRuleInitial.Metadata.Workspace,
		securityGroupRuleInitial.Metadata.Name,
	)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(workspace.Metadata.Name).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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

	// Get the created workspace
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspaceResponse, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Security group rule
	securityGroupRuleInitialResponse, err := builders.NewSecurityGroupRuleBuilder().
		Name(securityGroupRuleInitial.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(securityGroupRuleInitial.Metadata.Tenant).Workspace(securityGroupRuleInitial.Metadata.Workspace).Region(securityGroupRuleInitial.Metadata.Region).
		Spec(&securityGroupRuleInitial.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a security group rule
	if err := configurator.ConfigureCreateSecurityGroupRuleStub(securityGroupRuleInitialResponse, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created security group rule
	if err := configurator.ConfigureGetCreatingSecurityGroupRuleStub(securityGroupRuleInitialResponse, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupRuleStub(securityGroupRuleInitialResponse, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}

	// Update the security group rule
	securityGroupRuleUpdatedResponse, err := builders.NewSecurityGroupRuleBuilder().
		Name(securityGroupRuleUpdated.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(securityGroupRuleUpdated.Metadata.Tenant).Workspace(securityGroupRuleUpdated.Metadata.Workspace).Region(securityGroupRuleUpdated.Metadata.Region).
		Spec(&securityGroupRuleUpdated.Spec).
		Build()
	if err != nil {
		return err
	}

	if err := configurator.ConfigureUpdateSecurityGroupRuleStub(securityGroupRuleUpdatedResponse, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated security group rule
	if err := configurator.ConfigureGetUpdatingSecurityGroupRuleStub(securityGroupRuleUpdatedResponse, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupRuleStub(securityGroupRuleUpdatedResponse, securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the security group rule
	if err := configurator.ConfigureDeleteStub(securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(securityGroupRuleURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete the workspace
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
