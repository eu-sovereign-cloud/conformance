package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigurePublicIpConstraintsViolationsV1(scenario *mockscenarios.Scenario, p params.PublicIpConstraintsViolationsV1Params) error {
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
		generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, p.OverLengthNamePublicIp.Metadata.Tenant, p.OverLengthNamePublicIp.Metadata.Workspace, p.OverLengthNamePublicIp.Metadata.Name),
		generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, p.InvalidPatternNamePublicIp.Metadata.Tenant, p.InvalidPatternNamePublicIp.Metadata.Workspace, p.InvalidPatternNamePublicIp.Metadata.Name),
		generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, p.OverLengthLabelValuePublicIp.Metadata.Tenant, p.OverLengthLabelValuePublicIp.Metadata.Workspace, p.OverLengthLabelValuePublicIp.Metadata.Name),
		generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, p.OverLengthAnnotationPublicIp.Metadata.Tenant, p.OverLengthAnnotationPublicIp.Metadata.Workspace, p.OverLengthAnnotationPublicIp.Metadata.Name),
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
