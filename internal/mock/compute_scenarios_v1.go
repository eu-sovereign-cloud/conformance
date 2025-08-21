package mock

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wiremock/go-wiremock"
)

func CreateComputeLifecycleScenarioV1(scenario string, params ComputeParamsV1) (*wiremock.Client, error) {
	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	storageSkuUrl := fmt.Sprintf(storageSkuURLV1, params.Tenant, params.StorageSku.Name)
	blockStorageUrl := fmt.Sprintf(blockStorageURLV1, params.Tenant, params.Workspace, params.BlockStorage.Name)
	instanceSkuUrl := fmt.Sprintf(instanceSkuURLV1, params.Tenant, params.InstanceSku.Name)
	instanceUrl := fmt.Sprintf(instanceURLV1, params.Tenant, params.Workspace, params.Instance.Name)

	// Get storage sku
	storageSkuResponse := storageSkuResponseV1{
		Metadata: metadataResponse{
			Name:            params.StorageSku.Name,
			Provider:        storageProviderV1,
			Resource:        fmt.Sprintf(skuResource, params.Tenant, params.StorageSku.Name),
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
			Provider:   storageProviderV1,
			Resource:   fmt.Sprintf(blockStorageResource, params.Tenant, params.Workspace, params.BlockStorage.Name),
			ApiVersion: version1,
			Kind:       blockStorageKind,
			Tenant:     params.Tenant,
			Region:     params.Region,
		},
	}

	// Create a block storage
	blockResponse.Metadata.Verb = http.MethodPut
	blockResponse.Status.State = creatingStatusState
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
	blockResponse.Status.State = activeStatusState
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
			Provider:        computeProviderV1,
			Resource:        fmt.Sprintf(skuResource, params.Tenant, params.InstanceSku.Name),
			Verb:            http.MethodGet,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            instanceSkuKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            activeStatusState,
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
			Provider:   computeProviderV1,
			Resource:   fmt.Sprintf(instanceResource, params.Tenant, params.Workspace, params.Instance.Name),
			ApiVersion: version1,
			Kind:       instanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		SkuRef:        params.Instance.SkuRef,
		Zone:          params.Instance.CreatedZone,
		BootDeviceRef: params.Instance.BootDeviceRef,
	}

	// Create an instance
	instResponse.Metadata.Verb = http.MethodPut
	instResponse.Status.State = creatingStatusState
	instResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	instResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	instResponse.Metadata.ResourceVersion = 1
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
	instResponse.Status.State = activeStatusState
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
	instResponse.Status.State = updatingStatusState
	instResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	instResponse.Metadata.ResourceVersion = instResponse.Metadata.ResourceVersion + 1
	instResponse.Zone = params.Instance.UpdatedZone
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
	instResponse.Status.State = activeStatusState
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
	instResponse.Status.State = suspendedStatusState
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
	instResponse.Status.State = activeStatusState
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
	instResponse.Status.State = activeStatusState
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
