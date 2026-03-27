package generators

import (
	"fmt"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Refs

func GenerateRegionRef(provider, name string) string {
	return fmt.Sprintf(regionRef, provider, name)
}

func GenerateSkuRef(provider, tenant, name string) string {
	return fmt.Sprintf(skuRef, provider, tenant, name)
}

func GenerateRoleRef(provider, tenant, name string) string {
	return fmt.Sprintf(roleRef, provider, tenant, name)
}

func GenerateRoleAssignmentRef(provider, tenant, name string) string {
	return fmt.Sprintf(roleAssignmentRef, provider, tenant, name)
}

func GenerateWorkspaceRef(provider, tenant, name string) string {
	return fmt.Sprintf(workspaceRef, provider, tenant, name)
}

func GenerateInstanceRef(provider, tenant, workspace, name string) string {
	return fmt.Sprintf(instanceRef, provider, tenant, workspace, name)
}

func GenerateBlockStorageRef(provider, tenant, workspace, name string) string {
	return fmt.Sprintf(blockStorageRef, provider, tenant, workspace, name)
}

func GenerateImageRef(provider, tenant, name string) string {
	return fmt.Sprintf(imageRef, provider, tenant, name)
}

func GenerateNetworkRef(provider, tenant, workspace, name string) string {
	return fmt.Sprintf(networkRef, provider, tenant, workspace, name)
}

func GenerateInternetGatewayRef(provider, tenant, workspace, name string) string {
	return fmt.Sprintf(internetGatewayRef, provider, tenant, workspace, name)
}

func GenerateNicRef(provider, tenant, workspace, name string) string {
	return fmt.Sprintf(nicRef, provider, tenant, workspace, name)
}

func GenerateRouteTableRef(provider, tenant, workspace, network, name string) string {
	return fmt.Sprintf(routeTableRef, provider, tenant, workspace, network, name)
}

func GenerateSubnetRef(provider, tenant, workspace, network, name string) string {
	return fmt.Sprintf(subnetRef, provider, tenant, workspace, network, name)
}

func GeneratePublicIpRef(provider, tenant, workspace, name string) string {
	return fmt.Sprintf(publicIpRef, provider, tenant, workspace, name)
}

func GenerateSecurityGroupRuleRef(provider, tenant, workspace, name string) string {
	return fmt.Sprintf(securityGroupRuleRef, provider, tenant, workspace, name)
}

func GenerateSecurityGroupRef(provider, tenant, workspace, name string) string {
	return fmt.Sprintf(securityGroupRef, provider, tenant, workspace, name)
}

// RefObjects

func GenerateSkuRefObject(provider, tenant, name string) *schema.Reference {
	urn := GenerateSkuRef(provider, tenant, name)
	return &schema.Reference{Resource: urn}
}

func GenerateInstanceRefObject(provider, tenant, workspace, name string) *schema.Reference {
	urn := GenerateInstanceRef(provider, tenant, workspace, name)
	return &schema.Reference{Resource: urn}
}

func GenerateBlockStorageRefObject(provider, tenant, workspace, name string) *schema.Reference {
	urn := GenerateBlockStorageRef(provider, tenant, workspace, name)
	return &schema.Reference{Resource: urn}
}

func GenerateNetworkRefObject(provider, tenant, workspace, name string) *schema.Reference {
	urn := GenerateNetworkRef(provider, tenant, workspace, name)
	return &schema.Reference{Resource: urn}
}

func GenerateInternetGatewayRefObject(provider, tenant, workspace, name string) *schema.Reference {
	urn := GenerateInternetGatewayRef(provider, tenant, workspace, name)
	return &schema.Reference{Resource: urn}
}

func GenerateRouteTableRefObject(provider, tenant, workspace, network, name string) *schema.Reference {
	urn := GenerateRouteTableRef(provider, tenant, workspace, network, name)
	return &schema.Reference{Resource: urn}
}

func GenerateSubnetRefObject(provider, tenant, workspace, network, name string) *schema.Reference {
	urn := GenerateSubnetRef(provider, tenant, workspace, network, name)
	return &schema.Reference{Resource: urn}
}

func GeneratePublicIpRefObject(provider, tenant, workspace, name string) *schema.Reference {
	urn := GeneratePublicIpRef(provider, tenant, workspace, name)
	return &schema.Reference{Resource: urn}
}

func GenerateSecurityGroupRuleRefObject(provider, tenant, workspace, name string) *schema.Reference {
	urn := GenerateSecurityGroupRuleRef(provider, tenant, workspace, name)
	return &schema.Reference{Resource: urn}
}
