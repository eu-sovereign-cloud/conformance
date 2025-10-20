package mock

import (
	"log/slog"
	"net/http"

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
	blockUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := secalib.GenerateInstanceURL(params.Tenant, params.Workspace.Name, params.Instance.Name)

	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)
	blockResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace.Name, params.Instance.Name)

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

	// Get created workspace
	setWorkspaceStatusState(workspaceResponse.Status, secalib.ActiveStatusState)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, response: workspaceResponse, currentState: "GetCreatedWorkspace", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse := newBlockStorageResponse(params.BlockStorage.Name, secalib.ComputeProviderV1, blockResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.BlockStorage.InitialSpec)

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = newBlockStorageStatus(secalib.CreatingStatusState)
	blockResponse.Spec.SizeGB = params.BlockStorage.InitialSpec.SizeGB
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

	// Instance
	instanceResponse := newInstanceResponse(params.Instance.Name, secalib.ComputeProviderV1, instanceResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.Instance.InitialSpec)

	// Create an instance
	setCreatedRegionalWorkspaceResourceMetadata(instanceResponse.Metadata)
	instanceResponse.Status = newInstanceStatus(secalib.CreatingStatusState)
	instanceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "CreateInstance", nextState: "GetCreatedInstance"}); err != nil {
		return nil, err
	}

	// Get created instance
	setInstanceStatusState(instanceResponse.Status, secalib.ActiveStatusState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "GetCreatedInstance", nextState: "UpdateInstance"}); err != nil {
		return nil, err
	}

	// Update the instance
	setModifiedRegionalWorkspaceResourceMetadata(instanceResponse.Metadata)
	setInstanceStatusState(instanceResponse.Status, secalib.UpdatingStatusState)
	instanceResponse.Spec = *params.Instance.UpdatedSpec
	instanceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "UpdateInstance", nextState: "GetUpdatedInstance"}); err != nil {
		return nil, err
	}

	// Get updated instance
	setInstanceStatusState(instanceResponse.Status, secalib.ActiveStatusState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "GetUpdatedInstance", nextState: "StopInstance"}); err != nil {
		return nil, err
	}

	// Stop the instance
	instanceResponse.Metadata.Verb = http.MethodPost
	if err := configurePostSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl + "/stop", params: params, response: instanceResponse, currentState: "StopInstance", nextState: "GetStoppedInstance"}); err != nil {
		return nil, err
	}

	// Get the stopped instance
	setInstanceStatusState(instanceResponse.Status, secalib.SuspendedStatusState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "GetStoppedInstance", nextState: "StartInstance"}); err != nil {
		return nil, err
	}

	// Start the instance
	instanceResponse.Metadata.Verb = http.MethodPost
	if err := configurePostSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl + "/start", params: params, response: instanceResponse, currentState: "StartInstance", nextState: "GetStartedInstance"}); err != nil {
		return nil, err
	}

	// Get the started instance
	setInstanceStatusState(instanceResponse.Status, secalib.ActiveStatusState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "GetStartedInstance", nextState: "RestartInstance"}); err != nil {
		return nil, err
	}

	// Restart the instance
	instanceResponse.Metadata.Verb = http.MethodPost
	if err := configurePostSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl + "/restart", params: params, response: instanceResponse, currentState: "RestartInstance", nextState: "GetRestartedInstance"}); err != nil {
		return nil, err
	}

	// Get the restarted instance
	setInstanceStatusState(instanceResponse.Status, secalib.ActiveStatusState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "GetRestartedInstance", nextState: "DeleteInstance"}); err != nil {
		return nil, err
	}

	// Delete the instance
	instanceResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "DeleteInstance", nextState: "GetDeletedInstance"}); err != nil {
		return nil, err
	}

	// Get deleted instance
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "GetDeletedInstance", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}
