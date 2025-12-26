package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigureComputeLifecycleScenarioV1(scenario string, params *ComputeLifeCycleParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := generators.GenerateInstanceURL(computeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	instanceStartUrl := generators.GenerateInstanceStartURL(computeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	instanceStopUrl := generators.GenerateInstanceStopURL(computeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	instanceRestartUrl := generators.GenerateInstanceRestartURL(computeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.configureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.configureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.configureGetActiveBlockStorageStub(blockResponse, blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(params.Instance.Name).
		Provider(computeProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Instance.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an instance
	if err := configurator.configureCreateInstanceStub(instanceResponse, instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Update the instance
	instanceResponse.Spec = *params.Instance.UpdatedSpec
	if err := configurator.configureUpdateInstanceStub(instanceResponse, instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Stop the instance
	if err := configurator.configureInstanceOperationStub(instanceResponse, instanceStopUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the stopped instance
	if err := configurator.configureGetSuspendedInstanceStub(instanceResponse, instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Start the instance
	if err := configurator.configureInstanceOperationStub(instanceResponse, instanceStartUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the started instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Restart the instance
	if err := configurator.configureInstanceOperationStub(instanceResponse, instanceRestartUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the restarted instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the instance
	if err := configurator.configureDeleteStub(instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.configureGetNotFoundStub(instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.configureGetNotFoundStub(blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	return configurator.client, nil
}

func ConfigureComputeListScenarioV1(scenario string, params *ComputeListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	instanceListUrl := generators.GenerateInstanceListURL(computeProviderV1, params.Tenant, params.Workspace.Name)
	skuListUrl := generators.GenerateInstanceSkuListURL(computeProviderV1, params.Tenant)
	blockListUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)

	// Workspace

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.configureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.configureCreateBlockStorageStub(blockResponse, blockListUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Instance
	instancesList, err := bulkCreateInstancesStubV1(configurator, params.getBaseParams(), params.Workspace.Name, params.Instances)
	if err != nil {
		return nil, err
	}
	instanceResponse, err := builders.NewInstanceIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(instancesList).
		Build()
	if err != nil {
		return nil, err
	}

	// List instances
	if err := configurator.configureGetListInstanceStub(instanceResponse, instanceListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List roles with limit 1
	instanceResponse.Items = instancesList[:1]
	if err := configurator.configureGetListInstanceStub(instanceResponse, instanceListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List instances with label
	instancesWithLabel := func(instancesList []schema.Instance) []schema.Instance {
		var filteredInstances []schema.Instance
		for _, instance := range instancesList {
			if val, ok := instance.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredInstances = append(filteredInstances, instance)
			}
		}
		return filteredInstances
	}
	instanceResponse.Items = instancesWithLabel(instancesList)
	if err := configurator.configureGetListInstanceStub(instanceResponse, instanceListUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List instances with limit and label
	instanceResponse.Items = instancesWithLabel(instancesList)[:1]
	if err := configurator.configureGetListInstanceStub(instanceResponse, instanceListUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create skus
	skusList := generateInstanceSkusV1(params.getBaseParams().Tenant)
	skuResponse, err := builders.NewInstanceSkuIteratorBuilder().Provider(storageProviderV1).Tenant(params.Tenant).Items(skusList).Build()
	if err != nil {
		return nil, err
	}

	// List skus
	if err := configurator.configureGetListSkuStub(skuResponse, skuListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List skus with limit 1
	if err := configurator.configureGetListSkuStub(skuResponse, skuListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List skus with label
	skusWithLabel := func(skusList []schema.InstanceSku) []schema.InstanceSku {
		var filteredSkus []schema.InstanceSku
		for _, sku := range skusList {
			if val, ok := sku.Labels["tier"]; ok && val == "D2XS" {
				filteredSkus = append(filteredSkus, sku)
			}
		}
		return filteredSkus
	}
	skuResponse.Items = skusWithLabel(skusList)
	if err := configurator.configureGetListSkuStub(skuResponse, skuListUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List sku with limit and label
	if err := configurator.configureGetListSkuStub(skuResponse, skuListUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Delete instances
	for _, instance := range instancesList {
		url := generators.GenerateInstanceURL(computeProviderV1, params.Tenant, params.Workspace.Name, instance.Metadata.Name)
		if err := configurator.configureDeleteStub(url, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted instance
		if err := configurator.configureGetNotFoundStub(url, params.getBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockListUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(blockListUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	return configurator.client, nil
}
