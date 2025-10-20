package secatest

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"

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

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	blockStorageName := secalib.GenerateBlockStorageName()
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
	var workResp *schema.Workspace

	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		ws := &schema.Workspace{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		}
		workResp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workResp)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *workResp.Status.State)
	})

	t.WithNewStep("Get created workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		workResp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workResp)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *workResp.Status.State)
	})

	// Network
	var networkResp *schema.Network
	var expectedNetworkMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedNetworkSpec *schema.NetworkSpec

	t.WithNewStep("Create network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNetwork", workspaceName)

		net := &schema.Network{
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
		networkResp, err = suite.client.NetworkV1.CreateOrUpdateNetwork(ctx, net)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta = secalib.NewRegionalWorkspaceResourceMetadata(networkName, secalib.NetworkProviderV1, networkResource, secalib.ApiVersion1, secalib.NetworkKind,
			suite.tenant, workspaceName, suite.region)
		expectedNetworkMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		expectedNetworkSpec = &schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}
		suite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, &networkResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *networkResp.Status.State)
	})

	t.WithNewStep("Get created network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      networkName,
		}
		networkResp, err = suite.client.NetworkV1.GetNetwork(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		suite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, &networkResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *networkResp.Status.State)
	})

	t.WithNewStep("Update network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNetwork", workspaceName)

		networkResp, err = suite.client.NetworkV1.CreateOrUpdateNetwork(ctx, networkResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		expectedNetworkSpec.SkuRef = *networkSkuRef2Obj
		suite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, &networkResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingStatusState, *networkResp.Status.State)
	})

	t.WithNewStep("Get updated network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      networkName,
		}
		networkResp, err = suite.client.NetworkV1.GetNetwork(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		suite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, &networkResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *networkResp.Status.State)
	})

	// Internet gateway
	var gatewayResp *schema.InternetGateway
	var expectedGatewayMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedGatewaySpec *schema.InternetGatewaySpec

	t.WithNewStep("Create internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateInternetGateway", workspaceName)

		gtw := &schema.InternetGateway{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      internetGatewayName,
			},
		}
		gatewayResp, err = suite.client.NetworkV1.CreateOrUpdateInternetGateway(ctx, gtw)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta = secalib.NewRegionalWorkspaceResourceMetadata(internetGatewayName, secalib.NetworkProviderV1, internetGatewayResource, secalib.ApiVersion1, secalib.InternetGatewayKind,
			suite.tenant, workspaceName, suite.region)
		expectedGatewayMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		expectedGatewaySpec = &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)}
		suite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, &gatewayResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *gatewayResp.Status.State)
	})

	t.WithNewStep("Get created internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      internetGatewayName,
		}
		gatewayResp, err = suite.client.NetworkV1.GetInternetGateway(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		suite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, &gatewayResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *gatewayResp.Status.State)
	})

	t.WithNewStep("Update internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateInternetGateway", workspaceName)

		gatewayResp, err = suite.client.NetworkV1.CreateOrUpdateInternetGateway(ctx, gatewayResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		expectedGatewaySpec.EgressOnly = ptr.To(true)
		suite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, &gatewayResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingStatusState, *gatewayResp.Status.State)
	})

	t.WithNewStep("Get updated internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      internetGatewayName,
		}
		gatewayResp, err = suite.client.NetworkV1.GetInternetGateway(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		suite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, &gatewayResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *gatewayResp.Status.State)
	})

	// Route table
	var routeResp *schema.RouteTable
	var expectedRouteMeta *schema.RegionalNetworkResourceMetadata
	var expectedRouteSpec *schema.RouteTableSpec

	t.WithNewStep("Create route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateRouteTable", workspaceName)

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
		routeResp, err = suite.client.NetworkV1.CreateOrUpdateRouteTable(ctx, route)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta = secalib.NewRegionalNetworkResourceMetadata(routeTableName, secalib.NetworkProviderV1, routeTableResource, secalib.ApiVersion1, secalib.RouteTableKind,
			suite.tenant, workspaceName, networkName, suite.region)
		expectedRouteMeta.Verb = http.MethodPut
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		expectedRouteSpec = &schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}
		suite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, &routeResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *routeResp.Status.State)
	})

	t.WithNewStep("Get created route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", workspaceName)

		nref := secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      routeTableName,
		}
		routeResp, err = suite.client.NetworkV1.GetRouteTable(ctx, nref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta.Verb = http.MethodGet
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		suite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, &routeResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *routeResp.Status.State)
	})

	t.WithNewStep("Update route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateRouteTable", workspaceName)

		routeResp, err = suite.client.NetworkV1.CreateOrUpdateRouteTable(ctx, routeResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta.Verb = http.MethodPut
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		expectedRouteSpec = &schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *instanceRefObj},
			},
		}
		suite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, &routeResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingStatusState, *routeResp.Status.State)
	})

	t.WithNewStep("Get updated route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", workspaceName)

		nref := secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      routeTableName,
		}
		routeResp, err = suite.client.NetworkV1.GetRouteTable(ctx, nref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta.Verb = http.MethodGet
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		suite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, &routeResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *routeResp.Status.State)
	})

	// Subnet
	var subnetResp *schema.Subnet
	var expectedSubnetMeta *schema.RegionalNetworkResourceMetadata
	var expectedSubnetSpec *schema.SubnetSpec

	t.WithNewStep("Create subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSubnet", workspaceName)

		sub := &schema.Subnet{
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
		subnetResp, err = suite.client.NetworkV1.CreateOrUpdateSubnet(ctx, sub)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta = secalib.NewRegionalNetworkResourceMetadata(subnetName, secalib.NetworkProviderV1, subnetResource, secalib.ApiVersion1, secalib.SubnetKind,
			suite.tenant, workspaceName, networkName, suite.region)
		expectedSubnetMeta.Verb = http.MethodPut
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		expectedSubnetSpec = &schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone1,
		}
		suite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, &subnetResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *subnetResp.Status.State)
	})

	t.WithNewStep("Get created subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", workspaceName)

		nref := secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      subnetName,
		}
		subnetResp, err = suite.client.NetworkV1.GetSubnet(ctx, nref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta.Verb = http.MethodGet
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		suite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, &subnetResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *subnetResp.Status.State)
	})

	t.WithNewStep("Update subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSubnet", workspaceName)

		subnetResp, err = suite.client.NetworkV1.CreateOrUpdateSubnet(ctx, subnetResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta.Verb = http.MethodPut
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		expectedSubnetSpec.Zone = zone2
		suite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, &subnetResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingStatusState, *subnetResp.Status.State)
	})

	t.WithNewStep("Get updated subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSubnet", workspaceName)

		nref := secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      subnetName,
		}
		subnetResp, err = suite.client.NetworkV1.GetSubnet(ctx, nref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta.Verb = http.MethodGet
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		suite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, &subnetResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *subnetResp.Status.State)
	})

	// Public ip
	var publicIpResp *schema.PublicIp
	var expectedIpMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedIpSpec *schema.PublicIpSpec

	t.WithNewStep("Create public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdatePublicIp", workspaceName)

		ip := &schema.PublicIp{
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
		publicIpResp, err = suite.client.NetworkV1.CreateOrUpdatePublicIp(ctx, ip)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta = secalib.NewRegionalWorkspaceResourceMetadata(publicIpName, secalib.NetworkProviderV1, publicIpResource, secalib.ApiVersion1, secalib.PublicIpKind,
			suite.tenant, workspaceName, suite.region)
		expectedIpMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		expectedIpSpec = &schema.PublicIpSpec{
			Address: &publicIpAddress1,
			Version: secalib.IpVersion4,
		}
		suite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, &publicIpResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *publicIpResp.Status.State)
	})

	t.WithNewStep("Get created public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIpName,
		}
		publicIpResp, err = suite.client.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		suite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, &publicIpResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *publicIpResp.Status.State)
	})

	t.WithNewStep("Update public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdatePublicIp", workspaceName)

		publicIpResp, err = suite.client.NetworkV1.CreateOrUpdatePublicIp(ctx, publicIpResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		expectedIpSpec.Address = ptr.To(publicIpAddress2)
		suite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, &publicIpResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingStatusState, *publicIpResp.Status.State)
	})

	t.WithNewStep("Get updated public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIpName,
		}
		publicIpResp, err = suite.client.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		suite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, &publicIpResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *publicIpResp.Status.State)
	})

	// Nic
	var nicResp *schema.Nic
	var expectedNicMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedNicSpec *schema.NicSpec

	t.WithNewStep("Create nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNic", workspaceName)

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
		nicResp, err = suite.client.NetworkV1.CreateOrUpdateNic(ctx, nic)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta = secalib.NewRegionalWorkspaceResourceMetadata(nicName, secalib.NetworkProviderV1, nicResource, secalib.ApiVersion1, secalib.NicKind,
			suite.tenant, workspaceName, suite.region)
		expectedNicMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		expectedNicSpec = &schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}
		suite.verifyNicSpecStep(sCtx, expectedNicSpec, &nicResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *nicResp.Status.State)
	})

	t.WithNewStep("Get created nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNic", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nicName,
		}
		nicResp, err = suite.client.NetworkV1.GetNic(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		suite.verifyNicSpecStep(sCtx, expectedNicSpec, &nicResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *nicResp.Status.State)
	})

	t.WithNewStep("Update nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNic", workspaceName)

		nicResp, err = suite.client.NetworkV1.CreateOrUpdateNic(ctx, nicResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		expectedNicSpec.Addresses = []string{nicAddress2}
		suite.verifyNicSpecStep(sCtx, expectedNicSpec, &nicResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingStatusState, *nicResp.Status.State)
	})

	t.WithNewStep("Get updated nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNic", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nicName,
		}
		nicResp, err = suite.client.NetworkV1.GetNic(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		suite.verifyNicSpecStep(sCtx, expectedNicSpec, &nicResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *nicResp.Status.State)
	})

	// Security Group
	var groupResp *schema.SecurityGroup
	var instanceResp *schema.Instance
	var expectedGroupMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedGroupSpec *schema.SecurityGroupSpec

	t.WithNewStep("Create security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSecurityGroup", workspaceName)

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
		groupResp, err = suite.client.NetworkV1.CreateOrUpdateSecurityGroup(ctx, group)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta = secalib.NewRegionalWorkspaceResourceMetadata(securityGroupName, secalib.NetworkProviderV1, securityGroupResource, secalib.ApiVersion1, secalib.SecurityGroupKind,
			suite.tenant, workspaceName, suite.region)
		expectedGroupMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		expectedGroupSpec = &schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{
				{Direction: secalib.SecurityRuleDirectionIngress},
			},
		}
		suite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, &groupResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *groupResp.Status.State)
	})

	t.WithNewStep("Get created security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      securityGroupName,
		}
		groupResp, err = suite.client.NetworkV1.GetSecurityGroup(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		suite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, &groupResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *groupResp.Status.State)
	})

	t.WithNewStep("Update security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSecurityGroup", workspaceName)

		groupResp, err = suite.client.NetworkV1.CreateOrUpdateSecurityGroup(ctx, groupResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		expectedGroupSpec.Rules[0] = schema.SecurityGroupRuleSpec{Direction: secalib.SecurityRuleDirectionEgress}
		suite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, &groupResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingStatusState, *groupResp.Status.State)
	})

	t.WithNewStep("Get updated security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      securityGroupName,
		}
		groupResp, err = suite.client.NetworkV1.GetSecurityGroup(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		suite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, &groupResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *groupResp.Status.State)
	})

	// Block storage
	var blockResp *schema.BlockStorage

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", workspaceName)

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
		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, block)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *blockResp.Status.State)
	})

	t.WithNewStep("Get created block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		blockResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *blockResp.Status.State)
	})

	t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "CreateOrUpdateInstance", workspaceName)

		inst := &schema.Instance{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      instanceName,
			},
			Spec: schema.InstanceSpec{
				SkuRef: *instanceSkuRefObj,
				Zone:   zone1,
			},
		}
		inst.Spec.BootVolume.DeviceRef = *blockStorageRefObj

		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, inst)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *instanceResp.Status.State)
	})

	t.WithNewStep("Get created instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *instanceResp.Status.State)
	})

	t.WithNewStep("Delete instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "DeleteInstance", workspaceName)

		err = suite.client.ComputeV1.DeleteInstance(ctx, instanceResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		_, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteBlockStorage", workspaceName)

		err = suite.client.StorageV1.DeleteBlockStorage(ctx, blockResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		_, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSecurityGroup", workspaceName)

		err = suite.client.NetworkV1.DeleteSecurityGroup(ctx, groupResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      securityGroupName,
		}
		_, err = suite.client.NetworkV1.GetSecurityGroup(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNic", workspaceName)

		err = suite.client.NetworkV1.DeleteNic(ctx, nicResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNic", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nicName,
		}
		_, err = suite.client.NetworkV1.GetNic(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeletePublicIp", workspaceName)

		err = suite.client.NetworkV1.DeletePublicIp(ctx, publicIpResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIpName,
		}
		_, err = suite.client.NetworkV1.GetPublicIp(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSubnet", workspaceName)

		err = suite.client.NetworkV1.DeleteSubnet(ctx, subnetResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSubnet", workspaceName)

		wref := secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      subnetName,
		}
		_, err = suite.client.NetworkV1.GetSubnet(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteRouteTable", workspaceName)

		err = suite.client.NetworkV1.DeleteRouteTable(ctx, routeResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", workspaceName)

		wref := secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      routeTableName,
		}
		_, err = suite.client.NetworkV1.GetRouteTable(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteInternetGateway", workspaceName)

		err = suite.client.NetworkV1.DeleteInternetGateway(ctx, gatewayResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      internetGatewayName,
		}
		_, err = suite.client.NetworkV1.GetInternetGateway(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNetwork", workspaceName)

		err = suite.client.NetworkV1.DeleteNetwork(ctx, networkResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      networkName,
		}
		_, err = suite.client.NetworkV1.GetNetwork(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *NetworkV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
