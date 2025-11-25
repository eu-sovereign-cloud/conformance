package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	storageV1 "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/eu-sovereign-cloud/conformance/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigStorageLifecycleScenarioV1(scenario string, params *StorageParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)

	blockUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, (*params.BlockStorage)[0].Name)
	imageUrl := secalib.GenerateImageURL(params.Tenant, (*params.Image)[0].Name)

	blockResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, (*params.BlockStorage)[0].Name)
	imageResource := secalib.GenerateImageResource(params.Tenant, (*params.Image)[0].Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(secalib.WorkspaceProviderV1).
		Resource(workspaceResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	setCreatedRegionalResourceMetadata(workspaceResponse.Metadata)
	workspaceResponse.Status = secalib.NewWorkspaceStatus(schema.ResourceStateCreating)
	workspaceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	secalib.SetWorkspaceStatusState(workspaceResponse.Status, schema.ResourceStateActive)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: "GetCreatedWorkspace", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(secalib.StorageProviderV1).
		Resource(blockResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = secalib.NewBlockStorageStatus(schema.ResourceStateCreating)
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the created block storage
	secalib.SetBlockStorageStatusState(blockResponse.Status, schema.ResourceStateActive)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "UpdateBlockStorage"}); err != nil {
		return nil, err
	}

	// Update the block storage
	setModifiedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	secalib.SetBlockStorageStatusState(blockResponse.Status, schema.ResourceStateUpdating)
	blockResponse.Spec = *params.BlockStorage.UpdatedSpec
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "UpdateBlockStorage", nextState: "GetUpdatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the updated block storage
	secalib.SetBlockStorageStatusState(blockResponse.Status, schema.ResourceStateActive)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetUpdatedBlockStorage", nextState: "CreateImage"}); err != nil {
		return nil, err
	}

	// Image
	imageResponse, err := builders.NewImageBuilder().
		Name(params.Image.Name).
		Provider(secalib.StorageProviderV1).
		Resource(imageResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Region(params.Region).
		Spec(params.Image.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create an image
	setCreatedRegionalResourceMetadata(imageResponse.Metadata)
	imageResponse.Status = secalib.NewImageStatus(schema.ResourceStateCreating)
	imageResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "CreateImage", nextState: "GetCreatedImage"}); err != nil {
		return nil, err
	}

	// Get the created image
	secalib.SetImageStatusState(imageResponse.Status, schema.ResourceStateActive)
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "GetCreatedImage", nextState: "UpdateImage"}); err != nil {
		return nil, err
	}

	// Update the image
	setModifiedRegionalResourceMetadata(imageResponse.Metadata)
	secalib.SetImageStatusState(imageResponse.Status, secalib.UpdatingResourceState)
	imageResponse.Spec = *(*params.Image)[0].UpdatedSpec
	imageResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "UpdateImage", nextState: "GetUpdatedImage"}); err != nil {
		return nil, err
	}

	// Get the updated image
	secalib.SetImageStatusState(imageResponse.Status, schema.ResourceStateActive)
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "GetUpdatedImage", nextState: "DeleteImage"}); err != nil {
		return nil, err
	}

	// Delete the image
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, currentState: "DeleteImage", nextState: "GetDeletedImage"}); err != nil {
		return nil, err
	}

	// Get the deleted image
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, currentState: "GetDeletedImage", nextState: "DeleteBlockStorage"}); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, currentState: "DeleteBlockStorage", nextState: "GetDeletedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, currentState: "GetDeletedBlockStorage", nextState: "DeleteWorkspace"}); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, currentState: "DeleteWorkspace", nextState: "GetDeletedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, currentState: "GetDeletedWorkspace", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}

func ConfigStorageListLifecycleScenarioV1(scenario string, params *StorageParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)

	blockResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, (*params.BlockStorage)[0].Name)

	// Workspace
	workspaceResponse := newWorkspaceResponse(params.Workspace.Name, secalib.WorkspaceProviderV1, workspaceResource, secalib.ApiVersion1,
		params.Tenant, params.Region,
		params.Workspace.InitialLabels)

	// Create a workspace
	setCreatedRegionalResourceMetadata(workspaceResponse.Metadata)
	workspaceResponse.Status = secalib.NewWorkspaceStatus(secalib.CreatingResourceState)
	workspaceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: startedScenarioState, nextState: (*params.BlockStorage)[0].Name}); err != nil {
		return nil, err
	}

	// Block storage
	var blockList []schema.BlockStorage
	for i := range *params.BlockStorage {
		blockResource = secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, (*params.BlockStorage)[i].Name)
		blockResponse := newBlockStorageResponse((*params.BlockStorage)[i].Name, secalib.StorageProviderV1, blockResource, secalib.ApiVersion1,
			params.Tenant, params.Workspace.Name, params.Region,
			(*params.BlockStorage)[i].InitialLabels,
			(*params.BlockStorage)[i].InitialSpec)

		var nextState string
		if i < len(*params.BlockStorage)-1 {
			nextState = (*params.BlockStorage)[i+1].Name
		} else {
			nextState = "GetBlockStorageList"
		}
		// Create a block storage
		setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
		blockResponse.Status = secalib.NewBlockStorageStatus(secalib.CreatingResourceState)
		blockResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{
				url: secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, (*params.BlockStorage)[i].Name), params: params, responseBody: blockResponse,
				currentState: (*params.BlockStorage)[i].Name, nextState: nextState,
			}); err != nil {
			return nil, err
		}
		blockList = append(blockList, *blockResponse)
	}

	// List
	blockListResource := secalib.GenerateBlockStorageListResource(params.Tenant, params.Workspace.Name)
	blockListResponse := &storageV1.BlockStorageIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.StorageProviderV1,
			Resource: blockListResource,
			Verb:     http.MethodGet,
		},
	}
	blockListResponse.Items = blockList

	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateBlockStorageListURL(params.Tenant, params.Workspace.Name), params: params, responseBody: blockListResponse,
			currentState: "GetBlockStorageList", nextState: "GetBlockStorageListWithLimit",
		}); err != nil {
		return nil, err
	}
	// List with limit

	blockListResponse.Items = blockList[:1]

	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateBlockStorageListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimit("1"), responseBody: blockListResponse,
			currentState: "GetBlockStorageListWithLimit", nextState: "GetBlockStorageListWithLabel",
		}); err != nil {
		return nil, err
	}

	// List with Label

	blocksWithLabel := func(blockList []schema.BlockStorage) []schema.BlockStorage {
		var filteredInstances []schema.BlockStorage
		for _, instance := range blockList {
			if val, ok := instance.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformance {
				filteredInstances = append(filteredInstances, instance)
			}
		}
		return filteredInstances
	}
	blockListResponse.Items = blocksWithLabel(blockList)

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateBlockStorageListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformance), responseBody: blockListResponse, currentState: "GetBlockStorageListWithLabel", nextState: "GetBlockStorageListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List with limit & label

	blockListResponse.Items = blocksWithLabel(blockList)[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateBlockStorageListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformance), responseBody: blockListResponse, currentState: "GetBlockStorageListWithLimitAndLabel", nextState: (*params.Image)[0].Name}); err != nil {
		return nil, err
	}
	// image
	var imageList []schema.Image
	for i := range *params.Image {
		imageResource := secalib.GenerateImageResource(params.Tenant, (*params.Image)[i].Name)
		imageResp := newImageResponse((*params.Image)[i].Name, secalib.StorageProviderV1, imageResource, secalib.ApiVersion1,
			params.Tenant, params.Region,
			&(*params.Image)[i].InitialLabels,
			(*params.Image)[i].InitialSpec)

		var nextState string
		if i < len(*params.Image)-1 {
			nextState = (*params.Image)[i+1].Name
		} else {
			nextState = "GetImageList"
		}

		// Create image
		imageResp.Status = secalib.NewImageStatus(secalib.CreatingResourceState)
		imageResp.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateImageURL(params.Tenant, (*params.Image)[i].Name), params: params, responseBody: imageResp, currentState: (*params.Image)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		imageList = append(imageList, *imageResp)
	}
	// List
	imageListResource := secalib.GenerateImageListResource(params.Tenant)
	imageListResponse := &storageV1.ImageIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.StorageProviderV1,
			Resource: imageListResource,
			Verb:     http.MethodGet,
		},
	}
	imageListResponse.Items = imageList
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateImageListURL(params.Tenant), params: params, responseBody: imageListResponse, currentState: "GetImageList", nextState: "GetImageListWithLimit"}); err != nil {
		return nil, err
	}
	// List with limit

	imageListWithLimitResponse := &storageV1.ImageIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.StorageProviderV1,
			Resource: imageListResource,
			Verb:     http.MethodGet,
		},
	}
	imageListWithLimitResponse.Items = imageList[:1]

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateImageListURL(params.Tenant), params: params, pathParams: pathParamsLimit("1"), responseBody: imageListWithLimitResponse, currentState: "GetImageListWithLimit", nextState: "GetImageListWithLabel"}); err != nil {
		return nil, err
	}
	// List with Label
	imageListWithLabelResponse := &storageV1.ImageIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.StorageProviderV1,
			Resource: imageListResource,
			Verb:     http.MethodGet,
		},
	}

	imagesWithLabel := func(imageList []schema.Image) []schema.Image {
		var filteredImages []schema.Image
		for _, image := range imageList {
			if val, ok := image.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformance {
				filteredImages = append(filteredImages, image)
			}
		}
		return filteredImages
	}
	imageListWithLabelResponse.Items = imagesWithLabel(imageList)

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateImageListURL(params.Tenant), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformance), responseBody: imageListWithLabelResponse, currentState: "GetImageListWithLabel", nextState: "GetImageListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List with limit & label

	imageListWithLimitAndLabelResponse := &storageV1.ImageIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.StorageProviderV1,
			Resource: imageListResource,
			Verb:     http.MethodGet,
		},
	}

	imageListWithLimitAndLabelResponse.Items = imagesWithLabel(imageList)[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateImageListURL(params.Tenant), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformance), responseBody: imageListWithLimitAndLabelResponse, currentState: "GetImageListWithLimitAndLabel", nextState: "GetSkuList"}); err != nil {
		return nil, err
	}

	// Storage Sku

	skuList := []schema.StorageSku{
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "RD100",
				Provider: secalib.StorageProviderV1,
				Resource: blockResource,
				Verb:     http.MethodGet,
				Kind:     secalib.StorageSkuKind,
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
				Provider: secalib.StorageProviderV1,
				Resource: blockResource,
				Verb:     http.MethodGet,
				Kind:     secalib.StorageSkuKind,
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
				Provider: secalib.StorageProviderV1,
				Resource: blockResource,
				Verb:     http.MethodGet,
				Kind:     secalib.StorageSkuKind,
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
	skuResponse := &storageV1.SkuIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.StorageProviderV1,
			Resource: secalib.GenerateSkuListResource(params.Tenant),
			Verb:     http.MethodGet,
		},
		Items: skuList,
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateStorageSkuListURL(params.Tenant), params: params, responseBody: skuResponse, currentState: "GetSkuList", nextState: "GetSkuListWithLimit"}); err != nil {
		return nil, err
	}
	// List with limit

	skuWithLimitResponse := &storageV1.SkuIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.StorageProviderV1,
			Resource: secalib.GenerateSkuListResource(params.Tenant),
			Verb:     http.MethodGet,
		},
	}
	skuWithLimitResponse.Items = skuList[:1]

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateStorageSkuListURL(params.Tenant), params: params, pathParams: pathParamsLimit("1"), responseBody: skuWithLimitResponse, currentState: "GetSkuListWithLimit", nextState: "GetSkuListWithLabel"}); err != nil {
		return nil, err
	}
	// List with Label
	skuWithLabelResponse := &storageV1.SkuIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.StorageProviderV1,
			Resource: secalib.GenerateStorageSkuListURL(params.Tenant),
			Verb:     http.MethodGet,
		},
	}

	skusWithLabel := func(skuList []schema.StorageSku) []schema.StorageSku {
		var filteredSkus []schema.StorageSku
		for _, sku := range skuList {
			if val, ok := sku.Labels["tier"]; ok && val == "RD500" {
				filteredSkus = append(filteredSkus, sku)
			}
		}
		return filteredSkus
	}
	skuWithLabelResponse.Items = skusWithLabel(skuList)

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateStorageSkuListURL(params.Tenant), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformance), responseBody: skuWithLabelResponse, currentState: "GetSkuListWithLabel", nextState: "GetSkuListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List with limit & label

	skuWithLimitAndLabelResponse := &storageV1.SkuIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.StorageProviderV1,
			Resource: secalib.GenerateSkuListResource(params.Tenant),
			Verb:     http.MethodGet,
		},
	}

	skuWithLimitAndLabelResponse.Items = skusWithLabel(skuList)[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateStorageSkuListURL(params.Tenant), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformance), responseBody: skuWithLimitAndLabelResponse, currentState: "GetSkuListWithLimitAndLabel", nextState: "DeleteImage_" + (*params.Image)[0].Name}); err != nil {
		return nil, err
	}

	// Delete Images
	for i := range *params.Image {
		imageUrl := secalib.GenerateImageURL(params.Tenant, (*params.Image)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteImage_" + (*params.Image)[i].Name
		} else {
			currentState = "GetDeletedImage_" + (*params.Image)[i-1].Name
		}

		nextState = "DeleteImage_" + (*params.Image)[i].Name

		// Delete the image
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: imageUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted image (should return 404)
		nextState = func() string {
			if i < len(*params.Image)-1 {
				return "GetDeletedImage_" + (*params.Image)[i].Name
			} else {
				return "DeleteBlockStorage_" + (*params.BlockStorage)[0].Name
			}
		}()

		if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
			&stubConfig{url: imageUrl, params: params, currentState: "DeleteImage_" + (*params.Image)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
	}

	// Delete BlockStorages
	for i := range *params.BlockStorage {
		blockUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, (*params.BlockStorage)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteBlockStorage_" + (*params.BlockStorage)[i].Name
		} else {
			currentState = "GetDeletedBlockStorage_" + (*params.BlockStorage)[i-1].Name
		}

		nextState = "DeleteBlockStorage_" + (*params.BlockStorage)[i].Name

		// Delete the block storage
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: blockUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted block storage (should return 404)
		nextState = func() string {
			if i < len(*params.BlockStorage)-1 {
				return "GetDeletedBlockStorage_" + (*params.BlockStorage)[i].Name
			} else {
				return "DeleteWorkspace"
			}
		}()

		if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
			&stubConfig{url: blockUrl, params: params, currentState: "DeleteBlockStorage_" + (*params.BlockStorage)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
	}

	// Delete the workspace
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, currentState: "DeleteWorkspace", nextState: "GetDeletedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: workspaceUrl, params: params, currentState: "GetDeletedWorkspace", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}
