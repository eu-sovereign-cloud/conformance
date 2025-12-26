//nolint:dupl
package mock

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Authorization

func bulkCreateRolesStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams,
	roleParams []ResourceParams[schema.RoleSpec],
) ([]schema.Role, error) {
	var roles []schema.Role

	for _, role := range roleParams {
		roleUrl := generators.GenerateRoleURL(authorizationProviderV1, baseParams.Tenant, role.Name)

		roleResponse, err := builders.NewRoleBuilder().
			Name(role.Name).
			Provider(authorizationProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).
			Labels(role.InitialLabels).
			Spec(role.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a role
		if err := configurator.configureCreateRoleStub(roleResponse, roleUrl, baseParams); err != nil {
			return nil, err
		}
		roles = append(roles, *roleResponse)
	}

	return roles, nil
}

func bulkCreateRoleAssignmentsStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams,
	roleAssignmentParams []ResourceParams[schema.RoleAssignmentSpec],
) ([]schema.RoleAssignment, error) {
	var assignments []schema.RoleAssignment

	for _, roleAssignment := range roleAssignmentParams {
		roleAssignmentUrl := generators.GenerateRoleAssignmentURL(authorizationProviderV1, baseParams.Tenant, roleAssignment.Name)
		roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
			Name(roleAssignment.Name).
			Provider(authorizationProviderV1).
			ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).
			Labels(roleAssignment.InitialLabels).
			Spec(roleAssignment.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create a role assignment
		if err := configurator.configureCreateRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, baseParams); err != nil {
			return nil, err
		}

		assignments = append(assignments, *roleAssignResponse)
	}

	return assignments, nil
}

// Workspace

func bulkCreateWorkspacesStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams,
	workspaceParams []ResourceParams[schema.WorkspaceSpec],
) ([]schema.Workspace, error) {
	var workspaces []schema.Workspace

	for _, workspace := range workspaceParams {
		url := generators.GenerateWorkspaceURL(workspaceProviderV1, baseParams.Tenant, workspace.Name)

		response, err := builders.NewWorkspaceBuilder().
			Name(workspace.Name).
			Provider(workspaceProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).Region(baseParams.Region).
			Labels(workspace.InitialLabels).
			Build()
		if err != nil {
			return nil, err
		}

		// Create a workspace
		if err := configurator.configureCreateWorkspaceStub(response, url, baseParams); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, *response)
	}
	return workspaces, nil
}

// Compute

func bulkCreateInstancesStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams, workspace string,
	instanceParams []ResourceParams[schema.InstanceSpec],
) ([]schema.Instance, error) {
	var instances []schema.Instance

	for _, instance := range instanceParams {
		instanceUrl := generators.GenerateInstanceURL(computeProviderV1, baseParams.Tenant, workspace, instance.Name)
		instanceResponse, err := builders.NewInstanceBuilder().
			Name(instance.Name).
			Provider(computeProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(instance.InitialLabels).
			Spec(instance.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create an instance
		if err := configurator.configureCreateInstanceStub(instanceResponse, instanceUrl, baseParams); err != nil {
			return nil, err
		}
		instances = append(instances, *instanceResponse)
	}

	return instances, nil
}

// Storage

func bulkCreateBlockStoragesStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams, workspace string,
	blockStorageParams []ResourceParams[schema.BlockStorageSpec],
) ([]schema.BlockStorage, error) {
	var blocks []schema.BlockStorage

	for _, block := range blockStorageParams {
		blockResponse, err := builders.NewBlockStorageBuilder().
			Name(block.Name).
			Provider(storageProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(block.InitialLabels).
			Spec(block.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create a block storage
		blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, baseParams.Tenant, workspace, block.Name)
		if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, baseParams); err != nil {
			return nil, err
		}
		blocks = append(blocks, *blockResponse)
	}

	return blocks, nil
}

func bulkCreateImagesStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams,
	imageParams []ResourceParams[schema.ImageSpec],
) ([]schema.Image, error) {
	var images []schema.Image

	for _, image := range imageParams {
		imageUrl := generators.GenerateImageURL(storageProviderV1, baseParams.Tenant, image.Name)
		imageResponse, err := builders.NewImageBuilder().
			Name(image.Name).
			Provider(storageProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).Region(baseParams.Region).
			Labels(image.InitialLabels).
			Spec(image.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create an image
		if err := configurator.configureCreateImageStub(imageResponse, imageUrl, baseParams); err != nil {
			return nil, err
		}
		images = append(images, *imageResponse)
	}

	return images, nil
}

// Network

func bulkCreateNetworksStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams, workspace string,
	networkParams []ResourceParams[schema.NetworkSpec],
) ([]schema.Network, error) {
	var networks []schema.Network

	for _, network := range networkParams {
		networkUrl := generators.GenerateNetworkURL(networkProviderV1, baseParams.Tenant, workspace, network.Name)

		networkResponse, err := builders.NewNetworkBuilder().
			Name(network.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(network.InitialLabels).
			Spec(network.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a network
		if err := configurator.configureCreateNetworkStub(networkResponse, networkUrl, baseParams); err != nil {
			return nil, err
		}
		networks = append(networks, *networkResponse)
	}

	return networks, nil
}

func bulkCreateInternetGatewaysStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams, workspace string,
	internetGatewayParams []ResourceParams[schema.InternetGatewaySpec],
) ([]schema.InternetGateway, error) {
	var gateways []schema.InternetGateway

	for _, gateway := range internetGatewayParams {
		gatewayUrl := generators.GenerateInternetGatewayURL(networkProviderV1, baseParams.Tenant, workspace, gateway.Name)
		gatewayResponse, err := builders.NewInternetGatewayBuilder().
			Name(gateway.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(gateway.InitialLabels).
			Spec(gateway.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create an internet gateway
		if err := configurator.configureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, baseParams); err != nil {
			return nil, err
		}
		gateways = append(gateways, *gatewayResponse)
	}

	return gateways, nil
}

func bulkCreateRouteTableStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams, workspace, network string,
	routeTableParams []ResourceParams[schema.RouteTableSpec],
) ([]schema.RouteTable, error) {
	var routeTables []schema.RouteTable

	for _, routeTable := range routeTableParams {
		routeTableUrl := generators.GenerateRouteTableURL(networkProviderV1, baseParams.Tenant, workspace, network, routeTable.Name)
		routeTableResponse, err := builders.NewRouteTableBuilder().
			Name(routeTable.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Network(network).Region(baseParams.Region).
			Labels(routeTable.InitialLabels).
			Spec(routeTable.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a route table
		if err := configurator.configureCreateRouteTableStub(routeTableResponse, routeTableUrl, baseParams); err != nil {
			return nil, err
		}
		routeTables = append(routeTables, *routeTableResponse)
	}

	return routeTables, nil
}

func bulkCreateSubnetsStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams, workspace, network string,
	subnetParams []ResourceParams[schema.SubnetSpec],
) ([]schema.Subnet, error) {
	var subnets []schema.Subnet

	for _, subnet := range subnetParams {
		subnetUrl := generators.GenerateSubnetURL(networkProviderV1, baseParams.Tenant, workspace, network, subnet.Name)
		subnetResponse, err := builders.NewSubnetBuilder().
			Name(subnet.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Network(network).Region(baseParams.Region).
			Labels(subnet.InitialLabels).
			Spec(subnet.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a RouteTable
		if err := configurator.configureCreateSubnetStub(subnetResponse, subnetUrl, baseParams); err != nil {
			return nil, err
		}
		subnets = append(subnets, *subnetResponse)
	}

	return subnets, nil
}

func bulkCreatePublicIpsStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams, workspace string,
	publicIpParams []ResourceParams[schema.PublicIpSpec],
) ([]schema.PublicIp, error) {
	var publicIps []schema.PublicIp

	for _, publicIp := range publicIpParams {
		publicIpUrl := generators.GeneratePublicIpURL(networkProviderV1, baseParams.Tenant, workspace, publicIp.Name)
		publicIpResponse, err := builders.NewPublicIpBuilder().
			Name(publicIp.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(publicIp.InitialLabels).
			Spec(publicIp.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a public ip
		if err := configurator.configureCreatePublicIpStub(publicIpResponse, publicIpUrl, baseParams.getBaseParams()); err != nil {
			return nil, err
		}
		publicIps = append(publicIps, *publicIpResponse)
	}

	return publicIps, nil
}

func bulkCreateNicsStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams, workspace string,
	nicParams []ResourceParams[schema.NicSpec],
) ([]schema.Nic, error) {
	var nics []schema.Nic

	for _, nic := range nicParams {
		nicUrl := generators.GenerateNicURL(networkProviderV1, baseParams.Tenant, workspace, nic.Name)
		nicResponse, err := builders.NewNicBuilder().
			Name(nic.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(nic.InitialLabels).
			Spec(nic.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a nic
		if err := configurator.configureCreateNicStub(nicResponse, nicUrl, baseParams.getBaseParams()); err != nil {
			return nil, err
		}
		nics = append(nics, *nicResponse)
	}

	return nics, nil
}

func bulkCreateSecurityGroupsStubV1(configurator *scenarioConfigurator,
	baseParams *BaseParams, workspace string,
	securityGroupParams []ResourceParams[schema.SecurityGroupSpec],
) ([]schema.SecurityGroup, error) {
	var securityGroups []schema.SecurityGroup

	for _, securityGroup := range securityGroupParams {
		securityGroupUrl := generators.GenerateSecurityGroupURL(networkProviderV1, baseParams.Tenant, workspace, securityGroup.Name)
		securityGroupResponse, err := builders.NewSecurityGroupBuilder().
			Name(securityGroup.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(baseParams.Tenant).Workspace(workspace).Region(baseParams.Region).
			Labels(securityGroup.InitialLabels).
			Spec(securityGroup.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		// Create a security group
		if err := configurator.configureCreateSecurityGroupStub(securityGroupResponse, securityGroupUrl, baseParams.getBaseParams()); err != nil {
			return nil, err
		}
		securityGroups = append(securityGroups, *securityGroupResponse)
	}

	return securityGroups, nil
}
