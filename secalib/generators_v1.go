package secalib

import (
	"fmt"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// URLs

// Remove the v1 dependency to support multiple API versions
func GenerateRoleURL(tenant string, role string) string {
	return fmt.Sprintf(roleURLV1, tenant, role)
}

func GenerateRoleListURL(tenant string) string {
	return fmt.Sprintf(roleListURLV1, tenant)
}

func GenerateRoleAssignmentURL(tenant string, roleAssignment string) string {
	return fmt.Sprintf(roleAssignmentURLV1, tenant, roleAssignment)
}

func GenerateRoleAssignmentListURL(tenant string) string {
	return fmt.Sprintf(roleAssignmentListURLV1, tenant)
}

func GenerateRegionURL(region string) string {
	return fmt.Sprintf(regionURLV1, region)
}

func GenerateWorkspaceURL(tenant string, workspace string) string {
	return fmt.Sprintf(workspaceURLV1, tenant, workspace)
}

func GenerateWorkspaceListURL(tenant string) string {
	return fmt.Sprintf(workspaceListURLV1, tenant)
}

func GenerateStorageSkuURL(tenant string, sku string) string {
	return fmt.Sprintf(storageSkuURLV1, tenant, sku)
}

func GenerateStorageSkuListURL(tenant string) string {
	return fmt.Sprintf(storageSkuListURLV1, tenant)
}

func GenerateBlockStorageURL(tenant string, workspace string, blockStorage string) string {
	return fmt.Sprintf(blockStorageURLV1, tenant, workspace, blockStorage)
}

func GenerateBlockStorageListURL(tenant string, workspace string) string {
	return fmt.Sprintf(blockStorageListURLV1, tenant, workspace)
}

func GenerateImageURL(tenant string, image string) string {
	return fmt.Sprintf(imageURLV1, tenant, image)
}

func GenerateImageListURL(tenant string) string {
	return fmt.Sprintf(imageListURLV1, tenant)
}

func GenerateInstanceSkuURL(tenant string, sku string) string {
	return fmt.Sprintf(instanceSkuURLV1, tenant, sku)
}

func GenerateInstanceURL(tenant string, workspace string, instance string) string {
	return fmt.Sprintf(instanceURLV1, tenant, workspace, instance)
}

func GenerateInstanceListURL(tenant string, workspace string) string {
	return fmt.Sprintf(instanceListURLV1, tenant, workspace)
}

func GenerateNetworkURL(tenant string, workspace string, network string) string {
	return fmt.Sprintf(networkURLV1, tenant, workspace, network)
}

func GenerateNetworkListURL(tenant string, workspace string) string {
	return fmt.Sprintf(networkListURLV1, tenant, workspace)
}

func GenerateNetworkSkuURL(tenant string, sku string) string {
	return fmt.Sprintf(networkSkuURLV1, tenant, sku)
}

func GenerateInternetGatewayURL(tenant string, workspace string, internetGateway string) string {
	return fmt.Sprintf(internetGatewayURLV1, tenant, workspace, internetGateway)
}

func GenerateInternetGatewayListURL(tenant string, workspace string) string {
	return fmt.Sprintf(internetGatewayListURLV1, tenant, workspace)
}

func GenerateNicURL(tenant string, workspace string, nic string) string {
	return fmt.Sprintf(nicURLV1, tenant, workspace, nic)
}

func GenerateNicListURL(tenant string, workspace string) string {
	return fmt.Sprintf(nicListURLV1, tenant, workspace)
}

func GeneratePublicIpURL(tenant string, workspace string, publicIp string) string {
	return fmt.Sprintf(publicIpURLV1, tenant, workspace, publicIp)
}

func GeneratePublicIpListURL(tenant string, workspace string) string {
	return fmt.Sprintf(publicIpListURLV1, tenant, workspace)
}

func GenerateRouteTableURL(tenant string, workspace string, network string, routeTable string) string {
	return fmt.Sprintf(routeTableURLV1, tenant, workspace, network, routeTable)
}

func GenerateRouteTableListURL(tenant string, workspace string, network string) string {
	return fmt.Sprintf(routeTableListURLV1, tenant, workspace, network)
}

func GenerateSubnetURL(tenant string, workspace string, network string, subnet string) string {
	return fmt.Sprintf(subnetURLV1, tenant, workspace, network, subnet)
}

func GenerateSubnetListURL(tenant string, workspace string, network string) string {
	return fmt.Sprintf(subnetListURLV1, tenant, workspace, network)
}

func GenerateSecurityGroupURL(tenant string, workspace string, securityGroup string) string {
	return fmt.Sprintf(securityGroupURLV1, tenant, workspace, securityGroup)
}

func GenerateSecurityGroupListURL(tenant string, workspace string) string {
	return fmt.Sprintf(securityGroupListURLV1, tenant, workspace)
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
