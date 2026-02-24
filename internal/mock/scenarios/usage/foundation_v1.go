package mockusage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureFoundationScenarioV1(scenario *mockscenarios.Scenario, params *params.FoundationUsageV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	role := *params.Role
	roleAssignment := *params.RoleAssignment
	workspace := *params.Workspace
	blockStorage := *params.BlockStorage
	image := *params.Image
	instance := *params.Instance
	network := *params.Network
	gateway := *params.InternetGateway
	nic := *params.Nic
	publicIp := *params.PublicIp
	routeTable := *params.RouteTable
	subnet := *params.Subnet
	securityGroup := *params.SecurityGroup

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

	// Role
	roleResponse, err := builders.NewRoleBuilder().
		Name(role.Metadata.Name).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(role.Metadata.Tenant).
		Spec(&role.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a role
	if err := configurator.ConfigureCreateRoleStub(roleResponse, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created role
	if err := configurator.ConfigureGetCreatingRoleStub(roleResponse, roleUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleStub(roleResponse, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	// Role assignment
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignment.Metadata.Name).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(roleAssignment.Metadata.Tenant).
		Spec(&roleAssignment.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a role assignment
	if err := configurator.ConfigureCreateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created role assignment
	if err := configurator.ConfigureGetCreatingRoleAssignmentStub(roleAssignResponse, roleAssignUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, scenario.MockParams); err != nil {
		return err
	}

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(workspace.Metadata.Name).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(workspace.Metadata.Tenant).Region(workspace.Metadata.Region).
		Labels(workspace.Labels).
		Build()
	if err != nil {
		return err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspaceResponse, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Storage

	// Image
	imageResponse, err := builders.NewImageBuilder().
		Name(image.Metadata.Name).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(image.Metadata.Tenant).Region(image.Metadata.Region).
		Spec(&image.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create an image
	if err := configurator.ConfigureCreateImageStub(imageResponse, imageUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created image
	if err := configurator.ConfigureGetCreatingImageStub(imageResponse, imageUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveImageStub(imageResponse, imageUrl, scenario.MockParams); err != nil {
		return err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(blockStorage.Metadata.Name).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(blockStorage.Metadata.Tenant).Workspace(blockStorage.Metadata.Workspace).Region(blockStorage.Metadata.Region).
		Spec(&blockStorage.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetCreatingBlockStorageStub(blockResponse, blockUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Network

	// Network
	networkResponse, err := builders.NewNetworkBuilder().
		Name(network.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(network.Metadata.Tenant).Workspace(network.Metadata.Workspace).Region(network.Metadata.Region).
		Spec(&network.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a network
	if err := configurator.ConfigureCreateNetworkStub(networkResponse, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created network
	if err := configurator.ConfigureGetCreatingNetworkStub(networkResponse, networkUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(networkResponse, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Internet gateway
	gatewayResponse, err := builders.NewInternetGatewayBuilder().
		Name(gateway.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(gateway.Metadata.Tenant).Workspace(gateway.Metadata.Workspace).Region(gateway.Metadata.Region).
		Spec(&gateway.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create an internet gateway
	if err := configurator.ConfigureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created internet gateway
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(gatewayResponse, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Route table
	routeResponse, err := builders.NewRouteTableBuilder().
		Name(routeTable.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(routeTable.Metadata.Tenant).Workspace(routeTable.Metadata.Workspace).Network(routeTable.Metadata.Network).Region(routeTable.Metadata.Region).
		Spec(&routeTable.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a route table
	if err := configurator.ConfigureCreateRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created route table
	if err := configurator.ConfigureGetCreatingRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Subnet
	subnetResponse, err := builders.NewSubnetBuilder().
		Name(subnet.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(subnet.Metadata.Tenant).Workspace(subnet.Metadata.Workspace).Network(subnet.Metadata.Network).Region(subnet.Metadata.Region).
		Spec(&subnet.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a subnet
	if err := configurator.ConfigureCreateSubnetStub(subnetResponse, subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created subnet
	if err := configurator.ConfigureGetCreatingSubnetStub(subnetResponse, subnetUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(subnetResponse, subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Security group
	groupResponse, err := builders.NewSecurityGroupBuilder().
		Name(securityGroup.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(securityGroup.Metadata.Tenant).Workspace(securityGroup.Metadata.Workspace).Region(securityGroup.Metadata.Region).
		Spec(&securityGroup.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a security group
	if err := configurator.ConfigureCreateSecurityGroupStub(groupResponse, groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created security group
	if err := configurator.ConfigureGetCreatingSecurityGroupStub(groupResponse, groupUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupStub(groupResponse, groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Public ip
	publicIpResponse, err := builders.NewPublicIpBuilder().
		Name(publicIp.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(publicIp.Metadata.Tenant).Workspace(publicIp.Metadata.Workspace).Region(publicIp.Metadata.Region).
		Spec(&publicIp.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a public ip
	if err := configurator.ConfigureCreatePublicIpStub(publicIpResponse, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created public ip
	if err := configurator.ConfigureGetCreatingPublicIpStub(publicIpResponse, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActivePublicIpStub(publicIpResponse, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// NIC
	nicResponse, err := builders.NewNicBuilder().
		Name(nic.Metadata.Name).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(nic.Metadata.Tenant).Workspace(nic.Metadata.Workspace).Region(nic.Metadata.Region).
		Spec(&nic.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a nic
	if err := configurator.ConfigureCreateNicStub(nicResponse, nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created nic
	if err := configurator.ConfigureGetCreatingNicStub(nicResponse, nicUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNicStub(nicResponse, nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Compute

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(instance.Metadata.Name).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(instance.Metadata.Tenant).Workspace(instance.Metadata.Workspace).Region(instance.Metadata.Region).
		Spec(&instance.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create an instance
	if err := configurator.ConfigureCreateInstanceStub(instanceResponse, instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created instance
	if err := configurator.ConfigureGetCreatingInstanceStub(instanceResponse, instanceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInstanceStub(instanceResponse, instanceUrl, scenario.MockParams); err != nil {
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
