package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/generators"
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
	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := generators.GenerateInstanceURL(computeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkUrl := generators.GenerateNetworkURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name)
	gatewayUrl := generators.GenerateInternetGatewayURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicUrl := generators.GenerateNicURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Nic.Name)
	publicIpUrl := generators.GeneratePublicIpURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeUrl := generators.GenerateRouteTableURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name, params.RouteTable.Name)
	subnetUrl := generators.GenerateSubnetURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name, params.Subnet.Name)
	groupUrl := generators.GenerateSecurityGroupURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
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
		Name(params.Network.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Network.InitialSpec).
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
	networkResponse.Spec = *params.Network.UpdatedSpec
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
		Name(params.InternetGateway.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.InternetGateway.InitialSpec).
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
	gatewayResponse.Spec = *params.InternetGateway.UpdatedSpec
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
		Name(params.RouteTable.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(params.Network.Name).Region(params.Region).
		Spec(params.RouteTable.InitialSpec).
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
	routeResponse.Spec = *params.RouteTable.UpdatedSpec
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
		Name(params.Subnet.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(params.Network.Name).Region(params.Region).
		Spec(params.Subnet.InitialSpec).
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
	subnetResponse.Spec = *params.Subnet.UpdatedSpec
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
		Name(params.PublicIp.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.PublicIp.InitialSpec).
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
	publicIpResponse.Spec = *params.PublicIp.UpdatedSpec
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
		Name(params.Nic.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Nic.InitialSpec).
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
	nicResponse.Spec = *params.Nic.UpdatedSpec
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
		Name(params.SecurityGroup.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.SecurityGroup.InitialSpec).
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
	groupResponse.Spec = *params.SecurityGroup.UpdatedSpec
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
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
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
		Provider(computeProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
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
