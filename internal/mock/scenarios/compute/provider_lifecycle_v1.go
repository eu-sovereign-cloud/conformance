package mockcompute

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
)

func ConfigureProviderLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.ComputeProviderLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}
	workspace := *params.Workspace
	blockStorage := *params.BlockStorage
	instance := *params.InitialInstance

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, blockStorage.Metadata.Tenant, workspace.Metadata.Name, blockStorage.Metadata.Name)
	instanceUrl := generators.GenerateInstanceURL(constants.ComputeProviderV1, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)
	instanceStartUrl := generators.GenerateInstanceStartURL(constants.ComputeProviderV1, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)
	instanceStopUrl := generators.GenerateInstanceStopURL(constants.ComputeProviderV1, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)
	instanceRestartUrl := generators.GenerateInstanceRestartURL(constants.ComputeProviderV1, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(&workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(&workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(&blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(&blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create an instance
	if err := configurator.ConfigureCreateInstanceStub(&instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the instance
	instance.Spec = params.UpdatedInstance.Spec
	if err := configurator.ConfigureUpdateInstanceStub(&instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Stop the instance
	if err := configurator.ConfigureInstanceOperationStub(&instance, instanceStopUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the stopped instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Start the instance
	if err := configurator.ConfigureInstanceOperationStub(&instance, instanceStartUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the started instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Restart the instance
	if err := configurator.ConfigureInstanceOperationStub(&instance, instanceRestartUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the restarted instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the instance
	if err := configurator.ConfigureDeleteStub(instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted instance
	if err := configurator.ConfigureGetNotFoundStub(instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
