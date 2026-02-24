package mockcompute

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func ConfigureProviderQueriesV1(scenario *mockscenarios.Scenario, params *params.ComputeProviderQueriesV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := params.Workspace
	instances := params.Instances
	blockStorage := params.BlockStorage

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
		return err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(blockStorage.Metadata.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(blockStorage.Metadata.Tenant).Workspace(blockStorage.Metadata.Workspace).Region(blockStorage.Metadata.Region).
		Spec(&blockStorage.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Instance
	err = stubs.BulkCreateInstancesStubV1(configurator, scenario.MockParams, instances)
	if err != nil {
		return err
	}
	instanceResponse, err := builders.NewInstanceIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(workspace.Metadata.Tenant).Workspace(workspace.Metadata.Name).
		Items(instances).
		Build()
	if err != nil {
		return err
	}

	// List instances
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List roles with limit 1
	instanceResponse.Items = instances[:1]
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
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
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List instances with limit and label
	instanceResponse.Items = instancesWithLabel(instances)[:1]
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Create skus
	skusList := steps.GenerateInstanceSkusV1(workspace.Metadata.Tenant)
	skuResponse, err := builders.NewInstanceSkuIteratorBuilder().Provider(constants.StorageProviderV1).Tenant(workspace.Metadata.Tenant).Items(skusList).Build()
	if err != nil {
		return err
	}

	// List skus
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List skus with limit 1
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
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
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, scenario.MockParams, mock.PathParamsLabel(constants.TierLabel, constants.TierSkuD2XSLabel)); err != nil {
		return err
	}

	// List sku with limit and label
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.TierLabel, constants.TierSkuD2XSLabel)); err != nil {
		return err
	}

	// Delete instances
	for _, instance := range instances {
		url := generators.GenerateInstanceURL(constants.ComputeProviderV1, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)
		if err := configurator.ConfigureDeleteStub(url, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted instance
		if err := configurator.ConfigureGetNotFoundStub(url, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted workspace
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
