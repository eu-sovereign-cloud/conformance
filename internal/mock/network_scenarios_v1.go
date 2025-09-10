package mock

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/wiremock/go-wiremock"
)

type NetworkScenariosV1 struct {
	Scenarios
}

func NewNetworkScenariosV1(authToken string, tenant string, region string, mockURL string) *NetworkScenariosV1 {
	return &NetworkScenariosV1{
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

func (scenarios *NetworkScenariosV1) ConfigureLifecycleScenario(id string, params *secalib.NetworkLifeCycleParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to Network Lifecycle Scenario")

	name := "NetworkLifecycleV1_" + id

	wm, err := scenarios.newClient()
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := secalib.GenerateWorkspaceURL(scenarios.params.Tenant, params.Workspace.Name)
	blockStorageUrl := secalib.GenerateBlockStorageURL(scenarios.params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := secalib.GenerateInstanceURL(scenarios.params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkUrl := secalib.GenerateNetworkURL(scenarios.params.Tenant, params.Workspace.Name, params.Network.Name)
	internetGatewayUrl := secalib.GenerateInternetGatewayURL(scenarios.params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicUrl := secalib.GenerateNicURL(scenarios.params.Tenant, params.Workspace.Name, params.Nic.Name)
	publicIPUrl := secalib.GeneratePublicIPURL(scenarios.params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeTableUrl := secalib.GenerateRouteTableURL(scenarios.params.Tenant, params.Workspace.Name, params.RouteTable.Name)
	subnetUrl := secalib.GenerateSubnetURL(scenarios.params.Tenant, params.Workspace.Name, params.Subnet.Name)
	securityGroupUrl := secalib.GenerateSecurityGroupURL(scenarios.params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Generate resources
	workspaceResource := secalib.GenerateWorkspaceResource(scenarios.params.Tenant, params.Workspace.Name)
	blockStorageResource := secalib.GenerateBlockStorageResource(scenarios.params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceResource := secalib.GenerateInstanceResource(scenarios.params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkResource := secalib.GenerateNetworkResource(scenarios.params.Tenant, params.Workspace.Name, params.Network.Name)
	internetGatewayResource := secalib.GenerateInternetGatewayResource(scenarios.params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicResource := secalib.GenerateNicResource(scenarios.params.Tenant, params.Workspace.Name, params.Nic.Name)
	publicIPResource := secalib.GeneratePublicIPResource(scenarios.params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeTableResource := secalib.GenerateRouteTableResource(scenarios.params.Tenant, params.Workspace.Name, params.RouteTable.Name)
	subnetResource := secalib.GenerateSubnetResource(scenarios.params.Tenant, params.Workspace.Name, params.Subnet.Name)
	securityGroupResource := secalib.GenerateSecurityGroupResource(scenarios.params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Workspace
	workResponse := &resourceResponse[secalib.WorkspaceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Workspace.Name,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     scenarios.params.Tenant,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Labels: &[]secalib.Label{},
	}

	// Create a workspace
	workResponse.Metadata.Verb = http.MethodPut
	workResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	workResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	workResponse.Metadata.ResourceVersion = 1
	workResponse.Status.State = secalib.CreatingStatusState
	workResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          workspaceUrl,
		params:       scenarios.params,
		response:     workResponse,
		template:     workspaceResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    "GetCreatedWorkspace",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get created workspace
	workResponse.Metadata.Verb = http.MethodGet
	workResponse.Status.State = secalib.ActiveStatusState
	workResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          workspaceUrl,
		params:       scenarios.params,
		response:     workResponse,
		template:     workspaceResponseTemplateV1,
		currentState: "GetCreatedWorkspace",
		nextState:    "CreateNetwork",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Network
	networkResponse := &resourceResponse[secalib.NetworkSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Network.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   networkResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NetworkKind,
			Tenant:     scenarios.params.Tenant,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.NetworkSpecV1{
			Cidr: &secalib.NetworkSpecCIDRV1{
				Ipv4: params.Network.InitialSpec.Cidr.Ipv4,
				Ipv6: params.Network.InitialSpec.Cidr.Ipv6,
			},
			SkuRef:        params.Network.InitialSpec.SkuRef,
			RouteTableRef: params.Network.InitialSpec.RouteTableRef,
		},
	}

	// Create a network
	networkResponse.Metadata.Verb = http.MethodPut
	networkResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	networkResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	networkResponse.Metadata.ResourceVersion = 1
	networkResponse.Status.State = secalib.CreatingStatusState
	networkResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          networkUrl,
		params:       scenarios.params,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "CreateNetwork",
		nextState:    "GetNetwork",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get the created network
	networkResponse.Metadata.Verb = http.MethodGet
	networkResponse.Status.State = secalib.ActiveStatusState
	networkResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          networkUrl,
		params:       scenarios.params,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "GetNetwork",
		nextState:    "UpdateNetwork",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the network
	networkResponse.Metadata.Verb = http.MethodPut
	networkResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	networkResponse.Spec = params.Network.UpdatedSpec
	networkResponse.Status.State = secalib.UpdatingStatusState
	networkResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          networkUrl,
		params:       scenarios.params,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "UpdateNetwork",
		nextState:    "GetNetwork2x",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get the updated network
	networkResponse.Metadata.Verb = http.MethodGet
	networkResponse.Status.State = secalib.ActiveStatusState
	if err := configureGetStub(wm, name, stubConfig{
		url:          networkUrl,
		params:       scenarios.params,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "GetNetwork2x",
		nextState:    "CreateInternetGateway",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Internet gateway
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

	// Create an internet gateway
	internetGatewayResponse.Metadata.Verb = http.MethodPut
	internetGatewayResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.ResourceVersion = 1
	internetGatewayResponse.Status.State = secalib.CreatingStatusState
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          internetGatewayUrl,
		params:       scenarios.params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "CreateInternetGateway",
		nextState:    "GetInternetGateway",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get the created internet gateway
	internetGatewayResponse.Metadata.Verb = http.MethodGet
	internetGatewayResponse.Status.State = secalib.ActiveStatusState
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          internetGatewayUrl,
		params:       scenarios.params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "GetInternetGateway",
		nextState:    "UpdateInternetGateway",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the internet gateway
	internetGatewayResponse.Metadata.Verb = http.MethodPut
	internetGatewayResponse.Status.State = secalib.UpdatingStatusState
	internetGatewayResponse.Metadata.ResourceVersion = internetGatewayResponse.Metadata.ResourceVersion + 1
	internetGatewayResponse.Spec = params.InternetGateway.UpdatedSpec
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          internetGatewayUrl,
		params:       scenarios.params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "UpdateInternetGateway",
		nextState:    "GetInternetGateway2x",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get the updated internet gateway
	internetGatewayResponse.Metadata.Verb = http.MethodGet
	internetGatewayResponse.Status.State = secalib.ActiveStatusState
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          internetGatewayUrl,
		params:       scenarios.params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "GetInternetGateway2x",
		nextState:    "CreateRouteTable",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Route table
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

	// Create a route table
	routeTableResponse.Metadata.Verb = http.MethodPut
	routeTableResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	routeTableResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	routeTableResponse.Metadata.ResourceVersion = 1
	routeTableResponse.Status.State = secalib.CreatingStatusState
	routeTableResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          routeTableUrl,
		params:       scenarios.params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "CreateRouteTable",
		nextState:    "GetRouteTable",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get the created route table
	routeTableResponse.Metadata.Verb = http.MethodGet
	routeTableResponse.Status.State = secalib.ActiveStatusState
	routeTableResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          routeTableUrl,
		params:       scenarios.params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "GetRouteTable",
		nextState:    "UpdateRouteTable",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the route table
	routeTableResponse.Metadata.Verb = http.MethodPut
	routeTableResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	routeTableResponse.Metadata.ResourceVersion = routeTableResponse.Metadata.ResourceVersion + 1
	routeTableResponse.Spec = params.RouteTable.UpdatedSpec
	routeTableResponse.Status.State = secalib.UpdatingStatusState
	routeTableResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          routeTableUrl,
		params:       scenarios.params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "UpdateRouteTable",
		nextState:    "GetRouteTableUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get the updated route table
	routeTableResponse.Metadata.Verb = http.MethodGet
	routeTableResponse.Status.State = secalib.ActiveStatusState
	routeTableResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          routeTableUrl,
		params:       scenarios.params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "GetRouteTableUpdated",
		nextState:    "CreateSubnet",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Subnet
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

	// Create a subnet
	subnetResponse.Metadata.Verb = http.MethodPut
	subnetResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	subnetResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	subnetResponse.Metadata.ResourceVersion = 1
	subnetResponse.Status.State = secalib.CreatingStatusState
	subnetResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          subnetUrl,
		params:       scenarios.params,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "CreateSubnet",
		nextState:    "GetSubnet",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get the created subnet
	subnetResponse.Metadata.Verb = http.MethodGet
	subnetResponse.Status.State = secalib.ActiveStatusState
	subnetResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          subnetUrl,
		params:       scenarios.params,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "GetSubnet",
		nextState:    "UpdateSubnet",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the subnet
	subnetResponse.Metadata.Verb = http.MethodPut
	subnetResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	subnetResponse.Metadata.ResourceVersion = subnetResponse.Metadata.ResourceVersion + 1
	subnetResponse.Spec = params.Subnet.UpdatedSpec
	subnetResponse.Status.State = secalib.UpdatingStatusState
	subnetResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          subnetUrl,
		params:       scenarios.params,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "UpdateSubnet",
		nextState:    "GetSubnetUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get the updated subnet
	subnetResponse.Metadata.Verb = http.MethodGet
	subnetResponse.Status.State = secalib.ActiveStatusState
	subnetResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          subnetUrl,
		params:       scenarios.params,
		response:     subnetResponse,
		template:     subnetResponseTemplateV1,
		currentState: "GetSubnetUpdated",
		nextState:    "CreatePublicIP",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Public ip
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

	// Create a public ip
	publicIPResponse.Metadata.Verb = http.MethodPut
	publicIPResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.ResourceVersion = 1
	publicIPResponse.Status.State = secalib.CreatingStatusState
	publicIPResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          publicIPUrl,
		params:       scenarios.params,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "CreatePublicIP",
		nextState:    "GetPublicIP",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get the created public ip
	publicIPResponse.Metadata.Verb = http.MethodGet
	publicIPResponse.Status.State = secalib.ActiveStatusState
	publicIPResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          publicIPUrl,
		params:       scenarios.params,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "GetPublicIP",
		nextState:    "UpdatePublicIP",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the public ip
	publicIPResponse.Metadata.Verb = http.MethodPut
	publicIPResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.ResourceVersion = publicIPResponse.Metadata.ResourceVersion + 1
	publicIPResponse.Spec = params.PublicIp.UpdatedSpec
	publicIPResponse.Status.State = secalib.UpdatingStatusState
	if err := configurePutStub(wm, name, stubConfig{
		url:          publicIPUrl,
		params:       scenarios.params,
		response:     publicIPResponse,
		template:     publicIPResponseTemplateV1,
		currentState: "UpdatePublicIP",
		nextState:    "GetPublicIPUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get the updated public ip
	publicIPResponse.Metadata.Verb = http.MethodGet
	publicIPResponse.Status.State = secalib.ActiveStatusState
	publicIPResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          publicIPUrl,
		params:       scenarios.params,
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

	// Create a nic
	nicResponse.Metadata.Verb = http.MethodPut
	nicResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.ResourceVersion = 1
	nicResponse.Status.State = secalib.CreatingStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          nicUrl,
		params:       scenarios.params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "CreateNIC",
		nextState:    "GetNIC",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get the created nic
	nicResponse.Metadata.Verb = http.MethodGet
	nicResponse.Status.State = secalib.ActiveStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          nicUrl,
		params:       scenarios.params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "GetNIC",
		nextState:    "UpdateNIC",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the nic
	nicResponse.Metadata.Verb = http.MethodPut
	nicResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.ResourceVersion = nicResponse.Metadata.ResourceVersion + 1
	nicResponse.Spec = params.Nic.UpdatedSpec
	nicResponse.Status.State = secalib.UpdatingStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          nicUrl,
		params:       scenarios.params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "UpdateNIC",
		nextState:    "GetNICUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get the updated nic
	nicResponse.Metadata.Verb = http.MethodGet
	nicResponse.Status.State = secalib.ActiveStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          nicUrl,
		params:       scenarios.params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "GetNICUpdated",
		nextState:    "CreateSecurityGroup",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Security group
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

	// Create a security group
	securityGroupResponse.Metadata.Verb = http.MethodPut
	securityGroupResponse.Status.State = secalib.CreatingStatusState
	securityGroupResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	securityGroupResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	securityGroupResponse.Metadata.ResourceVersion = 1
	securityGroupResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          securityGroupUrl,
		params:       scenarios.params,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "CreateSecurityGroup",
		nextState:    "GetSecurityGroup",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get the created security group
	securityGroupResponse.Metadata.Verb = http.MethodGet
	securityGroupResponse.Status.State = secalib.ActiveStatusState
	securityGroupResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          securityGroupUrl,
		params:       scenarios.params,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "GetSecurityGroup",
		nextState:    "UpdateSecurityGroup",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Update the security group
	securityGroupResponse.Metadata.Verb = http.MethodPut
	securityGroupResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	securityGroupResponse.Metadata.ResourceVersion = securityGroupResponse.Metadata.ResourceVersion + 1
	securityGroupResponse.Spec = params.SecurityGroup.UpdatedSpec
	securityGroupResponse.Status.State = secalib.UpdatingStatusState
	securityGroupResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          securityGroupUrl,
		params:       scenarios.params,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "UpdateSecurityGroup",
		nextState:    "GetSecurityGroupUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get the updated security group
	securityGroupResponse.Metadata.Verb = http.MethodGet
	securityGroupResponse.Status.State = secalib.ActiveStatusState
	securityGroupResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, name, stubConfig{
		url:          securityGroupUrl,
		params:       scenarios.params,
		response:     securityGroupResponse,
		template:     securityGroupResponseTemplateV1,
		currentState: "GetSecurityGroupUpdated",
		nextState:    "CreateBlockStorage",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse := &resourceResponse[secalib.BlockStorageSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.BlockStorage.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     scenarios.params.Tenant,
			Region:     scenarios.params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.BlockStorageSpecV1{
			SkuRef: params.BlockStorage.InitialSpec.SkuRef,
		},
	}

	// Create a block storage
	blockResponse.Metadata.Verb = http.MethodPut
	blockResponse.Status.State = secalib.CreatingStatusState
	blockResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.ResourceVersion = 1
	blockResponse.Spec.SizeGB = params.BlockStorage.InitialSpec.SizeGB
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, name, stubConfig{
		url:          blockStorageUrl,
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
		url:          blockStorageUrl,
		params:       scenarios.params,
		response:     blockResponse,
		template:     blockStorageResponseTemplateV1,
		currentState: "GetCreatedBlockStorage",
		nextState:    "CreateInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Instance
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
		url:          instanceUrl,
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
		url:          instanceUrl,
		params:       scenarios.params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: "GetCreatedInstance",
		nextState:    "DeleteInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Delete the instance
	instResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          instanceUrl,
		params:       scenarios.params,
		response:     instResponse,
		currentState: "DeleteInstance",
		nextState:    "GetDeletedInstance",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted instance
	instResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          instanceUrl,
		params:       scenarios.params,
		currentState: "GetDeletedInstance",
		nextState:    "DeleteBlockStorage",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the block storage
	blockResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          blockStorageUrl,
		params:       scenarios.params,
		response:     instResponse,
		currentState: "DeleteBlockStorage",
		nextState:    "GetDeletedBlockStorage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted block storage
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          blockStorageUrl,
		params:       scenarios.params,
		currentState: "GetDeletedBlockStorage",
		nextState:    "DeleteSecurityGroup",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the security group
	securityGroupResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          securityGroupUrl,
		params:       scenarios.params,
		response:     securityGroupResponse,
		currentState: "DeleteSecurityGroup",
		nextState:    "GetDeletedSecurityGroup",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted security group
	securityGroupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          securityGroupUrl,
		params:       scenarios.params,
		currentState: "GetDeletedSecurityGroup",
		nextState:    "DeleteNic",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the nic
	nicResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          nicUrl,
		params:       scenarios.params,
		response:     nicResponse,
		currentState: "DeleteNic",
		nextState:    "GetDeletedNic",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted nic
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          nicUrl,
		params:       scenarios.params,
		currentState: "GetDeletedNic",
		nextState:    "DeletePublicIP",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the public ip
	publicIPResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          publicIPUrl,
		params:       scenarios.params,
		response:     publicIPResponse,
		currentState: "DeletePublicIP",
		nextState:    "GetDeletedPublicIP",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted public ip
	publicIPResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          publicIPUrl,
		params:       scenarios.params,
		currentState: "GetDeletedPublicIP",
		nextState:    "DeleteSubnet",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the subnet
	subnetResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          subnetUrl,
		params:       scenarios.params,
		response:     subnetResponse,
		currentState: "DeleteSubnet",
		nextState:    "GetDeletedSubnet",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted subnet
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          subnetUrl,
		params:       scenarios.params,
		currentState: "GetDeletedSubnet",
		nextState:    "DeleteRouteTable",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the route table
	routeTableResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          routeTableUrl,
		params:       scenarios.params,
		response:     routeTableResponse,
		currentState: "DeleteRouteTable",
		nextState:    "GetDeletedRouteTable",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted route table
	routeTableResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          routeTableUrl,
		params:       scenarios.params,
		currentState: "GetDeletedRouteTable",
		nextState:    "DeleteInternetGateway",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	internetGatewayResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          internetGatewayUrl,
		params:       scenarios.params,
		response:     internetGatewayResponse,
		currentState: "DeleteInternetGateway",
		nextState:    "GetDeletedInternetGateway",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted internet gateway
	internetGatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          internetGatewayUrl,
		params:       scenarios.params,
		currentState: "GetDeletedInternetGateway",
		nextState:    "DeleteNetwork",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the network
	networkResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, name, stubConfig{
		url:          networkUrl,
		params:       scenarios.params,
		response:     networkResponse,
		currentState: "DeleteNetwork",
		nextState:    "GetDeletedNetwork",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted network
	networkResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, name, stubConfig{
		url:          networkUrl,
		params:       scenarios.params,
		currentState: "GetDeletedNetwork",
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
