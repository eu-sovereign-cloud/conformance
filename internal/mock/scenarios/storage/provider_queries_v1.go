package mockstorage

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

func ConfigureProviderQueriesV1(scenario *mockscenarios.Scenario, params *params.StorageProviderQueriesV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := params.Workspace
	blockStorages := params.BlockStorages
	images := params.Images

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
		return err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create block storages
	err = stubs.BulkCreateBlockStoragesStubV1(configurator, scenario.MockParams, blockStorages)
	if err != nil {
		return err
	}
	blockListResponse, err := builders.NewBlockStorageIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(workspace.Metadata.Tenant).Workspace(workspace.Metadata.Name).
		Items(blockStorages).
		Build()
	if err != nil {
		return err
	}

	// List block storages
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with Limit 1
	blockListResponse.Items = blockStorages[:1]
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
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
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with Limit and Label
	blockListResponse.Items = blocksWithLabel(blockStorages)[:1]
	if err := configurator.ConfigureGetListBlockStorageStub(*blockListResponse, blockListUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Create images
	err = stubs.BulkCreateImagesStubV1(configurator, scenario.MockParams, images)
	if err != nil {
		return err
	}
	imageListResponse, err := builders.NewImageIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(workspace.Metadata.Tenant).
		Items(images).
		Build()
	if err != nil {
		return err
	}

	// List images
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with Limit 1
	imageListResponse.Items = images[:1]
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
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
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with Limit and Label
	imageListResponse.Items = imagesWithLabel(images)[:1]
	if err := configurator.ConfigureGetListImageStub(*imageListResponse, imageListUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Create storage skus
	skuList := steps.GenerateStorageSkusV1(workspace.Metadata.Tenant)
	skuResponse, err := builders.NewStorageSkuIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(workspace.Metadata.Tenant).
		Items(skuList).
		Build()
	if err != nil {
		return err
	}

	// List
	if err := configurator.ConfigureGetListStorageSkuStub(*skuResponse, skuListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with Limit 1
	skuResponse.Items = skuList[:1]
	if err := configurator.ConfigureGetListStorageSkuStub(*skuResponse, skuListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
	}

	// Delete Lifecycle

	// Delete Images
	for _, image := range images {
		url := generators.GenerateImageURL(constants.StorageProviderV1, image.Metadata.Tenant, image.Metadata.Name)

		// Delete the Image
		if err := configurator.ConfigureDeleteStub(url, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted Image
		if err := configurator.ConfigureGetNotFoundStub(url, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete BlockStorages
	for _, block := range blockStorages {
		url := generators.GenerateBlockStorageURL(constants.StorageProviderV1, block.Metadata.Tenant, block.Metadata.Workspace, block.Metadata.Name)

		// Delete the BlockStorages
		if err := configurator.ConfigureDeleteStub(url, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted BlockStorages
		if err := configurator.ConfigureGetNotFoundStub(url, scenario.MockParams); err != nil {
			return err
		}
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
