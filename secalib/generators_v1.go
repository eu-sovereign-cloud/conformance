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

func GenerateRoleAssignmentURL(tenant string, roleAssignment string) string {
	return fmt.Sprintf(RoleAssignmentURLV1, tenant, roleAssignment)
}

func GenerateRegionURL(region string) string {
	return fmt.Sprintf(RegionURLV1, region)
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

func GeneratePublicIpURL(tenant string, workspace string, publicIp string) string {
	return fmt.Sprintf(PublicIpURLV1, tenant, workspace, publicIp)
}

func GenerateRouteTableURL(tenant string, workspace string, network string, routeTable string) string {
	return fmt.Sprintf(RouteTableURLV1, tenant, workspace, network, routeTable)
}

func GenerateSubnetURL(tenant string, workspace string, network string, subnet string) string {
	return fmt.Sprintf(SubnetURLV1, tenant, workspace, network, subnet)
}

func GenerateSecurityGroupURL(tenant string, workspace string, securityGroup string) string {
	return fmt.Sprintf(SecurityGroupURLV1, tenant, workspace, securityGroup)
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
