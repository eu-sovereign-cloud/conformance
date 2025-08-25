package mock

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wiremock/go-wiremock"
)

func CreateWorkspaceScenarioV1(scenario string, params WorkspaceParamsV1) (*wiremock.Client, error) {
	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	params.MockURL = fmt.Sprintf(workspaceURLV1, params.Tenant, params.Name)

	response := workspaceResponseV1{
		Metadata: metadataResponse{
			Name:       params.Name,
			Provider:   workspaceProviderV1,
			Resource:   fmt.Sprintf(workspaceResource, params.Tenant, params.Name),
			ApiVersion: version1,
			Kind:       workspaceKind,
			Tenant:     params.Tenant,
			Region:     params.Region,
		},
	}

	// Create a workspace
	response.Metadata.Verb = http.MethodPut
	response.Status.State = creatingStatusState
	response.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = 1
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		template:     workspaceResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    creatingScenarioState,
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get created workspace
	response.Metadata.Verb = http.MethodGet
	response.Status.State = activeStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		template:     workspaceResponseTemplateV1,
		currentState: creatingScenarioState,
		nextState:    createdScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Update the workspace
	response.Metadata.Verb = http.MethodPut
	response.Status.State = updatingStatusState
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = response.Metadata.ResourceVersion + 1
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		template:     workspaceResponseTemplateV1,
		currentState: createdScenarioState,
		nextState:    updatingScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get updated workspace
	response.Metadata.Verb = http.MethodGet
	response.Status.State = activeStatusState
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		template:     workspaceResponseTemplateV1,
		currentState: updatingScenarioState,
		nextState:    updatedScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Delete the workspace
	response.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		currentState: updatedScenarioState,
		nextState:    deletingScenarioState,
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Re-delete the workspace
	response.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		currentState: deletingScenarioState,
		nextState:    redeletingScenarioState,
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get deleted workspace
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		currentState: redeletingScenarioState,
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}

func CreateComputeScenarioV1(scenario string, params ComputeParamsV1) (*wiremock.Client, error) {
	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	// Sku
	params.MockURL = fmt.Sprintf(instanceSkuURLV1, params.Tenant, params.Workspace, params.Sku.Name)

	// Get sku
	skuResponse := instanceSkuResponseV1{
		metadata: metadataResponse{
			Name:            params.Sku.Name,
			Provider:        computeProviderV1,
			Resource:        fmt.Sprintf(instanceSkuResource, params.Tenant, params.Workspace, params.Sku.Name),
			Verb:            http.MethodGet,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            instanceSkuKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		status: statusResponse{
			State:            activeStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
		architecture: params.Sku.Architecture,
		provider:     params.Sku.Provider,
		tier:         params.Sku.Tier,
		ram:          params.Sku.RAM,
		vCPU:         params.Sku.VCPU,
	}
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     skuResponse,
		template:     instanceSkuResponseTemplateV1,
		currentState: creatingScenarioState,
		nextState:    createdScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Instance
	params.MockURL = fmt.Sprintf(instanceURLV1, params.Tenant, params.Workspace, params.Instance.Name)

	instResponse := instanceResponseV1{
		metadata: metadataResponse{
			Name:       params.Instance.Name,
			Provider:   computeProviderV1,
			Resource:   fmt.Sprintf(instanceResource, params.Tenant, params.Workspace, params.Instance.Name),
			ApiVersion: version1,
			Kind:       instanceKind,
			Tenant:     params.Tenant,
			Region:     params.Region,
		},
	}

	// Create an instance
	instResponse.metadata.Verb = http.MethodPut
	instResponse.status.State = creatingStatusState
	instResponse.metadata.CreatedAt = time.Now().Format(time.RFC3339)
	instResponse.metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	instResponse.metadata.ResourceVersion = 1
	instResponse.status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    creatingScenarioState,
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get created instance
	instResponse.metadata.Verb = http.MethodGet
	instResponse.status.State = activeStatusState
	instResponse.status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: creatingScenarioState,
		nextState:    createdScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Update the instance
	instResponse.metadata.Verb = http.MethodPut
	instResponse.status.State = updatingStatusState
	instResponse.metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	instResponse.metadata.ResourceVersion = instResponse.metadata.ResourceVersion + 1
	instResponse.status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: createdScenarioState,
		nextState:    updatingScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get updated instance
	instResponse.metadata.Verb = http.MethodGet
	instResponse.status.State = activeStatusState
	instResponse.status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: updatingScenarioState,
		nextState:    updatedScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Power-off the instance
	instResponse.metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		currentState: updatingScenarioState,
		nextState:    updatedScenarioState,
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get the powered-off instance

	// Power-on the instance
	instResponse.metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		currentState: updatingScenarioState,
		nextState:    updatedScenarioState,
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get the powered-on instance

	// Restart the instance
	instResponse.metadata.Verb = http.MethodPost
	if err := configurePostStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		currentState: updatingScenarioState,
		nextState:    updatedScenarioState,
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get the restarted instance

	// Delete the instance
	instResponse.metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		currentState: updatedScenarioState,
		nextState:    deletingScenarioState,
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get deleted instance
	instResponse.metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		currentState: redeletingScenarioState,
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}

func CreateStorageScenarioV1(scenario string, mockParams Params) (*wiremock.Client, error) {
	return nil, nil
}

func CreateAuthorizationScenarioV1(scenario string, params AuthorizationParamsV1) (*wiremock.Client, error) {
	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	params.MockURL = fmt.Sprintf(rolesURLV1, params.Tenant, params.roles.Name)

	response := rolesResponseV1{
		Metadata: metadataResponse{
			Name:            params.roles.Name,
			Provider:        authorizationProviderV1,
			Resource:        fmt.Sprintf(rolesResource, params.Tenant, params.roles.Name),
			Verb:            http.MethodPut,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            rolesKind,
			Tenant:          params.Tenant,
		},
		Permissions: []PermissionsParamsV1{
			{
				Provider:  authorizationProviderV1,
				Resources: []string{fmt.Sprintf(rolesResource, params.Tenant, params.roles.Name)},
				Verbs:     []string{http.MethodPut},
			},
		},
	}

	// Create a Role
	response.Metadata.Verb = http.MethodPut
	response.Status.State = creatingStatusState
	response.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = 1
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		template:     roleResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    "Get created Role",
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		log.Printf("Failed to configure Create a Role PUT stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	// Get created Role
	response.Metadata.Verb = http.MethodGet
	response.Status.State = activeStatusState
	response.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = 1
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		template:     roleResponseTemplateV1,
		currentState: "Get created Role",
		nextState:    "Update the Role",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		log.Printf("Failed to configure Get created Role stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	// Update the Role
	response.Metadata.Verb = http.MethodPut
	response.Status.State = updatingStatusState
	response.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = 1
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		template:     roleResponseTemplateV1,
		currentState: "Update the Role",
		nextState:    "Get updated Role",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		log.Printf("Failed to configure Update the Role stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	// Get updated Role
	response.Metadata.Verb = http.MethodGet
	response.Status.State = activeStatusState
	response.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	response.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	response.Metadata.ResourceVersion = 1
	response.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		template:     roleResponseTemplateV1,
		currentState: "Get updated Role",
		nextState:    "Delete the Role",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		log.Printf("Failed to configure Get updated Role stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	// Delete the Role
	response.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		currentState: "Delete the Role",
		nextState:    "Get deleted Role",
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		log.Printf("Failed to configure Delete the Role stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	// Get deleted Role
	response.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     response,
		currentState: "Get deleted Role",
		nextState:    "Create a Role-Assignment",
		httpStatus:   http.StatusNotFound,
		priority:     defaultScenarioPriority,
	}); err != nil {
		log.Printf("Failed to configure Get deleted Role stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	params.MockURL = fmt.Sprintf(roleAssignmentURLV1, params.Tenant, params.roleAssignment.Name)

	responseRA := roleAssignmentResponseV1{
		Metadata: metadataResponse{
			Name:            params.roleAssignment.Name,
			Provider:        authorizationProviderV1,
			Resource:        fmt.Sprintf(roleAssignmentResource, params.Tenant, params.roleAssignment.Name),
			Verb:            http.MethodPut,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            roleAssignmentKind,
			Tenant:          params.Tenant,
		},
		Subs: []string{
			"sub1",
			"sub2",
		},
		Roles: []string{
			"role1",
		},
		Scopes: []Scopes{
			{
				Tenants:    []string{"tenant1", "tenant2"},
				Regions:    []string{"region1", "region2"},
				Workspaces: []string{"workspace1", "workspace2"},
			},
		},
	}

	// Create a Role-Assignment
	responseRA.Metadata.Verb = http.MethodPut
	responseRA.Status.State = creatingStatusState
	responseRA.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	responseRA.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	responseRA.Metadata.ResourceVersion = 1
	responseRA.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     responseRA,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "Create a Role-Assignment",
		nextState:    "Get created Role-Assignment",
		httpStatus:   http.StatusCreated,
		priority:     2,
	}); err != nil {
		log.Printf("Failed to configure Create a Role-assignment PUT stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	// Get created Role-Assignment
	responseRA.Metadata.Verb = http.MethodGet
	responseRA.Status.State = activeStatusState
	responseRA.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	responseRA.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	responseRA.Metadata.ResourceVersion = 1
	responseRA.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     responseRA,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "Get created Role-Assignment",
		nextState:    "Update the Role-Assignment",
		httpStatus:   http.StatusOK,
		priority:     2,
	}); err != nil {
		log.Printf("Failed to configure Get created Role-Assignment Get stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	// Update the Role-Assignment
	responseRA.Metadata.Verb = http.MethodPut
	responseRA.Status.State = updatingStatusState
	responseRA.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	responseRA.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	responseRA.Metadata.ResourceVersion = 1
	responseRA.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     responseRA,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "Update the Role-Assignment",
		nextState:    "Get updated Role-Assignment",
		httpStatus:   http.StatusOK,
		priority:     2,
	}); err != nil {
		log.Printf("Failed to configure Update the Role-Assignment PUT stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	// Get updated Role-Assignment
	responseRA.Metadata.Verb = http.MethodGet
	responseRA.Status.State = activeStatusState
	responseRA.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	responseRA.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	responseRA.Metadata.ResourceVersion = 1
	responseRA.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     responseRA,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "Get updated Role-Assignment",
		nextState:    "Delete the Role-Assignment",
		httpStatus:   http.StatusOK,
		priority:     2,
	}); err != nil {
		log.Printf("Failed to configure Get updated Role-Assignment stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	// Delete the Role-Assignment
	responseRA.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     responseRA,
		currentState: "Delete the Role-Assignment",
		nextState:    "Get deleted Role-Assignment",
		httpStatus:   http.StatusAccepted,
		priority:     2,
	}); err != nil {
		log.Printf("Failed to configure DELETE Role-Assignment stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	// Get deleted Role-Assignment
	responseRA.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     responseRA,
		currentState: "Get deleted Role-Assignment",
		httpStatus:   http.StatusNotFound,
		priority:     2,
	}); err != nil {
		log.Printf("Failed to configure GET deleted Role-Assignment stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	return wm, nil
}

func CreateNetworkScenarioV1(scenario string, params NetworkParamsV1) (*wiremock.Client, error) {
	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	// Network Sku
	params.MockURL = fmt.Sprintf(networkSkuURLV1, params.Tenant, params.Workspace, params.NetworkSku.Name)
	networkSkuResponse := networkSkuResponseV1{
		Metadata: metadataResponse{
			Name:   params.NetworkSku.Name,
			Tenant: params.Tenant,
		},
	}
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     networkSkuResponse,
		template:     networkSkuResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    "Get network",
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}
	// Network
	params.MockURL = fmt.Sprintf(networkURLV1, params.Tenant, params.Workspace, params.Network.Name)
	// Create  Network
	networkResponse := networkResponseV1{
		Metadata: metadataResponse{
			Name:            params.Network.Name,
			Provider:        networkProviderV1,
			Resource:        fmt.Sprintf(instanceSkuResource, params.Tenant, params.Workspace, params.Network.Name),
			Verb:            http.MethodPut,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            networkKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            activeStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
		SkuRef:        params.Network.SkuRef,
		RouteTableRef: params.Network.RouteTableRef,
	}

	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "Create network",
		nextState:    "Get network",
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	//Get network
	networkResponse.Metadata.Verb = http.MethodGet
	networkResponse.Status.State = activeStatusState
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "Get network",
		nextState:    "Update network",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	//Update Network
	networkResponse.Metadata.Verb = http.MethodPut
	networkResponse.Status.State = updatingScenarioState
	networkResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	networkResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "Update network",
		nextState:    "Get network2x",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}
	//Get network 2x time
	networkResponse.Metadata.Verb = http.MethodGet
	networkResponse.Status.State = activeStatusState
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "Get network2x",
		nextState:    "Create internet-gateway",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// internet-Gateway
	internetGatewayResponse := internetGatewayResponseV1{
		Metadata: metadataResponse{
			Name:            params.InternetGateway.Name,
			Provider:        networkProviderV1,
			Resource:        fmt.Sprintf(internetGatewayURLV1, params.Tenant, params.Workspace, params.InternetGateway.Name),
			Verb:            http.MethodPut,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            internetGatewayKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            activeStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
		EgressOnly: params.InternetGateway.EgressOnly,
	}

	// Create internet-Gateway
	internetGatewayResponse.Metadata.Verb = http.MethodPut
	internetGatewayResponse.Status.State = activeStatusState
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "Create internet-gateway",
		nextState:    "Get internet-gateway",
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get internet-Gateway
	internetGatewayResponse.Metadata.Verb = http.MethodGet
	internetGatewayResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "Get internet-gateway",
		nextState:    "Update internet-gateway",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	//Update internet-gateway
	internetGatewayResponse.Metadata.Verb = http.MethodPut
	internetGatewayResponse.Status.State = updatingScenarioState
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "Update internet-gateway",
		nextState:    "Get internet-gateway updated",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get internet-gateway after update
	internetGatewayResponse.Metadata.Verb = http.MethodGet
	internetGatewayResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "Get internet-gateway2x",
		nextState:    "Create route-table",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Route-Table
	routeTableResponse := routeTableResponseV1{
		Metadata: metadataResponse{
			Name:            params.RouteTable.Name,
			Provider:        networkProviderV1,
			Resource:        fmt.Sprintf(routeTableURLV1, params.Tenant, params.Workspace, params.RouteTable.Name),
			Verb:            http.MethodPost,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            routeTableKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            activeStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
	}

	// Create route-Table
	routeTableResponse.Metadata.Verb = http.MethodPost
	routeTableResponse.Status.State = activeStatusState
	if err := configurePostStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "Create route-table",
		nextState:    "Get route-table",
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get route-Table
	routeTableResponse.Metadata.Verb = http.MethodGet
	routeTableResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "Get route-table",
		nextState:    "Update route-table",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Update route-Table
	routeTableResponse.Metadata.Verb = http.MethodPut
	routeTableResponse.Status.State = updatingScenarioState
	routeTableResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	routeTableResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "Update route-table",
		nextState:    "Get route-table updated",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get route-table after update
	routeTableResponse.Metadata.Verb = http.MethodGet
	routeTableResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "Get route-table updated",
		nextState:    "Create subnet",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	//subnet

	subnetResponse := subnetResponseV1{
		Metadata: metadataResponse{
			Name:            params.Subnet.Name,
			Provider:        networkProviderV1,
			Resource:        fmt.Sprintf(subnetURLV1, params.Tenant, params.Workspace, params.Subnet.Name),
			Verb:            http.MethodPost,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            subnetKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            activeStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
	}

	// Create subnet
	subnetResponse.Metadata.Verb = http.MethodPost
	subnetResponse.Status.State = activeStatusState
	if err := configurePostStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "Create subnet",
		nextState:    "Get subnet",
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get subnet
	subnetResponse.Metadata.Verb = http.MethodGet
	subnetResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "Get subnet",
		nextState:    "Update subnet",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Update subnet
	subnetResponse.Metadata.Verb = http.MethodPut
	subnetResponse.Status.State = updatingScenarioState
	subnetResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	subnetResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "Update subnet",
		nextState:    "Get subnet updated",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get subnet after update
	subnetResponse.Metadata.Verb = http.MethodGet
	subnetResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "Get subnet updated",
		nextState:    "Create public-ip",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Public-IP
	publicIPResponse := publicIPResponseV1{
		Metadata: metadataResponse{
			Name:            params.PublicIP.Name,
			Provider:        networkProviderV1,
			Resource:        fmt.Sprintf(publicIPURLV1, params.Tenant, params.Workspace, params.PublicIP.Name),
			Verb:            http.MethodPost,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            publicIPKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            activeStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
	}

	// Create public-IP
	publicIPResponse.Metadata.Verb = http.MethodPost
	publicIPResponse.Status.State = activeStatusState
	if err := configurePostStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "Create public-ip",
		nextState:    "Get public-ip",
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get public-IP
	publicIPResponse.Metadata.Verb = http.MethodGet
	publicIPResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "Get public-ip",
		nextState:    "Update public-ip",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Update public-IP
	publicIPResponse.Metadata.Verb = http.MethodPut
	publicIPResponse.Status.State = updatingScenarioState
	publicIPResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "Update public-ip",
		nextState:    "Get public-ip updated",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get public-IP after update
	publicIPResponse.Metadata.Verb = http.MethodGet
	publicIPResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "Get public-ip updated",
		nextState:    "Create nic",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// NIC
	nicResponse := nicResponseV1{
		Metadata: metadataResponse{
			Name:            params.NIC.Name,
			Provider:        networkProviderV1,
			Resource:        fmt.Sprintf(nicURLV1, params.Tenant, params.Workspace, params.NIC.Name),
			Verb:            http.MethodPost,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            nicKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            activeStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
	}

	// Create NIC
	nicResponse.Metadata.Verb = http.MethodPost
	nicResponse.Status.State = activeStatusState
	if err := configurePostStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "Create nic",
		nextState:    "Get nic",
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get NIC
	nicResponse.Metadata.Verb = http.MethodGet
	nicResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "Get nic",
		nextState:    "Update nic",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Update NIC
	nicResponse.Metadata.Verb = http.MethodPut
	nicResponse.Status.State = updatingScenarioState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "Update nic",
		nextState:    "Get nic updated",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get NIC after update
	nicResponse.Metadata.Verb = http.MethodGet
	nicResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "Get nic updated",
		nextState:    "Create security-group",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Security-group
	securityGroupResponse := securityGroupResponseV1{
		Metadata: metadataResponse{
			Name:            params.SecurityGroup.Name,
			Provider:        networkProviderV1,
			Resource:        fmt.Sprintf(securityGroupURLV1, params.Tenant, params.Workspace, params.SecurityGroup.Name),
			Verb:            http.MethodPost,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            SecurityGroupKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            activeStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
	}

	// Create Security-group
	securityGroupResponse.Metadata.Verb = http.MethodPost
	securityGroupResponse.Status.State = activeStatusState
	if err := configurePostStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "Create security-group",
		nextState:    "Get security-group",
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get Security-group
	securityGroupResponse.Metadata.Verb = http.MethodGet
	securityGroupResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "Get security-group",
		nextState:    "Update security-group",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Update Security-group
	securityGroupResponse.Metadata.Verb = http.MethodPut
	securityGroupResponse.Status.State = updatingScenarioState
	securityGroupResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	securityGroupResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "Update security-group",
		nextState:    "Get security-group updated",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get Security-group after update
	securityGroupResponse.Metadata.Verb = http.MethodGet
	securityGroupResponse.Status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "Get security-group updated",
		nextState:    "Get instance sku",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Instance
	// Sku
	params.MockURL = fmt.Sprintf(instanceSkuURLV1, params.Tenant, params.Workspace, params.InstanceSku.Name)

	// Get sku
	skuResponse := instanceSkuResponseV1{
		metadata: metadataResponse{
			Name:            params.InstanceSku.Name,
			Provider:        computeProviderV1,
			Resource:        fmt.Sprintf(instanceSkuResource, params.Tenant, params.Workspace, params.InstanceSku.Name),
			Verb:            http.MethodGet,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      version1,
			Kind:            instanceSkuKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		status: statusResponse{
			State:            activeStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
		architecture: params.InstanceSku.Architecture,
		provider:     params.InstanceSku.Provider,
		tier:         params.InstanceSku.Tier,
		ram:          params.InstanceSku.RAM,
		vCPU:         params.InstanceSku.VCPU,
	}
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     skuResponse,
		template:     instanceSkuResponseTemplateV1,
		currentState: "Get instance sku",
		nextState:    "Create instance",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Instance
	params.MockURL = fmt.Sprintf(instanceURLV1, params.Tenant, params.Workspace, params.Instance.Name)

	instResponse := instanceResponseV1{
		metadata: metadataResponse{
			Name:       params.Instance.Name,
			Provider:   computeProviderV1,
			Resource:   fmt.Sprintf(instanceResource, params.Tenant, params.Workspace, params.Instance.Name),
			ApiVersion: version1,
			Kind:       instanceKind,
			Tenant:     params.Tenant,
			Region:     params.Region,
		},
	}

	// Create an instance
	instResponse.metadata.Verb = http.MethodPut
	instResponse.status.State = creatingStatusState
	instResponse.metadata.CreatedAt = time.Now().Format(time.RFC3339)
	instResponse.metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	instResponse.metadata.ResourceVersion = 1
	instResponse.status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: "Create instance",
		nextState:    "Get instance",
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Get Instance
	instResponse.metadata.Verb = http.MethodGet
	instResponse.status.State = activeStatusState
	if err := configureGetStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: "Get instance",
		nextState:    "Delete instance",
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	//Delete Instance
	instResponse.metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     instResponse,
		currentState: "Delete instance",
		nextState:    "Delete Security Group",
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Delete Security Group
	securityGroupResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     securityGroupResponse,
		currentState: "Delete Security Group",
		nextState:    "Delete Nic",
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Delete Nic
	nicResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     nicResponse,
		currentState: "Delete Nic",
		nextState:    "Delete Public IP",
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Delete public ip
	publicIPResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     publicIPResponse,
		currentState: "Delete Public IP",
		nextState:    "Delete Subnet",
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Delete subnet
	subnetResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     subnetResponse,
		currentState: "Delete Subnet",
		nextState:    "Delete Route-table",
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Delete Route-table
	routeTableResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     routeTableResponse,
		currentState: "Delete Route-table",
		nextState:    "Delete Internet-gateway",
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Delete Internet-gateway
	internetGatewayResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     internetGatewayResponse,
		currentState: "Delete Internet-gateway",
		nextState:    "Delete Network",
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}

	// Delete Network
	networkResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenarioConfig{
		name:         scenario,
		params:       params,
		response:     networkResponse,
		currentState: "Delete Network",
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
	}); err != nil {
		return nil, err
	}
	// End
	return wm, nil
}
