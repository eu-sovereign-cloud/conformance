package storage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/conformance/pkg/mock/scenarios"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureImageConstraintsViolationsV1 sets up mock stubs for the image constraints
// violations suite. Each image in the params targets a different constraint violation,
// all returning 422 Unprocessable Entity.
func ConfigureImageConstraintsViolationsV1(scenario *scenarios.Scenario, p params.ImageConstraintsValidationV1Params) error {
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

	// Create block storage
	blockStorage := p.BlockStorage
	blockURL := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockStorage.Metadata.Tenant, blockStorage.Metadata.Workspace, blockStorage.Metadata.Name)

	if err := configurator.ConfigureCreateBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length name violation
	overLengthNameURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		p.OverLengthNameImage.Metadata.Tenant,
		p.OverLengthNameImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name violation
	invalidPatternNameURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		p.InvalidPatternNameImage.Metadata.Tenant,
		p.InvalidPatternNameImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value violation
	overLengthLabelURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		p.OverLengthLabelValueImage.Metadata.Tenant,
		p.OverLengthLabelValueImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value violation
	overLengthAnnotationURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		p.OverLengthAnnotationImage.Metadata.Tenant,
		p.OverLengthAnnotationImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid cpuArchitecture enum violation
	invalidCpuArchURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		p.InvalidCpuArchitectureImage.Metadata.Tenant,
		p.InvalidCpuArchitectureImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidCpuArchURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid initializer enum violation
	invalidInitializerURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		p.InvalidInitializerImage.Metadata.Tenant,
		p.InvalidInitializerImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidInitializerURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid boot enum violation
	invalidBootURL := generators.GenerateImageURL(
		sdkconsts.StorageProviderV1Name,
		p.InvalidBootImage.Metadata.Tenant,
		p.InvalidBootImage.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidBootURL, scenario.MockParams); err != nil {
		return err
	}

	// Delete block storage teardown
	if err := configurator.ConfigureDeleteStub(blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetDeletingBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(blockURL, scenario.MockParams); err != nil {
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
