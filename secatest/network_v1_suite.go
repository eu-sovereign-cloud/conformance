package secatest

import (
	"log/slog"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	network "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

const (
	routeTableDefaultDestination = "0.0.0.0/0"
)

type NetworkV1TestSuite struct {
	regionalTestSuite
}

func (suite *NetworkV1TestSuite) TestNetworkV1(t provider.T) {
	t.Title("Network Lifecycle Test")
	configureTags(t, secalib.NetworkProviderV1, secalib.NetworkKind, secalib.InternetGatewayKind, secalib.NicKind, secalib.PublicIPKind, secalib.RouteTableKind,
		secalib.SubnetKind, secalib.SecurityGroupKind)

	// TODO Receive via configuration the network cidr ranges and calculate the cidr
	networkCIDR := "10.1.0.0/16"

	// TODO Calculate the subnet cidr from network cidr
	subnetCIDR := "10.1.0.0/24"

	// TODO Calculate the nic cidr from subnet cidr
	nicAddress := "10.1.0.1"

	// TODO Receive via configuration the public ip address range and calculate an ip address
	publicIPAddress := "192.168.0.1"

	// TODO Load from region
	zones := []string{"zone-a", "zone-b"}

	// TODO Get from list storage skus endpoint
	storageSkus := []string{"LD100"}

	// TODO Get from list instance skus endpoint
	instanceSkus := []string{"SXL"}

	// TODO Get from list network skus endpoint
	networkSkus := []string{"N1K"}

	// Select skus
	storageSkuName := storageSkus[rand.Intn(len(storageSkus))]
	instanceSkuName := instanceSkus[rand.Intn(len(instanceSkus))]
	networkSkuName := networkSkus[rand.Intn(len(networkSkus))]
	zone := zones[rand.Intn(len(zones))]

	// TODO Dynamically create before the scenario
	workspaceName := secalib.GenerateWorkspaceName()

	// Generate scenario data
	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)

	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	instanceName := secalib.GenerateInstanceName()

	networkSkuRef := secalib.GenerateSkuRef(networkSkuName)
	networkName := secalib.GenerateNetworkName()
	networkRef := secalib.GenerateNetworkRef(networkName)

	internetGatewayName := secalib.GenerateInternetGatewayName()
	internetGatewayRef := secalib.GenerateInternetGatewayRef(internetGatewayName)

	routeTableName := secalib.GenerateRouteTableName()
	routeTableRef := secalib.GenerateRouteTableRef(routeTableName)

	subnetName := secalib.GenerateSubnetName()
	subnetRef := secalib.GenerateSubnetRef(subnetName)

	nicName := secalib.GenerateNicName()

	publicIPName := secalib.GeneratePublicIPName()
	publicIPRef := secalib.GeneratePublicIPRef(publicIPName)

	securityGroupName := secalib.GenerateSecurityGroupName()

	blockStorageSize := secalib.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		wm, err := mock.CreateNetworkLifecycleScenarioV1("Network Lifecycle",
			mock.NetworkParamsV1{
				Params: &mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
					Workspace: workspaceName,
					Region:    suite.region,
				},
				BlockStorage: &mock.ResourceParams[secalib.BlockStorageSpecV1]{
					Name: blockStorageName,
					InitialSpec: &secalib.BlockStorageSpecV1{
						SkuRef: storageSkuRef,
						SizeGB: blockStorageSize,
					},
				},
				Instance: &mock.ResourceParams[secalib.InstanceSpecV1]{
					Name: instanceName,
					InitialSpec: &secalib.InstanceSpecV1{
						SkuRef:        instanceSkuRef,
						Zone:          zone,
						BootDeviceRef: blockStorageRef,
					},
				},
				Network: &mock.ResourceParams[secalib.NetworkSpecV1]{
					Name: networkName,
					InitialSpec: &secalib.NetworkSpecV1{
						Cidr: &secalib.NetworkSpecCIDRV1{
							Ipv4: networkCIDR},
						SkuRef:        networkSkuRef,
						RouteTableRef: routeTableRef,
					},
				},
				InternetGateway: &mock.ResourceParams[secalib.InternetGatewaySpecV1]{
					Name: internetGatewayName,
					InitialSpec: &secalib.InternetGatewaySpecV1{
						EgressOnly: false,
					},
				},
				RouteTable: &mock.ResourceParams[secalib.RouteTableSpecV1]{
					Name: routeTableName,
					InitialSpec: &secalib.RouteTableSpecV1{
						LocalRef: networkRef,
						Routes: []*secalib.RouteTableRouteV1{
							{
								DestinationCidrBlock: routeTableDefaultDestination,
								TargetRef:            internetGatewayRef,
							},
						},
					},
				},
				Subnet: &mock.ResourceParams[secalib.SubnetSpecV1]{
					Name: subnetName,
					InitialSpec: &secalib.SubnetSpecV1{
						Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCIDR},
						Zone: zone,
					},
				},
				NIC: &mock.ResourceParams[secalib.NICSpecV1]{
					Name: nicName,
					InitialSpec: &secalib.NICSpecV1{
						Addresses:    []string{nicAddress},
						PublicIpRefs: []string{publicIPRef},
						SubnetRef:    subnetRef,
					},
				},
				PublicIP: &mock.ResourceParams[secalib.PublicIPSpecV1]{
					Name: publicIPName,
					InitialSpec: &secalib.PublicIPSpecV1{
						Version: secalib.IPVersion4,
						Address: publicIPAddress,
					},
				},
				SecurityGroup: &mock.ResourceParams[secalib.SecurityGroupSpecV1]{
					Name: securityGroupName,
					InitialSpec: &secalib.SecurityGroupSpecV1{
						Rules: []*secalib.SecurityGroupRule{
							{
								Direction: secalib.SecurityRuleDirectionIngress},
						},
					},
				},
			})
		if err != nil {
			slog.Error("Failed to create network scenario", "error", err)
			return
		}
		suite.mockClient = wm
	}
}

func (suite *NetworkV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}

func verifyNetworkZonalMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *network.ZonalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Workspace:  *metadata.Workspace,
		Region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyNetworkRegionalMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *network.RegionalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Workspace:  *metadata.Workspace,
		Region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}
