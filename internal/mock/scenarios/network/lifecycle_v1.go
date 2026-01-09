package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"

	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, params *params.NetworkLifeCycleParamsV1) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	configurator, err := stubs.NewStubConfigurator(scenario, params.MockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := generators.GenerateInstanceURL(constants.ComputeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkUrl := generators.GenerateNetworkURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name)
	gatewayUrl := generators.GenerateInternetGatewayURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicUrl := generators.GenerateNicURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.Nic.Name)
	publicIpUrl := generators.GeneratePublicIpURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeUrl := generators.GenerateRouteTableURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name, params.RouteTable.Name)
	subnetUrl := generators.GenerateSubnetURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name, params.Subnet.Name)
	groupUrl := generators.GenerateSecurityGroupURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Network
	networkResponse, err := builders.NewNetworkBuilder().
		Name(params.Network.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Network.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a network
	if err := configurator.ConfigureCreateNetworkStub(networkResponse, networkUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created network
	if err := configurator.ConfigureGetActiveNetworkStub(networkResponse, networkUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the network
	networkResponse.Spec = *params.Network.UpdatedSpec
	if err := configurator.ConfigureUpdateNetworkStub(networkResponse, networkUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated network
	if err := configurator.ConfigureGetActiveNetworkStub(networkResponse, networkUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Internet gateway
	gatewayResponse, err := builders.NewInternetGatewayBuilder().
		Name(params.InternetGateway.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.InternetGateway.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an internet gateway
	if err := configurator.ConfigureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created internet gateway
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the internet gateway
	gatewayResponse.Spec = *params.InternetGateway.UpdatedSpec
	if err := configurator.ConfigureUpdateInternetGatewayStub(gatewayResponse, gatewayUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated internet gateway
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Route table
	routeResponse, err := builders.NewRouteTableBuilder().
		Name(params.RouteTable.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(params.Network.Name).Region(params.Region).
		Spec(params.RouteTable.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a route table
	if err := configurator.ConfigureCreateRouteTableStub(routeResponse, routeUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created route table
	if err := configurator.ConfigureGetActiveRouteTableStub(routeResponse, routeUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the route table
	routeResponse.Spec = *params.RouteTable.UpdatedSpec
	if err := configurator.ConfigureUpdateRouteTableStub(routeResponse, routeUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated route table
	if err := configurator.ConfigureGetActiveRouteTableStub(routeResponse, routeUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Subnet
	subnetResponse, err := builders.NewSubnetBuilder().
		Name(params.Subnet.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(params.Network.Name).Region(params.Region).
		Spec(params.Subnet.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a subnet
	if err := configurator.ConfigureCreateSubnetStub(subnetResponse, subnetUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created subnet
	if err := configurator.ConfigureGetActiveSubnetStub(subnetResponse, subnetUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the subnet
	subnetResponse.Spec = *params.Subnet.UpdatedSpec
	if err := configurator.ConfigureUpdateSubnetStub(subnetResponse, subnetUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated subnet
	if err := configurator.ConfigureGetActiveSubnetStub(subnetResponse, subnetUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Public ip
	publicIpResponse, err := builders.NewPublicIpBuilder().
		Name(params.PublicIp.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.PublicIp.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a public ip
	if err := configurator.ConfigureCreatePublicIpStub(publicIpResponse, publicIpUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created public ip
	if err := configurator.ConfigureGetActivePublicIpStub(publicIpResponse, publicIpUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the public ip
	publicIpResponse.Spec = *params.PublicIp.UpdatedSpec
	if err := configurator.ConfigureUpdatePublicIpStub(publicIpResponse, publicIpUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated public ip
	if err := configurator.ConfigureGetActivePublicIpStub(publicIpResponse, publicIpUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Nic
	nicResponse, err := builders.NewNicBuilder().
		Name(params.Nic.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Nic.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a nic
	if err := configurator.ConfigureCreateNicStub(nicResponse, nicUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created nic
	if err := configurator.ConfigureGetActiveNicStub(nicResponse, nicUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the nic
	nicResponse.Spec = *params.Nic.UpdatedSpec
	if err := configurator.ConfigureUpdateNicStub(nicResponse, nicUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated nic
	if err := configurator.ConfigureGetActiveNicStub(nicResponse, nicUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Security group
	groupResponse, err := builders.NewSecurityGroupBuilder().
		Name(params.SecurityGroup.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.SecurityGroup.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a security group
	if err := configurator.ConfigureCreateSecurityGroupStub(groupResponse, groupUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created security group
	if err := configurator.ConfigureGetActiveSecurityGroupStub(groupResponse, groupUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Update the security group
	groupResponse.Spec = *params.SecurityGroup.UpdatedSpec
	if err := configurator.ConfigureUpdateSecurityGroupStub(groupResponse, groupUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the updated security group
	if err := configurator.ConfigureGetActiveSecurityGroupStub(groupResponse, groupUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(params.Instance.Name).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Instance.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an instance
	if err := configurator.ConfigureCreateInstanceStub(instanceResponse, instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.ConfigureGetActiveInstanceStub(instanceResponse, instanceUrl, params.MockParams); err != nil {
		return nil, err
	}
	// Delete the instance
	if err := configurator.ConfigureDeleteStub(instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.ConfigureGetNotFoundStub(instanceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the security group
	if err := configurator.ConfigureDeleteStub(groupUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted security group
	if err := configurator.ConfigureGetNotFoundStub(groupUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the nic
	if err := configurator.ConfigureDeleteStub(nicUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted nic
	if err := configurator.ConfigureGetNotFoundStub(nicUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the public ip
	if err := configurator.ConfigureDeleteStub(publicIpUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted public ip
	if err := configurator.ConfigureGetNotFoundStub(publicIpUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the subnet
	if err := configurator.ConfigureDeleteStub(subnetUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted subnet
	if err := configurator.ConfigureGetNotFoundStub(subnetUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the route table
	if err := configurator.ConfigureDeleteStub(routeUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted route table
	if err := configurator.ConfigureGetNotFoundStub(routeUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	if err := configurator.ConfigureDeleteStub(gatewayUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted internet gateway
	if err := configurator.ConfigureGetNotFoundStub(gatewayUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the network
	if err := configurator.ConfigureDeleteStub(networkUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted network
	if err := configurator.ConfigureGetNotFoundStub(networkUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, params.MockParams); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
