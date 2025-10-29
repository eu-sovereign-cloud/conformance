package secatest

import (
	"context"
	"log/slog"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"k8s.io/utils/ptr"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

const (
	routeTableDefaultDestination = "0.0.0.0/0"
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
	ctx := context.Background()
	var err error
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.NetworkProviderV1, secalib.NetworkKind, secalib.InternetGatewayKind, secalib.NicKind, secalib.PublicIpKind, secalib.RouteTableKind,
		secalib.SubnetKind, secalib.SecurityGroupKind)

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
					secalib.EnvLabel: secalib.EnvDevelopmentLabel,
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
			NIC: &mock.ResourceParams[schema.NicSpec]{
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
					Version: secalib.IpVersion4,
					Address: ptr.To(publicIpAddress1),
				},
				UpdatedSpec: &schema.PublicIpSpec{
					Version: secalib.IpVersion4,
					Address: ptr.To(publicIpAddress2),
				},
			},
			SecurityGroup: &mock.ResourceParams[schema.SecurityGroupSpec]{
				Name: securityGroupName,
				InitialSpec: &schema.SecurityGroupSpec{
					Rules: []schema.SecurityGroupRuleSpec{{Direction: secalib.SecurityRuleDirectionIngress}},
				},
				UpdatedSpec: &schema.SecurityGroupSpec{
					Rules: []schema.SecurityGroupRuleSpec{{Direction: secalib.SecurityRuleDirectionEgress}},
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
			secalib.EnvLabel: secalib.EnvDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}

	expectMeta := secalib.NewRegionalResourceMetadata(workspaceName,
		secalib.WorkspaceProviderV1,
		workspaceResource,
		secalib.ApiVersion1,
		secalib.WorkspaceKind,
		suite.tenant, suite.region)
	expectLabels := schema.Labels{secalib.EnvLabel: secalib.EnvDevelopmentLabel}
	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, ctx, suite.client.WorkspaceV1, workspace, expectMeta, expectLabels, secalib.CreatingResourceState)

	// Get the created Workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	suite.getWorkspaceV1Step("Get the created workspace", t, ctx, suite.client.WorkspaceV1, *workspaceTRef, expectMeta, expectLabels, secalib.ActiveResourceState)

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
	expectNetworkMeta := secalib.NewRegionalWorkspaceResourceMetadata(networkName,
		secalib.NetworkProviderV1,
		networkResource,
		secalib.ApiVersion1,
		secalib.NetworkKind,
		suite.tenant,
		workspaceName,
		suite.region,
	)
	expectNetworkSpec := &schema.NetworkSpec{
		Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
		SkuRef:        *networkSkuRefObj,
		RouteTableRef: *routeTableRefObj,
	}
	suite.createOrUpdateNetworkV1Step("Create a network", t, ctx, suite.client.NetworkV1, network,
		expectNetworkMeta, expectNetworkSpec, secalib.CreatingResourceState)

	// Get the created network
	networkWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      networkName,
	}
	suite.getNetworkV1Step("Get the created network", t, ctx, suite.client.NetworkV1, *networkWRef,
		expectNetworkMeta, expectNetworkSpec, secalib.ActiveResourceState)

	// Update the network
	network.Spec.SkuRef = *networkSkuRef2Obj
	expectNetworkSpec.SkuRef = network.Spec.SkuRef
	suite.createOrUpdateNetworkV1Step("Update the network", t, ctx, suite.client.NetworkV1, network,
		expectNetworkMeta, expectNetworkSpec, secalib.UpdatingResourceState)

	// Get the updated network
	suite.getNetworkV1Step("Get the updated network", t, ctx, suite.client.NetworkV1, *networkWRef,
		expectNetworkMeta, expectNetworkSpec, secalib.ActiveResourceState)

	// Internet gateway

	// Create an internet gateway
	gateway := &schema.InternetGateway{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      internetGatewayName,
		},
	}
	expectGatewayMeta := secalib.NewRegionalWorkspaceResourceMetadata(internetGatewayName,
		secalib.NetworkProviderV1,
		internetGatewayResource,
		secalib.ApiVersion1,
		secalib.InternetGatewayKind,
		suite.tenant,
		workspaceName,
		suite.region,
	)
	expectGatewaySpec := &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)}
	suite.createOrUpdateInternetGatewayV1Step("Create a internet gateway", t, ctx, suite.client.NetworkV1, gateway,
		expectGatewayMeta, expectGatewaySpec, secalib.CreatingResourceState)

	// Get the created internet gateway
	gatewayWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      internetGatewayName,
	}
	suite.getInternetGatewayV1Step("Get the created internet gateway", t, ctx, suite.client.NetworkV1, *gatewayWRef,
		expectGatewayMeta, expectGatewaySpec, secalib.ActiveResourceState)

	// Update the internet gateway
	gateway.Spec.EgressOnly = ptr.To(true)
	expectGatewaySpec.EgressOnly = gateway.Spec.EgressOnly
	suite.createOrUpdateInternetGatewayV1Step("Update the internet gateway", t, ctx, suite.client.NetworkV1, gateway,
		expectGatewayMeta, expectGatewaySpec, secalib.UpdatingResourceState)

	// Get the updated internet gateway
	suite.getInternetGatewayV1Step("Get the updated internet gateway", t, ctx, suite.client.NetworkV1, *gatewayWRef,
		expectGatewayMeta, expectGatewaySpec, secalib.ActiveResourceState)

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
	expectRouteMeta := secalib.NewRegionalNetworkResourceMetadata(routeTableName,
		secalib.NetworkProviderV1,
		routeTableResource,
		secalib.ApiVersion1,
		secalib.RouteTableKind,
		suite.tenant,
		workspaceName,
		networkName,
		suite.region)
	expectRouteSpec := &schema.RouteTableSpec{
		Routes: []schema.RouteSpec{
			{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
		},
	}
	suite.createOrUpdateRouteTableV1Step("Create a route table", t, ctx, suite.client.NetworkV1, route,
		expectRouteMeta, expectRouteSpec, secalib.CreatingResourceState)

	// Get the created route table
	routeNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      routeTableName,
	}
	suite.getRouteTableV1Step("Get the created route table", t, ctx, suite.client.NetworkV1, *routeNRef,
		expectRouteMeta, expectRouteSpec, secalib.ActiveResourceState)

	// Update the route table
	route.Spec.Routes = []schema.RouteSpec{
		{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *instanceRefObj},
	}
	expectRouteSpec.Routes = route.Spec.Routes
	suite.createOrUpdateRouteTableV1Step("Update the route table", t, ctx, suite.client.NetworkV1, route,
		expectRouteMeta, expectRouteSpec, secalib.UpdatingResourceState)

	// Get the updated route table
	suite.getRouteTableV1Step("Get the updated route table", t, ctx, suite.client.NetworkV1, *routeNRef,
		expectRouteMeta, expectRouteSpec, secalib.ActiveResourceState)

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
	expectSubnetMeta := secalib.NewRegionalNetworkResourceMetadata(subnetName,
		secalib.NetworkProviderV1,
		subnetResource,
		secalib.ApiVersion1,
		secalib.SubnetKind,
		suite.tenant,
		workspaceName,
		networkName,
		suite.region,
	)
	expectSubnetSpec := &schema.SubnetSpec{
		Cidr: schema.Cidr{Ipv4: &subnetCidr},
		Zone: zone1,
	}
	suite.createOrUpdateSubnetV1Step("Create a subnet", t, ctx, suite.client.NetworkV1, subnet,
		expectSubnetMeta, expectSubnetSpec, secalib.CreatingResourceState)

	// Get the created subnet
	subnetNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      subnetName,
	}
	suite.getSubnetV1Step("Get the created subnet", t, ctx, suite.client.NetworkV1, *subnetNRef,
		expectSubnetMeta, expectSubnetSpec, secalib.ActiveResourceState)

	// Update the subnet
	subnet.Spec.Zone = zone2
	expectSubnetSpec.Zone = subnet.Spec.Zone
	suite.createOrUpdateSubnetV1Step("Update the subnet", t, ctx, suite.client.NetworkV1, subnet,
		expectSubnetMeta, expectSubnetSpec, secalib.UpdatingResourceState)

	// Get the updated subnet
	suite.getSubnetV1Step("Get the updated subnet", t, ctx, suite.client.NetworkV1, *subnetNRef,
		expectSubnetMeta, expectSubnetSpec, secalib.ActiveResourceState)

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
			Version: secalib.IpVersion4,
		},
	}
	expectPublicIpMeta := secalib.NewRegionalWorkspaceResourceMetadata(publicIpName,
		secalib.NetworkProviderV1,
		publicIpResource,
		secalib.ApiVersion1,
		secalib.PublicIpKind,
		suite.tenant,
		workspaceName,
		suite.region,
	)
	expectPublicIpSpec := &schema.PublicIpSpec{
		Address: &publicIpAddress1,
		Version: secalib.IpVersion4,
	}
	suite.createOrUpdatePublicIpV1Step("Create a public ip", t, ctx, suite.client.NetworkV1, publicIp,
		expectPublicIpMeta, expectPublicIpSpec, secalib.CreatingResourceState)

	// Get the created public ip
	publicIpWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      publicIpName,
	}
	suite.getPublicIpV1Step("Get the created public ip", t, ctx, suite.client.NetworkV1, *publicIpWRef,
		expectPublicIpMeta, expectPublicIpSpec, secalib.ActiveResourceState)

	// Update the public ip
	publicIp.Spec.Address = ptr.To(publicIpAddress2)
	expectPublicIpSpec.Address = publicIp.Spec.Address
	suite.createOrUpdatePublicIpV1Step("Update the public ip", t, ctx, suite.client.NetworkV1, publicIp,
		expectPublicIpMeta, expectPublicIpSpec, secalib.UpdatingResourceState)

	// Get the updated public ip
	suite.getPublicIpV1Step("Get the updated public ip", t, ctx, suite.client.NetworkV1, *publicIpWRef,
		expectPublicIpMeta, expectPublicIpSpec, secalib.ActiveResourceState)

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
	expectNicMeta := secalib.NewRegionalWorkspaceResourceMetadata(nicName,
		secalib.NetworkProviderV1,
		nicResource,
		secalib.ApiVersion1,
		secalib.NicKind,
		suite.tenant, workspaceName, suite.region)
	expectNicSpec := &schema.NicSpec{
		Addresses:    []string{nicAddress1},
		PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
		SubnetRef:    *subnetRefObj,
	}
	suite.createOrUpdateNicV1Step("Create a nic", t, ctx, suite.client.NetworkV1, nic,
		expectNicMeta, expectNicSpec, secalib.CreatingResourceState)

	// Get the created nic
	nicWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      nicName,
	}
	suite.getNicV1Step("Get the created nic", t, ctx, suite.client.NetworkV1, *nicWRef,
		expectNicMeta, expectNicSpec, secalib.ActiveResourceState)

	// Update the nic
	nic.Spec.Addresses = []string{nicAddress2}
	expectNicSpec.Addresses = nic.Spec.Addresses
	suite.createOrUpdateNicV1Step("Update the nic", t, ctx, suite.client.NetworkV1, nic,
		expectNicMeta, expectNicSpec, secalib.UpdatingResourceState)

	// Get the updated nic
	suite.getNicV1Step("Get the updated nic", t, ctx, suite.client.NetworkV1, *nicWRef,
		expectNicMeta, expectNicSpec, secalib.ActiveResourceState)

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
				{Direction: secalib.SecurityRuleDirectionIngress},
			},
		},
	}
	expectGroupMeta := secalib.NewRegionalWorkspaceResourceMetadata(securityGroupName,
		secalib.NetworkProviderV1,
		securityGroupResource,
		secalib.ApiVersion1,
		secalib.SecurityGroupKind,
		suite.tenant,
		workspaceName,
		suite.region,
	)
	expectGroupSpec := &schema.SecurityGroupSpec{
		Rules: []schema.SecurityGroupRuleSpec{
			{Direction: secalib.SecurityRuleDirectionIngress},
		},
	}
	suite.createOrUpdateSecurityGroupV1Step("Create a security group", t, ctx, suite.client.NetworkV1, group,
		expectGroupMeta, expectGroupSpec, secalib.CreatingResourceState)

	// Get the created security group
	groupWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      securityGroupName,
	}
	suite.getSecurityGroupV1Step("Get the created security group", t, ctx, suite.client.NetworkV1, *groupWRef,
		expectGroupMeta, expectGroupSpec, secalib.ActiveResourceState)

	// Update the security group
	group.Spec.Rules[0] = schema.SecurityGroupRuleSpec{Direction: secalib.SecurityRuleDirectionEgress}
	expectGroupSpec.Rules = group.Spec.Rules
	suite.createOrUpdateSecurityGroupV1Step("Update the security group", t, ctx, suite.client.NetworkV1, group,
		expectGroupMeta, expectGroupSpec, secalib.UpdatingResourceState)

	// Get the updated security group
	suite.getSecurityGroupV1Step("Get the updated security group", t, ctx, suite.client.NetworkV1, *groupWRef,
		expectGroupMeta, expectGroupSpec, secalib.ActiveResourceState)

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
	expectedBlockMeta := secalib.NewRegionalWorkspaceResourceMetadata(blockStorageName,
		secalib.StorageProviderV1,
		blockStorageResource,
		secalib.ApiVersion1,
		secalib.BlockStorageKind,
		suite.tenant,
		workspaceName,
		suite.region)
	expectedBlockSpec := &schema.BlockStorageSpec{
		SizeGB: blockStorageSize,
		SkuRef: *storageSkuRefObj,
	}
	suite.createOrUpdateBlockStorageV1Step("Create a block storage", t, ctx, suite.client.StorageV1, block, expectedBlockMeta, expectedBlockSpec, secalib.CreatingResourceState)

	// Get the created block storage
	blockWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	suite.getBlockStorageV1Step("Get the created block storage", t, ctx, suite.client.StorageV1, *blockWRef, expectedBlockMeta, expectedBlockSpec, secalib.ActiveResourceState)

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

	expectInstanceMeta := secalib.NewRegionalWorkspaceResourceMetadata(instanceName,
		secalib.ComputeProviderV1,
		instanceResource,
		secalib.ApiVersion1,
		secalib.InstanceKind,
		suite.tenant,
		workspaceName,
		suite.region)
	expectInstanceSpec := &schema.InstanceSpec{
		SkuRef: *instanceSkuRefObj,
		Zone:   zone1,
		BootVolume: schema.VolumeReference{
			DeviceRef: *blockStorageRefObj,
		},
	}
	suite.createOrUpdateInstanceV1Step("Create an instance", t, ctx, suite.client.ComputeV1, instance, expectInstanceMeta, expectInstanceSpec, secalib.CreatingResourceState)

	// Get the created instance
	instanceWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      instanceName,
	}
	instance = suite.getInstanceV1Step("Get the created instance", t, ctx, suite.client.ComputeV1, *instanceWRef, expectInstanceMeta, expectInstanceSpec, secalib.ActiveResourceState)

	// Delete the instance
	suite.deleteInstanceV1Step("Delete the instance", t, ctx, suite.client.ComputeV1, instance)

	// Get the deleted instance
	suite.getInstanceWithErrorV1Step("Get the deleted instance", t, ctx, suite.client.ComputeV1, *instanceWRef, secapi.ErrResourceNotFound)

	// Delete the block storage
	suite.deleteBlockStorageV1Step("Delete the block storage", t, ctx, suite.client.StorageV1, block)

	// Get the deleted block storage
	suite.getBlockStorageWithErrorV1Step("Get the deleted block storage", t, ctx, suite.client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)

	// Delete the security group
	suite.deleteSecurityGroupV1Step("Delete the security group", t, ctx, suite.client.NetworkV1, group)

	// Get deleted security group
	suite.getSecurityGroupWithErrorV1Step("Get deleted security group", t, ctx, suite.client.NetworkV1, *groupWRef, secapi.ErrResourceNotFound)

	// Delete the nic
	suite.deleteNicV1Step("Delete the nic", t, ctx, suite.client.NetworkV1, nic)

	// Get the deleted nic
	suite.getNicWithErrorV1Step("Get deleted nic", t, ctx, suite.client.NetworkV1, *nicWRef, secapi.ErrResourceNotFound)

	// Delete the public ip
	suite.deletePublicIpV1Step("Delete the public ip", t, ctx, suite.client.NetworkV1, publicIp)

	// Get the deleted public ip
	suite.getPublicIpWithErrorV1Step("Get deleted public ip", t, ctx, suite.client.NetworkV1, *publicIpWRef, secapi.ErrResourceNotFound)

	// Delete the subnet
	suite.deleteSubnetV1Step("Delete the subnet", t, ctx, suite.client.NetworkV1, subnet)

	// Get the  deleted subnet
	suite.getSubnetWithErrorV1Step("Get deleted subnet", t, ctx, suite.client.NetworkV1, *subnetNRef, secapi.ErrResourceNotFound)

	// Delete the route table
	suite.deleteRouteTableV1Step("Delete the route table", t, ctx, suite.client.NetworkV1, route)

	// Get the  deleted route table
	suite.getRouteTableWithErrorV1Step("Get deleted route table", t, ctx, suite.client.NetworkV1, *routeNRef, secapi.ErrResourceNotFound)

	// Delete the internet gateway
	suite.deleteInternetGatewayV1Step("Delete the internet gateway", t, ctx, suite.client.NetworkV1, gateway)

	// Get the deleted internet gateway
	suite.getInternetGatewayWithErrorV1Step("Get deleted internet gateway", t, ctx, suite.client.NetworkV1, *gatewayWRef, secapi.ErrResourceNotFound)

	// Delete the network
	suite.deleteNetworkV1Step("Delete the network", t, ctx, suite.client.NetworkV1, network)

	// Get the deleted network
	suite.getNetworkWithErrorV1Step("Get deleted network", t, ctx, suite.client.NetworkV1, *networkWRef, secapi.ErrResourceNotFound)

	// Delete the workspace
	suite.deleteWorkspaceV1Step("Delete the workspace", t, ctx, suite.client.WorkspaceV1, workspace)

	// Get the deleted workspace
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, ctx, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *NetworkV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
