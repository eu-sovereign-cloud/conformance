package generators

import (
	"fmt"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
)

func generateSkuRef(name string) string {
	return fmt.Sprintf(skuRef, name)
}

func GenerateSkuRefObject(name string) (*schema.Reference, error) {
	urn := generateSkuRef(name)

	ref, err := secapi.BuildReferenceFromURN(urn)
	if err != nil {
		return nil, fmt.Errorf("error building reference %s: %s", urn, err)
	}
	return ref, nil
}

func generateInstanceRef(instanceName string) string {
	return fmt.Sprintf(instanceRef, instanceName)
}

func GenerateInstanceRefObject(name string) (*schema.Reference, error) {
	urn := generateInstanceRef(name)

	ref, err := secapi.BuildReferenceFromURN(urn)
	if err != nil {
		return nil, fmt.Errorf("error building reference %s: %s", urn, err)
	}
	return ref, nil
}

func generateBlockStorageRef(blockStorageName string) string {
	return fmt.Sprintf(blockStorageRef, blockStorageName)
}

func GenerateBlockStorageRefObject(name string) (*schema.Reference, error) {
	urn := generateBlockStorageRef(name)

	ref, err := secapi.BuildReferenceFromURN(urn)
	if err != nil {
		return nil, fmt.Errorf("error building reference %s: %s", urn, err)
	}
	return ref, nil
}

func generateInternetGatewayRef(internetGatewayName string) string {
	return fmt.Sprintf(internetGatewayRef, internetGatewayName)
}

func GenerateInternetGatewayRefObject(name string) (*schema.Reference, error) {
	urn := generateInternetGatewayRef(name)

	ref, err := secapi.BuildReferenceFromURN(urn)
	if err != nil {
		return nil, fmt.Errorf("error building reference %s: %s", urn, err)
	}
	return ref, nil
}

func generateNetworkRef(networkName string) string {
	return fmt.Sprintf(networkRef, networkName)
}

func GenerateNetworkRefObject(name string) (*schema.Reference, error) {
	urn := generateNetworkRef(name)

	ref, err := secapi.BuildReferenceFromURN(urn)
	if err != nil {
		return nil, fmt.Errorf("error building reference %s: %s", urn, err)
	}
	return ref, nil
}

func generateRouteTableRef(routeTableName string) string {
	return fmt.Sprintf(routeTableRef, routeTableName)
}

func GenerateRouteTableRefObject(name string) (*schema.Reference, error) {
	urn := generateRouteTableRef(name)

	ref, err := secapi.BuildReferenceFromURN(urn)
	if err != nil {
		return nil, fmt.Errorf("error building reference %s: %s", urn, err)
	}
	return ref, nil
}

func generateSubnetRef(subnetName string) string {
	return fmt.Sprintf(subnetRef, subnetName)
}

func GenerateSubnetRefObject(name string) (*schema.Reference, error) {
	urn := generateSubnetRef(name)

	ref, err := secapi.BuildReferenceFromURN(urn)
	if err != nil {
		return nil, fmt.Errorf("error building reference %s: %s", urn, err)
	}
	return ref, nil
}

func generatePublicIpRef(publicIpName string) string {
	return fmt.Sprintf(publicIpRef, publicIpName)
}

func GeneratePublicIpRefObject(name string) (*schema.Reference, error) {
	urn := generatePublicIpRef(name)

	ref, err := secapi.BuildReferenceFromURN(urn)
	if err != nil {
		return nil, fmt.Errorf("error building reference %s: %s", urn, err)
	}
	return ref, nil
}
