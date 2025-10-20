package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/wiremock/go-wiremock"
)

func CreateNetworkLifecycleScenarioV1(scenario string, params *NetworkParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	blockUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := secalib.GenerateInstanceURL(params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkUrl := secalib.GenerateNetworkURL(params.Tenant, params.Workspace.Name, params.Network.Name)
	gatewayUrl := secalib.GenerateInternetGatewayURL(params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicUrl := secalib.GenerateNicURL(params.Tenant, params.Workspace.Name, params.NIC.Name)
	publicIpUrl := secalib.GeneratePublicIpURL(params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeUrl := secalib.GenerateRouteTableURL(params.Tenant, params.Workspace.Name, params.Network.Name, params.RouteTable.Name)
	subnetUrl := secalib.GenerateSubnetURL(params.Tenant, params.Workspace.Name, params.Network.Name, params.Subnet.Name)
	groupUrl := secalib.GenerateSecurityGroupURL(params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Generate resources
	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)
	blockResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkResource := secalib.GenerateNetworkResource(params.Tenant, params.Workspace.Name, params.Network.Name)
	gatewayResource := secalib.GenerateInternetGatewayResource(params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicResource := secalib.GenerateNicResource(params.Tenant, params.Workspace.Name, params.NIC.Name)
	publicIpResource := secalib.GeneratePublicIpResource(params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeResource := secalib.GenerateRouteTableResource(params.Tenant, params.Workspace.Name, params.Network.Name, params.RouteTable.Name)
	subnetResource := secalib.GenerateSubnetResource(params.Tenant, params.Workspace.Name, params.Network.Name, params.Subnet.Name)
	groupResource := secalib.GenerateSecurityGroupResource(params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Workspace
	workResponse := newWorkspaceResponse(params.Workspace.Name, secalib.WorkspaceProviderV1, workspaceResource, secalib.ApiVersion1, params.Tenant, params.Region,
		params.Workspace.InitialLabels)

	// Create a workspace
	setCreatedRegionalResourceMetadata(workResponse.Metadata)
	workResponse.Status = newWorkspaceStatus(secalib.CreatingStatusState)
	workResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, response: workResponse, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get created workspace
	setWorkspaceStatusState(workResponse.Status, secalib.ActiveStatusState)
	workResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, response: workResponse, currentState: "GetCreatedWorkspace", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	// Network
	networkResponse := newNetworkResponse(params.Network.Name, secalib.NetworkProviderV1, networkResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.Network.InitialSpec)

	// Create a network
	setCreatedRegionalWorkspaceResourceMetadata(networkResponse.Metadata)
	networkResponse.Status = newNetworkStatus(secalib.CreatingStatusState)
	networkResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, response: networkResponse, currentState: "CreateNetwork", nextState: "GetNetwork"}); err != nil {
		return nil, err
	}

	// Get the created network
	setNetworkStatusState(networkResponse.Status, secalib.ActiveStatusState)
	networkResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, response: networkResponse, currentState: "GetNetwork", nextState: "UpdateNetwork"}); err != nil {
		return nil, err
	}

	// Update the network
	setNetworkStatusState(networkResponse.Status, secalib.UpdatingStatusState)
	networkResponse.Spec = *params.Network.UpdatedSpec
	networkResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, response: networkResponse, currentState: "UpdateNetwork", nextState: "GetNetwork2x"}); err != nil {
		return nil, err
	}

	// Get the updated network
	setNetworkStatusState(networkResponse.Status, secalib.ActiveStatusState)
	networkResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, response: networkResponse, currentState: "GetNetwork2x", nextState: "CreateInternetGateway"}); err != nil {
		return nil, err
	}

	// Internet gateway
	gatewayResponse := newInternetGatewayResponse(params.InternetGateway.Name, secalib.NetworkProviderV1, gatewayResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.InternetGateway.InitialSpec)

	// Create an internet gateway
	setCreatedRegionalWorkspaceResourceMetadata(gatewayResponse.Metadata)
	gatewayResponse.Status = newStatus(secalib.CreatingStatusState)
	gatewayResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, response: gatewayResponse, currentState: "CreateInternetGateway", nextState: "GetInternetGateway"}); err != nil {
		return nil, err
	}

	// Get the created internet gateway
	setStatusState(gatewayResponse.Status, secalib.ActiveStatusState)
	gatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, response: gatewayResponse, currentState: "GetInternetGateway", nextState: "UpdateInternetGateway"}); err != nil {
		return nil, err
	}

	// Update the internet gateway
	setModifiedRegionalWorkspaceResourceMetadata(gatewayResponse.Metadata)
	setStatusState(gatewayResponse.Status, secalib.UpdatingStatusState)
	gatewayResponse.Spec = *params.InternetGateway.UpdatedSpec
	gatewayResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, response: gatewayResponse, currentState: "UpdateInternetGateway", nextState: "GetInternetGateway2x"}); err != nil {
		return nil, err
	}

	// Get the updated internet gateway
	setStatusState(gatewayResponse.Status, secalib.ActiveStatusState)
	gatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, response: gatewayResponse, currentState: "GetInternetGateway2x", nextState: "CreateRouteTable"}); err != nil {
		return nil, err
	}

	// Route table
	routeResponse := newRouteTableResponse(params.RouteTable.Name, secalib.NetworkProviderV1, routeResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Network.Name, params.Region,
		params.RouteTable.InitialSpec)

	// Create a route table
	setCreatedRegionalNetworkResourceMetadata(routeResponse.Metadata)
	routeResponse.Status = newRouteTableStatus(secalib.CreatingStatusState)
	routeResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, response: routeResponse, currentState: "CreateRouteTable", nextState: "GetRouteTable"}); err != nil {
		return nil, err
	}

	// Get the created route table
	setRouteTableStatusState(routeResponse.Status, secalib.ActiveStatusState)
	routeResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, response: routeResponse, currentState: "GetRouteTable", nextState: "UpdateRouteTable"}); err != nil {
		return nil, err
	}

	// Update the route table
	setModifiedRegionalNetworkResourceMetadata(routeResponse.Metadata)
	setRouteTableStatusState(routeResponse.Status, secalib.UpdatingStatusState)
	routeResponse.Spec = *params.RouteTable.UpdatedSpec
	routeResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, response: routeResponse, currentState: "UpdateRouteTable", nextState: "GetRouteTableUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated route table
	setRouteTableStatusState(routeResponse.Status, secalib.ActiveStatusState)
	routeResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, response: routeResponse, currentState: "GetRouteTableUpdated", nextState: "CreateSubnet"}); err != nil {
		return nil, err
	}

	// Subnet
	subnetResponse := newSubnetResponse(params.Subnet.Name, secalib.NetworkProviderV1, subnetResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Network.Name, params.Region,
		params.Subnet.InitialSpec)

	// Create a subnet
	setCreatedRegionalNetworkResourceMetadata(subnetResponse.Metadata)
	subnetResponse.Status = newSubnetStatus(secalib.CreatingStatusState)
	subnetResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, response: subnetResponse, currentState: "CreateSubnet", nextState: "GetSubnet"}); err != nil {
		return nil, err
	}

	// Get the created subnet
	setSubnetStatusState(subnetResponse.Status, secalib.ActiveStatusState)
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, response: subnetResponse, currentState: "GetSubnet", nextState: "UpdateSubnet"}); err != nil {
		return nil, err
	}

	// Update the subnet
	setModifiedRegionalNetworkResourceMetadata(subnetResponse.Metadata)
	setSubnetStatusState(subnetResponse.Status, secalib.UpdatingStatusState)
	subnetResponse.Spec = *params.Subnet.UpdatedSpec
	subnetResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, response: subnetResponse, currentState: "UpdateSubnet", nextState: "GetSubnetUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated subnet
	setSubnetStatusState(subnetResponse.Status, secalib.ActiveStatusState)
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, response: subnetResponse, currentState: "GetSubnetUpdated", nextState: "CreatePublicIp"}); err != nil {
		return nil, err
	}

	// Public ip
	publicIpResponse := newPublicIpResponse(params.PublicIp.Name, secalib.NetworkProviderV1, publicIpResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.PublicIp.InitialSpec)

	// Create a public ip
	setCreatedRegionalWorkspaceResourceMetadata(publicIpResponse.Metadata)
	publicIpResponse.Status = newPublicIpStatus(secalib.CreatingStatusState)
	publicIpResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, response: publicIpResponse, currentState: "CreatePublicIp", nextState: "GetPublicIp"}); err != nil {
		return nil, err
	}

	// Get the created public ip
	setPublicIpStatusState(publicIpResponse.Status, secalib.ActiveStatusState)
	publicIpResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, response: publicIpResponse, currentState: "GetPublicIp", nextState: "UpdatePublicIp"}); err != nil {
		return nil, err
	}

	// Update the public ip
	setModifiedRegionalWorkspaceResourceMetadata(publicIpResponse.Metadata)
	setPublicIpStatusState(publicIpResponse.Status, secalib.UpdatingStatusState)
	publicIpResponse.Spec = *params.PublicIp.UpdatedSpec
	publicIpResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, response: publicIpResponse, currentState: "UpdatePublicIp", nextState: "GetPublicIpUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated public ip
	setPublicIpStatusState(publicIpResponse.Status, secalib.ActiveStatusState)
	publicIpResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, response: publicIpResponse, currentState: "GetPublicIpUpdated", nextState: "CreateNIC"}); err != nil {
		return nil, err
	}

	// NIC
	nicResponse := newNicResponse(params.NIC.Name, secalib.NetworkProviderV1, nicResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.NIC.InitialSpec)

	// Create a nic
	setCreatedRegionalWorkspaceResourceMetadata(nicResponse.Metadata)
	nicResponse.Status = newNicStatus(secalib.CreatingStatusState)
	nicResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, response: nicResponse, currentState: "CreateNIC", nextState: "GetNIC"}); err != nil {
		return nil, err
	}

	// Get the created nic
	setNicStatusState(nicResponse.Status, secalib.ActiveStatusState)
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, response: nicResponse, currentState: "GetNIC", nextState: "UpdateNIC"}); err != nil {
		return nil, err
	}

	// Update the nic
	setModifiedRegionalWorkspaceResourceMetadata(nicResponse.Metadata)
	setNicStatusState(nicResponse.Status, secalib.UpdatingStatusState)
	nicResponse.Spec = *params.NIC.UpdatedSpec
	nicResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, response: nicResponse, currentState: "UpdateNIC", nextState: "GetNICUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated nic
	setNicStatusState(nicResponse.Status, secalib.ActiveStatusState)
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, response: nicResponse, currentState: "GetNICUpdated", nextState: "CreateSecurityGroup"}); err != nil {
		return nil, err
	}

	// Security group
	groupResponse := newSecurityGroupResponse(params.SecurityGroup.Name, secalib.NetworkProviderV1, secalib.ApiVersion1, groupResource,
		params.Tenant, params.Workspace.Name, params.Region,
		params.SecurityGroup.InitialSpec)

	// Create a security group
	setCreatedRegionalWorkspaceResourceMetadata(groupResponse.Metadata)
	groupResponse.Status = newSecurityGroupStatus(secalib.CreatingStatusState)
	groupResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, response: groupResponse, currentState: "CreateSecurityGroup", nextState: "GetSecurityGroup"}); err != nil {
		return nil, err
	}

	// Get the created security group
	setSecurityGroupStatusState(groupResponse.Status, secalib.ActiveStatusState)
	groupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, response: groupResponse, currentState: "GetSecurityGroup", nextState: "UpdateSecurityGroup"}); err != nil {
		return nil, err
	}

	// Update the security group
	setModifiedRegionalWorkspaceResourceMetadata(groupResponse.Metadata)
	setSecurityGroupStatusState(groupResponse.Status, secalib.UpdatingStatusState)
	groupResponse.Spec = *params.SecurityGroup.UpdatedSpec
	groupResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, response: groupResponse, currentState: "UpdateSecurityGroup", nextState: "GetSecurityGroupUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated security group
	setSecurityGroupStatusState(groupResponse.Status, secalib.ActiveStatusState)
	groupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, response: groupResponse, currentState: "GetSecurityGroupUpdated", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse := newBlockStorageResponse(params.BlockStorage.Name, secalib.ComputeProviderV1, blockResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.BlockStorage.InitialSpec)

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = newBlockStorageStatus(secalib.CreatingStatusState)
	blockResponse.Spec.SizeGB = params.BlockStorage.InitialSpec.SizeGB
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, response: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get created block storage
	setBlockStorageStatusState(blockResponse.Status, secalib.ActiveStatusState)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, response: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "CreateInstance"}); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse := newInstanceResponse(params.Instance.Name, secalib.ComputeProviderV1, instanceResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.Instance.InitialSpec)

	// Create an instance
	setCreatedRegionalWorkspaceResourceMetadata(instanceResponse.Metadata)
	instanceResponse.Status = newInstanceStatus(secalib.CreatingStatusState)
	instanceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "CreateInstance", nextState: "GetCreatedInstance"}); err != nil {
		return nil, err
	}

	// Get created instance
	setInstanceStatusState(instanceResponse.Status, secalib.ActiveStatusState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "GetCreatedInstance", nextState: "DeleteInstance"}); err != nil {
		return nil, err
	}

	// Delete the instance
	instanceResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "DeleteInstance", nextState: "GetDeletedInstance"}); err != nil {
		return nil, err
	}

	// Get deleted instance
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: instanceUrl, params: params, currentState: "GetDeletedInstance", nextState: "DeleteBlockStorage"}); err != nil {
		return nil, err
	}

	// Delete the block storage
	blockResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, response: instanceResponse, currentState: "DeleteBlockStorage", nextState: "GetDeletedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get deleted block storage
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: blockUrl, params: params, currentState: "GetDeletedBlockStorage", nextState: "DeleteSecurityGroup"}); err != nil {
		return nil, err
	}

	// Delete the security group
	groupResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, response: groupResponse, currentState: "DeleteSecurityGroup", nextState: "GetDeletedSecurityGroup"}); err != nil {
		return nil, err
	}

	// Get deleted security group
	groupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: groupUrl, params: params, currentState: "GetDeletedSecurityGroup", nextState: "DeleteNic"}); err != nil {
		return nil, err
	}

	// Delete the nic
	nicResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, response: nicResponse, currentState: "DeleteNic", nextState: "GetDeletedNic"}); err != nil {
		return nil, err
	}

	// Get deleted nic
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: nicUrl, params: params, currentState: "GetDeletedNic", nextState: "DeletePublicIp"}); err != nil {
		return nil, err
	}

	// Delete the public ip
	publicIpResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, response: publicIpResponse, currentState: "DeletePublicIp", nextState: "GetDeletedPublicIp"}); err != nil {
		return nil, err
	}

	// Get deleted public ip
	publicIpResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: publicIpUrl, params: params, currentState: "GetDeletedPublicIp", nextState: "DeleteSubnet"}); err != nil {
		return nil, err
	}

	// Delete the subnet
	subnetResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, response: subnetResponse, currentState: "DeleteSubnet", nextState: "GetDeletedSubnet"}); err != nil {
		return nil, err
	}

	// Get deleted subnet
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: subnetUrl, params: params, currentState: "GetDeletedSubnet", nextState: "DeleteRouteTable"}); err != nil {
		return nil, err
	}

	// Delete the route table
	routeResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, response: routeResponse, currentState: "DeleteRouteTable", nextState: "GetDeletedRouteTable"}); err != nil {
		return nil, err
	}

	// Get deleted route table
	routeResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: routeUrl, params: params, currentState: "GetDeletedRouteTable", nextState: "DeleteInternetGateway"}); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	gatewayResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, response: gatewayResponse, currentState: "DeleteInternetGateway", nextState: "GetDeletedInternetGateway"}); err != nil {
		return nil, err
	}

	// Get deleted internet gateway
	gatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: gatewayUrl, params: params, currentState: "GetDeletedInternetGateway", nextState: "DeleteNetwork"}); err != nil {
		return nil, err
	}

	// Delete the network
	networkResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, response: networkResponse, currentState: "DeleteNetwork", nextState: "GetDeletedNetwork"}); err != nil {
		return nil, err
	}

	// Get deleted network
	networkResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario, http.StatusNotFound,
		&stubConfig{url: networkUrl, params: params, currentState: "GetDeletedNetwork", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}
