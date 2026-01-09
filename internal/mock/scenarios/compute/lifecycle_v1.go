package mockcompute

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"

	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, params *params.ComputeLifeCycleParamsV1) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	workspaceResponse := *params.Workspace
	blockStorageResponse := *params.BlockStorage
	instanceResponse := *params.CreatedInstance

	configurator, err := stubs.NewStubConfigurator(scenario, params.MockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, params.Tenant, params.Workspace.Metadata.Name)
	blockUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, params.Tenant, params.Workspace.Metadata.Name, params.BlockStorage.Metadata.Name)
	instanceUrl := generators.GenerateInstanceURL(constants.ComputeProviderV1, params.Tenant, params.Workspace.Metadata.Name, params.CreatedInstance.Metadata.Name)
	instanceStartUrl := generators.GenerateInstanceStartURL(constants.ComputeProviderV1, params.Tenant, params.Workspace.Metadata.Name, params.CreatedInstance.Metadata.Name)
	instanceStopUrl := generators.GenerateInstanceStopURL(constants.ComputeProviderV1, params.Tenant, params.Workspace.Metadata.Name, params.CreatedInstance.Metadata.Name)
	instanceRestartUrl := generators.GenerateInstanceRestartURL(constants.ComputeProviderV1, params.Tenant, params.Workspace.Metadata.Name, params.CreatedInstance.Metadata.Name)

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(&workspaceResponse, workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(&workspaceResponse, workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(&blockStorageResponse, blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(&blockStorageResponse, blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Create an instance
	if err := configurator.ConfigureCreateInstanceStub(&instanceResponse, instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instanceResponse, instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the instance
	instanceResponse.Spec = params.UpdatedInstance.Spec
	if err := configurator.ConfigureUpdateInstanceStub(&instanceResponse, instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instanceResponse, instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Stop the instance
	if err := configurator.ConfigureInstanceOperationStub(&instanceResponse, instanceStopUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the stopped instance
	if err := configurator.ConfigureGetSuspendedInstanceStub(&instanceResponse, instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Start the instance
	if err := configurator.ConfigureInstanceOperationStub(&instanceResponse, instanceStartUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the started instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instanceResponse, instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Restart the instance
	if err := configurator.ConfigureInstanceOperationStub(&instanceResponse, instanceRestartUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the restarted instance
	if err := configurator.ConfigureGetActiveInstanceStub(&instanceResponse, instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the instance
	if err := configurator.ConfigureDeleteStub(instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.ConfigureGetNotFoundStub(instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
