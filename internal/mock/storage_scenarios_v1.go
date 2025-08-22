package mock

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wiremock/go-wiremock"
)

func CreateStorageLifecycleScenarioV1(scenario string, params StorageParamsV1) (*wiremock.Client, error) {
	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	storageSkuUrl := fmt.Sprintf(storageSkuURLV1, params.Tenant, params.Sku.Name)
	blockStorageUrl := fmt.Sprintf(blockStorageURLV1, params.Tenant, params.Workspace, params.BlockStorage.Name)
	imageUrl := fmt.Sprintf(imageURLV1, params.Tenant, params.Image.Name)

	// Storage sku
	skuResponse := storageSkuResponseV1{
		Metadata: metadataResponse{
			Name:            params.Sku.Name,
			Provider:        storageProviderV1,
			Resource:        fmt.Sprintf(skuResource, params.Tenant, params.Sku.Name),
			Verb:            http.MethodGet,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            storageSkuKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            activeStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
		Provider:      params.Sku.Provider,
		Tier:          params.Sku.Tier,
		Iops:          params.Sku.Iops,
		StorageType:   params.Sku.StorageType,
		MinVolumeSize: params.Sku.MinVolumeSize,
	}

	// Get storage sku
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          storageSkuUrl,
		params:       params,
		response:     skuResponse,
		template:     storageSkuResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    "CreateBlockStorage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse := blockStorageResponseV1{
		Metadata: metadataResponse{
			Name:       params.BlockStorage.Name,
			Provider:   storageProviderV1,
			Resource:   fmt.Sprintf(blockStorageResource, params.Tenant, params.Workspace, params.BlockStorage.Name),
			ApiVersion: version1,
			Kind:       blockStorageKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		SkuRef: params.BlockStorage.SkuRef,
	}

	// Create a block storage
	blockResponse.Metadata.Verb = http.MethodPut
	blockResponse.Status.State = creatingStatusState
	blockResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.ResourceVersion = 1
	blockResponse.SizeGB = params.BlockStorage.SizeGBInitial
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
		template:     blockStorageResponseTemplateV1,
		currentState: "CreateBlockStorage",
		nextState:    "GetCreatedBlockStorage",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created block storage
	blockResponse.Metadata.Verb = http.MethodGet
	blockResponse.Status.State = activeStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
		template:     blockStorageResponseTemplateV1,
		currentState: "GetCreatedBlockStorage",
		nextState:    "UpdateBlockStorage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the block storage
	blockResponse.Metadata.Verb = http.MethodPut
	blockResponse.Status.State = updatingStatusState
	blockResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.ResourceVersion = blockResponse.Metadata.ResourceVersion + 1
	blockResponse.SizeGB = params.BlockStorage.SizeGBUpdated
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
		template:     blockStorageResponseTemplateV1,
		currentState: "UpdateBlockStorage",
		nextState:    "GetUpdatedBlockStorage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get updated block storage
	blockResponse.Metadata.Verb = http.MethodGet
	blockResponse.Status.State = activeStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
		template:     blockStorageResponseTemplateV1,
		currentState: "GetUpdatedBlockStorage",
		nextState:    "CreateImage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// image
	imageResponse := imageResponseV1{
		Metadata: metadataResponse{
			Name:       params.Image.Name,
			Provider:   storageProviderV1,
			Resource:   fmt.Sprintf(imageResource, params.Tenant, params.Image.Name),
			ApiVersion: version1,
			Kind:       imageKind,
			Tenant:     params.Tenant,
			Region:     params.Region,
		},
		BlockStorageRef: params.Image.BlockStorageRef,
		CpuArchitecture: params.Image.CpuArchitectureInitial,
	}

	// Create an image
	imageResponse.Metadata.Verb = http.MethodPut
	imageResponse.Status.State = creatingStatusState
	imageResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	imageResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	imageResponse.Metadata.ResourceVersion = 1
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          imageUrl,
		params:       params,
		response:     imageResponse,
		template:     imageResponseTemplateV1,
		currentState: "CreateImage",
		nextState:    "GetCreatedImage",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created image
	imageResponse.Metadata.Verb = http.MethodGet
	imageResponse.Status.State = activeStatusState
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          imageUrl,
		params:       params,
		response:     imageResponse,
		template:     imageResponseTemplateV1,
		currentState: "GetCreatedImage",
		nextState:    "UpdateImage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the image
	imageResponse.Metadata.Verb = http.MethodPut
	imageResponse.Status.State = updatingStatusState
	imageResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	imageResponse.Metadata.ResourceVersion = imageResponse.Metadata.ResourceVersion + 1
	imageResponse.CpuArchitecture = params.Image.CpuArchitectureUpdated
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          imageUrl,
		params:       params,
		response:     imageResponse,
		template:     imageResponseTemplateV1,
		currentState: "UpdateImage",
		nextState:    "GetUpdatedImage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get updated image
	imageResponse.Metadata.Verb = http.MethodGet
	imageResponse.Status.State = activeStatusState
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          imageUrl,
		params:       params,
		response:     imageResponse,
		template:     imageResponseTemplateV1,
		currentState: "GetUpdatedImage",
		nextState:    "DeleteImage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Delete the image
	imageResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          imageUrl,
		params:       params,
		response:     imageResponse,
		currentState: "DeleteImage",
		nextState:    "GetDeletedImage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted image (not found)
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          imageUrl,
		params:       params,
		response:     imageResponse,
		currentState: "GetDeletedImage",
		nextState:    "DeleteBlockStorage",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the block storage
	blockResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
		currentState: "DeleteBlockStorage",
		nextState:    "GetDeletedBlockStorage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted block storage (not found)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
		currentState: "GetDeletedBlockStorage",
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
