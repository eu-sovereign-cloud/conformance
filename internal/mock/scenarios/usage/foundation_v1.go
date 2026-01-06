package usage

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/wiremock/go-wiremock"
)

func ConfigureFoundationScenarioV1(scenario string, params *mock.FoundationUsageParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := stubs.NewScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	roleUrl := generators.GenerateRoleURL(constants.AuthorizationProviderV1, params.Tenant, params.Role.Name)
	roleAssignUrl := generators.GenerateRoleAssignmentURL(constants.AuthorizationProviderV1, params.Tenant, params.RoleAssignment.Name)
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageUrl := generators.GenerateImageURL(constants.StorageProviderV1, params.Tenant, params.Image.Name)
	instanceUrl := generators.GenerateInstanceURL(constants.ComputeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkUrl := generators.GenerateNetworkURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name)
	gatewayUrl := generators.GenerateInternetGatewayURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicUrl := generators.GenerateNicURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.Nic.Name)
	publicIpUrl := generators.GeneratePublicIpURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeUrl := generators.GenerateRouteTableURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name, params.RouteTable.Name)
	subnetUrl := generators.GenerateSubnetURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name, params.Subnet.Name)
	groupUrl := generators.GenerateSecurityGroupURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Authorization

	// Role
	roleResponse, err := builders.NewRoleBuilder().
		Name(params.Role.Name).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).
		Spec(params.Role.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a role
	if err := configurator.ConfigureCreateRoleStub(roleResponse, roleUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created role
	if err := configurator.ConfigureGetActiveRoleStub(roleResponse, roleUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Role assignment
	roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
		Name(params.RoleAssignment.Name).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).
		Spec(params.RoleAssignment.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a role assignment
	if err := configurator.ConfigureCreateRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created role assignment
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignResponse, roleAssignUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

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
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Storage

	// Image
	imageResponse, err := builders.NewImageBuilder().
		Name(params.Image.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Spec(params.Image.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an image
	if err := configurator.ConfigureCreateImageStub(imageResponse, imageUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created image
	if err := configurator.ConfigureGetActiveImageStub(imageResponse, imageUrl, params.GetBaseParams()); err != nil {
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
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Network

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
	if err := configurator.ConfigureCreateNetworkStub(networkResponse, networkUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created network
	if err := configurator.ConfigureGetActiveNetworkStub(networkResponse, networkUrl, params.GetBaseParams()); err != nil {
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
	if err := configurator.ConfigureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created internet gateway
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, params.GetBaseParams()); err != nil {
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
	if err := configurator.ConfigureCreateRouteTableStub(routeResponse, routeUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created route table
	if err := configurator.ConfigureGetActiveRouteTableStub(routeResponse, routeUrl, params.GetBaseParams()); err != nil {
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
	if err := configurator.ConfigureCreateSubnetStub(subnetResponse, subnetUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created subnet
	if err := configurator.ConfigureGetActiveSubnetStub(subnetResponse, subnetUrl, params.GetBaseParams()); err != nil {
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
	if err := configurator.ConfigureCreateSecurityGroupStub(groupResponse, groupUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created security group
	if err := configurator.ConfigureGetActiveSecurityGroupStub(groupResponse, groupUrl, params.GetBaseParams()); err != nil {
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
	if err := configurator.ConfigureCreatePublicIpStub(publicIpResponse, publicIpUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created public ip
	if err := configurator.ConfigureGetActivePublicIpStub(publicIpResponse, publicIpUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// NIC
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
	if err := configurator.ConfigureCreateNicStub(nicResponse, nicUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created nic
	if err := configurator.ConfigureGetActiveNicStub(nicResponse, nicUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Compute

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
	if err := configurator.ConfigureCreateInstanceStub(instanceResponse, instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.ConfigureGetActiveInstanceStub(instanceResponse, instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete all

	// Delete the instance
	if err := configurator.ConfigureDeleteStub(instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the security Group
	if err := configurator.ConfigureDeleteStub(groupUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the nic
	if err := configurator.ConfigureDeleteStub(nicUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the public ip
	if err := configurator.ConfigureDeleteStub(publicIpUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the subnet
	if err := configurator.ConfigureDeleteStub(subnetUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the route-table
	if err := configurator.ConfigureDeleteStub(routeUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	if err := configurator.ConfigureDeleteStub(gatewayUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the network
	if err := configurator.ConfigureDeleteStub(networkUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the image
	if err := configurator.ConfigureDeleteStub(imageUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the role assignment
	if err := configurator.ConfigureDeleteStub(roleAssignUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the role
	if err := configurator.ConfigureDeleteStub(roleUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	return configurator.Client, err
}
