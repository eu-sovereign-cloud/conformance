//nolint:dupl
package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// TODO Review if the right package of these methods are stubs

// Authorization

func BulkCreateRolesStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams,
	roleParams []mock.ResourceParams[schema.RoleSpec],
) ([]schema.Role, error) {
	var roles []schema.Role

	for _, role := range roleParams {
		roleUrl := generators.GenerateRoleURL(constants.AuthorizationProviderV1, baseParams.Tenant, role.Name)

		roleResponse, err := builders.NewRoleBuilder().
			Name(role.Name).
			Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).
			Labels(role.InitialLabels).
			Spec(role.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a role
		if err := configurator.ConfigureCreateRoleStub(roleResponse, roleUrl, baseParams); err != nil {
			return nil, err
		}
		roles = append(roles, *roleResponse)
	}

	return roles, nil
}

func BulkCreateRoleAssignmentsStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams,
	roleAssignmentParams []mock.ResourceParams[schema.RoleAssignmentSpec],
) ([]schema.RoleAssignment, error) {
	var assignments []schema.RoleAssignment

	for _, roleAssignment := range roleAssignmentParams {
		roleAssignmentUrl := generators.GenerateRoleAssignmentURL(constants.AuthorizationProviderV1, baseParams.Tenant, roleAssignment.Name)
		roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
			Name(roleAssignment.Name).
			Provider(constants.AuthorizationProviderV1).
			ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).
			Labels(roleAssignment.InitialLabels).
			Spec(roleAssignment.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create a role assignment
		if err := configurator.ConfigureCreateRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, baseParams); err != nil {
			return nil, err
		}

		assignments = append(assignments, *roleAssignResponse)
	}

	return assignments, nil
}

// Workspace

func BulkCreateWorkspacesStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams,
	workspaceParams []mock.ResourceParams[schema.WorkspaceSpec],
) ([]schema.Workspace, error) {
	var workspaces []schema.Workspace

	for _, workspace := range workspaceParams {
		url := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, baseParams.Tenant, workspace.Name)

		response, err := builders.NewWorkspaceBuilder().
			Name(workspace.Name).
			Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).Region(baseParams.Region).
			Labels(workspace.InitialLabels).
			Build()
		if err != nil {
			return nil, err
		}

		// Create a workspace
		if err := configurator.ConfigureCreateWorkspaceStub(response, url, baseParams); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, *response)
	}
	return workspaces, nil
}

// Compute

func BulkCreateInstancesStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams, workspace string,
	instanceParams []mock.ResourceParams[schema.InstanceSpec],
) ([]schema.Instance, error) {
	var instances []schema.Instance

	for _, instance := range instanceParams {
		instanceUrl := generators.GenerateInstanceURL(constants.ComputeProviderV1, baseParams.Tenant, workspace, instance.Name)
		instanceResponse, err := builders.NewInstanceBuilder().
			Name(instance.Name).
			Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(instance.InitialLabels).
			Spec(instance.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create an instance
		if err := configurator.ConfigureCreateInstanceStub(instanceResponse, instanceUrl, baseParams); err != nil {
			return nil, err
		}
		instances = append(instances, *instanceResponse)
	}

	return instances, nil
}

// Storage

func BulkCreateBlockStoragesStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams, workspace string,
	blockStorageParams []mock.ResourceParams[schema.BlockStorageSpec],
) ([]schema.BlockStorage, error) {
	var blocks []schema.BlockStorage

	for _, block := range blockStorageParams {
		blockResponse, err := builders.NewBlockStorageBuilder().
			Name(block.Name).
			Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(block.InitialLabels).
			Spec(block.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create a block storage
		blockUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, baseParams.Tenant, workspace, block.Name)
		if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, baseParams); err != nil {
			return nil, err
		}
		blocks = append(blocks, *blockResponse)
	}

	return blocks, nil
}

func BulkCreateImagesStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams,
	imageParams []mock.ResourceParams[schema.ImageSpec],
) ([]schema.Image, error) {
	var images []schema.Image

	for _, image := range imageParams {
		imageUrl := generators.GenerateImageURL(constants.StorageProviderV1, baseParams.Tenant, image.Name)
		imageResponse, err := builders.NewImageBuilder().
			Name(image.Name).
			Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).Region(baseParams.Region).
			Labels(image.InitialLabels).
			Spec(image.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create an image
		if err := configurator.ConfigureCreateImageStub(imageResponse, imageUrl, baseParams); err != nil {
			return nil, err
		}
		images = append(images, *imageResponse)
	}

	return images, nil
}

// Network

func BulkCreateNetworksStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams, workspace string,
	networkParams []mock.ResourceParams[schema.NetworkSpec],
) ([]schema.Network, error) {
	var networks []schema.Network

	for _, network := range networkParams {
		networkUrl := generators.GenerateNetworkURL(constants.NetworkProviderV1, baseParams.Tenant, workspace, network.Name)

		networkResponse, err := builders.NewNetworkBuilder().
			Name(network.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(network.InitialLabels).
			Spec(network.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a network
		if err := configurator.ConfigureCreateNetworkStub(networkResponse, networkUrl, baseParams); err != nil {
			return nil, err
		}
		networks = append(networks, *networkResponse)
	}

	return networks, nil
}

func BulkCreateInternetGatewaysStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams, workspace string,
	internetGatewayParams []mock.ResourceParams[schema.InternetGatewaySpec],
) ([]schema.InternetGateway, error) {
	var gateways []schema.InternetGateway

	for _, gateway := range internetGatewayParams {
		gatewayUrl := generators.GenerateInternetGatewayURL(constants.NetworkProviderV1, baseParams.Tenant, workspace, gateway.Name)
		gatewayResponse, err := builders.NewInternetGatewayBuilder().
			Name(gateway.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(gateway.InitialLabels).
			Spec(gateway.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create an internet gateway
		if err := configurator.ConfigureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, baseParams); err != nil {
			return nil, err
		}
		gateways = append(gateways, *gatewayResponse)
	}

	return gateways, nil
}

func BulkCreateRouteTableStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams, workspace, network string,
	routeTableParams []mock.ResourceParams[schema.RouteTableSpec],
) ([]schema.RouteTable, error) {
	var routeTables []schema.RouteTable

	for _, routeTable := range routeTableParams {
		routeTableUrl := generators.GenerateRouteTableURL(constants.NetworkProviderV1, baseParams.Tenant, workspace, network, routeTable.Name)
		routeTableResponse, err := builders.NewRouteTableBuilder().
			Name(routeTable.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Network(network).Region(baseParams.Region).
			Labels(routeTable.InitialLabels).
			Spec(routeTable.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a route table
		if err := configurator.ConfigureCreateRouteTableStub(routeTableResponse, routeTableUrl, baseParams); err != nil {
			return nil, err
		}
		routeTables = append(routeTables, *routeTableResponse)
	}

	return routeTables, nil
}

func BulkCreateSubnetsStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams, workspace, network string,
	subnetParams []mock.ResourceParams[schema.SubnetSpec],
) ([]schema.Subnet, error) {
	var subnets []schema.Subnet

	for _, subnet := range subnetParams {
		subnetUrl := generators.GenerateSubnetURL(constants.NetworkProviderV1, baseParams.Tenant, workspace, network, subnet.Name)
		subnetResponse, err := builders.NewSubnetBuilder().
			Name(subnet.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Network(network).Region(baseParams.Region).
			Labels(subnet.InitialLabels).
			Spec(subnet.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a RouteTable
		if err := configurator.ConfigureCreateSubnetStub(subnetResponse, subnetUrl, baseParams); err != nil {
			return nil, err
		}
		subnets = append(subnets, *subnetResponse)
	}

	return subnets, nil
}

func BulkCreatePublicIpsStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams, workspace string,
	publicIpParams []mock.ResourceParams[schema.PublicIpSpec],
) ([]schema.PublicIp, error) {
	var publicIps []schema.PublicIp

	for _, publicIp := range publicIpParams {
		publicIpUrl := generators.GeneratePublicIpURL(constants.NetworkProviderV1, baseParams.Tenant, workspace, publicIp.Name)
		publicIpResponse, err := builders.NewPublicIpBuilder().
			Name(publicIp.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(publicIp.InitialLabels).
			Spec(publicIp.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a public ip
		if err := configurator.ConfigureCreatePublicIpStub(publicIpResponse, publicIpUrl, baseParams.GetBaseParams()); err != nil {
			return nil, err
		}
		publicIps = append(publicIps, *publicIpResponse)
	}

	return publicIps, nil
}

func BulkCreateNicsStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams, workspace string,
	nicParams []mock.ResourceParams[schema.NicSpec],
) ([]schema.Nic, error) {
	var nics []schema.Nic

	for _, nic := range nicParams {
		nicUrl := generators.GenerateNicURL(constants.NetworkProviderV1, baseParams.Tenant, workspace, nic.Name)
		nicResponse, err := builders.NewNicBuilder().
			Name(nic.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(nic.InitialLabels).
			Spec(nic.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a nic
		if err := configurator.ConfigureCreateNicStub(nicResponse, nicUrl, baseParams.GetBaseParams()); err != nil {
			return nil, err
		}
		nics = append(nics, *nicResponse)
	}

	return nics, nil
}

func BulkCreateSecurityGroupsStubV1(configurator *stubConfigurator,
	baseParams *mock.BaseParams, workspace string,
	securityGroupParams []mock.ResourceParams[schema.SecurityGroupSpec],
) ([]schema.SecurityGroup, error) {
	var securityGroups []schema.SecurityGroup

	for _, securityGroup := range securityGroupParams {
		securityGroupUrl := generators.GenerateSecurityGroupURL(constants.NetworkProviderV1, baseParams.Tenant, workspace, securityGroup.Name)
		securityGroupResponse, err := builders.NewSecurityGroupBuilder().
			Name(securityGroup.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(securityGroup.InitialLabels).
			Spec(securityGroup.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create a security group
		if err := configurator.ConfigureCreateSecurityGroupStub(securityGroupResponse, securityGroupUrl, baseParams.GetBaseParams()); err != nil {
			return nil, err
		}
		securityGroups = append(securityGroups, *securityGroupResponse)
	}

	return securityGroups, nil
}
