package secatest

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

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

func (suite *NetworkV1TestSuite) TestNetworkV1(t provider.T) {
	slog.Info("Starting Network Lifecycle Test")

	t.Title("Network Lifecycle Test")
	configureTags(t, secalib.NetworkProviderV1, secalib.NetworkKind, secalib.InternetGatewayKind, secalib.NicKind, secalib.PublicIPKind, secalib.RouteTableKind,
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

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)

	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	instanceName := secalib.GenerateInstanceName()
	instanceRef := secalib.GenerateInstanceRef(instanceName)

	networkSkuRef1 := secalib.GenerateSkuRef(networkSkuName1)
	networkSkuRef2 := secalib.GenerateSkuRef(networkSkuName2)

	networkName := secalib.GenerateNetworkName()
	networkRef := secalib.GenerateNetworkRef(networkName)
	networkResource := secalib.GenerateNetworkResource(suite.tenant, workspaceName, networkName)

	internetGatewayName := secalib.GenerateInternetGatewayName()
	internetGatewayRef := secalib.GenerateInternetGatewayRef(internetGatewayName)
	internetGatewayResource := secalib.GenerateInternetGatewayResource(suite.tenant, workspaceName, internetGatewayName)

	routeTableName := secalib.GenerateRouteTableName()
	routeTableRef := secalib.GenerateRouteTableRef(routeTableName)
	routeTableResource := secalib.GenerateRouteTableResource(suite.tenant, workspaceName, networkName, routeTableName)

	subnetName := secalib.GenerateSubnetName()
	subnetRef := secalib.GenerateSubnetRef(subnetName)
	subnetResource := secalib.GenerateSubnetResource(suite.tenant, workspaceName, networkName, subnetName)

	nicName := secalib.GenerateNicName()
	nicResource := secalib.GenerateNicResource(suite.tenant, workspaceName, nicName)

	publicIPName := secalib.GeneratePublicIPName()
	publicIPRef := secalib.GeneratePublicIPRef(publicIPName)
	publicIPResource := secalib.GeneratePublicIPResource(suite.tenant, workspaceName, publicIPName)

	securityGroupName := secalib.GenerateSecurityGroupName()
	securityGroupResource := secalib.GenerateSecurityGroupResource(suite.tenant, workspaceName, securityGroupName)

	blockStorageSize := secalib.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		wm, err := mock.CreateNetworkLifecycleScenarioV1("Network Lifecycle",
			mock.NetworkParamsV1{
				Params: &mock.Params{
					MockURL:   *suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
					Region:    suite.region,
				},
				Workspace: &mock.ResourceParams[secalib.WorkspaceSpecV1]{
					Name: workspaceName,
					InitialSpec: &secalib.WorkspaceSpecV1{
						Labels: &[]secalib.Label{
							{
								Name:  secalib.EnvLabel,
								Value: secalib.EnvDevelopmentLabel,
							},
						},
					},
				},
				BlockStorage: &mock.ResourceParams[secalib.BlockStorageSpecV1]{
					Name: blockStorageName,
					InitialSpec: &secalib.BlockStorageSpecV1{
						SkuRef: storageSkuRef,
						SizeGB: blockStorageSize,
					},
				},
				Instance: &mock.ResourceParams[secalib.InstanceSpecV1]{
					Name: instanceName,
					InitialSpec: &secalib.InstanceSpecV1{
						SkuRef:        instanceSkuRef,
						Zone:          zone1,
						BootDeviceRef: blockStorageRef,
					},
				},
				Network: &mock.ResourceParams[secalib.NetworkSpecV1]{
					Name: networkName,
					InitialSpec: &secalib.NetworkSpecV1{
						Cidr:          &secalib.NetworkSpecCIDRV1{Ipv4: suite.networkCidr},
						SkuRef:        networkSkuRef1,
						RouteTableRef: routeTableRef,
					},
					UpdatedSpec: &secalib.NetworkSpecV1{
						Cidr:          &secalib.NetworkSpecCIDRV1{Ipv4: suite.networkCidr},
						SkuRef:        networkSkuRef2,
						RouteTableRef: routeTableRef,
					},
				},
				InternetGateway: &mock.ResourceParams[secalib.InternetGatewaySpecV1]{
					Name:        internetGatewayName,
					InitialSpec: &secalib.InternetGatewaySpecV1{EgressOnly: false},
					UpdatedSpec: &secalib.InternetGatewaySpecV1{EgressOnly: true},
				},
				RouteTable: &mock.ResourceParams[secalib.RouteTableSpecV1]{
					Name: routeTableName,
					InitialSpec: &secalib.RouteTableSpecV1{
						LocalRef: networkRef,
						Routes: []*secalib.RouteTableRouteV1{
							{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: internetGatewayRef},
						},
					},
					UpdatedSpec: &secalib.RouteTableSpecV1{
						LocalRef: networkRef,
						Routes: []*secalib.RouteTableRouteV1{
							{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: instanceRef},
						},
					},
				},
				Subnet: &mock.ResourceParams[secalib.SubnetSpecV1]{
					Name: subnetName,
					InitialSpec: &secalib.SubnetSpecV1{
						Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCidr},
						Zone: zone1,
					},
					UpdatedSpec: &secalib.SubnetSpecV1{
						Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCidr},
						Zone: zone2,
					},
				},
				NIC: &mock.ResourceParams[secalib.NICSpecV1]{
					Name: nicName,
					InitialSpec: &secalib.NICSpecV1{
						Addresses:    []string{nicAddress1},
						PublicIpRefs: []string{publicIPRef},
						SubnetRef:    subnetRef,
					},
					UpdatedSpec: &secalib.NICSpecV1{
						Addresses:    []string{nicAddress2},
						PublicIpRefs: []string{publicIPRef},
						SubnetRef:    subnetRef,
					},
				},
				PublicIP: &mock.ResourceParams[secalib.PublicIpSpecV1]{
					Name: publicIPName,
					InitialSpec: &secalib.PublicIpSpecV1{
						Version: secalib.IPVersion4,
						Address: publicIpAddress1,
					},
					UpdatedSpec: &secalib.PublicIpSpecV1{
						Version: secalib.IPVersion4,
						Address: publicIpAddress2,
					},
				},
				SecurityGroup: &mock.ResourceParams[secalib.SecurityGroupSpecV1]{
					Name: securityGroupName,
					InitialSpec: &secalib.SecurityGroupSpecV1{
						Rules: []*secalib.SecurityGroupRuleV1{{Direction: secalib.SecurityRuleDirectionIngress}},
					},
					UpdatedSpec: &secalib.SecurityGroupSpecV1{
						Rules: []*secalib.SecurityGroupRuleV1{{Direction: secalib.SecurityRuleDirectionEgress}},
					},
				},
			})
		if err != nil {
			t.Fatalf("Failed to create network scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var workResp *schema.Workspace
	var networkResp *schema.Network
	var gatewayResp *schema.InternetGateway
	var routeResp *schema.RouteTable
	var subnetResp *schema.Subnet
	var publicIpResp *schema.PublicIp
	var nicResp *schema.Nic
	var groupResp *schema.SecurityGroup
	var blockResp *schema.BlockStorage
	var instanceResp *schema.Instance

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

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*workResp.Status.State)},
		)
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

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*workResp.Status.State)},
		)
	})

	// Network
	var expectedNetworkMeta *secalib.Metadata
	var expectedNetworkSpec *secalib.NetworkSpecV1

	t.WithNewStep("Create network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNetwork", workspaceName)

		networkSkuURN, err := suite.client.NetworkV1.BuildReferenceURN(networkSkuRef1)
		if err != nil {
			t.Fatal(err)
		}

		routeTableURN, err := suite.client.NetworkV1.BuildReferenceURN(routeTableRef)
		if err != nil {
			t.Fatal(err)
		}

		net := &schema.Network{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      networkName,
			},
			Spec: schema.NetworkSpec{
				Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
				SkuRef:        *networkSkuURN,
				RouteTableRef: *routeTableURN,
			},
		}
		networkResp, err = suite.client.NetworkV1.CreateOrUpdateNetwork(ctx, net)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta = &secalib.Metadata{
			Name:       networkName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   networkResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NetworkKind,
			Tenant:     suite.tenant,
			Workspace:  &workspaceName,
			Region:     &suite.region,
		}
		suite.verifyWorkspaceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		expectedNetworkSpec = &secalib.NetworkSpecV1{
			Cidr:          &secalib.NetworkSpecCIDRV1{Ipv4: suite.networkCidr},
			SkuRef:        networkSkuRef1,
			RouteTableRef: routeTableRef,
		}
		suite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
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
		suite.verifyWorkspaceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		suite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	t.WithNewStep("Update network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNetwork", workspaceName)

		networkResp, err = suite.client.NetworkV1.CreateOrUpdateNetwork(ctx, networkResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta.Verb = http.MethodPut
		suite.verifyWorkspaceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		expectedNetworkSpec.SkuRef = networkSkuRef2
		suite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
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
		suite.verifyWorkspaceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		suite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	// Internet gateway
	var expectedGatewayMeta *secalib.Metadata
	var expectedGatewaySpec *secalib.InternetGatewaySpecV1

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

		expectedGatewayMeta = &secalib.Metadata{
			Name:       internetGatewayName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   internetGatewayResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InternetGatewayKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.verifyWorkspaceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		expectedGatewaySpec = &secalib.InternetGatewaySpecV1{
			EgressOnly: false,
		}
		suite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
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
		suite.verifyWorkspaceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		suite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	t.WithNewStep("Update internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateInternetGateway", workspaceName)

		gatewayResp, err = suite.client.NetworkV1.CreateOrUpdateInternetGateway(ctx, gatewayResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta.Verb = http.MethodPut
		suite.verifyWorkspaceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		expectedGatewaySpec.EgressOnly = true
		suite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
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
		suite.verifyWorkspaceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		suite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	// Route table
	var expectedRouteMeta *secalib.Metadata
	var expectedRouteSpec *secalib.RouteTableSpecV1

	t.WithNewStep("Create route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateRouteTable", workspaceName)

		route := &schema.RouteTable{
			Metadata: &schema.RegionalNetworkResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Network:   networkName,
				Name:      routeTableName,
			},
		}
		routeResp, err = suite.client.NetworkV1.CreateOrUpdateRouteTable(ctx, route)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta = &secalib.Metadata{
			Name:       routeTableName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   routeTableResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RouteTableKind,
			Tenant:     suite.tenant,
			Workspace:  &workspaceName,
			Region:     &suite.region,
		}
		suite.verifyNetworkMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		expectedRouteSpec = &secalib.RouteTableSpecV1{
			LocalRef: networkRef,
			Routes: []*secalib.RouteTableRouteV1{
				{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: internetGatewayRef},
			},
		}
		suite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
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
		suite.verifyNetworkMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		suite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	t.WithNewStep("Update route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateRouteTable", workspaceName)

		routeResp, err = suite.client.NetworkV1.CreateOrUpdateRouteTable(ctx, routeResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta.Verb = http.MethodPut
		suite.verifyNetworkMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		expectedRouteSpec.Routes = []*secalib.RouteTableRouteV1{{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: instanceRef}}
		suite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
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
		suite.verifyNetworkMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		suite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	// Subnet
	var expectedSubnetMeta *secalib.Metadata
	var expectedSubnetSpec *secalib.SubnetSpecV1

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

		expectedSubnetMeta = &secalib.Metadata{
			Name:       subnetName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   subnetResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SubnetKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.verifyNetworkMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		expectedSubnetSpec = &secalib.SubnetSpecV1{
			Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCidr},
			Zone: zone1,
		}
		suite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
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
		suite.verifyNetworkMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		suite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	t.WithNewStep("Update subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSubnet", workspaceName)

		subnetResp, err = suite.client.NetworkV1.CreateOrUpdateSubnet(ctx, subnetResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta.Verb = http.MethodPut
		suite.verifyNetworkMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		expectedSubnetSpec.Zone = zone2
		suite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
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
		suite.verifyNetworkMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		suite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	// Public ip
	var expectedIpMeta *secalib.Metadata
	var expectedIpSpec *secalib.PublicIpSpecV1

	t.WithNewStep("Create public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdatePublicIp", workspaceName)

		ip := &schema.PublicIp{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      publicIPName,
			},
			Spec: schema.PublicIpSpec{
				Address: &publicIpAddress1,
				Version: secalib.IPVersion4,
			},
		}
		publicIpResp, err = suite.client.NetworkV1.CreateOrUpdatePublicIp(ctx, ip)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta = &secalib.Metadata{
			Name:       publicIPName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   publicIPResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.PublicIPKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.verifyWorkspaceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		expectedIpSpec = &secalib.PublicIpSpecV1{
			Version: secalib.IPVersion4,
			Address: publicIpAddress1,
		}
		suite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	t.WithNewStep("Get created public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIP", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIPName,
		}
		publicIpResp, err = suite.client.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodGet
		suite.verifyWorkspaceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		suite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	t.WithNewStep("Update public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdatePublicIp", workspaceName)

		publicIpResp, err = suite.client.NetworkV1.CreateOrUpdatePublicIp(ctx, publicIpResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodPut
		suite.verifyWorkspaceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		expectedIpSpec.Address = publicIpAddress2
		suite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIPName,
		}
		publicIpResp, err = suite.client.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodGet
		suite.verifyWorkspaceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		suite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	// Nic
	var expectedNicMeta *secalib.Metadata
	var expectedNicSpec *secalib.NICSpecV1

	t.WithNewStep("Create nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNic", workspaceName)

		publicIPURN, err := suite.client.NetworkV1.BuildReferenceURN(publicIPRef)
		if err != nil {
			t.Fatal(err)
		}

		subnetURN, err := suite.client.NetworkV1.BuildReferenceURN(subnetRef)
		if err != nil {
			t.Fatal(err)
		}

		nic := &schema.Nic{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      nicName,
			},
			Spec: schema.NicSpec{
				Addresses:    []string{nicAddress1},
				PublicIpRefs: &[]schema.Reference{*publicIPURN},
				SubnetRef:    *subnetURN,
			},
		}
		nicResp, err = suite.client.NetworkV1.CreateOrUpdateNic(ctx, nic)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta = &secalib.Metadata{
			Name:       nicName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   nicResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NicKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.verifyWorkspaceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		expectedNicSpec = &secalib.NICSpecV1{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: []string{publicIPRef},
			SubnetRef:    subnetRef,
		}
		suite.verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
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
		suite.verifyWorkspaceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		suite.verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	t.WithNewStep("Update nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNic", workspaceName)

		nicResp, err = suite.client.NetworkV1.CreateOrUpdateNic(ctx, nicResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta.Verb = http.MethodPut
		suite.verifyWorkspaceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		expectedNicSpec.Addresses = []string{nicAddress2}
		suite.verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
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
		suite.verifyWorkspaceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		suite.verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	// Security Group
	var expectedGroupMeta *secalib.Metadata
	var expectedGroupSpec *secalib.SecurityGroupSpecV1

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
					{
						Direction: secalib.SecurityRuleDirectionIngress,
					},
				},
			},
		}
		groupResp, err = suite.client.NetworkV1.CreateOrUpdateSecurityGroup(ctx, group)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta = &secalib.Metadata{
			Name:       securityGroupName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   securityGroupResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SecurityGroupKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.verifyWorkspaceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		expectedGroupSpec = &secalib.SecurityGroupSpecV1{
			Rules: []*secalib.SecurityGroupRuleV1{{Direction: secalib.SecurityRuleDirectionIngress}},
		}
		suite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
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
		suite.verifyWorkspaceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		suite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
	})

	t.WithNewStep("Update security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSecurityGroup", workspaceName)

		groupResp, err = suite.client.NetworkV1.CreateOrUpdateSecurityGroup(ctx, groupResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta.Verb = http.MethodPut
		suite.verifyWorkspaceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		expectedGroupSpec.Rules[0] = &secalib.SecurityGroupRuleV1{Direction: secalib.SecurityRuleDirectionEgress}
		suite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
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
		suite.verifyWorkspaceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		suite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
	})

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", workspaceName)

		storageSkuURN, err := suite.client.StorageV1.BuildReferenceURN(storageSkuRef)
		if err != nil {
			t.Fatal(err)
		}

		block := &schema.BlockStorage{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      blockStorageName,
			},
			Spec: schema.BlockStorageSpec{
				SizeGB: blockStorageSize,
				SkuRef: *storageSkuURN,
			},
		}
		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, block)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
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

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "CreateOrUpdateInstance", workspaceName)

		instanceSkuURN, err := suite.client.ComputeV1.BuildReferenceURN(instanceSkuRef)
		if err != nil {
			t.Fatal(err)
		}

		blockStorageURN, err := suite.client.ComputeV1.BuildReferenceURN(blockStorageRef)
		if err != nil {
			t.Fatal(err)
		}

		inst := &schema.Instance{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      instanceName,
			},
			Spec: schema.InstanceSpec{
				SkuRef: *instanceSkuURN,
				Zone:   zone1,
			},
		}
		inst.Spec.BootVolume.DeviceRef = *blockStorageURN

		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, inst)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
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

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
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
			Name:      publicIPName,
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

	slog.Info("Finishing Network Lifecycle Test")
}

func (suite *NetworkV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}

func (suite *NetworkV1TestSuite) verifyNetworkMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *schema.RegionalNetworkResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Workspace:  &metadata.Workspace,
		Network:    &metadata.Network,
		Region:     &metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func (suite *NetworkV1TestSuite) verifyWorkspaceMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *schema.RegionalWorkspaceResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Workspace:  &metadata.Workspace,
		Region:     &metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func (suite *NetworkV1TestSuite) verifyNetworkSpecStep(ctx provider.StepCtx, expected *secalib.NetworkSpecV1, actual schema.NetworkSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		if actual.Cidr.Ipv4 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv4, *actual.Cidr.Ipv4, "Cidr.Ipv4 should match expected")
		}
		if actual.Cidr.Ipv6 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv6, *actual.Cidr.Ipv6, "Cidr.Ipv6 should match expected")
		}

		skuRef, err := asNetworkReferenceURN(actual.SkuRef)
		if err != nil {
			ctx.Error(err)
		}
		stepCtx.Require().Equal(expected.SkuRef, skuRef, "SkuRef should match expected")

		// TODO Convert this to equals string/Reference function
		routeTableRef, err := asNetworkReferenceURN(actual.RouteTableRef)
		if err != nil {
			ctx.Error(err)
		}
		stepCtx.Require().Equal(expected.RouteTableRef, routeTableRef, "RouteTableRef should match expected")
	})
}

func (suite *NetworkV1TestSuite) verifyInternetGatewaySpecStep(ctx provider.StepCtx, expected *secalib.InternetGatewaySpecV1, actual schema.InternetGatewaySpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		if actual.EgressOnly != nil {
			stepCtx.Require().Equal(expected.EgressOnly, *actual.EgressOnly, "EgressOnly should match expected")
		}
	})
}

func (suite *NetworkV1TestSuite) verifyRouteTableSpecStep(ctx provider.StepCtx, expected *secalib.RouteTableSpecV1, actual schema.RouteTableSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(len(expected.Routes), len(actual.Routes), "Route list length should match expected")
		for i := 0; i < len(expected.Routes); i++ {
			expectedRoute := expected.Routes[i]
			actualRoute := actual.Routes[i]
			stepCtx.Require().Equal(expectedRoute.DestinationCidrBlock, actualRoute.DestinationCidrBlock, fmt.Sprintf("Route [%d] DestinationCidrBlock should match expected", i))

			targetRef, err := asNetworkReferenceURN(actualRoute.TargetRef)
			if err != nil {
				ctx.Error(err)
			}
			stepCtx.Require().Equal(expectedRoute.TargetRef, targetRef, fmt.Sprintf("Route [%d] TargetRef should match expected", i))
		}
	})
}

func (suite *NetworkV1TestSuite) verifySubNetSpecStep(ctx provider.StepCtx, expected *secalib.SubnetSpecV1, actual schema.SubnetSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		if actual.Cidr.Ipv4 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv4, *actual.Cidr.Ipv4, "Cidr.Ipv4 should match expected")
		}
		if actual.Cidr.Ipv6 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv6, *actual.Cidr.Ipv6, "Cidr.Ipv6 should match expected")
		}
		stepCtx.Require().Equal(expected.Zone, actual.Zone, "Zone should match expected")
	})
}

func (suite *NetworkV1TestSuite) verifyPublicIpSpecStep(ctx provider.StepCtx, expected *secalib.PublicIpSpecV1, actual schema.PublicIpSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Version, string(actual.Version), "Version should match expected")
		if actual.Address != nil {
			stepCtx.Require().Equal(expected.Address, *actual.Address, "Address should match expected")
		}
	})
}

func (suite *NetworkV1TestSuite) verifyNicSpecStep(ctx provider.StepCtx, expected *secalib.NICSpecV1, actual schema.NicSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Addresses, actual.Addresses, "Addresses should match expected")
		if actual.PublicIpRefs != nil {
			stepCtx.Require().Equal(expected.PublicIpRefs, *actual.PublicIpRefs, "PublicIpRefs should match expected")
		}

		subnetRef, err := asNetworkReferenceURN(actual.SubnetRef)
		if err != nil {
			ctx.Error(err)
		}
		stepCtx.Require().Equal(expected.SubnetRef, subnetRef, "SubnetRef should match expected")
	})
}

func (suite *NetworkV1TestSuite) verifySecurityGroupSpecStep(ctx provider.StepCtx, expected *secalib.SecurityGroupSpecV1, actual schema.SecurityGroupSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(len(expected.Rules), len(actual.Rules), "Rule list length should match expected")
		for i := 0; i < len(expected.Rules); i++ {
			expectedRule := expected.Rules[i]
			actualRule := actual.Rules[i]
			stepCtx.Require().Equal(expectedRule.Direction, string(actualRule.Direction), fmt.Sprintf("Rule [%d] Direction should match expected", i))
		}
	})
}

func asNetworkReferenceURN(ref schema.Reference) (string, error) {
	urn, err := ref.AsReferenceURN()
	if err != nil {
		return "", fmt.Errorf("error extracting URN from reference: %w", err)
	}
	return string(urn), nil
}
