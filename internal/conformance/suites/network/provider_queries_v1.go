package network

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockNetwork "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/network"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"k8s.io/utils/ptr"
)

type ProviderQueriesV1TestSuite struct {
	suites.RegionalTestSuite

	config *ProviderQueriesV1Config
	params *params.NetworkProviderQueriesV1Params
}

type ProviderQueriesV1Config struct {
	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	StorageSkus    []string
	InstanceSkus   []string
	NetworkSkus    []string
}

func CreateProviderQueriesV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *ProviderQueriesV1Config) *ProviderQueriesV1TestSuite {
	suite := &ProviderQueriesV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.NetworkProviderQueriesV1SuiteName.String()
	return suite
}

func (suite *ProviderQueriesV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Network")

	// Generate the subnet cidr
	subnetCidr, err := generators.GenerateSubnetCidr(suite.config.NetworkCidr, 8, 1)
	if err != nil {
		t.Fatalf("Failed to generate subnet cidr: %v", err)
	}

	// Generate the nic addresses
	nicAddress1, err := generators.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}

	// Generate the public ips
	publicIpAddress1, err := generators.GeneratePublicIp(suite.config.PublicIpsRange, 1)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}

	// Select zones
	zone := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]

	// Select skus
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	networkSkuName1 := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj, err := generators.GenerateSkuRefObject(storageSkuName)
	if err != nil {
		t.Fatal(err)
	}

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj, err := generators.GenerateBlockStorageRefObject(blockStorageName)
	if err != nil {
		t.Fatal(err)
	}

	instanceSkuRefObj, err := generators.GenerateSkuRefObject(instanceSkuName)
	if err != nil {
		t.Fatal(err)
	}

	instanceName := generators.GenerateInstanceName()

	networkSkuRefObj, err := generators.GenerateSkuRefObject(networkSkuName1)
	if err != nil {
		t.Fatal(err)
	}

	networkName := generators.GenerateNetworkName()
	networkName2 := generators.GenerateNetworkName()

	internetGatewayName := generators.GenerateInternetGatewayName()
	internetGatewayName2 := generators.GenerateInternetGatewayName()
	internetGatewayRefObj, err := generators.GenerateInternetGatewayRefObject(internetGatewayName)
	if err != nil {
		t.Fatal(err)
	}

	routeTableName := generators.GenerateRouteTableName()
	routeTableName2 := generators.GenerateRouteTableName()
	routeTableRefObj, err := generators.GenerateRouteTableRefObject(routeTableName)
	if err != nil {
		t.Fatal(err)
	}

	subnetName := generators.GenerateSubnetName()
	subnetName2 := generators.GenerateSubnetName()
	subnetRefObj, err := generators.GenerateSubnetRefObject(subnetName)
	if err != nil {
		t.Fatal(err)
	}

	nicName := generators.GenerateNicName()
	nicName2 := generators.GenerateNicName()

	publicIpName := generators.GeneratePublicIpName()
	publicIpName2 := generators.GeneratePublicIpName()
	publicIpRefObj, err := generators.GeneratePublicIpRefObject(publicIpName)
	if err != nil {
		t.Fatal(err)
	}

	securityGroupName := generators.GenerateSecurityGroupName()
	securityGroupName2 := generators.GenerateSecurityGroupName()

	blockStorageSize := generators.GenerateBlockStorageSize()

	// Workspace
	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: blockStorageSize,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}
	instance, err := builders.NewInstanceBuilder().
		Name(instanceName).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	network, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: ptr.To(suite.config.NetworkCidr)},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	network2, err := builders.NewNetworkBuilder().
		Name(networkName2).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: ptr.To(suite.config.NetworkCidr)},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	networks := []schema.Network{*network, *network2}

	internetGateway, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: ptr.To(false),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	internetGateway2, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName2).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: ptr.To(false),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	internetGateways := []schema.InternetGateway{*internetGateway, *internetGateway2}

	routeTable, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Route Table: %v", err)
	}

	routeTable2, err := builders.NewRouteTableBuilder().
		Name(routeTableName2).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Route Table: %v", err)
	}

	routeTables := []schema.RouteTable{*routeTable, *routeTable2}

	subnet, err := builders.NewSubnetBuilder().
		Name(subnetName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	subnet2, err := builders.NewSubnetBuilder().
		Name(subnetName2).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	subnets := []schema.Subnet{*subnet, *subnet2}

	nic, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	nic2, err := builders.NewNicBuilder().
		Name(nicName2).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	nics := []schema.Nic{*nic, *nic2}

	publicIp, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: ptr.To(publicIpAddress1),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	publicIp2, err := builders.NewPublicIpBuilder().
		Name(publicIpName2).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: ptr.To(publicIpAddress1),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	publicIps := []schema.PublicIp{*publicIp, *publicIp2}

	securityGroup, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	securityGroup2, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName2).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	securityGroups := []schema.SecurityGroup{*securityGroup, *securityGroup2}

	params := &params.NetworkProviderQueriesV1Params{
		Workspace:        workspace,
		BlockStorage:     blockStorage,
		Instance:         instance,
		Networks:         networks,
		InternetGateways: internetGateways,
		RouteTables:      routeTables,
		Subnets:          subnets,
		Nics:             nics,
		PublicIps:        publicIps,
		SecurityGroups:   securityGroups,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureProviderQueriesV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ProviderQueriesV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.NetworkProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace
	workspace := suite.params.Workspace

	// Create a workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Network
	networks := suite.params.Networks

	// Create networks
	for _, network := range networks {
		expectNetworkMeta := network.Metadata
		expectNetworkSpec := &network.Spec
		stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", suite.Client.NetworkV1, &network,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
				Metadata:      expectNetworkMeta,
				Spec:          expectNetworkSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List networks
	wref := secapi.WorkspaceReference{
		Name:      workspace.Metadata.Name,
		Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
		Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
	}
	stepsBuilder.GetListNetworkV1Step("List Network", suite.Client.NetworkV1, wref, nil)

	// List networks with limit
	stepsBuilder.GetListNetworkV1Step("Get list of Network with limit", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List networks with label
	stepsBuilder.GetListNetworkV1Step("Get list of Network with label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List networks with limit and label
	stepsBuilder.GetListNetworkV1Step("Get list of Network with limit and label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Skus

	// List skus
	stepsBuilder.GetListNetworkSkusV1Step("List skus", suite.Client.NetworkV1, secapi.TenantReference{Tenant: secapi.TenantID(workspace.Metadata.Tenant)}, nil)

	// List skus with limit
	stepsBuilder.GetListNetworkSkusV1Step("Get list of skus", suite.Client.NetworkV1, secapi.TenantReference{Tenant: secapi.TenantID(workspace.Metadata.Tenant)},
		secapi.NewListOptions().WithLimit(1))

	// Internet gateway
	gateways := suite.params.InternetGateways

	// Create internet gateways
	for _, gateway := range gateways {
		expectGatewayMeta := gateway.Metadata
		expectGatewaySpec := &gateway.Spec
		stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create a internet gateway", suite.Client.NetworkV1, &gateway,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
				Metadata:      expectGatewayMeta,
				Spec:          expectGatewaySpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

	}

	// List internet gateways
	stepsBuilder.GetListInternetGatewayV1Step("List Internet Gateway", suite.Client.NetworkV1, wref, nil)

	// List internet gateways with limit
	stepsBuilder.GetListInternetGatewayV1Step("Get list of Internet Gateway with limit", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List internet gateways with label
	stepsBuilder.GetListInternetGatewayV1Step("Get list of Internet Gateway with label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List internet gateways with limit and label
	stepsBuilder.GetListInternetGatewayV1Step("Get list of Internet Gateway with limit and label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Route table
	routes := suite.params.RouteTables

	// Create route tables
	for _, route := range routes {
		expectRouteMeta := route.Metadata
		expectRouteSpec := &route.Spec
		stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", suite.Client.NetworkV1, &route,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
				Metadata:      expectRouteMeta,
				Spec:          expectRouteSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List route tables
	nref := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		Workspace: secapi.WorkspaceID((workspace.Metadata.Name)),
		Network:   secapi.NetworkID(networks[0].Metadata.Name),
		Name:      routes[0].Metadata.Name,
	}
	stepsBuilder.GetListRouteTableV1Step("List Route table", suite.Client.NetworkV1, *nref, nil)

	// List route tables with limit
	stepsBuilder.GetListRouteTableV1Step("Get list of Route table with limit", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1))

	// List route tables with label
	stepsBuilder.GetListRouteTableV1Step("Get list of Route table with label", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List route tables with limit and label
	stepsBuilder.GetListRouteTableV1Step("Get list of Route table with limit and label", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Subnet
	subnets := suite.params.Subnets

	// Create subnets
	for _, subnet := range subnets {
		expectSubnetMeta := subnet.Metadata
		expectSubnetSpec := &subnet.Spec
		stepsBuilder.CreateOrUpdateSubnetV1Step("Create a subnet", suite.Client.NetworkV1, &subnet,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
				Metadata:      expectSubnetMeta,
				Spec:          expectSubnetSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List subnets
	stepsBuilder.GetListSubnetV1Step("List Subnet", suite.Client.NetworkV1, *nref, nil)

	// List subnets with limit
	stepsBuilder.GetListSubnetV1Step("Get list of Subnet with limit", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1))

	// List subnets with label
	stepsBuilder.GetListSubnetV1Step("Get list of Subnet with label", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List subnets with limit and label
	stepsBuilder.GetListSubnetV1Step("Get list of Subnet with limit and label", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Public ip
	publicIps := suite.params.PublicIps

	// Create public ips
	for _, publicIp := range publicIps {
		expectPublicIpMeta := publicIp.Metadata
		expectPublicIpSpec := &publicIp.Spec
		stepsBuilder.CreateOrUpdatePublicIpV1Step("Create a public ip", suite.Client.NetworkV1, &publicIp,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
				Metadata:      expectPublicIpMeta,
				Spec:          expectPublicIpSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List public ips
	stepsBuilder.GetListPublicIpV1Step("List PublicIP", suite.Client.NetworkV1, wref, nil)

	// List public ips with limit
	stepsBuilder.GetListPublicIpV1Step("Get list of PublicIP with limit", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List public ips with label
	stepsBuilder.GetListPublicIpV1Step("Get list of PublicIP with label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List public ips with limit and label
	stepsBuilder.GetListPublicIpV1Step("Get list of PublicIP with limit and label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Nic
	nics := suite.params.Nics

	// Create nics
	for _, nic := range nics {
		expectNicMeta := nic.Metadata
		expectNicSpec := &nic.Spec
		stepsBuilder.CreateOrUpdateNicV1Step("Create a nic", suite.Client.NetworkV1, &nic,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
				Metadata:      expectNicMeta,
				Spec:          expectNicSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List nics
	stepsBuilder.GetListNicV1Step("List Nic", suite.Client.NetworkV1, wref, nil)

	// List nics with limit
	stepsBuilder.GetListNicV1Step("Get list of Nic with limit", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List nics with label
	stepsBuilder.GetListNicV1Step("Get list of Nic with label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List nics with limit and label
	stepsBuilder.GetListNicV1Step("Get list of Nic with limit and label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Security Group
	groups := suite.params.SecurityGroups

	// Create security groups
	for _, group := range groups {
		expectGroupMeta := group.Metadata
		expectGroupSpec := &group.Spec
		stepsBuilder.CreateOrUpdateSecurityGroupV1Step("Create a security group", suite.Client.NetworkV1, &group,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
				Metadata:      expectGroupMeta,
				Spec:          expectGroupSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List security groups
	stepsBuilder.GetListSecurityGroupV1Step("List Security Group", suite.Client.NetworkV1, wref, nil)

	// List security groups with limit
	stepsBuilder.GetListSecurityGroupV1Step("Get list of Security Group with limit", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List security groups with label
	stepsBuilder.GetListSecurityGroupV1Step("Get list of Security Group with label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List security groups with limit and label
	stepsBuilder.GetListSecurityGroupV1Step("Get list of Security Group with limit and label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Delete all security groups
	for _, group := range groups {
		stepsBuilder.DeleteSecurityGroupV1Step("Delete the security group", suite.Client.NetworkV1, &group)

		// Get deleted security group
		groupWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(group.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(group.Metadata.Workspace),
			Name:      group.Metadata.Name,
		}
		stepsBuilder.GetSecurityGroupWithErrorV1Step("Get deleted security group", suite.Client.NetworkV1, *groupWRef, secapi.ErrResourceNotFound)
	}

	// Delete all nics
	for _, nic := range nics {
		stepsBuilder.DeleteNicV1Step("Delete the nic", suite.Client.NetworkV1, &nic)

		// Get the deleted nic
		nicWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(nic.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(nic.Metadata.Workspace),
			Name:      nic.Metadata.Name,
		}
		stepsBuilder.GetNicWithErrorV1Step("Get deleted nic", suite.Client.NetworkV1, *nicWRef, secapi.ErrResourceNotFound)
	}

	// Delete all public ips
	for _, publicIp := range publicIps {
		stepsBuilder.DeletePublicIpV1Step("Delete the public ip", suite.Client.NetworkV1, &publicIp)

		// Get the deleted public ip
		publicIpWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(publicIp.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(publicIp.Metadata.Workspace),
			Name:      publicIp.Metadata.Name,
		}
		stepsBuilder.GetPublicIpWithErrorV1Step("Get deleted public ip", suite.Client.NetworkV1, *publicIpWRef, secapi.ErrResourceNotFound)
	}

	// Delete all subnets
	for _, subnet := range subnets {
		stepsBuilder.DeleteSubnetV1Step("Delete the subnet", suite.Client.NetworkV1, &subnet)

		// Get the deleted subnet
		subnetNRef := &secapi.NetworkReference{
			Tenant:    secapi.TenantID(subnet.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(subnet.Metadata.Workspace),
			Network:   secapi.NetworkID(subnet.Metadata.Network),
			Name:      subnet.Metadata.Name,
		}
		stepsBuilder.GetSubnetWithErrorV1Step("Get deleted subnet", suite.Client.NetworkV1, *subnetNRef, secapi.ErrResourceNotFound)
	}

	// Delete all route tables
	for _, route := range routes {
		stepsBuilder.DeleteRouteTableV1Step("Delete the route table", suite.Client.NetworkV1, &route)

		// Get the deleted route table
		routeNRef := &secapi.NetworkReference{
			Tenant:    secapi.TenantID(route.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(route.Metadata.Workspace),
			Network:   secapi.NetworkID(route.Metadata.Network),
			Name:      route.Metadata.Name,
		}
		stepsBuilder.GetRouteTableWithErrorV1Step("Get deleted route table", suite.Client.NetworkV1, *routeNRef, secapi.ErrResourceNotFound)
	}

	// Delete all internet gateways
	for _, gateway := range gateways {
		stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, &gateway)

		// Get the deleted internet gateway
		gatewayWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(gateway.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(gateway.Metadata.Workspace),
			Name:      gateway.Metadata.Name,
		}
		stepsBuilder.GetInternetGatewayWithErrorV1Step("Get deleted internet gateway", suite.Client.NetworkV1, *gatewayWRef, secapi.ErrResourceNotFound)
	}

	// Delete all networks
	for _, network := range networks {
		stepsBuilder.DeleteNetworkV1Step("Delete the network", suite.Client.NetworkV1, &network)

		// Get the deleted network
		networkWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(network.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
			Name:      network.Metadata.Name,
		}
		stepsBuilder.GetNetworkWithErrorV1Step("Get deleted network", suite.Client.NetworkV1, *networkWRef, secapi.ErrResourceNotFound)
	}

	// Delete the workspace
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)

	// Get the deleted workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *ProviderQueriesV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
