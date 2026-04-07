package steps

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Authorization

func BulkCreateRolesStepsV1(configurator *StepsConfigurator, suite suites.GlobalTestSuite, stepName string, roles []schema.Role) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, role := range roles {
			expectRoleMeta := role.Metadata
			expectRoleSpec := role.Spec

			// Create a role
			configurator.CreateOrUpdateRoleV1Step("Create a role", sCtx, suite.Client.AuthorizationV1, &role,
				StepResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
					Metadata:       expectRoleMeta,
					Spec:           &expectRoleSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteRolesStepsV1(configurator *StepsConfigurator, suite suites.GlobalTestSuite, stepName string, roles []schema.Role) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, role := range roles {
			configurator.DeleteRoleV1Step("Delete the role", sCtx, suite.Client.AuthorizationV1, &role)

			// Get the deleted role
			roleTRef := secapi.TenantReference{
				Tenant: secapi.TenantID(suite.Tenant),
				Name:   role.Metadata.Name,
			}
			configurator.WatchRoleUntilDeletedV1Step("Watch the role deletion", sCtx, suite.Client.AuthorizationV1, roleTRef)
		}
	})
}

func BulkCreateRoleAssignmentsStepsV1(configurator *StepsConfigurator, suite suites.GlobalTestSuite, stepName string, roleAssignments []schema.RoleAssignment) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, roleAssign := range roleAssignments {
			expectRoleAssignMeta := roleAssign.Metadata
			expectRoleAssignSpec := &roleAssign.Spec

			// Create a role assignment
			configurator.CreateOrUpdateRoleAssignmentV1Step("Create a role assignment", sCtx, suite.Client.AuthorizationV1, &roleAssign,
				StepResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
					Metadata:       expectRoleAssignMeta,
					Spec:           expectRoleAssignSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteRoleAssignmentsStepsV1(configurator *StepsConfigurator, suite suites.GlobalTestSuite, stepName string, roleAssignments []schema.RoleAssignment) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, roleAssign := range roleAssignments {
			configurator.DeleteRoleAssignmentV1Step("Delete the role assignment", sCtx, suite.Client.AuthorizationV1, &roleAssign)

			// Get the deleted role assignment
			roleAssignTRef := secapi.TenantReference{
				Tenant: secapi.TenantID(suite.Tenant),
				Name:   roleAssign.Metadata.Name,
			}
			configurator.WatchRoleAssignmentUntilDeletedV1Step("Watch the role assignment deletion", sCtx, suite.Client.AuthorizationV1, roleAssignTRef)
		}
	})
}

// Workspace

func BulkCreateWorkspacesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, workspaces []schema.Workspace) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, workspace := range workspaces {
			expectWorkspaceMeta := workspace.Metadata
			expectWorkspaceSpec := &workspace.Spec

			// Create a workspace
			configurator.CreateOrUpdateWorkspaceV1Step("Create a workspace", sCtx, suite.Client.WorkspaceV1, &workspace,
				StepResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
					Metadata:       expectWorkspaceMeta,
					Spec:           expectWorkspaceSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteWorkspacesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, workspaces []schema.Workspace) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, workspace := range workspaces {
			// Delete the workspace
			configurator.DeleteWorkspaceV1Step("Delete the workspace", sCtx, suite.Client.WorkspaceV1, &workspace)

			// Get the deleted workspace
			tRef := secapi.TenantReference{
				Tenant: secapi.TenantID(workspace.Metadata.Tenant),
				Name:   workspace.Metadata.Name,
			}
			configurator.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", sCtx, suite.Client.WorkspaceV1, tRef)
		}
	})
}

// Compute

func BulkCreateInstancesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, instances []schema.Instance) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, instance := range instances {
			expectInstanceMeta := instance.Metadata
			expectInstanceSpec := &instance.Spec

			// Create an instance
			configurator.CreateOrUpdateInstanceV1Step("Create an instance", sCtx, suite.Client.ComputeV1, &instance,
				StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
					Metadata:       expectInstanceMeta,
					Spec:           expectInstanceSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteInstancesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, instances []schema.Instance) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, instance := range instances {
			// Delete the instance
			configurator.DeleteInstanceV1Step("Delete the instance", sCtx, suite.Client.ComputeV1, &instance)

			// Get the deleted instance
			wRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(instance.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(instance.Metadata.Workspace),
				Name:      instance.Metadata.Name,
			}
			configurator.WatchInstanceUntilDeletedV1Step("Watch the instance deletion", sCtx, suite.Client.ComputeV1, wRef)
		}
	})
}

// Storage

func BulkCreateBlockStoragesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, blockStorages []schema.BlockStorage) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, blockStorage := range blockStorages {
			expectedBlockMeta := blockStorage.Metadata
			expectedBlockSpec := &blockStorage.Spec

			// Create a block storage
			configurator.CreateOrUpdateBlockStorageV1Step("Create a block storage", sCtx, suite.Client.StorageV1, &blockStorage,
				StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
					Metadata:       expectedBlockMeta,
					Spec:           expectedBlockSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteBlockStoragesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, blockStorages []schema.BlockStorage) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, blockStorage := range blockStorages {
			// Delete the block storage
			configurator.DeleteBlockStorageV1Step("Delete block storage", sCtx, suite.Client.StorageV1, &blockStorage)

			// Get the deleted block storage
			blockWRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(blockStorage.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(blockStorage.Metadata.Workspace),
				Name:      blockStorage.Metadata.Name,
			}
			configurator.WatchBlockStorageUntilDeletedV1Step("Watch the block storage deletion", sCtx, suite.Client.StorageV1, blockWRef)
		}
	})
}

func BulkCreateImagesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, images []schema.Image) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, image := range images {
			expectedImageMeta := image.Metadata
			expectedImageSpec := &image.Spec

			// Create an image
			configurator.CreateOrUpdateImageV1Step("Create an image", sCtx, suite.Client.StorageV1, &image,
				StepResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
					Metadata:       expectedImageMeta,
					Spec:           expectedImageSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteImagesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, images []schema.Image) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, image := range images {
			// Delete the image
			configurator.DeleteImageV1Step("Delete the image", sCtx, suite.Client.StorageV1, &image)

			// Get the deleted image
			imageTRef := secapi.TenantReference{
				Tenant: secapi.TenantID(image.Metadata.Tenant),
				Name:   image.Metadata.Name,
			}
			configurator.WatchImageUntilDeletedV1Step("Watch the image deletion", sCtx, suite.Client.StorageV1, imageTRef)
		}
	})
}

// Network

func BulkCreateNetworksStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, networks []schema.Network) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, network := range networks {
			expectNetworkMeta := network.Metadata
			expectNetworkSpec := &network.Spec
			configurator.CreateOrUpdateNetworkV1Step("Create a network", sCtx, suite.Client.NetworkV1, &network,
				StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
					Metadata:       expectNetworkMeta,
					Spec:           expectNetworkSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteNetworksStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, networks []schema.Network) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, network := range networks {
			// Delete the network
			configurator.DeleteNetworkV1Step("Delete the network", sCtx, suite.Client.NetworkV1, &network)

			// Get the deleted network
			wRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(network.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
				Name:      network.Metadata.Name,
			}
			configurator.WatchNetworkUntilDeletedV1Step("Watch the network deletion", sCtx, suite.Client.NetworkV1, wRef)
		}
	})
}

func BulkCreateInternetGatewaysStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, internetGateways []schema.InternetGateway) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, internetGateway := range internetGateways {
			expectGatewayMeta := internetGateway.Metadata
			expectGatewaySpec := &internetGateway.Spec

			// Create an internet gateway
			configurator.CreateOrUpdateInternetGatewayV1Step("Create an internet gateway", sCtx, suite.Client.NetworkV1, &internetGateway,
				StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
					Metadata:       expectGatewayMeta,
					Spec:           expectGatewaySpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteInternetGatewaysStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, internetGateways []schema.InternetGateway) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, internetGateway := range internetGateways {
			// Delete the internet gateway
			configurator.DeleteInternetGatewayV1Step("Delete the internet gateway", sCtx, suite.Client.NetworkV1, &internetGateway)

			// Get the deleted internet gateway
			wRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(internetGateway.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(internetGateway.Metadata.Workspace),
				Name:      internetGateway.Metadata.Name,
			}
			configurator.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", sCtx, suite.Client.NetworkV1, wRef)
		}
	})
}

func BulkCreateRouteTablesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, routeTables []schema.RouteTable) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, routeTable := range routeTables {
			expectRouteMeta := routeTable.Metadata
			expectRouteSpec := &routeTable.Spec

			// Create a route table
			configurator.CreateOrUpdateRouteTableV1Step("Create a route table", sCtx, suite.Client.NetworkV1, &routeTable,
				StepResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
					Metadata:       expectRouteMeta,
					Spec:           expectRouteSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteRouteTablesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, routeTables []schema.RouteTable) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, routeTable := range routeTables {
			// Delete the route table
			configurator.DeleteRouteTableV1Step("Delete the route table", sCtx, suite.Client.NetworkV1, &routeTable)

			// Get the deleted route table
			nRef := secapi.NetworkReference{
				Tenant:    secapi.TenantID(routeTable.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(routeTable.Metadata.Workspace),
				Network:   secapi.NetworkID(routeTable.Metadata.Network),
				Name:      routeTable.Metadata.Name,
			}
			configurator.WatchRouteTableUntilDeletedV1Step("Watch the route table deletion", sCtx, suite.Client.NetworkV1, nRef)
		}
	})
}

func BulkCreateSubnetsStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, subnets []schema.Subnet) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, subnet := range subnets {
			expectSubnetMeta := subnet.Metadata
			expectSubnetSpec := &subnet.Spec

			// Create a subnet
			configurator.CreateOrUpdateSubnetV1Step("Create a subnet", sCtx, suite.Client.NetworkV1, &subnet,
				StepResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
					Metadata:       expectSubnetMeta,
					Spec:           expectSubnetSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteSubnetsStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, subnets []schema.Subnet) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, subnet := range subnets {
			// Delete the subnet
			configurator.DeleteSubnetV1Step("Delete the subnet", sCtx, suite.Client.NetworkV1, &subnet)

			// Get the deleted subnet
			nRef := secapi.NetworkReference{
				Tenant:    secapi.TenantID(subnet.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(subnet.Metadata.Workspace),
				Network:   secapi.NetworkID(subnet.Metadata.Network),
				Name:      subnet.Metadata.Name,
			}
			configurator.WatchSubnetUntilDeletedV1Step("Watch the subnet deletion", sCtx, suite.Client.NetworkV1, nRef)
		}
	})
}

func BulkCreatePublicIpsStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, publicIps []schema.PublicIp) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, publicIp := range publicIps {
			expectPublicIpMeta := publicIp.Metadata
			expectPublicIpSpec := &publicIp.Spec

			// Create a public ip
			configurator.CreateOrUpdatePublicIpV1Step("Create a public ip", sCtx, suite.Client.NetworkV1, &publicIp,
				StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
					Metadata:       expectPublicIpMeta,
					Spec:           expectPublicIpSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeletePublicIpsStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, publicIps []schema.PublicIp) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, publicIp := range publicIps {
			// Delete the public ip
			configurator.DeletePublicIpV1Step("Delete the public ip", sCtx, suite.Client.NetworkV1, &publicIp)

			// Get the deleted public ip
			wRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(publicIp.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(publicIp.Metadata.Workspace),
				Name:      publicIp.Metadata.Name,
			}
			configurator.WatchPublicIpUntilDeletedV1Step("Watch the public ip deletion", sCtx, suite.Client.NetworkV1, wRef)
		}
	})
}

func BulkCreateNicsStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, nics []schema.Nic) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, nic := range nics {
			expectNicMeta := nic.Metadata
			expectNicSpec := &nic.Spec

			// Create a nic
			configurator.CreateOrUpdateNicV1Step("Create a nic", sCtx, suite.Client.NetworkV1, &nic,
				StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
					Metadata:       expectNicMeta,
					Spec:           expectNicSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteNicsStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, nics []schema.Nic) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, nic := range nics {
			// Delete the nic
			configurator.DeleteNicV1Step("Delete the nic", sCtx, suite.Client.NetworkV1, &nic)

			// Get the deleted nic
			wRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(nic.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(nic.Metadata.Workspace),
				Name:      nic.Metadata.Name,
			}
			configurator.WatchNicUntilDeletedV1Step("Watch the nic deletion", sCtx, suite.Client.NetworkV1, wRef)
		}
	})
}

func BulkCreateSecurityGroupRulesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, securityGroupRules []schema.SecurityGroupRule) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, securityGroupRule := range securityGroupRules {
			expectRuleMeta := securityGroupRule.Metadata
			expectRuleSpec := &securityGroupRule.Spec

			// Create a security group rule
			configurator.CreateOrUpdateSecurityGroupRuleV1Step("Create a security group rule", sCtx, suite.Client.NetworkV1, &securityGroupRule,
				StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
					Metadata:       expectRuleMeta,
					Spec:           expectRuleSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteSecurityGroupRulesStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, securityGroupRules []schema.SecurityGroupRule) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, securityGroupRule := range securityGroupRules {
			// Delete the security group rule
			configurator.DeleteSecurityGroupRuleV1Step("Delete the security group rule", sCtx, suite.Client.NetworkV1, &securityGroupRule)

			// Get the deleted security group rule
			wRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(securityGroupRule.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(securityGroupRule.Metadata.Workspace),
				Name:      securityGroupRule.Metadata.Name,
			}
			configurator.WatchSecurityGroupRuleUntilDeletedV1Step("Watch the security group rule deletion", sCtx, suite.Client.NetworkV1, wRef)
		}
	})
}

func BulkCreateSecurityGroupsStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, securityGroups []schema.SecurityGroup) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, securityGroup := range securityGroups {
			expectGroupMeta := securityGroup.Metadata
			expectGroupSpec := &securityGroup.Spec

			// Create a security group
			configurator.CreateOrUpdateSecurityGroupV1Step("Create a security group", sCtx, suite.Client.NetworkV1, &securityGroup,
				StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
					Metadata:       expectGroupMeta,
					Spec:           expectGroupSpec,
					ResourceStates: suites.CreatedResourceExpectedStates,
				},
			)
		}
	})
}

func BulkDeleteSecurityGroupsStepsV1(configurator *StepsConfigurator, suite suites.RegionalTestSuite, stepName string, securityGroups []schema.SecurityGroup) {
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		for _, securityGroup := range securityGroups {
			// Delete the security group
			configurator.DeleteSecurityGroupV1Step("Delete the security group", sCtx, suite.Client.NetworkV1, &securityGroup)

			// Get the deleted security group
			wRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(securityGroup.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(securityGroup.Metadata.Workspace),
				Name:      securityGroup.Metadata.Name,
			}
			configurator.WatchSecurityGroupUntilDeletedV1Step("Watch the security group deletion", sCtx, suite.Client.NetworkV1, wRef)
		}
	})
}
