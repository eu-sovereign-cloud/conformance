package secalib

import (
	"fmt"
	"math"
	"math/rand"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

// Names

func GenerateRoleName() string {
	return fmt.Sprintf("role-%d", rand.Intn(math.MaxInt32))
}

func GenerateRoleAssignmentName() string {
	return fmt.Sprintf("role-assignment-%d", rand.Intn(math.MaxInt32))
}

func GenerateRegionName() string {
	return fmt.Sprintf("region-%d", rand.Intn(math.MaxInt32))
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

func GenerateNetworkName() string {
	return fmt.Sprintf("network-%d", rand.Intn(math.MaxInt32))
}

func GenerateInternetGatewayName() string {
	return fmt.Sprintf("internet-gateway-%d", rand.Intn(math.MaxInt32))
}

func GenerateRouteTableName() string {
	return fmt.Sprintf("route-table-%d", rand.Intn(math.MaxInt32))
}

func GenerateSubnetName() string {
	return fmt.Sprintf("subnet-%d", rand.Intn(math.MaxInt32))
}

func GeneratePublicIpName() string {
	return fmt.Sprintf("public-ip-%d", rand.Intn(math.MaxInt32))
}

func GenerateNicName() string {
	return fmt.Sprintf("nic-%d", rand.Intn(math.MaxInt32))
}

func GenerateSecurityGroupName() string {
	return fmt.Sprintf("security-group-%d", rand.Intn(math.MaxInt32))
}

func GenerateRegionProviderUrl(provider string) string {
	return fmt.Sprintf("{{request.scheme}}://{{request.host}}:{{request.port}}%s%s", urlProvidersPrefix, provider)
}

// Resources

func GenerateSkuResource(tenant string, sku string) string {
	return fmt.Sprintf(skuResource, tenant, sku)
}

func GenerateSkuListResource(tenant string) string {
	return fmt.Sprintf(SkuListResource, tenant)
}

func GenerateRoleResource(tenant string, role string) string {
	return fmt.Sprintf(roleResource, tenant, role)
}

func GenerateRolesResource(tenant string) string {
	return fmt.Sprintf(RolesResource, tenant)
}

func GenerateRoleAssignmentResource(tenant string, roleAssignment string) string {
	return fmt.Sprintf(roleAssignmentResource, tenant, roleAssignment)
}

func GenerateRoleAssignmentsResource(tenant string) string {
	return fmt.Sprintf(RoleAssignmentsResource, tenant)
}

func GenerateRegionResource(region string) string {
	return fmt.Sprintf(regionResource, region)
}

func GenerateWorkspaceResource(tenant string, workspace string) string {
	return fmt.Sprintf(workspaceResource, tenant, workspace)
}

func GenerateWorkspaceListResource(tenant string) string {
	return fmt.Sprintf(WorkspaceListResource, tenant)
}

func GenerateBlockStorageResource(tenant string, workspace string, blockStorage string) string {
	return fmt.Sprintf(blockStorageResource, tenant, workspace, blockStorage)
}

func GenerateBlockStorageListResource(tenant string, workspace string) string {
	return fmt.Sprintf(BlockStorageListResource, tenant, workspace)
}

func GenerateImageResource(tenant string, image string) string {
	return fmt.Sprintf(imageResource, tenant, image)
}

func GenerateImageListResource(tenant string) string {
	return fmt.Sprintf(ImageListResource, tenant)
}

func GenerateInstanceResource(tenant string, workspace string, instance string) string {
	return fmt.Sprintf(instanceResource, tenant, workspace, instance)
}

func GenerateInstanceListResource(tenant string, workspace string) string {
	return fmt.Sprintf(InstanceListResource, tenant, workspace)
}

func GenerateNetworkResource(tenant string, workspace string, network string) string {
	return fmt.Sprintf(networkResource, tenant, workspace, network)
}

func GenerateNetworkListResource(tenant string, workspace string) string {
	return fmt.Sprintf(NetworkListResource, tenant, workspace)
}

func GenerateInternetGatewayResource(tenant string, workspace string, internetGateway string) string {
	return fmt.Sprintf(internetGatewayResource, tenant, workspace, internetGateway)
}

func GenerateInternetGatewayListResource(tenant string, workspace string) string {
	return fmt.Sprintf(InternetGatewayListResource, tenant, workspace)
}

func GenerateNicResource(tenant string, workspace string, nic string) string {
	return fmt.Sprintf(nicResource, tenant, workspace, nic)
}

func GenerateNicListResource(tenant string, workspace string) string {
	return fmt.Sprintf(NicListResource, tenant, workspace)
}

func GeneratePublicIpResource(tenant string, workspace string, publicIp string) string {
	return fmt.Sprintf(publicIpResource, tenant, workspace, publicIp)
}

func GeneratePublicIpListResource(tenant string, workspace string) string {
	return fmt.Sprintf(PublicIpListResource, tenant, workspace)
}

func GenerateRouteTableResource(tenant string, workspace string, network string, routeTable string) string {
	return fmt.Sprintf(routeTableResource, tenant, workspace, network, routeTable)
}

func GenerateRouteTableListResource(tenant string, workspace string, network string) string {
	return fmt.Sprintf(RouteTableListResource, tenant, workspace, network)
}

func GenerateSubnetResource(tenant string, workspace string, network string, subnet string) string {
	return fmt.Sprintf(subnetResource, tenant, workspace, network, subnet)
}

func GenerateSubnetListResource(tenant string, workspace string, network string) string {
	return fmt.Sprintf(SubnetListResource, tenant, workspace, network)
}

func GenerateSecurityGroupResource(tenant string, workspace string, securityGroup string) string {
	return fmt.Sprintf(securityGroupResource, tenant, workspace, securityGroup)
}

func GenerateSecurityGroupListResource(tenant string, workspace string) string {
	return fmt.Sprintf(SecurityGroupListResource, tenant, workspace)
}

// References

func GenerateSkuRef(name string) string {
	return fmt.Sprintf(skuRef, name)
}

func GenerateInstanceRef(instanceName string) string {
	return fmt.Sprintf(instanceRef, instanceName)
}

func GenerateBlockStorageRef(blockStorageName string) string {
	return fmt.Sprintf(blockStorageRef, blockStorageName)
}

func GenerateInternetGatewayRef(internetGatewayName string) string {
	return fmt.Sprintf(internetGatewayRef, internetGatewayName)
}

func GenerateNetworkRef(networkName string) string {
	return fmt.Sprintf(networkRef, networkName)
}

func GenerateRouteTableRef(routeTableName string) string {
	return fmt.Sprintf(routeTableRef, routeTableName)
}

func GenerateSubnetRef(subnetName string) string {
	return fmt.Sprintf(subnetRef, subnetName)
}

func GeneratePublicIpRef(publicIpName string) string {
	return fmt.Sprintf(publicIpRef, publicIpName)
}

// Random

func GenerateBlockStorageSize() int {
	return rand.Intn(maxBlockStorageSize)
}

// Network

func GenerateSubnetCidr(networkCidr string, size int, netNum int) (string, error) {
	_, network, err := net.ParseCIDR(networkCidr)
	if err != nil {
		return "", err
	}

	subnet, err := cidr.Subnet(network, size, netNum)
	if err != nil {
		return "", err
	}

	return subnet.String(), nil
}

func GenerateNicAddress(subnetCidr string, hostNum int) (string, error) {
	_, network, err := net.ParseCIDR(subnetCidr)
	if err != nil {
		return "", err
	}

	ip, err := cidr.Host(network, hostNum)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}

func GeneratePublicIp(publicIpRange string, hostNum int) (string, error) {
	_, network, err := net.ParseCIDR(publicIpRange)
	if err != nil {
		return "", err
	}

	ip, err := cidr.Host(network, hostNum)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}
