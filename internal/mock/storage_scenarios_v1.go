package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigStorageLifecycleScenarioV1(scenario string, params *StorageParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	workspaceResource := generators.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)

	blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageUrl := generators.GenerateImageURL(storageProviderV1, params.Tenant, params.Image.Name)

	blockResource := generators.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageResource := generators.GenerateImageResource(params.Tenant, params.Image.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).Resource(workspaceResource).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	setCreatedRegionalResourceMetadata(workspaceResponse.Metadata)
	workspaceResponse.Status = newWorkspaceStatus(schema.ResourceStateCreating)
	workspaceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	setWorkspaceState(workspaceResponse.Status, schema.ResourceStateActive)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: "GetCreatedWorkspace", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).Resource(blockResource).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = newBlockStorageStatus(schema.ResourceStateCreating)
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the created block storage
	setBlockStorageState(blockResponse.Status, schema.ResourceStateActive)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "UpdateBlockStorage"}); err != nil {
		return nil, err
	}

	// Update the block storage
	setModifiedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	setBlockStorageState(blockResponse.Status, schema.ResourceStateUpdating)
	blockResponse.Spec = *params.BlockStorage.UpdatedSpec
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "UpdateBlockStorage", nextState: "GetUpdatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the updated block storage
	setBlockStorageState(blockResponse.Status, schema.ResourceStateActive)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetUpdatedBlockStorage", nextState: "CreateImage"}); err != nil {
		return nil, err
	}

	// Image
	imageResponse, err := builders.NewImageBuilder().
		Name(params.Image.Name).Resource(imageResource).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Spec(params.Image.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create an image
	setCreatedRegionalResourceMetadata(imageResponse.Metadata)
	imageResponse.Status = newImageStatus(schema.ResourceStateCreating)
	imageResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "CreateImage", nextState: "GetCreatedImage"}); err != nil {
		return nil, err
	}

	// Get the created image
	setImageState(imageResponse.Status, schema.ResourceStateActive)
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "GetCreatedImage", nextState: "UpdateImage"}); err != nil {
		return nil, err
	}

	// Update the image
	setModifiedRegionalResourceMetadata(imageResponse.Metadata)
	setImageState(imageResponse.Status, schema.ResourceStateUpdating)
	imageResponse.Spec = *params.Image.UpdatedSpec
	imageResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "UpdateImage", nextState: "GetUpdatedImage"}); err != nil {
		return nil, err
	}

	// Get the updated image
	setImageState(imageResponse.Status, schema.ResourceStateActive)
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
