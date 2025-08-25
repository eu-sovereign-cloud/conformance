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

	//Generate URLs
	networkSkuURL := secalib.GenerateNetworkSkuURL(params.Tenant, params.Workspace, params.NetworkSku.Name)
	networkURL := secalib.GenerateNetworkURL(params.Tenant, params.Workspace, params.Network.Name)
	internetGatewayURL := secalib.GenerateInternetGatewayURL(params.Tenant, params.Workspace, params.InternetGateway.Name)
	nicURL := secalib.GenerateNicURL(params.Tenant, params.Workspace, params.NIC.Name)
	publicIPURL := secalib.GeneratePublicIPURL(params.Tenant, params.Workspace, params.PublicIP.Name)
	routeTableURL := secalib.GenerateRouteTableURL(params.Tenant, params.Workspace, params.RouteTable.Name)
	subnetURL := secalib.GenerateSubnetURL(params.Tenant, params.Workspace, params.Subnet.Name)
	securityGroupURL := secalib.GenerateSecurityGroupURL(params.Tenant, params.Workspace, params.SecurityGroup.Name)
	instanceSkuURL := secalib.GenerateInstanceSkuURL(params.Tenant, params.InstanceSku.Name)
	instanceURL := secalib.GenerateInstanceURL(params.Tenant, params.Workspace, params.Instance.Name)
	blockStorageURL := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace, params.BlockStorage.Name)

	//GenerateResources
	networkSkuResource := secalib.GenerateNetworkResource(params.Tenant, params.Workspace, params.NetworkSku.Name)
	networkResource := secalib.GenerateNetworkResource(params.Tenant, params.Workspace, params.Network.Name)
	internetGatewayResource := secalib.GenerateInternetGatewayResource(params.Tenant, params.Workspace, params.InternetGateway.Name)
	nicResource := secalib.GenerateNicResource(params.Tenant, params.Workspace, params.NIC.Name)
	publicIPResource := secalib.GeneratePublicIPResource(params.Tenant, params.Workspace, params.PublicIP.Name)
	routeTableResource := secalib.GenerateRouteTableResource(params.Tenant, params.Workspace, params.RouteTable.Name)
	subnetResource := secalib.GenerateSubnetResource(params.Tenant, params.Workspace, params.Subnet.Name)
	securityGroupResource := secalib.GenerateSecurityGroupResource(params.Tenant, params.Workspace, params.SecurityGroup.Name)
	instanceSkuResource := secalib.GenerateSkuResource(params.Tenant, params.InstanceSku.Name)
	instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace, params.Instance.Name)
	blockStorageResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace, params.BlockStorage.Name)

	// Network Sku
	networkSkuResponse := networkSkuResponseV1{
		Metadata: metadataResponse{
			Name:            params.NetworkSku.Name,
			Provider:        secalib.NetworkProviderV1,
			Resource:        networkSkuResource,
			Verb:            http.MethodGet,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      secalib.ApiVersion1,
			Kind:            secalib.NetworkSkuKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Bandwidth: params.NetworkSku.Bandwidth,
		Packets:   params.NetworkSku.Packets,
	}
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          networkSkuURL,
		params:       params,
		response:     networkSkuResponse,
		template:     networkSkuResponseTemplateV1,
		currentState: startedScenarioState,
		nextState:    "GetNetwork",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}
	// Network
	networkResponse := networkResponseV1{
		Metadata: metadataResponse{
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
		Status: statusResponse{
			State:            secalib.ActiveStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
		SkuRef:        params.Network.SkuRef,
		RouteTableRef: params.Network.RouteTableRef,
	}

	// Create  Network
	networkResponse.Metadata.Verb = http.MethodPut
	networkResponse.Status.State = secalib.ActiveStatusState
	networkResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	networkResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	networkResponse.Metadata.ResourceVersion = 1
	networkResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          networkURL,
		params:       params,
		response:     networkResponse,
		template:     networkResponseTemplateV1,
		currentState: "CreateNetwork",
		nextState:    "GetNetwork",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	//Get network
	networkResponse.Metadata.Verb = http.MethodGet
	networkResponse.Status.State = secalib.ActiveStatusState
	networkResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
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

	//Update Network
	networkResponse.Metadata.Verb = http.MethodPut
	networkResponse.Status.State = secalib.UpdatingStatusState
	networkResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	networkResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
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
	//Get network 2x time
	networkResponse.Metadata.Verb = http.MethodGet
	networkResponse.Status.State = secalib.ActiveStatusState
	if err := configurePutStub(wm, scenario, scenarioConfig{
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
	internetGatewayResponse := internetGatewayResponseV1{
		Metadata: metadataResponse{
			Name:       params.InternetGateway.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   internetGatewayResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		EgressOnly: params.InternetGateway.EgressOnly,
	}

	// Create internet-Gateway
	internetGatewayResponse.Metadata.Verb = http.MethodPut
	internetGatewayResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.ResourceVersion = 1
	internetGatewayResponse.Status.State = secalib.CreatingStatusState
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          internetGatewayURL,
		params:       params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "CreateInternet-gateway",
		nextState:    "GetInternet-gateway",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get internet-Gateway
	internetGatewayResponse.Metadata.Verb = http.MethodGet
	internetGatewayResponse.Status.State = secalib.ActiveStatusState
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          internetGatewayURL,
		params:       params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "GetInternet-gateway",
		nextState:    "UpdateInternet-gateway",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	//Update internet-gateway
	internetGatewayResponse.Metadata.Verb = http.MethodPut
	internetGatewayResponse.Status.State = secalib.UpdatingStatusState
	internetGatewayResponse.Metadata.ResourceVersion = internetGatewayResponse.Metadata.ResourceVersion + 1
	internetGatewayResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	internetGatewayResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          internetGatewayURL,
		params:       params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "UpdateInternet-gateway",
		nextState:    "GetInternet-gateway updated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get internet-gateway after update
	internetGatewayResponse.Metadata.Verb = http.MethodGet
	internetGatewayResponse.Status.State = secalib.ActiveStatusState
	if err := configureGetStub(wm, scenario, scenarioConfig{
		params:       params,
		response:     internetGatewayResponse,
		template:     internetGatewayResponseTemplateV1,
		currentState: "GetInternet-gateway2x",
		nextState:    "CreateRoute-table",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Route-Table
	routeTableResponse := routeTableResponseV1{
		Metadata: metadataResponse{
			Name:       params.RouteTable.Name,
			Provider:   networkProviderV1,
			Resource:   routeTableResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		LocalRef: params.RouteTable.LocalRef,
	}

	for _, routes := range params.RouteTable.Routes {
		routeTableResponse.Routes = append(routeTableResponse.Routes, Routes{
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
	if err := configurePostStub(wm, scenario, scenarioConfig{
		url:          routeTableURL,
		params:       params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "CreateRoute-table",
		nextState:    "GetRoute-table",
		httpStatus:   http.StatusCreated,
	}); err != nil {
		return nil, err
	}

	// Get route-Table
	routeTableResponse.Metadata.Verb = http.MethodGet
	routeTableResponse.Status.State = secalib.ActiveStatusState
	routeTableResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          routeTableURL,
		params:       params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "GetRoute-table",
		nextState:    "UpdateRoute-table",
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
	if err := configurePutStub(wm, scenario, scenarioConfig{
		url:          routeTableURL,
		params:       params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "UpdateRoute-table",
		nextState:    "GetRoute-tableUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get route-table after update
	routeTableResponse.Metadata.Verb = http.MethodGet
	routeTableResponse.Status.State = secalib.ActiveStatusState
	routeTableResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          routeTableURL,
		params:       params,
		response:     routeTableResponse,
		template:     routeTableResponseTemplateV1,
		currentState: "GetRoute-tableUpdated",
		nextState:    "CreateSubnet",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	//subnet
	subnetResponse := subnetResponseV1{
		Metadata: metadataResponse{
			Name:       params.Subnet.Name,
			Provider:   networkProviderV1,
			Resource:   subnetResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		Cidr: params.Subnet.Cidr,
		Zone: params.Subnet.Zone,
	}

	// Create subnet
	subnetResponse.Metadata.Verb = http.MethodPost
	subnetResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	subnetResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	subnetResponse.Metadata.ResourceVersion = 1
	subnetResponse.Status.State = secalib.CreatingStatusState
	subnetResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePostStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
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
	if err := configurePutStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
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
	publicIPResponse := publicIPResponseV1{
		Metadata: metadataResponse{
			Name:       params.PublicIP.Name,
			Provider:   networkProviderV1,
			Resource:   publicIPResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		Version: params.PublicIP.Version,
		Address: params.PublicIP.Address,
	}

	// Create public-IP
	publicIPResponse.Metadata.Verb = http.MethodPost
	publicIPResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	publicIPResponse.Metadata.ResourceVersion = 1
	publicIPResponse.Status.State = secalib.CreatingStatusState
	publicIPResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)

	if err := configurePostStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
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
	if err := configurePutStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
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
	nicResponse := nicResponseV1{
		Metadata: metadataResponse{
			Name:       params.NIC.Name,
			Provider:   networkProviderV1,
			Resource:   nicResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		Addresses: params.NIC.Addresses,
		SubnetRef: params.NIC.SubnetRef,
	}

	// Create NIC
	nicResponse.Metadata.Verb = http.MethodPost
	nicResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	nicResponse.Metadata.ResourceVersion = 1
	nicResponse.Status.State = secalib.CreatingStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePostStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          nicURL,
		params:       params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "GetNic",
		nextState:    "UpdateNic",
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
	if err := configurePutStub(wm, scenario, scenarioConfig{
		params:       params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "UpdateNic",
		nextState:    "GetNicUpdated",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Get NIC after update
	nicResponse.Metadata.Verb = http.MethodGet
	nicResponse.Status.State = secalib.ActiveStatusState
	nicResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          nicURL,
		params:       params,
		response:     nicResponse,
		template:     nicResponseTemplateV1,
		currentState: "GetNicUpdated",
		nextState:    "CreateSecurityGroup",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Security-group
	securityGroupResponse := securityGroupResponseV1{
		Metadata: metadataResponse{
			Name:       params.SecurityGroup.Name,
			Provider:   networkProviderV1,
			Resource:   securityGroupResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
	}

	for _, rules := range params.SecurityGroup.rules {
		securityGroupResponse.Rules = append(securityGroupResponse.Rules, Rules{
			Direction: rules.Direction,
		})
	}
	// Create Security-group
	securityGroupResponse.Metadata.Verb = http.MethodPost
	securityGroupResponse.Status.State = secalib.CreatingStatusState
	securityGroupResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	securityGroupResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	securityGroupResponse.Metadata.ResourceVersion = 1
	securityGroupResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePostStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
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
	if err := configurePutStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
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
	blockResponse := blockStorageResponseV1{
		Metadata: metadataResponse{
			Name:       params.BlockStorage.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     params.Tenant,
			Region:     params.Region,
		},
	}

	// Create a block storage
	blockResponse.Metadata.Verb = http.MethodPut
	blockResponse.Status.State = secalib.CreatingStatusState
	blockResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	blockResponse.Metadata.ResourceVersion = 1
	blockResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          blockStorageURL,
		params:       params,
		response:     blockResponse,
		template:     blockStorageResponseTemplateV1,
		currentState: "GetCreatedBlockStorage",
		nextState:    "GetInstanceSku",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Instance
	// Sku
	instSkuResponse := instanceSkuResponseV1{
		Metadata: metadataResponse{
			Name:            params.InstanceSku.Name,
			Provider:        secalib.ComputeProviderV1,
			Resource:        instanceSkuResource,
			Verb:            http.MethodGet,
			CreatedAt:       time.Now().Format(time.RFC3339),
			LastModifiedAt:  time.Now().Format(time.RFC3339),
			ResourceVersion: 1,
			ApiVersion:      secalib.ApiVersion1,
			Kind:            secalib.InstanceSkuKind,
			Tenant:          params.Tenant,
			Region:          params.Region,
		},
		Status: statusResponse{
			State:            secalib.ActiveStatusState,
			LastTransitionAt: time.Now().Format(time.RFC3339),
		},
		Architecture: params.InstanceSku.Architecture,
		Provider:     params.InstanceSku.Provider,
		Tier:         params.InstanceSku.Tier,
		RAM:          params.InstanceSku.RAM,
		VCPU:         params.InstanceSku.VCPU,
	}
	// Get sku
	if err := configureGetStub(wm, scenario, scenarioConfig{
		url:          instanceSkuURL,
		params:       params,
		response:     instSkuResponse,
		template:     instanceSkuResponseTemplateV1,
		currentState: "GetInstanceSku",
		nextState:    "CreateInstance",
		httpStatus:   http.StatusOK,
	}); err != nil {
		return nil, err
	}

	// Instance
	instResponse := instanceResponseV1{
		Metadata: metadataResponse{
			Name:       params.Instance.Name,
			Provider:   secalib.ComputeProviderV1,
			Resource:   instanceResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     params.Tenant,
			Workspace:  params.Workspace,
			Region:     params.Region,
		},
		SkuRef:        params.Instance.SkuRef,
		Zone:          params.Instance.ZoneInitial,
		BootDeviceRef: params.Instance.BootDeviceRef,
	}

	// Create an instance
	instResponse.Metadata.Verb = http.MethodPut
	instResponse.Metadata.CreatedAt = time.Now().Format(time.RFC3339)
	instResponse.Metadata.LastModifiedAt = time.Now().Format(time.RFC3339)
	instResponse.Metadata.ResourceVersion = 1
	instResponse.Status.State = secalib.CreatingStatusState
	instResponse.Status.LastTransitionAt = time.Now().Format(time.RFC3339)
	if err := configurePutStub(wm, scenario, scenarioConfig{
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
	if err := configureGetStub(wm, scenario, scenarioConfig{
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

	//Delete Instance
	instResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          instanceURL,
		params:       params,
		response:     instResponse,
		currentState: "DeleteInstance",
		nextState:    "DeleteBlockStorage",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	//Delete Block Storage
	instResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
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
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
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
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
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
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
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
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
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
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          routeTableURL,
		params:       params,
		response:     routeTableResponse,
		currentState: "DeleteRoutTeable",
		nextState:    "DeleteInternetGateway",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}

	// Delete Internet-gateway
	internetGatewayResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
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
	if err := configureDeleteStub(wm, scenario, scenarioConfig{
		url:          networkURL,
		params:       params,
		response:     networkResponse,
		currentState: "DeleteNetwork",
		httpStatus:   http.StatusAccepted,
	}); err != nil {
		return nil, err
	}
	// End
	return wm, nil
}
