package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/conformance/secalib/builders"
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
	instanceUrl := secalib.GenerateInstanceURL(params.Tenant, params.Workspace.Name, params.Instance.Name)

	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)
	blockResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace.Name, params.Instance.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(secalib.WorkspaceProviderV1).
		Resource(workspaceResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		BuildResponse()
	if err != nil {
		return nil, err
	}

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
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(secalib.StorageProviderV1).
		Resource(blockResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

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
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(params.Instance.Name).
		Provider(secalib.ComputeProviderV1).
		Resource(instanceResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec(params.Instance.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

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
	instanceResponse.Spec = *params.Instance.UpdatedSpec
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
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, currentState: "GetDeletedInstance", nextState: "DeleteBlockStorage"}); err != nil {
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
