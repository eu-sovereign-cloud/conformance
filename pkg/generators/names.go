package generators

import (
	"fmt"
	"math"
	"math/rand"
)

func GenerateRoleName() string {
	return fmt.Sprintf(roleName, rand.Intn(math.MaxInt32))
}

func GenerateRoleAssignmentName() string {
	return fmt.Sprintf(roleAssignmentName, rand.Intn(math.MaxInt32))
}

func GenerateRegionName() string {
	return fmt.Sprintf(regionName, rand.Intn(math.MaxInt32))
}

func GenerateWorkspaceName() string {
	return fmt.Sprintf(workspaceName, rand.Intn(math.MaxInt32))
}

func GenerateBlockStorageName() string {
	return fmt.Sprintf(blockStorageName, rand.Intn(math.MaxInt32))
}

func GenerateImageName() string {
	return fmt.Sprintf(imageName, rand.Intn(math.MaxInt32))
}

func GenerateInstanceName() string {
	return fmt.Sprintf(instanceName, rand.Intn(math.MaxInt32))
}

func GenerateNetworkName() string {
	return fmt.Sprintf(networkName, rand.Intn(math.MaxInt32))
}

func GenerateInternetGatewayName() string {
	return fmt.Sprintf(internetGatewayName, rand.Intn(math.MaxInt32))
}

func GenerateRouteTableName() string {
	return fmt.Sprintf(routeTableName, rand.Intn(math.MaxInt32))
}

func GenerateSubnetName() string {
	return fmt.Sprintf(subnetName, rand.Intn(math.MaxInt32))
}

func GeneratePublicIpName() string {
	return fmt.Sprintf(publicIpName, rand.Intn(math.MaxInt32))
}

func GenerateNicName() string {
	return fmt.Sprintf(nicName, rand.Intn(math.MaxInt32))
}

func GenerateSecurityGroupRuleName() string {
	return fmt.Sprintf(securityGroupRulesName, rand.Intn(math.MaxInt32))
}


func GenerateSecurityGroupName() string {
	return fmt.Sprintf(securityGroupName, rand.Intn(math.MaxInt32))
}
