package mockstorage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureBlockStorageConstraintsV1 sets up mock stubs for the block storage
// constraints violations suite. Creates a valid workspace environment before testing
// violations, all invalid block storage requests returning 422 Unprocessable Entity.
func ConfigureBlockStorageConstraintsV1(scenario *mockscenarios.Scenario, p params.BlockStorageConstraintsValidationV1Params) error {
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
	overLengthNameURL := generators.GenerateBlockStorageURL(
		sdkconsts.StorageProviderV1Name,
		p.OverLengthNameBlockStorage.Metadata.Tenant,
		p.OverLengthNameBlockStorage.Metadata.Workspace,
		p.OverLengthNameBlockStorage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name violation
	invalidPatternNameURL := generators.GenerateBlockStorageURL(
		sdkconsts.StorageProviderV1Name,
		p.InvalidPatternNameBlockStorage.Metadata.Tenant,
		p.InvalidPatternNameBlockStorage.Metadata.Workspace,
		p.InvalidPatternNameBlockStorage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value violation
	overLengthLabelURL := generators.GenerateBlockStorageURL(
		sdkconsts.StorageProviderV1Name,
		p.OverLengthLabelValueBlockStorage.Metadata.Tenant,
		p.OverLengthLabelValueBlockStorage.Metadata.Workspace,
		p.OverLengthLabelValueBlockStorage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value violation
	overLengthAnnotationURL := generators.GenerateBlockStorageURL(
		sdkconsts.StorageProviderV1Name,
		p.OverLengthAnnotationBlockStorage.Metadata.Tenant,
		p.OverLengthAnnotationBlockStorage.Metadata.Workspace,
		p.OverLengthAnnotationBlockStorage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
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
