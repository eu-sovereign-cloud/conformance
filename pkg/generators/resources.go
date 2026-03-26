package generators

import (
	"fmt"
)

func GenerateSkuResource(name string) string {
	return fmt.Sprintf(skuResource, name)
}

func GenerateSkuListResource() string {
	return skuListResource
}

func GenerateRoleResource(name string) string {
	return fmt.Sprintf(roleResource, name)
}

func GenerateRoleListResource() string {
	return roleListResource
}

func GenerateRoleAssignmentResource(name string) string {
	return fmt.Sprintf(roleAssignmentResource, name)
}

func GenerateRoleAssignmentListResource() string {
	return roleAssignmentListResource
}

func GenerateRegionResource(name string) string {
	return fmt.Sprintf(regionResource, name)
}

func GenerateRegionListResource() string {
	return regionListResource
}

func GenerateWorkspaceResource(name string) string {
	return fmt.Sprintf(workspaceResource, name)
}

func GenerateWorkspaceListResource() string {
	return workspaceListResource
}

func GenerateBlockStorageResource(name string) string {
	return fmt.Sprintf(blockStorageResource, name)
}

func GenerateBlockStorageListResource() string {
	return blockStorageListResource
}

func GenerateImageResource(name string) string {
	return fmt.Sprintf(imageResource, name)
}

func GenerateImageListResource() string {
	return imageListResource
}

func GenerateInstanceResource(name string) string {
	return fmt.Sprintf(instanceResource, name)
}

func GenerateInstanceListResource() string {
	return instanceListResource
}

func GenerateNetworkResource(name string) string {
	return fmt.Sprintf(networkResource, name)
}

func GenerateNetworkListResource() string {
	return networkListResource
}

func GenerateInternetGatewayResource(name string) string {
	return fmt.Sprintf(internetGatewayResource, name)
}

func GenerateInternetGatewayListResource() string {
	return internetGatewayListResource
}

func GenerateNicResource(name string) string {
	return fmt.Sprintf(nicResource, name)
}

func GenerateNicListResource() string {
	return nicListResource
}

func GeneratePublicIpResource(name string) string {
	return fmt.Sprintf(publicIpResource, name)
}

func GeneratePublicIpListResource() string {
	return publicIpListResource
}

func GenerateRouteTableResource(network, name string) string {
	return fmt.Sprintf(routeTableResource, network, name)
}

func GenerateRouteTableListResource(network string) string {
	return fmt.Sprintf(routeTableListResource, network)
}

func GenerateSubnetResource(network, name string) string {
	return fmt.Sprintf(subnetResource, network, name)
}

func GenerateSubnetListResource(network string) string {
	return fmt.Sprintf(subnetListResource, network)
}

func GenerateSecurityGroupRuleResource(name string) string {
	return fmt.Sprintf(securityGroupRuleResource, name)
}

func GenerateSecurityGroupResource(name string) string {
	return fmt.Sprintf(securityGroupResource, name)
}

func GenerateSecurityGroupListResource() string {
	return securityGroupListResource
}
