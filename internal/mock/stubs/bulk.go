//nolint:dupl
package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Authorization

// TODO Find a better package to it
func BulkCreateRolesStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	rolesParams []schema.Role,
) error {
	for _, role := range rolesParams {
		roleUrl := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, role.Metadata.Tenant, role.Metadata.Name)

		roleResponse, err := builders.NewRoleBuilder().
			Name(role.Metadata.Name).
			Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(role.Metadata.Tenant).
			Labels(role.Labels).
			Spec(&role.Spec).
			Build()
		if err != nil {
			return err
		}
		// Create a role
		if err := configurator.ConfigureCreateRoleStub(roleResponse, roleUrl, mockParams); err != nil {
			return err
		}
	}

	return nil
}

func BulkCreateRoleAssignmentsStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	roleAssignmentParams []schema.RoleAssignment,
) error {
	for _, roleAssignment := range roleAssignmentParams {
		roleAssignmentUrl := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, roleAssignment.Metadata.Tenant, roleAssignment.Metadata.Name)
		roleAssignResponse, err := builders.NewRoleAssignmentBuilder().
			Name(roleAssignment.Metadata.Name).
			Provider(sdkconsts.AuthorizationProviderV1Name).
			ApiVersion(sdkconsts.ApiVersion1).
			Tenant(roleAssignment.Metadata.Tenant).
			Labels(roleAssignment.Labels).
			Spec(&roleAssignment.Spec).
			Build()
		if err != nil {
			return err
		}

		// Create a role assignment
		if err := configurator.ConfigureCreateRoleAssignmentStub(roleAssignResponse, roleAssignmentUrl, mockParams); err != nil {
			return err
		}
	}

	return nil
}

// Workspace

func BulkCreateWorkspacesStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	workspaceParams []schema.Workspace,
) error {
	for _, workspace := range workspaceParams {
		url := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)

		response, err := builders.NewWorkspaceBuilder().
			Name(workspace.Metadata.Name).
			Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(workspace.Metadata.Tenant).Region(workspace.Metadata.Region).
			Labels(workspace.Labels).
			Build()
		if err != nil {
			return err
		}

		// Create a workspace
		if err := configurator.ConfigureCreateWorkspaceStub(response, url, mockParams); err != nil {
			return err
		}
	}
	return nil
}

// Compute

func BulkCreateInstancesStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	instanceParams []schema.Instance,
) error {
	for _, instance := range instanceParams {
		instanceUrl := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, instance.Metadata.Tenant, instance.Metadata.Workspace, instance.Metadata.Name)
		instanceResponse, err := builders.NewInstanceBuilder().
			Name(instance.Metadata.Name).
			Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(instance.Metadata.Tenant).Workspace(instance.Metadata.Workspace).Region(instance.Metadata.Region).
			Labels(instance.Labels).
			Spec(&instance.Spec).
			Build()
		if err != nil {
			return err
		}

		// Create an instance
		if err := configurator.ConfigureCreateInstanceStub(instanceResponse, instanceUrl, mockParams); err != nil {
			return err
		}
	}

	return nil
}

// Storage

func BulkCreateBlockStoragesStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	blockStorageParams []schema.BlockStorage,
) error {
	for _, block := range blockStorageParams {
		blockResponse, err := builders.NewBlockStorageBuilder().
			Name(block.Metadata.Name).
			Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(block.Metadata.Tenant).Workspace(block.Metadata.Workspace).Region(block.Metadata.Region).
			Labels(block.Labels).
			Spec(&block.Spec).
			Build()
		if err != nil {
			return err
		}

		// Create a block storage
		blockUrl := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, block.Metadata.Tenant, block.Metadata.Workspace, block.Metadata.Name)
		if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, mockParams); err != nil {
			return err
		}

	}

	return nil
}

func BulkCreateImagesStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	imageParams []schema.Image,
) error {
	for _, image := range imageParams {
		imageUrl := generators.GenerateImageURL(sdkconsts.StorageProviderV1Name, image.Metadata.Tenant, image.Metadata.Name)
		imageResponse, err := builders.NewImageBuilder().
			Name(image.Metadata.Name).
			Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(image.Metadata.Tenant).Region(image.Metadata.Region).
			Labels(image.Labels).
			Spec(&image.Spec).
			Build()
		if err != nil {
			return err
		}

		// Create an image
		if err := configurator.ConfigureCreateImageStub(imageResponse, imageUrl, mockParams); err != nil {
			return err
		}
	}

	return nil
}

// Network

func BulkCreateNetworksStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	networkParams []schema.Network,
) error {
	for _, network := range networkParams {
		networkUrl := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)

		networkResponse, err := builders.NewNetworkBuilder().
			Name(network.Metadata.Name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(network.Metadata.Tenant).Workspace(network.Metadata.Workspace).Region(network.Metadata.Region).
			Labels(network.Labels).
			Spec(&network.Spec).
			Build()
		if err != nil {
			return err
		}
		// Create a network
		if err := configurator.ConfigureCreateNetworkStub(networkResponse, networkUrl, mockParams); err != nil {
			return err
		}
	}

	return nil
}

func BulkCreateInternetGatewaysStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	internetGatewayParams []schema.InternetGateway,
) error {
	for _, gateway := range internetGatewayParams {
		gatewayUrl := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, gateway.Metadata.Tenant, gateway.Metadata.Workspace, gateway.Metadata.Name)
		gatewayResponse, err := builders.NewInternetGatewayBuilder().
			Name(gateway.Metadata.Name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(gateway.Metadata.Tenant).Workspace(gateway.Metadata.Workspace).Region(gateway.Metadata.Region).
			Labels(gateway.Labels).
			Spec(&gateway.Spec).
			Build()
		if err != nil {
			return err
		}

		// Create an internet gateway
		if err := configurator.ConfigureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, mockParams); err != nil {
			return err
		}

	}
	return nil
}

func BulkCreateRouteTableStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	routeTableParams []schema.RouteTable,
) error {
	for _, routeTable := range routeTableParams {
		routeTableUrl := generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)
		routeTableResponse, err := builders.NewRouteTableBuilder().
			Name(routeTable.Metadata.Name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(routeTable.Metadata.Tenant).Workspace(routeTable.Metadata.Workspace).Network(routeTable.Metadata.Network).Region(routeTable.Metadata.Region).
			Labels(routeTable.Labels).
			Spec(&routeTable.Spec).
			Build()
		if err != nil {
			return err
		}
		// Create a route table
		if err := configurator.ConfigureCreateRouteTableStub(routeTableResponse, routeTableUrl, mockParams); err != nil {
			return err
		}
	}

	return nil
}

func BulkCreateSubnetsStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	subnetParams []schema.Subnet,
) error {
	for _, subnet := range subnetParams {
		subnetUrl := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, subnet.Metadata.Tenant, subnet.Metadata.Workspace, subnet.Metadata.Network, subnet.Metadata.Name)
		subnetResponse, err := builders.NewSubnetBuilder().
			Name(subnet.Metadata.Name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(subnet.Metadata.Tenant).Workspace(subnet.Metadata.Workspace).Network(subnet.Metadata.Network).Region(subnet.Metadata.Region).
			Labels(subnet.Labels).
			Spec(&subnet.Spec).
			Build()
		if err != nil {
			return err
		}
		// Create a RouteTable
		if err := configurator.ConfigureCreateSubnetStub(subnetResponse, subnetUrl, mockParams); err != nil {
			return err
		}
	}

	return nil
}

func BulkCreatePublicIpsStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	publicIpParams []schema.PublicIp,
) error {
	for _, publicIp := range publicIpParams {
		publicIpUrl := generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, publicIp.Metadata.Tenant, publicIp.Metadata.Workspace, publicIp.Metadata.Name)
		publicIpResponse, err := builders.NewPublicIpBuilder().
			Name(publicIp.Metadata.Name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(publicIp.Metadata.Tenant).Workspace(publicIp.Metadata.Workspace).Region(publicIp.Metadata.Region).
			Labels(publicIp.Labels).
			Spec(&publicIp.Spec).
			Build()
		if err != nil {
			return err
		}
		// Create a public ip
		if err := configurator.ConfigureCreatePublicIpStub(publicIpResponse, publicIpUrl, mockParams); err != nil {
			return err
		}
	}

	return nil
}

func BulkCreateNicsStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	nicParams []schema.Nic,
) error {
	for _, nic := range nicParams {
		nicUrl := generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, nic.Metadata.Tenant, nic.Metadata.Workspace, nic.Metadata.Name)
		nicResponse, err := builders.NewNicBuilder().
			Name(nic.Metadata.Name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(nic.Metadata.Tenant).Workspace(nic.Metadata.Workspace).Region(nic.Metadata.Region).
			Labels(nic.Labels).
			Spec(&nic.Spec).
			Build()
		if err != nil {
			return err
		}
		// Create a nic
		if err := configurator.ConfigureCreateNicStub(nicResponse, nicUrl, mockParams); err != nil {
			return err
		}
	}

	return nil
}

func BulkCreateSecurityGroupsStubV1(configurator *Configurator,
	mockParams *mock.MockParams,
	securityGroupParams []schema.SecurityGroup,
) error {
	for _, securityGroup := range securityGroupParams {
		securityGroupUrl := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, securityGroup.Metadata.Tenant, securityGroup.Metadata.Workspace, securityGroup.Metadata.Name)
		securityGroupResponse, err := builders.NewSecurityGroupBuilder().
			Name(securityGroup.Metadata.Name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(securityGroup.Metadata.Tenant).Workspace(securityGroup.Metadata.Workspace).Region(securityGroup.Metadata.Region).
			Labels(securityGroup.Labels).
			Spec(&securityGroup.Spec).
			Build()
		if err != nil {
			return err
		}

		// Create a security group
		if err := configurator.ConfigureCreateSecurityGroupStub(securityGroupResponse, securityGroupUrl, mockParams); err != nil {
			return err
		}
	}

	return nil
}
