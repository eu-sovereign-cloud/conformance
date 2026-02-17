package mockcompute

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/wiremock/go-wiremock"
)

func ConfigureListScenarioV1(scenario string, mockParams *mock.MockParams, suiteParams *params.ComputeListV1Params) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	workspace := suiteParams.Workspace
	instances := suiteParams.Instances
	blockStorage := suiteParams.BlockStorage

	configurator, err := stubs.NewStubConfigurator(scenario, mockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	instanceListUrl := generators.GenerateInstanceListURL(constants.ComputeProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	skuListUrl := generators.GenerateInstanceSkuListURL(constants.ComputeProviderV1, workspace.Metadata.Tenant)
	blockUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, blockStorage.Metadata.Tenant, blockStorage.Metadata.Workspace, blockStorage.Metadata.Name)

	// Workspace

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(workspace.Metadata.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(workspace.Metadata.Tenant).Region(workspace.Metadata.Region).
		Labels(workspace.Labels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(blockStorage.Metadata.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(blockStorage.Metadata.Tenant).Workspace(blockStorage.Metadata.Workspace).Region(blockStorage.Metadata.Region).
		Spec(&blockStorage.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Instance
	err = stubs.BulkCreateInstancesStubV1(configurator, mockParams, instances)
	if err != nil {
		return nil, err
	}
	instanceResponse, err := builders.NewInstanceIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(workspace.Metadata.Tenant).Workspace(workspace.Metadata.Name).
		Items(instances).
		Build()
	if err != nil {
		return nil, err
	}

	// List instances
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, mockParams, nil); err != nil {
		return nil, err
	}

	// List roles with limit 1
	instanceResponse.Items = instances[:1]
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, mockParams, mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List instances with label
	instancesWithLabel := func(instancesList []schema.Instance) []schema.Instance {
		var filteredInstances []schema.Instance
		for _, instance := range instancesList {
			if val, ok := instance.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredInstances = append(filteredInstances, instance)
			}
		}
		return filteredInstances
	}
	instanceResponse.Items = instancesWithLabel(instances)
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, mockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List instances with limit and label
	instanceResponse.Items = instancesWithLabel(instances)[:1]
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, mockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create skus
	skusList := mock.GenerateInstanceSkusV1(workspace.Metadata.Tenant)
	skuResponse, err := builders.NewInstanceSkuIteratorBuilder().Provider(constants.StorageProviderV1).Tenant(workspace.Metadata.Tenant).Items(skusList).Build()
	if err != nil {
		return nil, err
	}

	// List skus
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, mockParams, nil); err != nil {
		return nil, err
	}

	// List skus with limit 1
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, mockParams, mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List skus with label
	skusWithLabel := func(skusList []schema.InstanceSku) []schema.InstanceSku {
		var filteredSkus []schema.InstanceSku
		for _, sku := range skusList {
			if val, ok := sku.Labels[constants.TierLabel]; ok && val == constants.TierSkuD2XSLabel {
				filteredSkus = append(filteredSkus, sku)
			}
		}
		return filteredSkus
	}
	skuResponse.Items = skusWithLabel(skusList)
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, mockParams, mock.PathParamsLabel(constants.TierLabel, constants.TierSkuD2XSLabel)); err != nil {
		return nil, err
	}

	// List sku with limit and label
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, mockParams, mock.PathParamsLimitAndLabel("1", constants.TierLabel, constants.TierSkuD2XSLabel)); err != nil {
		return nil, err
	}

	// Delete instances
	for _, instance := range instances {
		url := generators.GenerateInstanceURL(constants.ComputeProviderV1, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)
		if err := configurator.ConfigureDeleteStub(url, mockParams); err != nil {
			return nil, err
		}

		// Get the deleted instance
		if err := configurator.ConfigureGetNotFoundStub(url, mockParams); err != nil {
			return nil, err
		}
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted workspace
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

	// Finish the stubs configuration
	if client, err := configurator.Finish(); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}
