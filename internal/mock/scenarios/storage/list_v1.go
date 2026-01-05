package storage

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigureListScenarioV1(scenario string, params *mock.StorageListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := stubs.NewScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(mock.WorkspaceProviderV1, params.Tenant, params.Workspace.Name)
	skuListUrl := generators.GenerateStorageSkuListURL(mock.StorageProviderV1, params.Tenant)
	blockListUrl := generators.GenerateBlockStorageListURL(mock.StorageProviderV1, params.Tenant, params.Workspace.Name)
	imageListUrl := generators.GenerateImageListURL(mock.StorageProviderV1, params.Tenant)

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

	// Create block storages
	blockList, err := stubs.BulkCreateBlockStoragesStubV1(configurator, params.GetBaseParams(), params.Workspace.Name, params.BlockStorages)
	if err != nil {
		return nil, err
	}
	blockListResponse, err := builders.NewBlockStorageIteratorBuilder().
		Provider(mock.StorageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(blockList).
		Build()
	if err != nil {
		return nil, err
	}

	// List block storages
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	blockListResponse.Items = blockList[:1]
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	blocksWithLabel := func(blockList []schema.BlockStorage) []schema.BlockStorage {
		var filteredInstances []schema.BlockStorage
		for _, instance := range blockList {
			if val, ok := instance.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredInstances = append(filteredInstances, instance)
			}
		}
		return filteredInstances
	}
	blockListResponse.Items = blocksWithLabel(blockList)
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, params.GetBaseParams(), mock.PathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	blockListResponse.Items = blocksWithLabel(blockList)[:1]
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create images
	imageList, err := stubs.BulkCreateImagesStubV1(configurator, params.GetBaseParams(), params.Images)
	if err != nil {
		return nil, err
	}
	imageListResponse, err := builders.NewImageIteratorBuilder().
		Provider(mock.StorageProviderV1).
		Tenant(params.Tenant).
		Items(imageList).
		Build()
	if err != nil {
		return nil, err
	}

	// List images
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	imageListResponse.Items = imageList[:1]
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	imagesWithLabel := func(imageList []schema.Image) []schema.Image {
		var filteredImages []schema.Image
		for _, image := range imageList {
			if val, ok := image.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredImages = append(filteredImages, image)
			}
		}
		return filteredImages
	}

	imageListResponse.Items = imagesWithLabel(imageList)
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, params.GetBaseParams(), mock.PathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	imageListResponse.Items = imagesWithLabel(imageList)[:1]
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create storage skus
	skuList := steps.GenerateStorageSkusV1(params.Tenant)
	skuResponse, err := builders.NewStorageSkuIteratorBuilder().
		Provider(mock.StorageProviderV1).
		Tenant(params.Tenant).
		Items(skuList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.ConfigureGetListStorageSkuStub(*skuResponse, skuListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	skuResponse.Items = skuList[:1]
	if err := configurator.ConfigureGetListStorageSkuStub(*skuResponse, skuListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// Delete Lifecycle

	// Delete Images
	for _, image := range params.Images {
		url := generators.GenerateImageURL(mock.StorageProviderV1, params.Tenant, image.Name)

		// Delete the Image
		if err := configurator.ConfigureDeleteStub(url, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted Image
		if err := configurator.ConfigureGetNotFoundStub(url, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete BlockStorages
	for _, block := range params.BlockStorages {
		url := generators.GenerateBlockStorageURL(mock.StorageProviderV1, params.Tenant, params.Workspace.Name, block.Name)

		// Delete the BlockStorages
		if err := configurator.ConfigureDeleteStub(url, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted BlockStorages
		if err := configurator.ConfigureGetNotFoundStub(url, params.GetBaseParams()); err != nil {
			return nil, err
		}
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
