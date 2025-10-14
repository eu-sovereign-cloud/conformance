package mock

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/wiremock/go-wiremock"
)

func CreateNetworkLifecycleScenarioV1(scenario string, params NetworkParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to Network Lifecycle Scenario")

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	blockStorageUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := secalib.GenerateInstanceURL(params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkUrl := secalib.GenerateNetworkURL(params.Tenant, params.Workspace.Name, params.Network.Name)
	internetGatewayUrl := secalib.GenerateInternetGatewayURL(params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicUrl := secalib.GenerateNicURL(params.Tenant, params.Workspace.Name, params.NIC.Name)
	publicIPUrl := secalib.GeneratePublicIPURL(params.Tenant, params.Workspace.Name, params.PublicIP.Name)
	routeTableUrl := secalib.GenerateRouteTableURL(params.Tenant, params.Workspace.Name, params.Network.Name, params.RouteTable.Name)
	subnetUrl := secalib.GenerateSubnetURL(params.Tenant, params.Workspace.Name, params.Network.Name, params.Subnet.Name)
	securityGroupUrl := secalib.GenerateSecurityGroupURL(params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Generate resources
	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)
	blockStorageResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkResource := secalib.GenerateNetworkResource(params.Tenant, params.Workspace.Name, params.Network.Name)
	internetGatewayResource := secalib.GenerateInternetGatewayResource(params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicResource := secalib.GenerateNicResource(params.Tenant, params.Workspace.Name, params.NIC.Name)
	publicIPResource := secalib.GeneratePublicIPResource(params.Tenant, params.Workspace.Name, params.PublicIP.Name)
	routeTableResource := secalib.GenerateRouteTableResource(params.Tenant, params.Workspace.Name, params.Network.Name, params.RouteTable.Name)
	subnetResource := secalib.GenerateSubnetResource(params.Tenant, params.Workspace.Name, params.Network.Name, params.Subnet.Name)
	securityGroupResource := secalib.GenerateSecurityGroupResource(params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Workspace
	workResponse := &resourceResponse[secalib.WorkspaceSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.Workspace.Name,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     params.Tenant,
			Region:     &params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          workspaceUrl,
		params:       params,
		response:     workResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          workspaceUrl,
		params:       params,
		response:     workResponse,
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
			Tenant:     params.Tenant,
			Workspace:  &params.Workspace.Name,
			Region:     &params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          networkUrl,
		params:       params,
		response:     networkResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          networkUrl,
		params:       params,
		response:     networkResponse,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          networkUrl,
		params:       params,
		response:     networkResponse,
		currentState: "UpdateNetwork",
		nextState:    "GetNetwork2x",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get the updated network
	networkResponse.Metadata.Verb = http.MethodGet
	networkResponse.Status.State = secalib.ActiveStatusState
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          networkUrl,
		params:       params,
		response:     networkResponse,
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
			Tenant:     params.Tenant,
			Workspace:  &params.Workspace.Name,
			Region:     &params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          internetGatewayUrl,
		params:       params,
		response:     internetGatewayResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          internetGatewayUrl,
		params:       params,
		response:     internetGatewayResponse,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          internetGatewayUrl,
		params:       params,
		response:     internetGatewayResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          internetGatewayUrl,
		params:       params,
		response:     internetGatewayResponse,
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
			Tenant:     params.Tenant,
			Workspace:  &params.Workspace.Name,
			Network:    &params.Network.Name,
			Region:     &params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          routeTableUrl,
		params:       params,
		response:     routeTableResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          routeTableUrl,
		params:       params,
		response:     routeTableResponse,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          routeTableUrl,
		params:       params,
		response:     routeTableResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          routeTableUrl,
		params:       params,
		response:     routeTableResponse,
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
			Tenant:     params.Tenant,
			Workspace:  &params.Workspace.Name,
			Network:    &params.Network.Name,
			Region:     &params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          subnetUrl,
		params:       params,
		response:     subnetResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          subnetUrl,
		params:       params,
		response:     subnetResponse,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          subnetUrl,
		params:       params,
		response:     subnetResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          subnetUrl,
		params:       params,
		response:     subnetResponse,
		currentState: "GetSubnetUpdated",
		nextState:    "CreatePublicIP",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Public ip
	publicIPResponse := &resourceResponse[secalib.PublicIpSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.PublicIP.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   publicIPResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.PublicIPKind,
			Tenant:     params.Tenant,
			Workspace:  &params.Workspace.Name,
			Region:     &params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.PublicIpSpecV1{
			Version: params.PublicIP.InitialSpec.Version,
			Address: params.PublicIP.InitialSpec.Address,
		},
	}

	// Create a public ip
	publicIPResponse.Metadata.Verb = http.MethodPut
	publicIPResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.ResourceVersion = 1
	publicIPResponse.Status.State = secalib.CreatingStatusState
	publicIPResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          publicIPUrl,
		params:       params,
		response:     publicIPResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          publicIPUrl,
		params:       params,
		response:     publicIPResponse,
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
	publicIPResponse.Spec = params.PublicIP.UpdatedSpec
	publicIPResponse.Status.State = secalib.UpdatingStatusState
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          publicIPUrl,
		params:       params,
		response:     publicIPResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          publicIPUrl,
		params:       params,
		response:     publicIPResponse,
		currentState: "GetPublicIPUpdated",
		nextState:    "CreateNIC",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// NIC
	nicResponse := &resourceResponse[secalib.NICSpecV1]{
		Metadata: &secalib.Metadata{
			Name:       params.NIC.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   nicResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NicKind,
			Tenant:     params.Tenant,
			Workspace:  &params.Workspace.Name,
			Region:     &params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.NICSpecV1{
			Addresses: params.NIC.InitialSpec.Addresses,
			SubnetRef: params.NIC.InitialSpec.SubnetRef,
		},
	}

	// Create a nic
	nicResponse.Metadata.Verb = http.MethodPut
	nicResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.ResourceVersion = 1
	nicResponse.Status.State = secalib.CreatingStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          nicUrl,
		params:       params,
		response:     nicResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          nicUrl,
		params:       params,
		response:     nicResponse,
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
	nicResponse.Spec = params.NIC.UpdatedSpec
	nicResponse.Status.State = secalib.UpdatingStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          nicUrl,
		params:       params,
		response:     nicResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          nicUrl,
		params:       params,
		response:     nicResponse,
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
			Tenant:     params.Tenant,
			Workspace:  &params.Workspace.Name,
			Region:     &params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          securityGroupUrl,
		params:       params,
		response:     securityGroupResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          securityGroupUrl,
		params:       params,
		response:     securityGroupResponse,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          securityGroupUrl,
		params:       params,
		response:     securityGroupResponse,
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
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          securityGroupUrl,
		params:       params,
		response:     securityGroupResponse,
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
			Tenant:     params.Tenant,
			Region:     &params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
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
		url:          blockStorageUrl,
		params:       params,
		response:     blockResponse,
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
			Tenant:     params.Tenant,
			Workspace:  &params.Workspace.Name,
			Region:     &params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
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
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "GetCreatedInstance",
		nextState:    "DeleteInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Delete the instance
	instResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		response:     instResponse,
		currentState: "DeleteInstance",
		nextState:    "GetDeletedInstance",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted instance
	instResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          instanceUrl,
		params:       params,
		currentState: "GetDeletedInstance",
		nextState:    "DeleteBlockStorage",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the block storage
	blockResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          blockStorageUrl,
		params:       params,
		response:     instResponse,
		currentState: "DeleteBlockStorage",
		nextState:    "GetDeletedBlockStorage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted block storage
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          blockStorageUrl,
		params:       params,
		currentState: "GetDeletedBlockStorage",
		nextState:    "DeleteSecurityGroup",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the security group
	securityGroupResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          securityGroupUrl,
		params:       params,
		response:     securityGroupResponse,
		currentState: "DeleteSecurityGroup",
		nextState:    "GetDeletedSecurityGroup",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted security group
	securityGroupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          securityGroupUrl,
		params:       params,
		currentState: "GetDeletedSecurityGroup",
		nextState:    "DeleteNic",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the nic
	nicResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          nicUrl,
		params:       params,
		response:     nicResponse,
		currentState: "DeleteNic",
		nextState:    "GetDeletedNic",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted nic
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          nicUrl,
		params:       params,
		currentState: "GetDeletedNic",
		nextState:    "DeletePublicIP",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the public ip
	publicIPResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          publicIPUrl,
		params:       params,
		response:     publicIPResponse,
		currentState: "DeletePublicIP",
		nextState:    "GetDeletedPublicIP",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted public ip
	publicIPResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          publicIPUrl,
		params:       params,
		currentState: "GetDeletedPublicIP",
		nextState:    "DeleteSubnet",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the subnet
	subnetResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          subnetUrl,
		params:       params,
		response:     subnetResponse,
		currentState: "DeleteSubnet",
		nextState:    "GetDeletedSubnet",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted subnet
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          subnetUrl,
		params:       params,
		currentState: "GetDeletedSubnet",
		nextState:    "DeleteRouteTable",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the route table
	routeTableResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          routeTableUrl,
		params:       params,
		response:     routeTableResponse,
		currentState: "DeleteRouteTable",
		nextState:    "GetDeletedRouteTable",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted route table
	routeTableResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          routeTableUrl,
		params:       params,
		currentState: "GetDeletedRouteTable",
		nextState:    "DeleteInternetGateway",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	internetGatewayResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          internetGatewayUrl,
		params:       params,
		response:     internetGatewayResponse,
		currentState: "DeleteInternetGateway",
		nextState:    "GetDeletedInternetGateway",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted internet gateway
	internetGatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          internetGatewayUrl,
		params:       params,
		currentState: "GetDeletedInternetGateway",
		nextState:    "DeleteNetwork",
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	// Delete the network
	networkResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          networkUrl,
		params:       params,
		response:     networkResponse,
		currentState: "DeleteNetwork",
		nextState:    "GetDeletedNetwork",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Get deleted network
	networkResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, stubConfig{
		url:          networkUrl,
		params:       params,
		currentState: "GetDeletedNetwork",
		nextState:    startedScenarioState,
		httpStatus:   http.StatusNotFound,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
