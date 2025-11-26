package secalib

import (
	"fmt"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// URLs

// Remove the v1 dependency to support multiple API versions
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

func GenerateStorageSkuListURL(tenant string) string {
	return fmt.Sprintf(StorageSkuListURLV1, tenant)
}

func GenerateBlockStorageURL(tenant string, workspace string, blockStorage string) string {
	return fmt.Sprintf(BlockStorageURLV1, tenant, workspace, blockStorage)
}

func GenerateBlockStorageListURL(tenant string, workspace string) string {
	return fmt.Sprintf(BlockStorageListURLV1, tenant, workspace)
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
