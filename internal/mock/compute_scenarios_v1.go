package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	compute "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.compute.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
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
		Build()
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
		Build()
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
	if err := configurator.configureInstanceOperationStub(instanceResponse, instanceStopUrl, params); err != nil {
		return nil, err
	}

	// Get the stopped instance
	if err := configurator.configureGetSuspendedInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Start the instance
	if err := configurator.configureInstanceOperationStub(instanceResponse, instanceStartUrl, params); err != nil {
		return nil, err
	}

	// Get the started instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Restart the instance
	if err := configurator.configureInstanceOperationStub(instanceResponse, instanceRestartUrl, params); err != nil {
		return nil, err
	}

	// Get the restarted instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Delete the instance
	if err := configurator.configureDeleteStub(instanceUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.configureGetNotFoundStub(instanceUrl, params); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.configureGetNotFoundStub(blockUrl, params); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(workspaceUrl, params); err != nil {
		return nil, err
	}

	return configurator.client, nil
}

func ConfigComputeListLifecycleScenarioV1(scenario string, params *ComputeListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)

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
	if err := configurator.configureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
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
	if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Instance
	var instanceList []schema.Instance
	for _, instance := range *params.Instance {
		instanceUrl := generators.GenerateInstanceListURL(computeProviderV1, params.Tenant, instance.Name)
		instanceResponse, err := builders.NewInstanceBuilder().
			Name(instance.Name).
			Provider(computeProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
			Spec(instance.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create an instance
		if err := configurator.configureCreateInstanceStub(instanceResponse, instanceUrl, params); err != nil {
			return nil, err
		}
		instanceList = append(instanceList, *instanceResponse)
	}

	// List Instances
	instanceUrl := generators.GenerateInstanceListURL(computeProviderV1, params.Tenant, params.Workspace.Name)
	instancesResource := generators.GenerateInstanceListResource(params.Tenant, params.Workspace.Name)
	instanceResponse := &compute.InstanceIterator{
		Metadata: schema.ResponseMetadata{
			Provider: computeProviderV1,
			Resource: instancesResource,
			Verb:     http.MethodGet,
		},
	}
	instanceResponse.Items = instanceList

	if err := configurator.configureGetListInstanceStub(instanceResponse, instanceUrl, params, nil); err != nil {
		return nil, err
	}

	// List Roles with limit 1

	instanceResponse.Items = instanceList[:1]
	if err := configurator.configureGetListInstanceStub(instanceResponse, instanceUrl, params, pathParamsLimit("1")); err != nil {
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
	instanceResponse.Items = instancesWithLabel(instanceList)
	if err := configurator.configureGetListInstanceStub(instanceResponse, instanceUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}
	// List instances with limit and label

	instanceResponse.Items = instancesWithLabel(instanceList)[:1]
	if err := configurator.configureGetListInstanceStub(instanceResponse, instanceUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List Sku
	instanceSku := []schema.InstanceSku{
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "D2XS",
				Provider: computeProviderV1,
				Resource: instancesResource,
				Verb:     http.MethodGet,
				Tenant:   params.Tenant,
			},
			Labels: schema.Labels{
				"architecture": "amd64",
				"provider":     "seca",
				"tier":         "D2XS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  1,
				VCPU: 1,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "DXS",
				Provider: computeProviderV1,
				Resource: instancesResource,
				Verb:     http.MethodGet,
				Tenant:   params.Tenant,
			},
			Labels: schema.Labels{
				"architecture": "amd64",
				"provider":     "seca",
				"tier":         "DXS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  2,
				VCPU: 1,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "DS",
				Provider: computeProviderV1,
				Resource: instancesResource,
				Verb:     http.MethodGet,
				Tenant:   params.Tenant,
			},
			Labels: schema.Labels{
				"architecture": "amd64",
				"provider":     "seca",
				"tier":         "DS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  4,
				VCPU: 2,
			},
		},
	}
	skuResource := generators.GenerateInstanceSkuListURL(computeProviderV1, params.Tenant)
	skuUrl := generators.GenerateInstanceSkuListURL(computeProviderV1, params.Tenant)
	skuResponse := &compute.SkuIterator{
		Metadata: schema.ResponseMetadata{
			Provider: computeProviderV1,
			Resource: skuResource,
			Verb:     http.MethodGet,
		},
		Items: instanceSku,
	}

	if err := configurator.configureGetListSkuStub(skuResponse, skuUrl, params, nil); err != nil {
		return nil, err
	}

	// List Sku with limit 1

	if err := configurator.configureGetListSkuStub(skuResponse, skuUrl, params, pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List Sku with label

	skusWithLabel := func(skusList []schema.InstanceSku) []schema.InstanceSku {
		var filteredSkus []schema.InstanceSku
		for _, sku := range skusList {
			if val, ok := sku.Labels["tier"]; ok && val == "D2XS" {
				filteredSkus = append(filteredSkus, sku)
			}
		}
		return filteredSkus
	}
	skuResponse.Items = skusWithLabel(instanceSku)
	if err := configurator.configureGetListSkuStub(skuResponse, skuUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List Sku with limit and label
	if err := configurator.configureGetListSkuStub(skuResponse, skuUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Delete Instance
	for _, instance := range instanceList {
		url := generators.GenerateInstanceURL(computeProviderV1, params.Tenant, params.Workspace.Name, instance.Metadata.Name)
		if err := configurator.configureDeleteStub(url, params); err != nil {
			return nil, err
		}

		// Get the deleted instance
		if err := configurator.configureGetNotFoundStub(url, params); err != nil {
			return nil, err
		}
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(blockUrl, params); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(workspaceUrl, params); err != nil {
		return nil, err
	}

	return configurator.client, nil
}
