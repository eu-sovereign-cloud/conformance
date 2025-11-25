package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/conformance/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
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
	roleResponse, err := builders.NewRoleBuilder().
		Name(params.Role.Name).
		Provider(secalib.AuthorizationProviderV1).
		Resource(roleResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Spec(params.Role.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create Role
	setCreatedGlobalTenantResourceMetadata(roleResponse.Metadata)
	roleResponse.Status = secalib.NewResourceStatus(schema.ResourceStateCreating)
	roleResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: startedScenarioState, nextState: "GetCreatedRole"}); err != nil {
		return nil, err
	}

	// Get the created role
	secalib.SetStatusState(roleResponse.Status, schema.ResourceStateActive)
	roleResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: roleUrl, params: params, responseBody: roleResponse, currentState: "GetCreatedRole", nextState: "CreateRoleAssignment"}); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name(params.RoleAssignment.Name).
		Provider(secalib.AuthorizationProviderV1).
		Resource(roleAssignResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Spec(params.RoleAssignment.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a role assignment
	setCreatedGlobalTenantResourceMetadata(roleAssignResponse.Metadata)
	roleAssignResponse.Status = secalib.NewResourceStatus(schema.ResourceStateCreating)
	roleAssignResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: "CreateRoleAssignment", nextState: "GetCreatedRoleAssignment"}); err != nil {
		return nil, err
	}

	// Get the created role assignment
	secalib.SetStatusState(roleAssignResponse.Status, schema.ResourceStateActive)
	roleAssignResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: roleAssignUrl, params: params, responseBody: roleAssignResponse, currentState: "GetCreatedRoleAssignment", nextState: "CreateWorkspace"}); err != nil {
		return nil, err
	}

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(secalib.WorkspaceProviderV1).
		Resource(workspaceResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	setCreatedRegionalResourceMetadata(workspaceResponse.Metadata)
	workspaceResponse.Status = secalib.NewWorkspaceStatus(schema.ResourceStateCreating)
	workspaceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: "CreateWorkspace", nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	secalib.SetWorkspaceStatusState(workspaceResponse.Status, schema.ResourceStateActive)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: "GetCreatedWorkspace", nextState: "CreateImage"}); err != nil {
		return nil, err
	}

	// Storage

	// Image
	imageResponse, err := builders.NewImageBuilder().
		Name(params.Image.Name).
		Provider(secalib.StorageProviderV1).
		Resource(imageResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Region(params.Region).
		Spec(params.Image.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create an image
	setCreatedRegionalResourceMetadata(imageResponse.Metadata)
	imageResponse.Status = secalib.NewImageStatus(schema.ResourceStateCreating)
	imageResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "CreateImage", nextState: "GetCreatedImage"}); err != nil {
		return nil, err
	}

	// Get the created image
	secalib.SetImageStatusState(imageResponse.Status, schema.ResourceStateActive)
	imageResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: imageUrl, params: params, responseBody: imageResponse, currentState: "GetCreatedImage", nextState: "CreateBlockStorage"}); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(secalib.StorageProviderV1).
		Resource(blockResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	setCreatedRegionalWorkspaceResourceMetadata(blockResponse.Metadata)
	blockResponse.Status = secalib.NewBlockStorageStatus(schema.ResourceStateCreating)
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the created block storage
	secalib.SetBlockStorageStatusState(blockResponse.Status, schema.ResourceStateActive)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "CreateNetwork"}); err != nil {
		return nil, err
	}

	// Network

	// Network
	networkResponse, err := builders.NewNetworkBuilder().
		Name(params.Network.Name).
		Provider(secalib.NetworkProviderV1).
		Resource(networkResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec(params.Network.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create  Network
	setCreatedRegionalWorkspaceResourceMetadata(networkResponse.Metadata)
	networkResponse.Status = secalib.NewNetworkStatus(schema.ResourceStateCreating)
	networkResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "CreateNetwork", nextState: "GetNetwork"}); err != nil {
		return nil, err
	}

	// Get network
	secalib.SetNetworkStatusState(networkResponse.Status, schema.ResourceStateActive)
	networkResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "GetNetwork", nextState: "CreateInternetGateway"}); err != nil {
		return nil, err
	}

	// Internet gateway
	gatewayResponse, err := builders.NewInternetGatewayBuilder().
		Name(params.InternetGateway.Name).
		Provider(secalib.NetworkProviderV1).
		Resource(gatewayResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec(params.InternetGateway.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create internet gateway
	setCreatedRegionalWorkspaceResourceMetadata(gatewayResponse.Metadata)
	gatewayResponse.Status = secalib.NewResourceStatus(schema.ResourceStateCreating)
	gatewayResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "CreateInternetGateway", nextState: "GetInternetGateway"}); err != nil {
		return nil, err
	}

	// Get internet-gateway
	secalib.SetStatusState(gatewayResponse.Status, schema.ResourceStateActive)
	gatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "GetInternetGateway", nextState: "CreateRouteTable"}); err != nil {
		return nil, err
	}

	// Route table
	routeResponse, err := builders.NewRouteTableBuilder().
		Name(params.RouteTable.Name).
		Provider(secalib.NetworkProviderV1).
		Resource(routeResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Network(params.Network.Name).
		Region(params.Region).
		Spec(params.RouteTable.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create route-table
	setCreatedRegionalNetworkResourceMetadata(routeResponse.Metadata)
	routeResponse.Status = secalib.NewRouteTableStatus(schema.ResourceStateCreating)
	routeResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "CreateRouteTable", nextState: "GetRouteTable"}); err != nil {
		return nil, err
	}

	// Get route-table
	secalib.SetRouteTableStatusState(routeResponse.Status, schema.ResourceStateActive)
	routeResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "GetRouteTable", nextState: "CreateSubnet"}); err != nil {
		return nil, err
	}

	// Subnet
	subnetResponse, err := builders.NewSubnetBuilder().
		Name(params.Subnet.Name).
		Provider(secalib.NetworkProviderV1).
		Resource(subnetResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Network(params.Network.Name).
		Region(params.Region).
		Spec(params.Subnet.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create subnet
	setCreatedRegionalNetworkResourceMetadata(subnetResponse.Metadata)
	subnetResponse.Status = secalib.NewSubnetStatus(schema.ResourceStateCreating)
	subnetResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "CreateSubnet", nextState: "GetSubnet"}); err != nil {
		return nil, err
	}

	// Get subnet
	secalib.SetSubnetStatusState(subnetResponse.Status, schema.ResourceStateActive)
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "GetSubnet", nextState: "CreateSecurityGroup"}); err != nil {
		return nil, err
	}

	// Security group
	groupResponse, err := builders.NewSecurityGroupBuilder().
		Name(params.SecurityGroup.Name).
		Provider(secalib.NetworkProviderV1).
		Resource(groupResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec(params.SecurityGroup.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create security-group
	setCreatedRegionalWorkspaceResourceMetadata(groupResponse.Metadata)
	groupResponse.Status = secalib.NewSecurityGroupStatus(schema.ResourceStateCreating)
	groupResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "CreateSecurityGroup", nextState: "GetSecurityGroup"}); err != nil {
		return nil, err
	}

	// Get security-group
	secalib.SetSecurityGroupStatusState(groupResponse.Status, schema.ResourceStateActive)
	groupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "GetSecurityGroup", nextState: "CreatePublicIp"}); err != nil {
		return nil, err
	}

	// Public-ip
	publicIpResponse, err := builders.NewPublicIpBuilder().
		Name(params.PublicIp.Name).
		Provider(secalib.NetworkProviderV1).
		Resource(publicIpResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec(params.PublicIp.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create public-ip
	setCreatedRegionalWorkspaceResourceMetadata(publicIpResponse.Metadata)
	publicIpResponse.Status = secalib.NewPublicIpStatus(schema.ResourceStateCreating)
	publicIpResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "CreatePublicIp", nextState: "GetPublicIp"}); err != nil {
		return nil, err
	}

	// Get public-ip
	secalib.SetPublicIpStatusState(publicIpResponse.Status, schema.ResourceStateActive)
	publicIpResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "GetPublicIp", nextState: "CreateNIC"}); err != nil {
		return nil, err
	}

	// NIC
	nicResponse, err := builders.NewNicBuilder().
		Name(params.Nic.Name).
		Provider(secalib.NetworkProviderV1).
		Resource(nicResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec(params.Nic.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create NIC
	setCreatedRegionalWorkspaceResourceMetadata(nicResponse.Metadata)
	nicResponse.Status = secalib.NewNicStatus(schema.ResourceStateCreating)
	nicResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "CreateNIC", nextState: "GetNIC"}); err != nil {
		return nil, err
	}

	// Get NIC
	secalib.SetNicStatusState(nicResponse.Status, schema.ResourceStateActive)
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "GetNIC", nextState: "CreateInstance"}); err != nil {
		return nil, err
	}

	// Compute

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(params.Instance.Name).
		Provider(secalib.ComputeProviderV1).
		Resource(instanceResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec(params.Instance.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create an instance
	setCreatedRegionalWorkspaceResourceMetadata(instanceResponse.Metadata)
	instanceResponse.Status = secalib.NewInstanceStatus(schema.ResourceStateCreating)
	instanceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, responseBody: instanceResponse, currentState: "CreateInstance", nextState: "GetCreatedInstance"}); err != nil {
		return nil, err
	}

	// Get the created instance
	secalib.SetInstanceStatusState(instanceResponse.Status, schema.ResourceStateActive)
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
