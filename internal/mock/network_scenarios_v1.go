package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	networkV1 "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
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
	workspaceResponse.Status = newWorkspaceStatus(schema.ResourceStateCreating)
	workspaceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: startedScenarioState, nextState: "GetCreatedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the created workspace
	setWorkspaceState(workspaceResponse.Status, schema.ResourceStateActive)
	workspaceResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workspaceResponse, currentState: "GetCreatedWorkspace", nextState: "CreateNetwork"}); err != nil {
		return nil, err
	}

	// Network
	networkResponse, err := builders.NewNetworkBuilder().
		Name((*params.Network)[0].Name).
		Provider(secalib.NetworkProviderV1).
		Resource(networkResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec((*params.Network)[0].InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a network
	setCreatedRegionalWorkspaceResourceMetadata(networkResponse.Metadata)
	networkResponse.Status = newNetworkStatus(schema.ResourceStateCreating)
	networkResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "CreateNetwork", nextState: "GetNetwork"}); err != nil {
		return nil, err
	}

	// Get the created network
	setNetworkState(networkResponse.Status, schema.ResourceStateActive)
	networkResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "GetNetwork", nextState: "UpdateNetwork"}); err != nil {
		return nil, err
	}

	// Update the network
	setNetworkState(networkResponse.Status, schema.ResourceStateUpdating)
	networkResponse.Spec = *(*params.Network)[0].UpdatedSpec
	networkResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "UpdateNetwork", nextState: "GetNetwork2x"}); err != nil {
		return nil, err
	}

	// Get the updated network
	setNetworkState(networkResponse.Status, schema.ResourceStateActive)
	networkResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, responseBody: networkResponse, currentState: "GetNetwork2x", nextState: "CreateInternetGateway"}); err != nil {
		return nil, err
	}

	// Internet gateway
	gatewayResponse, err := builders.NewInternetGatewayBuilder().
		Name((*params.InternetGateway)[0].Name).
		Provider(secalib.NetworkProviderV1).
		Resource(gatewayResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec((*params.InternetGateway)[0].InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create an internet gateway
	setCreatedRegionalWorkspaceResourceMetadata(gatewayResponse.Metadata)
	gatewayResponse.Status = newResourceStatus(schema.ResourceStateCreating)
	gatewayResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "CreateInternetGateway", nextState: "GetInternetGateway"}); err != nil {
		return nil, err
	}

	// Get the created internet gateway
	setResourceState(gatewayResponse.Status, schema.ResourceStateActive)
	gatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "GetInternetGateway", nextState: "UpdateInternetGateway"}); err != nil {
		return nil, err
	}

	// Update the internet gateway
	setModifiedRegionalWorkspaceResourceMetadata(gatewayResponse.Metadata)
	setResourceState(gatewayResponse.Status, schema.ResourceStateUpdating)
	gatewayResponse.Spec = *(*params.InternetGateway)[0].UpdatedSpec
	gatewayResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "UpdateInternetGateway", nextState: "GetInternetGateway2x"}); err != nil {
		return nil, err
	}

	// Get the updated internet gateway
	setResourceState(gatewayResponse.Status, schema.ResourceStateActive)
	gatewayResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, responseBody: gatewayResponse, currentState: "GetInternetGateway2x", nextState: "CreateRouteTable"}); err != nil {
		return nil, err
	}

	// Route table
	routeResponse, err := builders.NewRouteTableBuilder().
		Name((*params.RouteTable)[0].Name).
		Provider(secalib.NetworkProviderV1).
		Resource(routeResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Network((*params.Network)[0].Name).
		Region(params.Region).
		Spec((*params.RouteTable)[0].InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a route table
	setCreatedRegionalNetworkResourceMetadata(routeResponse.Metadata)
	routeResponse.Status = newRouteTableStatus(schema.ResourceStateCreating)
	routeResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "CreateRouteTable", nextState: "GetRouteTable"}); err != nil {
		return nil, err
	}

	// Get the created route table
	setRouteTableState(routeResponse.Status, schema.ResourceStateActive)
	routeResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "GetRouteTable", nextState: "UpdateRouteTable"}); err != nil {
		return nil, err
	}

	// Update the route table
	setModifiedRegionalNetworkResourceMetadata(routeResponse.Metadata)
	setRouteTableState(routeResponse.Status, schema.ResourceStateUpdating)
	routeResponse.Spec = *(*params.RouteTable)[0].UpdatedSpec
	routeResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "UpdateRouteTable", nextState: "GetRouteTableUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated route table
	setRouteTableState(routeResponse.Status, schema.ResourceStateActive)
	routeResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, responseBody: routeResponse, currentState: "GetRouteTableUpdated", nextState: "CreateSubnet"}); err != nil {
		return nil, err
	}

	// Subnet
	subnetResponse, err := builders.NewSubnetBuilder().
		Name((*params.Subnet)[0].Name).
		Provider(secalib.NetworkProviderV1).
		Resource(subnetResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Network((*params.Network)[0].Name).
		Region(params.Region).
		Spec((*params.Subnet)[0].InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a subnet
	setCreatedRegionalNetworkResourceMetadata(subnetResponse.Metadata)
	subnetResponse.Status = newSubnetStatus(schema.ResourceStateCreating)
	subnetResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "CreateSubnet", nextState: "GetSubnet"}); err != nil {
		return nil, err
	}

	// Get the created subnet
	setSubnetState(subnetResponse.Status, schema.ResourceStateActive)
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "GetSubnet", nextState: "UpdateSubnet"}); err != nil {
		return nil, err
	}

	// Update the subnet
	setModifiedRegionalNetworkResourceMetadata(subnetResponse.Metadata)
	setSubnetState(subnetResponse.Status, schema.ResourceStateUpdating)
	subnetResponse.Spec = *(*params.Subnet)[0].UpdatedSpec
	subnetResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "UpdateSubnet", nextState: "GetSubnetUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated subnet
	setSubnetState(subnetResponse.Status, schema.ResourceStateActive)
	subnetResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, responseBody: subnetResponse, currentState: "GetSubnetUpdated", nextState: "CreatePublicIp"}); err != nil {
		return nil, err
	}

	// Public ip
	publicIpResponse, err := builders.NewPublicIpBuilder().
		Name((*params.PublicIp)[0].Name).
		Provider(secalib.NetworkProviderV1).
		Resource(publicIpResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec((*params.PublicIp)[0].InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a public ip
	setCreatedRegionalWorkspaceResourceMetadata(publicIpResponse.Metadata)
	publicIpResponse.Status = newPublicIpStatus(schema.ResourceStateCreating)
	publicIpResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "CreatePublicIp", nextState: "GetPublicIp"}); err != nil {
		return nil, err
	}

	// Get the created public ip
	setPublicIpState(publicIpResponse.Status, schema.ResourceStateActive)
	publicIpResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "GetPublicIp", nextState: "UpdatePublicIp"}); err != nil {
		return nil, err
	}

	// Update the public ip
	setModifiedRegionalWorkspaceResourceMetadata(publicIpResponse.Metadata)
	setPublicIpState(publicIpResponse.Status, schema.ResourceStateUpdating)
	publicIpResponse.Spec = *(*params.PublicIp)[0].UpdatedSpec
	publicIpResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "UpdatePublicIp", nextState: "GetPublicIpUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated public ip
	setPublicIpState(publicIpResponse.Status, schema.ResourceStateActive)
	publicIpResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, responseBody: publicIpResponse, currentState: "GetPublicIpUpdated", nextState: "CreateNIC"}); err != nil {
		return nil, err
	}

	// Nic
	nicResponse, err := builders.NewNicBuilder().
		Name((*params.NIC)[0].Name).
		Provider(secalib.NetworkProviderV1).
		Resource(nicResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec((*params.NIC)[0].InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a nic
	setCreatedRegionalWorkspaceResourceMetadata(nicResponse.Metadata)
	nicResponse.Status = newNicStatus(schema.ResourceStateCreating)
	nicResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "CreateNIC", nextState: "GetNIC"}); err != nil {
		return nil, err
	}

	// Get the created nic
	setNicState(nicResponse.Status, schema.ResourceStateActive)
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "GetNIC", nextState: "UpdateNIC"}); err != nil {
		return nil, err
	}

	// Update the nic
	setModifiedRegionalWorkspaceResourceMetadata(nicResponse.Metadata)
	setNicState(nicResponse.Status, schema.ResourceStateUpdating)
	nicResponse.Spec = *(*params.NIC)[0].UpdatedSpec
	nicResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "UpdateNIC", nextState: "GetNICUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated nic
	setNicState(nicResponse.Status, schema.ResourceStateActive)
	nicResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, responseBody: nicResponse, currentState: "GetNICUpdated", nextState: "CreateSecurityGroup"}); err != nil {
		return nil, err
	}

	// Security group
	groupResponse, err := builders.NewSecurityGroupBuilder().
		Name((*params.SecurityGroup)[0].Name).
		Provider(secalib.NetworkProviderV1).
		Resource(groupResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Workspace(params.Workspace.Name).
		Region(params.Region).
		Spec((*params.SecurityGroup)[0].InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a security group
	setCreatedRegionalWorkspaceResourceMetadata(groupResponse.Metadata)
	groupResponse.Status = newSecurityGroupStatus(schema.ResourceStateCreating)
	groupResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "CreateSecurityGroup", nextState: "GetSecurityGroup"}); err != nil {
		return nil, err
	}

	// Get the created security group
	setSecurityGroupState(groupResponse.Status, schema.ResourceStateActive)
	groupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "GetSecurityGroup", nextState: "UpdateSecurityGroup"}); err != nil {
		return nil, err
	}

	// Update the security group
	setModifiedRegionalWorkspaceResourceMetadata(groupResponse.Metadata)
	setSecurityGroupState(groupResponse.Status, schema.ResourceStateUpdating)
	groupResponse.Spec = *(*params.SecurityGroup)[0].UpdatedSpec
	groupResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "UpdateSecurityGroup", nextState: "GetSecurityGroupUpdated"}); err != nil {
		return nil, err
	}

	// Get the updated security group
	setSecurityGroupState(groupResponse.Status, schema.ResourceStateActive)
	groupResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, responseBody: groupResponse, currentState: "GetSecurityGroupUpdated", nextState: "CreateBlockStorage"}); err != nil {
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
	blockResponse.Status = newBlockStorageStatus(schema.ResourceStateCreating)
	blockResponse.Spec.SizeGB = params.BlockStorage.InitialSpec.SizeGB
	blockResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "CreateBlockStorage", nextState: "GetCreatedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the created block storage
	setBlockStorageState(blockResponse.Status, schema.ResourceStateActive)
	blockResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, responseBody: blockResponse, currentState: "GetCreatedBlockStorage", nextState: "CreateInstance"}); err != nil {
		return nil, err
	}

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
	instanceResponse.Status = newInstanceStatus(schema.ResourceStateCreating)
	instanceResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, responseBody: instanceResponse, currentState: "CreateInstance", nextState: "GetCreatedInstance"}); err != nil {
		return nil, err
	}

	// Get the created instance
	setInstanceState(instanceResponse.Status, schema.ResourceStateActive)
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
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: instanceUrl, params: params, currentState: "GetDeletedInstance", nextState: "DeleteBlockStorage"}); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, currentState: "DeleteBlockStorage", nextState: "GetDeletedBlockStorage"}); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: blockUrl, params: params, currentState: "GetDeletedBlockStorage", nextState: "DeleteSecurityGroup"}); err != nil {
		return nil, err
	}

	// Delete the security group
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, currentState: "DeleteSecurityGroup", nextState: "GetDeletedSecurityGroup"}); err != nil {
		return nil, err
	}

	// Get the deleted security group
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: groupUrl, params: params, currentState: "GetDeletedSecurityGroup", nextState: "DeleteNic"}); err != nil {
		return nil, err
	}

	// Delete the nic
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, currentState: "DeleteNic", nextState: "GetDeletedNic"}); err != nil {
		return nil, err
	}

	// Get the deleted nic
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: nicUrl, params: params, currentState: "GetDeletedNic", nextState: "DeletePublicIp"}); err != nil {
		return nil, err
	}

	// Delete the public ip
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, currentState: "DeletePublicIp", nextState: "GetDeletedPublicIp"}); err != nil {
		return nil, err
	}

	// Get the deleted public ip
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: publicIpUrl, params: params, currentState: "GetDeletedPublicIp", nextState: "DeleteSubnet"}); err != nil {
		return nil, err
	}

	// Delete the subnet
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, currentState: "DeleteSubnet", nextState: "GetDeletedSubnet"}); err != nil {
		return nil, err
	}

	// Get the deleted subnet
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: subnetUrl, params: params, currentState: "GetDeletedSubnet", nextState: "DeleteRouteTable"}); err != nil {
		return nil, err
	}

	// Delete the route table
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, currentState: "DeleteRouteTable", nextState: "GetDeletedRouteTable"}); err != nil {
		return nil, err
	}

	// Get the deleted route table
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: routeUrl, params: params, currentState: "GetDeletedRouteTable", nextState: "DeleteInternetGateway"}); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, currentState: "DeleteInternetGateway", nextState: "GetDeletedInternetGateway"}); err != nil {
		return nil, err
	}

	// Get the deleted internet gateway
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: gatewayUrl, params: params, currentState: "GetDeletedInternetGateway", nextState: "DeleteNetwork"}); err != nil {
		return nil, err
	}

	// Delete the network
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, currentState: "DeleteNetwork", nextState: "GetDeletedNetwork"}); err != nil {
		return nil, err
	}

	// Get the deleted network
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: networkUrl, params: params, currentState: "GetDeletedNetwork", nextState: "DeleteWorkspace"}); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, currentState: "DeleteWorkspace", nextState: "GetDeletedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, currentState: "GetDeletedWorkspace", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}

func ConfigNetworkListLifecycleScenarioV1(scenario string, params *NetworkParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)

	// Generate resources
	workspaceResource := secalib.GenerateWorkspaceResource(params.Tenant, params.Workspace.Name)

	// Workspace
	workResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(secalib.WorkspaceProviderV1).
		Resource(workspaceResource).
		ApiVersion(secalib.ApiVersion1).
		Tenant(params.Tenant).
		Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		BuildResponse()
	// Create a workspace
	setCreatedRegionalResourceMetadata(workResponse.Metadata)
	workResponse.Status = newWorkspaceStatus(schema.ResourceStateCreating)
	workResponse.Metadata.Verb = http.MethodPut
	if err := configurePutStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, responseBody: workResponse, currentState: startedScenarioState, nextState: (*params.Network)[0].Name}); err != nil {
		return nil, err
	}

	// Network
	var networkList []schema.Network
	for i := range *params.Network {
		networkResource := secalib.GenerateNetworkResource(params.Tenant, params.Workspace.Name, (*params.Network)[i].Name)

		networkResponse, err := builders.NewNetworkBuilder().
			Name((*params.Network)[0].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(networkResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Region(params.Region).
			Labels((*params.Network)[i].InitialLabels).
			Spec((*params.Network)[i].InitialSpec).
			BuildResponse()

		if err != nil {
			return nil, err
		}
		var nextState string
		if i < len(*params.Network)-1 {
			nextState = (*params.Network)[i+1].Name
		} else {
			nextState = "GetNetworkList"
		}
		// Create a network
		setCreatedRegionalWorkspaceResourceMetadata(networkResponse.Metadata)
		networkResponse.Status = newNetworkStatus(schema.ResourceStateCreating)
		networkResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateNetworkURL(params.Tenant, params.Workspace.Name, (*params.Network)[i].Name), params: params, responseBody: networkResponse, currentState: (*params.Network)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		networkList = append(networkList, *networkResponse)
	}

	// List
	networkResource := secalib.GenerateNetworkListResource(params.Tenant, params.Workspace.Name)
	networkListResponse := &networkV1.NetworkIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: networkResource,
			Verb:     http.MethodGet,
		},
		Items: networkList,
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateNetworkListURL(params.Tenant, params.Workspace.Name), params: params, responseBody: networkListResponse,
			currentState: "GetNetworkList", nextState: "GetNetworkListWithLimit",
		}); err != nil {
		return nil, err
	}

	// List limit
	networkListResponse.Items = networkList[:1]

	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateNetworkListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimit("1"), responseBody: networkListResponse,
			currentState: "GetNetworkListWithLimit", nextState: "GetNetworkListWithLabel",
		}); err != nil {
		return nil, err
	}
	// List label

	networkWithLabel := func(networkList []schema.Network) []schema.Network {
		var filteredNetworks []schema.Network
		for _, network := range networkList {
			if val, ok := network.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformanceLabel {
				filteredNetworks = append(filteredNetworks, network)
			}
		}
		return filteredNetworks
	}

	networkListResponse = &networkV1.NetworkIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: networkResource,
			Verb:     http.MethodGet,
		},
		Items: networkWithLabel(networkList),
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateNetworkListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: networkListResponse,
			currentState: "GetNetworkListWithLabel", nextState: "GetNetworkListWithLimitAndLabel",
		}); err != nil {
		return nil, err
	}

	// List limit & label

	networkListResponse = &networkV1.NetworkIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: networkResource,
			Verb:     http.MethodGet,
		},
		Items: networkWithLabel(networkList)[:1],
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateNetworkListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: networkListResponse,
			currentState: "GetNetworkListWithLimitAndLabel", nextState: (*params.InternetGateway)[0].Name,
		}); err != nil {
		return nil, err
	}

	// Internet gateway
	var gatewayList []schema.InternetGateway
	for i := range *params.InternetGateway {
		gatewayResource := secalib.GenerateInternetGatewayResource(params.Tenant, params.Workspace.Name, (*params.InternetGateway)[i].Name)
		gatewayResponse, err := builders.NewInternetGatewayBuilder().
			Name((*params.InternetGateway)[0].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(gatewayResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Region(params.Region).
			Labels((*params.InternetGateway)[i].InitialLabels).
			Spec((*params.InternetGateway)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}
		var nextState string
		if i < len(*params.InternetGateway)-1 {
			nextState = (*params.InternetGateway)[i+1].Name
		} else {
			nextState = "GetInternetGatewayList"
		}
		// Create an internet gateway
		setCreatedRegionalWorkspaceResourceMetadata(gatewayResponse.Metadata)
		gatewayResponse.Status = newResourceStatus(schema.ResourceStateCreating)
		gatewayResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{
				url: secalib.GenerateInternetGatewayURL(params.Tenant, params.Workspace.Name, (*params.InternetGateway)[i].Name), params: params, responseBody: gatewayResponse,
				currentState: (*params.InternetGateway)[i].Name, nextState: nextState,
			}); err != nil {
			return nil, err
		}
		gatewayList = append(gatewayList, *gatewayResponse)
	}

	// List
	gatewayListResource := secalib.GenerateInternetGatewayListResource(params.Tenant, params.Workspace.Name)
	gatewayListResponse := &networkV1.InternetGatewayIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: gatewayListResource,
			Verb:     http.MethodGet,
		},
		Items: gatewayList,
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateInternetGatewayListURL(params.Tenant, params.Workspace.Name), params: params, responseBody: gatewayListResponse, currentState: "GetInternetGatewayList", nextState: "GetInternetGatewayListWithLimit"}); err != nil {
		return nil, err
	}
	// List with limit
	gatewayListWithLimitResponse := &networkV1.InternetGatewayIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: gatewayListResource,
			Verb:     http.MethodGet,
		},
		Items: gatewayList[:1],
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateInternetGatewayListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimit("1"), responseBody: gatewayListWithLimitResponse, currentState: "GetInternetGatewayListWithLimit", nextState: "GetInternetGatewayListWithLabel"}); err != nil {
		return nil, err
	}
	// List with label

	gatewayWithLabel := func(gatewayList []schema.InternetGateway) []schema.InternetGateway {
		var filteredGateway []schema.InternetGateway
		for _, gateway := range gatewayList {
			if val, ok := gateway.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformanceLabel {
				filteredGateway = append(filteredGateway, gateway)
			}
		}
		return filteredGateway
	}

	gatewayListResponse = &networkV1.InternetGatewayIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: gatewayListResource,
			Verb:     http.MethodGet,
		},
		Items: gatewayWithLabel(gatewayList),
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateInternetGatewayListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: gatewayListResponse, currentState: "GetInternetGatewayListWithLabel", nextState: "GetInternetGatewayListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List with limit & label

	gatewayListResponse = &networkV1.InternetGatewayIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: gatewayListResource,
			Verb:     http.MethodGet,
		},
		Items: gatewayWithLabel(gatewayList)[:1],
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateInternetGatewayListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: gatewayListResponse,
			currentState: "GetInternetGatewayListWithLimitAndLabel", nextState: (*params.RouteTable)[0].Name,
		}); err != nil {
		return nil, err
	}

	// Route table

	var routeTableList []schema.RouteTable
	for i := range *params.RouteTable {
		routeTableResource := secalib.GenerateRouteTableResource(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.RouteTable)[i].Name)
		routeTableResponse, err := builders.NewRouteTableBuilder().
			Name((*params.RouteTable)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(routeTableResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Network((*params.Network)[i].Name).
			Region(params.Region).
			Labels((*params.RouteTable)[i].InitialLabels).
			Spec((*params.RouteTable)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}
		var nextState string
		if i < len(*params.RouteTable)-1 {
			nextState = (*params.RouteTable)[i+1].Name
		} else {
			nextState = "GetRouteTableList"
		}
		// Create an internet gateway
		setCreatedRegionalNetworkResourceMetadata(routeTableResponse.Metadata)
		routeTableResponse.Status = newRouteTableStatus(schema.ResourceStateCreating)
		routeTableResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateRouteTableURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.RouteTable)[i].Name), params: params, responseBody: routeTableResponse, currentState: (*params.RouteTable)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		routeTableList = append(routeTableList, *routeTableResponse)
	}

	// List
	routeTableListResource := secalib.GenerateRouteTableListResource(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name)
	routeTableListResponse := &networkV1.RouteTableIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: routeTableListResource,
			Verb:     http.MethodGet,
		},
		Items: routeTableList,
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateRouteTableListURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name), params: params, responseBody: routeTableListResponse,
			currentState: "GetRouteTableList", nextState: "GetRouteTableListWithLimit",
		}); err != nil {
		return nil, err
	}
	// List with limit

	routeTableListResponse.Items = routeTableList[:1]

	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateRouteTableListURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name), params: params, pathParams: pathParamsLimit("1"), responseBody: routeTableListResponse,
			currentState: "GetRouteTableListWithLimit", nextState: "GetRouteTableListWithLabel",
		}); err != nil {
		return nil, err
	}
	// List with label
	routeTableWithLabel := func(routeTableList []schema.RouteTable) []schema.RouteTable {
		var filteredRouteTables []schema.RouteTable
		for _, routeTable := range routeTableList {
			if val, ok := routeTable.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformanceLabel {
				filteredRouteTables = append(filteredRouteTables, routeTable)
			}
		}
		return filteredRouteTables
	}

	routeTableListResponse.Items = routeTableWithLabel(routeTableList)
	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateRouteTableListURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: routeTableListResponse,
			currentState: "GetRouteTableListWithLabel", nextState: "GetRouteTableListWithLimitAndLabel",
		}); err != nil {
		return nil, err
	}
	// List with limit & label

	routeTableListResponse.Items = routeTableWithLabel(routeTableList)[:1]
	if err := configureGetStub(wm, scenario,
		&stubConfig{
			url: secalib.GenerateRouteTableListURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: routeTableListResponse,
			currentState: "GetRouteTableListWithLimitAndLabel", nextState: (*params.Subnet)[0].Name,
		}); err != nil {
		return nil, err
	}

	// Subnet

	var subnetList []schema.Subnet
	for i := range *params.Subnet {
		subnetResource := secalib.GenerateSubnetResource(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.Subnet)[i].Name)
		subnetResponse, err := builders.NewSubnetBuilder().
			Name((*params.Subnet)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(subnetResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Network((*params.Network)[i].Name).
			Region(params.Region).
			Labels((*params.Subnet)[i].InitialLabels).
			Spec((*params.Subnet)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}
		var nextState string
		if i < len(*params.Subnet)-1 {
			nextState = (*params.Subnet)[i+1].Name
		} else {
			nextState = "GetSubnetList"
		}
		// Create an internet gateway
		setCreatedRegionalNetworkResourceMetadata(subnetResponse.Metadata)
		subnetResponse.Status = newSubnetStatus(schema.ResourceStateCreating)
		subnetResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateSubnetURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.Subnet)[i].Name), params: params, responseBody: subnetResponse, currentState: (*params.Subnet)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		subnetList = append(subnetList, *subnetResponse)
	}

	// List
	subnetListResource := secalib.GenerateSubnetListResource(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name)
	subnetListResponse := &networkV1.SubnetIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: subnetListResource,
			Verb:     http.MethodGet,
		},
		Items: subnetList,
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateSubnetListURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name), params: params, responseBody: subnetListResponse, currentState: "GetSubnetList", nextState: "GetSubnetListWithLimit"}); err != nil {
		return nil, err
	}
	// List with limit

	subnetListResponse.Items = subnetList[:1]

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateSubnetListURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name), params: params, pathParams: pathParamsLimit("1"), responseBody: subnetListResponse, currentState: "GetSubnetListWithLimit", nextState: "GetSubnetListWithLabel"}); err != nil {
		return nil, err
	}
	// List with label

	subnetWithLabel := func(subnetList []schema.Subnet) []schema.Subnet {
		var filteredInstances []schema.Subnet
		for _, subnet := range subnetList {
			if val, ok := subnet.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformanceLabel {
				filteredInstances = append(filteredInstances, subnet)
			}
		}
		return filteredInstances
	}

	subnetListResponse.Items = subnetWithLabel(subnetList)

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateSubnetListURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: subnetListResponse, currentState: "GetSubnetListWithLabel", nextState: "GetSubnetListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List with limit & label

	subnetListResponse = &networkV1.SubnetIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: subnetListResource,
			Verb:     http.MethodGet,
		},
		Items: subnetWithLabel(subnetList)[:1],
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateSubnetListURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: subnetListResponse, currentState: "GetSubnetListWithLimitAndLabel", nextState: (*params.PublicIp)[0].Name}); err != nil {
		return nil, err
	}

	// Public ip

	var publicIpList []schema.PublicIp
	for i := range *params.PublicIp {
		publicIpResource := secalib.GeneratePublicIpResource(params.Tenant, params.Workspace.Name, (*params.PublicIp)[i].Name)
		publicIpResponse, err := builders.NewPublicIpBuilder().
			Name((*params.PublicIp)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(publicIpResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Region(params.Region).
			Labels((*params.PublicIp)[i].InitialLabels).
			Spec((*params.PublicIp)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}
		var nextState string
		if i < len(*params.PublicIp)-1 {
			nextState = (*params.PublicIp)[i+1].Name
		} else {
			nextState = "GetPublicIpList"
		}
		// Create a PublicIp
		setCreatedRegionalWorkspaceResourceMetadata(publicIpResponse.Metadata)
		publicIpResponse.Status = newPublicIpStatus(schema.ResourceStateCreating)
		publicIpResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GeneratePublicIpURL(params.Tenant, params.Workspace.Name, (*params.PublicIp)[i].Name), params: params, responseBody: publicIpResponse, currentState: (*params.PublicIp)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		publicIpList = append(publicIpList, *publicIpResponse)
	}

	// List
	publicIpListResource := secalib.GeneratePublicIpListResource(params.Tenant, params.Workspace.Name)
	publicIpListResponse := &networkV1.PublicIpIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: publicIpListResource,
			Verb:     http.MethodGet,
		},
		Items: publicIpList,
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GeneratePublicIpListURL(params.Tenant, params.Workspace.Name), params: params, responseBody: publicIpListResponse, currentState: "GetPublicIpList", nextState: "GetPublicIpListWithLimit"}); err != nil {
		return nil, err
	}
	// List with limit
	publicIpListResponse = &networkV1.PublicIpIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: publicIpListResource,
			Verb:     http.MethodGet,
		},
		Items: publicIpList[:1],
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GeneratePublicIpListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimit("1"), responseBody: publicIpListResponse, currentState: "GetPublicIpListWithLimit", nextState: "GetPublicIpListWithLabel"}); err != nil {
		return nil, err
	}
	// List with label

	publicIpWithLabel := func(publicIpList []schema.PublicIp) []schema.PublicIp {
		var filteredInstances []schema.PublicIp
		for _, publicIp := range publicIpList {
			if val, ok := publicIp.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformanceLabel {
				filteredInstances = append(filteredInstances, publicIp)
			}
		}
		return filteredInstances
	}

	publicIpListResponse = &networkV1.PublicIpIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: publicIpListResource,
			Verb:     http.MethodGet,
		},
		Items: publicIpWithLabel(publicIpList),
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GeneratePublicIpListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: publicIpListResponse, currentState: "GetPublicIpListWithLabel", nextState: "GetPublicIpListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List with limit & label

	publicIpListResponse = &networkV1.PublicIpIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: publicIpListResource,
			Verb:     http.MethodGet,
		},
		Items: publicIpWithLabel(publicIpList)[:1],
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GeneratePublicIpListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: publicIpListResponse, currentState: "GetPublicIpListWithLimitAndLabel", nextState: (*params.NIC)[0].Name}); err != nil {
		return nil, err
	}

	// Nic

	var nicList []schema.Nic
	for i := range *params.NIC {
		nicResource := secalib.GenerateNicResource(params.Tenant, params.Workspace.Name, (*params.NIC)[i].Name)
		nicResponse, err := builders.NewNicBuilder().
			Name((*params.NIC)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(nicResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Region(params.Region).
			Labels((*params.NIC)[i].InitialLabels).
			Spec((*params.NIC)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}

		var nextState string
		if i < len(*params.NIC)-1 {
			nextState = (*params.NIC)[i+1].Name
		} else {
			nextState = "GetNICList"
		}
		// Create an internet gateway
		setCreatedRegionalWorkspaceResourceMetadata(nicResponse.Metadata)
		nicResponse.Status = newNicStatus(schema.ResourceStateCreating)
		nicResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateNicURL(params.Tenant, params.Workspace.Name, (*params.NIC)[i].Name), params: params, responseBody: nicResponse, currentState: (*params.NIC)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		nicList = append(nicList, *nicResponse)
	}

	// List
	nicListResource := secalib.GenerateNicListResource(params.Tenant, params.Workspace.Name)
	nicListResponse := &networkV1.NicIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: nicListResource,
			Verb:     http.MethodGet,
		},
		Items: nicList,
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateNicListURL(params.Tenant, params.Workspace.Name), params: params, responseBody: nicListResponse, currentState: "GetNICList", nextState: "GetNICListWithLimit"}); err != nil {
		return nil, err
	}
	// List with limit

	nicListResponse = &networkV1.NicIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: nicListResource,
			Verb:     http.MethodGet,
		},
		Items: nicList[:1],
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateNicListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimit("1"), responseBody: nicListResponse, currentState: "GetNICListWithLimit", nextState: "GetNICListWithLabel"}); err != nil {
		return nil, err
	}
	// List with label

	nicWithLabel := func(nicList []schema.Nic) []schema.Nic {
		var filteredInstances []schema.Nic
		for _, nic := range nicList {
			if val, ok := nic.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformanceLabel {
				filteredInstances = append(filteredInstances, nic)
			}
		}
		return filteredInstances
	}

	nicListResponse = &networkV1.NicIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: nicListResource,
			Verb:     http.MethodGet,
		},
		Items: nicWithLabel(nicList),
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateNicListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: nicListResponse, currentState: "GetNICListWithLabel", nextState: "GetNICListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List with limit & label

	nicListResponse = &networkV1.NicIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: nicListResource,
			Verb:     http.MethodGet,
		},
		Items: nicWithLabel(nicList)[:1],
	}
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateNicListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: nicListResponse, currentState: "GetNICListWithLimitAndLabel", nextState: (*params.SecurityGroup)[0].Name}); err != nil {
		return nil, err
	}

	// Security group
	securityGroupList := []schema.SecurityGroup{}
	for i := range *params.SecurityGroup {
		securityGroupResource := secalib.GenerateSecurityGroupResource(params.Tenant, params.Workspace.Name, (*params.SecurityGroup)[i].Name)
		securityGroupResponse, err := builders.NewSecurityGroupBuilder().
			Name((*params.SecurityGroup)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(securityGroupResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Region(params.Region).
			Labels((*params.SecurityGroup)[i].InitialLabels).
			Spec((*params.SecurityGroup)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}

		var nextState string
		if i < len(*params.SecurityGroup)-1 {
			nextState = (*params.SecurityGroup)[i+1].Name
		} else {
			nextState = "GetSecurityGroupList"
		}
		// Create an internet gateway
		setCreatedRegionalWorkspaceResourceMetadata(securityGroupResponse.Metadata)
		securityGroupResponse.Status = newSecurityGroupStatus(schema.ResourceStateCreating)
		securityGroupResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateSecurityGroupURL(params.Tenant, params.Workspace.Name, (*params.SecurityGroup)[i].Name), params: params, responseBody: securityGroupResponse, currentState: (*params.SecurityGroup)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		securityGroupList = append(securityGroupList, *securityGroupResponse)
	}

	// List
	securityGroupListResource := secalib.GenerateSecurityGroupListResource(params.Tenant, params.Workspace.Name)
	securityGroupListResponse := &networkV1.SecurityGroupIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.NetworkProviderV1,
			Resource: securityGroupListResource,
			Verb:     http.MethodGet,
		},
		Items: securityGroupList,
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateSecurityGroupListURL(params.Tenant, params.Workspace.Name), params: params, responseBody: securityGroupListResponse, currentState: "GetSecurityGroupList", nextState: "GetSecurityGroupListWithLimit"}); err != nil {
		return nil, err
	}
	// List with limit
	securityGroupListResponse.Items = securityGroupList[:1]

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateSecurityGroupListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimit("1"), responseBody: securityGroupListResponse, currentState: "GetSecurityGroupListWithLimit", nextState: "GetSecurityGroupListWithLabel"}); err != nil {
		return nil, err
	}
	// List with label

	securityGroupWithLabel := func(securityGroupList []schema.SecurityGroup) []schema.SecurityGroup {
		var filteredInstances []schema.SecurityGroup
		for _, securityGroup := range securityGroupList {
			if val, ok := securityGroup.Labels[secalib.EnvLabel]; ok && val == secalib.EnvConformanceLabel {
				filteredInstances = append(filteredInstances, securityGroup)
			}
		}
		return filteredInstances
	}

	securityGroupListResponse.Items = securityGroupWithLabel(securityGroupList)

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateSecurityGroupListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLabel(secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: securityGroupListResponse, currentState: "GetSecurityGroupListWithLabel", nextState: "GetSecurityGroupListWithLimitAndLabel"}); err != nil {
		return nil, err
	}
	// List with limit & label

	securityGroupListResponse.Items = securityGroupWithLabel(securityGroupList)[:1]

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.GenerateSecurityGroupListURL(params.Tenant, params.Workspace.Name), params: params, pathParams: pathParamsLimitAndLabel("1", secalib.EnvLabel, secalib.EnvConformanceLabel), responseBody: securityGroupListResponse, currentState: "GetSecurityGroupListWithLimitAndLabel", nextState: "DeleteSecurityGroup_" + (*params.SecurityGroup)[0].Name}); err != nil {
		return nil, err
	}

	// Delete SecurityGroups
	for i := range *params.SecurityGroup {
		securityGroupUrl := secalib.GenerateSecurityGroupURL(params.Tenant, params.Workspace.Name, (*params.SecurityGroup)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteSecurityGroup_" + (*params.SecurityGroup)[i].Name
		} else {
			currentState = "GetDeletedSecurityGroup_" + (*params.SecurityGroup)[i-1].Name
		}

		nextState = "DeleteSecurityGroup_" + (*params.SecurityGroup)[i].Name

		// Delete the security group
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: securityGroupUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted security group (should return 404)
		nextState = func() string {
			if i < len(*params.SecurityGroup)-1 {
				return "GetDeletedSecurityGroup_" + (*params.SecurityGroup)[i].Name
			} else {
				return "DeleteNic_" + (*params.NIC)[0].Name
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: securityGroupUrl, params: params, currentState: "DeleteSecurityGroup_" + (*params.SecurityGroup)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
	}

	// Delete NICs
	for i := range *params.NIC {
		nicUrl := secalib.GenerateNicURL(params.Tenant, params.Workspace.Name, (*params.NIC)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteNic_" + (*params.NIC)[i].Name
		} else {
			currentState = "GetDeletedNic_" + (*params.NIC)[i-1].Name
		}

		nextState = "DeleteNic_" + (*params.NIC)[i].Name

		// Delete the NIC
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: nicUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted NIC (should return 404)
		nextState = func() string {
			if i < len(*params.NIC)-1 {
				return "GetDeletedNic_" + (*params.NIC)[i].Name
			} else {
				return "DeletePublicIp_" + (*params.PublicIp)[0].Name
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: nicUrl, params: params, currentState: "DeleteNic_" + (*params.NIC)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
	}

	// Delete Public IPs
	for i := range *params.PublicIp {
		publicIpUrl := secalib.GeneratePublicIpURL(params.Tenant, params.Workspace.Name, (*params.PublicIp)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeletePublicIp_" + (*params.PublicIp)[i].Name
		} else {
			currentState = "GetDeletedPublicIp_" + (*params.PublicIp)[i-1].Name
		}

		nextState = "DeletePublicIp_" + (*params.PublicIp)[i].Name

		// Delete the public IP
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: publicIpUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted public IP (should return 404)
		nextState = func() string {
			if i < len(*params.PublicIp)-1 {
				return "GetDeletedPublicIp_" + (*params.PublicIp)[i].Name
			} else {
				return "DeleteSubnet_" + (*params.Subnet)[0].Name
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: publicIpUrl, params: params, currentState: "DeletePublicIp_" + (*params.PublicIp)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
	}

	// Delete Subnets
	for i := range *params.Subnet {
		subnetUrl := secalib.GenerateSubnetURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.Subnet)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteSubnet_" + (*params.Subnet)[i].Name
		} else {
			currentState = "GetDeletedSubnet_" + (*params.Subnet)[i-1].Name
		}

		nextState = "DeleteSubnet_" + (*params.Subnet)[i].Name

		// Delete the subnet
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: subnetUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted subnet (should return 404)
		nextState = func() string {
			if i < len(*params.Subnet)-1 {
				return "GetDeletedSubnet_" + (*params.Subnet)[i].Name
			} else {
				return "DeleteRouteTable_" + (*params.RouteTable)[0].Name
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: subnetUrl, params: params, currentState: "DeleteSubnet_" + (*params.Subnet)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
	}

	// Delete Route Tables
	for i := range *params.RouteTable {
		routeUrl := secalib.GenerateRouteTableURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.RouteTable)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteRouteTable_" + (*params.RouteTable)[i].Name
		} else {
			currentState = "GetDeletedRouteTable_" + (*params.RouteTable)[i-1].Name
		}

		nextState = "DeleteRouteTable_" + (*params.RouteTable)[i].Name

		// Delete the route table
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: routeUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted route table (should return 404)
		nextState = func() string {
			if i < len(*params.RouteTable)-1 {
				return "GetDeletedRouteTable_" + (*params.RouteTable)[i].Name
			} else {
				return "DeleteGateway_" + (*params.InternetGateway)[0].Name
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: routeUrl, params: params, currentState: "DeleteRouteTable_" + (*params.RouteTable)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
	}

	// Delete Internet Gateways
	for i := range *params.InternetGateway {
		gatewayUrl := secalib.GenerateInternetGatewayURL(params.Tenant, params.Workspace.Name, (*params.InternetGateway)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteGateway_" + (*params.InternetGateway)[i].Name
		} else {
			currentState = "GetDeletedGateway_" + (*params.InternetGateway)[i-1].Name
		}

		nextState = "DeleteGateway_" + (*params.InternetGateway)[i].Name

		// Delete the gateway
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: gatewayUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted gateway (should return 404)
		nextState = func() string {
			if i < len(*params.InternetGateway)-1 {
				return "GetDeletedGateway_" + (*params.InternetGateway)[i].Name
			} else {
				return "DeleteNetwork_" + (*params.Network)[0].Name
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: gatewayUrl, params: params, currentState: "DeleteGateway_" + (*params.InternetGateway)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
	}

	// Delete Networks
	for i := range *params.Network {
		networkUrl := secalib.GenerateNetworkURL(params.Tenant, params.Workspace.Name, (*params.Network)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteNetwork_" + (*params.Network)[i].Name
		} else {
			currentState = "GetDeletedNetwork_" + (*params.Network)[i-1].Name
		}

		nextState = "DeleteNetwork_" + (*params.Network)[i].Name

		// Delete the network
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: networkUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return nil, err
		}

		// Get the deleted network (should return 404)
		nextState = func() string {
			if i < len(*params.Network)-1 {
				return "GetDeletedNetwork_" + (*params.Network)[i].Name
			} else {
				return "DeleteWorkspace"
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: networkUrl, params: params, currentState: "DeleteNetwork_" + (*params.Network)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
	}

	// Delete the workspace
	workspaceUrl = secalib.GenerateWorkspaceURL(params.Tenant, params.Workspace.Name)
	if err := configureDeleteStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, currentState: "DeleteWorkspace", nextState: "GetDeletedWorkspace"}); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configureGetNotFoundStub(wm, scenario,
		&stubConfig{url: workspaceUrl, params: params, currentState: "GetDeletedWorkspace", nextState: startedScenarioState}); err != nil {
		return nil, err
	}

	return wm, nil
}
