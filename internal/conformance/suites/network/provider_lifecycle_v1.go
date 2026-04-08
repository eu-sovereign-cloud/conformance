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

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ProviderLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	config *ProviderLifeCycleV1Config
	params *params.NetworkProviderLifeCycleV1Params
}

type ProviderLifeCycleV1Config struct {
	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	StorageSkus    []string
	InstanceSkus   []string
	NetworkSkus    []string
}

func CreateProviderLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *ProviderLifeCycleV1Config) *ProviderLifeCycleV1TestSuite {
	suite := &ProviderLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.NetworkProviderLifeCycleV1SuiteName.String()
	return suite
}

func (suite *ProviderLifeCycleV1TestSuite) BeforeAll(t provider.T) {
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
	nicAddress2, err := generators.GenerateNicAddress(subnetCidr, 2)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}

	// Generate the public ips
	publicIpAddress1, err := generators.GeneratePublicIp(suite.config.PublicIpsRange, 1)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}
	publicIpAddress2, err := generators.GeneratePublicIp(suite.config.PublicIpsRange, 2)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}

	// Select zones
	zone1 := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]
	zone2 := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]

	// Select skus
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	networkSkuName1 := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]
	networkSkuName2 := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)

	instanceSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, instanceSkuName)
	instanceName := generators.GenerateInstanceName()
	instanceRefObj := generators.GenerateInstanceRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, workspaceName, instanceName)

	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName1)
	networkSkuRef2Obj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName2)

	networkName := generators.GenerateNetworkName()

	internetGatewayName := generators.GenerateInternetGatewayName()
	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)

	routeTableName := generators.GenerateRouteTableName()
	routeTableRefObj := generators.GenerateRouteTableRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, routeTableName)

	subnetName := generators.GenerateSubnetName()
	subnetRefObj := generators.GenerateSubnetRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, subnetName)

	nicName := generators.GenerateNicName()

	publicIpName := generators.GeneratePublicIpName()
	publicIpRefObj := generators.GeneratePublicIpRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, publicIpName)

	securityGroupRuleName := generators.GenerateSecurityGroupRuleName()
	securityGroupRuleRefObj := generators.GenerateSecurityGroupRuleRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, securityGroupRuleName)

	securityGroupName := generators.GenerateSecurityGroupName()

	blockStorageSize := constants.BlockStorageInitialSize

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
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
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone1,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		},
		).Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	networkInitial, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: suite.config.NetworkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	networkUpdated, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: suite.config.NetworkCidr},
			SkuRef:        *networkSkuRef2Obj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	internetGatInitial, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: false,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	internetGatUpdated, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: true,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	routeTableInitial, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Route Table: %v", err)
	}

	routeTableUpdated, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *instanceRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Route Table: %v", err)
	}

	subnetInitial, err := builders.NewSubnetBuilder().
		Name(subnetName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: subnetCidr},
			Zone: zone1,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	subnetUpdated, err := builders.NewSubnetBuilder().
		Name(subnetName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: subnetCidr},
			Zone: zone2,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	nicInitial, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: []schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	nicUpdated, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NicSpec{
			Addresses:    []string{nicAddress2},
			PublicIpRefs: []schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	publicIpInitial, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: publicIpAddress1,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	publicIpUpdated, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: publicIpAddress2,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	securityGroupRuleInitial, err := builders.NewSecurityGroupRuleBuilder().
		Name(securityGroupRuleName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group Rule: %v", err)
	}

	securityGroupRuleUpdated, err := builders.NewSecurityGroupRuleBuilder().
		Name(securityGroupRuleName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionEgress}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group Rule: %v", err)
	}

	securityGroupInitial, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.SecurityGroupSpec{
			RuleRefs: []schema.Reference{*securityGroupRuleRefObj},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	securityGroupUpdated, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	params := &params.NetworkProviderLifeCycleV1Params{
		Workspace:                workspace,
		BlockStorage:             blockStorage,
		Instance:                 instance,
		NetworkInitial:           networkInitial,
		NetworkUpdated:           networkUpdated,
		InternetGatewayInitial:   internetGatInitial,
		InternetGatewayUpdated:   internetGatUpdated,
		RouteTableInitial:        routeTableInitial,
		RouteTableUpdated:        routeTableUpdated,
		SubnetInitial:            subnetInitial,
		SubnetUpdated:            subnetUpdated,
		NicInitial:               nicInitial,
		NicUpdated:               nicUpdated,
		PublicIpInitial:          publicIpInitial,
		PublicIpUpdated:          publicIpUpdated,
		SecurityGroupInitial:     securityGroupInitial,
		SecurityGroupUpdated:     securityGroupUpdated,
		SecurityGroupRuleInitial: securityGroupRuleInitial,
		SecurityGroupRuleUpdated: securityGroupRuleUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureProviderLifecycleScenarioV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ProviderLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
	)
	suite.ConfigureDepends(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalResourceMetadataKindResourceKindInstance),
	)

	stepsConfigurator := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsConfigurator.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsConfigurator.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Network

	// Create a network
	network := suite.params.NetworkInitial
	expectNetworkMeta := network.Metadata
	expectNetworkSpec := &network.Spec
	stepsConfigurator.CreateOrUpdateNetworkV1Step("Create a network", t, suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	networkWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(network.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
		Name:      network.Metadata.Name,
	}

	// Internet gateway

	// Create an internet gateway
	gateway := suite.params.InternetGatewayInitial
	expectGatewayMeta := gateway.Metadata
	expectGatewaySpec := &gateway.Spec
	stepsConfigurator.CreateOrUpdateInternetGatewayV1Step("Create a internet gateway", t, suite.Client.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:       expectGatewayMeta,
			Spec:           expectGatewaySpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	gatewayWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(gateway.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(gateway.Metadata.Workspace),
		Name:      gateway.Metadata.Name,
	}

	// Get the created internet gateway
	stepsConfigurator.GetInternetGatewayV1Step("Get the created internet gateway", suite.Client.NetworkV1, gatewayWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:       expectGatewayMeta,
			Spec:           expectGatewaySpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the internet gateway
	gateway = suite.params.InternetGatewayUpdated
	expectGatewaySpec.EgressOnly = gateway.Spec.EgressOnly
	stepsConfigurator.CreateOrUpdateInternetGatewayV1Step("Update the internet gateway", t, suite.Client.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:       expectGatewayMeta,
			Spec:           expectGatewaySpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated internet gateway
	stepsConfigurator.GetInternetGatewayV1Step("Get the updated internet gateway", suite.Client.NetworkV1, gatewayWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:       expectGatewayMeta,
			Spec:           expectGatewaySpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Route table

	// Create a route table
	route := suite.params.RouteTableInitial
	expectRouteMeta := route.Metadata
	expectRouteSpec := &route.Spec
	stepsConfigurator.CreateOrUpdateRouteTableV1Step("Create a route table", t, suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	routeNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(route.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(route.Metadata.Workspace),
		Network:   secapi.NetworkID(route.Metadata.Network),
		Name:      route.Metadata.Name,
	}

	// Get the created route table
	stepsConfigurator.GetRouteTableV1Step("Get the created route table", suite.Client.NetworkV1, routeNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the created network
	stepsConfigurator.GetNetworkV1Step("Get the created network", suite.Client.NetworkV1, networkWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the route table
	route = suite.params.RouteTableUpdated
	expectRouteSpec.Routes = route.Spec.Routes
	stepsConfigurator.CreateOrUpdateRouteTableV1Step("Update the route table", t, suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated route table
	stepsConfigurator.GetRouteTableV1Step("Get the updated route table", suite.Client.NetworkV1, routeNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Network

	// Update the network
	network = suite.params.NetworkUpdated
	expectNetworkSpec.SkuRef = network.Spec.SkuRef
	stepsConfigurator.CreateOrUpdateNetworkV1Step("Update the network", t, suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated network
	stepsConfigurator.GetNetworkV1Step("Get the updated network", suite.Client.NetworkV1, networkWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Subnet

	// Create a subnet
	subnet := suite.params.SubnetInitial
	expectSubnetMeta := subnet.Metadata
	expectSubnetSpec := &subnet.Spec
	stepsConfigurator.CreateOrUpdateSubnetV1Step("Create a subnet", t, suite.Client.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:       expectSubnetMeta,
			Spec:           expectSubnetSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	subnetNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(subnet.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(subnet.Metadata.Workspace),
		Network:   secapi.NetworkID(subnet.Metadata.Network),
		Name:      subnet.Metadata.Name,
	}

	// Get the created subnet
	stepsConfigurator.GetSubnetV1Step("Get the created subnet", suite.Client.NetworkV1, subnetNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:       expectSubnetMeta,
			Spec:           expectSubnetSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the subnet
	subnet = suite.params.SubnetUpdated
	expectSubnetSpec.Zone = subnet.Spec.Zone
	stepsConfigurator.CreateOrUpdateSubnetV1Step("Update the subnet", t, suite.Client.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:       expectSubnetMeta,
			Spec:           expectSubnetSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated subnet
	stepsConfigurator.GetSubnetV1Step("Get the updated subnet", suite.Client.NetworkV1, subnetNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:       expectSubnetMeta,
			Spec:           expectSubnetSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)
	// Nic

	// Create a nic
	nic := suite.params.NicInitial
	expectNicMeta := nic.Metadata
	expectNicSpec := &nic.Spec
	stepsConfigurator.CreateOrUpdateNicV1Step("Create a nic", t, suite.Client.NetworkV1, nic,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:       expectNicMeta,
			Spec:           expectNicSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	nicWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(nic.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(nic.Metadata.Workspace),
		Name:      nic.Metadata.Name,
	}

	// Get the created nic
	stepsConfigurator.GetNicV1Step("Get the created nic", suite.Client.NetworkV1, nicWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:       expectNicMeta,
			Spec:           expectNicSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the nic
	nic = suite.params.NicUpdated
	expectNicSpec.Addresses = nic.Spec.Addresses
	stepsConfigurator.CreateOrUpdateNicV1Step("Update the nic", t, suite.Client.NetworkV1, nic,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:       expectNicMeta,
			Spec:           expectNicSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated nic
	stepsConfigurator.GetNicV1Step("Get the updated nic", suite.Client.NetworkV1, nicWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:       expectNicMeta,
			Spec:           expectNicSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)
	// Public ip

	// Create a public ip
	publicIp := suite.params.PublicIpInitial
	expectPublicIpMeta := publicIp.Metadata
	expectPublicIpSpec := &publicIp.Spec
	stepsConfigurator.CreateOrUpdatePublicIpV1Step("Create a public ip", t, suite.Client.NetworkV1, publicIp,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:       expectPublicIpMeta,
			Spec:           expectPublicIpSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created public ip
	publicIpWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(publicIp.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(publicIp.Metadata.Workspace),
		Name:      publicIp.Metadata.Name,
	}
	stepsConfigurator.GetPublicIpV1Step("Get the created public ip", suite.Client.NetworkV1, publicIpWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:       expectPublicIpMeta,
			Spec:           expectPublicIpSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the public ip
	publicIp = suite.params.PublicIpUpdated
	expectPublicIpSpec.Address = publicIp.Spec.Address
	stepsConfigurator.CreateOrUpdatePublicIpV1Step("Update the public ip", t, suite.Client.NetworkV1, publicIp,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:       expectPublicIpMeta,
			Spec:           expectPublicIpSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated public ip
	stepsConfigurator.GetPublicIpV1Step("Get the updated public ip", suite.Client.NetworkV1, publicIpWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:       expectPublicIpMeta,
			Spec:           expectPublicIpSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Security Group Rule

	// Create a security group rule
	rule := suite.params.SecurityGroupRuleInitial
	expectRuleMeta := rule.Metadata
	expectRuleSpec := &rule.Spec
	stepsConfigurator.CreateOrUpdateSecurityGroupRuleV1Step("Create a security group rule", t, suite.Client.NetworkV1, rule,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created security group rule
	ruleWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(rule.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(rule.Metadata.Workspace),
		Name:      rule.Metadata.Name,
	}
	stepsConfigurator.GetSecurityGroupRuleV1Step("Get the created security group rule", suite.Client.NetworkV1, ruleWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the security group rule
	rule = suite.params.SecurityGroupRuleUpdated
	expectRuleSpec.Direction = rule.Spec.Direction
	stepsConfigurator.CreateOrUpdateSecurityGroupRuleV1Step("Update the security group rule", t, suite.Client.NetworkV1, rule,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated security group rule
	stepsConfigurator.GetSecurityGroupRuleV1Step("Get the updated security group rule", suite.Client.NetworkV1, ruleWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Security Group

	// Create a security group
	group := suite.params.SecurityGroupInitial
	expectGroupMeta := group.Metadata
	expectGroupSpec := &group.Spec
	stepsConfigurator.CreateOrUpdateSecurityGroupV1Step("Create a security group", t, suite.Client.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created security group
	groupWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(group.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(group.Metadata.Workspace),
		Name:      group.Metadata.Name,
	}
	stepsConfigurator.GetSecurityGroupV1Step("Get the created security group", suite.Client.NetworkV1, groupWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the security group

	group.Spec.Rules = []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}}
	expectGroupSpec = &group.Spec

	stepsConfigurator.CreateOrUpdateSecurityGroupV1Step("Update the security group", t, suite.Client.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated security group
	stepsConfigurator.GetSecurityGroupV1Step("Get the updated security group", suite.Client.NetworkV1, groupWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Block storage

	// Create a block storage
	block := suite.params.BlockStorage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	stepsConfigurator.CreateOrUpdateBlockStorageV1Step("Create a block storage", t, suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Spec:           expectedBlockSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created block storage
	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(block.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
		Name:      block.Metadata.Name,
	}
	stepsConfigurator.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1, blockWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Spec:           expectedBlockSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Instance

	// Create an instance
	instance := suite.params.Instance
	expectInstanceMeta := instance.Metadata
	expectInstanceSpec := &instance.Spec
	stepsConfigurator.CreateOrUpdateInstanceV1Step("Create an instance", t, suite.Client.ComputeV1, instance,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:       expectInstanceMeta,
			Spec:           expectInstanceSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created instance
	instanceWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(instance.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(instance.Metadata.Workspace),
		Name:      instance.Metadata.Name,
	}
	instance = stepsConfigurator.GetInstanceV1Step("Get the created instance", suite.Client.ComputeV1, instanceWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:       expectInstanceMeta,
			Spec:           expectInstanceSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Resources deletion

	stepsConfigurator.DeleteInstanceV1Step("Delete the instance", t, suite.Client.ComputeV1, instance)
	stepsConfigurator.WatchInstanceUntilDeletedV1Step("Watch the instance deletion", t, suite.Client.ComputeV1, instanceWRef)

	stepsConfigurator.DeleteBlockStorageV1Step("Delete the block storage", t, suite.Client.StorageV1, block)
	stepsConfigurator.WatchBlockStorageUntilDeletedV1Step("Watch the block storage deletion", t, suite.Client.StorageV1, blockWRef)

	stepsConfigurator.DeleteSecurityGroupRuleV1Step("Delete the security group rule", t, suite.Client.NetworkV1, rule)
	stepsConfigurator.WatchSecurityGroupRuleUntilDeletedV1Step("Watch the security group rule deletion", t, suite.Client.NetworkV1, ruleWRef)

	stepsConfigurator.DeleteSecurityGroupV1Step("Delete the security group", t, suite.Client.NetworkV1, group)
	stepsConfigurator.WatchSecurityGroupUntilDeletedV1Step("Watch the security group deletion", t, suite.Client.NetworkV1, groupWRef)

	stepsConfigurator.DeleteNicV1Step("Delete the nic", t, suite.Client.NetworkV1, nic)
	stepsConfigurator.WatchNicUntilDeletedV1Step("Watch the nic deletion", t, suite.Client.NetworkV1, nicWRef)

	stepsConfigurator.DeletePublicIpV1Step("Delete the public ip", t, suite.Client.NetworkV1, publicIp)
	stepsConfigurator.WatchPublicIpUntilDeletedV1Step("Watch the public ip deletion", t, suite.Client.NetworkV1, publicIpWRef)

	stepsConfigurator.DeleteSubnetV1Step("Delete the subnet", t, suite.Client.NetworkV1, subnet)
	stepsConfigurator.WatchSubnetUntilDeletedV1Step("Watch the subnet deletion", t, suite.Client.NetworkV1, subnetNRef)

	stepsConfigurator.DeleteRouteTableV1Step("Delete the route table", t, suite.Client.NetworkV1, route)
	stepsConfigurator.WatchRouteTableUntilDeletedV1Step("Watch the route table deletion", t, suite.Client.NetworkV1, routeNRef)

	stepsConfigurator.DeleteInternetGatewayV1Step("Delete the internet gateway", t, suite.Client.NetworkV1, gateway)
	stepsConfigurator.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", t, suite.Client.NetworkV1, gatewayWRef)

	stepsConfigurator.DeleteNetworkV1Step("Delete the network", t, suite.Client.NetworkV1, network)
	stepsConfigurator.WatchNetworkUntilDeletedV1Step("Watch the network deletion", t, suite.Client.NetworkV1, networkWRef)

	stepsConfigurator.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsConfigurator.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *ProviderLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
