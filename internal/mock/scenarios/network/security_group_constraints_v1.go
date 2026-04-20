package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureSecurityGroupConstraintsViolationsV1(scenario *mockscenarios.Scenario, p params.SecurityGroupConstraintsViolationsV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)

	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	for _, url := range []string{
		generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, p.OverLengthNameSecurityGroup.Metadata.Tenant, p.OverLengthNameSecurityGroup.Metadata.Workspace, p.OverLengthNameSecurityGroup.Metadata.Name),
		generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, p.InvalidPatternNameSecurityGroup.Metadata.Tenant, p.InvalidPatternNameSecurityGroup.Metadata.Workspace, p.InvalidPatternNameSecurityGroup.Metadata.Name),
		generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, p.OverLengthLabelValueSecurityGroup.Metadata.Tenant, p.OverLengthLabelValueSecurityGroup.Metadata.Workspace, p.OverLengthLabelValueSecurityGroup.Metadata.Name),
		generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, p.OverLengthAnnotationSecurityGroup.Metadata.Tenant, p.OverLengthAnnotationSecurityGroup.Metadata.Workspace, p.OverLengthAnnotationSecurityGroup.Metadata.Name),
	} {
		if err := configurator.ConfigurePutUnprocessableEntityStub(url, scenario.MockParams); err != nil {
			return err
		}
	}

	if err := configurator.ConfigureDeleteStub(workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
