package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/wiremock/go-wiremock"
)

func ConfigNetworkLifecycleScenarioV1(scenario string, params *NetworkParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	blockUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := secalib.GenerateInstanceURL(params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkUrl := secalib.GenerateNetworkURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name)
	gatewayUrl := secalib.GenerateInternetGatewayURL(params.Tenant, params.Workspace.Name, (*params.InternetGateway)[0].Name)
	nicUrl := secalib.GenerateNicURL(params.Tenant, params.Workspace.Name, (*params.NIC)[0].Name)
	publicIpUrl := secalib.GeneratePublicIpURL(params.Tenant, params.Workspace.Name, (*params.PublicIp)[0].Name)
	routeUrl := secalib.GenerateRouteTableURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.RouteTable)[0].Name)
	subnetUrl := secalib.GenerateSubnetURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.Subnet)[0].Name)
	groupUrl := secalib.GenerateSecurityGroupURL(params.Tenant, params.Workspace.Name, (*params.SecurityGroup)[0].Name)

	// Generate resources
	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)
	blockResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkResource := secalib.GenerateNetworkResource(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name)
	gatewayResource := secalib.GenerateInternetGatewayResource(params.Tenant, params.Workspace.Name, (*params.InternetGateway)[0].Name)
	nicResource := secalib.GenerateNicResource(params.Tenant, params.Workspace.Name, (*params.NIC)[0].Name)
	publicIpResource := secalib.GeneratePublicIpResource(params.Tenant, params.Workspace.Name, (*params.PublicIp)[0].Name)
	routeResource := secalib.GenerateRouteTableResource(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.RouteTable)[0].Name)
	subnetResource := secalib.GenerateSubnetResource(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.Subnet)[0].Name)
	groupResource := secalib.GenerateSecurityGroupResource(params.Tenant, params.Workspace.Name, (*params.SecurityGroup)[0].Name)

	// Workspace
	workResponse := newWorkspaceResponse(params.Workspace.Name, secalib.WorkspaceProviderV1, workspaceResource, secalib.ApiVersion1, params.Tenant, params.Region,
		params.Workspace.InitialLabels)

	// Create a workspace
	setCreatedRegionalResourceMetadata(workResponse.Metadata)
	workResponse.Status = secalib.NewWorkspaceStatus(secalib.CreatingResourceState)
	workResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workResponse, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	secalib.SetWorkspaceStatusState(workResponse.Status, secalib.ActiveResourceState)
	workResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workResponse, currentState: "GetCreatedWorkspace", nextState: "CreateNetwork"}); err != nil {
		return nil, err
	}

	// Network
	networkResponse := newNetworkResponse((*params.Network)[0].Name, secalib.NetworkProviderV1, networkResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		(*params.Network)[0].InitialSpec)

	// Create a network
	setCreatedRegionalWorkspaceResourceMetadata(networkResponse.Metadata)
	networkResponse.Status = secalib.NewNetworkStatus(secalib.CreatingResourceState)
	networkResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "CreateNetwork", nextState: "GetNetwork"}); err != nil {
		return nil, err
	}

	// Get the created network
	secalib.SetNetworkStatusState(networkResponse.Status, secalib.ActiveResourceState)
	networkResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "GetNetwork", nextState: "UpdateNetwork"}); err != nil {
		return nil, err
	}

	// Update the network
	secalib.SetNetworkStatusState(networkResponse.Status, secalib.UpdatingResourceState)
	networkResponse.Spec = *(*params.Network)[0].UpdatedSpec
	networkResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "UpdateNetwork", nextState: "GetNetwork2x"}); err != nil {
		return nil, err
	}

	// Get the updated network
	secalib.SetNetworkStatusState(networkResponse.Status, secalib.ActiveResourceState)
	networkResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "GetNetwork2x", nextState: "CreateInternetGateway"}); err != nil {
		return nil, err
	}

	// Internet gateway
	gatewayResponse := newInternetGatewayResponse((*params.InternetGateway)[0].Name, secalib.NetworkProviderV1, gatewayResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		(*params.InternetGateway)[0].InitialSpec)

	// Create an internet gateway
	setCreatedRegionalWorkspaceResourceMetadata(gatewayResponse.Metadata)
	gatewayResponse.Status = secalib.NewResourceStatus(secalib.CreatingResourceState)
	gatewayResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "CreateInternetGateway", nextState: "GetInternetGateway"}); err != nil {
		return nil, err
	}

	// Get the created internet gateway
	secalib.SetStatusState(gatewayResponse.Status, secalib.ActiveResourceState)
	gatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "GetInternetGateway", nextState: "UpdateInternetGateway"}); err != nil {
		return nil, err
	}

	// Update the internet gateway
	setModifiedRegionalWorkspaceResourceMetadata(gatewayResponse.Metadata)
	secalib.SetStatusState(gatewayResponse.Status, secalib.UpdatingResourceState)
	gatewayResponse.Spec = *(*params.InternetGateway)[0].UpdatedSpec
	gatewayResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "UpdateInternetGateway", nextState: "GetInternetGateway2x"}); err != nil {
		return nil, err
	}

	// Get the updated internet gateway
	secalib.SetStatusState(gatewayResponse.Status, secalib.ActiveResourceState)
	gatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "GetInternetGateway2x", nextState: "CreateRouteTable"}); err != nil {
		return nil, err
	}

	// Route table
	routeResponse := newRouteTableResponse((*params.RouteTable)[0].Name, secalib.NetworkProviderV1, routeResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, params.Region,
		(*params.RouteTable)[0].InitialSpec)

	// Create a route table
	setCreatedRegionalNetworkResourceMetadata(routeResponse.Metadata)
	routeResponse.Status = secalib.NewRouteTableStatus(secalib.CreatingResourceState)
	routeResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "CreateRouteTable", nextState: "GetRouteTable"}); err != nil {
		return nil, err
	}

	// Get the created route table
	secalib.SetRouteTableStatusState(routeResponse.Status, secalib.ActiveResourceState)
	routeResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "GetRouteTable", nextState: "UpdateRouteTable"}); err != nil {
		return nil, err
	}

	// Update the route table
	setModifiedRegionalNetworkResourceMetadata(routeResponse.Metadata)
	secalib.SetRouteTableStatusState(routeResponse.Status, secalib.UpdatingResourceState)
	routeResponse.Spec = *(*params.RouteTable)[0].UpdatedSpec
	routeResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "UpdateRouteTable", nextState: "GetRouteTableUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated route table
	secalib.SetRouteTableStatusState(routeResponse.Status, secalib.ActiveResourceState)
	routeResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "GetRouteTableUpdated", nextState: "CreateSubnet"}); err != nil {
		return nil, err
	}

	// Subnet
	subnetResponse := newSubnetResponse((*params.Subnet)[0].Name, secalib.NetworkProviderV1, subnetResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, params.Region,
		(*params.Subnet)[0].InitialSpec)

	// Create a subnet
	setCreatedRegionalNetworkResourceMetadata(subnetResponse.Metadata)
	subnetResponse.Status = secalib.NewSubnetStatus(secalib.CreatingResourceState)
	subnetResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "CreateSubnet", nextState: "GetSubnet"}); err != nil {
		return nil, err
	}

	// Get the created subnet
	secalib.SetSubnetStatusState(subnetResponse.Status, secalib.ActiveResourceState)
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "GetSubnet", nextState: "UpdateSubnet"}); err != nil {
		return nil, err
	}

	// Update the subnet
	setModifiedRegionalNetworkResourceMetadata(subnetResponse.Metadata)
	secalib.SetSubnetStatusState(subnetResponse.Status, secalib.UpdatingResourceState)
	subnetResponse.Spec = *(*params.Subnet)[0].UpdatedSpec
	subnetResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "UpdateSubnet", nextState: "GetSubnetUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated subnet
	secalib.SetSubnetStatusState(subnetResponse.Status, secalib.ActiveResourceState)
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "GetSubnetUpdated", nextState: "CreatePublicIp"}); err != nil {
		return nil, err
	}

	// Public ip
	publicIpResponse := newPublicIpResponse((*params.PublicIp)[0].Name, secalib.NetworkProviderV1, publicIpResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		(*params.PublicIp)[0].InitialSpec)

	// Create a public ip
	setCreatedRegionalWorkspaceResourceMetadata(publicIpResponse.Metadata)
	publicIpResponse.Status = secalib.NewPublicIpStatus(secalib.CreatingResourceState)
	publicIpResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "CreatePublicIp", nextState: "GetPublicIp"}); err != nil {
		return nil, err
	}

	// Get the created public ip
	secalib.SetPublicIpStatusState(publicIpResponse.Status, secalib.ActiveResourceState)
	publicIpResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "GetPublicIp", nextState: "UpdatePublicIp"}); err != nil {
		return nil, err
	}

	// Update the public ip
	setModifiedRegionalWorkspaceResourceMetadata(publicIpResponse.Metadata)
	secalib.SetPublicIpStatusState(publicIpResponse.Status, secalib.UpdatingResourceState)
	publicIpResponse.Spec = *(*params.PublicIp)[0].UpdatedSpec
	publicIpResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "UpdatePublicIp", nextState: "GetPublicIpUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated public ip
	secalib.SetPublicIpStatusState(publicIpResponse.Status, secalib.ActiveResourceState)
	publicIpResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "GetPublicIpUpdated", nextState: "CreateNIC"}); err != nil {
		return nil, err
	}

	// Nic
	nicResponse := newNicResponse((*params.NIC)[0].Name, secalib.NetworkProviderV1, nicResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		(*params.NIC)[0].InitialSpec)

	// Create a nic
	setCreatedRegionalWorkspaceResourceMetadata(nicResponse.Metadata)
	nicResponse.Status = secalib.NewNicStatus(secalib.CreatingResourceState)
	nicResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "CreateNIC", nextState: "GetNIC"}); err != nil {
		return nil, err
	}

	// Get the created nic
	secalib.SetNicStatusState(nicResponse.Status, secalib.ActiveResourceState)
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "GetNIC", nextState: "UpdateNIC"}); err != nil {
		return nil, err
	}

	// Update the nic
	setModifiedRegionalWorkspaceResourceMetadata(nicResponse.Metadata)
	secalib.SetNicStatusState(nicResponse.Status, secalib.UpdatingResourceState)
	nicResponse.Spec = *(*params.NIC)[0].UpdatedSpec
	nicResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "UpdateNIC", nextState: "GetNICUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated nic
	secalib.SetNicStatusState(nicResponse.Status, secalib.ActiveResourceState)
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "GetNICUpdated", nextState: "CreateSecurityGroup"}); err != nil {
		return nil, err
	}

	// Security group
	groupResponse := newSecurityGroupResponse((*params.SecurityGroup)[0].Name, secalib.NetworkProviderV1, groupResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		(*params.SecurityGroup)[0].InitialSpec)

	// Create a security group
	setCreatedRegionalWorkspaceResourceMetadata(groupResponse.Metadata)
	groupResponse.Status = secalib.SewSecurityGroupStatus(secalib.CreatingResourceState)
	groupResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "CreateSecurityGroup", nextState: "GetSecurityGroup"}); err != nil {
		return nil, err
	}

	// Get the created security group
	secalib.SetSecurityGroupStatusState(groupResponse.Status, secalib.ActiveResourceState)
	groupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "GetSecurityGroup", nextState: "UpdateSecurityGroup"}); err != nil {
		return nil, err
	}

	// Update the security group
	setModifiedRegionalWorkspaceResourceMetadata(groupResponse.Metadata)
	secalib.SetSecurityGroupStatusState(groupResponse.Status, secalib.UpdatingResourceState)
	groupResponse.Spec = *(*params.SecurityGroup)[0].UpdatedSpec
	groupResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "UpdateSecurityGroup", nextState: "GetSecurityGroupUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated security group
	secalib.SetSecurityGroupStatusState(groupResponse.Status, secalib.ActiveResourceState)
	groupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "GetSecurityGroupUpdated", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse := newBlockStorageResponse(params.BlockStorage.Name, secalib.ComputeProviderV1, blockResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.BlockStorage.InitialSpec)

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = secalib.NewBlockStorageStatus(secalib.CreatingResourceState)
	blockResponse.Spec.SizeGB = params.BlockStorage.InitialSpec.SizeGB
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the created block storage
	secalib.SetBlockStorageStatusState(blockResponse.Status, secalib.ActiveResourceState)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "CreateInstance"}); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse := newInstanceResponse(params.Instance.Name, secalib.ComputeProviderV1, instanceResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.Instance.InitialSpec)

	// Create an instance
	setCreatedRegionalWorkspaceResourceMetadata(instanceResponse.Metadata)
	instanceResponse.Status = secalib.NewInstanceStatus(secalib.CreatingResourceState)
	instanceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, responseBody: instanceResponse, currentState: "CreateInstance", nextState: "GetCreatedInstance"}); err != nil {
		return nil, err
	}

	// Get the created instance
	secalib.SetInstanceStatusState(instanceResponse.Status, secalib.ActiveResourceState)
	instanceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, responseBody: instanceResponse, currentState: "GetCreatedInstance", nextState: "DeleteInstance"}); err != nil {
		return nil, err
	}

	// Delete the instance
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, currentState: "DeleteInstance", nextState: "GetDeletedInstance"}); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: instanceUrl, params: params, currentState: "GetDeletedInstance", nextState: "DeleteBlockStorage"}); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, currentState: "DeleteBlockStorage", nextState: "GetDeletedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: blockUrl, params: params, currentState: "GetDeletedBlockStorage", nextState: "DeleteSecurityGroup"}); err != nil {
		return nil, err
	}

	// Delete the security group
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, currentState: "DeleteSecurityGroup", nextState: "GetDeletedSecurityGroup"}); err != nil {
		return nil, err
	}

	// Get the deleted security group
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: groupUrl, params: params, currentState: "GetDeletedSecurityGroup", nextState: "DeleteNic"}); err != nil {
		return nil, err
	}

	// Delete the nic
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, currentState: "DeleteNic", nextState: "GetDeletedNic"}); err != nil {
		return nil, err
	}

	// Get the deleted nic
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: nicUrl, params: params, currentState: "GetDeletedNic", nextState: "DeletePublicIp"}); err != nil {
		return nil, err
	}

	// Delete the public ip
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, currentState: "DeletePublicIp", nextState: "GetDeletedPublicIp"}); err != nil {
		return nil, err
	}

	// Get the deleted public ip
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: publicIpUrl, params: params, currentState: "GetDeletedPublicIp", nextState: "DeleteSubnet"}); err != nil {
		return nil, err
	}

	// Delete the subnet
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, currentState: "DeleteSubnet", nextState: "GetDeletedSubnet"}); err != nil {
		return nil, err
	}

	// Get the deleted subnet
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: subnetUrl, params: params, currentState: "GetDeletedSubnet", nextState: "DeleteRouteTable"}); err != nil {
		return nil, err
	}

	// Delete the route table
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, currentState: "DeleteRouteTable", nextState: "GetDeletedRouteTable"}); err != nil {
		return nil, err
	}

	// Get the deleted route table
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: routeUrl, params: params, currentState: "GetDeletedRouteTable", nextState: "DeleteInternetGateway"}); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, currentState: "DeleteInternetGateway", nextState: "GetDeletedInternetGateway"}); err != nil {
		return nil, err
	}

	// Get the deleted internet gateway
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: gatewayUrl, params: params, currentState: "GetDeletedInternetGateway", nextState: "DeleteNetwork"}); err != nil {
		return nil, err
	}

	// Delete the network
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, currentState: "DeleteNetwork", nextState: "GetDeletedNetwork"}); err != nil {
		return nil, err
	}

	// Get the deleted network
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: networkUrl, params: params, currentState: "GetDeletedNetwork", nextState: "DeleteWorkspace"}); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, currentState: "DeleteWorkspace", nextState: "GetDeletedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configureGetStubWithStatus(wm, scenario, http.StatusNotFound,
		&stubConfig{url: workspaceUrl, params: params, currentState: "GetDeletedWorkspace", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}
