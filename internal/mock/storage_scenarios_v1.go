package mock

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

func CreateStorageLifecycleScenarioV1(scenario string, params StorageParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to Storage Lifecycle Scenario")

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	blockStorageUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace, params.BlockStorage.Name)
	imageUrl := secalib.GenerateImageURL(params.Tenant, params.Image.Name)

	blockStorageResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace, params.BlockStorage.Name)
	imageResource := secalib.GenerateImageResource(params.Tenant, params.Image.Name)

	// Block storage
	blockResponse := &resourceResponse[secalib.BlockStorageSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.BlockStorage.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.BlockStorageSpecV1{
			SkuRef: params.BlockStorage.InitialSpec.SkuRef,
		},
	}

	// Create a block storage
	blockResponse.Metadata.Verb = http.MethodPut
	blockResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.ResourceVersion = 1
	blockResponse.Spec.SizeGB = params.BlockStorage.InitialSpec.SizeGB
	blockResponse.Status.State = secalib.CreatingStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
		template:     blockStorageResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    "GetCreatedBlockStorage",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created block storage
	blockResponse.Metadata.Verb = http.MethodGet
	blockResponse.Status.State = secalib.ActiveStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
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
	blockResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.ResourceVersion = blockResponse.Metadata.ResourceVersion + 1
	blockResponse.Spec.SizeGB = params.BlockStorage.UpdatedSpec.SizeGB
	blockResponse.Status.State = secalib.UpdatingStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
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
	blockResponse.Status.State = secalib.ActiveStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
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
	imageResponse := &resourceResponse[secalib.ImageSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Image.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   imageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.ImageKind,
			Tenant:     params.Tenant,
			Region:     params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.ImageSpecV1{
			BlockStorageRef: params.Image.InitialSpec.BlockStorageRef,
			CpuArchitecture: params.Image.InitialSpec.CpuArchitecture,
		},
	}

	// Create an image
	imageResponse.Metadata.Verb = http.MethodPut
	imageResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	imageResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	imageResponse.Metadata.ResourceVersion = 1
	imageResponse.Status.State = secalib.CreatingStatusState
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
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
	imageResponse.Status.State = secalib.ActiveStatusState
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
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
	imageResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	imageResponse.Metadata.ResourceVersion = imageResponse.Metadata.ResourceVersion + 1
	imageResponse.Spec.CpuArchitecture = params.Image.UpdatedSpec.CpuArchitecture
	imageResponse.Status.State = secalib.UpdatingStatusState
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
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
	imageResponse.Status.State = secalib.ActiveStatusState
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
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
	if err := configureDeleteStub(wm, scenario, stubConfig{
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
	if err := configureGetStub(wm, scenario, stubConfig{
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
	if err := configureDeleteStub(wm, scenario, stubConfig{
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
		currentState: "GetDeletedBlockStorage",
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	slog.Info("Configured mock to Storage Lifecycle Scenario")
	return wm, nil
}
