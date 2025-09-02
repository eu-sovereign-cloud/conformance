package mock

import (
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/wiremock/go-wiremock"
)

func TestUsageScenario(scenario string, paramsUsage UsageParamsV1) (*wiremock.Client, error) {

	wm, err := newClient(paramsUsage.Workspace.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs

	// Authorization
	roleUrl := secalib.GenerateRoleURL(paramsUsage.Authorization.Tenant, paramsUsage.Authorization.Role.Name)
	roleAssignmentUrl := secalib.GenerateRoleAssignmentURL(paramsUsage.Authorization.Tenant, paramsUsage.Authorization.RoleAssignment.Name)

	//workspace
	workspaceURL := secalib.GenerateWorkspaceURL(paramsUsage.Workspace.Tenant, paramsUsage.Workspace.Workspace.Name)

	//Storage
	blockStorageURL := secalib.GenerateBlockStorageURL(paramsUsage.Storage.Tenant, paramsUsage.Storage.Workspace, paramsUsage.Storage.BlockStorage.Name)
	imageURL := secalib.GenerateImageURL(paramsUsage.Storage.Tenant, paramsUsage.Storage.Image.Name)
	//Compute
	instanceURL := secalib.GenerateInstanceURL(paramsUsage.Compute.Tenant, paramsUsage.Compute.Workspace, paramsUsage.Compute.Instance.Name)

	//Network
	networkURL := secalib.GenerateNetworkURL(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.Network.Name)
	internetGatewayURL := secalib.GenerateInternetGatewayURL(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.InternetGateway.Name)
	nicURL := secalib.GenerateNicURL(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.NIC.Name)
	publicIPURL := secalib.GeneratePublicIPURL(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.PublicIP.Name)
	routeTableURL := secalib.GenerateRouteTableURL(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.RouteTable.Name)
	subnetURL := secalib.GenerateSubnetURL(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.Subnet.Name)
	securityGroupURL := secalib.GenerateSecurityGroupURL(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.SecurityGroup.Name)

	// GenerateResources
	// Authorization
	roleResource := secalib.GenerateRoleResource(paramsUsage.Authorization.Tenant, paramsUsage.Authorization.Role.Name)
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(paramsUsage.Authorization.Tenant, paramsUsage.Authorization.RoleAssignment.Name)

	//Workspace

	workspaceResource := secalib.GenerateWorkspaceResource(paramsUsage.Workspace.Tenant, paramsUsage.Workspace.Workspace.Name)

	//Storage
	blockStorageResource := secalib.GenerateBlockStorageResource(paramsUsage.Storage.Tenant, paramsUsage.Storage.Workspace, paramsUsage.Storage.BlockStorage.Name)
	imageResource := secalib.GenerateImageResource(paramsUsage.Storage.Tenant, paramsUsage.Storage.Image.Name)
	// Compute
	instanceResource := secalib.GenerateInstanceResource(paramsUsage.Compute.Tenant, paramsUsage.Compute.Workspace, paramsUsage.Compute.Instance.Name)

	//Network
	networkResource := secalib.GenerateNetworkResource(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.Network.Name)
	internetGatewayResource := secalib.GenerateInternetGatewayResource(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.InternetGateway.Name)
	nicResource := secalib.GenerateNicResource(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.NIC.Name)
	publicIPResource := secalib.GeneratePublicIPResource(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.PublicIP.Name)
	routeTableResource := secalib.GenerateRouteTableResource(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.RouteTable.Name)
	subnetResource := secalib.GenerateSubnetResource(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.Subnet.Name)
	securityGroupResource := secalib.GenerateSecurityGroupResource(paramsUsage.Network.Tenant, paramsUsage.Network.Workspace, paramsUsage.Network.SecurityGroup.Name)

	//Authorization
	roleResponse := &resourceResponse[secalib.RoleSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsUsage.Authorization.Role.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleKind,
			Tenant:     paramsUsage.Authorization.Tenant,
		},
		Status: &secalib.Status{},
		Spec:   &secalib.RoleSpecV1{},
	}
	for _, perm := range paramsUsage.Authorization.Role.InitialSpec.Permissions {
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
		params:       paramsUsage.Authorization,
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
		params:       paramsUsage.Authorization,
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
	for i, perm := range paramsUsage.Authorization.Role.UpdatedSpec.Permissions {
		roleResponse.Spec.Permissions[i].Verb = append([]string{}, perm.Verb...)
	}
	roleResponse.Status.State = secalib.UpdatingStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          roleUrl,
		params:       paramsUsage.Authorization,
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
		params:       paramsUsage.Authorization,
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
			Name:       paramsUsage.Authorization.RoleAssignment.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleAssignmentResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleAssignmentKind,
			Tenant:     paramsUsage.Authorization.Tenant,
		},
		Status: &secalib.Status{},
		Spec: &secalib.RoleAssignmentSpecV1{
			Subs:  paramsUsage.Authorization.RoleAssignment.InitialSpec.Subs,
			Roles: paramsUsage.Authorization.RoleAssignment.InitialSpec.Roles,
		},
	}
	for _, scope := range paramsUsage.Authorization.RoleAssignment.InitialSpec.Scopes {
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
		params:       paramsUsage.Authorization,
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
		params:       paramsUsage.Authorization,
		response:     roleAssignmentResponse,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "GetCreatedRoleAssignment",
		nextState:    "CreateWorkspace",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	//Workspace
	workspaceResponse := &resourceResponse[secalib.WorkspaceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsUsage.Workspace.Workspace.Name,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     paramsUsage.Workspace.Tenant,
			Region:     paramsUsage.Workspace.Region,
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
		params:       paramsUsage.Workspace,
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
		params:       paramsUsage.Workspace,
		response:     workspaceResponse,
		template:     workspaceResponseTemplateV1,
		currentState: "GetCreatedWorkspace",
		nextState:    "CreateBlockStorage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	//Skus

	//Storage
	blockResponse := &resourceResponse[secalib.BlockStorageSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsUsage.Storage.BlockStorage.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     paramsUsage.Storage.Tenant,
			Workspace:  paramsUsage.Storage.Workspace,
			Region:     paramsUsage.Storage.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.BlockStorageSpecV1{
			SkuRef: paramsUsage.Storage.BlockStorage.InitialSpec.SkuRef,
		},
	}

	// Create a block storage
	blockResponse.Metadata.Verb = http.MethodPut
	blockResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.ResourceVersion = 1
	blockResponse.Spec.SizeGB = paramsUsage.Storage.BlockStorage.InitialSpec.SizeGB
	blockResponse.Status.State = secalib.CreatingStatusState
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          blockStorageURL,
		params:       paramsUsage.Storage,
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
		params:       paramsUsage.Storage,
		response:     blockResponse,
		template:     blockStorageResponseTemplateV1,
		currentState: "GetCreatedBlockStorage",
		nextState:    "UpdateBlockStorage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// image
	imageResponse := &resourceResponse[secalib.ImageSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsUsage.Storage.Image.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   imageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.ImageKind,
			Tenant:     paramsUsage.Storage.Tenant,
			Region:     paramsUsage.Storage.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.ImageSpecV1{
			BlockStorageRef: paramsUsage.Storage.Image.InitialSpec.BlockStorageRef,
			CpuArchitecture: paramsUsage.Storage.Image.InitialSpec.CpuArchitecture,
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
		params:       paramsUsage.Storage,
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
		params:       paramsUsage.Storage,
		response:     imageResponse,
		template:     imageResponseTemplateV1,
		currentState: "GetCreatedImage",
		nextState:    "CreateNetwork",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	//Network
	networkResponse := &resourceResponse[secalib.NetworkSpecV1]{
		Metadata: &secalib.Metadata{
			Name:            paramsUsage.Network.Network.Name,
			Provider:        secalib.NetworkProviderV1,
			Resource:        networkResource,
			Verb:            http.MethodPut,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      secalib.ApiVersion1,
			Kind:            secalib.NetworkKind,
			Tenant:          paramsUsage.Network.Tenant,
			Region:          paramsUsage.Network.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.NetworkSpecV1{
			Cidr: &secalib.NetworkSpecCIDRV1{
				Ipv4: paramsUsage.Network.Network.InitialSpec.Cidr.Ipv4,
			},
			SkuRef:        paramsUsage.Network.Network.InitialSpec.SkuRef,
			RouteTableRef: paramsUsage.Network.Network.InitialSpec.RouteTableRef,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "GetNetwork",
		nextState:    "CreateInternetGateway",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// internet-Gateway
	internetGatewayResponse := &resourceResponse[secalib.InternetGatewaySpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsUsage.Network.InternetGateway.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   internetGatewayResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsUsage.Network.Tenant,
			Workspace:  paramsUsage.Network.Workspace,
			Region:     paramsUsage.Network.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.InternetGatewaySpecV1{
			EgressOnly: paramsUsage.Network.InternetGateway.InitialSpec.EgressOnly,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "GetInternetGateway",
		nextState:    "CreateRouteTable",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Route-Table
	routeTableResponse := &resourceResponse[secalib.RouteTableSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsUsage.Network.RouteTable.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   routeTableResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsUsage.Network.Tenant,
			Workspace:  paramsUsage.Network.Workspace,
			Region:     paramsUsage.Network.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.RouteTableSpecV1{
			LocalRef: paramsUsage.Network.RouteTable.InitialSpec.LocalRef,
		},
	}

	for _, routes := range paramsUsage.Network.RouteTable.InitialSpec.Routes {
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "GetRouteTable",
		nextState:    "CreateSubnet",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// subnet
	subnetResponse := &resourceResponse[secalib.SubnetSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsUsage.Network.Subnet.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   subnetResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsUsage.Network.Tenant,
			Workspace:  paramsUsage.Network.Workspace,
			Region:     paramsUsage.Network.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.SubnetSpecV1{
			Cidr: paramsUsage.Network.Subnet.InitialSpec.Cidr,
			Zone: paramsUsage.Network.Subnet.InitialSpec.Zone,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "GetSubnet",
		nextState:    "UpdateSubnet",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Public-IP
	publicIPResponse := &resourceResponse[secalib.PublicIPSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsUsage.Network.PublicIP.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   publicIPResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsUsage.Network.Tenant,
			Workspace:  paramsUsage.Network.Workspace,
			Region:     paramsUsage.Network.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.PublicIPSpecV1{
			Version: paramsUsage.Network.PublicIP.InitialSpec.Version,
			Address: paramsUsage.Network.PublicIP.InitialSpec.Address,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "GetPublicIP",
		nextState:    "UpdatePublicIP",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// NIC
	nicResponse := &resourceResponse[secalib.NICSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsUsage.Network.NIC.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   nicResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsUsage.Network.Tenant,
			Workspace:  paramsUsage.Network.Workspace,
			Region:     paramsUsage.Network.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.NICSpecV1{
			Addresses: paramsUsage.Network.NIC.InitialSpec.Addresses,
			SubnetRef: paramsUsage.Network.NIC.InitialSpec.SubnetRef,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "GetNIC",
		nextState:    "CreateSecurityGroup",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Security-group
	securityGroupResponse := &resourceResponse[secalib.SecurityGroupSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsUsage.Network.SecurityGroup.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   securityGroupResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsUsage.Network.Tenant,
			Workspace:  paramsUsage.Network.Workspace,
			Region:     paramsUsage.Network.Region,
		},
		Status: &secalib.Status{},
		Spec:   &secalib.SecurityGroupSpecV1{},
	}

	for _, rules := range paramsUsage.Network.SecurityGroup.InitialSpec.Rules {
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "GetSecurityGroup",
		nextState:    "CreateInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	//Compute
	instResponse := &resourceResponse[secalib.InstanceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       paramsUsage.Compute.Instance.Name,
			Provider:   secalib.ComputeProviderV1,
			Resource:   instanceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     paramsUsage.Compute.Tenant,
			Workspace:  paramsUsage.Compute.Workspace,
			Region:     paramsUsage.Compute.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.InstanceSpecV1{
			SkuRef:        paramsUsage.Compute.Instance.InitialSpec.SkuRef,
			Zone:          paramsUsage.Compute.Instance.InitialSpec.Zone,
			BootDeviceRef: paramsUsage.Compute.Instance.InitialSpec.BootDeviceRef,
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
		params:       paramsUsage.Compute,
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
		params:       paramsUsage.Compute,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: "GetCreatedInstance",
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
		params:       paramsUsage.Compute,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Network,
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
		params:       paramsUsage.Storage,
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
		params:       paramsUsage.Storage,
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
		params:       paramsUsage.Workspace,
		response:     workspaceResponse,
		currentState: "DeleteWorkspace",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	return wm, err
}
