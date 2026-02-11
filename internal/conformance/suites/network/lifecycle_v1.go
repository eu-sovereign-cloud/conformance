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
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"k8s.io/utils/ptr"
)

type NetworkLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	config *NetworkLifeCycleV1Config
	params *params.NetworkLifeCycleV1Params
}

type NetworkLifeCycleV1Config struct {
	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	StorageSkus    []string
	InstanceSkus   []string
	NetworkSkus    []string
}

func CreateLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *NetworkLifeCycleV1Config) *NetworkLifeCycleV1TestSuite {
	suite := &NetworkLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.NetworkV1LifeCycleSuiteName
	return suite
}

func (suite *NetworkLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	var err error

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

	storageSkuRefObj, err := generators.GenerateSkuRefObject(storageSkuName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj, err := generators.GenerateBlockStorageRefObject(blockStorageName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	instanceSkuRefObj, err := generators.GenerateSkuRefObject(instanceSkuName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	instanceName := generators.GenerateInstanceName()

	instanceRefObj, err := generators.GenerateInstanceRefObject(instanceName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	networkSkuRefObj, err := generators.GenerateSkuRefObject(networkSkuName1)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	networkSkuRef2Obj, err := generators.GenerateSkuRefObject(networkSkuName2)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	networkName := generators.GenerateNetworkName()

	internetGatewayName := generators.GenerateInternetGatewayName()
	internetGatewayRefObj, err := generators.GenerateInternetGatewayRefObject(internetGatewayName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	routeTableName := generators.GenerateRouteTableName()
	routeTableRefObj, err := generators.GenerateRouteTableRefObject(routeTableName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	subnetName := generators.GenerateSubnetName()
	subnetRefObj, err := generators.GenerateSubnetRefObject(subnetName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	nicName := generators.GenerateNicName()

	publicIpName := generators.GeneratePublicIpName()
	publicIpRefObj, err := generators.GeneratePublicIpRefObject(publicIpName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	securityGroupName := generators.GenerateSecurityGroupName()

	blockStorageSize := generators.GenerateBlockStorageSize()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: ptr.To(suite.config.NetworkCidr)},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	networkUpdated, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: ptr.To(suite.config.NetworkCidr)},
			SkuRef:        *networkSkuRef2Obj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	internetGatInitial, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: ptr.To(false),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	internetGatUpdated, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: ptr.To(true),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	routeTableInitial, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone1,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	subnetUpdated, err := builders.NewSubnetBuilder().
		Name(subnetName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone2,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	nicInitial, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	nicUpdated, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NicSpec{
			Addresses:    []string{nicAddress2},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	publicIpInitial, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: ptr.To(publicIpAddress1),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	publicIpUpdated, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: ptr.To(publicIpAddress2),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	securityGroupInitial, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	securityGroupUpdated, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionEgress}},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	params := &params.NetworkLifeCycleV1Params{
		Workspace:              workspace,
		BlockStorage:           blockStorage,
		Instance:               instance,
		NetworkInitial:         networkInitial,
		NetworkUpdated:         networkUpdated,
		InternetGatewayInitial: internetGatInitial,
		InternetGatewayUpdated: internetGatUpdated,
		RouteTableInitial:      routeTableInitial,
		RouteTableUpdated:      routeTableUpdated,
		SubnetInitial:          subnetInitial,
		SubnetUpdated:          subnetUpdated,
		NicInitial:             nicInitial,
		NicUpdated:             nicUpdated,
		PublicIpInitial:        publicIpInitial,
		PublicIpUpdated:        publicIpUpdated,
		SecurityGroupInitial:   securityGroupInitial,
		SecurityGroupUpdated:   securityGroupUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureLifecycleScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

//nolint:dupl
func (suite *NetworkLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.NetworkProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindInstance),
	)

	// Prepare shared resources and references
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}

	network := suite.params.NetworkInitial
	expectNetworkMeta := network.Metadata
	expectNetworkSpec := &network.Spec
	networkWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(network.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
		Name:      network.Metadata.Name,
	}

	gateway := suite.params.InternetGatewayInitial
	expectGatewayMeta := gateway.Metadata
	expectGatewaySpec := &gateway.Spec
	gatewayWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(gateway.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(gateway.Metadata.Workspace),
		Name:      gateway.Metadata.Name,
	}

	route := suite.params.RouteTableInitial
	expectRouteMeta := route.Metadata
	expectRouteSpec := &route.Spec
	routeNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(route.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(route.Metadata.Workspace),
		Network:   secapi.NetworkID(route.Metadata.Network),
		Name:      route.Metadata.Name,
	}

	subnet := suite.params.SubnetInitial
	expectSubnetMeta := subnet.Metadata
	expectSubnetSpec := &subnet.Spec
	subnetNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(subnet.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(subnet.Metadata.Workspace),
		Network:   secapi.NetworkID(subnet.Metadata.Network),
		Name:      subnet.Metadata.Name,
	}

	publicIp := suite.params.PublicIpInitial
	expectPublicIpMeta := publicIp.Metadata
	expectPublicIpSpec := &publicIp.Spec
	publicIpWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(publicIp.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(publicIp.Metadata.Workspace),
		Name:      publicIp.Metadata.Name,
	}

	nic := suite.params.NicInitial
	expectNicMeta := nic.Metadata
	expectNicSpec := &nic.Spec
	nicWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(nic.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(nic.Metadata.Workspace),
		Name:      nic.Metadata.Name,
	}

	group := suite.params.SecurityGroupInitial
	expectGroupMeta := group.Metadata
	expectGroupSpec := &group.Spec
	groupWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(group.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(group.Metadata.Workspace),
		Name:      group.Metadata.Name,
	}

	block := suite.params.BlockStorage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(block.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
		Name:      block.Metadata.Name,
	}

	instance := suite.params.Instance
	expectInstanceMeta := instance.Metadata
	expectInstanceSpec := &instance.Spec
	instanceWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(instance.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(instance.Metadata.Workspace),
		Name:      instance.Metadata.Name,
	}

	// Workspace lifecycle
	t.WithNewStep("Workspace", func(wsCtx provider.StepCtx) {
		wsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, wsCtx)

		wsSteps.CreateOrUpdateWorkspaceV1Step("Create", suite.Client.WorkspaceV1, workspace,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectWorkspaceLabels,
				Metadata:      expectWorkspaceMeta,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		wsSteps.GetWorkspaceV1Step("Get", suite.Client.WorkspaceV1, workspaceTRef,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectWorkspaceLabels,
				Metadata:      expectWorkspaceMeta,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Network lifecycle
	t.WithNewStep("Network", func(netCtx provider.StepCtx) {
		netSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, netCtx)

		netSteps.CreateOrUpdateNetworkV1Step("Create", suite.Client.NetworkV1, network,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
				Metadata:      expectNetworkMeta,
				Spec:          expectNetworkSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		netSteps.GetNetworkV1Step("Get", suite.Client.NetworkV1, networkWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
				Metadata:      expectNetworkMeta,
				Spec:          expectNetworkSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		network.Spec = suite.params.NetworkUpdated.Spec
		expectNetworkSpec.SkuRef = network.Spec.SkuRef
		netSteps.CreateOrUpdateNetworkV1Step("Update", suite.Client.NetworkV1, network,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
				Metadata:      expectNetworkMeta,
				Spec:          expectNetworkSpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		netSteps.GetNetworkV1Step("GetUpdated", suite.Client.NetworkV1, networkWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
				Metadata:      expectNetworkMeta,
				Spec:          expectNetworkSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Internet gateway lifecycle
	t.WithNewStep("InternetGateway", func(gwCtx provider.StepCtx) {
		gwSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, gwCtx)

		gwSteps.CreateOrUpdateInternetGatewayV1Step("Create", suite.Client.NetworkV1, gateway,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
				Metadata:      expectGatewayMeta,
				Spec:          expectGatewaySpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		gwSteps.GetInternetGatewayV1Step("Get", suite.Client.NetworkV1, gatewayWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
				Metadata:      expectGatewayMeta,
				Spec:          expectGatewaySpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		gateway.Spec = suite.params.InternetGatewayUpdated.Spec
		expectGatewaySpec.EgressOnly = gateway.Spec.EgressOnly
		gwSteps.CreateOrUpdateInternetGatewayV1Step("Update", suite.Client.NetworkV1, gateway,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
				Metadata:      expectGatewayMeta,
				Spec:          expectGatewaySpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		gwSteps.GetInternetGatewayV1Step("GetUpdated", suite.Client.NetworkV1, gatewayWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
				Metadata:      expectGatewayMeta,
				Spec:          expectGatewaySpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Route table lifecycle
	t.WithNewStep("RouteTable", func(rtCtx provider.StepCtx) {
		rtSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, rtCtx)

		rtSteps.CreateOrUpdateRouteTableV1Step("Create", suite.Client.NetworkV1, route,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
				Metadata:      expectRouteMeta,
				Spec:          expectRouteSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		rtSteps.GetRouteTableV1Step("Get", suite.Client.NetworkV1, routeNRef,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
				Metadata:      expectRouteMeta,
				Spec:          expectRouteSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		route.Spec = suite.params.RouteTableUpdated.Spec
		expectRouteSpec.Routes = route.Spec.Routes
		rtSteps.CreateOrUpdateRouteTableV1Step("Update", suite.Client.NetworkV1, route,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
				Metadata:      expectRouteMeta,
				Spec:          expectRouteSpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		rtSteps.GetRouteTableV1Step("GetUpdated", suite.Client.NetworkV1, routeNRef,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
				Metadata:      expectRouteMeta,
				Spec:          expectRouteSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Subnet lifecycle
	t.WithNewStep("Subnet", func(snCtx provider.StepCtx) {
		snSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, snCtx)

		snSteps.CreateOrUpdateSubnetV1Step("Create", suite.Client.NetworkV1, subnet,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
				Metadata:      expectSubnetMeta,
				Spec:          expectSubnetSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		snSteps.GetSubnetV1Step("Get", suite.Client.NetworkV1, subnetNRef,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
				Metadata:      expectSubnetMeta,
				Spec:          expectSubnetSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		subnet.Spec = suite.params.SubnetUpdated.Spec
		expectSubnetSpec.Zone = subnet.Spec.Zone
		snSteps.CreateOrUpdateSubnetV1Step("Update", suite.Client.NetworkV1, subnet,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
				Metadata:      expectSubnetMeta,
				Spec:          expectSubnetSpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		snSteps.GetSubnetV1Step("GetUpdated", suite.Client.NetworkV1, subnetNRef,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
				Metadata:      expectSubnetMeta,
				Spec:          expectSubnetSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Public IP lifecycle
	t.WithNewStep("PublicIp", func(ipCtx provider.StepCtx) {
		ipSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, ipCtx)

		ipSteps.CreateOrUpdatePublicIpV1Step("Create", suite.Client.NetworkV1, publicIp,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
				Metadata:      expectPublicIpMeta,
				Spec:          expectPublicIpSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		ipSteps.GetPublicIpV1Step("Get", suite.Client.NetworkV1, publicIpWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
				Metadata:      expectPublicIpMeta,
				Spec:          expectPublicIpSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		publicIp.Spec = suite.params.PublicIpUpdated.Spec
		expectPublicIpSpec.Address = publicIp.Spec.Address
		ipSteps.CreateOrUpdatePublicIpV1Step("Update", suite.Client.NetworkV1, publicIp,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
				Metadata:      expectPublicIpMeta,
				Spec:          expectPublicIpSpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		ipSteps.GetPublicIpV1Step("GetUpdated", suite.Client.NetworkV1, publicIpWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
				Metadata:      expectPublicIpMeta,
				Spec:          expectPublicIpSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Nic lifecycle
	t.WithNewStep("Nic", func(nicCtx provider.StepCtx) {
		nicSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, nicCtx)

		nicSteps.CreateOrUpdateNicV1Step("Create", suite.Client.NetworkV1, nic,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
				Metadata:      expectNicMeta,
				Spec:          expectNicSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		nicSteps.GetNicV1Step("Get", suite.Client.NetworkV1, nicWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
				Metadata:      expectNicMeta,
				Spec:          expectNicSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		nic.Spec = suite.params.NicUpdated.Spec
		expectNicSpec = &nic.Spec
		nicSteps.CreateOrUpdateNicV1Step("Update", suite.Client.NetworkV1, nic,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
				Metadata:      expectNicMeta,
				Spec:          expectNicSpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		nicSteps.GetNicV1Step("GetUpdated", suite.Client.NetworkV1, nicWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
				Metadata:      expectNicMeta,
				Spec:          expectNicSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Security group lifecycle
	t.WithNewStep("SecurityGroup", func(sgCtx provider.StepCtx) {
		sgSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, sgCtx)

		sgSteps.CreateOrUpdateSecurityGroupV1Step("Create", suite.Client.NetworkV1, group,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
				Metadata:      expectGroupMeta,
				Spec:          expectGroupSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		sgSteps.GetSecurityGroupV1Step("Get", suite.Client.NetworkV1, groupWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
				Metadata:      expectGroupMeta,
				Spec:          expectGroupSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		group.Spec.Rules[0] = schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionEgress}
		expectGroupSpec.Rules = group.Spec.Rules
		sgSteps.CreateOrUpdateSecurityGroupV1Step("Update", suite.Client.NetworkV1, group,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
				Metadata:      expectGroupMeta,
				Spec:          expectGroupSpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		sgSteps.GetSecurityGroupV1Step("GetUpdated", suite.Client.NetworkV1, groupWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
				Metadata:      expectGroupMeta,
				Spec:          expectGroupSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Block storage lifecycle
	t.WithNewStep("BlockStorage", func(bsCtx provider.StepCtx) {
		bsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, bsCtx)

		bsSteps.CreateOrUpdateBlockStorageV1Step("Create", suite.Client.StorageV1, block,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		bsSteps.GetBlockStorageV1Step("Get", suite.Client.StorageV1, blockWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Instance lifecycle
	t.WithNewStep("Instance", func(instCtx provider.StepCtx) {
		instSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, instCtx)

		instSteps.CreateOrUpdateInstanceV1Step("Create", suite.Client.ComputeV1, instance,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		instance = instSteps.GetInstanceV1Step("Get", suite.Client.ComputeV1, instanceWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Resources deletion
	t.WithNewStep("Delete", func(delCtx provider.StepCtx) {
		delSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, delCtx)
		delSteps.DeleteInstanceV1Step("Delete the instance", suite.Client.ComputeV1, instance)
		delSteps.GetInstanceWithErrorV1Step("Get the deleted instance", suite.Client.ComputeV1, instanceWRef, secapi.ErrResourceNotFound)

		delSteps.DeleteBlockStorageV1Step("Delete the block storage", suite.Client.StorageV1, block)
		delSteps.GetBlockStorageWithErrorV1Step("Get the deleted block storage", suite.Client.StorageV1, blockWRef, secapi.ErrResourceNotFound)

		delSteps.DeleteSecurityGroupV1Step("Delete the security group", suite.Client.NetworkV1, group)
		delSteps.GetSecurityGroupWithErrorV1Step("Get deleted security group", suite.Client.NetworkV1, groupWRef, secapi.ErrResourceNotFound)

		delSteps.DeleteNicV1Step("Delete the nic", suite.Client.NetworkV1, nic)
		delSteps.GetNicWithErrorV1Step("Get deleted nic", suite.Client.NetworkV1, nicWRef, secapi.ErrResourceNotFound)

		delSteps.DeletePublicIpV1Step("Delete the public ip", suite.Client.NetworkV1, publicIp)
		delSteps.GetPublicIpWithErrorV1Step("Get deleted public ip", suite.Client.NetworkV1, publicIpWRef, secapi.ErrResourceNotFound)

		delSteps.DeleteSubnetV1Step("Delete the subnet", suite.Client.NetworkV1, subnet)
		delSteps.GetSubnetWithErrorV1Step("Get deleted subnet", suite.Client.NetworkV1, subnetNRef, secapi.ErrResourceNotFound)

		delSteps.DeleteRouteTableV1Step("Delete the route table", suite.Client.NetworkV1, route)
		delSteps.GetRouteTableWithErrorV1Step("Get deleted route table", suite.Client.NetworkV1, routeNRef, secapi.ErrResourceNotFound)

		delSteps.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, gateway)
		delSteps.GetInternetGatewayWithErrorV1Step("Get deleted internet gateway", suite.Client.NetworkV1, gatewayWRef, secapi.ErrResourceNotFound)

		delSteps.DeleteNetworkV1Step("Delete the network", suite.Client.NetworkV1, network)
		delSteps.GetNetworkWithErrorV1Step("Get deleted network", suite.Client.NetworkV1, networkWRef, secapi.ErrResourceNotFound)

		delSteps.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
		delSteps.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)
	})

	suite.FinishScenario()
}

func (suite *NetworkLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
