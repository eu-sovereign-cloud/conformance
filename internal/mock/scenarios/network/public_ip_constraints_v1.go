package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/pkg/mock/scenarios"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigurePublicIpConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.PublicIpConstraintsValidationV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)

	// Create workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length name violation
	overLengthNameURL := generators.GeneratePublicIpURL(
		sdkconsts.NetworkProviderV1Name,
		p.OverLengthNamePublicIp.Metadata.Tenant,
		p.OverLengthNamePublicIp.Metadata.Workspace,
		p.OverLengthNamePublicIp.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name violation
	invalidPatternNameURL := generators.GeneratePublicIpURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidPatternNamePublicIp.Metadata.Tenant,
		p.InvalidPatternNamePublicIp.Metadata.Workspace,
		p.InvalidPatternNamePublicIp.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value violation
	overLengthLabelURL := generators.GeneratePublicIpURL(
		sdkconsts.NetworkProviderV1Name,
		p.OverLengthLabelValuePublicIp.Metadata.Tenant,
		p.OverLengthLabelValuePublicIp.Metadata.Workspace,
		p.OverLengthLabelValuePublicIp.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value violation
	overLengthAnnotationURL := generators.GeneratePublicIpURL(
		sdkconsts.NetworkProviderV1Name,
		p.OverLengthAnnotationPublicIp.Metadata.Tenant,
		p.OverLengthAnnotationPublicIp.Metadata.Workspace,
		p.OverLengthAnnotationPublicIp.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length address violation
	overLengthAddressURL := generators.GeneratePublicIpURL(
		sdkconsts.NetworkProviderV1Name,
		p.OverLengthAddressPublicIp.Metadata.Tenant,
		p.OverLengthAddressPublicIp.Metadata.Workspace,
		p.OverLengthAddressPublicIp.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAddressURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete workspace teardown
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
