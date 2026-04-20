package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureNicConstraintsViolationsV1(scenario *mockscenarios.Scenario, p params.NicConstraintsViolationsV1Params) error {
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
		generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, p.OverLengthNameNic.Metadata.Tenant, p.OverLengthNameNic.Metadata.Workspace, p.OverLengthNameNic.Metadata.Name),
		generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, p.InvalidPatternNameNic.Metadata.Tenant, p.InvalidPatternNameNic.Metadata.Workspace, p.InvalidPatternNameNic.Metadata.Name),
		generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, p.OverLengthLabelValueNic.Metadata.Tenant, p.OverLengthLabelValueNic.Metadata.Workspace, p.OverLengthLabelValueNic.Metadata.Name),
		generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, p.OverLengthAnnotationNic.Metadata.Tenant, p.OverLengthAnnotationNic.Metadata.Workspace, p.OverLengthAnnotationNic.Metadata.Name),
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
