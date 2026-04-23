package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureSecurityGroupConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.SecurityGroupConstraintsValidationV1Params) error {
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

	// Over-length name Validation
	overLengthNameSecurityGroup := p.OverLengthNameSecurityGroup
	overLengthNameURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, overLengthNameSecurityGroup.Metadata.Tenant, overLengthNameSecurityGroup.Metadata.Workspace, overLengthNameSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name Validation
	invalidPatternNameSecurityGroup := p.InvalidPatternNameSecurityGroup
	invalidPatternNameURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, invalidPatternNameSecurityGroup.Metadata.Tenant, invalidPatternNameSecurityGroup.Metadata.Workspace, invalidPatternNameSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value Validation
	overLengthLabelSecurityGroup := p.OverLengthLabelValueSecurityGroup
	overLengthLabelURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, overLengthLabelSecurityGroup.Metadata.Tenant, overLengthLabelSecurityGroup.Metadata.Workspace, overLengthLabelSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value Validation
	overLengthAnnotationSecurityGroup := p.OverLengthAnnotationSecurityGroup
	overLengthAnnotationURL := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, overLengthAnnotationSecurityGroup.Metadata.Tenant, overLengthAnnotationSecurityGroup.Metadata.Workspace, overLengthAnnotationSecurityGroup.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
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
