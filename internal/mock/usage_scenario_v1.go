package mock

import (
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/wiremock/go-wiremock"
)

type UsageScenariosV1 struct {
	Scenarios
}

func NewUsageScenariosV1(authToken string, tenant string, region string, mockURL string) *UsageScenariosV1 {
	return &UsageScenariosV1{
		Scenarios: Scenarios{
			params: secalib.GeneralParams{
				AuthToken: authToken,
				Tenant:    tenant,
				Region:    region,
			},
			mockURL: mockURL,
		},
	}
}

func (scenarios *UsageScenariosV1) ConfigureUsageScenario(id string, params *secalib.UsageParamsV1) (*wiremock.Client, error) {
	wm, err := scenarios.newClient()
	if err != nil {
		return nil, err
	}

	name := "UsageV1_" + id

	// Generate URLs

	// Authorization
	roleUrl := secalib.GenerateRoleURL(scenarios.params.Tenant, params.Role.Name)
	roleAssignmentUrl := secalib.GenerateRoleAssignmentURL(scenarios.params.Tenant, params.RoleAssignment.Name)

	// workspace
	workspaceURL := secalib.GenerateWorkspaceURL(scenarios.params.Tenant, params.Workspace.Name)

	// Storage
	blockStorageURL := secalib.GenerateBlockStorageURL(scenarios.params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageURL := secalib.GenerateImageURL(scenarios.params.Tenant, params.Image.Name)
	// Compute
	instanceURL := secalib.GenerateInstanceURL(scenarios.params.Tenant, params.Workspace.Name, params.Instance.Name)

	// Network
	networkURL := secalib.GenerateNetworkURL(scenarios.params.Tenant, params.Workspace.Name, params.Network.Name)
	internetGatewayURL := secalib.GenerateInternetGatewayURL(scenarios.params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicURL := secalib.GenerateNicURL(scenarios.params.Tenant, params.Workspace.Name, params.Nic.Name)
	publicIPURL := secalib.GeneratePublicIPURL(scenarios.params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeTableURL := secalib.GenerateRouteTableURL(scenarios.params.Tenant, params.Workspace.Name, params.RouteTable.Name)
	subnetURL := secalib.GenerateSubnetURL(scenarios.params.Tenant, params.Workspace.Name, params.Subnet.Name)
	securityGroupURL := secalib.GenerateSecurityGroupURL(scenarios.params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// GenerateResources
	// Authorization
	roleResource := secalib.GenerateRoleResource(scenarios.params.Tenant, params.Role.Name)
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(scenarios.params.Tenant, params.RoleAssignment.Name)

	// Workspace

	workspaceResource := secalib.GenerateWorkspaceResource(scenarios.params.Tenant, params.Workspace.Name)

	// Storage
	blockStorageResource := secalib.GenerateBlockStorageResource(scenarios.params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageResource := secalib.GenerateImageResource(scenarios.params.Tenant, params.Image.Name)
	// Compute
	instanceResource := secalib.GenerateInstanceResource(scenarios.params.Tenant, params.Workspace.Name, params.Instance.Name)

	// Network
	networkResource := secalib.GenerateNetworkResource(scenarios.params.Tenant, params.Workspace.Name, params.Network.Name)
	internetGatewayResource := secalib.GenerateInternetGatewayResource(scenarios.params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicResource := secalib.GenerateNicResource(scenarios.params.Tenant, params.Workspace.Name, params.Nic.Name)
	publicIPResource := secalib.GeneratePublicIPResource(scenarios.params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeTableResource := secalib.GenerateRouteTableResource(scenarios.params.Tenant, params.Workspace.Name, params.RouteTable.Name)
	subnetResource := secalib.GenerateSubnetResource(scenarios.params.Tenant, params.Workspace.Name, params.Subnet.Name)
	securityGroupResource := secalib.GenerateSecurityGroupResource(scenarios.params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Authorization
	roleResponse := &resourceResponse[secalib.RoleSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Role.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleKind,
			Tenant:     scenarios.params.Tenant,
		},
		Status: &secalib.Status{},
		Spec:   &secalib.RoleSpecV1{},
	}
	for _, perm := range params.Role.InitialSpec.Permissions {
		roleResponse.Spec.Permissions = append(roleResponse.Spec.Permissions, &secalib.RoleSpecPermissionV1{
			Provider:  perm.Provider,
			Resources: append([]string{}, perm.Resources...),
			Verb:      append([]string{}, perm.Verb...),
		})
	}

	// Create Role
	roleResponse.Metadata.Verb = http.MethodPut
	roleResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	roleResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	roleResponse.Metadata.ResourceVersion = 1
	roleResponse.Status.State = secalib.CreatingStatusState
	roleResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          roleUrl,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          roleUrl,
		params:       scenarios.params,
		response:     roleResponse,
		template:     roleResponseTemplateV1,
		currentState: "GetCreatedRole",
		nextState:    "CreateRoleAssignment",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignmentResponse := &resourceResponse[secalib.RoleAssignmentSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.RoleAssignment.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleAssignmentResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleAssignmentKind,
			Tenant:     scenarios.params.Tenant,
		},
		Status: &secalib.Status{},
		Spec: &secalib.RoleAssignmentSpecV1{
			Subs:  params.RoleAssignment.InitialSpec.Subs,
			Roles: params.RoleAssignment.InitialSpec.Roles,
		},
	}
	for _, scope := range params.RoleAssignment.InitialSpec.Scopes {
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
	if err := configurePutStub(wm, name, stubConfig{
		url:          roleAssignmentUrl,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          roleAssignmentUrl,
		params:       scenarios.params,
		response:     roleAssignmentResponse,
		template:     roleAssignmentResponseTemplateV1,
		currentState: "GetCreatedRoleAssignment",
		nextState:    "CreateWorkspace",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Workspace
	workspaceResponse := &resourceResponse[secalib.WorkspaceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Workspace.Name,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     scenarios.params.Tenant,
			Region:     scenarios.params.Region,
		},
		Labels: &[]secalib.Label{},
		Status: &secalib.Status{},
	}

	// Create a workspace
	workspaceResponse.Metadata.Verb = http.MethodPut
	workspaceResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	workspaceResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	workspaceResponse.Metadata.ResourceVersion = 1
	workspaceResponse.Status.State = secalib.CreatingStatusState
	workspaceResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          workspaceURL,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          workspaceURL,
		params:       scenarios.params,
		response:     workspaceResponse,
		template:     workspaceResponseTemplateV1,
		currentState: "GetCreatedWorkspace",
		nextState:    "CreateImage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Skus

	// Storage

	// image
	imageResponse := &resourceResponse[secalib.ImageSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Image.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   imageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.ImageKind,
			Tenant:     scenarios.params.Tenant,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.ImageSpecV1{
			BlockStorageRef: params.Image.InitialSpec.BlockStorageRef,
			CpuArchitecture: params.Image.InitialSpec.CpuArchitecture,
		},
	}

	// Create an image
	imageResponse.Metadata.Verb = http.MethodPut
	imageResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	imageResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	imageResponse.Metadata.ResourceVersion = 1
	imageResponse.Status.State = secalib.CreatingStatusState
	imageResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          imageURL,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          imageURL,
		params:       scenarios.params,
		response:     imageResponse,
		template:     imageResponseTemplateV1,
		currentState: "GetCreatedImage",
		nextState:    "CreateBlockStorage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	blockResponse := &resourceResponse[secalib.BlockStorageSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.BlockStorage.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     scenarios.params.Tenant,
			Workspace:  params.Workspace.Name,
			Region:     scenarios.params.Region,
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
	if err := configurePutStub(wm, name, stubConfig{
		url:          blockStorageURL,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          blockStorageURL,
		params:       scenarios.params,
		response:     blockResponse,
		template:     blockStorageResponseTemplateV1,
		currentState: "GetCreatedBlockStorage",
		nextState:    "CreateNetwork",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Network
	networkResponse := &resourceResponse[secalib.NetworkSpecV1]{
		Metadata: &secalib.Metadata{
			Name:            params.Network.Name,
			Provider:        secalib.NetworkProviderV1,
			Resource:        networkResource,
			Verb:            http.MethodPut,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      secalib.ApiVersion1,
			Kind:            secalib.NetworkKind,
			Tenant:          scenarios.params.Tenant,
			Region:          scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.NetworkSpecV1{
			Cidr: &secalib.NetworkSpecCIDRV1{
				Ipv4: params.Network.InitialSpec.Cidr.Ipv4,
			},
			SkuRef:        params.Network.InitialSpec.SkuRef,
			RouteTableRef: params.Network.InitialSpec.RouteTableRef,
		},
	}

	// Create  Network
	networkResponse.Metadata.Verb = http.MethodPut
	networkResponse.Status.State = secalib.CreatingStatusState
	networkResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	networkResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	networkResponse.Metadata.ResourceVersion = 1
	networkResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          networkURL,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          networkURL,
		params:       scenarios.params,
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
			Name:       params.InternetGateway.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   internetGatewayResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InternetGatewayKind,
			Tenant:     scenarios.params.Tenant,
			Workspace:  params.Workspace.Name,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.InternetGatewaySpecV1{
			EgressOnly: params.InternetGateway.InitialSpec.EgressOnly,
		},
	}

	// Create internet-Gateway
	internetGatewayResponse.Metadata.Verb = http.MethodPut
	internetGatewayResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.ResourceVersion = 1
	internetGatewayResponse.Status.State = secalib.CreatingStatusState
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          internetGatewayURL,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          internetGatewayURL,
		params:       scenarios.params,
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
			Name:       params.RouteTable.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   routeTableResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RouteTableKind,
			Tenant:     scenarios.params.Tenant,
			Workspace:  params.Workspace.Name,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.RouteTableSpecV1{
			LocalRef: params.RouteTable.InitialSpec.LocalRef,
		},
	}

	for _, routes := range params.RouteTable.InitialSpec.Routes {
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
	if err := configurePutStub(wm, name, stubConfig{
		url:          routeTableURL,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          routeTableURL,
		params:       scenarios.params,
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
			Name:       params.Subnet.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   subnetResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SubnetKind,
			Tenant:     scenarios.params.Tenant,
			Workspace:  params.Workspace.Name,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.SubnetSpecV1{
			Cidr: params.Subnet.InitialSpec.Cidr,
			Zone: params.Subnet.InitialSpec.Zone,
		},
	}

	// Create subnet
	subnetResponse.Metadata.Verb = http.MethodPut
	subnetResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	subnetResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	subnetResponse.Metadata.ResourceVersion = 1
	subnetResponse.Status.State = secalib.CreatingStatusState
	subnetResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          subnetURL,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          subnetURL,
		params:       scenarios.params,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "GetSubnet",
		nextState:    "CreateSecurityGroup",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Security-group
	securityGroupResponse := &resourceResponse[secalib.SecurityGroupSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.SecurityGroup.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   securityGroupResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SecurityGroupKind,
			Tenant:     scenarios.params.Tenant,
			Workspace:  params.Workspace.Name,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Spec:   &secalib.SecurityGroupSpecV1{},
	}

	for _, rules := range params.SecurityGroup.InitialSpec.Rules {
		securityGroupResponse.Spec.Rules = append(securityGroupResponse.Spec.Rules, &secalib.SecurityGroupRuleV1{
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
	if err := configurePutStub(wm, name, stubConfig{
		url:          securityGroupURL,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          securityGroupURL,
		params:       scenarios.params,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "GetSecurityGroup",
		nextState:    "CreatePublicIP",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Public-IP
	publicIPResponse := &resourceResponse[secalib.PublicIpSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.PublicIp.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   publicIPResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.PublicIPKind,
			Tenant:     scenarios.params.Tenant,
			Workspace:  params.Workspace.Name,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.PublicIpSpecV1{
			Version: params.PublicIp.InitialSpec.Version,
			Address: params.PublicIp.InitialSpec.Address,
		},
	}

	// Create public-IP
	publicIPResponse.Metadata.Verb = http.MethodPut
	publicIPResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.ResourceVersion = 1
	publicIPResponse.Status.State = secalib.CreatingStatusState
	publicIPResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)

	if err := configurePutStub(wm, name, stubConfig{
		url:          publicIPURL,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          publicIPURL,
		params:       scenarios.params,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "GetPublicIP",
		nextState:    "CreateNIC",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// NIC
	nicResponse := &resourceResponse[secalib.NICSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Nic.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   nicResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NicKind,
			Tenant:     scenarios.params.Tenant,
			Workspace:  params.Workspace.Name,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.NICSpecV1{
			Addresses: params.Nic.InitialSpec.Addresses,
			SubnetRef: params.Nic.InitialSpec.SubnetRef,
		},
	}
	// Create NIC
	nicResponse.Metadata.Verb = http.MethodPut
	nicResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.ResourceVersion = 1
	nicResponse.Status.State = secalib.CreatingStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          nicURL,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          nicURL,
		params:       scenarios.params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "GetNIC",
		nextState:    "CreateInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Compute
	instResponse := &resourceResponse[secalib.InstanceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Instance.Name,
			Provider:   secalib.ComputeProviderV1,
			Resource:   instanceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     scenarios.params.Tenant,
			Workspace:  params.Workspace.Name,
			Region:     scenarios.params.Region,
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
	if err := configurePutStub(wm, name, stubConfig{
		url:          instanceURL,
		params:       scenarios.params,
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
	if err := configureGetStub(wm, name, stubConfig{
		url:          instanceURL,
		params:       scenarios.params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: "GetCreatedInstance",
		nextState:    "DeleteInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Delete all
	// Delete Instance
	instResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          instanceURL,
		params:       scenarios.params,
		response:     instResponse,
		currentState: "DeleteInstance",
		nextState:    "DeleteSecurityGroup",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Security Group
	securityGroupResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          securityGroupURL,
		params:       scenarios.params,
		response:     securityGroupResponse,
		currentState: "DeleteSecurityGroup",
		nextState:    "DeleteNic",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Nic
	nicResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          nicURL,
		params:       scenarios.params,
		response:     nicResponse,
		currentState: "DeleteNic",
		nextState:    "DeletePublicIP",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete public ip
	publicIPResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          publicIPURL,
		params:       scenarios.params,
		response:     publicIPResponse,
		currentState: "DeletePublicIP",
		nextState:    "DeleteSubnet",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete subnet
	subnetResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          subnetURL,
		params:       scenarios.params,
		response:     subnetResponse,
		currentState: "DeleteSubnet",
		nextState:    "DeleteRouteTable",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Route-table
	routeTableResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          routeTableURL,
		params:       scenarios.params,
		response:     routeTableResponse,
		currentState: "DeleteRouteTable",
		nextState:    "DeleteInternetGateway",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Internet-gateway
	internetGatewayResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          internetGatewayURL,
		params:       scenarios.params,
		response:     internetGatewayResponse,
		currentState: "DeleteInternetGateway",
		nextState:    "DeleteNetwork",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Network
	networkResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          networkURL,
		params:       scenarios.params,
		response:     networkResponse,
		currentState: "DeleteNetwork",
		nextState:    "DeleteBlockStorage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete BlockStorage
	blockResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          blockStorageURL,
		params:       scenarios.params,
		response:     blockResponse,
		currentState: "DeleteBlockStorage",
		nextState:    "DeleteImage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Image
	imageResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          imageURL,
		params:       scenarios.params,
		response:     imageResponse,
		currentState: "DeleteImage",
		nextState:    "DeleteWorkspace",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Workspace
	workspaceResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          workspaceURL,
		params:       scenarios.params,
		response:     workspaceResponse,
		currentState: "DeleteWorkspace",
		nextState:    "DeleteRoleAssignment",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}
	// Delete Role assignment
	roleAssignmentResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          roleAssignmentUrl,
		params:       scenarios.params,
		response:     roleAssignmentResponse,
		currentState: "DeleteRoleAssignment",
		nextState:    "DeleteRole",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}
	// Delete Role
	roleResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          roleUrl,
		params:       scenarios.params,
		response:     roleResponse,
		currentState: "DeleteRole",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	return wm, err
}
