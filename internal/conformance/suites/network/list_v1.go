package network

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockNetwork "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/network"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"k8s.io/utils/ptr"
)

type ListV1TestSuite struct {
	suites.RegionalTestSuite

	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	StorageSkus    []string
	InstanceSkus   []string
	NetworkSkus    []string
}

func (suite *ListV1TestSuite) TestScenario(t provider.T) {
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

	var err error

	// Generate the subnet cidr
	subnetCidr, err := generators.GenerateSubnetCidr(suite.NetworkCidr, 8, 1)
	if err != nil {
		t.Fatalf("Failed to generate subnet cidr: %v", err)
	}

	// Generate the nic addresses
	nicAddress1, err := generators.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}

	// Generate the public ips
	publicIpAddress1, err := generators.GeneratePublicIp(suite.PublicIpsRange, 1)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}

	// Select zones
	zone1 := suite.RegionZones[rand.Intn(len(suite.RegionZones))]

	// Select skus
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]
	instanceSkuName := suite.InstanceSkus[rand.Intn(len(suite.InstanceSkus))]
	networkSkuName1 := suite.NetworkSkus[rand.Intn(len(suite.NetworkSkus))]

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

	// Setup mock, if configured to use
	if suite.MockEnabled {
		mockParams := &params.NetworkListParamsV1{
			BaseParams: &params.BaseParams{
				Tenant: suite.Tenant,
				Region: suite.Region,
				MockParams: &mock.MockParams{
					ServerURL: *suite.MockServerURL,
					AuthToken: suite.AuthToken,
				},
			},
			Workspace: &params.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					constants.EnvLabel: constants.EnvConformanceLabel,
				},
			},
			BlockStorage: &params.ResourceParams[schema.BlockStorageSpec]{
				Name: blockStorageName,
				InitialLabels: schema.Labels{
					constants.EnvLabel: constants.EnvConformanceLabel,
				},
				InitialSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: blockStorageSize,
				},
			},
			Instance: &params.ResourceParams[schema.InstanceSpec]{
				Name: instanceName,
				InitialLabels: schema.Labels{
					constants.EnvLabel: constants.EnvConformanceLabel,
				},
				InitialSpec: &schema.InstanceSpec{
					SkuRef: *instanceSkuRefObj,
					Zone:   zone1,
					BootVolume: schema.VolumeReference{
						DeviceRef: *blockStorageRefObj,
					},
				},
			},
			Networks: []params.ResourceParams[schema.NetworkSpec]{
				{
					Name: networkName,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.NetworkSpec{
						Cidr:          schema.Cidr{Ipv4: ptr.To(suite.NetworkCidr)},
						SkuRef:        *networkSkuRefObj,
						RouteTableRef: *routeTableRefObj,
					},
				},
				{
					Name: networkName2,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.NetworkSpec{
						Cidr:          schema.Cidr{Ipv4: ptr.To(suite.NetworkCidr)},
						SkuRef:        *networkSkuRefObj,
						RouteTableRef: *routeTableRefObj,
					},
				},
			},
			InternetGateways: []params.ResourceParams[schema.InternetGatewaySpec]{
				{
					Name: internetGatewayName,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)},
				},
				{
					Name: internetGatewayName2,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)},
				},
			},
			RouteTables: []params.ResourceParams[schema.RouteTableSpec]{
				{
					Name: routeTableName,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.RouteTableSpec{
						Routes: []schema.RouteSpec{
							{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
						},
					},
				},
				{
					Name: routeTableName2,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.RouteTableSpec{
						Routes: []schema.RouteSpec{
							{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
						},
					},
				},
			},
			Subnets: []params.ResourceParams[schema.SubnetSpec]{
				{
					Name: subnetName,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.SubnetSpec{
						Cidr: schema.Cidr{Ipv4: &subnetCidr},
						Zone: zone1,
					},
				}, {
					Name: subnetName2,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.SubnetSpec{
						Cidr: schema.Cidr{Ipv4: &subnetCidr},
						Zone: zone1,
					},
				},
			},
			Nics: []params.ResourceParams[schema.NicSpec]{
				{
					Name: nicName,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.NicSpec{
						Addresses:    []string{nicAddress1},
						PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
						SubnetRef:    *subnetRefObj,
					},
				},
				{
					Name: nicName2,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.NicSpec{
						Addresses:    []string{nicAddress1},
						PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
						SubnetRef:    *subnetRefObj,
					},
				},
			},
			PublicIps: []params.ResourceParams[schema.PublicIpSpec]{
				{
					Name: publicIpName,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.PublicIpSpec{
						Version: schema.IPVersionIPv4,
						Address: ptr.To(publicIpAddress1),
					},
				},
				{
					Name: publicIpName2,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.PublicIpSpec{
						Version: schema.IPVersionIPv4,
						Address: ptr.To(publicIpAddress1),
					},
				},
			},
			SecurityGroups: []params.ResourceParams[schema.SecurityGroupSpec]{
				{
					Name: securityGroupName,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.SecurityGroupSpec{
						Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
					},
				},
				{
					Name: securityGroupName2,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.SecurityGroupSpec{
						Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
					},
				},
			},
		}
		wm, err := mockNetwork.ConfigureListScenarioV1(suite.ScenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.MockClient = wm
	}

	stepsBuilder := steps.NewStepsConfigurator(&suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.Tenant,
			Name:   workspaceName,
		},
	}
	expectWorkspaceMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)
	// Network

	// Create a network
	networks := &[]schema.Network{
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      networkName,
			},
			Spec: schema.NetworkSpec{
				Cidr:          schema.Cidr{Ipv4: &suite.NetworkCidr},
				SkuRef:        *networkSkuRefObj,
				RouteTableRef: *routeTableRefObj,
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      networkName2,
			},
			Spec: schema.NetworkSpec{
				Cidr:          schema.Cidr{Ipv4: &suite.NetworkCidr},
				SkuRef:        *networkSkuRefObj,
				RouteTableRef: *routeTableRefObj,
			},
		},
	}

	for _, network := range *networks {

		expectNetworkMeta, err := builders.NewNetworkMetadataBuilder().
			Name(network.Metadata.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}
		expectNetworkSpec := &schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: &suite.NetworkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}
		stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", suite.Client.NetworkV1, &network,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
				Metadata:      expectNetworkMeta,
				Spec:          expectNetworkSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	wref := secapi.WorkspaceReference{
		Name:      workspaceName,
		Workspace: secapi.WorkspaceID(workspaceName),
		Tenant:    secapi.TenantID(suite.Tenant),
	}

	// List Network
	stepsBuilder.GetListNetworkV1Step("List Network", suite.Client.NetworkV1, wref, nil)

	// List Network with limit
	stepsBuilder.GetListNetworkV1Step("Get list of Network with limit", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List Network with Label
	stepsBuilder.GetListNetworkV1Step("Get list of Network with label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Network with Limit and label
	stepsBuilder.GetListNetworkV1Step("Get list of Network with limit and label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Network Skus
	// List SKUS
	stepsBuilder.GetListNetworkSkusV1Step("List skus", suite.Client.NetworkV1, secapi.TenantReference{Tenant: secapi.TenantID(suite.Tenant)}, nil)

	// List SKUS with limit
	stepsBuilder.GetListNetworkSkusV1Step("Get list of skus", suite.Client.NetworkV1, secapi.TenantReference{Tenant: secapi.TenantID(suite.Tenant)},
		secapi.NewListOptions().WithLimit(1))

	// Internet gateway

	// Create an internet gateway
	gateways := &[]schema.InternetGateway{
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      internetGatewayName,
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      internetGatewayName2,
			},
		},
	}

	for _, gateway := range *gateways {
		expectGatewayMeta, err := builders.NewInternetGatewayMetadataBuilder().
			Name(gateway.Metadata.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}
		expectGatewaySpec := &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)}
		stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create a internet gateway", suite.Client.NetworkV1, &gateway,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
				Metadata:      expectGatewayMeta,
				Spec:          expectGatewaySpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

	}

	// List Internet Gateway
	stepsBuilder.GetListInternetGatewayV1Step("List Internet Gateway", suite.Client.NetworkV1, wref, nil)

	// List Internet Gateway with limit
	stepsBuilder.GetListInternetGatewayV1Step("Get list of Internet Gateway with limit", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List Internet Gateway with Label
	stepsBuilder.GetListInternetGatewayV1Step("Get list of Internet Gateway with label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Internet Gateway with Limit and label
	stepsBuilder.GetListInternetGatewayV1Step("Get list of Internet Gateway with limit and label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Route table

	// Create a route table
	routes := &[]schema.RouteTable{
		{
			Metadata: &schema.RegionalNetworkResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Network:   networkName,
				Name:      routeTableName,
			},
			Spec: schema.RouteTableSpec{
				Routes: []schema.RouteSpec{
					{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
				},
			},
		},
		{
			Metadata: &schema.RegionalNetworkResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Network:   networkName,
				Name:      routeTableName2,
			},
			Spec: schema.RouteTableSpec{
				Routes: []schema.RouteSpec{
					{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
				},
			},
		},
	}
	for _, route := range *routes {
		expectRouteMeta, err := builders.NewRouteTableMetadataBuilder().
			Name(route.Metadata.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Network(networkName).Region(suite.Region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}
		expectRouteSpec := &schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}
		stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", suite.Client.NetworkV1, &route,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
				Metadata:      expectRouteMeta,
				Spec:          expectRouteSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}
	// Get the created route table
	nref := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      routeTableName,
	}
	// List Route table
	stepsBuilder.GetListRouteTableV1Step("List Route table", suite.Client.NetworkV1, *nref, nil)

	// List Route table with limit
	stepsBuilder.GetListRouteTableV1Step("Get list of Route table with limit", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1))

	// List Route table with Label
	stepsBuilder.GetListRouteTableV1Step("Get list of Route table with label", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Route table with Limit and label
	stepsBuilder.GetListRouteTableV1Step("Get list of Route table with limit and label", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	// Subnet

	// Create a subnet
	subnets := &[]schema.Subnet{
		{
			Metadata: &schema.RegionalNetworkResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Network:   networkName,
				Name:      subnetName,
			},
			Spec: schema.SubnetSpec{
				Cidr: schema.Cidr{Ipv4: &subnetCidr},
				Zone: zone1,
			},
		},
		{
			Metadata: &schema.RegionalNetworkResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Network:   networkName,
				Name:      subnetName2,
			},
			Spec: schema.SubnetSpec{
				Cidr: schema.Cidr{Ipv4: &subnetCidr},
				Zone: zone1,
			},
		},
	}

	for _, subnet := range *subnets {

		expectSubnetMeta, err := builders.NewSubnetMetadataBuilder().
			Name(subnet.Metadata.Name).
			Provider(constants.NetworkProviderV1).
			ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Network(networkName).Region(suite.Region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}
		expectSubnetSpec := &schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone1,
		}
		stepsBuilder.CreateOrUpdateSubnetV1Step("Create a subnet", suite.Client.NetworkV1, &subnet,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
				Metadata:      expectSubnetMeta,
				Spec:          expectSubnetSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List Subnet
	stepsBuilder.GetListSubnetV1Step("List Subnet", suite.Client.NetworkV1, *nref, nil)

	// List Subnet with limit
	stepsBuilder.GetListSubnetV1Step("Get list of Subnet with limit", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1))

	// List Subnet with Label
	stepsBuilder.GetListSubnetV1Step("Get list of Subnet with label", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Subnet with Limit and label
	stepsBuilder.GetListSubnetV1Step("Get list of Subnet with limit and label", suite.Client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Public ip

	// Create a public ip
	publicIps := &[]schema.PublicIp{
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      publicIpName,
			},
			Spec: schema.PublicIpSpec{
				Address: &publicIpAddress1,
				Version: schema.IPVersionIPv4,
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      publicIpName2,
			},
			Spec: schema.PublicIpSpec{
				Address: &publicIpAddress1,
				Version: schema.IPVersionIPv4,
			},
		},
	}

	for _, publicIp := range *publicIps {
		expectPublicIpMeta, err := builders.NewPublicIpMetadataBuilder().
			Name(publicIp.Metadata.Name).
			Provider(constants.NetworkProviderV1).
			ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}
		expectPublicIpSpec := &schema.PublicIpSpec{
			Address: &publicIpAddress1,
			Version: schema.IPVersionIPv4,
		}
		stepsBuilder.CreateOrUpdatePublicIpV1Step("Create a public ip", suite.Client.NetworkV1, &publicIp,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
				Metadata:      expectPublicIpMeta,
				Spec:          expectPublicIpSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List PublicIP
	stepsBuilder.GetListPublicIpV1Step("List PublicIP", suite.Client.NetworkV1, wref, nil)

	// List PublicIP with limit
	stepsBuilder.GetListPublicIpV1Step("Get list of PublicIP with limit", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List PublicIP with Label
	stepsBuilder.GetListPublicIpV1Step("Get list of PublicIP with label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List PublicIP with Limit and label
	stepsBuilder.GetListPublicIpV1Step("Get list of PublicIP with limit and label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Nic

	// Create a nic
	nics := &[]schema.Nic{
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      nicName,
			},
			Spec: schema.NicSpec{
				Addresses:    []string{nicAddress1},
				PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
				SubnetRef:    *subnetRefObj,
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      nicName2,
			},
			Spec: schema.NicSpec{
				Addresses:    []string{nicAddress1},
				PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
				SubnetRef:    *subnetRefObj,
			},
		},
	}

	for _, nic := range *nics {

		expectNicMeta, err := builders.NewNicMetadataBuilder().
			Name(nic.Metadata.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Build()
		expectNicSpec := &schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}
		stepsBuilder.CreateOrUpdateNicV1Step("Create a nic", suite.Client.NetworkV1, &nic,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
				Metadata:      expectNicMeta,
				Spec:          expectNicSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}
	// List Nic
	stepsBuilder.GetListNicV1Step("List Nic", suite.Client.NetworkV1, wref, nil)

	// List Nic with limit
	stepsBuilder.GetListNicV1Step("Get list of Nic with limit", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List Nic with Label
	stepsBuilder.GetListNicV1Step("Get list of Nic with label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Nic with Limit and label
	stepsBuilder.GetListNicV1Step("Get list of Nic with limit and label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Security Group

	// Create a security group
	groups := &[]schema.SecurityGroup{
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      securityGroupName,
			},
			Spec: schema.SecurityGroupSpec{
				Rules: []schema.SecurityGroupRuleSpec{
					{Direction: schema.SecurityGroupRuleDirectionIngress},
				},
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      securityGroupName2,
			},
			Spec: schema.SecurityGroupSpec{
				Rules: []schema.SecurityGroupRuleSpec{
					{Direction: schema.SecurityGroupRuleDirectionIngress},
				},
			},
		},
	}

	for _, group := range *groups {
		expectGroupMeta, err := builders.NewSecurityGroupMetadataBuilder().
			Name(group.Metadata.Name).
			Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Build()
		expectGroupSpec := &schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{
				{Direction: schema.SecurityGroupRuleDirectionIngress},
			},
		}
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}
		stepsBuilder.CreateOrUpdateSecurityGroupV1Step("Create a security group", suite.Client.NetworkV1, &group,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
				Metadata:      expectGroupMeta,
				Spec:          expectGroupSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}
	// List Security Group
	stepsBuilder.GetListSecurityGroupV1Step("List Security Group", suite.Client.NetworkV1, wref, nil)

	// List Security Group with limit
	stepsBuilder.GetListSecurityGroupV1Step("Get list of Security Group with limit", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List Security Group with Label
	stepsBuilder.GetListSecurityGroupV1Step("Get list of Security Group with label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Nic with Limit and label
	stepsBuilder.GetListSecurityGroupV1Step("Get list of Security Group with limit and label", suite.Client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Delete all security groups
	for _, group := range *groups {
		stepsBuilder.DeleteSecurityGroupV1Step("Delete the security group", suite.Client.NetworkV1, &group)

		// Get deleted security group
		groupWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.Tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      group.Metadata.Name,
		}
		stepsBuilder.GetSecurityGroupWithErrorV1Step("Get deleted security group", suite.Client.NetworkV1, *groupWRef, secapi.ErrResourceNotFound)
	}

	// Delete all NICs
	for _, nic := range *nics {
		stepsBuilder.DeleteNicV1Step("Delete the nic", suite.Client.NetworkV1, &nic)

		// Get the deleted nic
		nicWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.Tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nic.Metadata.Name,
		}
		stepsBuilder.GetNicWithErrorV1Step("Get deleted nic", suite.Client.NetworkV1, *nicWRef, secapi.ErrResourceNotFound)
	}

	// Delete all public IPs
	for _, publicIp := range *publicIps {
		stepsBuilder.DeletePublicIpV1Step("Delete the public ip", suite.Client.NetworkV1, &publicIp)

		// Get the deleted public ip
		publicIpWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.Tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIp.Metadata.Name,
		}
		stepsBuilder.GetPublicIpWithErrorV1Step("Get deleted public ip", suite.Client.NetworkV1, *publicIpWRef, secapi.ErrResourceNotFound)
	}

	// Delete all subnets
	for _, subnet := range *subnets {
		stepsBuilder.DeleteSubnetV1Step("Delete the subnet", suite.Client.NetworkV1, &subnet)

		// Get the deleted subnet
		subnetNRef := &secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.Tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      subnet.Metadata.Name,
		}
		stepsBuilder.GetSubnetWithErrorV1Step("Get deleted subnet", suite.Client.NetworkV1, *subnetNRef, secapi.ErrResourceNotFound)
	}

	// Delete all route tables
	for _, route := range *routes {
		stepsBuilder.DeleteRouteTableV1Step("Delete the route table", suite.Client.NetworkV1, &route)

		// Get the deleted route table
		routeNRef := &secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.Tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      route.Metadata.Name,
		}
		stepsBuilder.GetRouteTableWithErrorV1Step("Get deleted route table", suite.Client.NetworkV1, *routeNRef, secapi.ErrResourceNotFound)
	}

	// Delete all internet gateways
	for _, gateway := range *gateways {
		stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, &gateway)

		// Get the deleted internet gateway
		gatewayWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.Tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      gateway.Metadata.Name,
		}
		stepsBuilder.GetInternetGatewayWithErrorV1Step("Get deleted internet gateway", suite.Client.NetworkV1, *gatewayWRef, secapi.ErrResourceNotFound)
	}

	// Delete all networks
	for _, network := range *networks {
		stepsBuilder.DeleteNetworkV1Step("Delete the network", suite.Client.NetworkV1, &network)

		// Get the deleted network
		networkWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.Tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      network.Metadata.Name,
		}
		stepsBuilder.GetNetworkWithErrorV1Step("Get deleted network", suite.Client.NetworkV1, *networkWRef, secapi.ErrResourceNotFound)
	}

	// Delete the workspace
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)

	// Get the deleted workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   workspaceName,
	}
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *ListV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
