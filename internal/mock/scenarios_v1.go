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
		nextState:    createdScenarioState,
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
		currentState: startedScenarioState,
		nextState:    creatingScenarioState,
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
		currentState: creatingScenarioState,
		nextState:    updatingScenarioState,
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
		currentState: updatingScenarioState,
		nextState:    updatedScenarioState,
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
		currentState: updatedScenarioState,
		nextState:    deletingScenarioState,
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
		currentState: redeletingScenarioState,
		nextState:    startedScenarioState,
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
		currentState: startedScenarioState,
		nextState:    createdScenarioState,
		httpStatus:   http.StatusCreated,
		priority:     defaultScenarioPriority,
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
		currentState: startedScenarioState,
		nextState:    creatingScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
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
		currentState: creatingScenarioState,
		nextState:    updatingScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
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
		currentState: updatingScenarioState,
		nextState:    updatedScenarioState,
		httpStatus:   http.StatusOK,
		priority:     defaultScenarioPriority,
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
		currentState: updatedScenarioState,
		nextState:    deletingScenarioState,
		httpStatus:   http.StatusAccepted,
		priority:     defaultScenarioPriority,
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
		currentState: redeletingScenarioState,
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
		priority:     defaultScenarioPriority,
	}); err != nil {
		log.Printf("Failed to configure GET deleted Role-Assignment stub for scenario %q: %v", scenario, err)
		return nil, err
	}

	return wm, nil
}
