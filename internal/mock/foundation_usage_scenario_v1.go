package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/wiremock/go-wiremock"
)

func ConfigFoundationUsageScenario(scenario string, params *FoundationUsageParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	roleUrl := generators.GenerateRoleURL(authorizationProviderV1, params.Tenant, params.Role.Name)
	roleAssignUrl := generators.GenerateRoleAssignmentURL(authorizationProviderV1, params.Tenant, params.RoleAssignment.Name)
	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageUrl := generators.GenerateImageURL(storageProviderV1, params.Tenant, params.Image.Name)
	instanceUrl := generators.GenerateInstanceURL(computeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkUrl := generators.GenerateNetworkURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name)
	gatewayUrl := generators.GenerateInternetGatewayURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicUrl := generators.GenerateNicURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Nic.Name)
	publicIpUrl := generators.GeneratePublicIpURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeUrl := generators.GenerateRouteTableURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name, params.RouteTable.Name)
	subnetUrl := generators.GenerateSubnetURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name, params.Subnet.Name)
	groupUrl := generators.GenerateSecurityGroupURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Authorization

	// Role
	roleResponse, err := builders.NewRoleBuilder().
		Name(params.Role.Name).
		Provider(authorizationProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).
		Spec(params.Role.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a role
	if err := configurator.configureCreateRoleStub(roleResponse, roleUrl, params); err != nil {
		return nil, err
	}

	// Get the created role
	if err := configurator.configureGetActiveRoleStub(roleResponse, roleUrl, params); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name(params.RoleAssignment.Name).
		Provider(authorizationProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).
		Spec(params.RoleAssignment.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create a role assignment
	if err := configurator.configureCreateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params); err != nil {
		return nil, err
	}

	// Get the created role assignment
	if err := configurator.configureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params); err != nil {
		return nil, err
	}

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
	if err := configurator.configureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.configureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
		return nil, err
	}

	// Storage

	// Image
	imageResponse, err := builders.NewImageBuilder().
		Name(params.Image.Name).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Spec(params.Image.InitialSpec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	// Create an image
	if err := configurator.configureCreateImageStub(imageResponse, imageUrl, params); err != nil {
		return nil, err
	}

	// Get the created image
	if err := configurator.configureGetActiveImageStub(imageResponse, imageUrl, params); err != nil {
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
	if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.configureGetActiveBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Network

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
	if err := configurator.configureCreateNetworkStub(networkResponse, networkUrl, params); err != nil {
		return nil, err
	}

	// Get the created network
	if err := configurator.configureGetActiveNetworkStub(networkResponse, networkUrl, params); err != nil {
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
	if err := configurator.configureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, params); err != nil {
		return nil, err
	}

	// Get the created internet gateway
	if err := configurator.configureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, params); err != nil {
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
	if err := configurator.configureCreateRouteTableStub(routeResponse, routeUrl, params); err != nil {
		return nil, err
	}

	// Get the created route table
	if err := configurator.configureGetActiveRouteTableStub(routeResponse, routeUrl, params); err != nil {
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
	if err := configurator.configureCreateSubnetStub(subnetResponse, subnetUrl, params); err != nil {
		return nil, err
	}

	// Get the created subnet
	if err := configurator.configureGetActiveSubnetStub(subnetResponse, subnetUrl, params); err != nil {
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
	if err := configurator.configureCreateSecurityGroupStub(groupResponse, groupUrl, params); err != nil {
		return nil, err
	}

	// Get the created security group
	if err := configurator.configureGetActiveSecurityGroupStub(groupResponse, groupUrl, params); err != nil {
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
	if err := configurator.configureCreatePublicIpStub(publicIpResponse, publicIpUrl, params); err != nil {
		return nil, err
	}

	// Get the created public ip
	if err := configurator.configureGetActivePublicIpStub(publicIpResponse, publicIpUrl, params); err != nil {
		return nil, err
	}

	// NIC
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
	if err := configurator.configureCreateNicStub(nicResponse, nicUrl, params); err != nil {
		return nil, err
	}

	// Get the created nic
	if err := configurator.configureGetActiveNicStub(nicResponse, nicUrl, params); err != nil {
		return nil, err
	}

	// Compute

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
	if err := configurator.configureCreateInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Delete all

	// Delete the instance
	if err := configurator.configureDeleteStub(instanceUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the security Group
	if err := configurator.configureDeleteStub(groupUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the nic
	if err := configurator.configureDeleteStub(nicUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the public ip
	if err := configurator.configureDeleteStub(publicIpUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the subnet
	if err := configurator.configureDeleteStub(subnetUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the route-table
	if err := configurator.configureDeleteStub(routeUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	if err := configurator.configureDeleteStub(gatewayUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the network
	if err := configurator.configureDeleteStub(networkUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the image
	if err := configurator.configureDeleteStub(imageUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the role assignment
	if err := configurator.configureDeleteStub(roleAssignUrl, params, false); err != nil {
		return nil, err
	}

	// Delete the role
	if err := configurator.configureDeleteStub(roleUrl, params, true); err != nil {
		return nil, err
	}

	return configurator.client, err
}
