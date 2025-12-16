package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigStorageLifecycleScenarioV1(scenario string, params *StorageParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageUrl := generators.GenerateImageURL(storageProviderV1, params.Tenant, params.Image.Name)

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

	// Update the block storage
	blockResponse.Spec = *params.BlockStorage.UpdatedSpec
	if err := configurator.configureUpdateBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Get the updated block storage
	if err := configurator.configureGetActiveBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Image
	imageResponse, err := builders.NewImageBuilder().
		Name(params.Image.Name).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Spec(params.Image.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an image
	if err := configurator.configureCreateImageStub(imageResponse, imageUrl, params); err != nil {
		return nil, err
	}

	// Get the created image
	if err := configurator.configureGetActiveImageStub(imageResponse, imageUrl, params); err != nil {
		return nil, err
	}

	// Update the image
	imageResponse.Spec = *params.Image.UpdatedSpec
	if err := configurator.configureUpdateImageStub(imageResponse, imageUrl, params); err != nil {
		return nil, err
	}

	// Get the updated image
	if err := configurator.configureGetActiveImageStub(imageResponse, imageUrl, params); err != nil {
		return nil, err
	}

	// Delete the image
	if err := configurator.configureDeleteStub(imageUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted image
	if err := configurator.configureGetNotFoundStub(imageUrl, params); err != nil {
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

func ConfigStorageListAndFilterScenarioV1(scenario string, params *StorageListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)

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

	var blockList []schema.BlockStorage
	for _, block := range *params.BlockStorage {
		// Block storage
		blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, block.Name)
		blockResponse, err := builders.NewBlockStorageBuilder().
			Name(block.Name).
			Provider(storageProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
			Labels(block.InitialLabels).
			Spec(block.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create a block storage
		if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, params); err != nil {
			return nil, err
		}

		blockList = append(blockList, *blockResponse)
	}

	// List block storages
	blockListResource := generators.GenerateBlockStorageListResource(params.Tenant, params.Workspace.Name)
	blockListUrl := generators.GenerateBlockStorageListURL(storageProviderV1, params.Tenant, params.Workspace.Name)
	blockListResponse := &storage.BlockStorageIterator{
		Metadata: schema.ResponseMetadata{
			Provider: storageProviderV1,
			Resource: blockListResource,
			Verb:     http.MethodGet,
		},
		Items: blockList,
	}

	if err := configurator.configureGetListBlockStorageStub(*blockListResponse, blockListUrl, params, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	blockListResponse.Items = blockList[:1]
	if err := configurator.configureGetListBlockStorageStub(*blockListResponse, blockListUrl, params, pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListBlockStorageStub(*blockListResponse, blockListUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	blockListResponse.Items = blocksWithLabel(blockList)[:1]
	if err := configurator.configureGetListBlockStorageStub(*blockListResponse, blockListUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Image

	var imageList []schema.Image
	for _, image := range *params.Image {
		imageUrl := generators.GenerateImageURL(storageProviderV1, params.Tenant, image.Name)
		imageResponse, err := builders.NewImageBuilder().
			Name(image.Name).
			Provider(storageProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).Region(params.Region).
			Labels(image.InitialLabels).
			Spec(image.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create an image
		if err := configurator.configureCreateImageStub(imageResponse, imageUrl, params); err != nil {
			return nil, err
		}
		imageList = append(imageList, *imageResponse)
	}

	imageListResource := generators.GenerateImageListResource(params.Tenant)
	imageListUrl := generators.GenerateImageListURL(storageProviderV1, params.Tenant)
	imageListResponse := &storage.ImageIterator{
		Metadata: schema.ResponseMetadata{
			Provider: storageProviderV1,
			Resource: imageListResource,
			Verb:     http.MethodGet,
		},
		Items: imageList,
	}
	if err := configurator.configureGetListImageStub(*imageListResponse, imageListUrl, params, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	imageListResponse.Items = imageList[:1]
	if err := configurator.configureGetListImageStub(*imageListResponse, imageListUrl, params, pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListImageStub(*imageListResponse, imageListUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	imageListResponse.Items = imagesWithLabel(imageList)[:1]
	if err := configurator.configureGetListImageStub(*imageListResponse, imageListUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Storage Sku
	block := blockList[:1]
	blockResource := generators.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, block[0].Metadata.Name)
	skuList := []schema.StorageSku{
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "RD100",
				Provider: storageProviderV1,
				Resource: blockResource,
				Verb:     http.MethodGet,
				Tenant:   params.Tenant,
			},
			Labels: schema.Labels{
				"provider": "seca",
				"tier":     "RD100",
			},
			Spec: &schema.StorageSkuSpec{
				Iops:          100,
				MinVolumeSize: 50,
				Type:          "remote-durable",
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "RD500",
				Provider: storageProviderV1,
				Resource: blockResource,
				Verb:     http.MethodGet,
				Tenant:   params.Tenant,
			},
			Labels: schema.Labels{
				"provider": "seca",
				"tier":     "RD500",
			},
			Spec: &schema.StorageSkuSpec{
				Iops:          500,
				MinVolumeSize: 50,
				Type:          "remote-durable",
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "RD2K",
				Provider: storageProviderV1,
				Resource: blockResource,
				Verb:     http.MethodGet,
				Tenant:   params.Tenant,
			},
			Labels: schema.Labels{
				"provider": "seca",
				"tier":     "RD2k",
			},
			Spec: &schema.StorageSkuSpec{
				Iops:          2000,
				MinVolumeSize: 50,
				Type:          "remote-durable",
			},
		},
	}

	// List
	skuListUrl := generators.GenerateStorageSkuListURL(storageProviderV1, params.Tenant)
	skuResponse := &storage.SkuIterator{
		Metadata: schema.ResponseMetadata{
			Provider: storageProviderV1,
			Resource: generators.GenerateSkuListResource(params.Tenant),
			Verb:     http.MethodGet,
		},
		Items: skuList,
	}

	if err := configurator.configureGetListStorageSkuStub(*skuResponse, skuListUrl, params, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	skuResponse.Items = skuList[:1]
	if err := configurator.configureGetListStorageSkuStub(*skuResponse, skuListUrl, params, pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	skusWithLabel := func(skuList []schema.StorageSku) []schema.StorageSku {
		var filteredSkus []schema.StorageSku
		for _, sku := range skuList {
			if val, ok := sku.Labels["tier"]; ok && val == "RD500" {
				filteredSkus = append(filteredSkus, sku)
			}
		}
		return filteredSkus
	}

	skuResponse.Items = skusWithLabel(skuList)
	if err := configurator.configureGetListStorageSkuStub(*skuResponse, skuListUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	skuResponse.Items = skusWithLabel(skuList)[:1]
	if err := configurator.configureGetListStorageSkuStub(*skuResponse, skuListUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Delete Lifecycle

	// Delete Images
	for _, image := range *params.Image {
		url := generators.GenerateImageURL(storageProviderV1, params.Tenant, image.Name)

		// Delete the Image
		if err := configurator.configureDeleteStub(url, params); err != nil {
			return nil, err
		}

		// Get the deleted Image
		if err := configurator.configureGetNotFoundStub(url, params); err != nil {
			return nil, err
		}
	}

	// Delete BlockStorages
	for _, block := range *params.BlockStorage {
		url := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, block.Name)

		// Delete the BlockStorages
		if err := configurator.configureDeleteStub(url, params); err != nil {
			return nil, err
		}

		// Get the deleted BlockStorages
		if err := configurator.configureGetNotFoundStub(url, params); err != nil {
			return nil, err
		}
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
