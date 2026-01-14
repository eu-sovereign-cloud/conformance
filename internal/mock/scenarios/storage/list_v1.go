package mockstorage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigureListScenarioV1(scenario string, mockParams *mock.MockParams, suiteParams *params.StorageListParamsV1) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)
	workspace := suiteParams.Workspace
	blockStorages := suiteParams.BlockStorages
	images := suiteParams.Images
	configurator, err := stubs.NewStubConfigurator(scenario, mockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	skuListUrl := generators.GenerateStorageSkuListURL(constants.StorageProviderV1, workspace.Metadata.Tenant)
	blockListUrl := generators.GenerateBlockStorageListURL(constants.StorageProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	imageListUrl := generators.GenerateImageListURL(constants.StorageProviderV1, workspace.Metadata.Tenant)

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

	// Create block storages
	err = stubs.BulkCreateBlockStoragesStubV1(configurator, mockParams, blockStorages)
	if err != nil {
		return nil, err
	}
	blockListResponse, err := builders.NewBlockStorageIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(workspace.Metadata.Tenant).Workspace(workspace.Metadata.Name).
		Items(blockStorages).
		Build()
	if err != nil {
		return nil, err
	}

	// List block storages
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, mockParams, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	blockListResponse.Items = blockStorages[:1]
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, mockParams, mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	blocksWithLabel := func(blockList []schema.BlockStorage) []schema.BlockStorage {
		var filteredInstances []schema.BlockStorage
		for _, instance := range blockList {
			if val, ok := instance.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredInstances = append(filteredInstances, instance)
			}
		}
		return filteredInstances
	}
	blockListResponse.Items = blocksWithLabel(blockStorages)
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, mockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	blockListResponse.Items = blocksWithLabel(blockStorages)[:1]
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, mockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create images
	err = stubs.BulkCreateImagesStubV1(configurator, mockParams, images)
	if err != nil {
		return nil, err
	}
	imageListResponse, err := builders.NewImageIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(workspace.Metadata.Tenant).
		Items(images).
		Build()
	if err != nil {
		return nil, err
	}

	// List images
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, mockParams, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	imageListResponse.Items = images[:1]
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, mockParams, mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	imagesWithLabel := func(imageList []schema.Image) []schema.Image {
		var filteredImages []schema.Image
		for _, image := range imageList {
			if val, ok := image.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredImages = append(filteredImages, image)
			}
		}
		return filteredImages
	}

	imageListResponse.Items = imagesWithLabel(images)
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, mockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	imageListResponse.Items = imagesWithLabel(images)[:1]
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, mockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create storage skus
	skuList := steps.GenerateStorageSkusV1(workspace.Metadata.Tenant)
	skuResponse, err := builders.NewStorageSkuIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(workspace.Metadata.Tenant).
		Items(skuList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.ConfigureGetListStorageSkuStub(*skuResponse, skuListUrl, mockParams, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	skuResponse.Items = skuList[:1]
	if err := configurator.ConfigureGetListStorageSkuStub(*skuResponse, skuListUrl, mockParams, mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// Delete Lifecycle

	// Delete Images
	for _, image := range images {
		url := generators.GenerateImageURL(constants.StorageProviderV1, image.Metadata.Tenant, image.Metadata.Name)

		// Delete the Image
		if err := configurator.ConfigureDeleteStub(url, mockParams); err != nil {
			return nil, err
		}

		// Get the deleted Image
		if err := configurator.ConfigureGetNotFoundStub(url, mockParams); err != nil {
			return nil, err
		}
	}

	// Delete BlockStorages
	for _, block := range blockStorages {
		url := generators.GenerateBlockStorageURL(constants.StorageProviderV1, block.Metadata.Tenant, block.Metadata.Workspace, block.Metadata.Name)

		// Delete the BlockStorages
		if err := configurator.ConfigureDeleteStub(url, mockParams); err != nil {
			return nil, err
		}

		// Get the deleted BlockStorages
		if err := configurator.ConfigureGetNotFoundStub(url, mockParams); err != nil {
			return nil, err
		}
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, mockParams); err != nil {
		return nil, err
	}
	return configurator.Client, nil
}
