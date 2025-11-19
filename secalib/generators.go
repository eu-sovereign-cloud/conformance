package secalib

import (
	"fmt"
	"math"
	"math/rand"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
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
	return fmt.Sprintf("{{request.scheme}}://{{request.host}}:{{request.port}}%s%s", UrlProvidersPrefix, provider)
}

// Resources
func GenerateSkuResource(tenant string, sku string) string {
	return fmt.Sprintf(SkuResource, tenant, sku)
}

func GenerateSkuListResource(tenant string) string {
	return fmt.Sprintf(SkuListResource, tenant)
}

func GenerateRoleResource(tenant string, role string) string {
	return fmt.Sprintf(RoleResource, tenant, role)
}

func GenerateRolesResource(tenant string) string {
	return fmt.Sprintf(RolesResource, tenant)
}

func GenerateRoleAssignmentResource(tenant string, roleAssignment string) string {
	return fmt.Sprintf(RoleAssignmentResource, tenant, roleAssignment)
}

func GenerateRoleAssignmentsResource(tenant string) string {
	return fmt.Sprintf(RoleAssignmentsResource, tenant)
}

func GenerateRegionResource(region string) string {
	return fmt.Sprintf(RegionResource, region)
}

func GenerateWorkspaceResource(tenant string, workspace string) string {
	return fmt.Sprintf(WorkspaceResource, tenant, workspace)
}

func GenerateWorkspaceListResource(tenant string) string {
	return fmt.Sprintf(WorkspaceListResource, tenant)
}

func GenerateBlockStorageResource(tenant string, workspace string, blockStorage string) string {
	return fmt.Sprintf(BlockStorageResource, tenant, workspace, blockStorage)
}

func GenerateBlockStorageListResource(tenant string, workspace string) string {
	return fmt.Sprintf(BlockStorageListResource, tenant, workspace)
}

func GenerateImageResource(tenant string, image string) string {
	return fmt.Sprintf(ImageResource, tenant, image)
}

func GenerateImageListResource(tenant string) string {
	return fmt.Sprintf(ImageListResource, tenant)
}

func GenerateInstanceResource(tenant string, workspace string, instance string) string {
	return fmt.Sprintf(InstanceResource, tenant, workspace, instance)
}

func GenerateInstanceListResource(tenant string, workspace string) string {
	return fmt.Sprintf(InstanceListResource, tenant, workspace)
}

func GenerateNetworkResource(tenant string, workspace string, network string) string {
	return fmt.Sprintf(NetworkResource, tenant, workspace, network)
}

func GenerateNetworkListResource(tenant string, workspace string) string {
	return fmt.Sprintf(NetworkListResource, tenant, workspace)
}

func GenerateInternetGatewayResource(tenant string, workspace string, internetGateway string) string {
	return fmt.Sprintf(InternetGatewayResource, tenant, workspace, internetGateway)
}

func GenerateInternetGatewayListResource(tenant string, workspace string) string {
	return fmt.Sprintf(InternetGatewayListResource, tenant, workspace)
}

func GenerateNicResource(tenant string, workspace string, nic string) string {
	return fmt.Sprintf(NicResource, tenant, workspace, nic)
}

func GenerateNicListResource(tenant string, workspace string) string {
	return fmt.Sprintf(NicListResource, tenant, workspace)
}

func GeneratePublicIpResource(tenant string, workspace string, publicIp string) string {
	return fmt.Sprintf(PublicIpResource, tenant, workspace, publicIp)
}

func GeneratePublicIpListResource(tenant string, workspace string) string {
	return fmt.Sprintf(PublicIpListResource, tenant, workspace)
}

func GenerateRouteTableResource(tenant string, workspace string, network string, routeTable string) string {
	return fmt.Sprintf(RouteTableResource, tenant, workspace, network, routeTable)
}

func GenerateRouteTableListResource(tenant string, workspace string, network string) string {
	return fmt.Sprintf(RouteTableListResource, tenant, workspace, network)
}

func GenerateSubnetResource(tenant string, workspace string, network string, subnet string) string {
	return fmt.Sprintf(SubnetResource, tenant, workspace, network, subnet)
}

func GenerateSubnetListResource(tenant string, workspace string, network string) string {
	return fmt.Sprintf(SubnetListResource, tenant, workspace, network)
}

func GenerateSecurityGroupResource(tenant string, workspace string, securityGroup string) string {
	return fmt.Sprintf(SecurityGroupResource, tenant, workspace, securityGroup)
}

func GenerateSecurityGroupListResource(tenant string, workspace string) string {
	return fmt.Sprintf(SecurityGroupListResource, tenant, workspace)
}

// References
func GenerateSkuRef(name string) string {
	return fmt.Sprintf(SkuRef, name)
}

func GenerateInstanceRef(instanceName string) string {
	return fmt.Sprintf(InstanceRef, instanceName)
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

func GeneratePublicIpRef(publicIpName string) string {
	return fmt.Sprintf(PublicIpRef, publicIpName)
}

// URLs
func GenerateRoleURL(tenant string, role string) string {
	return fmt.Sprintf(RoleURLV1, tenant, role)
}

func GenerateRoleListURL(tenant string) string {
	return fmt.Sprintf(RoleListURLV1, tenant)
}

func GenerateRoleAssignmentURL(tenant string, roleAssignment string) string {
	return fmt.Sprintf(RoleAssignmentURLV1, tenant, roleAssignment)
}

func GenerateRoleAssignmentListURL(tenant string) string {
	return fmt.Sprintf(RoleAssignmentListURLV1, tenant)
}

func GenerateRegionURL(region string) string {
	return fmt.Sprintf(RegionURLV1, region)
}

func GenerateWorkspaceURL(tenant string, workspace string) string {
	return fmt.Sprintf(WorkspaceURLV1, tenant, workspace)
}

func GenerateWorkspaceListURL(tenant string) string {
	return fmt.Sprintf(WorkspaceListURLV1, tenant)
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

func GenerateImageListURL(tenant string) string {
	return fmt.Sprintf(ImageListURLV1, tenant)
}

func GenerateInstanceSkuURL(tenant string, sku string) string {
	return fmt.Sprintf(InstanceSkuURLV1, tenant, sku)
}

func GenerateInstanceURL(tenant string, workspace string, instance string) string {
	return fmt.Sprintf(InstanceURLV1, tenant, workspace, instance)
}

func GenerateInstanceListURL(tenant string, workspace string) string {
	return fmt.Sprintf(InstanceListURLV1, tenant, workspace)
}

func GenerateNetworkURL(tenant string, workspace string, network string) string {
	return fmt.Sprintf(NetworkURLV1, tenant, workspace, network)
}

func GenerateNetworkListURL(tenant string, workspace string) string {
	return fmt.Sprintf(NetworkListURLV1, tenant, workspace)
}

func GenerateNetworkSkuURL(tenant string, sku string) string {
	return fmt.Sprintf(NetworkSkuURLV1, tenant, sku)
}

func GenerateInternetGatewayURL(tenant string, workspace string, internetGateway string) string {
	return fmt.Sprintf(InternetGatewayURLV1, tenant, workspace, internetGateway)
}

func GenerateInternetGatewayListURL(tenant string, workspace string) string {
	return fmt.Sprintf(InternetGatewayListURLV1, tenant, workspace)
}

func GenerateNicURL(tenant string, workspace string, nic string) string {
	return fmt.Sprintf(NicURLV1, tenant, workspace, nic)
}

func GenerateNicListURL(tenant string, workspace string) string {
	return fmt.Sprintf(NicListURLV1, tenant, workspace)
}

func GeneratePublicIpURL(tenant string, workspace string, publicIp string) string {
	return fmt.Sprintf(PublicIpURLV1, tenant, workspace, publicIp)
}

func GeneratePublicIpListURL(tenant string, workspace string) string {
	return fmt.Sprintf(PublicIpListURLV1, tenant, workspace)
}

func GenerateRouteTableURL(tenant string, workspace string, network string, routeTable string) string {
	return fmt.Sprintf(RouteTableURLV1, tenant, workspace, network, routeTable)
}

func GenerateRouteTableListURL(tenant string, workspace string, network string) string {
	return fmt.Sprintf(RouteTableListURLV1, tenant, workspace, network)
}

func GenerateSubnetURL(tenant string, workspace string, network string, subnet string) string {
	return fmt.Sprintf(SubnetURLV1, tenant, workspace, network, subnet)
}

func GenerateSubnetListURL(tenant string, workspace string, network string) string {
	return fmt.Sprintf(SubnetListURLV1, tenant, workspace, network)
}

func GenerateSecurityGroupURL(tenant string, workspace string, securityGroup string) string {
	return fmt.Sprintf(SecurityGroupURLV1, tenant, workspace, securityGroup)
}

func GenerateSecurityGroupListURL(tenant string, workspace string) string {
	return fmt.Sprintf(SecurityGroupListURLV1, tenant, workspace)
}

// Providers
func GenerateProviderSpec() []schema.Provider {
	return []schema.Provider{
		{
			Name:    AuthorizationProvider,
			Version: ApiVersion1,
			Url:     GenerateRegionProviderUrl(AuthorizationProvider),
		},
		{
			Name:    ComputeProvider,
			Version: ApiVersion1,
			Url:     GenerateRegionProviderUrl(ComputeProvider),
		},
		{
			Name:    NetworkProvider,
			Version: ApiVersion1,
			Url:     GenerateRegionProviderUrl(NetworkProvider),
		},
		{
			Name:    StorageProvider,
			Version: ApiVersion1,
			Url:     GenerateRegionProviderUrl(StorageProvider),
		},
		{
			Name:    WorkspaceProvider,
			Version: ApiVersion1,
			Url:     GenerateRegionProviderUrl(WorkspaceProvider),
		},
	}
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
