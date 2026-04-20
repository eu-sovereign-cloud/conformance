package mockcompute

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureInstanceConstraintsViolationsV1 sets up mock stubs for the instance constraints
// violations suite. Creates a valid workspace + block storage environment before testing
// violations, all invalid instance requests returning 422 Unprocessable Entity.
func ConfigureInstanceConstraintsViolationsV1(scenario *mockscenarios.Scenario, p params.InstanceConstraintsViolationsV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	blockStorage := p.BlockStorage

	// Generate URLs
	workspaceURL := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockURL := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockStorage.Metadata.Tenant, workspace.Metadata.Name, blockStorage.Metadata.Name)

	// Create workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceURL, scenario.MockParams); err != nil {
		return err
	}

	// Create block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetCreatingBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockStorage, blockURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length name violation
	overLengthNameURL := generators.GenerateInstanceURL(
		sdkconsts.ComputeProviderV1Name,
		p.OverLengthNameInstance.Metadata.Tenant,
		p.OverLengthNameInstance.Metadata.Workspace,
		p.OverLengthNameInstance.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name violation
	invalidPatternNameURL := generators.GenerateInstanceURL(
		sdkconsts.ComputeProviderV1Name,
		p.InvalidPatternNameInstance.Metadata.Tenant,
		p.InvalidPatternNameInstance.Metadata.Workspace,
		p.InvalidPatternNameInstance.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value violation
	overLengthLabelURL := generators.GenerateInstanceURL(
		sdkconsts.ComputeProviderV1Name,
		p.OverLengthLabelValueInstance.Metadata.Tenant,
		p.OverLengthLabelValueInstance.Metadata.Workspace,
		p.OverLengthLabelValueInstance.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value violation
	overLengthAnnotationURL := generators.GenerateInstanceURL(
		sdkconsts.ComputeProviderV1Name,
		p.OverLengthAnnotationInstance.Metadata.Tenant,
		p.OverLengthAnnotationInstance.Metadata.Workspace,
		p.OverLengthAnnotationInstance.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
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
