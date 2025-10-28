package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"

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
	workspaceResponse.Status = secalib.NewWorkspaceStatus(secalib.CreatingResourceState)
	workspaceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	secalib.SetWorkspaceStatusState(workspaceResponse.Status, secalib.ActiveResourceState)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: "GetCreatedWorkspace", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse := newBlockStorageResponse(params.BlockStorage.Name, secalib.StorageProviderV1, blockResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.BlockStorage.InitialSpec)

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = secalib.NewBlockStorageStatus(secalib.CreatingResourceState)
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the created block storage
	secalib.SetBlockStorageStatusState(blockResponse.Status, secalib.ActiveResourceState)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "UpdateBlockStorage"}); err != nil {
		return nil, err
	}

	// Update the block storage
	setModifiedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	secalib.SetBlockStorageStatusState(blockResponse.Status, secalib.UpdatingResourceState)
	blockResponse.Spec = *params.BlockStorage.UpdatedSpec
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "UpdateBlockStorage", nextState: "GetUpdatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the updated block storage
	secalib.SetBlockStorageStatusState(blockResponse.Status, secalib.ActiveResourceState)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetUpdatedBlockStorage", nextState: "CreateImage"}); err != nil {
		return nil, err
	}

	// image
	imageResponse := newImageResponse(params.Image.Name, secalib.StorageProviderV1, imageResource, secalib.ApiVersion1,
		params.Tenant, params.Region,
		params.Image.InitialSpec)

	// Create an image
	setCreatedRegionalResourceMetadata(imageResponse.Metadata)
	imageResponse.Status = secalib.NewImageStatus(secalib.CreatingResourceState)
	imageResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "CreateImage", nextState: "GetCreatedImage"}); err != nil {
		return nil, err
	}

	// Get the created image
	secalib.SetImageStatusState(imageResponse.Status, secalib.ActiveResourceState)
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "GetCreatedImage", nextState: "UpdateImage"}); err != nil {
		return nil, err
	}

	// Update the image
	setModifiedRegionalResourceMetadata(imageResponse.Metadata)
	secalib.SetImageStatusState(imageResponse.Status, secalib.UpdatingResourceState)
	imageResponse.Spec = *params.Image.UpdatedSpec
	imageResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "UpdateImage", nextState: "GetUpdatedImage"}); err != nil {
		return nil, err
	}

	// Get the updated image
	secalib.SetImageStatusState(imageResponse.Status, secalib.ActiveResourceState)
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
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: imageUrl, params: params, currentState: "GetDeletedImage", nextState: "DeleteBlockStorage"}); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, currentState: "DeleteBlockStorage", nextState: "GetDeletedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: blockUrl, params: params, currentState: "GetDeletedBlockStorage", nextState: "DeleteWorkspace"}); err != nil {
		return nil, err
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

func ConfigStorageWithWaitingScenarioV1(scenario string, params *StorageParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)

	blockUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)

	blockResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)

	// Workspace
	workspaceResponse := newWorkspaceResponse(params.Workspace.Name, secalib.WorkspaceProviderV1, workspaceResource, secalib.ApiVersion1,
		params.Tenant, params.Region,
		params.Workspace.InitialLabels)

	// Create a workspace
	setCreatedRegionalResourceMetadata(workspaceResponse.Metadata)
	workspaceResponse.Status = secalib.NewWorkspaceStatus(secalib.CreatingResourceState)
	workspaceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: startedScenarioState, nextState: "GetCreatedStatusNotReadyTry1Workspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	secalib.SetWorkspaceStatusState(workspaceResponse.Status, secalib.CreatingResourceState)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: "GetCreatedStatusNotReadyTry1Workspace", nextState: "GetCreatedStatusNotReadyTry2Workspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	secalib.SetWorkspaceStatusState(workspaceResponse.Status, secalib.CreatingResourceState)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: "GetCreatedStatusNotReadyTry2Workspace", nextState: "GetCreatedStatusNotReadyTry3Workspace"}); err != nil {
		return nil, err
	}

	secalib.SetWorkspaceStatusState(workspaceResponse.Status, secalib.ActiveResourceState)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: "GetCreatedStatusNotReadyTry3Workspace", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}
	// Block storage
	blockResponse := newBlockStorageResponse(params.BlockStorage.Name, secalib.StorageProviderV1, blockResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.BlockStorage.InitialSpec)

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = secalib.NewBlockStorageStatus(secalib.CreatingResourceState)
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the created block storage
	secalib.SetBlockStorageStatusState(blockResponse.Status, secalib.ActiveResourceState)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "end"}); err != nil {
		return nil, err
	}

	return wm, nil
}
