package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/wiremock/go-wiremock"
)

func ConfigNetworkLifecycleScenarioV1(scenario string, params *NetworkParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
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
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.configureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.configureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
		return nil, err
	}

	// Network
	networkResponse, err := builders.NewNetworkBuilder().
		Name(params.Network.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Network.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a network
	if err := configurator.configureCreateNetworkStub(networkResponse, networkUrl, params); err != nil {
		return nil, err
	}

	// Get the created network
	if err := configurator.configureGetActiveNetworkStub(networkResponse, networkUrl, params); err != nil {
		return nil, err
	}

	// Update the network
	networkResponse.Spec = *params.Network.UpdatedSpec
	if err := configurator.configureUpdateNetworkStub(networkResponse, networkUrl, params); err != nil {
		return nil, err
	}

	// Get the updated network
	if err := configurator.configureGetActiveNetworkStub(networkResponse, networkUrl, params); err != nil {
		return nil, err
	}

	// Internet gateway
	gatewayResponse, err := builders.NewInternetGatewayBuilder().
		Name(params.InternetGateway.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.InternetGateway.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an internet gateway
	if err := configurator.configureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, params); err != nil {
		return nil, err
	}

	// Get the created internet gateway
	if err := configurator.configureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, params); err != nil {
		return nil, err
	}

	// Update the internet gateway
	gatewayResponse.Spec = *params.InternetGateway.UpdatedSpec
	if err := configurator.configureUpdateInternetGatewayStub(gatewayResponse, gatewayUrl, params); err != nil {
		return nil, err
	}

	// Get the updated internet gateway
	if err := configurator.configureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, params); err != nil {
		return nil, err
	}

	// Route table
	routeResponse, err := builders.NewRouteTableBuilder().
		Name(params.RouteTable.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(params.Network.Name).Region(params.Region).
		Spec(params.RouteTable.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a route table
	if err := configurator.configureCreateRouteTableStub(routeResponse, routeUrl, params); err != nil {
		return nil, err
	}

	// Get the created route table
	if err := configurator.configureGetActiveRouteTableStub(routeResponse, routeUrl, params); err != nil {
		return nil, err
	}

	// Update the route table
	routeResponse.Spec = *params.RouteTable.UpdatedSpec
	if err := configurator.configureUpdateRouteTableStub(routeResponse, routeUrl, params); err != nil {
		return nil, err
	}

	// Get the updated route table
	if err := configurator.configureGetActiveRouteTableStub(routeResponse, routeUrl, params); err != nil {
		return nil, err
	}

	// Subnet
	subnetResponse, err := builders.NewSubnetBuilder().
		Name(params.Subnet.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(params.Network.Name).Region(params.Region).
		Spec(params.Subnet.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a subnet
	if err := configurator.configureCreateSubnetStub(subnetResponse, subnetUrl, params); err != nil {
		return nil, err
	}

	// Get the created subnet
	if err := configurator.configureGetActiveSubnetStub(subnetResponse, subnetUrl, params); err != nil {
		return nil, err
	}

	// Update the subnet
	subnetResponse.Spec = *params.Subnet.UpdatedSpec
	if err := configurator.configureUpdateSubnetStub(subnetResponse, subnetUrl, params); err != nil {
		return nil, err
	}

	// Get the updated subnet
	if err := configurator.configureGetActiveSubnetStub(subnetResponse, subnetUrl, params); err != nil {
		return nil, err
	}

	// Public ip
	publicIpResponse, err := builders.NewPublicIpBuilder().
		Name(params.PublicIp.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.PublicIp.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a public ip
	if err := configurator.configureCreatePublicIpStub(publicIpResponse, publicIpUrl, params); err != nil {
		return nil, err
	}

	// Get the created public ip
	if err := configurator.configureGetActivePublicIpStub(publicIpResponse, publicIpUrl, params); err != nil {
		return nil, err
	}

	// Update the public ip
	publicIpResponse.Spec = *params.PublicIp.UpdatedSpec
	if err := configurator.configureUpdatePublicIpStub(publicIpResponse, publicIpUrl, params); err != nil {
		return nil, err
	}

	// Get the updated public ip
	if err := configurator.configureGetActivePublicIpStub(publicIpResponse, publicIpUrl, params); err != nil {
		return nil, err
	}

	// Nic
	nicResponse, err := builders.NewNicBuilder().
		Name(params.Nic.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Nic.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a nic
	if err := configurator.configureCreateNicStub(nicResponse, nicUrl, params); err != nil {
		return nil, err
	}

	// Get the created nic
	if err := configurator.configureGetActiveNicStub(nicResponse, nicUrl, params); err != nil {
		return nil, err
	}

	// Update the nic
	nicResponse.Spec = *params.Nic.UpdatedSpec
	if err := configurator.configureUpdateNicStub(nicResponse, nicUrl, params); err != nil {
		return nil, err
	}

	// Get the updated nic
	if err := configurator.configureGetActiveNicStub(nicResponse, nicUrl, params); err != nil {
		return nil, err
	}

	// Security group
	groupResponse, err := builders.NewSecurityGroupBuilder().
		Name(params.SecurityGroup.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.SecurityGroup.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a security group
	if err := configurator.configureCreateSecurityGroupStub(groupResponse, groupUrl, params); err != nil {
		return nil, err
	}

	// Get the created security group
	if err := configurator.configureGetActiveSecurityGroupStub(groupResponse, groupUrl, params); err != nil {
		return nil, err
	}

	// Update the security group
	groupResponse.Spec = *params.SecurityGroup.UpdatedSpec
	if err := configurator.configureUpdateSecurityGroupStub(groupResponse, groupUrl, params); err != nil {
		return nil, err
	}

	// Get the updated security group
	if err := configurator.configureGetActiveSecurityGroupStub(groupResponse, groupUrl, params); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.configureGetActiveBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(params.Instance.Name).
		Provider(computeProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Instance.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an instance
	if err := configurator.configureCreateInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}
	// Delete the instance
	if err := configurator.configureDeleteStub(instanceUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.configureGetNotFoundStub(instanceUrl, params); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.configureGetNotFoundStub(blockUrl, params); err != nil {
		return nil, err
	}

	// Delete the security group
	if err := configurator.configureDeleteStub(groupUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted security group
	if err := configurator.configureGetNotFoundStub(groupUrl, params); err != nil {
		return nil, err
	}

	// Delete the nic
	if err := configurator.configureDeleteStub(nicUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted nic
	if err := configurator.configureGetNotFoundStub(nicUrl, params); err != nil {
		return nil, err
	}

	// Delete the public ip
	if err := configurator.configureDeleteStub(publicIpUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted public ip
	if err := configurator.configureGetNotFoundStub(publicIpUrl, params); err != nil {
		return nil, err
	}

	// Delete the subnet
	if err := configurator.configureDeleteStub(subnetUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted subnet
	if err := configurator.configureGetNotFoundStub(subnetUrl, params); err != nil {
		return nil, err
	}

	// Delete the route table
	if err := configurator.configureDeleteStub(routeUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted route table
	if err := configurator.configureGetNotFoundStub(routeUrl, params); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	if err := configurator.configureDeleteStub(gatewayUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted internet gateway
	if err := configurator.configureGetNotFoundStub(gatewayUrl, params); err != nil {
		return nil, err
	}

	// Delete the network
	if err := configurator.configureDeleteStub(networkUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted network
	if err := configurator.configureGetNotFoundStub(networkUrl, params); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(workspaceUrl, params); err != nil {
		return nil, err
	}

	return configurator.client, nil
}
