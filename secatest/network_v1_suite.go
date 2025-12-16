package secatest

import (
	"log/slog"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"k8s.io/utils/ptr"
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
	configureTags(t, networkProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
	)

	// Generate the subnet cidr
	subnetCidr, err := generators.GenerateSubnetCidr(suite.networkCidr, 8, 1)
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
	zone1 := suite.regionZones[rand.Intn(len(suite.regionZones))]
	zone2 := suite.regionZones[rand.Intn(len(suite.regionZones))]

	// Select skus
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	networkSkuName1 := suite.networkSkus[rand.Intn(len(suite.networkSkus))]
	networkSkuName2 := suite.networkSkus[rand.Intn(len(suite.networkSkus))]

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
	expectWorkspaceMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(workspaceName).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Region(suite.region).
		Build()
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
	expectNetworkMeta, err := builders.NewNetworkMetadataBuilder().
		Name(networkName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		Build()
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
	expectGatewayMeta, err := builders.NewInternetGatewayMetadataBuilder().
		Name(internetGatewayName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		Build()
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
	expectRouteMeta, err := builders.NewRouteTableMetadataBuilder().
		Name(routeTableName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Network(networkName).Region(suite.region).
		Build()
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
	expectSubnetMeta, err := builders.NewSubnetMetadataBuilder().
		Name(subnetName).
		Provider(networkProviderV1).
		ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Network(networkName).Region(suite.region).
		Build()
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
	expectPublicIpMeta, err := builders.NewPublicIpMetadataBuilder().
		Name(publicIpName).
		Provider(networkProviderV1).
		ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		Build()
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
	expectNicMeta, err := builders.NewNicMetadataBuilder().
		Name(nicName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		Build()
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
	expectGroupMeta, err := builders.NewSecurityGroupMetadataBuilder().
		Name(securityGroupName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		Build()
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
	expectedBlockMeta, err := builders.NewBlockStorageMetadataBuilder().
		Name(blockStorageName).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		Build()
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
	expectInstanceMeta, err := builders.NewInstanceMetadataBuilder().
		Name(instanceName).
		Provider(computeProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		Build()
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

	// Resources deletion

	suite.deleteInstanceV1Step("Delete the instance", t, suite.client.ComputeV1, instance)
	suite.getInstanceWithErrorV1Step("Get the deleted instance", t, suite.client.ComputeV1, *instanceWRef, secapi.ErrResourceNotFound)

	suite.deleteBlockStorageV1Step("Delete the block storage", t, suite.client.StorageV1, block)
	suite.getBlockStorageWithErrorV1Step("Get the deleted block storage", t, suite.client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)

	suite.deleteSecurityGroupV1Step("Delete the security group", t, suite.client.NetworkV1, group)
	suite.getSecurityGroupWithErrorV1Step("Get deleted security group", t, suite.client.NetworkV1, *groupWRef, secapi.ErrResourceNotFound)

	suite.deleteNicV1Step("Delete the nic", t, suite.client.NetworkV1, nic)
	suite.getNicWithErrorV1Step("Get deleted nic", t, suite.client.NetworkV1, *nicWRef, secapi.ErrResourceNotFound)

	suite.deletePublicIpV1Step("Delete the public ip", t, suite.client.NetworkV1, publicIp)
	suite.getPublicIpWithErrorV1Step("Get deleted public ip", t, suite.client.NetworkV1, *publicIpWRef, secapi.ErrResourceNotFound)

	suite.deleteSubnetV1Step("Delete the subnet", t, suite.client.NetworkV1, subnet)
	suite.getSubnetWithErrorV1Step("Get deleted subnet", t, suite.client.NetworkV1, *subnetNRef, secapi.ErrResourceNotFound)

	suite.deleteRouteTableV1Step("Delete the route table", t, suite.client.NetworkV1, route)
	suite.getRouteTableWithErrorV1Step("Get deleted route table", t, suite.client.NetworkV1, *routeNRef, secapi.ErrResourceNotFound)

	suite.deleteInternetGatewayV1Step("Delete the internet gateway", t, suite.client.NetworkV1, gateway)
	suite.getInternetGatewayWithErrorV1Step("Get deleted internet gateway", t, suite.client.NetworkV1, *gatewayWRef, secapi.ErrResourceNotFound)

	suite.deleteNetworkV1Step("Delete the network", t, suite.client.NetworkV1, network)
	suite.getNetworkWithErrorV1Step("Get deleted network", t, suite.client.NetworkV1, *networkWRef, secapi.ErrResourceNotFound)

	suite.deleteWorkspaceV1Step("Delete the workspace", t, suite.client.WorkspaceV1, workspace)
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *NetworkV1TestSuite) TestListSuite(t provider.T) {
	var err error
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, networkProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
	)

	// Generate the subnet cidr
	subnetCidr, err := generators.GenerateSubnetCidr(suite.networkCidr, 8, 1)
	if err != nil {
		t.Fatalf("Failed to generate subnet cidr: %v", err)
	}

	// Generate the nic addresses
	nicAddress1, err := generators.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}

	// Generate the public ips
	publicIpAddress1, err := generators.GeneratePublicIp(suite.publicIpsRange, 1)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}

	// Select zones
	zone1 := suite.regionZones[rand.Intn(len(suite.regionZones))]

	// Select skus
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	networkSkuName1 := suite.networkSkus[rand.Intn(len(suite.networkSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRef := generators.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRef := generators.GenerateBlockStorageRef(blockStorageName)
	blockStorageRefObj, err := secapi.BuildReferenceFromURN(blockStorageRef)
	if err != nil {
		t.Fatal(err)
	}

	instanceSkuRef := generators.GenerateSkuRef(instanceSkuName)
	instanceSkuRefObj, err := secapi.BuildReferenceFromURN(instanceSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	instanceName := generators.GenerateInstanceName()

	networkSkuRef1 := generators.GenerateSkuRef(networkSkuName1)
	networkSkuRefObj, err := secapi.BuildReferenceFromURN(networkSkuRef1)
	if err != nil {
		t.Fatal(err)
	}

	networkName := generators.GenerateNetworkName()
	networkName2 := generators.GenerateNetworkName()

	internetGatewayName := generators.GenerateInternetGatewayName()
	internetGatewayName2 := generators.GenerateInternetGatewayName()
	internetGatewayRef := generators.GenerateInternetGatewayRef(internetGatewayName)
	internetGatewayRefObj, err := secapi.BuildReferenceFromURN(internetGatewayRef)
	if err != nil {
		t.Fatal(err)
	}

	routeTableName := generators.GenerateRouteTableName()
	routeTableName2 := generators.GenerateRouteTableName()
	routeTableRef := generators.GenerateRouteTableRef(routeTableName)
	routeTableRefObj, err := secapi.BuildReferenceFromURN(routeTableRef)
	if err != nil {
		t.Fatal(err)
	}

	subnetName := generators.GenerateSubnetName()
	subnetName2 := generators.GenerateSubnetName()
	subnetRef := generators.GenerateSubnetRef(subnetName)
	subnetRefObj, err := secapi.BuildReferenceFromURN(subnetRef)
	if err != nil {
		t.Fatal(err)
	}

	nicName := generators.GenerateNicName()
	nicName2 := generators.GenerateNicName()

	publicIpName := generators.GeneratePublicIpName()
	publicIpName2 := generators.GeneratePublicIpName()
	publicIpRef := generators.GeneratePublicIpRef(publicIpName)
	publicIpRefObj, err := secapi.BuildReferenceFromURN(publicIpRef)
	if err != nil {
		t.Fatal(err)
	}

	securityGroupName := generators.GenerateSecurityGroupName()
	securityGroupName2 := generators.GenerateSecurityGroupName()

	blockStorageSize := generators.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.NetworkListParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					generators.EnvLabel: generators.EnvConformanceLabel,
				},
			},
			BlockStorage: &mock.ResourceParams[schema.BlockStorageSpec]{
				Name: blockStorageName,
				InitialLabels: schema.Labels{
					generators.EnvLabel: generators.EnvConformanceLabel,
				},
				InitialSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: blockStorageSize,
				},
			},
			Instance: &mock.ResourceParams[schema.InstanceSpec]{
				Name: instanceName,
				InitialLabels: schema.Labels{
					generators.EnvLabel: generators.EnvConformanceLabel,
				},
				InitialSpec: &schema.InstanceSpec{
					SkuRef: *instanceSkuRefObj,
					Zone:   zone1,
					BootVolume: schema.VolumeReference{
						DeviceRef: *blockStorageRefObj,
					},
				},
			},
			Network: &[]mock.ResourceParams[schema.NetworkSpec]{
				{
					Name: networkName,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.NetworkSpec{
						Cidr:          schema.Cidr{Ipv4: ptr.To(suite.networkCidr)},
						SkuRef:        *networkSkuRefObj,
						RouteTableRef: *routeTableRefObj,
					},
				},
				{
					Name: networkName2,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.NetworkSpec{
						Cidr:          schema.Cidr{Ipv4: ptr.To(suite.networkCidr)},
						SkuRef:        *networkSkuRefObj,
						RouteTableRef: *routeTableRefObj,
					},
				},
			},
			InternetGateway: &[]mock.ResourceParams[schema.InternetGatewaySpec]{
				{
					Name: internetGatewayName,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)},
				},
				{
					Name: internetGatewayName2,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)},
				},
			},
			RouteTable: &[]mock.ResourceParams[schema.RouteTableSpec]{
				{
					Name: routeTableName,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.RouteTableSpec{
						Routes: []schema.RouteSpec{
							{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
						},
					},
				},
				{
					Name: routeTableName2,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.RouteTableSpec{
						Routes: []schema.RouteSpec{
							{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
						},
					},
				},
			},
			Subnet: &[]mock.ResourceParams[schema.SubnetSpec]{
				{
					Name: subnetName,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.SubnetSpec{
						Cidr: schema.Cidr{Ipv4: &subnetCidr},
						Zone: zone1,
					},
				}, {
					Name: subnetName2,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.SubnetSpec{
						Cidr: schema.Cidr{Ipv4: &subnetCidr},
						Zone: zone1,
					},
				},
			},
			Nic: &[]mock.ResourceParams[schema.NicSpec]{
				{
					Name: nicName,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
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
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.NicSpec{
						Addresses:    []string{nicAddress1},
						PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
						SubnetRef:    *subnetRefObj,
					},
				},
			},
			PublicIp: &[]mock.ResourceParams[schema.PublicIpSpec]{
				{
					Name: publicIpName,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.PublicIpSpec{
						Version: schema.IPVersionIPv4,
						Address: ptr.To(publicIpAddress1),
					},
				},
				{
					Name: publicIpName2,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.PublicIpSpec{
						Version: schema.IPVersionIPv4,
						Address: ptr.To(publicIpAddress1),
					},
				},
			},
			SecurityGroup: &[]mock.ResourceParams[schema.SecurityGroupSpec]{
				{
					Name: securityGroupName,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.SecurityGroupSpec{
						Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
					},
				},
				{
					Name: securityGroupName2,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.SecurityGroupSpec{
						Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
					},
				},
			},
		}
		wm, err := mock.ConfigNetworkListAndFilterScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			generators.EnvLabel: generators.EnvConformanceLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	expectWorkspaceMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(workspaceName).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Region(suite.region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{generators.EnvLabel: generators.EnvConformanceLabel}
	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, suite.client.WorkspaceV1, workspace,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectWorkspaceLabels,
			metadata:      expectWorkspaceMeta,
			resourceState: schema.ResourceStateCreating,
		},
	)
	// Network

	// Create a network
	networks := &[]schema.Network{
		{
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
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      networkName2,
			},
			Spec: schema.NetworkSpec{
				Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
				SkuRef:        *networkSkuRefObj,
				RouteTableRef: *routeTableRefObj,
			},
		},
	}

	for _, network := range *networks {

		expectNetworkMeta, err := builders.NewNetworkMetadataBuilder().
			Name(network.Metadata.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		expectNetworkSpec := &schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}
		suite.createOrUpdateNetworkV1Step("Create a network", t, suite.client.NetworkV1, &network,
			responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
				metadata:      expectNetworkMeta,
				spec:          expectNetworkSpec,
				resourceState: schema.ResourceStateCreating,
			},
		)
	}

	wref := secapi.WorkspaceReference{
		Name:      workspaceName,
		Workspace: secapi.WorkspaceID(workspaceName),
		Tenant:    secapi.TenantID(suite.tenant),
	}

	// List Network
	suite.getListNetworkV1Step("List Network", t, suite.client.NetworkV1, wref, nil)

	// List Network with limit
	suite.getListNetworkV1Step("Get list of Network with limit", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List Network with Label
	suite.getListNetworkV1Step("Get list of Network with label", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List Network with Limit and label
	suite.getListNetworkV1Step("Get list of Network with limit and label", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Internet gateway

	// Create an internet gateway
	gateways := &[]schema.InternetGateway{
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      internetGatewayName,
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      internetGatewayName2,
			},
		},
	}

	for _, gateway := range *gateways {
		expectGatewayMeta, err := builders.NewInternetGatewayMetadataBuilder().
			Name(gateway.Metadata.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		expectGatewaySpec := &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)}
		suite.createOrUpdateInternetGatewayV1Step("Create a internet gateway", t, suite.client.NetworkV1, &gateway,
			responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
				metadata:      expectGatewayMeta,
				spec:          expectGatewaySpec,
				resourceState: schema.ResourceStateCreating,
			},
		)

	}

	// List Internet Gateway
	suite.getListInternetGatewayV1Step("List Internet Gateway", t, suite.client.NetworkV1, wref, nil)

	// List Internet Gateway with limit
	suite.getListInternetGatewayV1Step("Get list of Internet Gateway with limit", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List Internet Gateway with Label
	suite.getListInternetGatewayV1Step("Get list of Internet Gateway with label", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List Internet Gateway with Limit and label
	suite.getListInternetGatewayV1Step("Get list of Internet Gateway with limit and label", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Route table

	// Create a route table
	routes := &[]schema.RouteTable{
		{
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
		},
		{
			Metadata: &schema.RegionalNetworkResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Network:   networkName,
				Name:      routeTableName2,
			},
			Spec: schema.RouteTableSpec{
				Routes: []schema.RouteSpec{
					{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
				},
			},
		},
	}
	for _, route := range *routes {
		expectRouteMeta, err := builders.NewRouteTableMetadataBuilder().
			Name(route.Metadata.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(suite.tenant).Workspace(workspaceName).Network(networkName).Region(suite.region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		expectRouteSpec := &schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}
		suite.createOrUpdateRouteTableV1Step("Create a route table", t, suite.client.NetworkV1, &route,
			responseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
				metadata:      expectRouteMeta,
				spec:          expectRouteSpec,
				resourceState: schema.ResourceStateCreating,
			},
		)
	}
	// Get the created route table
	nref := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      routeTableName,
	}
	// List Route table
	suite.getListRouteTableV1Step("List Route table", t, suite.client.NetworkV1, *nref, nil)

	// List Route table with limit
	suite.getListRouteTableV1Step("Get list of Route table with limit", t, suite.client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1))

	// List Route table with Label
	suite.getListRouteTableV1Step("Get list of Route table with label", t, suite.client.NetworkV1, *nref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List Route table with Limit and label
	suite.getListRouteTableV1Step("Get list of Route table with limit and label", t, suite.client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))
	// Subnet

	// Create a subnet
	subnets := &[]schema.Subnet{
		{
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
		},
		{
			Metadata: &schema.RegionalNetworkResourceMetadata{
				Tenant:    suite.tenant,
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
			Provider(networkProviderV1).
			ApiVersion(apiVersion1).
			Tenant(suite.tenant).Workspace(workspaceName).Network(networkName).Region(suite.region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		expectSubnetSpec := &schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone1,
		}
		suite.createOrUpdateSubnetV1Step("Create a subnet", t, suite.client.NetworkV1, &subnet,
			responseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
				metadata:      expectSubnetMeta,
				spec:          expectSubnetSpec,
				resourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List Subnet
	suite.getListSubnetV1Step("List Subnet", t, suite.client.NetworkV1, *nref, nil)

	// List Subnet with limit
	suite.getListSubnetV1Step("Get list of Subnet with limit", t, suite.client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1))

	// List Subnet with Label
	suite.getListSubnetV1Step("Get list of Subnet with label", t, suite.client.NetworkV1, *nref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List Subnet with Limit and label
	suite.getListSubnetV1Step("Get list of Subnet with limit and label", t, suite.client.NetworkV1, *nref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Public ip

	// Create a public ip
	publicIps := &[]schema.PublicIp{
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
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
				Tenant:    suite.tenant,
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
			Provider(networkProviderV1).
			ApiVersion(apiVersion1).
			Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		expectPublicIpSpec := &schema.PublicIpSpec{
			Address: &publicIpAddress1,
			Version: schema.IPVersionIPv4,
		}
		suite.createOrUpdatePublicIpV1Step("Create a public ip", t, suite.client.NetworkV1, &publicIp,
			responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
				metadata:      expectPublicIpMeta,
				spec:          expectPublicIpSpec,
				resourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List PublicIP
	suite.getListPublicIpV1Step("List PublicIP", t, suite.client.NetworkV1, wref, nil)

	// List PublicIP with limit
	suite.getListPublicIpV1Step("Get list of PublicIP with limit", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List PublicIP with Label
	suite.getListPublicIpV1Step("Get list of PublicIP with label", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List PublicIP with Limit and label
	suite.getListPublicIpV1Step("Get list of PublicIP with limit and label", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Nic

	// Create a nic
	nics := &[]schema.Nic{
		{
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
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
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
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
			Build()
		expectNicSpec := &schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		suite.createOrUpdateNicV1Step("Create a nic", t, suite.client.NetworkV1, &nic,
			responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
				metadata:      expectNicMeta,
				spec:          expectNicSpec,
				resourceState: schema.ResourceStateCreating,
			},
		)
	}
	// List Nic
	suite.getListNicV1Step("List Nic", t, suite.client.NetworkV1, wref, nil)

	// List Nic with limit
	suite.getListNicV1Step("Get list of Nic with limit", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List Nic with Label
	suite.getListNicV1Step("Get list of Nic with label", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List Nic with Limit and label
	suite.getListNicV1Step("Get list of Nic with limit and label", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Security Group

	// Create a security group
	groups := &[]schema.SecurityGroup{
		{
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
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
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
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
			Build()
		expectGroupSpec := &schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{
				{Direction: schema.SecurityGroupRuleDirectionIngress},
			},
		}
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		suite.createOrUpdateSecurityGroupV1Step("Create a security group", t, suite.client.NetworkV1, &group,
			responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
				metadata:      expectGroupMeta,
				spec:          expectGroupSpec,
				resourceState: schema.ResourceStateCreating,
			},
		)
	}
	// List Security Group
	suite.getListSecurityGroupV1Step("List Security Group", t, suite.client.NetworkV1, wref, nil)

	// List Security Group with limit
	suite.getListSecurityGroupV1Step("Get list of Security Group with limit", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List Security Group with Label
	suite.getListSecurityGroupV1Step("Get list of Security Group with label", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List Nic with Limit and label
	suite.getListSecurityGroupV1Step("Get list of Security Group with limit and label", t, suite.client.NetworkV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Delete all security groups
	for _, group := range *groups {
		suite.deleteSecurityGroupV1Step("Delete the security group", t, suite.client.NetworkV1, &group)

		// Get deleted security group
		groupWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      group.Metadata.Name,
		}
		suite.getSecurityGroupWithErrorV1Step("Get deleted security group", t, suite.client.NetworkV1, *groupWRef, secapi.ErrResourceNotFound)
	}

	// Delete all NICs
	for _, nic := range *nics {
		suite.deleteNicV1Step("Delete the nic", t, suite.client.NetworkV1, &nic)

		// Get the deleted nic
		nicWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nic.Metadata.Name,
		}
		suite.getNicWithErrorV1Step("Get deleted nic", t, suite.client.NetworkV1, *nicWRef, secapi.ErrResourceNotFound)
	}

	// Delete all public IPs
	for _, publicIp := range *publicIps {
		suite.deletePublicIpV1Step("Delete the public ip", t, suite.client.NetworkV1, &publicIp)

		// Get the deleted public ip
		publicIpWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIp.Metadata.Name,
		}
		suite.getPublicIpWithErrorV1Step("Get deleted public ip", t, suite.client.NetworkV1, *publicIpWRef, secapi.ErrResourceNotFound)
	}

	// Delete all subnets
	for _, subnet := range *subnets {
		suite.deleteSubnetV1Step("Delete the subnet", t, suite.client.NetworkV1, &subnet)

		// Get the deleted subnet
		subnetNRef := &secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      subnet.Metadata.Name,
		}
		suite.getSubnetWithErrorV1Step("Get deleted subnet", t, suite.client.NetworkV1, *subnetNRef, secapi.ErrResourceNotFound)
	}

	// Delete all route tables
	for _, route := range *routes {
		suite.deleteRouteTableV1Step("Delete the route table", t, suite.client.NetworkV1, &route)

		// Get the deleted route table
		routeNRef := &secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      route.Metadata.Name,
		}
		suite.getRouteTableWithErrorV1Step("Get deleted route table", t, suite.client.NetworkV1, *routeNRef, secapi.ErrResourceNotFound)
	}

	// Delete all internet gateways
	for _, gateway := range *gateways {
		suite.deleteInternetGatewayV1Step("Delete the internet gateway", t, suite.client.NetworkV1, &gateway)

		// Get the deleted internet gateway
		gatewayWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      gateway.Metadata.Name,
		}
		suite.getInternetGatewayWithErrorV1Step("Get deleted internet gateway", t, suite.client.NetworkV1, *gatewayWRef, secapi.ErrResourceNotFound)
	}

	// Delete all networks
	for _, network := range *networks {
		suite.deleteNetworkV1Step("Delete the network", t, suite.client.NetworkV1, &network)

		// Get the deleted network
		networkWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      network.Metadata.Name,
		}
		suite.getNetworkWithErrorV1Step("Get deleted network", t, suite.client.NetworkV1, *networkWRef, secapi.ErrResourceNotFound)
	}

	// Delete the workspace
	suite.deleteWorkspaceV1Step("Delete the workspace", t, suite.client.WorkspaceV1, workspace)

	// Get the deleted workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *NetworkV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
