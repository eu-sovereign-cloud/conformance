package secalib

import (
	"fmt"
	"math"
	"math/rand"
)

// Names
func GenerateRoleName() string {
	return fmt.Sprintf("role-%d", rand.Intn(math.MaxInt32))
}

func GenerateRoleAssignmentName() string {
	return fmt.Sprintf("role-assignment-%d", rand.Intn(math.MaxInt32))
}

func GenerateWorkspaceName() string {
	return fmt.Sprintf("workspace-%d", rand.Intn(math.MaxInt32))
}

func GenerateBlockStorageName() string {
	return fmt.Sprintf("disk-%d", rand.Intn(math.MaxInt32))
}

func GenerateImageName() string {
	return fmt.Sprintf("image-%d", rand.Intn(math.MaxInt32))
}

func GenerateInstanceName() string {
	return fmt.Sprintf("instance-%d", rand.Intn(math.MaxInt32))
}

// Resources
func GenerateSkuResource(tenant string, sku string) string {
	return fmt.Sprintf(SkuResource, tenant, sku)
}

func GenerateRoleResource(tenant string, role string) string {
	return fmt.Sprintf(RoleResource, tenant, role)
}

func GenerateRoleAssignmentResource(tenant string, roleAssignment string) string {
	return fmt.Sprintf(RoleAssignmentResource, tenant, roleAssignment)
}

func GenerateWorkspaceResource(tenant string, workspace string) string {
	return fmt.Sprintf(WorkspaceResource, tenant, workspace)
}

func GenerateBlockStorageResource(tenant string, workspace string, blockStorage string) string {
	return fmt.Sprintf(BlockStorageResource, tenant, workspace, blockStorage)
}

func GenerateImageResource(tenant string, image string) string {
	return fmt.Sprintf(ImageResource, tenant, image)
}

func GenerateInstanceResource(tenant string, workspace string, instance string) string {
	return fmt.Sprintf(InstanceResource, tenant, workspace, instance)
}
func GenerateNetworkResource(tenant string, workspace string, network string) string {
	return fmt.Sprintf(NetworkResource, tenant, workspace, network)
}
func GenerateInternetGatewayResource(tenant string, workspace string, internetGateway string) string {
	return fmt.Sprintf(InternetGatewayResource, tenant, workspace, internetGateway)
}
func GenerateNicResource(tenant string, workspace string, nic string) string {
	return fmt.Sprintf(NicResource, tenant, workspace, nic)
}
func GeneratePublicIPResource(tenant string, workspace string, publicIP string) string {
	return fmt.Sprintf(PublicIPResource, tenant, workspace, publicIP)
}
func GenerateRouteTableResource(tenant string, workspace string, routeTable string) string {
	return fmt.Sprintf(RouteTableResource, tenant, workspace, routeTable)
}
func GenerateSubnetResource(tenant string, workspace string, subnet string) string {
	return fmt.Sprintf(SubnetResource, tenant, workspace, subnet)
}
func GenerateSecurityGroupResource(tenant string, workspace string, securityGroup string) string {
	return fmt.Sprintf(SecurityGroupResource, tenant, workspace, securityGroup)
}

// References
func GenerateSkuRef(name string) string {
	return fmt.Sprintf(SkuRef, name)
}

func GenerateBlockStorageRef(blockStorageName string) string {
	return fmt.Sprintf(BlockStorageRef, blockStorageName)
}

func GenerateInternetGatewayRef(internetGatewayName string) string {
	return fmt.Sprintf(InternetGatewayRef, internetGatewayName)
}

func GenerateNetworkRef(networkName string) string {
	return fmt.Sprintf(NetworkRef, networkName)
}

func GenerateRouteTableRef(routeTableName string) string {
	return fmt.Sprintf(RouteTableRef, routeTableName)
}

func GenerateSubnetRef(subnetName string) string {
	return fmt.Sprintf(SubnetRef, subnetName)
}

func GeneratePublicIPRef(publicIPName string) string {
	return fmt.Sprintf(PublicIPRef, publicIPName)
}

// URLs
func GenerateRoleURL(tenant string, role string) string {
	return fmt.Sprintf(RoleURLV1, tenant, role)
}

func GenerateRoleAssignmentURL(tenant string, roleAssignment string) string {
	return fmt.Sprintf(RoleAssignmentURLV1, tenant, roleAssignment)
}

func GenerateWorkspaceURL(tenant string, workspace string) string {
	return fmt.Sprintf(WorkspaceURLV1, tenant, workspace)
}

func GenerateStorageSkuURL(tenant string, sku string) string {
	return fmt.Sprintf(StorageSkuURLV1, tenant, sku)
}

func GenerateBlockStorageURL(tenant string, workspace string, blockStorage string) string {
	return fmt.Sprintf(BlockStorageURLV1, tenant, workspace, blockStorage)
}

func GenerateImageURL(tenant string, image string) string {
	return fmt.Sprintf(ImageURLV1, tenant, image)
}

func GenerateInstanceSkuURL(tenant string, sku string) string {
	return fmt.Sprintf(InstanceSkuURLV1, tenant, sku)
}

func GenerateInstanceURL(tenant string, workspace string, instance string) string {
	return fmt.Sprintf(InstanceURLV1, tenant, workspace, instance)
}

func GenerateNetworkURL(tenant string, workspace string, network string) string {
	return fmt.Sprintf(NetworkURLV1, tenant, workspace, network)
}

func GenerateNetworkSkuURL(tenant string, sku string) string {
	return fmt.Sprintf(NetworkSkuURLV1, tenant, sku)
}

func GenerateInternetGatewayURL(tenant string, workspace string, internetGateway string) string {
	return fmt.Sprintf(InternetGatewayURLV1, tenant, workspace, internetGateway)
}

func GenerateNicURL(tenant string, workspace string, nic string) string {
	return fmt.Sprintf(NicURLV1, tenant, workspace, nic)
}

func GeneratePublicIPURL(tenant string, workspace string, publicIP string) string {
	return fmt.Sprintf(PublicIPURLV1, tenant, workspace, publicIP)
}

func GenerateRouteTableURL(tenant string, workspace string, routeTable string) string {
	return fmt.Sprintf(RouteTableURLV1, tenant, workspace, routeTable)
}

func GenerateSubnetURL(tenant string, workspace string, subnet string) string {
	return fmt.Sprintf(SubnetURLV1, tenant, workspace, subnet)
}

func GenerateSecurityGroupURL(tenant string, workspace string, securityGroup string) string {
	return fmt.Sprintf(SecurityGroupURLV1, tenant, workspace, securityGroup)
}

// Random

func GenerateBlockStorageSize() int {
	return rand.Intn(maxBlockStorageSize)
}
