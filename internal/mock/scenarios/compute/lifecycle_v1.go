package compute

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, params *mock.ComputeLifeCycleParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := stubs.NewScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(mock.WorkspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(mock.StorageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := generators.GenerateInstanceURL(mock.ComputeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	instanceStartUrl := generators.GenerateInstanceStartURL(mock.ComputeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	instanceStopUrl := generators.GenerateInstanceStopURL(mock.ComputeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	instanceRestartUrl := generators.GenerateInstanceRestartURL(mock.ComputeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(mock.WorkspaceProviderV1).ApiVersion(mock.ApiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(mock.StorageProviderV1).ApiVersion(mock.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(params.Instance.Name).
		Provider(mock.ComputeProviderV1).ApiVersion(mock.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Instance.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an instance
	if err := configurator.ConfigureCreateInstanceStub(instanceResponse, instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.ConfigureGetActiveInstanceStub(instanceResponse, instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Update the instance
	instanceResponse.Spec = *params.Instance.UpdatedSpec
	if err := configurator.ConfigureUpdateInstanceStub(instanceResponse, instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated instance
	if err := configurator.ConfigureGetActiveInstanceStub(instanceResponse, instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Stop the instance
	if err := configurator.ConfigureInstanceOperationStub(instanceResponse, instanceStopUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the stopped instance
	if err := configurator.ConfigureGetSuspendedInstanceStub(instanceResponse, instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Start the instance
	if err := configurator.ConfigureInstanceOperationStub(instanceResponse, instanceStartUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the started instance
	if err := configurator.ConfigureGetActiveInstanceStub(instanceResponse, instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Restart the instance
	if err := configurator.ConfigureInstanceOperationStub(instanceResponse, instanceRestartUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the restarted instance
	if err := configurator.ConfigureGetActiveInstanceStub(instanceResponse, instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the instance
	if err := configurator.ConfigureDeleteStub(instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.ConfigureGetNotFoundStub(instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
