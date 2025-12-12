package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/wiremock/go-wiremock"
)

func ConfigComputeLifecycleScenarioV1(scenario string, params *ComputeParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := generators.GenerateInstanceURL(computeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.configureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.configureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.configureGetActiveBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(params.Instance.Name).
		Provider(computeProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Instance.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create an instance
	if err := configurator.configureCreateInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Update the instance
	instanceResponse.Spec = *params.Instance.UpdatedSpec
	if err := configurator.configureUpdateInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Get the updated instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Stop the instance
	if err := configurator.configureInstanceOperationStub(instanceResponse, instanceUrl+"/stop", params); err != nil {
		return nil, err
	}

	// Get the stopped instance
	if err := configurator.configureGetSuspendedInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Start the instance
	if err := configurator.configureInstanceOperationStub(instanceResponse, instanceUrl+"/start", params); err != nil {
		return nil, err
	}

	// Get the started instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Restart the instance
	if err := configurator.configureInstanceOperationStub(instanceResponse, instanceUrl+"/restart", params); err != nil {
		return nil, err
	}

	// Get the restarted instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Delete the instance
	if err := configurator.configureDeleteStub(instanceUrl, params, false); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.configureGetNotFoundStub(instanceUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockUrl, params, false); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.configureGetNotFoundStub(blockUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params, false); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(workspaceUrl, params, true); err != nil {
		return nil, err
	}

	return configurator.client, nil
}
