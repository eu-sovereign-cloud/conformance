package mock

import (
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/wiremock/go-wiremock"
)

func CreateNetworkLifecycleScenarioV1(scenario string, params NetworkParamsV1) (*wiremock.Client, error) {
	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	blockStorageURL := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace, params.BlockStorage.Name)
	instanceURL := secalib.GenerateInstanceURL(params.Tenant, params.Workspace, params.Instance.Name)

	networkURL := secalib.GenerateNetworkURL(params.Tenant, params.Workspace, params.Network.Name)
	internetGatewayURL := secalib.GenerateInternetGatewayURL(params.Tenant, params.Workspace, params.InternetGateway.Name)
	nicURL := secalib.GenerateNicURL(params.Tenant, params.Workspace, params.NIC.Name)
	publicIPURL := secalib.GeneratePublicIPURL(params.Tenant, params.Workspace, params.PublicIP.Name)
	routeTableURL := secalib.GenerateRouteTableURL(params.Tenant, params.Workspace, params.RouteTable.Name)
	subnetURL := secalib.GenerateSubnetURL(params.Tenant, params.Workspace, params.Subnet.Name)
	securityGroupURL := secalib.GenerateSecurityGroupURL(params.Tenant, params.Workspace, params.SecurityGroup.Name)

	// GenerateResources
	blockStorageResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace, params.BlockStorage.Name)
	instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace, params.Instance.Name)

	networkResource := secalib.GenerateNetworkResource(params.Tenant, params.Workspace, params.Network.Name)
	internetGatewayResource := secalib.GenerateInternetGatewayResource(params.Tenant, params.Workspace, params.InternetGateway.Name)
	nicResource := secalib.GenerateNicResource(params.Tenant, params.Workspace, params.NIC.Name)
	publicIPResource := secalib.GeneratePublicIPResource(params.Tenant, params.Workspace, params.PublicIP.Name)
	routeTableResource := secalib.GenerateRouteTableResource(params.Tenant, params.Workspace, params.RouteTable.Name)
	subnetResource := secalib.GenerateSubnetResource(params.Tenant, params.Workspace, params.Subnet.Name)
	securityGroupResource := secalib.GenerateSecurityGroupResource(params.Tenant, params.Workspace, params.SecurityGroup.Name)

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
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.NetworkSpecV1{
			Cidr: &secalib.NetworkSpecCIDRV1{
				Ipv4: "10.0.0.0/16",
			},
			SkuRef:        params.Network.InitialSpec.SkuRef,
			RouteTableRef: params.Network.InitialSpec.RouteTableRef,
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
		params:       params,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: startedScenarioState,
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
		params:       params,
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
			Ipv4: "10.0.0.0/16",
		},
		SkuRef:        params.Network.UpdatedSpec.SkuRef,
		RouteTableRef: params.Network.UpdatedSpec.RouteTableRef,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          networkURL,
		params:       params,
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
		params:       params,
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
			Name:       params.InternetGateway.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   internetGatewayResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          internetGatewayURL,
		params:       params,
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
		params:       params,
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
		EgressOnly: params.InternetGateway.UpdatedSpec.EgressOnly,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          internetGatewayURL,
		params:       params,
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
		params:       params,
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
			Name:       params.RouteTable.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   routeTableResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          routeTableURL,
		params:       params,
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
		params:       params,
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
		LocalRef: params.RouteTable.UpdatedSpec.LocalRef,
		Routes:   params.RouteTable.UpdatedSpec.Routes,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          routeTableURL,
		params:       params,
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
		params:       params,
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
			Name:       params.Subnet.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   subnetResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
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
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          subnetURL,
		params:       params,
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
		params:       params,
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
		Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: params.Subnet.UpdatedSpec.Cidr.Ipv4},
		Zone: params.Subnet.UpdatedSpec.Zone,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          subnetURL,
		params:       params,
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
		params:       params,
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
			Name:       params.PublicIP.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   publicIPResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.PublicIPSpecV1{
			Version: params.PublicIP.InitialSpec.Version,
			Address: params.PublicIP.InitialSpec.Address,
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
		params:       params,
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
		params:       params,
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
		Version: params.PublicIP.UpdatedSpec.Version,
		Address: params.PublicIP.UpdatedSpec.Address,
	}

	if err := configurePutStub(wm, scenario, stubConfig{
		url:          publicIPURL,
		params:       params,
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
		params:       params,
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
			Name:       params.NIC.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   nicResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.NICSpecV1{
			Addresses: params.NIC.InitialSpec.Addresses,
			SubnetRef: params.NIC.InitialSpec.SubnetRef,
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
		params:       params,
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
		params:       params,
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
		Addresses: params.NIC.UpdatedSpec.Addresses,
		SubnetRef: params.NIC.UpdatedSpec.SubnetRef,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          nicURL,
		params:       params,
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
		params:       params,
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
			Name:       params.SecurityGroup.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   securityGroupResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		Status: &secalib.Status{},
		Spec:   &secalib.SecurityGroupSpecV1{},
	}

	for _, rules := range params.SecurityGroup.InitialSpec.Rules {
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
		params:       params,
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
		params:       params,
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
		Rules: params.SecurityGroup.UpdatedSpec.Rules,
	}
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          securityGroupURL,
		params:       params,
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
		params:       params,
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
			Tenant:     params.Tenant,
			Region:     params.Region,
		},
		Status: &secalib.Status{},
		Spec: &secalib.BlockStorageSpecV1{
			SkuRef: params.BlockStorage.InitialSpec.SkuRef,
			SizeGB: params.BlockStorage.InitialSpec.SizeGB,
		},
	}

	// Create a block storage
	blockResponse.Metadata.Verb = http.MethodPut
	blockResponse.Status.State = secalib.CreatingStatusState
	blockResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.ResourceVersion = 1
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, stubConfig{
		url:          blockStorageURL,
		params:       params,
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
		params:       params,
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
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
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
		url:          instanceURL,
		params:       params,
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
		params:       params,
		response:     instResponse,
		template:     instanceResponseTemplateV1,
		currentState: "GetCreatedInstance",
		nextState:    "DeleteInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Delete Instance
	instResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          instanceURL,
		params:       params,
		response:     instResponse,
		currentState: "DeleteInstance",
		nextState:    "DeleteBlockStorage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Block Storage
	instResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          blockStorageURL,
		params:       params,
		response:     instResponse,
		currentState: "DeleteBlockStorage",
		nextState:    "DeleteSecurityGroup",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Security Group
	securityGroupResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, stubConfig{
		url:          securityGroupURL,
		params:       params,
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
		params:       params,
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
		params:       params,
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
		params:       params,
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
		params:       params,
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
		params:       params,
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
		params:       params,
		response:     networkResponse,
		currentState: "DeleteNetwork",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}
	return wm, nil
}
