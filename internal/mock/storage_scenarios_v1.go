package mock

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

type StorageScenariosV1 struct {
	Scenarios
}

func NewStorageScenariosV1(authToken string, tenant string, region string, mockURL string) *StorageScenariosV1 {
	return &StorageScenariosV1{
		Scenarios: Scenarios{
			params: secalib.GeneralParams{
				AuthToken: authToken,
				Tenant:    tenant,
				Region:    region,
			},
			mockURL: mockURL,
		},
	}
}

func (scenarios *StorageScenariosV1) ConfigureLifecycleScenario(id string, params *secalib.StorageLifeCycleParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to Storage Lifecycle Scenario")

	name := "StorageLifecycleV1_" + id

	wm, err := scenarios.newClient()
	if err != nil {
		return nil, err
	}

	workspaceUrl := secalib.GenerateWorkspaceURL(scenarios.params.Tenant, params.Workspace.Name)
	workspaceResource := secalib.GenerateWorkspaceResource(scenarios.params.Tenant, params.Workspace.Name)

	blockStorageUrl := secalib.GenerateBlockStorageURL(scenarios.params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageUrl := secalib.GenerateImageURL(scenarios.params.Tenant, params.Image.Name)

	blockStorageResource := secalib.GenerateBlockStorageResource(scenarios.params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageResource := secalib.GenerateImageResource(scenarios.params.Tenant, params.Image.Name)

	// Workspace
	workResponse := &resourceResponse[secalib.WorkspaceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Workspace.Name,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     scenarios.params.Tenant,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Labels: &[]secalib.Label{},
	}

	// Create a workspace
	workResponse.Metadata.Verb = http.MethodPut
	workResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	workResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	workResponse.Metadata.ResourceVersion = 1
	workResponse.Status.State = secalib.CreatingStatusState
	workResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          workspaceUrl,
		params:       scenarios.params,
		response:     workResponse,
		template:     workspaceResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    "GetCreatedWorkspace",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created workspace
	workResponse.Metadata.Verb = http.MethodGet
	workResponse.Status.State = secalib.ActiveStatusState
	workResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          workspaceUrl,
		params:       scenarios.params,
		response:     workResponse,
		template:     workspaceResponseTemplateV1,
		currentState: "GetCreatedWorkspace",
		nextState:    "CreateBlockStorage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse := &resourceResponse[secalib.BlockStorageSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.BlockStorage.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     scenarios.params.Tenant,
			Workspace:  params.Workspace.Name,
			Region:     scenarios.params.Region,
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
	if err := configurePutStub(wm, name, stubConfig{
		url:          blockStorageUrl,
		params:       scenarios.params,
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
	blockResponse.Status.State = secalib.ActiveStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          blockStorageUrl,
		params:       scenarios.params,
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
	blockResponse.Spec = params.BlockStorage.UpdatedSpec
	blockResponse.Status.State = secalib.UpdatingStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          blockStorageUrl,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          blockStorageUrl,
		params:       scenarios.params,
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
			Tenant:     scenarios.params.Tenant,
			Region:     scenarios.params.Region,
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
	if err := configurePutStub(wm, name, stubConfig{
		url:          imageUrl,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          imageUrl,
		params:       scenarios.params,
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
	imageResponse.Spec = params.Image.UpdatedSpec
	imageResponse.Status.State = secalib.UpdatingStatusState
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          imageUrl,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          imageUrl,
		params:       scenarios.params,
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
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          imageUrl,
		params:       scenarios.params,
		response:     imageResponse,
		currentState: "DeleteImage",
		nextState:    "GetDeletedImage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted image
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          imageUrl,
		params:       scenarios.params,
		response:     imageResponse,
		currentState: "GetDeletedImage",
		nextState:    "DeleteBlockStorage",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the block storage
	blockResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          blockStorageUrl,
		params:       scenarios.params,
		response:     blockResponse,
		currentState: "DeleteBlockStorage",
		nextState:    "GetDeletedBlockStorage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted block storage
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          blockStorageUrl,
		params:       scenarios.params,
		response:     blockResponse,
		currentState: "GetDeletedBlockStorage",
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
