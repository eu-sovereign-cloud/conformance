package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	computeV1 "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.compute.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/wiremock/go-wiremock"
)

func ConfigComputeLifecycleScenarioV1(scenario string, params *ComputeParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	blockUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := secalib.GenerateInstanceURL(params.Tenant, params.Workspace.Name, (*params.Instance)[0].Name)

	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)
	blockResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace.Name, (*params.Instance)[0].Name)

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
	blockResponse := newBlockStorageResponse(params.BlockStorage.Name, secalib.ComputeProviderV1, blockResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.BlockStorage.InitialLabels,
		params.BlockStorage.InitialSpec)

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = secalib.NewBlockStorageStatus(secalib.CreatingResourceState)
	blockResponse.Spec.SizeGB = params.BlockStorage.InitialSpec.SizeGB
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the created block storage
	secalib.SetBlockStorageStatusState(blockResponse.Status, secalib.ActiveResourceState)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "CreateInstance"}); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse := newInstanceResponse((*params.Instance)[0].Name, secalib.ComputeProviderV1, instanceResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		&(*params.Instance)[0].InitialLabels,
		(*params.Instance)[0].InitialSpec)

	// Create an instance
	setCreatedRegionalWorkspaceResourceMetadata(instanceResponse.Metadata)
	instanceResponse.Status = secalib.NewInstanceStatus(secalib.CreatingResourceState)
	instanceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, responseBody: instanceResponse, currentState: "CreateInstance", nextState: "GetCreatedInstance"}); err != nil {
		return nil, err
	}

	// Get the created instance
	secalib.SetInstanceStatusState(instanceResponse.Status, secalib.ActiveResourceState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, responseBody: instanceResponse, currentState: "GetCreatedInstance", nextState: "UpdateInstance"}); err != nil {
		return nil, err
	}

	// Update the instance
	setModifiedRegionalWorkspaceResourceMetadata(instanceResponse.Metadata)
	secalib.SetInstanceStatusState(instanceResponse.Status, secalib.UpdatingResourceState)
	instanceResponse.Spec = *(*params.Instance)[0].UpdatedSpec
	instanceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, responseBody: instanceResponse, currentState: "UpdateInstance", nextState: "GetUpdatedInstance"}); err != nil {
		return nil, err
	}

	// Get the updated instance
	secalib.SetInstanceStatusState(instanceResponse.Status, secalib.ActiveResourceState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, responseBody: instanceResponse, currentState: "GetUpdatedInstance", nextState: "StopInstance"}); err != nil {
		return nil, err
	}

	// Stop the instance
	instanceResponse.Metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenario,
		&stubConfig{url: instanceUrl + "/stop", params: params, responseBody: instanceResponse, currentState: "StopInstance", nextState: "GetStoppedInstance"}); err != nil {
		return nil, err
	}

	// Get the stopped instance
	secalib.SetInstanceStatusState(instanceResponse.Status, secalib.SuspendedResourceState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, responseBody: instanceResponse, currentState: "GetStoppedInstance", nextState: "StartInstance"}); err != nil {
		return nil, err
	}

	// Start the instance
	instanceResponse.Metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenario,
		&stubConfig{url: instanceUrl + "/start", params: params, responseBody: instanceResponse, currentState: "StartInstance", nextState: "GetStartedInstance"}); err != nil {
		return nil, err
	}

	// Get the started instance
	secalib.SetInstanceStatusState(instanceResponse.Status, secalib.ActiveResourceState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, responseBody: instanceResponse, currentState: "GetStartedInstance", nextState: "RestartInstance"}); err != nil {
		return nil, err
	}

	// Restart the instance
	instanceResponse.Metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenario,
		&stubConfig{url: instanceUrl + "/restart", params: params, responseBody: instanceResponse, currentState: "RestartInstance", nextState: "GetRestartedInstance"}); err != nil {
		return nil, err
	}

	// Get the restarted instance
	secalib.SetInstanceStatusState(instanceResponse.Status, secalib.ActiveResourceState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, responseBody: instanceResponse, currentState: "GetRestartedInstance", nextState: "DeleteInstance"}); err != nil {
		return nil, err
	}

	// Delete the instance
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, currentState: "DeleteInstance", nextState: "GetDeletedInstance"}); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: instanceUrl, params: params, currentState: "GetDeletedInstance", nextState: "DeleteBlockStorage"}); err != nil {
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

func ConfigComputeListLifecycleScenarioV1(scenario string, params *ComputeParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	blockUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)

	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)
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
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: startedScenarioState, nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse := newBlockStorageResponse(params.BlockStorage.Name, secalib.ComputeProviderV1, blockResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.BlockStorage.InitialLabels,
		params.BlockStorage.InitialSpec)

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = secalib.NewBlockStorageStatus(secalib.CreatingResourceState)
	blockResponse.Spec.SizeGB = params.BlockStorage.InitialSpec.SizeGB
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "CreateBlockStorage", nextState: (*params.Instance)[0].Name}); err != nil {
		return nil, err
	}

	// Instance
	var instanceList []schema.Instance
	for i := range *params.Instance {
		instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace.Name, (*params.Instance)[i].Name)
		instanceResponse := newInstanceResponse((*params.Instance)[i].Name, secalib.ComputeProviderV1, instanceResource, secalib.ApiVersion1,
			params.Tenant, params.Workspace.Name, params.Region,
			&(*params.Instance)[i].InitialLabels,
			(*params.Instance)[i].InitialSpec)

		var nextState string
		if i < len(*params.Instance)-1 {
			nextState = (*params.Instance)[i+1].Name
		} else {
			nextState = "GetInstancesList"
		}
		// Create an instance
		setCreatedRegionalWorkspaceResourceMetadata(instanceResponse.Metadata)
		instanceResponse.Status = secalib.NewInstanceStatus(secalib.CreatingResourceState)
		instanceResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateInstanceURL(params.Tenant, params.Workspace.Name, (*params.Instance)[i].Name), params: params, responseBody: instanceResponse, currentState: (*params.Instance)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		instanceList = append(instanceList, *instanceResponse)
	}

	// List Instances
	instancesResource := secalib.GenerateInstanceListResource(params.Tenant, params.Workspace.Name)
	instanceResponse := &computeV1.InstanceIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.ComputeProviderV1,
			Resource: instancesResource,
			Verb:     http.MethodGet,
		},
	}
	instanceResponse.Items = instanceList
	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateInstanceListURL(params.Tenant, params.Workspace.Name), params: params, responseBody: instanceResponse,
			currentState: "GetInstancesList", nextState: "GetInstancesListWithLimit",
		}); err != nil {
		return nil, err
	}
	// List Roles with limit 1

	instanceResponse.Items = instanceList[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateInstanceListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimit("1"), responseBody: instanceResponse,
			currentState: "GetInstancesListWithLimit", nextState: "GetInstancesListWithLabel",
		}); err != nil {
		return nil, err
	}
	// List instances with label

	instancesWithLabel := func(instancesList []schema.Instance) []schema.Instance {
		var filteredInstances []schema.Instance
		for _, instance := range instancesList {
			if val, ok := instance.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformance {
				filteredInstances = append(filteredInstances, instance)
			}
		}
		return filteredInstances
	}
	instanceResponse.Items = instancesWithLabel(instanceList)
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateInstanceListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformance), responseBody: instanceResponse, currentState: "GetInstancesListWithLabel", nextState: "GetInstancesListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List instances with limit and label

	instanceResponse.Items = instancesWithLabel(instanceList)[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateInstanceListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformance), responseBody: instanceResponse, currentState: "GetInstancesListWithLimitAndLabel", nextState: "GetListSkus"}); err != nil {
		return nil, err
	}

	// List Sku
	instanceSku := []schema.InstanceSku{
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "D2XS",
				Provider: secalib.ComputeProviderV1,
				Resource: instancesResource,
				Verb:     http.MethodGet,
				Kind:     secalib.InstanceSkuKind,
				Tenant:   params.Tenant,
			},
			Labels: schema.Labels{
				"architecture": "amd64",
				"provider":     "seca",
				"tier":         "D2XS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  1,
				VCPU: 1,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "DXS",
				Provider: secalib.ComputeProviderV1,
				Resource: instancesResource,
				Verb:     http.MethodGet,
				Kind:     secalib.InstanceSkuKind,
				Tenant:   params.Tenant,
			},
			Labels: schema.Labels{
				"architecture": "amd64",
				"provider":     "seca",
				"tier":         "DXS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  2,
				VCPU: 1,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "DS",
				Provider: secalib.ComputeProviderV1,
				Resource: instancesResource,
				Verb:     http.MethodGet,
				Kind:     secalib.InstanceSkuKind,
				Tenant:   params.Tenant,
			},
			Labels: schema.Labels{
				"architecture": "amd64",
				"provider":     "seca",
				"tier":         "DS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  4,
				VCPU: 2,
			},
		},
	}

	skuResponse := &computeV1.SkuIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.ComputeProviderV1,
			Resource: secalib.GenerateSkuListResource(params.Tenant),
			Verb:     http.MethodGet,
		},
		Items: instanceSku,
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateInstanceListURL(params.Tenant, params.Workspace.Name), params: params, responseBody: skuResponse, currentState: "GetListSkus", nextState: "GetListSkusWithLimit"}); err != nil {
		return nil, err
	}

	// List Sku with limit 1

	skuWithLimitResponse := &computeV1.SkuIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.ComputeProviderV1,
			Resource: instancesResource,
			Verb:     http.MethodGet,
		},
		Items: instanceSku[:1],
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateInstanceSkuURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimit("1"), responseBody: skuWithLimitResponse, currentState: "GetListSkusWithLimit", nextState: "GetListSkusWithLabel"}); err != nil {
		return nil, err
	}

	// List Sku with label
	skuWithLabelResponse := &computeV1.SkuIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.ComputeProviderV1,
			Resource: instancesResource,
			Verb:     http.MethodGet,
		},
	}
	skusWithLabel := func(skusList []schema.InstanceSku) []schema.InstanceSku {
		var filteredSkus []schema.InstanceSku
		for _, sku := range skusList {
			if val, ok := sku.Labels["tier"]; ok && val == "D2XS" {
				filteredSkus = append(filteredSkus, sku)
			}
		}
		return filteredSkus
	}
	skuWithLabelResponse.Items = skusWithLabel(instanceSku)
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateInstanceSkuURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformance), responseBody: skuWithLabelResponse, currentState: "GetListSkusWithLabel", nextState: "GetListSkusWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List Sku with limit and label

	skuWithLimitAndLabelResponse := &computeV1.SkuIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.ComputeProviderV1,
			Resource: instancesResource,
			Verb:     http.MethodGet,
		},
	}

	skuWithLimitAndLabelResponse.Items = skusWithLabel(instanceSku)[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateInstanceListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformance), responseBody: skuWithLimitAndLabelResponse, currentState: "GetListSkusWithLimitAndLabel", nextState: "DeleteInstance_" + (*params.Instance)[0].Name}); err != nil {
		return nil, err
	}

	// Delete Instances
	for i := range *params.Instance {
		instanceUrl := secalib.GenerateInstanceURL(params.Tenant, params.Workspace.Name, (*params.Instance)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteInstance_" + (*params.Instance)[i].Name
		} else {
			currentState = "GetDeletedInstance_" + (*params.Instance)[i-1].Name
		}

		nextState = "DeleteInstance_" + (*params.Instance)[i].Name

		// Delete the instance
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: instanceUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted instance (should return 404)
		nextState = func() string {
			if i < len(*params.Instance)-1 {
				return "GetDeletedInstance_" + (*params.Instance)[i].Name
			} else {
				return "DeleteBlockStorage"
			}
		}()

		if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
			&stubConfig{url: instanceUrl, params: params, currentState: "DeleteInstance_" + (*params.Instance)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
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
