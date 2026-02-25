package generators

import (
	"fmt"
)

func GenerateSkuResource(tenant, sku string) string {
	return fmt.Sprintf(skuResource, tenant, sku)
}

func GenerateSkuListResource(tenant string) string {
	return fmt.Sprintf(skuListResource, tenant)
}

func GenerateRoleResource(tenant, role string) string {
	return fmt.Sprintf(roleResource, tenant, role)
}

func GenerateRoleListResource(tenant string) string {
	return fmt.Sprintf(roleListResource, tenant)
}

func GenerateRoleAssignmentResource(tenant, roleAssignment string) string {
	return fmt.Sprintf(roleAssignmentResource, tenant, roleAssignment)
}

func GenerateRoleAssignmentListResource(tenant string) string {
	return fmt.Sprintf(roleAssignmentListResource, tenant)
}

func GenerateRegionResource(region string) string {
	return fmt.Sprintf(regionResource, region)
}

func GenerateRegionListResource() string {
	return regionListResource
}

func GenerateWorkspaceResource(tenant, workspace string) string {
	return fmt.Sprintf(workspaceResource, tenant, workspace)
}

func GenerateWorkspaceListResource(tenant string) string {
	return fmt.Sprintf(workspaceListResource, tenant)
}

func GenerateBlockStorageResource(tenant, workspace, blockStorage string) string {
	return fmt.Sprintf(blockStorageResource, tenant, workspace, blockStorage)
}

func GenerateBlockStorageListResource(tenant, workspace string) string {
	return fmt.Sprintf(blockStorageListResource, tenant, workspace)
}

func GenerateImageResource(tenant, image string) string {
	return fmt.Sprintf(imageResource, tenant, image)
}

func GenerateImageListResource(tenant string) string {
	return fmt.Sprintf(imageListResource, tenant)
}

func GenerateInstanceResource(tenant, workspace, instance string) string {
	return fmt.Sprintf(instanceResource, tenant, workspace, instance)
}

func GenerateInstanceListResource(tenant, workspace string) string {
	return fmt.Sprintf(instanceListResource, tenant, workspace)
}

func GenerateNetworkResource(tenant, workspace, network string) string {
	return fmt.Sprintf(networkResource, tenant, workspace, network)
}

func GenerateNetworkListResource(tenant, workspace string) string {
	return fmt.Sprintf(networkListResource, tenant, workspace)
}

func GenerateInternetGatewayResource(tenant, workspace, internetGateway string) string {
	return fmt.Sprintf(internetGatewayResource, tenant, workspace, internetGateway)
}

func GenerateInternetGatewayListResource(tenant, workspace string) string {
	return fmt.Sprintf(internetGatewayListResource, tenant, workspace)
}

func GenerateNicResource(tenant, workspace, nic string) string {
	return fmt.Sprintf(nicResource, tenant, workspace, nic)
}

func GenerateNicListResource(tenant, workspace string) string {
	return fmt.Sprintf(nicListResource, tenant, workspace)
}

func GeneratePublicIpResource(tenant, workspace, publicIp string) string {
	return fmt.Sprintf(publicIpResource, tenant, workspace, publicIp)
}

func GeneratePublicIpListResource(tenant, workspace string) string {
	return fmt.Sprintf(publicIpListResource, tenant, workspace)
}

func GenerateRouteTableResource(tenant, workspace, network, routeTable string) string {
	return fmt.Sprintf(routeTableResource, tenant, workspace, network, routeTable)
}

func GenerateRouteTableListResource(tenant, workspace, network string) string {
	return fmt.Sprintf(routeTableListResource, tenant, workspace, network)
}

func GenerateSubnetResource(tenant, workspace, network, subnet string) string {
	return fmt.Sprintf(subnetResource, tenant, workspace, network, subnet)
}

func GenerateSubnetListResource(tenant, workspace, network string) string {
	return fmt.Sprintf(subnetListResource, tenant, workspace, network)
}

func GenerateSecurityGroupRuleResource(tenant, workspace, securityGroupRule string) string {
	return fmt.Sprintf(securityGroupRuleResource, tenant, workspace, securityGroupRule)
}

func GenerateSecurityGroupResource(tenant, workspace, securityGroup string) string {
	return fmt.Sprintf(securityGroupResource, tenant, workspace, securityGroup)
}

func GenerateSecurityGroupListResource(tenant, workspace string) string {
	return fmt.Sprintf(securityGroupListResource, tenant, workspace)
}
