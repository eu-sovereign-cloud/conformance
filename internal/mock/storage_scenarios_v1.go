package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

func CreateStorageLifecycleScenarioV1(scenario string, params *StorageParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)

	blockUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageUrl := secalib.GenerateImageURL(params.Tenant, params.Image.Name)

	blockResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageResource := secalib.GenerateImageResource(params.Tenant, params.Image.Name)

	// Workspace
	workspaceResponse := newWorkspaceResponse(params.Workspace.Name, secalib.WorkspaceProviderV1, workspaceResource, secalib.ApiVersion1,
		params.Tenant, params.Region,
		params.Workspace.InitialLabels)

	// Create a workspace
	setCreatedRegionalResourceMetadata(workspaceResponse.Metadata)
	workspaceResponse.Status = newWorkspaceStatus(secalib.CreatingStatusState)
	workspaceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, response: workspaceResponse, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse := newBlockStorageResponse(params.BlockStorage.Name, secalib.ComputeProviderV1, blockResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.BlockStorage.InitialSpec)

	// Get created workspace
	setWorkspaceStatusState(workspaceResponse.Status, secalib.ActiveStatusState)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, response: workspaceResponse, currentState: "GetCreatedWorkspace", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = newBlockStorageStatus(secalib.CreatingStatusState)
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, response: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get created block storage
	setBlockStorageStatusState(blockResponse.Status, secalib.ActiveStatusState)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, response: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "CreateInstance"}); err != nil {
		return nil, err
	}

	// Update the block storage
	setModifiedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	setBlockStorageStatusState(blockResponse.Status, secalib.UpdatingStatusState)
	blockResponse.Spec = *params.BlockStorage.UpdatedSpec
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, response: blockResponse, currentState: "UpdateBlockStorage", nextState: "GetUpdatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get updated block storage
	setBlockStorageStatusState(blockResponse.Status, secalib.ActiveStatusState)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, response: blockResponse, currentState: "GetUpdatedBlockStorage", nextState: "CreateImage"}); err != nil {
		return nil, err
	}

	// image
	imageResponse := newImageResponse(params.Image.Name, secalib.ComputeProviderV1, imageResource, secalib.ApiVersion1,
		params.Tenant, params.Region,
		params.Image.InitialSpec)

	// Create an image
	setCreatedRegionalResourceMetadata(imageResponse.Metadata)
	imageResponse.Status = newImageStatus(secalib.CreatingStatusState)
	imageResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, response: imageResponse, currentState: "CreateImage", nextState: "GetCreatedImage"}); err != nil {
		return nil, err
	}

	// Get created image
	setImageStatusState(imageResponse.Status, secalib.ActiveStatusState)
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, response: imageResponse, currentState: "GetCreatedImage", nextState: "UpdateImage"}); err != nil {
		return nil, err
	}

	// Update the image
	setModifiedRegionalResourceMetadata(imageResponse.Metadata)
	setImageStatusState(imageResponse.Status, secalib.UpdatingStatusState)
	imageResponse.Spec = *params.Image.UpdatedSpec
	imageResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, response: imageResponse, currentState: "UpdateImage", nextState: "GetUpdatedImage"}); err != nil {
		return nil, err
	}

	// Get updated image
	setImageStatusState(imageResponse.Status, secalib.ActiveStatusState)
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, response: imageResponse, currentState: "GetUpdatedImage", nextState: "DeleteImage"}); err != nil {
		return nil, err
	}

	// Delete the image
	imageResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, response: imageResponse, currentState: "DeleteImage", nextState: "GetDeletedImage"}); err != nil {
		return nil, err
	}

	// Get deleted image
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: imageUrl, params: params, response: imageResponse, currentState: "GetDeletedImage", nextState: "DeleteBlockStorage"}); err != nil {
		return nil, err
	}

	// Delete the block storage
	blockResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, response: blockResponse, currentState: "DeleteBlockStorage", nextState: "GetDeletedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get deleted block storage
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: blockUrl, params: params, response: blockResponse, currentState: "GetDeletedBlockStorage", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}
