package mock

import (
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/wiremock/go-wiremock"
)

func TestUsageScenario(scenario string, paramsAuthorization AuthorizationParamsV1, paramsCompute ComputeParamsV1, paramsNetwork NetworkParamsV1, paramsStorage StorageParamsV1, paramsWorkspace WorkspaceParamsV1) (*wiremock.Client, error) {

	wm, err := newClient(paramsWorkspace.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs

	// Authorization
	roleUrl := secalib.GenerateRoleURL(paramsAuthorization.Tenant, paramsAuthorization.Role.Name)
	roleAssignmentUrl := secalib.GenerateRoleAssignmentURL(paramsAuthorization.Tenant, paramsAuthorization.RoleAssignment.Name)

	//workspace
	workspaceURL := secalib.GenerateWorkspaceURL(paramsWorkspace.Tenant, paramsWorkspace.Workspace.Name)

	//Storage
	blockStorageURL := secalib.GenerateBlockStorageURL(paramsStorage.Tenant, paramsStorage.Workspace, paramsStorage.BlockStorage.Name)
	imageURL := secalib.GenerateImageURL(paramsStorage.Tenant, paramsStorage.Image.Name)

	//Compute
	instanceURL := secalib.GenerateInstanceURL(paramsCompute.Tenant, paramsCompute.Workspace, paramsCompute.Instance.Name)

	//Network
	networkURL := secalib.GenerateNetworkURL(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.Network.Name)
	internetGatewayURL := secalib.GenerateInternetGatewayURL(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.InternetGateway.Name)
	nicURL := secalib.GenerateNicURL(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.NIC.Name)
	publicIPURL := secalib.GeneratePublicIPURL(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.PublicIP.Name)
	routeTableURL := secalib.GenerateRouteTableURL(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.RouteTable.Name)
	subnetURL := secalib.GenerateSubnetURL(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.Subnet.Name)
	securityGroupURL := secalib.GenerateSecurityGroupURL(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.SecurityGroup.Name)

	// GenerateResources
	// Authorization
	roleResource := secalib.GenerateRoleResource(paramsAuthorization.Tenant, paramsAuthorization.Role.Name)
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(paramsAuthorization.Tenant, paramsAuthorization.RoleAssignment.Name)

	//Workspace

	workspaceResource := secalib.GenerateWorkspaceResource(paramsWorkspace.Tenant, paramsWorkspace.Workspace.Name)

	//Storage
	blockStorageResource := secalib.GenerateBlockStorageResource(paramsStorage.Tenant, paramsStorage.Workspace, paramsStorage.BlockStorage.Name)
	imageResource := secalib.GenerateImageResource(paramsStorage.Tenant, paramsStorage.Image.Name)
	// Compute
	instanceResource := secalib.GenerateInstanceResource(paramsCompute.Tenant, paramsCompute.Workspace, paramsCompute.Instance.Name)

	//Network
	networkResource := secalib.GenerateNetworkResource(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.Network.Name)
	internetGatewayResource := secalib.GenerateInternetGatewayResource(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.InternetGateway.Name)
	nicResource := secalib.GenerateNicResource(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.NIC.Name)
	publicIPResource := secalib.GeneratePublicIPResource(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.PublicIP.Name)
	routeTableResource := secalib.GenerateRouteTableResource(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.RouteTable.Name)
	subnetResource := secalib.GenerateSubnetResource(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.Subnet.Name)
	securityGroupResource := secalib.GenerateSecurityGroupResource(paramsNetwork.Tenant, paramsNetwork.Workspace, paramsNetwork.SecurityGroup.Name)

	//Authorization
	roleResponse := &resourceResponse[secalib.RoleSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsAuthorization.Role.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleKind,
			Tenant:     paramsAuthorization.Tenant,
		},
		Status: &secalib.Status{},
		Spec:   &secalib.RoleSpecV1{},
	}
	for _, perm := range paramsAuthorization.Role.InitialSpec.Permissions {
		roleResponse.Spec.Permissions = append(roleResponse.Spec.Permissions, &secalib.RoleSpecPermissionV1{
			Provider:  perm.Provider,
			Resources: append([]string{}, perm.Resources...),
			Verb:      append([]string{}, perm.Verb...),
		})
	}

	//Create Role
	roleResponse.Metadata.Verb = http.MethodPut
	roleResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	roleResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	roleResponse.Metadata.ResourceVersion = 1
	roleResponse.Status.State = secalib.CreatingStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          roleUrl,
		params:       paramsAuthorization,
		response:     roleResponse,
		template:     roleResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    "GetCreatedRole",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created role
	roleResponse.Metadata.Verb = http.MethodGet
	roleResponse.Status.State = secalib.ActiveStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          roleUrl,
		params:       paramsAuthorization,
		response:     roleResponse,
		template:     roleResponseTemplateV1,
		currentState: "GetCreatedRole",
		nextState:    "UpdateRole",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the role
	roleResponse.Metadata.Verb = http.MethodPut
	roleResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	roleResponse.Metadata.ResourceVersion = roleResponse.Metadata.ResourceVersion + 1
	for i, perm := range paramsAuthorization.Role.UpdatedSpec.Permissions {
		roleResponse.Spec.Permissions[i].Verb = append([]string{}, perm.Verb...)
	}
	roleResponse.Status.State = secalib.UpdatingStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          roleUrl,
		params:       paramsAuthorization,
		response:     roleResponse,
		template:     roleResponseTemplateV1,
		currentState: "UpdateRole",
		nextState:    "GetUpdatedRole",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get updated role
	roleResponse.Metadata.Verb = http.MethodGet
	roleResponse.Status.State = secalib.ActiveStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          roleUrl,
		params:       paramsAuthorization,
		response:     roleResponse,
		template:     roleResponseTemplateV1,
		currentState: "GetUpdatedRole",
		nextState:    "CreateRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignmentResponse := &resourceResponse[secalib.RoleAssignmentSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsAuthorization.RoleAssignment.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleAssignmentResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleAssignmentKind,
			Tenant:     paramsAuthorization.Tenant,
		},
		Status: &secalib.Status{},
		Spec: &secalib.RoleAssignmentSpecV1{
			Subs:  paramsAuthorization.RoleAssignment.InitialSpec.Subs,
			Roles: paramsAuthorization.RoleAssignment.InitialSpec.Roles,
		},
	}
	for _, scope := range paramsAuthorization.RoleAssignment.InitialSpec.Scopes {
		roleAssignmentResponse.Spec.Scopes = append(roleAssignmentResponse.Spec.Scopes, &secalib.RoleAssignmentSpecScopeV1{
			Tenants:    scope.Tenants,
			Regions:    scope.Regions,
			Workspaces: scope.Workspaces,
		})
	}

	// Create a role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodPut
	roleAssignmentResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	roleAssignmentResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	roleAssignmentResponse.Metadata.ResourceVersion = 1
	roleAssignmentResponse.Status.State = secalib.CreatingStatusState
	roleAssignmentResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          roleAssignmentUrl,
		params:       paramsAuthorization,
		response:     roleAssignmentResponse,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "CreateRoleAssignment",
		nextState:    "GetCreatedRoleAssignment",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodGet
	roleAssignmentResponse.Status.State = secalib.ActiveStatusState
	roleAssignmentResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          roleAssignmentUrl,
		params:       paramsAuthorization,
		response:     roleAssignmentResponse,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "GetCreatedRoleAssignment",
		nextState:    "UpdateRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodPut
	roleAssignmentResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	roleAssignmentResponse.Metadata.ResourceVersion = roleAssignmentResponse.Metadata.ResourceVersion + 1
	roleAssignmentResponse.Spec.Subs = paramsAuthorization.RoleAssignment.UpdatedSpec.Subs
	roleAssignmentResponse.Status.State = secalib.UpdatingStatusState
	roleAssignmentResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          roleAssignmentUrl,
		params:       paramsAuthorization,
		response:     roleAssignmentResponse,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "UpdateRoleAssignment",
		nextState:    "GetUpdatedRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get updated role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodGet
	roleAssignmentResponse.Status.State = secalib.ActiveStatusState
	roleAssignmentResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          roleAssignmentUrl,
		params:       paramsAuthorization,
		response:     roleAssignmentResponse,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "GetUpdatedRoleAssignment",
		nextState:    "CreateWorkspace",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	//Workspace
	workspaceResponse := &resourceResponse[secalib.WorkspaceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsWorkspace.Workspace.Name,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     paramsWorkspace.Tenant,
			Region:     paramsWorkspace.Region,
		},
		Status: &secalib.Status{},
	}

	// Create a workspace
	workspaceResponse.Metadata.Verb = http.MethodPut
	workspaceResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	workspaceResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	workspaceResponse.Metadata.ResourceVersion = 1
	workspaceResponse.Status.State = secalib.CreatingStatusState
	workspaceResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          workspaceURL,
		params:       paramsWorkspace,
		response:     workspaceResponse,
		template:     workspaceResponseTemplateV1,
		currentState: "CreateWorkspace",
		nextState:    "GetCreatedWorkspace",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created workspace
	workspaceResponse.Metadata.Verb = http.MethodGet
	workspaceResponse.Status.State = secalib.ActiveStatusState
	workspaceResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          workspaceURL,
		params:       paramsWorkspace,
		response:     workspaceResponse,
		template:     workspaceResponseTemplateV1,
		currentState: "GetCreatedWorkspace",
		nextState:    "UpdateWorkspace",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the workspace
	workspaceResponse.Metadata.Verb = http.MethodPut
	workspaceResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	workspaceResponse.Metadata.ResourceVersion = workspaceResponse.Metadata.ResourceVersion + 1
	workspaceResponse.Status.State = secalib.UpdatingStatusState
	workspaceResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          workspaceURL,
		params:       paramsWorkspace,
		response:     workspaceResponse,
		template:     workspaceResponseTemplateV1,
		currentState: "UpdateWorkspace",
		nextState:    "GetUpdatedWorkspace",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get updated workspace
	workspaceResponse.Metadata.Verb = http.MethodGet
	workspaceResponse.Status.State = secalib.ActiveStatusState
	workspaceResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          workspaceURL,
		params:       paramsWorkspace,
		response:     workspaceResponse,
		template:     workspaceResponseTemplateV1,
		currentState: "GetUpdatedWorkspace",
		nextState:    "CreateBlockStorage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	//Skus

	//Storage
	blockResponse := &resourceResponse[secalib.BlockStorageSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsStorage.BlockStorage.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     paramsStorage.Tenant,
			Workspace:  paramsStorage.Workspace,
			Region:     paramsStorage.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.BlockStorageSpecV1{
			SkuRef: paramsStorage.BlockStorage.InitialSpec.SkuRef,
		},
	}

	// Create a block storage
	blockResponse.Metadata.Verb = http.MethodPut
	blockResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.ResourceVersion = 1
	blockResponse.Spec.SizeGB = paramsStorage.BlockStorage.InitialSpec.SizeGB
	blockResponse.Status.State = secalib.CreatingStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          blockStorageURL,
		params:       paramsStorage,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          blockStorageURL,
		params:       paramsStorage,
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
	blockResponse.Spec.SizeGB = paramsStorage.BlockStorage.UpdatedSpec.SizeGB
	blockResponse.Status.State = secalib.UpdatingStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          blockStorageURL,
		params:       paramsStorage,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          blockStorageURL,
		params:       paramsStorage,
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
			Name:       paramsStorage.Image.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   imageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.ImageKind,
			Tenant:     paramsStorage.Tenant,
			Region:     paramsStorage.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.ImageSpecV1{
			BlockStorageRef: paramsStorage.Image.InitialSpec.BlockStorageRef,
			CpuArchitecture: paramsStorage.Image.InitialSpec.CpuArchitecture,
		},
	}

	// Create an image
	imageResponse.Metadata.Verb = http.MethodPut
	imageResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	imageResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	imageResponse.Metadata.ResourceVersion = 1
	imageResponse.Status.State = secalib.CreatingStatusState
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          imageURL,
		params:       paramsStorage,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          imageURL,
		params:       paramsStorage,
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
	imageResponse.Spec.CpuArchitecture = paramsStorage.Image.UpdatedSpec.CpuArchitecture
	imageResponse.Status.State = secalib.UpdatingStatusState
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          imageURL,
		params:       paramsStorage,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          imageURL,
		params:       paramsStorage,
		response:     imageResponse,
		template:     imageResponseTemplateV1,
		currentState: "GetUpdatedImage",
		nextState:    "CreateNetwork",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}
	//Network
	networkResponse := &resourceResponse[secalib.NetworkSpecV1]{
		Metadata: &secalib.Metadata{
			Name:            paramsNetwork.Network.Name,
			Provider:        secalib.NetworkProviderV1,
			Resource:        networkResource,
			Verb:            http.MethodPut,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      secalib.ApiVersion1,
			Kind:            secalib.NetworkKind,
			Tenant:          paramsNetwork.Tenant,
			Region:          paramsNetwork.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.NetworkSpecV1{
			Cidr: &secalib.NetworkSpecCIDRV1{
				Ipv4: paramsNetwork.Network.InitialSpec.Cidr.Ipv4,
			},
			SkuRef:        paramsNetwork.Network.InitialSpec.SkuRef,
			RouteTableRef: paramsNetwork.Network.InitialSpec.RouteTableRef,
		},
	}

	// Create  Network
	networkResponse.Metadata.Verb = http.MethodPut
	networkResponse.Status.State = secalib.ActiveStatusState
	networkResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	networkResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	networkResponse.Metadata.ResourceVersion = 1
	networkResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          networkURL,
		params:       paramsNetwork,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "CreateNetwork",
		nextState:    "GetNetwork",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get network
	networkResponse.Metadata.Verb = http.MethodGet
	networkResponse.Status.State = secalib.ActiveStatusState
	networkResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          networkURL,
		params:       paramsNetwork,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "GetNetwork",
		nextState:    "UpdateNetwork",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update Network
	networkResponse.Metadata.Verb = http.MethodPut
	networkResponse.Status.State = secalib.UpdatingStatusState
	networkResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	networkResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)

	networkResponse.Spec = &secalib.NetworkSpecV1{
		Cidr: &secalib.NetworkSpecCIDRV1{
			Ipv4: paramsNetwork.Network.UpdatedSpec.Cidr.Ipv4,
		},
		SkuRef:        paramsNetwork.Network.UpdatedSpec.SkuRef,
		RouteTableRef: paramsNetwork.Network.UpdatedSpec.RouteTableRef,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          networkURL,
		params:       paramsNetwork,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "UpdateNetwork",
		nextState:    "GetNetwork2x",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}
	// Get network 2x time
	networkResponse.Metadata.Verb = http.MethodGet
	networkResponse.Status.State = secalib.ActiveStatusState
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          networkURL,
		params:       paramsNetwork,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "GetNetwork2x",
		nextState:    "CreateInternetGateway",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// internet-Gateway
	internetGatewayResponse := &resourceResponse[secalib.InternetGatewaySpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsNetwork.InternetGateway.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   internetGatewayResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsNetwork.Tenant,
			Workspace:  paramsNetwork.Workspace,
			Region:     paramsNetwork.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.InternetGatewaySpecV1{
			EgressOnly: paramsNetwork.InternetGateway.InitialSpec.EgressOnly,
		},
	}

	// Create internet-Gateway
	internetGatewayResponse.Metadata.Verb = http.MethodPut
	internetGatewayResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.ResourceVersion = 1
	internetGatewayResponse.Status.State = secalib.CreatingStatusState
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          internetGatewayURL,
		params:       paramsNetwork,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "CreateInternetGateway",
		nextState:    "GetInternetGateway",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get internet-Gateway
	internetGatewayResponse.Metadata.Verb = http.MethodGet
	internetGatewayResponse.Status.State = secalib.ActiveStatusState
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          internetGatewayURL,
		params:       paramsNetwork,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "GetInternetGateway",
		nextState:    "UpdateInternetGateway",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update internet-gateway
	internetGatewayResponse.Metadata.Verb = http.MethodPut
	internetGatewayResponse.Status.State = secalib.UpdatingStatusState
	internetGatewayResponse.Metadata.ResourceVersion = internetGatewayResponse.Metadata.ResourceVersion + 1
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Spec = &secalib.InternetGatewaySpecV1{
		EgressOnly: paramsNetwork.InternetGateway.UpdatedSpec.EgressOnly,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          internetGatewayURL,
		params:       paramsNetwork,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "UpdateInternetGateway",
		nextState:    "GetInternetGateway2x",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get internet-gateway after update
	internetGatewayResponse.Metadata.Verb = http.MethodGet
	internetGatewayResponse.Status.State = secalib.ActiveStatusState
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          internetGatewayURL,
		params:       paramsNetwork,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "GetInternetGateway2x",
		nextState:    "CreateRouteTable",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Route-Table
	routeTableResponse := &resourceResponse[secalib.RouteTableSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsNetwork.RouteTable.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   routeTableResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsNetwork.Tenant,
			Workspace:  paramsNetwork.Workspace,
			Region:     paramsNetwork.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.RouteTableSpecV1{
			LocalRef: paramsNetwork.RouteTable.InitialSpec.LocalRef,
		},
	}

	for _, routes := range paramsNetwork.RouteTable.InitialSpec.Routes {
		routeTableResponse.Spec.Routes = append(routeTableResponse.Spec.Routes, &secalib.RouteTableRouteV1{
			DestinationCidrBlock: routes.DestinationCidrBlock,
			TargetRef:            routes.TargetRef,
		})
	}

	// Create route-Table
	routeTableResponse.Metadata.Verb = http.MethodPut
	routeTableResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	routeTableResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	routeTableResponse.Metadata.ResourceVersion = 1
	routeTableResponse.Status.State = secalib.CreatingStatusState
	routeTableResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          routeTableURL,
		params:       paramsNetwork,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "CreateRouteTable",
		nextState:    "GetRouteTable",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get route-Table
	routeTableResponse.Metadata.Verb = http.MethodGet
	routeTableResponse.Status.State = secalib.ActiveStatusState
	routeTableResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          routeTableURL,
		params:       paramsNetwork,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "GetRouteTable",
		nextState:    "UpdateRouteTable",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update route-Table
	routeTableResponse.Metadata.Verb = http.MethodPut
	routeTableResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	routeTableResponse.Metadata.ResourceVersion = routeTableResponse.Metadata.ResourceVersion + 1
	routeTableResponse.Status.State = secalib.UpdatingStatusState
	routeTableResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	routeTableResponse.Spec = &secalib.RouteTableSpecV1{
		LocalRef: paramsNetwork.RouteTable.UpdatedSpec.LocalRef,
		Routes:   paramsNetwork.RouteTable.UpdatedSpec.Routes,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          routeTableURL,
		params:       paramsNetwork,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "UpdateRouteTable",
		nextState:    "GetRouteTableUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get route-table after update
	routeTableResponse.Metadata.Verb = http.MethodGet
	routeTableResponse.Status.State = secalib.ActiveStatusState
	routeTableResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          routeTableURL,
		params:       paramsNetwork,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "GetRouteTableUpdated",
		nextState:    "CreateSubnet",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// subnet
	subnetResponse := &resourceResponse[secalib.SubnetSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsNetwork.Subnet.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   subnetResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsNetwork.Tenant,
			Workspace:  paramsNetwork.Workspace,
			Region:     paramsNetwork.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.SubnetSpecV1{
			Cidr: paramsNetwork.Subnet.InitialSpec.Cidr,
			Zone: paramsNetwork.Subnet.InitialSpec.Zone,
		},
	}

	// Create subnet
	subnetResponse.Metadata.Verb = http.MethodPut
	subnetResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	subnetResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	subnetResponse.Metadata.ResourceVersion = 1
	subnetResponse.Status.State = secalib.CreatingStatusState
	subnetResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          subnetURL,
		params:       paramsNetwork,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "CreateSubnet",
		nextState:    "GetSubnet",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get subnet
	subnetResponse.Metadata.Verb = http.MethodGet
	subnetResponse.Status.State = secalib.ActiveStatusState
	subnetResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          subnetURL,
		params:       paramsNetwork,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "GetSubnet",
		nextState:    "UpdateSubnet",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update subnet
	subnetResponse.Metadata.Verb = http.MethodPut
	subnetResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	subnetResponse.Metadata.ResourceVersion = subnetResponse.Metadata.ResourceVersion + 1
	subnetResponse.Status.State = secalib.UpdatingStatusState
	subnetResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	subnetResponse.Spec = &secalib.SubnetSpecV1{
		Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: paramsNetwork.Subnet.UpdatedSpec.Cidr.Ipv4},
		Zone: paramsNetwork.Subnet.UpdatedSpec.Zone,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          subnetURL,
		params:       paramsNetwork,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "UpdateSubnet",
		nextState:    "GetSubnetUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get subnet after update
	subnetResponse.Metadata.Verb = http.MethodGet
	subnetResponse.Status.State = secalib.ActiveStatusState
	subnetResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          subnetURL,
		params:       paramsNetwork,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "GetSubnetUpdated",
		nextState:    "CreatePublicIP",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Public-IP
	publicIPResponse := &resourceResponse[secalib.PublicIPSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsNetwork.PublicIP.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   publicIPResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsNetwork.Tenant,
			Workspace:  paramsNetwork.Workspace,
			Region:     paramsNetwork.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.PublicIPSpecV1{
			Version: paramsNetwork.PublicIP.InitialSpec.Version,
			Address: paramsNetwork.PublicIP.InitialSpec.Address,
		},
	}

	// Create public-IP
	publicIPResponse.Metadata.Verb = http.MethodPut
	publicIPResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.ResourceVersion = 1
	publicIPResponse.Status.State = secalib.CreatingStatusState
	publicIPResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)

	if err := configurePutStub(wm, scenario, stubConfig{
		url:          publicIPURL,
		params:       paramsNetwork,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "CreatePublicIP",
		nextState:    "GetPublicIP",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get public-IP
	publicIPResponse.Metadata.Verb = http.MethodGet
	publicIPResponse.Status.State = secalib.ActiveStatusState
	publicIPResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          publicIPURL,
		params:       paramsNetwork,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "GetPublicIP",
		nextState:    "UpdatePublicIP",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update public-IP
	publicIPResponse.Metadata.Verb = http.MethodPut
	publicIPResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.ResourceVersion = publicIPResponse.Metadata.ResourceVersion + 1
	publicIPResponse.Status.State = secalib.UpdatingStatusState
	publicIPResponse.Spec = &secalib.PublicIPSpecV1{
		Version: paramsNetwork.PublicIP.UpdatedSpec.Version,
		Address: paramsNetwork.PublicIP.UpdatedSpec.Address,
	}

	if err := configurePutStub(wm, scenario, stubConfig{
		url:          publicIPURL,
		params:       paramsNetwork,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "UpdatePublicIP",
		nextState:    "GetPublicIPUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get public-IP after update
	publicIPResponse.Metadata.Verb = http.MethodGet
	publicIPResponse.Status.State = secalib.ActiveStatusState
	publicIPResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          publicIPURL,
		params:       paramsNetwork,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "GetPublicIPUpdated",
		nextState:    "CreateNIC",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// NIC
	nicResponse := &resourceResponse[secalib.NICSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsNetwork.NIC.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   nicResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsNetwork.Tenant,
			Workspace:  paramsNetwork.Workspace,
			Region:     paramsNetwork.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.NICSpecV1{
			Addresses: paramsNetwork.NIC.InitialSpec.Addresses,
			SubnetRef: paramsNetwork.NIC.InitialSpec.SubnetRef,
		},
	}

	// Create NIC
	nicResponse.Metadata.Verb = http.MethodPut
	nicResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.ResourceVersion = 1
	nicResponse.Status.State = secalib.CreatingStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          nicURL,
		params:       paramsNetwork,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "CreateNIC",
		nextState:    "GetNIC",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get NIC
	nicResponse.Metadata.Verb = http.MethodGet
	nicResponse.Status.State = secalib.ActiveStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          nicURL,
		params:       paramsNetwork,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "GetNIC",
		nextState:    "UpdateNIC",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update NIC
	nicResponse.Metadata.Verb = http.MethodPut
	nicResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.ResourceVersion = nicResponse.Metadata.ResourceVersion + 1
	nicResponse.Status.State = secalib.UpdatingStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	nicResponse.Spec = &secalib.NICSpecV1{
		Addresses: paramsNetwork.NIC.UpdatedSpec.Addresses,
		SubnetRef: paramsNetwork.NIC.UpdatedSpec.SubnetRef,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          nicURL,
		params:       paramsNetwork,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "UpdateNIC",
		nextState:    "GetNICUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get NIC after update
	nicResponse.Metadata.Verb = http.MethodGet
	nicResponse.Status.State = secalib.ActiveStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          nicURL,
		params:       paramsNetwork,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "GetNICUpdated",
		nextState:    "CreateSecurityGroup",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Security-group
	securityGroupResponse := &resourceResponse[secalib.SecurityGroupSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsNetwork.SecurityGroup.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   securityGroupResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsNetwork.Tenant,
			Workspace:  paramsNetwork.Workspace,
			Region:     paramsNetwork.Region,
		},
		Status: &secalib.Status{},
		Spec:   &secalib.SecurityGroupSpecV1{},
	}

	for _, rules := range paramsNetwork.SecurityGroup.InitialSpec.Rules {
		securityGroupResponse.Spec.Rules = append(securityGroupResponse.Spec.Rules, &secalib.SecurityGroupRule{
			Direction: rules.Direction,
		})
	}
	// Create Security-group
	securityGroupResponse.Metadata.Verb = http.MethodPut
	securityGroupResponse.Status.State = secalib.CreatingStatusState
	securityGroupResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	securityGroupResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	securityGroupResponse.Metadata.ResourceVersion = 1
	securityGroupResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          securityGroupURL,
		params:       paramsNetwork,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "CreateSecurityGroup",
		nextState:    "GetSecurityGroup",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get Security-group
	securityGroupResponse.Metadata.Verb = http.MethodGet
	securityGroupResponse.Status.State = secalib.ActiveStatusState
	securityGroupResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          securityGroupURL,
		params:       paramsNetwork,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "GetSecurityGroup",
		nextState:    "UpdateSecurityGroup",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update Security-group
	securityGroupResponse.Metadata.Verb = http.MethodPut
	securityGroupResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	securityGroupResponse.Metadata.ResourceVersion = securityGroupResponse.Metadata.ResourceVersion + 1
	securityGroupResponse.Status.State = secalib.UpdatingStatusState
	securityGroupResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	securityGroupResponse.Spec = &secalib.SecurityGroupSpecV1{
		Rules: paramsNetwork.SecurityGroup.UpdatedSpec.Rules,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          securityGroupURL,
		params:       paramsNetwork,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "UpdateSecurityGroup",
		nextState:    "GetSecurityGroupUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get Security-group after update
	securityGroupResponse.Metadata.Verb = http.MethodGet
	securityGroupResponse.Status.State = secalib.ActiveStatusState
	securityGroupResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          securityGroupURL,
		params:       paramsNetwork,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "GetSecurityGroupUpdated",
		nextState:    "CreateBlockStorage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}
	//Compute
	instResponse := &resourceResponse[secalib.InstanceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsCompute.Instance.Name,
			Provider:   secalib.ComputeProviderV1,
			Resource:   instanceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsCompute.Tenant,
			Workspace:  paramsCompute.Workspace,
			Region:     paramsCompute.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.InstanceSpecV1{
			SkuRef:        paramsCompute.Instance.InitialSpec.SkuRef,
			Zone:          paramsCompute.Instance.InitialSpec.Zone,
			BootDeviceRef: paramsCompute.Instance.InitialSpec.BootDeviceRef,
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
		url:          instanceURL,
		params:       paramsCompute,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          instanceURL,
		params:       paramsCompute,
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
	instResponse.Spec.Zone = paramsCompute.Instance.UpdatedSpec.Zone
	instResponse.Status.State = secalib.UpdatingStatusState
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          instanceURL,
		params:       paramsCompute,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          instanceURL,
		params:       paramsCompute,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: "GetUpdatedInstance",
		nextState:    "DeleteInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}
	//Access

	//Delete all
	// Delete Instance
	instResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          instanceURL,
		params:       paramsCompute,
		response:     instResponse,
		currentState: "DeleteInstance",
		nextState:    "DeleteSecurityGroup",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Security Group
	securityGroupResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          securityGroupURL,
		params:       paramsNetwork,
		response:     securityGroupResponse,
		currentState: "DeleteSecurityGroup",
		nextState:    "DeleteNic",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Nic
	nicResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          nicURL,
		params:       paramsNetwork,
		response:     nicResponse,
		currentState: "DeleteNic",
		nextState:    "DeletePublicIP",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete public ip
	publicIPResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          publicIPURL,
		params:       paramsNetwork,
		response:     publicIPResponse,
		currentState: "DeletePublicIP",
		nextState:    "DeleteSubnet",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete subnet
	subnetResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          subnetURL,
		params:       paramsNetwork,
		response:     subnetResponse,
		currentState: "DeleteSubnet",
		nextState:    "DeleteRouteTable",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Route-table
	routeTableResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          routeTableURL,
		params:       paramsNetwork,
		response:     routeTableResponse,
		currentState: "DeleteRouteTable",
		nextState:    "DeleteInternetGateway",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Internet-gateway
	internetGatewayResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          internetGatewayURL,
		params:       paramsNetwork,
		response:     internetGatewayResponse,
		currentState: "DeleteInternetGateway",
		nextState:    "DeleteNetwork",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Network
	networkResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          networkURL,
		params:       paramsNetwork,
		response:     networkResponse,
		currentState: "DeleteNetwork",
		nextState:    "DeleteBlockStorage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	blockResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          blockStorageURL,
		params:       paramsStorage,
		response:     blockResponse,
		currentState: "DeleteBlockStorage",
		nextState:    "DeleteImage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}
	imageResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          imageURL,
		params:       paramsStorage,
		response:     imageResponse,
		currentState: "DeleteImage",
		nextState:    "DeleteWorkspace",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}
	workspaceResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          workspaceURL,
		params:       paramsWorkspace,
		response:     workspaceResponse,
		currentState: "DeleteWorkspace",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	return wm, err
}
