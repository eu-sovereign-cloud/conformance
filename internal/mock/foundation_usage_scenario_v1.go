package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/wiremock/go-wiremock"
)

func CreateFoundationUsageScenario(scenario string, params *FoundationUsageParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	roleUrl := secalib.GenerateRoleURL(params.Tenant, params.Role.Name)
	roleAssignUrl := secalib.GenerateRoleAssignmentURL(params.Tenant, params.RoleAssignment.Name)
	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	blockUrl := secalib.GenerateBlockStorageURL(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageUrl := secalib.GenerateImageURL(params.Tenant, params.Image.Name)
	instanceUrl := secalib.GenerateInstanceURL(params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkUrl := secalib.GenerateNetworkURL(params.Tenant, params.Workspace.Name, params.Network.Name)
	gatewayUrl := secalib.GenerateInternetGatewayURL(params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicUrl := secalib.GenerateNicURL(params.Tenant, params.Workspace.Name, params.NIC.Name)
	publicIpUrl := secalib.GeneratePublicIpURL(params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeUrl := secalib.GenerateRouteTableURL(params.Tenant, params.Workspace.Name, params.Network.Name, params.RouteTable.Name)
	subnetUrl := secalib.GenerateSubnetURL(params.Tenant, params.Workspace.Name, params.Network.Name, params.Subnet.Name)
	groupUrl := secalib.GenerateSecurityGroupURL(params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// GenerateResources
	roleResource := secalib.GenerateRoleResource(params.Tenant, params.Role.Name)
	roleAssignResource := secalib.GenerateRoleAssignmentResource(params.Tenant, params.RoleAssignment.Name)
	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)
	blockResource := secalib.GenerateBlockStorageResource(params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageResource := secalib.GenerateImageResource(params.Tenant, params.Image.Name)
	instanceResource := secalib.GenerateInstanceResource(params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkResource := secalib.GenerateNetworkResource(params.Tenant, params.Workspace.Name, params.Network.Name)
	gatewayResource := secalib.GenerateInternetGatewayResource(params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicResource := secalib.GenerateNicResource(params.Tenant, params.Workspace.Name, params.NIC.Name)
	publicIpResource := secalib.GeneratePublicIpResource(params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeResource := secalib.GenerateRouteTableResource(params.Tenant, params.Workspace.Name, params.Network.Name, params.RouteTable.Name)
	subnetResource := secalib.GenerateSubnetResource(params.Tenant, params.Workspace.Name, params.Network.Name, params.Subnet.Name)
	groupResource := secalib.GenerateSecurityGroupResource(params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Authorization

	// Role
	roleResponse := newRoleResponse(params.Role.Name, secalib.AuthorizationProviderV1, roleResource, secalib.ApiVersion1,
		params.Tenant,
		params.Role.InitialSpec)

	// Create Role
	setCreatedGlobalTenantResourceMetadata(roleResponse.Metadata)
	roleResponse.Status = newStatus(secalib.CreatingStatusState)
	roleResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, response: roleResponse, currentState: startedScenarioState, nextState: "GetCreatedRole"}); err != nil {
		return nil, err
	}

	// Get created role
	setStatusState(roleResponse.Status, secalib.ActiveStatusState)
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, response: roleResponse, currentState: "GetCreatedRole", nextState: "CreateRoleAssignment"}); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse := newRoleAssignmentResponse(params.RoleAssignment.Name, secalib.AuthorizationProviderV1, roleAssignResource, secalib.ApiVersion1,
		params.Tenant,
		params.RoleAssignment.InitialSpec)

	// Create a role assignment
	setCreatedGlobalTenantResourceMetadata(roleAssignResponse.Metadata)
	roleAssignResponse.Status = newStatus(secalib.CreatingStatusState)
	roleAssignResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, response: roleAssignResponse, currentState: "CreateRoleAssignment", nextState: "GetCreatedRoleAssignment"}); err != nil {
		return nil, err
	}

	// Get created role assignment
	setStatusState(roleAssignResponse.Status, secalib.ActiveStatusState)
	roleAssignResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, response: roleAssignResponse, currentState: "GetCreatedRoleAssignment", nextState: "CreateWorkspace"}); err != nil {
		return nil, err
	}

	// Workspace
	workspaceResponse := newWorkspaceResponse(params.Workspace.Name, secalib.WorkspaceProviderV1, workspaceResource, secalib.ApiVersion1,
		params.Tenant, params.Region,
		params.Workspace.InitialLabels)

	// Create a workspace
	setCreatedRegionalResourceMetadata(workspaceResponse.Metadata)
	workspaceResponse.Status = newWorkspaceStatus(secalib.CreatingStatusState)
	workspaceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, response: workspaceResponse, currentState: "CreateWorkspace", nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get created workspace
	setWorkspaceStatusState(workspaceResponse.Status, secalib.ActiveStatusState)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, response: workspaceResponse, currentState: "GetCreatedWorkspace", nextState: "CreateImage"}); err != nil {
		return nil, err
	}

	// Storage

	// Image
	imageResponse := newImageResponse(params.Image.Name, secalib.ComputeProviderV1, imageResource, secalib.ApiVersion1,
		params.Tenant, params.Region,
		params.Image.InitialSpec)

	// Create an image
	setCreatedRegionalResourceMetadata(imageResponse.Metadata)
	imageResponse.Status = newImageStatus(secalib.CreatingStatusState)
	imageResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, response: imageResponse, currentState: "CreateImage", nextState: "GetCreatedImage"}); err != nil {
		return nil, err
	}

	// Get created image
	setImageStatusState(imageResponse.Status, secalib.ActiveStatusState)
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, response: imageResponse, currentState: "GetCreatedImage", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	blockResponse := newBlockStorageResponse(params.BlockStorage.Name, secalib.ComputeProviderV1, blockResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.BlockStorage.InitialSpec)

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = newBlockStorageStatus(secalib.CreatingStatusState)
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, response: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get created block storage
	setBlockStorageStatusState(blockResponse.Status, secalib.ActiveStatusState)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, response: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "CreateNetwork"}); err != nil {
		return nil, err
	}

	// Network

	// Network
	networkResponse := newNetworkResponse(params.Network.Name, secalib.NetworkProviderV1, networkResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.Network.InitialSpec)

	// Create  Network
	setCreatedRegionalWorkspaceResourceMetadata(networkResponse.Metadata)
	networkResponse.Status = newNetworkStatus(secalib.CreatingStatusState)
	networkResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, response: networkResponse, currentState: "CreateNetwork", nextState: "GetNetwork"}); err != nil {
		return nil, err
	}

	// Get network
	setNetworkStatusState(networkResponse.Status, secalib.ActiveStatusState)
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, response: networkResponse, currentState: "GetNetwork", nextState: "CreateInternetGateway"}); err != nil {
		return nil, err
	}

	// Internet gateway
	gatewayResponse := newInternetGatewayResponse(params.InternetGateway.Name, secalib.NetworkProviderV1, gatewayResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.InternetGateway.InitialSpec)

	// Create internet gateway
	setCreatedRegionalWorkspaceResourceMetadata(gatewayResponse.Metadata)
	gatewayResponse.Status = newStatus(secalib.CreatingStatusState)
	gatewayResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, response: gatewayResponse, currentState: "CreateInternetGateway", nextState: "GetInternetGateway"}); err != nil {
		return nil, err
	}

	// Get internet-gateway
	setStatusState(gatewayResponse.Status, secalib.ActiveStatusState)
	gatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, response: gatewayResponse, currentState: "GetInternetGateway", nextState: "CreateRouteTable"}); err != nil {
		return nil, err
	}

	// Route table
	routeResponse := newRouteTableResponse(params.RouteTable.Name, secalib.NetworkProviderV1, routeResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Network.Name, params.Region,
		params.RouteTable.InitialSpec)

	// Create route-table
	setCreatedRegionalNetworkResourceMetadata(routeResponse.Metadata)
	routeResponse.Status = newRouteTableStatus(secalib.CreatingStatusState)
	routeResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, response: routeResponse, currentState: "CreateRouteTable", nextState: "GetRouteTable"}); err != nil {
		return nil, err
	}

	// Get route-table
	setRouteTableStatusState(routeResponse.Status, secalib.ActiveStatusState)
	routeResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, response: routeResponse, currentState: "GetRouteTable", nextState: "CreateSubnet"}); err != nil {
		return nil, err
	}

	// Subnet
	subnetResponse := newSubnetResponse(params.Subnet.Name, secalib.NetworkProviderV1, subnetResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Network.Name, params.Region,
		params.Subnet.InitialSpec)

	// Create subnet
	setCreatedRegionalNetworkResourceMetadata(subnetResponse.Metadata)
	subnetResponse.Status = newSubnetStatus(secalib.CreatingStatusState)
	subnetResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, response: subnetResponse, currentState: "CreateSubnet", nextState: "GetSubnet"}); err != nil {
		return nil, err
	}

	// Get subnet
	setSubnetStatusState(subnetResponse.Status, secalib.ActiveStatusState)
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, response: subnetResponse, currentState: "GetSubnet", nextState: "CreateSecurityGroup"}); err != nil {
		return nil, err
	}

	// Security group
	groupResponse := newSecurityGroupResponse(params.SecurityGroup.Name, secalib.NetworkProviderV1, secalib.ApiVersion1, groupResource,
		params.Tenant, params.Workspace.Name, params.Region,
		params.SecurityGroup.InitialSpec)

	// Create security-group
	setCreatedRegionalWorkspaceResourceMetadata(groupResponse.Metadata)
	groupResponse.Status = newSecurityGroupStatus(secalib.CreatingStatusState)
	groupResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, response: groupResponse, currentState: "CreateSecurityGroup", nextState: "GetSecurityGroup"}); err != nil {
		return nil, err
	}

	// Get security-group
	setSecurityGroupStatusState(groupResponse.Status, secalib.ActiveStatusState)
	groupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, response: groupResponse, currentState: "GetSecurityGroup", nextState: "CreatePublicIp"}); err != nil {
		return nil, err
	}

	// Public-ip
	publicIpResponse := newPublicIpResponse(params.PublicIp.Name, secalib.NetworkProviderV1, publicIpResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.PublicIp.InitialSpec)

	// Create public-ip
	setCreatedRegionalWorkspaceResourceMetadata(publicIpResponse.Metadata)
	publicIpResponse.Status = newPublicIpStatus(secalib.CreatingStatusState)
	publicIpResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, response: publicIpResponse, currentState: "CreatePublicIp", nextState: "GetPublicIp"}); err != nil {
		return nil, err
	}

	// Get public-ip
	setPublicIpStatusState(publicIpResponse.Status, secalib.ActiveStatusState)
	publicIpResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, response: publicIpResponse, currentState: "GetPublicIp", nextState: "CreateNIC"}); err != nil {
		return nil, err
	}

	// NIC
	nicResponse := newNicResponse(params.NIC.Name, secalib.NetworkProviderV1, nicResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.NIC.InitialSpec)

	// Create NIC
	setCreatedRegionalWorkspaceResourceMetadata(nicResponse.Metadata)
	nicResponse.Status = newNicStatus(secalib.CreatingStatusState)
	nicResponse.Metadata.Verb = http.MethodPut
	if err := configurePutSuccessStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, response: nicResponse, currentState: "CreateNIC", nextState: "GetNIC"}); err != nil {
		return nil, err
	}

	// Get NIC
	setNicStatusState(nicResponse.Status, secalib.ActiveStatusState)
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetSuccessStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, response: nicResponse, currentState: "GetNIC", nextState: "CreateInstance"}); err != nil {
		return nil, err
	}

	// Compute

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

	// Delete all

	// Delete instance
	instanceResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, response: instanceResponse, currentState: "DeleteInstance", nextState: "DeleteSecurityGroup"}); err != nil {
		return nil, err
	}

	// Delete Security Group
	groupResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, response: groupResponse, currentState: "DeleteSecurityGroup", nextState: "DeleteNic"}); err != nil {
		return nil, err
	}

	// Delete nic
	nicResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, response: nicResponse, currentState: "DeleteNic", nextState: "DeletePublicIp"}); err != nil {
		return nil, err
	}

	// Delete public ip
	publicIpResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, response: publicIpResponse, currentState: "DeletePublicIp", nextState: "DeleteSubnet"}); err != nil {
		return nil, err
	}

	// Delete subnet
	subnetResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, response: subnetResponse, currentState: "DeleteSubnet", nextState: "DeleteRouteTable"}); err != nil {
		return nil, err
	}

	// Delete route-table
	routeResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, response: routeResponse, currentState: "DeleteRouteTable", nextState: "DeleteInternetGateway"}); err != nil {
		return nil, err
	}

	// Delete internet-gateway
	gatewayResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, response: gatewayResponse, currentState: "DeleteInternetGateway", nextState: "DeleteNetwork"}); err != nil {
		return nil, err
	}

	// Delete network
	networkResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, response: networkResponse, currentState: "DeleteNetwork", nextState: "DeleteBlockStorage"}); err != nil {
		return nil, err
	}

	// Delete block-storage
	blockResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, response: blockResponse, currentState: "DeleteBlockStorage", nextState: "DeleteImage"}); err != nil {
		return nil, err
	}

	// Delete image
	imageResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, response: imageResponse, currentState: "DeleteImage", nextState: "DeleteWorkspace"}); err != nil {
		return nil, err
	}

	// Delete workspace
	workspaceResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, response: workspaceResponse, currentState: "DeleteWorkspace", nextState: "DeleteRoleAssignment"}); err != nil {
		return nil, err
	}

	// Delete role assignment
	roleAssignResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, response: roleAssignResponse, currentState: "DeleteRoleAssignment", nextState: "DeleteRole"}); err != nil {
		return nil, err
	}

	// Delete role
	roleResponse.Metadata.Verb = http.MethodDelete
	if err := configureDeleteSuccessStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, response: roleResponse, currentState: "DeleteRole"}); err != nil {
		return nil, err
	}

	return wm, err
}
