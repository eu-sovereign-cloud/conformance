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
	t.AddParentSuite(suites.NetworkParentSuite)

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
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetworkSku),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroupRule),
	)
	suite.ConfigureDepends(t, string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalResourceMetadataKindResourceKindInstance),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace
	workspace := suite.params.Workspace

	// Create a workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Network
	networks := suite.params.Networks

	// Create networks
	steps.BulkCreateNetworksStepsV1(stepsBuilder, suite.RegionalTestSuite, "Create networks", networks)

	wpath := secapi.WorkspacePath{
		Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
	}

	// List networks
	stepsBuilder.ListNetworkV1Step("List networks", suite.Client.NetworkV1, wpath, nil)

	// List networks with limit
	stepsBuilder.ListNetworkV1Step("List networks with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List networks with label
	stepsBuilder.ListNetworkV1Step("List networks with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List networks with limit and label
	stepsBuilder.ListNetworkV1Step("List networks with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Skus

	// List skus
	stepsBuilder.ListNetworkSkusV1Step("List skus", suite.Client.NetworkV1, secapi.TenantPath{Tenant: secapi.TenantID(workspace.Metadata.Tenant)}, nil)

	// List skus with limit
	stepsBuilder.ListNetworkSkusV1Step("List skus with limit", suite.Client.NetworkV1, secapi.TenantPath{Tenant: secapi.TenantID(workspace.Metadata.Tenant)},
		secapi.NewListOptions().WithLimit(1))

	// Internet gateway
	gateways := suite.params.InternetGateways

	// Create internet gateways
	steps.BulkCreateInternetGatewaysStepsV1(stepsBuilder, suite.RegionalTestSuite, "Create internet gateways", gateways)

	// List internet gateways
	stepsBuilder.ListInternetGatewayV1Step("List internet gateways", suite.Client.NetworkV1, wpath, nil)

	// List internet gateways with limit
	stepsBuilder.ListInternetGatewayV1Step("List internet gateways with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List internet gateways with label
	stepsBuilder.ListInternetGatewayV1Step("List internet gateways with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List internet gateways with limit and label
	stepsBuilder.ListInternetGatewayV1Step("List internet gateways with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Route table
	routes := suite.params.RouteTables

	// Create route tables
	steps.BulkCreateRouteTablesStepsV1(stepsBuilder, suite.RegionalTestSuite, "Create route tables", routes)

	npath := secapi.NetworkPath{
		Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		Workspace: secapi.WorkspaceID((workspace.Metadata.Name)),
		Network:   secapi.NetworkID(networks[0].Metadata.Name),
	}

	// List route tables
	stepsBuilder.ListRouteTableV1Step("List route tables", suite.Client.NetworkV1, npath, nil)

	// List route tables with limit
	stepsBuilder.ListRouteTableV1Step("List route tables with limit", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLimit(1))

	// List route tables with label
	stepsBuilder.ListRouteTableV1Step("List route tables with label", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List route tables with limit and label
	stepsBuilder.ListRouteTableV1Step("List route tables with limit and label", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Subnet
	subnets := suite.params.Subnets

	// Create subnets
	steps.BulkCreateSubnetsStepsV1(stepsBuilder, suite.RegionalTestSuite, "Create subnets", subnets)

	// List subnets
	stepsBuilder.ListSubnetV1Step("List subnets", suite.Client.NetworkV1, npath, nil)

	// List subnets with limit
	stepsBuilder.ListSubnetV1Step("List subnets with limit", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLimit(1))

	// List subnets with label
	stepsBuilder.ListSubnetV1Step("List subnets with label", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List subnets with limit and label
	stepsBuilder.ListSubnetV1Step("List subnets with limit and label", suite.Client.NetworkV1, npath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Public ip
	publicIps := suite.params.PublicIps

	// Create public ips
	steps.BulkCreatePublicIpsStepsV1(stepsBuilder, suite.RegionalTestSuite, "Create public ips", publicIps)

	// List public ips
	stepsBuilder.ListPublicIpV1Step("List public ips", suite.Client.NetworkV1, wpath, nil)

	// List public ips with limit
	stepsBuilder.ListPublicIpV1Step("List public ips with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List public ips with label
	stepsBuilder.ListPublicIpV1Step("List public ips with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List public ips with limit and label
	stepsBuilder.ListPublicIpV1Step("List public ips with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Nic
	nics := suite.params.Nics

	// Create nics
	steps.BulkCreateNicsStepsV1(stepsBuilder, suite.RegionalTestSuite, "Create nics", nics)

	// List nics
	stepsBuilder.ListNicV1Step("List nics", suite.Client.NetworkV1, wpath, nil)

	// List nics with limit
	stepsBuilder.ListNicV1Step("List nics with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List nics with label
	stepsBuilder.ListNicV1Step("List nics with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List nics with limit and label
	stepsBuilder.ListNicV1Step("List nics with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Security Group Rule
	rules := suite.params.SecurityGroupRules

	// Create security group rules
	steps.BulkCreateSecurityGroupRulesStepsV1(stepsBuilder, suite.RegionalTestSuite, "Create security group rules", rules)

	// List security group rules
	stepsBuilder.ListSecurityGroupRuleV1Step("List security group rules", suite.Client.NetworkV1, wpath, nil)

	// List security group rules with limit
	stepsBuilder.ListSecurityGroupRuleV1Step("List security group rules with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List security group rules with label
	stepsBuilder.ListSecurityGroupRuleV1Step("List security group rules with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List security group rules with limit and label
	stepsBuilder.ListSecurityGroupRuleV1Step("List security group rules with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Security Group
	groups := suite.params.SecurityGroups

	// Create security groups
	steps.BulkCreateSecurityGroupsStepsV1(stepsBuilder, suite.RegionalTestSuite, "Create security groups", groups)

	// List security groups
	stepsBuilder.ListSecurityGroupV1Step("List security groups", suite.Client.NetworkV1, wpath, nil)

	// List security groups with limit
	stepsBuilder.ListSecurityGroupV1Step("List security groups with limit", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List security groups with label
	stepsBuilder.ListSecurityGroupV1Step("List security groups with label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List security groups with limit and label
	stepsBuilder.ListSecurityGroupV1Step("List security groups with limit and label", suite.Client.NetworkV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Delete all security group rules
	steps.BulkDeleteSecurityGroupRulesStepsV1(stepsBuilder, suite.RegionalTestSuite, "Delete all security group rules", rules)

	// Delete all security groups
	steps.BulkDeleteSecurityGroupsStepsV1(stepsBuilder, suite.RegionalTestSuite, "Delete all security groups", groups)

	// Delete all nics
	steps.BulkDeleteNicsStepsV1(stepsBuilder, suite.RegionalTestSuite, "Delete all nics", nics)

	// Delete all public ips
	steps.BulkDeletePublicIpsStepsV1(stepsBuilder, suite.RegionalTestSuite, "Delete all public ips", publicIps)

	// Delete all subnets
	steps.BulkDeleteSubnetsStepsV1(stepsBuilder, suite.RegionalTestSuite, "Delete all subnets", subnets)

	// Delete all route tables
	steps.BulkDeleteRouteTablesStepsV1(stepsBuilder, suite.RegionalTestSuite, "Delete all route tables", routes)

	// Delete all internet gateways
	steps.BulkDeleteInternetGatewaysStepsV1(stepsBuilder, suite.RegionalTestSuite, "Delete all internet gateways", gateways)

	// Delete all networks
	steps.BulkDeleteNetworksStepsV1(stepsBuilder, suite.RegionalTestSuite, "Delete all networks", networks)

	// Delete the workspace
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)

	// Get the deleted workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *ProviderQueriesV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
