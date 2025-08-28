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
				Params: mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
					Workspace: workspaceName,
					Region:    suite.region,
				},
				BlockStorage: mock.BlockStorageParamsV1{
					Name:          blockStorageName,
					SkuRef:        storageSkuRef,
					SizeGBInitial: blockStorageSize,
				},
				Instance: mock.InstanceParamsV1{
					Name:          instanceName,
					SkuRef:        instanceSkuRef,
					ZoneInitial:   zone,
					BootDeviceRef: blockStorageRef,
				},
				Network: mock.NetworkInstanceParamsV1{
					Name:          networkName,
					Cidr:          mock.CIDR{Ipv4: networkCIDR},
					SkuRef:        networkSkuRef,
					RouteTableRef: routeTableRef,
				},
				InternetGateway: mock.InternetGatewayParamsV1{
					Name:       internetGatewayName,
					EgressOnly: false,
				},
				RouteTable: mock.RouteTableParamsV1{
					Name:     routeTableName,
					LocalRef: networkRef,
					Routes: []mock.RouteTableRoute{
						{
							DestinationCidrBlock: routeTableDefaultDestination,
							TargetRef:            internetGatewayRef,
						},
					},
				},
				Subnet: mock.SubnetParamsV1{
					Name: subnetName,
					Cidr: mock.CIDR{Ipv4: subnetCIDR},
					Zone: zone,
				},
				NIC: mock.NICParamsV1{
					Name:         nicName,
					Addresses:    []string{nicAddress},
					PublicIpRefs: []string{publicIPRef},
					SubnetRef:    subnetRef,
				},
				PublicIP: mock.PublicIPParamsV1{
					Name:    publicIPName,
					Version: secalib.IPVersion4,
					Address: publicIPAddress,
				},
				SecurityGroup: mock.SecurityGroupParamsV1{
					Name:  securityGroupName,
					Rules: []mock.SecurityGroupRule{{Direction: secalib.SecurityRuleDirectionIngress}},
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

func verifyNetworkZonalMetadataStep(ctx provider.StepCtx, expected verifyRegionalMetadataStepParams, metadata *network.ZonalResourceMetadata) {
	actualMetadata := verifyRegionalMetadataStepParams{
		name:       metadata.Name,
		provider:   metadata.Provider,
		verb:       metadata.Verb,
		resource:   metadata.Resource,
		apiVersion: metadata.ApiVersion,
		kind:       string(metadata.Kind),
		tenant:     metadata.Tenant,
		workspace:  *metadata.Workspace,
		region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyNetworkRegionalMetadataStep(ctx provider.StepCtx, expected verifyRegionalMetadataStepParams, metadata *network.RegionalResourceMetadata) {
	actualMetadata := verifyRegionalMetadataStepParams{
		name:       metadata.Name,
		provider:   metadata.Provider,
		verb:       metadata.Verb,
		resource:   metadata.Resource,
		apiVersion: metadata.ApiVersion,
		kind:       string(metadata.Kind),
		tenant:     metadata.Tenant,
		workspace:  *metadata.Workspace,
		region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}
