package generators

import (
	"fmt"
)

func GenerateSkuResource(provider, tenant, sku string) string {
	return fmt.Sprintf(skuResource, provider, tenant, sku)
}

func GenerateSkuListResource(provider, tenant string) string {
	return fmt.Sprintf(skuListResource, provider, tenant)
}

func GenerateRoleResource(provider, tenant, role string) string {
	return fmt.Sprintf(roleResource, provider, tenant, role)
}

func GenerateRoleListResource(provider, tenant string) string {
	return fmt.Sprintf(roleListResource, provider, tenant)
}

func GenerateRoleAssignmentResource(provider, tenant, roleAssignment string) string {
	return fmt.Sprintf(roleAssignmentResource, provider, tenant, roleAssignment)
}

func GenerateRoleAssignmentListResource(provider, tenant string) string {
	return fmt.Sprintf(roleAssignmentListResource, provider, tenant)
}

func GenerateRegionResource(provider, region string) string {
	return fmt.Sprintf(regionResource, provider, region)
}

func GenerateRegionListResource(provider string) string {
	return fmt.Sprintf(regionListResource, provider)
}

func GenerateWorkspaceResource(provider, tenant, workspace string) string {
	return fmt.Sprintf(workspaceResource, provider, tenant, workspace)
}

func GenerateWorkspaceListResource(provider, tenant string) string {
	return fmt.Sprintf(workspaceListResource, provider, tenant)
}

func GenerateBlockStorageResource(provider, tenant, workspace, blockStorage string) string {
	return fmt.Sprintf(blockStorageResource, provider, tenant, workspace, blockStorage)
}

func GenerateBlockStorageListResource(provider, tenant, workspace string) string {
	return fmt.Sprintf(blockStorageListResource, provider, tenant, workspace)
}

func GenerateImageResource(provider, tenant, image string) string {
	return fmt.Sprintf(imageResource, provider, tenant, image)
}

func GenerateImageListResource(provider, tenant string) string {
	return fmt.Sprintf(imageListResource, provider, tenant)
}

func GenerateInstanceResource(provider, tenant, workspace, instance string) string {
	return fmt.Sprintf(instanceResource, provider, tenant, workspace, instance)
}

func GenerateInstanceListResource(provider, tenant, workspace string) string {
	return fmt.Sprintf(instanceListResource, provider, tenant, workspace)
}

func GenerateNetworkResource(provider, tenant, workspace, network string) string {
	return fmt.Sprintf(networkResource, provider, tenant, workspace, network)
}

func GenerateNetworkListResource(provider, tenant, workspace string) string {
	return fmt.Sprintf(networkListResource, provider, tenant, workspace)
}

func GenerateInternetGatewayResource(provider, tenant, workspace, internetGateway string) string {
	return fmt.Sprintf(internetGatewayResource, provider, tenant, workspace, internetGateway)
}

func GenerateInternetGatewayListResource(provider, tenant, workspace string) string {
	return fmt.Sprintf(internetGatewayListResource, provider, tenant, workspace)
}

func GenerateNicResource(provider, tenant, workspace, nic string) string {
	return fmt.Sprintf(nicResource, provider, tenant, workspace, nic)
}

func GenerateNicListResource(provider, tenant, workspace string) string {
	return fmt.Sprintf(nicListResource, provider, tenant, workspace)
}

func GeneratePublicIpResource(provider, tenant, workspace, publicIp string) string {
	return fmt.Sprintf(publicIpResource, provider, tenant, workspace, publicIp)
}

func GeneratePublicIpListResource(provider, tenant, workspace string) string {
	return fmt.Sprintf(publicIpListResource, provider, tenant, workspace)
}

func GenerateRouteTableResource(provider, tenant, workspace, network, routeTable string) string {
	return fmt.Sprintf(routeTableResource, provider, tenant, workspace, network, routeTable)
}

func GenerateRouteTableListResource(provider, tenant, workspace, network string) string {
	return fmt.Sprintf(routeTableListResource, provider, tenant, workspace, network)
}

func GenerateSubnetResource(provider, tenant, workspace, network, subnet string) string {
	return fmt.Sprintf(subnetResource, provider, tenant, workspace, network, subnet)
}

func GenerateSubnetListResource(provider, tenant, workspace, network string) string {
	return fmt.Sprintf(subnetListResource, provider, tenant, workspace, network)
}

func GenerateSecurityGroupRuleResource(provider, tenant, workspace, securityGroupRule string) string {
	return fmt.Sprintf(securityGroupRuleResource, provider, tenant, workspace, securityGroupRule)
}

func GenerateSecurityGroupResource(provider, tenant, workspace, securityGroup string) string {
	return fmt.Sprintf(securityGroupResource, provider, tenant, workspace, securityGroup)
}

func GenerateSecurityGroupListResource(provider, tenant, workspace string) string {
	return fmt.Sprintf(securityGroupListResource, provider, tenant, workspace)
}
