package network

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/conformance/pkg/mock/scenarios"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureNicConstraintsValidationV1(scenario *scenarios.Scenario, p params.NicConstraintsValidationV1Params) error {
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
	overLengthNameURL := generators.GenerateNicURL(
		sdkconsts.NetworkProviderV1Name,
		p.OverLengthNameNic.Metadata.Tenant,
		p.OverLengthNameNic.Metadata.Workspace,
		p.OverLengthNameNic.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name violation
	invalidPatternNameURL := generators.GenerateNicURL(
		sdkconsts.NetworkProviderV1Name,
		p.InvalidPatternNameNic.Metadata.Tenant,
		p.InvalidPatternNameNic.Metadata.Workspace,
		p.InvalidPatternNameNic.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value violation
	overLengthLabelURL := generators.GenerateNicURL(
		sdkconsts.NetworkProviderV1Name,
		p.OverLengthLabelValueNic.Metadata.Tenant,
		p.OverLengthLabelValueNic.Metadata.Workspace,
		p.OverLengthLabelValueNic.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value violation
	overLengthAnnotationURL := generators.GenerateNicURL(
		sdkconsts.NetworkProviderV1Name,
		p.OverLengthAnnotationNic.Metadata.Tenant,
		p.OverLengthAnnotationNic.Metadata.Workspace,
		p.OverLengthAnnotationNic.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length address violation
	overLengthAddressURL := generators.GenerateNicURL(
		sdkconsts.NetworkProviderV1Name,
		p.OverLengthAddressNic.Metadata.Tenant,
		p.OverLengthAddressNic.Metadata.Workspace,
		p.OverLengthAddressNic.Metadata.Name,
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
