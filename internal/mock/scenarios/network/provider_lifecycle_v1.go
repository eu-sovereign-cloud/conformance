package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureProviderLifecycleScenarioV1(scenario *mockscenarios.Scenario, params params.NetworkProviderLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := params.Workspace
	blockStorage := params.BlockStorage
	instance := params.Instance
	network := params.NetworkInitial
	internetGateway := params.InternetGatewayInitial
	nic := params.NicInitial
	publicIp := params.PublicIpInitial
	routeTable := params.RouteTableInitial
	subnet := params.SubnetInitial
	securityGroup := params.SecurityGroupInitial

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockUrl := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockStorage.Metadata.Tenant, blockStorage.Metadata.Workspace, blockStorage.Metadata.Name)
	instanceUrl := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, instance.Metadata.Tenant, instance.Metadata.Workspace, instance.Metadata.Name)
	networkUrl := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)
	gatewayUrl := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, internetGateway.Metadata.Tenant, internetGateway.Metadata.Workspace, internetGateway.Metadata.Name)
	nicUrl := generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, nic.Metadata.Tenant, nic.Metadata.Workspace, nic.Metadata.Name)
	publicIpUrl := generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, publicIp.Metadata.Tenant, publicIp.Metadata.Workspace, publicIp.Metadata.Name)
	routeUrl := generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)
	subnetUrl := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, subnet.Metadata.Tenant, subnet.Metadata.Workspace, subnet.Metadata.Network, subnet.Metadata.Name)
	groupUrl := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, securityGroup.Metadata.Tenant, securityGroup.Metadata.Workspace, securityGroup.Metadata.Name)

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create a network
	if err := configurator.ConfigureCreateNetworkStub(network, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create an internet gateway
	if err := configurator.ConfigureCreateInternetGatewayStub(internetGateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created internet gateway
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(internetGateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(internetGateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the internet gateway
	internetGateway = params.InternetGatewayUpdated
	if err := configurator.ConfigureUpdateInternetGatewayStub(internetGateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated internet gateway
	if err := configurator.ConfigureGetUpdatingInternetGatewayStub(internetGateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(internetGateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create a route table
	if err := configurator.ConfigureCreateRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created route table
	if err := configurator.ConfigureGetCreatingRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created network
	if err := configurator.ConfigureGetCreatingNetworkStub(network, networkUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(network, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the route table
	routeTable = params.RouteTableUpdated
	if err := configurator.ConfigureUpdateRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated route table
	if err := configurator.ConfigureGetUpdatingRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the network
	network = params.NetworkUpdated
	if err := configurator.ConfigureUpdateNetworkStub(network, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated network
	if err := configurator.ConfigureGetUpdatingNetworkStub(network, networkUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(network, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create a subnet
	if err := configurator.ConfigureCreateSubnetStub(subnet, subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created subnet
	if err := configurator.ConfigureGetCreatingSubnetStub(subnet, subnetUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(subnet, subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the subnet
	subnet = params.SubnetUpdated
	if err := configurator.ConfigureUpdateSubnetStub(subnet, subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated subnet
	if err := configurator.ConfigureGetUpdatingSubnetStub(subnet, subnetUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(subnet, subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create a nic
	if err := configurator.ConfigureCreateNicStub(nic, nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created nic
	if err := configurator.ConfigureGetCreatingNicStub(nic, nicUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNicStub(nic, nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the nic
	nic = params.NicUpdated
	if err := configurator.ConfigureUpdateNicStub(nic, nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated nic
	if err := configurator.ConfigureGetUpdatingNicStub(nic, nicUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNicStub(nic, nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create a public ip
	if err := configurator.ConfigureCreatePublicIpStub(publicIp, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created public ip
	if err := configurator.ConfigureGetCreatingPublicIpStub(publicIp, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActivePublicIpStub(publicIp, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the public ip
	publicIp = params.PublicIpUpdated
	if err := configurator.ConfigureUpdatePublicIpStub(publicIp, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated public ip
	if err := configurator.ConfigureGetUpdatingPublicIpStub(publicIp, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActivePublicIpStub(publicIp, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create a security group
	if err := configurator.ConfigureCreateSecurityGroupStub(securityGroup, groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created security group
	if err := configurator.ConfigureGetCreatingSecurityGroupStub(securityGroup, groupUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupStub(securityGroup, groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the security group
	securityGroup = params.SecurityGroupUpdated
	if err := configurator.ConfigureUpdateSecurityGroupStub(securityGroup, groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated security group
	if err := configurator.ConfigureGetUpdatingSecurityGroupStub(securityGroup, groupUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupStub(securityGroup, groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetCreatingBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create an instance
	if err := configurator.ConfigureCreateInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created instance
	if err := configurator.ConfigureGetCreatingInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}
	// Delete the instance
	if err := configurator.ConfigureDeleteStub(instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted instance
	if err := configurator.ConfigureGetDeletingInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetDeletingBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the security group
	if err := configurator.ConfigureDeleteStub(groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted security group
	if err := configurator.ConfigureGetDeletingSecurityGroupStub(securityGroup, groupUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the nic
	if err := configurator.ConfigureDeleteStub(nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted nic
	if err := configurator.ConfigureGetDeletingNicStub(nic, nicUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the public ip
	if err := configurator.ConfigureDeleteStub(publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted public ip
	if err := configurator.ConfigureGetDeletingPublicIpStub(publicIp, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the subnet
	if err := configurator.ConfigureDeleteStub(subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted subnet
	if err := configurator.ConfigureGetDeletingSubnetStub(subnet, subnetUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the route table
	if err := configurator.ConfigureDeleteStub(routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted route table
	if err := configurator.ConfigureGetDeletingRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the internet gateway
	if err := configurator.ConfigureDeleteStub(gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted internet gateway
	if err := configurator.ConfigureGetDeletingInternetGatewayStub(internetGateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the network
	if err := configurator.ConfigureDeleteStub(networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted network
	if err := configurator.ConfigureGetDeletingNetworkStub(network, networkUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetDeletingWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
