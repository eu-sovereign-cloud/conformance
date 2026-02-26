package generators

import (
	"fmt"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Refs

func GenerateRegionRef(name string) string {
	return fmt.Sprintf(regionRef, name)
}

func GenerateSkuRef(name string) string {
	return fmt.Sprintf(skuRef, name)
}

func GenerateRoleRef(name string) string {
	return fmt.Sprintf(roleRef, name)
}

func GenerateRoleAssignmentRef(name string) string {
	return fmt.Sprintf(roleAssignmentRef, name)
}

func GenerateWorkspaceRef(name string) string {
	return fmt.Sprintf(workspaceRef, name)
}

func GenerateInstanceRef(name string) string {
	return fmt.Sprintf(instanceRef, name)
}

func GenerateBlockStorageRef(name string) string {
	return fmt.Sprintf(blockStorageRef, name)
}

func GenerateImageRef(name string) string {
	return fmt.Sprintf(imageRef, name)
}

func GenerateNetworkRef(name string) string {
	return fmt.Sprintf(networkRef, name)
}

func GenerateInternetGatewayRef(name string) string {
	return fmt.Sprintf(internetGatewayRef, name)
}

func GenerateNicRef(name string) string {
	return fmt.Sprintf(nicRef, name)
}

func GenerateRouteTableRef(name string) string {
	return fmt.Sprintf(routeTableRef, name)
}

func GenerateSubnetRef(name string) string {
	return fmt.Sprintf(subnetRef, name)
}

func GeneratePublicIpRef(name string) string {
	return fmt.Sprintf(publicIpRef, name)
}

func GenerateSecurityGroupRuleRef(name string) string {
	return fmt.Sprintf(securityGroupRuleRef, name)
}

func GenerateSecurityGroupRef(name string) string {
	return fmt.Sprintf(securityGroupRef, name)
}

// RefObjects

func GenerateSkuRefObject(name string) *schema.Reference {
	urn := GenerateSkuRef(name)
	return &schema.Reference{Resource: urn}
}

func GenerateInstanceRefObject(name string) *schema.Reference {
	urn := GenerateInstanceRef(name)
	return &schema.Reference{Resource: urn}
}

func GenerateBlockStorageRefObject(name string) *schema.Reference {
	urn := GenerateBlockStorageRef(name)
	return &schema.Reference{Resource: urn}
}

func GenerateNetworkRefObject(name string) *schema.Reference {
	urn := GenerateNetworkRef(name)
	return &schema.Reference{Resource: urn}
}

func GenerateInternetGatewayRefObject(name string) *schema.Reference {
	urn := GenerateInternetGatewayRef(name)
	return &schema.Reference{Resource: urn}
}

func GenerateRouteTableRefObject(name string) *schema.Reference {
	urn := GenerateRouteTableRef(name)
	return &schema.Reference{Resource: urn}
}

func GenerateSubnetRefObject(name string) *schema.Reference {
	urn := GenerateSubnetRef(name)
	return &schema.Reference{Resource: urn}
}

func GeneratePublicIpRefObject(name string) *schema.Reference {
	urn := GeneratePublicIpRef(name)
	return &schema.Reference{Resource: urn}
}
