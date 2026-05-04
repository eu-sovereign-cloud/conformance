package mockstorage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureImageConstraintsViolationsV1 sets up mock stubs for the image constraints
// violations suite. Each image in the params targets a different constraint violation,
// all returning 422 Unprocessable Entity.
func ConfigureImageConstraintsViolationsV1(scenario *mockscenarios.Scenario, params params.ImageConstraintsValidationV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := params.Workspace
	blockStorage := params.BlockStorage

	workspaceUrl := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockUrl := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockStorage.Metadata.Tenant, blockStorage.Metadata.Workspace, blockStorage.Metadata.Name)

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetCreatingBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Over-length name violation
	overLengthNameURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		params.OverLengthNameImage.Metadata.Tenant,
		params.OverLengthNameImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name violation
	invalidPatternNameURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		params.InvalidPatternNameImage.Metadata.Tenant,
		params.InvalidPatternNameImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value violation
	overLengthLabelURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		params.OverLengthLabelValueImage.Metadata.Tenant,
		params.OverLengthLabelValueImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value violation
	overLengthAnnotationURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		params.OverLengthAnnotationImage.Metadata.Tenant,
		params.OverLengthAnnotationImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
