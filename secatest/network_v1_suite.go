package secatest

import (
	"log/slog"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"k8s.io/utils/ptr"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type NetworkV1TestSuite struct {
	regionalTestSuite

	networkCidr    string
	publicIpsRange string
	regionZones    []string
	storageSkus    []string
	instanceSkus   []string
	networkSkus    []string
}

func (suite *NetworkV1TestSuite) TestSuite(t provider.T) {
	var err error
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.NetworkProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
	)

	// Generate the subnet cidr
	subnetCidr, err := secalib.GenerateSubnetCidr(suite.networkCidr, 8, 1)
	if err != nil {
		t.Fatalf("Failed to generate subnet cidr: %v", err)
	}

	// Generate the nic addresses
	nicAddress1, err := secalib.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}
	nicAddress2, err := secalib.GenerateNicAddress(subnetCidr, 2)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}

	// Generate the public ips
	publicIpAddress1, err := secalib.GeneratePublicIp(suite.publicIpsRange, 1)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}
	publicIpAddress2, err := secalib.GeneratePublicIp(suite.publicIpsRange, 2)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}

	// Select zones
	zone1 := suite.regionZones[rand.Intn(len(suite.regionZones))]
	zone2 := suite.regionZones[rand.Intn(len(suite.regionZones))]

	// Select skus
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	networkSkuName1 := suite.networkSkus[rand.Intn(len(suite.networkSkus))]
	networkSkuName2 := suite.networkSkus[rand.Intn(len(suite.networkSkus))]

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, workspaceName)
	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageResource := secalib.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName)
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)
	blockStorageRefObj, err := secapi.BuildReferenceFromURN(blockStorageRef)
	if err != nil {
		t.Fatal(err)
	}

	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	instanceSkuRefObj, err := secapi.BuildReferenceFromURN(instanceSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	instanceName := secalib.GenerateInstanceName()
	instanceResource := secalib.GenerateInstanceResource(suite.tenant, workspaceName, instanceName)
	instanceRef := secalib.GenerateInstanceRef(instanceName)
	instanceRefObj, err := secapi.BuildReferenceFromURN(instanceRef)
	if err != nil {
		t.Fatal(err)
	}

	networkSkuRef1 := secalib.GenerateSkuRef(networkSkuName1)
	networkSkuRefObj, err := secapi.BuildReferenceFromURN(networkSkuRef1)
	if err != nil {
		t.Fatal(err)
	}
	networkSkuRef2 := secalib.GenerateSkuRef(networkSkuName2)
	networkSkuRef2Obj, err := secapi.BuildReferenceFromURN(networkSkuRef2)
	if err != nil {
		t.Fatal(err)
	}

	networkName := secalib.GenerateNetworkName()
	networkResource := secalib.GenerateNetworkResource(suite.tenant, workspaceName, networkName)

	internetGatewayName := secalib.GenerateInternetGatewayName()
	internetGatewayRef := secalib.GenerateInternetGatewayRef(internetGatewayName)
	internetGatewayRefObj, err := secapi.BuildReferenceFromURN(internetGatewayRef)
	if err != nil {
		t.Fatal(err)
	}
	internetGatewayResource := secalib.GenerateInternetGatewayResource(suite.tenant, workspaceName, internetGatewayName)

	routeTableName := secalib.GenerateRouteTableName()
	routeTableRef := secalib.GenerateRouteTableRef(routeTableName)
	routeTableRefObj, err := secapi.BuildReferenceFromURN(routeTableRef)
	if err != nil {
		t.Fatal(err)
	}
	routeTableResource := secalib.GenerateRouteTableResource(suite.tenant, workspaceName, networkName, routeTableName)

	subnetName := secalib.GenerateSubnetName()
	subnetRef := secalib.GenerateSubnetRef(subnetName)
	subnetRefObj, err := secapi.BuildReferenceFromURN(subnetRef)
	if err != nil {
		t.Fatal(err)
	}
	subnetResource := secalib.GenerateSubnetResource(suite.tenant, workspaceName, networkName, subnetName)

	nicName := secalib.GenerateNicName()
	nicResource := secalib.GenerateNicResource(suite.tenant, workspaceName, nicName)

	publicIpName := secalib.GeneratePublicIpName()
	publicIpRef := secalib.GeneratePublicIpRef(publicIpName)
	publicIpRefObj, err := secapi.BuildReferenceFromURN(publicIpRef)
	if err != nil {
		t.Fatal(err)
	}
	publicIpResource := secalib.GeneratePublicIpResource(suite.tenant, workspaceName, publicIpName)

	securityGroupName := secalib.GenerateSecurityGroupName()
	securityGroupResource := secalib.GenerateSecurityGroupResource(suite.tenant, workspaceName, securityGroupName)

	blockStorageSize := secalib.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.NetworkParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					envLabel: envDevelopmentLabel,
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
					Cidr:          schema.Cidr{Ipv4: ptr.To(suite.networkCidr)},
					SkuRef:        *networkSkuRefObj,
					RouteTableRef: *routeTableRefObj,
				},
				UpdatedSpec: &schema.NetworkSpec{
					Cidr:          schema.Cidr{Ipv4: ptr.To(suite.networkCidr)},
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
						{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
					},
				},
				UpdatedSpec: &schema.RouteTableSpec{
					Routes: []schema.RouteSpec{
						{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *instanceRefObj},
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
		wm, err := mock.ConfigNetworkLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			envLabel: envDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	expectWorkspaceMeta, err := builders.NewRegionalResourceMetadataBuilder().
		Name(workspaceName).
		Provider(secalib.WorkspaceProviderV1).
		Resource(workspaceResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalResourceMetadataKindResourceKindWorkspace).
		Tenant(suite.tenant).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{envLabel: envDevelopmentLabel}
	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, suite.client.WorkspaceV1, workspace,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectWorkspaceLabels,
			metadata:      expectWorkspaceMeta,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	suite.getWorkspaceV1Step("Get the created workspace", t, suite.client.WorkspaceV1, *workspaceTRef,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectWorkspaceLabels,
			metadata:      expectWorkspaceMeta,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Network

	// Create a network
	network := &schema.Network{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      networkName,
		},
		Spec: schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		},
	}
	expectNetworkMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(networkName).
		Provider(secalib.NetworkProviderV1).
		Resource(networkResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectNetworkSpec := &schema.NetworkSpec{
		Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
		SkuRef:        *networkSkuRefObj,
		RouteTableRef: *routeTableRefObj,
	}
	suite.createOrUpdateNetworkV1Step("Create a network", t, suite.client.NetworkV1, network,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			metadata:      expectNetworkMeta,
			spec:          expectNetworkSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created network
	networkWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      networkName,
	}
	suite.getNetworkV1Step("Get the created network", t, suite.client.NetworkV1, *networkWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			metadata:      expectNetworkMeta,
			spec:          expectNetworkSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the network
	network.Spec.SkuRef = *networkSkuRef2Obj
	expectNetworkSpec.SkuRef = network.Spec.SkuRef
	suite.createOrUpdateNetworkV1Step("Update the network", t, suite.client.NetworkV1, network,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			metadata:      expectNetworkMeta,
			spec:          expectNetworkSpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated network
	suite.getNetworkV1Step("Get the updated network", t, suite.client.NetworkV1, *networkWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			metadata:      expectNetworkMeta,
			spec:          expectNetworkSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Internet gateway

	// Create an internet gateway
	gateway := &schema.InternetGateway{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      internetGatewayName,
		},
	}
	expectGatewayMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(internetGatewayName).
		Provider(secalib.NetworkProviderV1).
		Resource(internetGatewayResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectGatewaySpec := &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)}
	suite.createOrUpdateInternetGatewayV1Step("Create a internet gateway", t, suite.client.NetworkV1, gateway,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			metadata:      expectGatewayMeta,
			spec:          expectGatewaySpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created internet gateway
	gatewayWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      internetGatewayName,
	}
	suite.getInternetGatewayV1Step("Get the created internet gateway", t, suite.client.NetworkV1, *gatewayWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			metadata:      expectGatewayMeta,
			spec:          expectGatewaySpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the internet gateway
	gateway.Spec.EgressOnly = ptr.To(true)
	expectGatewaySpec.EgressOnly = gateway.Spec.EgressOnly
	suite.createOrUpdateInternetGatewayV1Step("Update the internet gateway", t, suite.client.NetworkV1, gateway,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			metadata:      expectGatewayMeta,
			spec:          expectGatewaySpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated internet gateway
	suite.getInternetGatewayV1Step("Get the updated internet gateway", t, suite.client.NetworkV1, *gatewayWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			metadata:      expectGatewayMeta,
			spec:          expectGatewaySpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Route table

	// Create a route table
	route := &schema.RouteTable{
		Metadata: &schema.RegionalNetworkResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Network:   networkName,
			Name:      routeTableName,
		},
		Spec: schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		},
	}
	expectRouteMeta, err := builders.NewRegionalNetworkResourceMetadataBuilder().
		Name(routeTableName).
		Provider(secalib.NetworkProviderV1).
		Resource(routeTableResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Network(networkName).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectRouteSpec := &schema.RouteTableSpec{
		Routes: []schema.RouteSpec{
			{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
		},
	}
	suite.createOrUpdateRouteTableV1Step("Create a route table", t, suite.client.NetworkV1, route,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			metadata:      expectRouteMeta,
			spec:          expectRouteSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created route table
	routeNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      routeTableName,
	}
	suite.getRouteTableV1Step("Get the created route table", t, suite.client.NetworkV1, *routeNRef,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			metadata:      expectRouteMeta,
			spec:          expectRouteSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the route table
	route.Spec.Routes = []schema.RouteSpec{
		{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *instanceRefObj},
	}
	expectRouteSpec.Routes = route.Spec.Routes
	suite.createOrUpdateRouteTableV1Step("Update the route table", t, suite.client.NetworkV1, route,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			metadata:      expectRouteMeta,
			spec:          expectRouteSpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated route table
	suite.getRouteTableV1Step("Get the updated route table", t, suite.client.NetworkV1, *routeNRef,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			metadata:      expectRouteMeta,
			spec:          expectRouteSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Subnet

	// Create a subnet
	subnet := &schema.Subnet{
		Metadata: &schema.RegionalNetworkResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Network:   networkName,
			Name:      subnetName,
		},
		Spec: schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone1,
		},
	}
	expectSubnetMeta, err := builders.NewRegionalNetworkResourceMetadataBuilder().
		Name(subnetName).
		Provider(secalib.NetworkProviderV1).
		Resource(subnetResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Network(networkName).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectSubnetSpec := &schema.SubnetSpec{
		Cidr: schema.Cidr{Ipv4: &subnetCidr},
		Zone: zone1,
	}
	suite.createOrUpdateSubnetV1Step("Create a subnet", t, suite.client.NetworkV1, subnet,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			metadata:      expectSubnetMeta,
			spec:          expectSubnetSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created subnet
	subnetNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      subnetName,
	}
	suite.getSubnetV1Step("Get the created subnet", t, suite.client.NetworkV1, *subnetNRef,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			metadata:      expectSubnetMeta,
			spec:          expectSubnetSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the subnet
	subnet.Spec.Zone = zone2
	expectSubnetSpec.Zone = subnet.Spec.Zone
	suite.createOrUpdateSubnetV1Step("Update the subnet", t, suite.client.NetworkV1, subnet,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			metadata:      expectSubnetMeta,
			spec:          expectSubnetSpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated subnet
	suite.getSubnetV1Step("Get the updated subnet", t, suite.client.NetworkV1, *subnetNRef,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			metadata:      expectSubnetMeta,
			spec:          expectSubnetSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Public ip

	// Create a public ip
	publicIp := &schema.PublicIp{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      publicIpName,
		},
		Spec: schema.PublicIpSpec{
			Address: &publicIpAddress1,
			Version: schema.IPVersionIPv4,
		},
	}
	expectPublicIpMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(publicIpName).
		Provider(secalib.NetworkProviderV1).
		Resource(publicIpResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectPublicIpSpec := &schema.PublicIpSpec{
		Address: &publicIpAddress1,
		Version: schema.IPVersionIPv4,
	}
	suite.createOrUpdatePublicIpV1Step("Create a public ip", t, suite.client.NetworkV1, publicIp,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			metadata:      expectPublicIpMeta,
			spec:          expectPublicIpSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created public ip
	publicIpWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      publicIpName,
	}
	suite.getPublicIpV1Step("Get the created public ip", t, suite.client.NetworkV1, *publicIpWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			metadata:      expectPublicIpMeta,
			spec:          expectPublicIpSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the public ip
	publicIp.Spec.Address = ptr.To(publicIpAddress2)
	expectPublicIpSpec.Address = publicIp.Spec.Address
	suite.createOrUpdatePublicIpV1Step("Update the public ip", t, suite.client.NetworkV1, publicIp,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			metadata:      expectPublicIpMeta,
			spec:          expectPublicIpSpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated public ip
	suite.getPublicIpV1Step("Get the updated public ip", t, suite.client.NetworkV1, *publicIpWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			metadata:      expectPublicIpMeta,
			spec:          expectPublicIpSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Nic

	// Create a nic
	nic := &schema.Nic{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      nicName,
		},
		Spec: schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		},
	}
	expectNicMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(nicName).
		Provider(secalib.NetworkProviderV1).
		Resource(nicResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Region(suite.region).
		BuildResponse()
	expectNicSpec := &schema.NicSpec{
		Addresses:    []string{nicAddress1},
		PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
		SubnetRef:    *subnetRefObj,
	}
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	suite.createOrUpdateNicV1Step("Create a nic", t, suite.client.NetworkV1, nic,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			metadata:      expectNicMeta,
			spec:          expectNicSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created nic
	nicWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      nicName,
	}
	suite.getNicV1Step("Get the created nic", t, suite.client.NetworkV1, *nicWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			metadata:      expectNicMeta,
			spec:          expectNicSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the nic
	nic.Spec.Addresses = []string{nicAddress2}
	expectNicSpec.Addresses = nic.Spec.Addresses
	suite.createOrUpdateNicV1Step("Update the nic", t, suite.client.NetworkV1, nic,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			metadata:      expectNicMeta,
			spec:          expectNicSpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated nic
	suite.getNicV1Step("Get the updated nic", t, suite.client.NetworkV1, *nicWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			metadata:      expectNicMeta,
			spec:          expectNicSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Security Group

	// Create a security group
	group := &schema.SecurityGroup{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      securityGroupName,
		},
		Spec: schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{
				{Direction: schema.SecurityGroupRuleDirectionIngress},
			},
		},
	}
	expectGroupMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(securityGroupName).
		Provider(secalib.NetworkProviderV1).
		Resource(securityGroupResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Region(suite.region).
		BuildResponse()
	expectGroupSpec := &schema.SecurityGroupSpec{
		Rules: []schema.SecurityGroupRuleSpec{
			{Direction: schema.SecurityGroupRuleDirectionIngress},
		},
	}
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	suite.createOrUpdateSecurityGroupV1Step("Create a security group", t, suite.client.NetworkV1, group,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			metadata:      expectGroupMeta,
			spec:          expectGroupSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created security group
	groupWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      securityGroupName,
	}
	suite.getSecurityGroupV1Step("Get the created security group", t, suite.client.NetworkV1, *groupWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			metadata:      expectGroupMeta,
			spec:          expectGroupSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the security group
	group.Spec.Rules[0] = schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionEgress}
	expectGroupSpec.Rules = group.Spec.Rules
	suite.createOrUpdateSecurityGroupV1Step("Update the security group", t, suite.client.NetworkV1, group,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			metadata:      expectGroupMeta,
			spec:          expectGroupSpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated security group
	suite.getSecurityGroupV1Step("Get the updated security group", t, suite.client.NetworkV1, *groupWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			metadata:      expectGroupMeta,
			spec:          expectGroupSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Block storage

	// Create a block storage
	block := &schema.BlockStorage{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      blockStorageName,
		},
		Spec: schema.BlockStorageSpec{
			SizeGB: blockStorageSize,
			SkuRef: *storageSkuRefObj,
		},
	}
	expectedBlockMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(blockStorageName).
		Provider(secalib.StorageProviderV1).
		Resource(blockStorageResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Region(suite.region).
		BuildResponse()
	expectedBlockSpec := &schema.BlockStorageSpec{
		SizeGB: blockStorageSize,
		SkuRef: *storageSkuRefObj,
	}
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	suite.createOrUpdateBlockStorageV1Step("Create a block storage", t, suite.client.StorageV1, block,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			metadata:      expectedBlockMeta,
			spec:          expectedBlockSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created block storage
	blockWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	suite.getBlockStorageV1Step("Get the created block storage", t, suite.client.StorageV1, *blockWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			metadata:      expectedBlockMeta,
			spec:          expectedBlockSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Instance

	// Create an instance
	instance := &schema.Instance{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
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
	expectInstanceMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(instanceName).
		Provider(secalib.ComputeProviderV1).
		Resource(instanceResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindInstance).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Region(suite.region).
		BuildResponse()
	expectInstanceSpec := &schema.InstanceSpec{
		SkuRef: *instanceSkuRefObj,
		Zone:   zone1,
		BootVolume: schema.VolumeReference{
			DeviceRef: *blockStorageRefObj,
		},
	}
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	suite.createOrUpdateInstanceV1Step("Create an instance", t, suite.client.ComputeV1, instance,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			metadata:      expectInstanceMeta,
			spec:          expectInstanceSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created instance
	instanceWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      instanceName,
	}
	instance = suite.getInstanceV1Step("Get the created instance", t, suite.client.ComputeV1, *instanceWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			metadata:      expectInstanceMeta,
			spec:          expectInstanceSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Delete the instance
	suite.deleteInstanceV1Step("Delete the instance", t, suite.client.ComputeV1, instance)

	// Get the deleted instance
	suite.getInstanceWithErrorV1Step("Get the deleted instance", t, suite.client.ComputeV1, *instanceWRef, secapi.ErrResourceNotFound)

	// Delete the block storage
	suite.deleteBlockStorageV1Step("Delete the block storage", t, suite.client.StorageV1, block)

	// Get the deleted block storage
	suite.getBlockStorageWithErrorV1Step("Get the deleted block storage", t, suite.client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)

	// Delete the security group
	suite.deleteSecurityGroupV1Step("Delete the security group", t, suite.client.NetworkV1, group)

	// Get deleted security group
	suite.getSecurityGroupWithErrorV1Step("Get deleted security group", t, suite.client.NetworkV1, *groupWRef, secapi.ErrResourceNotFound)

	// Delete the nic
	suite.deleteNicV1Step("Delete the nic", t, suite.client.NetworkV1, nic)

	// Get the deleted nic
	suite.getNicWithErrorV1Step("Get deleted nic", t, suite.client.NetworkV1, *nicWRef, secapi.ErrResourceNotFound)

	// Delete the public ip
	suite.deletePublicIpV1Step("Delete the public ip", t, suite.client.NetworkV1, publicIp)

	// Get the deleted public ip
	suite.getPublicIpWithErrorV1Step("Get deleted public ip", t, suite.client.NetworkV1, *publicIpWRef, secapi.ErrResourceNotFound)

	// Delete the subnet
	suite.deleteSubnetV1Step("Delete the subnet", t, suite.client.NetworkV1, subnet)

	// Get the  deleted subnet
	suite.getSubnetWithErrorV1Step("Get deleted subnet", t, suite.client.NetworkV1, *subnetNRef, secapi.ErrResourceNotFound)

	// Delete the route table
	suite.deleteRouteTableV1Step("Delete the route table", t, suite.client.NetworkV1, route)

	// Get the  deleted route table
	suite.getRouteTableWithErrorV1Step("Get deleted route table", t, suite.client.NetworkV1, *routeNRef, secapi.ErrResourceNotFound)

	// Delete the internet gateway
	suite.deleteInternetGatewayV1Step("Delete the internet gateway", t, suite.client.NetworkV1, gateway)

	// Get the deleted internet gateway
	suite.getInternetGatewayWithErrorV1Step("Get deleted internet gateway", t, suite.client.NetworkV1, *gatewayWRef, secapi.ErrResourceNotFound)

	// Delete the network
	suite.deleteNetworkV1Step("Delete the network", t, suite.client.NetworkV1, network)

	// Get the deleted network
	suite.getNetworkWithErrorV1Step("Get deleted network", t, suite.client.NetworkV1, *networkWRef, secapi.ErrResourceNotFound)

	// Delete the workspace
	suite.deleteWorkspaceV1Step("Delete the workspace", t, suite.client.WorkspaceV1, workspace)

	// Get the deleted workspace
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *NetworkV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
