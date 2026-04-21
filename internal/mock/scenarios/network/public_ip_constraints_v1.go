package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigurePublicIpConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.PublicIpConstraintsValidationV1Params) error {
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
	overLengthNamePublicIp := p.OverLengthNamePublicIp
	overLengthNameURL := generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, overLengthNamePublicIp.Metadata.Tenant, overLengthNamePublicIp.Metadata.Workspace, overLengthNamePublicIp.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name Validation
	invalidPatternNamePublicIp := p.InvalidPatternNamePublicIp
	invalidPatternNameURL := generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, invalidPatternNamePublicIp.Metadata.Tenant, invalidPatternNamePublicIp.Metadata.Workspace, invalidPatternNamePublicIp.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value Validation
	overLengthLabelPublicIp := p.OverLengthLabelValuePublicIp
	overLengthLabelURL := generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, overLengthLabelPublicIp.Metadata.Tenant, overLengthLabelPublicIp.Metadata.Workspace, overLengthLabelPublicIp.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value Validation
	overLengthAnnotationPublicIp := p.OverLengthAnnotationPublicIp
	overLengthAnnotationURL := generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, overLengthAnnotationPublicIp.Metadata.Tenant, overLengthAnnotationPublicIp.Metadata.Workspace, overLengthAnnotationPublicIp.Metadata.Name)
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
