package mockusage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureFoundationScenarioV1(scenario *mockscenarios.Scenario, params params.FoundationUsageV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	role := params.Role
	roleAssignment := params.RoleAssignment
	workspace := params.Workspace
	blockStorage := params.BlockStorage
	image := params.Image
	instance := params.Instance
	network := params.Network
	gateway := params.InternetGateway
	nic := params.Nic
	publicIp := params.PublicIp
	routeTable := params.RouteTable
	subnet := params.Subnet
	securityGroup := params.SecurityGroup

	// Generate URLs
	roleUrl := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, role.Metadata.Tenant, role.Metadata.Name)
	roleAssignUrl := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, roleAssignment.Metadata.Tenant, roleAssignment.Metadata.Name)
	workspaceUrl := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockUrl := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockStorage.Metadata.Tenant, blockStorage.Metadata.Workspace, blockStorage.Metadata.Name)
	imageUrl := generators.GenerateImageURL(sdkconsts.StorageProviderV1Name, image.Metadata.Tenant, image.Metadata.Name)
	instanceUrl := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, instance.Metadata.Tenant, instance.Metadata.Workspace, instance.Metadata.Name)
	networkUrl := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)
	gatewayUrl := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, gateway.Metadata.Tenant, gateway.Metadata.Workspace, gateway.Metadata.Name)
	nicUrl := generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, nic.Metadata.Tenant, nic.Metadata.Workspace, nic.Metadata.Name)
	publicIpUrl := generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, publicIp.Metadata.Tenant, publicIp.Metadata.Workspace, publicIp.Metadata.Name)
	routeUrl := generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)
	subnetUrl := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, subnet.Metadata.Tenant, subnet.Metadata.Workspace, subnet.Metadata.Network, subnet.Metadata.Name)
	groupUrl := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, securityGroup.Metadata.Tenant, securityGroup.Metadata.Workspace, securityGroup.Metadata.Name)

	// Authorization

	// Create a role
	if err := configurator.ConfigureCreateRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created role
	if err := configurator.ConfigureGetCreatingRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create a role assignment
	if err := configurator.ConfigureCreateRoleAssignmentStub(roleAssignment, roleAssignUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created role assignment
	if err := configurator.ConfigureGetCreatingRoleAssignmentStub(roleAssignment, roleAssignUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignment, roleAssignUrl, scenario.MockParams); err != nil {
		return err
	}

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

	// Storage

	// Create an image
	if err := configurator.ConfigureCreateImageStub(image, imageUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created image
	if err := configurator.ConfigureGetCreatingImageStub(image, imageUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveImageStub(image, imageUrl, scenario.MockParams); err != nil {
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

	// Network

	// Create a network
	if err := configurator.ConfigureCreateNetworkStub(network, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create an internet gateway
	if err := configurator.ConfigureCreateInternetGatewayStub(gateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created internet gateway
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(gateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gateway, gatewayUrl, scenario.MockParams); err != nil {
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

	// Compute

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

	// Delete all

	// Delete the instance
	if err := configurator.ConfigureDeleteStub(instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the security Group
	if err := configurator.ConfigureDeleteStub(groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the nic
	if err := configurator.ConfigureDeleteStub(nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the public ip
	if err := configurator.ConfigureDeleteStub(publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the subnet
	if err := configurator.ConfigureDeleteStub(subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the route-table
	if err := configurator.ConfigureDeleteStub(routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the internet gateway
	if err := configurator.ConfigureDeleteStub(gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the network
	if err := configurator.ConfigureDeleteStub(networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the image
	if err := configurator.ConfigureDeleteStub(imageUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the role assignment
	if err := configurator.ConfigureDeleteStub(roleAssignUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the role
	if err := configurator.ConfigureDeleteStub(roleUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
