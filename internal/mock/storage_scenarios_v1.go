package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigureStorageLifecycleScenarioV1(scenario string, params *StorageLifeCycleParamsV1) (*wiremock.Client, error) {
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

	// Update the block storage
	blockResponse.Spec = *params.BlockStorage.UpdatedSpec
	if err := configurator.configureUpdateBlockStorageStub(blockResponse, blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated block storage
	if err := configurator.configureGetActiveBlockStorageStub(blockResponse, blockUrl, params.getBaseParams()); err != nil {
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
	if err := configurator.configureCreateImageStub(imageResponse, imageUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created image
	if err := configurator.configureGetActiveImageStub(imageResponse, imageUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Update the image
	imageResponse.Spec = *params.Image.UpdatedSpec
	if err := configurator.configureUpdateImageStub(imageResponse, imageUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated image
	if err := configurator.configureGetActiveImageStub(imageResponse, imageUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the image
	if err := configurator.configureDeleteStub(imageUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted image
	if err := configurator.configureGetNotFoundStub(imageUrl, params.getBaseParams()); err != nil {
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

func ConfigureStorageListScenarioV1(scenario string, params *StorageListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	skuListUrl := generators.GenerateStorageSkuListURL(storageProviderV1, params.Tenant)
	blockListUrl := generators.GenerateBlockStorageListURL(storageProviderV1, params.Tenant, params.Workspace.Name)
	imageListUrl := generators.GenerateImageListURL(storageProviderV1, params.Tenant)

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

	// Create block storages
	blockList, err := bulkCreateBlockStoragesStubV1(configurator, params.getBaseParams(), params.Workspace.Name, params.BlockStorages)
	if err != nil {
		return nil, err
	}
	blockListResponse, err := builders.NewBlockStorageIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(blockList).
		Build()
	if err != nil {
		return nil, err
	}

	// List block storages
	if err := configurator.configureGetListBlockStorageStub(*blockListResponse, blockListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	blockListResponse.Items = blockList[:1]
	if err := configurator.configureGetListBlockStorageStub(*blockListResponse, blockListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListBlockStorageStub(*blockListResponse, blockListUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	blockListResponse.Items = blocksWithLabel(blockList)[:1]
	if err := configurator.configureGetListBlockStorageStub(*blockListResponse, blockListUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create images
	imageList, err := bulkCreateImagesStubV1(configurator, params.getBaseParams(), params.Images)
	if err != nil {
		return nil, err
	}
	imageListResponse, err := builders.NewImageIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).
		Items(imageList).
		Build()
	if err != nil {
		return nil, err
	}

	// List images
	if err := configurator.configureGetListImageStub(*imageListResponse, imageListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	imageListResponse.Items = imageList[:1]
	if err := configurator.configureGetListImageStub(*imageListResponse, imageListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListImageStub(*imageListResponse, imageListUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	imageListResponse.Items = imagesWithLabel(imageList)[:1]
	if err := configurator.configureGetListImageStub(*imageListResponse, imageListUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create storage skus
	skuList := generateStorageSkusV1(params.Tenant)
	skuResponse, err := builders.NewStorageSkuIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).
		Items(skuList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.configureGetListStorageSkuStub(*skuResponse, skuListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	skuResponse.Items = skuList[:1]
	if err := configurator.configureGetListStorageSkuStub(*skuResponse, skuListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// Delete Lifecycle

	// Delete Images
	for _, image := range params.Images {
		url := generators.GenerateImageURL(storageProviderV1, params.Tenant, image.Name)

		// Delete the Image
		if err := configurator.configureDeleteStub(url, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted Image
		if err := configurator.configureGetNotFoundStub(url, params.getBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete BlockStorages
	for _, block := range params.BlockStorages {
		url := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, block.Name)

		// Delete the BlockStorages
		if err := configurator.configureDeleteStub(url, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted BlockStorages
		if err := configurator.configureGetNotFoundStub(url, params.getBaseParams()); err != nil {
			return nil, err
		}
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
