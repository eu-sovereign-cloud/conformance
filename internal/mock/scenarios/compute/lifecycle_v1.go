package mockcompute

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"

	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, mockParams *mock.MockParams, suiteParams *params.ComputeLifeCycleV1Params) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	workspace := *suiteParams.Workspace
	blockStorage := *suiteParams.BlockStorage
	instance := *suiteParams.InitialInstance

	configurator, err := stubs.NewStubConfigurator(scenario, mockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, blockStorage.Metadata.Tenant, workspace.Metadata.Name, blockStorage.Metadata.Name)
	instanceUrl := generators.GenerateInstanceURL(constants.ComputeProviderV1, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)
	instanceStartUrl := generators.GenerateInstanceStartURL(constants.ComputeProviderV1, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)
	instanceStopUrl := generators.GenerateInstanceStopURL(constants.ComputeProviderV1, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)
	instanceRestartUrl := generators.GenerateInstanceRestartURL(constants.ComputeProviderV1, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(&workspace, workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(&workspace, workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(&blockStorage, blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(&blockStorage, blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Create an instance
	if err := configurator.ConfigureCreateInstanceStub(&instance, instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instance, instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the instance
	instance.Spec = suiteParams.UpdatedInstance.Spec
	if err := configurator.ConfigureUpdateInstanceStub(&instance, instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instance, instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Stop the instance
	if err := configurator.ConfigureInstanceOperationStub(&instance, instanceStopUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the stopped instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instance, instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Start the instance
	if err := configurator.ConfigureInstanceOperationStub(&instance, instanceStartUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the started instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instance, instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Restart the instance
	if err := configurator.ConfigureInstanceOperationStub(&instance, instanceRestartUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the restarted instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instance, instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the instance
	if err := configurator.ConfigureDeleteStub(instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.ConfigureGetNotFoundStub(instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
