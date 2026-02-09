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

type NetworkListV1TestSuite struct {
	suites.RegionalTestSuite

	config *NetworkListV1Config
	params *params.NetworkListV1Params
}

type NetworkListV1Config struct {
	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	StorageSkus    []string
	InstanceSkus   []string
	NetworkSkus    []string
}

func CreateListV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *NetworkListV1Config) *NetworkListV1TestSuite {
	suite := &NetworkListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.NetworkV1ListSuiteName
	return suite
}

func (suite *NetworkListV1TestSuite) BeforeAll(t provider.T) {
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

	params := &params.NetworkListV1Params{
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
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureListScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *NetworkListV1TestSuite) TestScenario(t provider.T) {
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

	workspace := suite.params.Workspace
	networks := suite.params.Networks
	gateways := suite.params.InternetGateways
	routes := suite.params.RouteTables
	subnets := suite.params.Subnets
	publicIps := suite.params.PublicIps
	nics := suite.params.Nics
	groups := suite.params.SecurityGroups

	t.WithNewStep("Workspace", func(wsCtx provider.StepCtx) {
		wsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, wsCtx)

		expectWorkspaceMeta := workspace.Metadata
		expectWorkspaceLabels := workspace.Labels
		wsSteps.CreateOrUpdateWorkspaceV1Step("Create", suite.Client.WorkspaceV1, workspace,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectWorkspaceLabels,
				Metadata:      expectWorkspaceMeta,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	})

	t.WithNewStep("Network", func(nwCtx provider.StepCtx) {
		nwSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, nwCtx)

		for _, network := range networks {
			n := network
			expectNetworkMeta := n.Metadata
			expectNetworkSpec := &n.Spec
			nwSteps.CreateOrUpdateNetworkV1Step("Create", suite.Client.NetworkV1, &n,
				steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
					Metadata:      expectNetworkMeta,
					Spec:          expectNetworkSpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		wref := secapi.WorkspaceReference{
			Name:      workspace.Metadata.Name,
			Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
			Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		}
		nwSteps.GetListNetworkV1Step("ListAll", suite.Client.NetworkV1, wref, nil)
		nwSteps.GetListNetworkV1Step("ListWithLimit", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLimit(1))
		nwSteps.GetListNetworkV1Step("ListWithLabel", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		nwSteps.GetListNetworkV1Step("ListWithLabelAndLimit", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("Skus", func(skuCtx provider.StepCtx) {
		skuSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, skuCtx)

		tenantRef := secapi.TenantReference{Tenant: secapi.TenantID(workspace.Metadata.Tenant)}
		skuSteps.GetListNetworkSkusV1Step("ListAll", suite.Client.NetworkV1, tenantRef, nil)
		skuSteps.GetListNetworkSkusV1Step("ListWithLimit", suite.Client.NetworkV1, tenantRef,
			secapi.NewListOptions().WithLimit(1))
	})

	t.WithNewStep("InternetGateway", func(igCtx provider.StepCtx) {
		igSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, igCtx)

		for _, gateway := range gateways {
			g := gateway
			expectGatewayMeta := g.Metadata
			expectGatewaySpec := &g.Spec
			igSteps.CreateOrUpdateInternetGatewayV1Step("Create", suite.Client.NetworkV1, &g,
				steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
					Metadata:      expectGatewayMeta,
					Spec:          expectGatewaySpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		wref := secapi.WorkspaceReference{
			Name:      workspace.Metadata.Name,
			Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
			Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		}
		igSteps.GetListInternetGatewayV1Step("ListAll", suite.Client.NetworkV1, wref, nil)
		igSteps.GetListInternetGatewayV1Step("ListWithLimit", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLimit(1))
		igSteps.GetListInternetGatewayV1Step("ListWithLabel", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		igSteps.GetListInternetGatewayV1Step("ListWithLabelAndLimit", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("RouteTable", func(rtCtx provider.StepCtx) {
		rtSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, rtCtx)

		for _, route := range routes {
			r := route
			expectRouteMeta := r.Metadata
			expectRouteSpec := &r.Spec
			rtSteps.CreateOrUpdateRouteTableV1Step("Create", suite.Client.NetworkV1, &r,
				steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
					Metadata:      expectRouteMeta,
					Spec:          expectRouteSpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		nref := secapi.NetworkReference{
			Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
			Network:   secapi.NetworkID(networks[0].Metadata.Name),
			Name:      routes[0].Metadata.Name,
		}
		rtSteps.GetListRouteTableV1Step("ListAll", suite.Client.NetworkV1, nref, nil)
		rtSteps.GetListRouteTableV1Step("ListWithLimit", suite.Client.NetworkV1, nref,
			secapi.NewListOptions().WithLimit(1))
		rtSteps.GetListRouteTableV1Step("ListWithLabel", suite.Client.NetworkV1, nref,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		rtSteps.GetListRouteTableV1Step("ListWithLabelAndLimit", suite.Client.NetworkV1, nref,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("Subnet", func(snCtx provider.StepCtx) {
		snSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, snCtx)

		for _, subnet := range subnets {
			s := subnet
			expectSubnetMeta := s.Metadata
			expectSubnetSpec := &s.Spec
			snSteps.CreateOrUpdateSubnetV1Step("Create", suite.Client.NetworkV1, &s,
				steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
					Metadata:      expectSubnetMeta,
					Spec:          expectSubnetSpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		nref := secapi.NetworkReference{
			Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
			Network:   secapi.NetworkID(networks[0].Metadata.Name),
			Name:      subnets[0].Metadata.Name,
		}
		snSteps.GetListSubnetV1Step("ListAll", suite.Client.NetworkV1, nref, nil)
		snSteps.GetListSubnetV1Step("ListWithLimit", suite.Client.NetworkV1, nref,
			secapi.NewListOptions().WithLimit(1))
		snSteps.GetListSubnetV1Step("ListWithLabel", suite.Client.NetworkV1, nref,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		snSteps.GetListSubnetV1Step("ListWithLabelAndLimit", suite.Client.NetworkV1, nref,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("PublicIp", func(piCtx provider.StepCtx) {
		piSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, piCtx)

		for _, publicIp := range publicIps {
			p := publicIp
			expectPublicIpMeta := p.Metadata
			expectPublicIpSpec := &p.Spec
			piSteps.CreateOrUpdatePublicIpV1Step("Create", suite.Client.NetworkV1, &p,
				steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
					Metadata:      expectPublicIpMeta,
					Spec:          expectPublicIpSpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		wref := secapi.WorkspaceReference{
			Name:      workspace.Metadata.Name,
			Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
			Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		}
		piSteps.GetListPublicIpV1Step("ListAll", suite.Client.NetworkV1, wref, nil)
		piSteps.GetListPublicIpV1Step("ListWithLimit", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLimit(1))
		piSteps.GetListPublicIpV1Step("ListWithLabel", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		piSteps.GetListPublicIpV1Step("ListWithLabelAndLimit", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("Nic", func(nicCtx provider.StepCtx) {
		nicSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, nicCtx)

		for _, nic := range nics {
			n := nic
			expectNicMeta := n.Metadata
			expectNicSpec := &n.Spec
			nicSteps.CreateOrUpdateNicV1Step("Create", suite.Client.NetworkV1, &n,
				steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
					Metadata:      expectNicMeta,
					Spec:          expectNicSpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		wref := secapi.WorkspaceReference{
			Name:      workspace.Metadata.Name,
			Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
			Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		}
		nicSteps.GetListNicV1Step("ListAll", suite.Client.NetworkV1, wref, nil)
		nicSteps.GetListNicV1Step("ListWithLimit", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLimit(1))
		nicSteps.GetListNicV1Step("ListWithLabel", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		nicSteps.GetListNicV1Step("ListWithLabelAndLimit", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("SecurityGroup", func(sgCtx provider.StepCtx) {
		gSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, sgCtx)

		for _, group := range groups {
			g := group
			expectGroupMeta := g.Metadata
			expectGroupSpec := &g.Spec
			gSteps.CreateOrUpdateSecurityGroupV1Step("Create", suite.Client.NetworkV1, &g,
				steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
					Metadata:      expectGroupMeta,
					Spec:          expectGroupSpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		wref := secapi.WorkspaceReference{
			Name:      workspace.Metadata.Name,
			Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
			Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		}
		gSteps.GetListSecurityGroupV1Step("ListAll", suite.Client.NetworkV1, wref, nil)
		gSteps.GetListSecurityGroupV1Step("ListWithLimit", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLimit(1))
		gSteps.GetListSecurityGroupV1Step("ListWithLabel", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		gSteps.GetListSecurityGroupV1Step("ListWithLabelAndLimit", suite.Client.NetworkV1, wref,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("Delete", func(delCtx provider.StepCtx) {
		delSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, delCtx)

		for _, group := range groups {
			g := group
			delSteps.DeleteSecurityGroupV1Step("SecurityGroup", suite.Client.NetworkV1, &g)

			groupWRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(g.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(g.Metadata.Workspace),
				Name:      g.Metadata.Name,
			}
			delSteps.GetSecurityGroupWithErrorV1Step("GetDeletedSecurityGroup", suite.Client.NetworkV1, groupWRef, secapi.ErrResourceNotFound)
		}

		for _, nic := range nics {
			n := nic
			delSteps.DeleteNicV1Step("Nic", suite.Client.NetworkV1, &n)

			nicWRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(n.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(n.Metadata.Workspace),
				Name:      n.Metadata.Name,
			}
			delSteps.GetNicWithErrorV1Step("GetDeletedNic", suite.Client.NetworkV1, nicWRef, secapi.ErrResourceNotFound)
		}

		for _, publicIp := range publicIps {
			p := publicIp
			delSteps.DeletePublicIpV1Step("PublicIp", suite.Client.NetworkV1, &p)

			publicIpWRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(p.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(p.Metadata.Workspace),
				Name:      p.Metadata.Name,
			}
			delSteps.GetPublicIpWithErrorV1Step("GetDeletedPublicIp", suite.Client.NetworkV1, publicIpWRef, secapi.ErrResourceNotFound)
		}

		for _, subnet := range subnets {
			s := subnet
			delSteps.DeleteSubnetV1Step("Subnet", suite.Client.NetworkV1, &s)

			subnetNRef := secapi.NetworkReference{
				Tenant:    secapi.TenantID(s.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(s.Metadata.Workspace),
				Network:   secapi.NetworkID(s.Metadata.Network),
				Name:      s.Metadata.Name,
			}
			delSteps.GetSubnetWithErrorV1Step("GetDeletedSubnet", suite.Client.NetworkV1, subnetNRef, secapi.ErrResourceNotFound)
		}

		for _, route := range routes {
			r := route
			delSteps.DeleteRouteTableV1Step("RouteTable", suite.Client.NetworkV1, &r)

			routeNRef := secapi.NetworkReference{
				Tenant:    secapi.TenantID(r.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(r.Metadata.Workspace),
				Network:   secapi.NetworkID(r.Metadata.Network),
				Name:      r.Metadata.Name,
			}
			delSteps.GetRouteTableWithErrorV1Step("GetDeletedRouteTable", suite.Client.NetworkV1, routeNRef, secapi.ErrResourceNotFound)
		}

		for _, gateway := range gateways {
			g := gateway
			delSteps.DeleteInternetGatewayV1Step("InternetGateway", suite.Client.NetworkV1, &g)

			gatewayWRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(g.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(g.Metadata.Workspace),
				Name:      g.Metadata.Name,
			}
			delSteps.GetInternetGatewayWithErrorV1Step("GetDeletedInternetGateway", suite.Client.NetworkV1, gatewayWRef, secapi.ErrResourceNotFound)
		}

		for _, network := range networks {
			n := network
			delSteps.DeleteNetworkV1Step("Network", suite.Client.NetworkV1, &n)

			networkWRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(n.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(n.Metadata.Workspace),
				Name:      n.Metadata.Name,
			}
			delSteps.GetNetworkWithErrorV1Step("GetDeletedNetwork", suite.Client.NetworkV1, networkWRef, secapi.ErrResourceNotFound)
		}

		workspaceTRef := secapi.TenantReference{
			Tenant: secapi.TenantID(workspace.Metadata.Tenant),
			Name:   workspace.Metadata.Name,
		}
		delSteps.DeleteWorkspaceV1Step("Workspace", suite.Client.WorkspaceV1, workspace)
		delSteps.GetWorkspaceWithErrorV1Step("GetDeletedWorkspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)
	})

	suite.FinishScenario()
}

func (suite *NetworkListV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
