package mock

import (
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

func CreateComputeLifecycleScenarioV1(scenario string, params ComputeParamsV1) (*wiremock.Client, error) {
	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	storageSkuUrl := secalib.GenerateStorageSkuURL(params.Tenant, params.StorageSku.Name)
	blockStorageUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace, params.BlockStorage.Name)
	instanceSkuUrl := secalib.GenerateInstanceSkuURL(params.Tenant, params.InstanceSku.Name)
	instanceUrl := secalib.GenerateInstanceURL(params.Tenant, params.Workspace, params.Instance.Name)

	storageResource := secalib.GenerateSkuResource(params.Tenant, params.StorageSku.Name)
	blockStorageResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace, params.BlockStorage.Name)
	instanceSkuResource := secalib.GenerateSkuResource(params.Tenant, params.InstanceSku.Name)
	instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace, params.Instance.Name)

	// Get storage sku
	storageSkuResponse := storageSkuResponseV1{
		Metadata: metadataResponse{
			Name:            params.StorageSku.Name,
			Provider:        secalib.StorageProviderV1,
			Resource:        storageResource,
			Verb:            http.MethodGet,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      secalib.ApiVersion1,
			Kind:            secalib.StorageSkuKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            secalib.ActiveStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
		Provider:      params.StorageSku.Provider,
		Tier:          params.StorageSku.Tier,
		Iops:          params.StorageSku.Iops,
		StorageType:   params.StorageSku.StorageType,
		MinVolumeSize: params.StorageSku.MinVolumeSize,
	}
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          storageSkuUrl,
		params:       params,
		response:     storageSkuResponse,
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
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     params.Tenant,
			Region:     params.Region,
		},
	}

	// Create a block storage
	blockResponse.Metadata.Verb = http.MethodPut
	blockResponse.Status.State = secalib.CreatingStatusState
	blockResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.ResourceVersion = 1
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
	blockResponse.Status.State = secalib.ActiveStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
		template:     blockStorageResponseTemplateV1,
		currentState: "GetCreatedBlockStorage",
		nextState:    "GetInstanceSku",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get instance sku
	instSkuResponse := instanceSkuResponseV1{
		Metadata: metadataResponse{
			Name:            params.InstanceSku.Name,
			Provider:        secalib.ComputeProviderV1,
			Resource:        instanceSkuResource,
			Verb:            http.MethodGet,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      secalib.ApiVersion1,
			Kind:            secalib.InstanceSkuKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            secalib.ActiveStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
		Architecture: params.InstanceSku.Architecture,
		Provider:     params.InstanceSku.Provider,
		Tier:         params.InstanceSku.Tier,
		RAM:          params.InstanceSku.RAM,
		VCPU:         params.InstanceSku.VCPU,
	}
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          instanceSkuUrl,
		params:       params,
		response:     instSkuResponse,
		template:     instanceSkuResponseTemplateV1,
		currentState: "GetInstanceSku",
		nextState:    "CreateInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Instance
	instResponse := instanceResponseV1{
		Metadata: metadataResponse{
			Name:       params.Instance.Name,
			Provider:   secalib.ComputeProviderV1,
			Resource:   instanceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		SkuRef:        params.Instance.SkuRef,
		Zone:          params.Instance.ZoneInitial,
		BootDeviceRef: params.Instance.BootDeviceRef,
	}

	// Create an instance
	instResponse.Metadata.Verb = http.MethodPut
	instResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	instResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	instResponse.Metadata.ResourceVersion = 1
	instResponse.Status.State = secalib.CreatingStatusState
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
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
	instResponse.Zone = params.Instance.ZoneUpdated
	instResponse.Status.State = secalib.UpdatingStatusState
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: "GetUpdatedInstance",
		nextState:    "StopInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Stop the instance
	instResponse.Metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: "GetStoppedInstance",
		nextState:    "StartInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Start the instance
	instResponse.Metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: "GetStartedInstance",
		nextState:    "RestartInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Restart the instance
	instResponse.Metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: "GetRestartedInstance",
		nextState:    "DeleteInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Delete the instance
	instResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "DeleteInstance",
		nextState:    "GetDeletedInstance",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted instance (not found)
	instResponse.Metadata.Verb = http.MethodGet
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
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
