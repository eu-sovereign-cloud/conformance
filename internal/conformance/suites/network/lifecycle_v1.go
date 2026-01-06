package network

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/network"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"k8s.io/utils/ptr"
)

type NetworkLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	NetworkCidr    string
	publicIpsRange string
	RegionZones    []string
	StorageSkus    []string
	InstanceSkus   []string
	NetworkSkus    []string
}

func (suite *NetworkLifeCycleV1TestSuite) TestLifeCycleScenario(t provider.T) {
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
	nicAddress2, err := generators.GenerateNicAddress(subnetCidr, 2)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}

	// Generate the public ips
	publicIpAddress1, err := generators.GeneratePublicIp(suite.publicIpsRange, 1)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}
	publicIpAddress2, err := generators.GeneratePublicIp(suite.publicIpsRange, 2)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}

	// Select zones
	zone1 := suite.RegionZones[rand.Intn(len(suite.RegionZones))]
	zone2 := suite.RegionZones[rand.Intn(len(suite.RegionZones))]

	// Select skus
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]
	instanceSkuName := suite.InstanceSkus[rand.Intn(len(suite.InstanceSkus))]
	networkSkuName1 := suite.NetworkSkus[rand.Intn(len(suite.NetworkSkus))]
	networkSkuName2 := suite.NetworkSkus[rand.Intn(len(suite.NetworkSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()
	storageSkuRef := generators.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRef := generators.GenerateBlockStorageRef(blockStorageName)
	blockStorageRefObj, err := secapi.BuildReferenceFromURN(blockStorageRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	instanceSkuRef := generators.GenerateSkuRef(instanceSkuName)
	instanceSkuRefObj, err := secapi.BuildReferenceFromURN(instanceSkuRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	instanceName := generators.GenerateInstanceName()
	instanceRef := generators.GenerateInstanceRef(instanceName)
	instanceRefObj, err := secapi.BuildReferenceFromURN(instanceRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	networkSkuRef1 := generators.GenerateSkuRef(networkSkuName1)
	networkSkuRefObj, err := secapi.BuildReferenceFromURN(networkSkuRef1)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}
	networkSkuRef2 := generators.GenerateSkuRef(networkSkuName2)
	networkSkuRef2Obj, err := secapi.BuildReferenceFromURN(networkSkuRef2)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	networkName := generators.GenerateNetworkName()

	internetGatewayName := generators.GenerateInternetGatewayName()
	internetGatewayRef := generators.GenerateInternetGatewayRef(internetGatewayName)
	internetGatewayRefObj, err := secapi.BuildReferenceFromURN(internetGatewayRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	routeTableName := generators.GenerateRouteTableName()
	routeTableRef := generators.GenerateRouteTableRef(routeTableName)
	routeTableRefObj, err := secapi.BuildReferenceFromURN(routeTableRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	subnetName := generators.GenerateSubnetName()
	subnetRef := generators.GenerateSubnetRef(subnetName)
	subnetRefObj, err := secapi.BuildReferenceFromURN(subnetRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	nicName := generators.GenerateNicName()

	publicIpName := generators.GeneratePublicIpName()
	publicIpRef := generators.GeneratePublicIpRef(publicIpName)
	publicIpRefObj, err := secapi.BuildReferenceFromURN(publicIpRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	securityGroupName := generators.GenerateSecurityGroupName()

	blockStorageSize := generators.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.MockEnabled {
		mockParams := &mock.NetworkLifeCycleParamsV1{
			BaseParams: &mock.BaseParams{
				MockURL:   *suite.MockServerURL,
				AuthToken: suite.AuthToken,
				Tenant:    suite.Tenant,
				Region:    suite.Region,
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					constants.EnvLabel: constants.EnvDevelopmentLabel,
				},
			},
			BlockStorage: &mock.ResourceParams[schema.BlockStorageSpec]{
				Name: blockStorageName,
				InitialSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: blockStorageSize,
				},
			},
			Instance: &mock.ResourceParams[schema.InstanceSpec]{
				Name: instanceName,
				InitialSpec: &schema.InstanceSpec{
					SkuRef: *instanceSkuRefObj,
					Zone:   zone1,
					BootVolume: schema.VolumeReference{
						DeviceRef: *blockStorageRefObj,
					},
				},
			},
			Network: &mock.ResourceParams[schema.NetworkSpec]{
				Name: networkName,
				InitialSpec: &schema.NetworkSpec{
					Cidr:          schema.Cidr{Ipv4: ptr.To(suite.NetworkCidr)},
					SkuRef:        *networkSkuRefObj,
					RouteTableRef: *routeTableRefObj,
				},
				UpdatedSpec: &schema.NetworkSpec{
					Cidr:          schema.Cidr{Ipv4: ptr.To(suite.NetworkCidr)},
					SkuRef:        *networkSkuRef2Obj,
					RouteTableRef: *routeTableRefObj,
				},
			},
			InternetGateway: &mock.ResourceParams[schema.InternetGatewaySpec]{
				Name:        internetGatewayName,
				InitialSpec: &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)},
				UpdatedSpec: &schema.InternetGatewaySpec{EgressOnly: ptr.To(true)},
			},
			RouteTable: &mock.ResourceParams[schema.RouteTableSpec]{
				Name: routeTableName,
				InitialSpec: &schema.RouteTableSpec{
					Routes: []schema.RouteSpec{
						{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
					},
				},
				UpdatedSpec: &schema.RouteTableSpec{
					Routes: []schema.RouteSpec{
						{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *instanceRefObj},
					},
				},
			},
			Subnet: &mock.ResourceParams[schema.SubnetSpec]{
				Name: subnetName,
				InitialSpec: &schema.SubnetSpec{
					Cidr: schema.Cidr{Ipv4: &subnetCidr},
					Zone: zone1,
				},
				UpdatedSpec: &schema.SubnetSpec{
					Cidr: schema.Cidr{Ipv4: &subnetCidr},
					Zone: zone2,
				},
			},
			Nic: &mock.ResourceParams[schema.NicSpec]{
				Name: nicName,
				InitialSpec: &schema.NicSpec{
					Addresses:    []string{nicAddress1},
					PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
					SubnetRef:    *subnetRefObj,
				},
				UpdatedSpec: &schema.NicSpec{
					Addresses:    []string{nicAddress2},
					PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
					SubnetRef:    *subnetRefObj,
				},
			},
			PublicIp: &mock.ResourceParams[schema.PublicIpSpec]{
				Name: publicIpName,
				InitialSpec: &schema.PublicIpSpec{
					Version: schema.IPVersionIPv4,
					Address: ptr.To(publicIpAddress1),
				},
				UpdatedSpec: &schema.PublicIpSpec{
					Version: schema.IPVersionIPv4,
					Address: ptr.To(publicIpAddress2),
				},
			},
			SecurityGroup: &mock.ResourceParams[schema.SecurityGroupSpec]{
				Name: securityGroupName,
				InitialSpec: &schema.SecurityGroupSpec{
					Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
				},
				UpdatedSpec: &schema.SecurityGroupSpec{
					Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionEgress}},
				},
			},
		}
		wm, err := network.ConfigureLifecycleScenarioV1(suite.ScenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.MockClient = wm
	}

	stepsBuilder := steps.NewBuilder(&suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
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
	expectWorkspaceLabels := schema.Labels{constants.EnvLabel: constants.EnvDevelopmentLabel}
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   workspaceName,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, *workspaceTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Network

	// Create a network
	network := &schema.Network{
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
	}
	expectNetworkMeta, err := builders.NewNetworkMetadataBuilder().
		Name(networkName).
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
	stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:      expectNetworkMeta,
			Spec:          expectNetworkSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created network
	networkWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      networkName,
	}
	stepsBuilder.GetNetworkV1Step("Get the created network", suite.Client.NetworkV1, *networkWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:      expectNetworkMeta,
			Spec:          expectNetworkSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the network
	network.Spec.SkuRef = *networkSkuRef2Obj
	expectNetworkSpec.SkuRef = network.Spec.SkuRef
	stepsBuilder.CreateOrUpdateNetworkV1Step("Update the network", suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:      expectNetworkMeta,
			Spec:          expectNetworkSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated network
	stepsBuilder.GetNetworkV1Step("Get the updated network", suite.Client.NetworkV1, *networkWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:      expectNetworkMeta,
			Spec:          expectNetworkSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Internet gateway

	// Create an internet gateway
	gateway := &schema.InternetGateway{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      internetGatewayName,
		},
	}
	expectGatewayMeta, err := builders.NewInternetGatewayMetadataBuilder().
		Name(internetGatewayName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectGatewaySpec := &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)}
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create a internet gateway", suite.Client.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:      expectGatewayMeta,
			Spec:          expectGatewaySpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created internet gateway
	gatewayWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      internetGatewayName,
	}
	stepsBuilder.GetInternetGatewayV1Step("Get the created internet gateway", suite.Client.NetworkV1, *gatewayWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:      expectGatewayMeta,
			Spec:          expectGatewaySpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the internet gateway
	gateway.Spec.EgressOnly = ptr.To(true)
	expectGatewaySpec.EgressOnly = gateway.Spec.EgressOnly
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Update the internet gateway", suite.Client.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:      expectGatewayMeta,
			Spec:          expectGatewaySpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated internet gateway
	stepsBuilder.GetInternetGatewayV1Step("Get the updated internet gateway", suite.Client.NetworkV1, *gatewayWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:      expectGatewayMeta,
			Spec:          expectGatewaySpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Route table

	// Create a route table
	route := &schema.RouteTable{
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
	}
	expectRouteMeta, err := builders.NewRouteTableMetadataBuilder().
		Name(routeTableName).
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
	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:      expectRouteMeta,
			Spec:          expectRouteSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created route table
	routeNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      routeTableName,
	}
	stepsBuilder.GetRouteTableV1Step("Get the created route table", suite.Client.NetworkV1, *routeNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:      expectRouteMeta,
			Spec:          expectRouteSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the route table
	route.Spec.Routes = []schema.RouteSpec{
		{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *instanceRefObj},
	}
	expectRouteSpec.Routes = route.Spec.Routes
	stepsBuilder.CreateOrUpdateRouteTableV1Step("Update the route table", suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:      expectRouteMeta,
			Spec:          expectRouteSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated route table
	stepsBuilder.GetRouteTableV1Step("Get the updated route table", suite.Client.NetworkV1, *routeNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:      expectRouteMeta,
			Spec:          expectRouteSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Subnet

	// Create a subnet
	subnet := &schema.Subnet{
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
	}
	expectSubnetMeta, err := builders.NewSubnetMetadataBuilder().
		Name(subnetName).
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
	stepsBuilder.CreateOrUpdateSubnetV1Step("Create a subnet", suite.Client.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:      expectSubnetMeta,
			Spec:          expectSubnetSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created subnet
	subnetNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      subnetName,
	}
	stepsBuilder.GetSubnetV1Step("Get the created subnet", suite.Client.NetworkV1, *subnetNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:      expectSubnetMeta,
			Spec:          expectSubnetSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the subnet
	subnet.Spec.Zone = zone2
	expectSubnetSpec.Zone = subnet.Spec.Zone
	stepsBuilder.CreateOrUpdateSubnetV1Step("Update the subnet", suite.Client.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:      expectSubnetMeta,
			Spec:          expectSubnetSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated subnet
	stepsBuilder.GetSubnetV1Step("Get the updated subnet", suite.Client.NetworkV1, *subnetNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:      expectSubnetMeta,
			Spec:          expectSubnetSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Public ip

	// Create a public ip
	publicIp := &schema.PublicIp{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      publicIpName,
		},
		Spec: schema.PublicIpSpec{
			Address: &publicIpAddress1,
			Version: schema.IPVersionIPv4,
		},
	}
	expectPublicIpMeta, err := builders.NewPublicIpMetadataBuilder().
		Name(publicIpName).
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
	stepsBuilder.CreateOrUpdatePublicIpV1Step("Create a public ip", suite.Client.NetworkV1, publicIp,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:      expectPublicIpMeta,
			Spec:          expectPublicIpSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created public ip
	publicIpWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      publicIpName,
	}
	stepsBuilder.GetPublicIpV1Step("Get the created public ip", suite.Client.NetworkV1, *publicIpWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:      expectPublicIpMeta,
			Spec:          expectPublicIpSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the public ip
	publicIp.Spec.Address = ptr.To(publicIpAddress2)
	expectPublicIpSpec.Address = publicIp.Spec.Address
	stepsBuilder.CreateOrUpdatePublicIpV1Step("Update the public ip", suite.Client.NetworkV1, publicIp,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:      expectPublicIpMeta,
			Spec:          expectPublicIpSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated public ip
	stepsBuilder.GetPublicIpV1Step("Get the updated public ip", suite.Client.NetworkV1, *publicIpWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:      expectPublicIpMeta,
			Spec:          expectPublicIpSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Nic

	// Create a nic
	nic := &schema.Nic{
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
	}
	expectNicMeta, err := builders.NewNicMetadataBuilder().
		Name(nicName).
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
	stepsBuilder.CreateOrUpdateNicV1Step("Create a nic", suite.Client.NetworkV1, nic,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:      expectNicMeta,
			Spec:          expectNicSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created nic
	nicWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      nicName,
	}
	stepsBuilder.GetNicV1Step("Get the created nic", suite.Client.NetworkV1, *nicWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:      expectNicMeta,
			Spec:          expectNicSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the nic
	nic.Spec.Addresses = []string{nicAddress2}
	expectNicSpec.Addresses = nic.Spec.Addresses
	stepsBuilder.CreateOrUpdateNicV1Step("Update the nic", suite.Client.NetworkV1, nic,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:      expectNicMeta,
			Spec:          expectNicSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated nic
	stepsBuilder.GetNicV1Step("Get the updated nic", suite.Client.NetworkV1, *nicWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:      expectNicMeta,
			Spec:          expectNicSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Security Group

	// Create a security group
	group := &schema.SecurityGroup{
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
	}
	expectGroupMeta, err := builders.NewSecurityGroupMetadataBuilder().
		Name(securityGroupName).
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
	stepsBuilder.CreateOrUpdateSecurityGroupV1Step("Create a security group", suite.Client.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:      expectGroupMeta,
			Spec:          expectGroupSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created security group
	groupWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      securityGroupName,
	}
	stepsBuilder.GetSecurityGroupV1Step("Get the created security group", suite.Client.NetworkV1, *groupWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:      expectGroupMeta,
			Spec:          expectGroupSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the security group
	group.Spec.Rules[0] = schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionEgress}
	expectGroupSpec.Rules = group.Spec.Rules
	stepsBuilder.CreateOrUpdateSecurityGroupV1Step("Update the security group", suite.Client.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:      expectGroupMeta,
			Spec:          expectGroupSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated security group
	stepsBuilder.GetSecurityGroupV1Step("Get the updated security group", suite.Client.NetworkV1, *groupWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:      expectGroupMeta,
			Spec:          expectGroupSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Block storage

	// Create a block storage
	block := &schema.BlockStorage{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      blockStorageName,
		},
		Spec: schema.BlockStorageSpec{
			SizeGB: blockStorageSize,
			SkuRef: *storageSkuRefObj,
		},
	}
	expectedBlockMeta, err := builders.NewBlockStorageMetadataBuilder().
		Name(blockStorageName).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Build()
	expectedBlockSpec := &schema.BlockStorageSpec{
		SizeGB: blockStorageSize,
		SkuRef: *storageSkuRefObj,
	}
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created block storage
	blockWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	stepsBuilder.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1, *blockWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Instance

	// Create an instance
	instance := &schema.Instance{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      instanceName,
		},
		Spec: schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone1,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		},
	}
	expectInstanceMeta, err := builders.NewInstanceMetadataBuilder().
		Name(instanceName).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Build()
	expectInstanceSpec := &schema.InstanceSpec{
		SkuRef: *instanceSkuRefObj,
		Zone:   zone1,
		BootVolume: schema.VolumeReference{
			DeviceRef: *blockStorageRefObj,
		},
	}
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	stepsBuilder.CreateOrUpdateInstanceV1Step("Create an instance", suite.Client.ComputeV1, instance,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:      expectInstanceMeta,
			Spec:          expectInstanceSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created instance
	instanceWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      instanceName,
	}
	instance = stepsBuilder.GetInstanceV1Step("Get the created instance", suite.Client.ComputeV1, *instanceWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:      expectInstanceMeta,
			Spec:          expectInstanceSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion

	stepsBuilder.DeleteInstanceV1Step("Delete the instance", suite.Client.ComputeV1, instance)
	stepsBuilder.GetInstanceWithErrorV1Step("Get the deleted instance", suite.Client.ComputeV1, *instanceWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", suite.Client.StorageV1, block)
	stepsBuilder.GetBlockStorageWithErrorV1Step("Get the deleted block storage", suite.Client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteSecurityGroupV1Step("Delete the security group", suite.Client.NetworkV1, group)
	stepsBuilder.GetSecurityGroupWithErrorV1Step("Get deleted security group", suite.Client.NetworkV1, *groupWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteNicV1Step("Delete the nic", suite.Client.NetworkV1, nic)
	stepsBuilder.GetNicWithErrorV1Step("Get deleted nic", suite.Client.NetworkV1, *nicWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeletePublicIpV1Step("Delete the public ip", suite.Client.NetworkV1, publicIp)
	stepsBuilder.GetPublicIpWithErrorV1Step("Get deleted public ip", suite.Client.NetworkV1, *publicIpWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteSubnetV1Step("Delete the subnet", suite.Client.NetworkV1, subnet)
	stepsBuilder.GetSubnetWithErrorV1Step("Get deleted subnet", suite.Client.NetworkV1, *subnetNRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", suite.Client.NetworkV1, route)
	stepsBuilder.GetRouteTableWithErrorV1Step("Get deleted route table", suite.Client.NetworkV1, *routeNRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, gateway)
	stepsBuilder.GetInternetGatewayWithErrorV1Step("Get deleted internet gateway", suite.Client.NetworkV1, *gatewayWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", suite.Client.NetworkV1, network)
	stepsBuilder.GetNetworkWithErrorV1Step("Get deleted network", suite.Client.NetworkV1, *networkWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}
