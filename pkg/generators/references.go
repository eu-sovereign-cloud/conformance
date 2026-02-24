package generators

import (
	"fmt"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func generateSkuRef(name string) string {
	return fmt.Sprintf(skuRef, name)
}

func GenerateSkuRefObject(name string) *schema.Reference {
	urn := generateSkuRef(name)
	return &schema.Reference{Resource: urn}
}

func generateInstanceRef(instanceName string) string {
	return fmt.Sprintf(instanceRef, instanceName)
}

func GenerateInstanceRefObject(name string) *schema.Reference {
	urn := generateInstanceRef(name)
	return &schema.Reference{Resource: urn}
}

func generateBlockStorageRef(blockStorageName string) string {
	return fmt.Sprintf(blockStorageRef, blockStorageName)
}

func GenerateBlockStorageRefObject(name string) *schema.Reference {
	urn := generateBlockStorageRef(name)
	return &schema.Reference{Resource: urn}
}

func generateInternetGatewayRef(internetGatewayName string) string {
	return fmt.Sprintf(internetGatewayRef, internetGatewayName)
}

func GenerateInternetGatewayRefObject(name string) *schema.Reference {
	urn := generateInternetGatewayRef(name)
	return &schema.Reference{Resource: urn}
}

func generateNetworkRef(networkName string) string {
	return fmt.Sprintf(networkRef, networkName)
}

func GenerateNetworkRefObject(name string) *schema.Reference {
	urn := generateNetworkRef(name)
	return &schema.Reference{Resource: urn}
}

func generateRouteTableRef(routeTableName string) string {
	return fmt.Sprintf(routeTableRef, routeTableName)
}

func GenerateRouteTableRefObject(name string) *schema.Reference {
	urn := generateRouteTableRef(name)
	return &schema.Reference{Resource: urn}
}

func generateSubnetRef(subnetName string) string {
	return fmt.Sprintf(subnetRef, subnetName)
}

func GenerateSubnetRefObject(name string) *schema.Reference {
	urn := generateSubnetRef(name)
	return &schema.Reference{Resource: urn}
}

func generatePublicIpRef(publicIpName string) string {
	return fmt.Sprintf(publicIpRef, publicIpName)
}

func GeneratePublicIpRefObject(name string) *schema.Reference {
	urn := generatePublicIpRef(name)
	return &schema.Reference{Resource: urn}
}
