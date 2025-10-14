package mock

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

func CreateComputeLifecycleScenarioV1(scenario string, params *ComputeParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)

	blockStorageUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := secalib.GenerateInstanceURL(params.Tenant, params.Workspace.Name, params.Instance.Name)

	blockStorageResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace.Name, params.Instance.Name)

	// Workspace
	workResponse := &resourceResponse[secalib.WorkspaceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Workspace.Name,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     params.Tenant,
			Region:     &params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          workspaceUrl,
		params:       params,
		response:     workResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          workspaceUrl,
		params:       params,
		response:     workResponse,
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
			Tenant:     params.Tenant,
			Region:     &params.Region,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
		currentState: "GetCreatedBlockStorage",
		nextState:    "CreateInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Instance
	instResponse := &resourceResponse[secalib.InstanceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Instance.Name,
			Provider:   secalib.ComputeProviderV1,
			Resource:   instanceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  &params.Workspace.Name,
			Region:     &params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.InstanceSpecV1{
			SkuRef:        params.Instance.InitialSpec.SkuRef,
			Zone:          params.Instance.InitialSpec.Zone,
			BootDeviceRef: params.Instance.InitialSpec.BootDeviceRef,
		},
	}

	// Create an instance
	instResponse.Metadata.Verb = http.MethodPut
	instResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	instResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	instResponse.Metadata.ResourceVersion = 1
	instResponse.Status.State = secalib.CreatingStatusState
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "CreateInstance",
		nextState:    "GetCreatedInstance",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created instance
	instResponse.Metadata.Verb = http.MethodGet
	instResponse.Status.State = secalib.ActiveStatusState
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "GetCreatedInstance",
		nextState:    "UpdateInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the instance
	instResponse.Metadata.Verb = http.MethodPut
	instResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	instResponse.Metadata.ResourceVersion = instResponse.Metadata.ResourceVersion + 1
	instResponse.Spec = params.Instance.UpdatedSpec
	instResponse.Status.State = secalib.UpdatingStatusState
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "UpdateInstance",
		nextState:    "GetUpdatedInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get updated instance
	instResponse.Metadata.Verb = http.MethodGet
	instResponse.Status.State = secalib.ActiveStatusState
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "GetUpdatedInstance",
		nextState:    "StopInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Stop the instance
	instResponse.Metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenario, stubConfig{
		url:          instanceUrl + "/stop",
		params:       params,
		response:     instResponse,
		currentState: "StopInstance",
		nextState:    "GetStoppedInstance",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get the stopped instance
	instResponse.Metadata.Verb = http.MethodGet
	instResponse.Status.State = secalib.SuspendedStatusState
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "GetStoppedInstance",
		nextState:    "StartInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Start the instance
	instResponse.Metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenario, stubConfig{
		url:          instanceUrl + "/start",
		params:       params,
		response:     instResponse,
		currentState: "StartInstance",
		nextState:    "GetStartedInstance",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get the started instance
	instResponse.Metadata.Verb = http.MethodGet
	instResponse.Status.State = secalib.ActiveStatusState
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "GetStartedInstance",
		nextState:    "RestartInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Restart the instance
	instResponse.Metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenario, stubConfig{
		url:          instanceUrl + "/restart",
		params:       params,
		response:     instResponse,
		currentState: "RestartInstance",
		nextState:    "GetRestartedInstance",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get the restarted instance
	instResponse.Metadata.Verb = http.MethodGet
	instResponse.Status.State = secalib.ActiveStatusState
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "GetRestartedInstance",
		nextState:    "DeleteInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Delete the instance
	instResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "DeleteInstance",
		nextState:    "GetDeletedInstance",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted instance
	instResponse.Metadata.Verb = http.MethodGet
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "GetDeletedInstance",
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
