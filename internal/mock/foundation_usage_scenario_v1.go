package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/wiremock/go-wiremock"
)

func ConfigFoundationUsageScenario(scenario string, params *FoundationUsageParamsV1) (*wiremock.Client, error) {
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
	nicUrl := secalib.GenerateNicURL(params.Tenant, params.Workspace.Name, params.Nic.Name)
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
	nicResource := secalib.GenerateNicResource(params.Tenant, params.Workspace.Name, params.Nic.Name)
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
	roleResponse.Status = secalib.NewResourceStatus(secalib.CreatingResourceState)
	roleResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: startedScenarioState, nextState: "GetCreatedRole"}); err != nil {
		return nil, err
	}

	// Get the created role
	secalib.SetStatusState(roleResponse.Status, secalib.ActiveResourceState)
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: "GetCreatedRole", nextState: "CreateRoleAssignment"}); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse := newRoleAssignmentResponse(params.RoleAssignment.Name, secalib.AuthorizationProviderV1, roleAssignResource, secalib.ApiVersion1,
		params.Tenant,
		params.RoleAssignment.InitialSpec)

	// Create a role assignment
	setCreatedGlobalTenantResourceMetadata(roleAssignResponse.Metadata)
	roleAssignResponse.Status = secalib.NewResourceStatus(secalib.CreatingResourceState)
	roleAssignResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: "CreateRoleAssignment", nextState: "GetCreatedRoleAssignment"}); err != nil {
		return nil, err
	}

	// Get the created role assignment
	secalib.SetStatusState(roleAssignResponse.Status, secalib.ActiveResourceState)
	roleAssignResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: "GetCreatedRoleAssignment", nextState: "CreateWorkspace"}); err != nil {
		return nil, err
	}

	// Workspace
	workspaceResponse := newWorkspaceResponse(params.Workspace.Name, secalib.WorkspaceProviderV1, workspaceResource, secalib.ApiVersion1,
		params.Tenant, params.Region,
		params.Workspace.InitialLabels)

	// Create a workspace
	setCreatedRegionalResourceMetadata(workspaceResponse.Metadata)
	workspaceResponse.Status = secalib.NewWorkspaceStatus(secalib.CreatingResourceState)
	workspaceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: "CreateWorkspace", nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	secalib.SetWorkspaceStatusState(workspaceResponse.Status, secalib.ActiveResourceState)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: "GetCreatedWorkspace", nextState: "CreateImage"}); err != nil {
		return nil, err
	}

	// Storage

	// Image
	imageResponse := newImageResponse(params.Image.Name, secalib.StorageProviderV1, imageResource, secalib.ApiVersion1,
		params.Tenant, params.Region,
		&params.Role.InitialLabels,
		params.Image.InitialSpec)

	// Create an image
	setCreatedRegionalResourceMetadata(imageResponse.Metadata)
	imageResponse.Status = secalib.NewImageStatus(secalib.CreatingResourceState)
	imageResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "CreateImage", nextState: "GetCreatedImage"}); err != nil {
		return nil, err
	}

	// Get the created image
	secalib.SetImageStatusState(imageResponse.Status, secalib.ActiveResourceState)
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "GetCreatedImage", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	blockResponse := newBlockStorageResponse(params.BlockStorage.Name, secalib.StorageProviderV1, blockResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.BlockStorage.InitialLabels,
		params.BlockStorage.InitialSpec)

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = secalib.NewBlockStorageStatus(secalib.CreatingResourceState)
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the created block storage
	secalib.SetBlockStorageStatusState(blockResponse.Status, secalib.ActiveResourceState)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "CreateNetwork"}); err != nil {
		return nil, err
	}

	// Network

	// Network
	networkResponse := newNetworkResponse(params.Network.Name, secalib.NetworkProviderV1, networkResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		&params.Network.InitialLabels,
		params.Network.InitialSpec)

	// Create  Network
	setCreatedRegionalWorkspaceResourceMetadata(networkResponse.Metadata)
	networkResponse.Status = secalib.NewNetworkStatus(secalib.CreatingResourceState)
	networkResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "CreateNetwork", nextState: "GetNetwork"}); err != nil {
		return nil, err
	}

	// Get network
	secalib.SetNetworkStatusState(networkResponse.Status, secalib.ActiveResourceState)
	networkResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "GetNetwork", nextState: "CreateInternetGateway"}); err != nil {
		return nil, err
	}

	// Internet gateway
	gatewayResponse := newInternetGatewayResponse(params.InternetGateway.Name, secalib.NetworkProviderV1, gatewayResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.InternetGateway.InitialSpec)

	// Create internet gateway
	setCreatedRegionalWorkspaceResourceMetadata(gatewayResponse.Metadata)
	gatewayResponse.Status = secalib.NewResourceStatus(secalib.CreatingResourceState)
	gatewayResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "CreateInternetGateway", nextState: "GetInternetGateway"}); err != nil {
		return nil, err
	}

	// Get internet-gateway
	secalib.SetStatusState(gatewayResponse.Status, secalib.ActiveResourceState)
	gatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "GetInternetGateway", nextState: "CreateRouteTable"}); err != nil {
		return nil, err
	}

	// Route table
	routeResponse := newRouteTableResponse(params.RouteTable.Name, secalib.NetworkProviderV1, routeResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Network.Name, params.Region,
		params.RouteTable.InitialSpec)

	// Create route-table
	setCreatedRegionalNetworkResourceMetadata(routeResponse.Metadata)
	routeResponse.Status = secalib.NewRouteTableStatus(secalib.CreatingResourceState)
	routeResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "CreateRouteTable", nextState: "GetRouteTable"}); err != nil {
		return nil, err
	}

	// Get route-table
	secalib.SetRouteTableStatusState(routeResponse.Status, secalib.ActiveResourceState)
	routeResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "GetRouteTable", nextState: "CreateSubnet"}); err != nil {
		return nil, err
	}

	// Subnet
	subnetResponse := newSubnetResponse(params.Subnet.Name, secalib.NetworkProviderV1, subnetResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Network.Name, params.Region,
		params.Subnet.InitialSpec)

	// Create subnet
	setCreatedRegionalNetworkResourceMetadata(subnetResponse.Metadata)
	subnetResponse.Status = secalib.NewSubnetStatus(secalib.CreatingResourceState)
	subnetResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "CreateSubnet", nextState: "GetSubnet"}); err != nil {
		return nil, err
	}

	// Get subnet
	secalib.SetSubnetStatusState(subnetResponse.Status, secalib.ActiveResourceState)
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "GetSubnet", nextState: "CreateSecurityGroup"}); err != nil {
		return nil, err
	}

	// Security group
	groupResponse := newSecurityGroupResponse(params.SecurityGroup.Name, secalib.NetworkProviderV1, groupResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.SecurityGroup.InitialSpec)

	// Create security-group
	setCreatedRegionalWorkspaceResourceMetadata(groupResponse.Metadata)
	groupResponse.Status = secalib.NewSecurityGroupStatus(secalib.CreatingResourceState)
	groupResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "CreateSecurityGroup", nextState: "GetSecurityGroup"}); err != nil {
		return nil, err
	}

	// Get security-group
	secalib.SetSecurityGroupStatusState(groupResponse.Status, secalib.ActiveResourceState)
	groupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "GetSecurityGroup", nextState: "CreatePublicIp"}); err != nil {
		return nil, err
	}

	// Public-ip
	publicIpResponse := newPublicIpResponse(params.PublicIp.Name, secalib.NetworkProviderV1, publicIpResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.PublicIp.InitialSpec)

	// Create public-ip
	setCreatedRegionalWorkspaceResourceMetadata(publicIpResponse.Metadata)
	publicIpResponse.Status = secalib.NewPublicIpStatus(secalib.CreatingResourceState)
	publicIpResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "CreatePublicIp", nextState: "GetPublicIp"}); err != nil {
		return nil, err
	}

	// Get public-ip
	secalib.SetPublicIpStatusState(publicIpResponse.Status, secalib.ActiveResourceState)
	publicIpResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "GetPublicIp", nextState: "CreateNIC"}); err != nil {
		return nil, err
	}

	// NIC
	nicResponse := newNicResponse(params.Nic.Name, secalib.NetworkProviderV1, nicResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		params.Nic.InitialSpec)

	// Create NIC
	setCreatedRegionalWorkspaceResourceMetadata(nicResponse.Metadata)
	nicResponse.Status = secalib.NewNicStatus(secalib.CreatingResourceState)
	nicResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "CreateNIC", nextState: "GetNIC"}); err != nil {
		return nil, err
	}

	// Get NIC
	secalib.SetNicStatusState(nicResponse.Status, secalib.ActiveResourceState)
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "GetNIC", nextState: "CreateInstance"}); err != nil {
		return nil, err
	}

	// Compute

	// Instance
	instanceResponse := newInstanceResponse(params.Instance.Name, secalib.ComputeProviderV1, instanceResource, secalib.ApiVersion1,
		params.Tenant, params.Workspace.Name, params.Region,
		&params.Instance.InitialLabels,
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

	// Delete all

	// Delete instance
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, currentState: "DeleteInstance", nextState: "DeleteSecurityGroup"}); err != nil {
		return nil, err
	}

	// Delete Security Group
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, currentState: "DeleteSecurityGroup", nextState: "DeleteNic"}); err != nil {
		return nil, err
	}

	// Delete nic
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, currentState: "DeleteNic", nextState: "DeletePublicIp"}); err != nil {
		return nil, err
	}

	// Delete public ip
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, currentState: "DeletePublicIp", nextState: "DeleteSubnet"}); err != nil {
		return nil, err
	}

	// Delete subnet
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, currentState: "DeleteSubnet", nextState: "DeleteRouteTable"}); err != nil {
		return nil, err
	}

	// Delete route-table
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, currentState: "DeleteRouteTable", nextState: "DeleteInternetGateway"}); err != nil {
		return nil, err
	}

	// Delete internet-gateway
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, currentState: "DeleteInternetGateway", nextState: "DeleteNetwork"}); err != nil {
		return nil, err
	}

	// Delete network
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, currentState: "DeleteNetwork", nextState: "DeleteBlockStorage"}); err != nil {
		return nil, err
	}

	// Delete block-storage
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, currentState: "DeleteBlockStorage", nextState: "DeleteImage"}); err != nil {
		return nil, err
	}

	// Delete image
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, currentState: "DeleteImage", nextState: "DeleteWorkspace"}); err != nil {
		return nil, err
	}

	// Delete workspace
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, currentState: "DeleteWorkspace", nextState: "DeleteRoleAssignment"}); err != nil {
		return nil, err
	}

	// Delete role assignment
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, currentState: "DeleteRoleAssignment", nextState: "DeleteRole"}); err != nil {
		return nil, err
	}

	// Delete role
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, currentState: "DeleteRole"}); err != nil {
		return nil, err
	}

	return wm, err
}
