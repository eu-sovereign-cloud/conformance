package compute

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigureListScenarioV1(scenario string, params *mock.ComputeListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := stubs.NewScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, params.Tenant, params.Workspace.Name)
	instanceListUrl := generators.GenerateInstanceListURL(constants.ComputeProviderV1, params.Tenant, params.Workspace.Name)
	skuListUrl := generators.GenerateInstanceSkuListURL(constants.ComputeProviderV1, params.Tenant)
	blockListUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)

	// Workspace

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
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

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockListUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Instance
	instancesList, err := stubs.BulkCreateInstancesStubV1(configurator, params.GetBaseParams(), params.Workspace.Name, params.Instances)
	if err != nil {
		return nil, err
	}
	instanceResponse, err := builders.NewInstanceIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(instancesList).
		Build()
	if err != nil {
		return nil, err
	}

	// List instances
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List roles with limit 1
	instanceResponse.Items = instancesList[:1]
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
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
	instanceResponse.Items = instancesWithLabel(instancesList)
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, params.GetBaseParams(), mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List instances with limit and label
	instanceResponse.Items = instancesWithLabel(instancesList)[:1]
	if err := configurator.ConfigureGetListInstanceStub(instanceResponse, instanceListUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create skus
	skusList := steps.GenerateInstanceSkusV1(params.GetBaseParams().Tenant)
	skuResponse, err := builders.NewInstanceSkuIteratorBuilder().Provider(constants.StorageProviderV1).Tenant(params.Tenant).Items(skusList).Build()
	if err != nil {
		return nil, err
	}

	// List skus
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List skus with limit 1
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
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
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, params.GetBaseParams(), mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List sku with limit and label
	if err := configurator.ConfigureGetListSkuStub(skuResponse, skuListUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Delete instances
	for _, instance := range instancesList {
		url := generators.GenerateInstanceURL(constants.ComputeProviderV1, params.Tenant, params.Workspace.Name, instance.Metadata.Name)
		if err := configurator.ConfigureDeleteStub(url, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted instance
		if err := configurator.ConfigureGetNotFoundStub(url, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockListUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(blockListUrl, params.GetBaseParams()); err != nil {
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
