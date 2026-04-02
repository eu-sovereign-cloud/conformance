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
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
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

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)

	instanceSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, instanceSkuName)
	instanceName := generators.GenerateInstanceName()

	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName1)
	networkName := generators.GenerateNetworkName()
	networkName2 := generators.GenerateNetworkName()

	internetGatewayName := generators.GenerateInternetGatewayName()
	internetGatewayName2 := generators.GenerateInternetGatewayName()
	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)

	routeTableName := generators.GenerateRouteTableName()
	routeTableName2 := generators.GenerateRouteTableName()
	routeTableRefObj := generators.GenerateRouteTableRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, routeTableName)

	subnetName := generators.GenerateSubnetName()
	subnetName2 := generators.GenerateSubnetName()
	subnetRefObj := generators.GenerateSubnetRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, subnetName)

	nicName := generators.GenerateNicName()
	nicName2 := generators.GenerateNicName()

	publicIpName := generators.GeneratePublicIpName()
	publicIpName2 := generators.GeneratePublicIpName()
	publicIpRefObj := generators.GeneratePublicIpRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, publicIpName)

	securityGroupRuleName := generators.GenerateSecurityGroupRuleName()
	securityGroupRuleName2 := generators.GenerateSecurityGroupRuleName()

	securityGroupName := generators.GenerateSecurityGroupName()
	securityGroupName2 := generators.GenerateSecurityGroupName()

	blockStorageSize := constants.BlockStorageInitialSize

	// Workspace
	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: suite.config.NetworkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	network2, err := builders.NewNetworkBuilder().
		Name(networkName2).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: suite.config.NetworkCidr},
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
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: false,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	internetGateway2, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName2).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: false,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	internetGateways := []schema.InternetGateway{*internetGateway, *internetGateway2}

	routeTable, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: subnetCidr},
			Zone: zone,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	subnet2, err := builders.NewSubnetBuilder().
		Name(subnetName2).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: subnetCidr},
			Zone: zone,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	subnets := []schema.Subnet{*subnet, *subnet2}

	nic, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: []schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	nic2, err := builders.NewNicBuilder().
		Name(nicName2).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: []schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	nics := []schema.Nic{*nic, *nic2}

	publicIp, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: publicIpAddress1,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	publicIp2, err := builders.NewPublicIpBuilder().
		Name(publicIpName2).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: publicIpAddress1,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	publicIps := []schema.PublicIp{*publicIp, *publicIp2}

	securityGroupRule, err := builders.NewSecurityGroupRuleBuilder().
		Name(securityGroupRuleName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group Rule: %v", err)
	}

	securityGroupRule2, err := builders.NewSecurityGroupRuleBuilder().
		Name(securityGroupRuleName2).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group Rule: %v", err)
	}

	securityGroupRules := []schema.SecurityGroupRule{*securityGroupRule, *securityGroupRule2}

	securityGroup, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Workspace:          workspace,
		BlockStorage:       blockStorage,
		Instance:           instance,
		Networks:           networks,
		InternetGateways:   internetGateways,
		RouteTables:        routeTables,
		Subnets:            subnets,
		Nics:               nics,
		PublicIps:          publicIps,
		SecurityGroupRules: securityGroupRules,
		SecurityGroups:     securityGroups,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureProviderQueriesV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ProviderQueriesV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.NetworkProviderV1Name,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroupRule),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace
	workspace := suite.params.Workspace

	// Create a workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.StepResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Network
	networks := suite.params.Networks

	// Create networks
	for _, network := range networks {
		expectNetworkMeta := network.Metadata
		expectNetworkSpec := &network.Spec
		stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", suite.Client.NetworkV1, &network,
			steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
				Metadata:       expectNetworkMeta,
				Spec:           expectNetworkSpec,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)
	}

	wpath := secapi.WorkspacePath{
		Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
	}

	// List networks
	stepsBuilder.ListNetworkV1Step("List Network", suite.Client.NetworkV1, wpath, nil)

	// List networks with limit
	stepsBuilder.ListNetworkV1Step("List Network with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List networks with label
	stepsBuilder.ListNetworkV1Step("List Network with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List networks with limit and label
	stepsBuilder.ListNetworkV1Step("List Network with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Skus

	// List skus
	stepsBuilder.ListNetworkSkusV1Step("List skus", suite.Client.NetworkV1, secapi.TenantPath{Tenant: secapi.TenantID(workspace.Metadata.Tenant)}, nil)

	// List skus with limit
	stepsBuilder.ListNetworkSkusV1Step("List skus", suite.Client.NetworkV1, secapi.TenantPath{Tenant: secapi.TenantID(workspace.Metadata.Tenant)},
		secapi.NewListOptions().WithLimit(1))

	// Internet gateway
	gateways := suite.params.InternetGateways

	// Create internet gateways
	for _, gateway := range gateways {
		expectGatewayMeta := gateway.Metadata
		expectGatewaySpec := &gateway.Spec
		stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create a internet gateway", suite.Client.NetworkV1, &gateway,
			steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
				Metadata:       expectGatewayMeta,
				Spec:           expectGatewaySpec,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)

	}

	// List internet gateways
	stepsBuilder.ListInternetGatewayV1Step("List Internet Gateway", suite.Client.NetworkV1, wpath, nil)

	// List internet gateways with limit
	stepsBuilder.ListInternetGatewayV1Step("List Internet Gateway with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List internet gateways with label
	stepsBuilder.ListInternetGatewayV1Step("List Internet Gateway with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List internet gateways with limit and label
	stepsBuilder.ListInternetGatewayV1Step("List Internet Gateway with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Route table
	routes := suite.params.RouteTables

	// Create route tables
	for _, route := range routes {
		expectRouteMeta := route.Metadata
		expectRouteSpec := &route.Spec
		stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", suite.Client.NetworkV1, &route,
			steps.StepResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
				Metadata:       expectRouteMeta,
				Spec:           expectRouteSpec,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)
	}

	npath := secapi.NetworkPath{
		Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		Workspace: secapi.WorkspaceID((workspace.Metadata.Name)),
		Network:   secapi.NetworkID(networks[0].Metadata.Name),
	}

	// List route tables
	stepsBuilder.ListRouteTableV1Step("List Route table", suite.Client.NetworkV1, npath, nil)

	// List route tables with limit
	stepsBuilder.ListRouteTableV1Step("List Route table with limit", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLimit(1))

	// List route tables with label
	stepsBuilder.ListRouteTableV1Step("List Route table with label", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List route tables with limit and label
	stepsBuilder.ListRouteTableV1Step("List Route table with limit and label", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Subnet
	subnets := suite.params.Subnets

	// Create subnets
	for _, subnet := range subnets {
		expectSubnetMeta := subnet.Metadata
		expectSubnetSpec := &subnet.Spec
		stepsBuilder.CreateOrUpdateSubnetV1Step("Create a subnet", suite.Client.NetworkV1, &subnet,
			steps.StepResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
				Metadata:       expectSubnetMeta,
				Spec:           expectSubnetSpec,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)
	}

	// List subnets
	stepsBuilder.ListSubnetV1Step("List Subnet", suite.Client.NetworkV1, npath, nil)

	// List subnets with limit
	stepsBuilder.ListSubnetV1Step("List Subnet with limit", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLimit(1))

	// List subnets with label
	stepsBuilder.ListSubnetV1Step("List Subnet with label", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List subnets with limit and label
	stepsBuilder.ListSubnetV1Step("List Subnet with limit and label", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Public ip
	publicIps := suite.params.PublicIps

	// Create public ips
	for _, publicIp := range publicIps {
		expectPublicIpMeta := publicIp.Metadata
		expectPublicIpSpec := &publicIp.Spec
		stepsBuilder.CreateOrUpdatePublicIpV1Step("Create a public ip", suite.Client.NetworkV1, &publicIp,
			steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
				Metadata:       expectPublicIpMeta,
				Spec:           expectPublicIpSpec,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)
	}

	// List public ips
	stepsBuilder.ListPublicIpV1Step("List PublicIP", suite.Client.NetworkV1, wpath, nil)

	// List public ips with limit
	stepsBuilder.ListPublicIpV1Step("List PublicIP with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List public ips with label
	stepsBuilder.ListPublicIpV1Step("List PublicIP with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List public ips with limit and label
	stepsBuilder.ListPublicIpV1Step("List PublicIP with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Nic
	nics := suite.params.Nics

	// Create nics
	for _, nic := range nics {
		expectNicMeta := nic.Metadata
		expectNicSpec := &nic.Spec
		stepsBuilder.CreateOrUpdateNicV1Step("Create a nic", suite.Client.NetworkV1, &nic,
			steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
				Metadata:       expectNicMeta,
				Spec:           expectNicSpec,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)
	}

	// List nics
	stepsBuilder.ListNicV1Step("List Nic", suite.Client.NetworkV1, wpath, nil)

	// List nics with limit
	stepsBuilder.ListNicV1Step("List Nic with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List nics with label
	stepsBuilder.ListNicV1Step("List Nic with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List nics with limit and label
	stepsBuilder.ListNicV1Step("List Nic with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Security Group Rule
	rules := suite.params.SecurityGroupRules

	// Create security group rules
	for _, rule := range rules {
		expectRuleMeta := rule.Metadata
		expectRuleSpec := &rule.Spec
		stepsBuilder.CreateOrUpdateSecurityGroupRuleV1Step("Create a security group rule", suite.Client.NetworkV1, &rule,
			steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
				Metadata:       expectRuleMeta,
				Spec:           expectRuleSpec,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)
	}

	// List security group rules
	stepsBuilder.ListSecurityGroupRuleV1Step("List Security Group Rule", suite.Client.NetworkV1, wpath, nil)

	// List security group rules with limit
	stepsBuilder.ListSecurityGroupRuleV1Step("List Security Group Rule with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List security group rules with label
	stepsBuilder.ListSecurityGroupRuleV1Step("List Security Group Rule with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List security group rules with limit and label
	stepsBuilder.ListSecurityGroupRuleV1Step("List Security Group Rule with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Security Group
	groups := suite.params.SecurityGroups

	// Create security groups
	for _, group := range groups {
		expectGroupMeta := group.Metadata
		expectGroupSpec := &group.Spec
		stepsBuilder.CreateOrUpdateSecurityGroupV1Step("Create a security group", suite.Client.NetworkV1, &group,
			steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
				Metadata:       expectGroupMeta,
				Spec:           expectGroupSpec,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)
	}

	// List security groups
	stepsBuilder.ListSecurityGroupV1Step("List Security Group", suite.Client.NetworkV1, wpath, nil)

	// List security groups with limit
	stepsBuilder.ListSecurityGroupV1Step("List Security Group with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List security groups with label
	stepsBuilder.ListSecurityGroupV1Step("List Security Group with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List security groups with limit and label
	stepsBuilder.ListSecurityGroupV1Step("List Security Group with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Delete all security group rules
	for _, rule := range rules {
		stepsBuilder.DeleteSecurityGroupRuleV1Step("Delete the security group rule", suite.Client.NetworkV1, &rule)

		ruleWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(rule.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(rule.Metadata.Workspace),
			Name:      rule.Metadata.Name,
		}
		stepsBuilder.WatchSecurityGroupRuleUntilDeletedV1Step("Watch the security group rule deletion", suite.Client.NetworkV1, ruleWRef)
	}

	// Delete all security groups
	for _, group := range groups {
		stepsBuilder.DeleteSecurityGroupV1Step("Delete the security group", suite.Client.NetworkV1, &group)

		// Get deleted security group
		groupWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(group.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(group.Metadata.Workspace),
			Name:      group.Metadata.Name,
		}
		stepsBuilder.WatchSecurityGroupUntilDeletedV1Step("Watch the security group deletion", suite.Client.NetworkV1, groupWRef)
	}

	// Delete all nics
	for _, nic := range nics {
		stepsBuilder.DeleteNicV1Step("Delete the nic", suite.Client.NetworkV1, &nic)

		// Get the deleted nic
		nicWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(nic.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(nic.Metadata.Workspace),
			Name:      nic.Metadata.Name,
		}
		stepsBuilder.WatchNicUntilDeletedV1Step("Watch the nic deletion", suite.Client.NetworkV1, nicWRef)
	}

	// Delete all public ips
	for _, publicIp := range publicIps {
		stepsBuilder.DeletePublicIpV1Step("Delete the public ip", suite.Client.NetworkV1, &publicIp)

		// Get the deleted public ip
		publicIpWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(publicIp.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(publicIp.Metadata.Workspace),
			Name:      publicIp.Metadata.Name,
		}
		stepsBuilder.WatchPublicIpUntilDeletedV1Step("Watch the public ip deletion", suite.Client.NetworkV1, publicIpWRef)
	}

	// Delete all subnets
	for _, subnet := range subnets {
		stepsBuilder.DeleteSubnetV1Step("Delete the subnet", suite.Client.NetworkV1, &subnet)

		// Get the deleted subnet
		subnetNRef := secapi.NetworkReference{
			Tenant:    secapi.TenantID(subnet.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(subnet.Metadata.Workspace),
			Network:   secapi.NetworkID(subnet.Metadata.Network),
			Name:      subnet.Metadata.Name,
		}
		stepsBuilder.WatchSubnetUntilDeletedV1Step("Watch the subnet deletion", suite.Client.NetworkV1, subnetNRef)
	}

	// Delete all route tables
	for _, route := range routes {
		stepsBuilder.DeleteRouteTableV1Step("Delete the route table", suite.Client.NetworkV1, &route)

		// Get the deleted route table
		routeNRef := secapi.NetworkReference{
			Tenant:    secapi.TenantID(route.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(route.Metadata.Workspace),
			Network:   secapi.NetworkID(route.Metadata.Network),
			Name:      route.Metadata.Name,
		}
		stepsBuilder.WatchRouteTableUntilDeletedV1Step("Watch the route table deletion", suite.Client.NetworkV1, routeNRef)
	}

	// Delete all internet gateways
	for _, gateway := range gateways {
		stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, &gateway)

		// Get the deleted internet gateway
		internetGatWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(gateway.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(gateway.Metadata.Workspace),
			Name:      gateway.Metadata.Name,
		}
		stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", suite.Client.NetworkV1, internetGatWRef)
	}

	// Delete all networks
	for _, network := range networks {
		stepsBuilder.DeleteNetworkV1Step("Delete the network", suite.Client.NetworkV1, &network)

		// Get the deleted network
		networkWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(network.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
			Name:      network.Metadata.Name,
		}
		stepsBuilder.WatchNetworkUntilDeletedV1Step("Watch the network deletion", suite.Client.NetworkV1, networkWRef)
	}

	// Delete the workspace
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)

	// Get the deleted workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *ProviderQueriesV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
