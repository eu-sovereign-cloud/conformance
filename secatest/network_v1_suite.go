package secatest

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	compute "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.compute.v1"
	network "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

const (
	routeTableDefaultDestination = "0.0.0.0/0"
)

type NetworkV1TestSuite struct {
	regionalTestSuite

	availableZones []string
	storageSkus    []string
	instanceSkus   []string
	networkSkus    []string
}

func (suite *NetworkV1TestSuite) TestNetworkV1(t provider.T) {
	slog.Info("Starting Network Lifecycle Test")

	t.Title("Network Lifecycle Test")
	configureTags(t, secalib.NetworkProviderV1, secalib.NetworkKind, secalib.InternetGatewayKind, secalib.NicKind, secalib.PublicIPKind, secalib.RouteTableKind,
		secalib.SubnetKind, secalib.SecurityGroupKind)

	// TODO Receive via configuration the network cidr ranges and calculate the cidr
	networkCIDR1 := "10.1.0.0/16"
	networkCIDR2 := "10.2.0.0/16"

	// TODO Calculate the subnet cidr from network cidr
	subnetCIDR := "10.1.0.0/24"

	// TODO Calculate the nic cidr from subnet cidr
	nicAddress := "10.1.0.1"

	// TODO Receive via configuration the public ip address range and calculate an ip address
	publicIPAddress := "192.168.0.1"

	// Select zone
	zone := suite.availableZones[rand.Intn(len(suite.availableZones))]

	// Select skus
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	networkSkuName := suite.networkSkus[rand.Intn(len(suite.networkSkus))]

	// TODO Dynamically create before the scenario
	workspaceName := secalib.GenerateWorkspaceName()

	// Generate scenario data
	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)

	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	instanceName := secalib.GenerateInstanceName()

	networkSkuRef := secalib.GenerateSkuRef(networkSkuName)
	networkName := secalib.GenerateNetworkName()
	networkRef := secalib.GenerateNetworkRef(networkName)
	networkResource := secalib.GenerateNetworkResource(suite.tenant, workspaceName, networkName)

	internetGatewayName := secalib.GenerateInternetGatewayName()
	internetGatewayRef := secalib.GenerateInternetGatewayRef(internetGatewayName)
	internetGatewayResource := secalib.GenerateInternetGatewayResource(suite.tenant, workspaceName, internetGatewayName)

	routeTableName := secalib.GenerateRouteTableName()
	routeTableRef := secalib.GenerateRouteTableRef(routeTableName)
	routeTableResource := secalib.GenerateRouteTableResource(suite.tenant, workspaceName, routeTableName)

	subnetName := secalib.GenerateSubnetName()
	subnetRef := secalib.GenerateSubnetRef(subnetName)
	subnetResource := secalib.GenerateSubnetResource(suite.tenant, workspaceName, subnetName)

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
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
					Workspace: workspaceName,
					Region:    suite.region,
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
						Zone:          zone,
						BootDeviceRef: blockStorageRef,
					},
				},
				Network: &mock.ResourceParams[secalib.NetworkSpecV1]{
					Name: networkName,
					InitialSpec: &secalib.NetworkSpecV1{
						Cidr:          &secalib.NetworkSpecCIDRV1{Ipv4: networkCIDR1},
						SkuRef:        networkSkuRef,
						RouteTableRef: routeTableRef,
					},
					UpdatedSpec: &secalib.NetworkSpecV1{
						Cidr:          &secalib.NetworkSpecCIDRV1{Ipv4: networkCIDR2},
						SkuRef:        networkSkuRef,
						RouteTableRef: routeTableRef,
					},
				},
				InternetGateway: &mock.ResourceParams[secalib.InternetGatewaySpecV1]{
					Name:        internetGatewayName,
					InitialSpec: &secalib.InternetGatewaySpecV1{EgressOnly: false},
				},
				RouteTable: &mock.ResourceParams[secalib.RouteTableSpecV1]{
					Name: routeTableName,
					InitialSpec: &secalib.RouteTableSpecV1{
						LocalRef: networkRef,
						Routes: []*secalib.RouteTableRouteV1{
							{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: internetGatewayRef},
						},
					},
				},
				Subnet: &mock.ResourceParams[secalib.SubnetSpecV1]{
					Name: subnetName,
					InitialSpec: &secalib.SubnetSpecV1{
						Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCIDR},
						Zone: zone,
					},
				},
				NIC: &mock.ResourceParams[secalib.NICSpecV1]{
					Name: nicName,
					InitialSpec: &secalib.NICSpecV1{
						Addresses:    []string{nicAddress},
						PublicIpRefs: []string{publicIPRef},
						SubnetRef:    subnetRef,
					},
				},
				PublicIP: &mock.ResourceParams[secalib.PublicIpSpecV1]{
					Name: publicIPName,
					InitialSpec: &secalib.PublicIpSpecV1{
						Version: secalib.IPVersion4,
						Address: publicIPAddress,
					},
				},
				SecurityGroup: &mock.ResourceParams[secalib.SecurityGroupSpecV1]{
					Name: securityGroupName,
					InitialSpec: &secalib.SecurityGroupSpecV1{
						Rules: []*secalib.SecurityGroupRuleV1{{Direction: secalib.SecurityRuleDirectionIngress}},
					},
				},
			})
		if err != nil {
			slog.Error("Failed to create network scenario", "error", err)
			return
		}
		suite.mockClient = wm
	}

	ctx := context.Background()

	var networkResp *network.Network
	var gatewayResp *network.InternetGateway
	var routeResp *network.RouteTable
	var subnetResp *network.Subnet
	var publicIpResp *network.PublicIp
	var nicResp *network.Nic
	var secGroupResp *network.SecurityGroup
	var blockResp *storage.BlockStorage
	var instanceResp *compute.Instance

	var err error

	var expectedNetworkMetadata *secalib.Metadata
	var expectedNetworkSpec *secalib.NetworkSpecV1

	// Step 1
	t.WithNewStep("Create network", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			// TODO Add the provider prefix in each operation
			operationStepParameter, "CreateOrUpdateNetwork",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      networkName,
		}
		net := &network.Network{
			Spec: network.NetworkSpec{
				Cidr:          network.Cidr{Ipv4: &networkCIDR1},
				SkuRef:        networkSkuRef,
				RouteTableRef: routeTableRef,
			},
		}
		networkResp, err = suite.client.NetworkV1.CreateOrUpdateNetwork(ctx, wref, net, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMetadata = &secalib.Metadata{
			Name:       networkName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   networkResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NetworkKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMetadata, networkResp.Metadata)

		expectedNetworkSpec = &secalib.NetworkSpecV1{
			Cidr:          &secalib.NetworkSpecCIDRV1{Ipv4: networkCIDR1},
			SkuRef:        networkSkuRef,
			RouteTableRef: routeTableRef,
		}
		verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	// Step 2
	t.WithNewStep("Get created network", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetNetwork",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      networkName,
		}
		networkResp, err = suite.client.NetworkV1.GetNetwork(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMetadata.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMetadata, networkResp.Metadata)

		verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	// Step 3
	t.WithNewStep("Update network", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateNetwork",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      networkName,
		}
		networkResp, err = suite.client.NetworkV1.CreateOrUpdateNetwork(ctx, wref, networkResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMetadata.Verb = http.MethodPut
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMetadata, networkResp.Metadata)

		expectedNetworkSpec.Cidr.Ipv4 = networkCIDR2
		verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	// Step 4
	t.WithNewStep("Get updated network", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetNetwork",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      networkName,
		}
		networkResp, err = suite.client.NetworkV1.GetNetwork(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMetadata.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMetadata, networkResp.Metadata)

		verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	var expectedGatewayMetadata *secalib.Metadata
	var expectedGatewaySpec *secalib.InternetGatewaySpecV1

	// Step 5
	t.WithNewStep("Create internet gateway", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateInternetGateway",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      internetGatewayName,
		}
		gtw := &network.InternetGateway{}
		gatewayResp, err = suite.client.NetworkV1.CreateOrUpdateInternetGateway(ctx, wref, gtw, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMetadata = &secalib.Metadata{
			Name:       internetGatewayName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   internetGatewayResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InternetGatewayKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMetadata, gatewayResp.Metadata)

		expectedGatewaySpec = &secalib.InternetGatewaySpecV1{
			EgressOnly: false,
		}
		verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	// Step 6
	t.WithNewStep("Get created internet gateway", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInternetGateway",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      internetGatewayName,
		}
		gatewayResp, err = suite.client.NetworkV1.GetInternetGateway(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMetadata.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMetadata, gatewayResp.Metadata)

		verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	// Step 7
	t.WithNewStep("Update internet gateway", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateInternetGateway",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      internetGatewayName,
		}
		gatewayResp, err = suite.client.NetworkV1.CreateOrUpdateInternetGateway(ctx, wref, gatewayResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMetadata.Verb = http.MethodPut
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMetadata, gatewayResp.Metadata)

		//expectedNetworkSpec.Cidr.Ipv4 = networkCIDR2
		verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	// Step 8
	t.WithNewStep("Get updated internet gateway", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInternetGateway",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      internetGatewayName,
		}
		gatewayResp, err = suite.client.NetworkV1.GetInternetGateway(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMetadata.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMetadata, gatewayResp.Metadata)

		verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	var expectedRouteMetadata *secalib.Metadata
	var expectedRouteSpec *secalib.RouteTableSpecV1

	// Step 9
	t.WithNewStep("Create route table", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateRouteTable",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      routeTableName,
		}
		route := &network.RouteTable{}
		routeResp, err = suite.client.NetworkV1.CreateOrUpdateRouteTable(ctx, wref, route, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMetadata = &secalib.Metadata{
			Name:       routeTableName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   routeTableResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RouteTableKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMetadata, routeResp.Metadata)

		expectedRouteSpec = &secalib.RouteTableSpecV1{
			LocalRef: networkRef,
			Routes: []*secalib.RouteTableRouteV1{
				{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: internetGatewayRef},
			},
		}
		verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	// Step 10
	t.WithNewStep("Get created route table", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetRouteTable",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      routeTableName,
		}
		routeResp, err = suite.client.NetworkV1.GetRouteTable(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMetadata.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMetadata, routeResp.Metadata)

		verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	// Step 11
	t.WithNewStep("Update route table", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateRouteTable",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      routeTableName,
		}
		routeResp, err = suite.client.NetworkV1.CreateOrUpdateRouteTable(ctx, wref, routeResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMetadata.Verb = http.MethodPut
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMetadata, routeResp.Metadata)

		//expectedRouteSpec.Cidr.Ipv4 = networkCIDR2
		verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	// Step 12
	t.WithNewStep("Get updated route table", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetRouteTable",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      routeTableName,
		}
		routeResp, err = suite.client.NetworkV1.GetRouteTable(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMetadata.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMetadata, routeResp.Metadata)

		verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	var expectedSubnetMetadata *secalib.Metadata
	var expectedSubnetSpec *secalib.SubnetSpecV1

	// Step 13
	t.WithNewStep("Create subnet", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateSubnet",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      subnetName,
		}
		sub := &network.Subnet{
			Spec: network.SubnetSpec{
				Cidr: network.Cidr{Ipv4: &subnetCIDR},
				Zone: zone,
			},
		}
		subnetResp, err = suite.client.NetworkV1.CreateOrUpdateSubnet(ctx, wref, sub, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMetadata = &secalib.Metadata{
			Name:       subnetName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   subnetResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SubnetKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMetadata, subnetResp.Metadata)

		expectedSubnetSpec = &secalib.SubnetSpecV1{
			Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCIDR},
			Zone: zone,
		}
		verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	// Step 14
	t.WithNewStep("Get created subnet", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetNetwork",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      subnetName,
		}
		subnetResp, err = suite.client.NetworkV1.GetSubnet(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMetadata.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMetadata, subnetResp.Metadata)

		verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	// Step 15
	t.WithNewStep("Update subnet", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateSubnet",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      subnetName,
		}
		subnetResp, err = suite.client.NetworkV1.CreateOrUpdateSubnet(ctx, wref, subnetResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMetadata.Verb = http.MethodPut
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMetadata, subnetResp.Metadata)

		//subnetResp.Cidr.Ipv4 = networkCIDR2
		verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	// Step 16
	t.WithNewStep("Get updated subnet", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetSubnet",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      subnetName,
		}
		subnetResp, err = suite.client.NetworkV1.GetSubnet(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMetadata.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMetadata, subnetResp.Metadata)

		verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	var expectedIpMetadata *secalib.Metadata
	var expectedIpSpec *secalib.PublicIpSpecV1

	// Step 17
	t.WithNewStep("Create public ip", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdatePublicIp",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIPName,
		}
		ip := &network.PublicIp{
			Spec: network.PublicIpSpec{
				Address: &publicIPAddress,
				Version: secalib.IPVersion4,
			},
		}
		publicIpResp, err = suite.client.NetworkV1.CreateOrUpdatePublicIp(ctx, wref, ip, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMetadata = &secalib.Metadata{
			Name:       publicIPName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   publicIPResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.PublicIPKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMetadata, publicIpResp.Metadata)

		expectedIpSpec = &secalib.PublicIpSpecV1{
			Version: secalib.IPVersion4,
			Address: publicIPAddress,
		}
		verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	// Step 18
	t.WithNewStep("Get created public ip", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetPublicIP",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIPName,
		}
		publicIpResp, err = suite.client.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMetadata.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMetadata, publicIpResp.Metadata)

		verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	// Step 19
	t.WithNewStep("Update public ip", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdatePublicIp",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIPName,
		}
		publicIpResp, err = suite.client.NetworkV1.CreateOrUpdatePublicIp(ctx, wref, publicIpResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMetadata.Verb = http.MethodPut
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMetadata, publicIpResp.Metadata)

		//expectedIpSpec.Cidr.Ipv4 = networkCIDR2
		verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	// Step 20
	t.WithNewStep("Get updated public ip", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetPublicIp",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIPName,
		}
		publicIpResp, err = suite.client.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMetadata.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMetadata, publicIpResp.Metadata)

		verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	var expectedNicMetadata *secalib.Metadata
	var expectedNicSpec *secalib.NICSpecV1

	// Step 21
	t.WithNewStep("Create nic", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateNic",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nicName,
		}
		nic := &network.Nic{
			Spec: network.NicSpec{
				Addresses:    []string{nicAddress},
				PublicIpRefs: &[]interface{}{publicIPRef},
				SubnetRef:    subnetRef,
			},
		}
		nicResp, err = suite.client.NetworkV1.CreateOrUpdateNic(ctx, wref, nic, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMetadata = &secalib.Metadata{
			Name:       nicName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   nicResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NicKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMetadata, nicResp.Metadata)

		expectedNicSpec = &secalib.NICSpecV1{
			Addresses:    []string{nicAddress},
			PublicIpRefs: []string{publicIPRef},
			SubnetRef:    subnetRef,
		}
		verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	// Step 22
	t.WithNewStep("Get created nic", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetNic",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nicName,
		}
		nicResp, err = suite.client.NetworkV1.GetNic(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMetadata.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMetadata, nicResp.Metadata)

		verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	// Step 23
	t.WithNewStep("Update nic", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateNic",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nicName,
		}
		nicResp, err = suite.client.NetworkV1.CreateOrUpdateNic(ctx, wref, nicResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMetadata.Verb = http.MethodPut
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMetadata, nicResp.Metadata)

		//expectedNicSpec.Cidr.Ipv4 = networkCIDR2
		verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	// Step 24
	t.WithNewStep("Get updated nic", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetNic",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nicName,
		}
		nicResp, err = suite.client.NetworkV1.GetNic(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMetadata.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMetadata, nicResp.Metadata)

		verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	var expectedSecGroupMetadata *secalib.Metadata
	var expectedSecGroupSpec *secalib.SecurityGroupSpecV1

	// Step 25
	t.WithNewStep("Create security group", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateSecurityGroup",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      securityGroupName,
		}
		secGroup := &network.SecurityGroup{
			Spec: network.SecurityGroupSpec{
				Rules: []network.SecurityGroupRuleSpec{
					{
						Direction: secalib.SecurityRuleDirectionIngress,
					},
				},
			},
		}
		secGroupResp, err = suite.client.NetworkV1.CreateOrUpdateSecurityGroup(ctx, wref, secGroup, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, secGroupResp)

		expectedSecGroupMetadata = &secalib.Metadata{
			Name:       securityGroupName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   securityGroupResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SecurityGroupKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedSecGroupMetadata, secGroupResp.Metadata)

		expectedSecGroupSpec = &secalib.SecurityGroupSpecV1{
			Rules: []*secalib.SecurityGroupRuleV1{{Direction: secalib.SecurityRuleDirectionIngress}},
		}
		verifySecurityGroupSpecStep(sCtx, expectedSecGroupSpec, secGroupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*secGroupResp.Status.State)},
		)
	})

	// Step 26
	t.WithNewStep("Get created security group", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetSecurityGroup",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      securityGroupName,
		}
		secGroupResp, err = suite.client.NetworkV1.GetSecurityGroup(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, secGroupResp)

		expectedSecGroupMetadata.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedSecGroupMetadata, secGroupResp.Metadata)

		verifySecurityGroupSpecStep(sCtx, expectedSecGroupSpec, secGroupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*secGroupResp.Status.State)},
		)
	})

	// Step 27
	t.WithNewStep("Update security group", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateSecurityGroup",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      securityGroupName,
		}
		secGroupResp, err = suite.client.NetworkV1.CreateOrUpdateSecurityGroup(ctx, wref, secGroupResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, secGroupResp)

		expectedSecGroupMetadata.Verb = http.MethodPut
		verifyNetworkRegionalMetadataStep(sCtx, expectedSecGroupMetadata, secGroupResp.Metadata)

		//expectedSecGroupSpec.Cidr.Ipv4 = networkCIDR2
		verifySecurityGroupSpecStep(sCtx, expectedSecGroupSpec, secGroupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*secGroupResp.Status.State)},
		)
	})

	// Step 28
	t.WithNewStep("Get updated security group", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetSecurityGroup",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      securityGroupName,
		}
		secGroupResp, err = suite.client.NetworkV1.GetSecurityGroup(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, secGroupResp)

		expectedSecGroupMetadata.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedSecGroupMetadata, secGroupResp.Metadata)

		verifySecurityGroupSpecStep(sCtx, expectedSecGroupSpec, secGroupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*secGroupResp.Status.State)},
		)
	})

	// Step 29
	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		block := &storage.BlockStorage{
			Spec: storage.BlockStorageSpec{
				SizeGB: blockStorageSize,
				SkuRef: storageSkuRef,
			},
		}
		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, wref, block, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	// Step 30
	t.WithNewStep("Get created block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

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

	// Step 31
	t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		inst := &compute.Instance{
			Spec: compute.InstanceSpec{
				SkuRef: instanceSkuRef,
				Zone:   zone,
			},
		}
		inst.Spec.BootVolume.DeviceRef = blockStorageRef

		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, wref, inst, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	// Step 32
	t.WithNewStep("Get created instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

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

	// Step 33
	t.WithNewStep("Delete instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteInstance",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.ComputeV1.DeleteInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 34
	t.WithNewStep("Get deleted instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		_, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	// Step 35
	t.WithNewStep("Delete block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteBlockStorage",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.StorageV1.DeleteBlockStorage(ctx, blockResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 36
	t.WithNewStep("Get deleted block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		_, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	// Step 37
	t.WithNewStep("Delete security group", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteSecurityGroup",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.NetworkV1.DeleteSecurityGroup(ctx, secGroupResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 38
	t.WithNewStep("Get deleted security group", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetSecurityGroup",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      securityGroupName,
		}
		_, err = suite.client.NetworkV1.GetSecurityGroup(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	// Step 39
	t.WithNewStep("Delete nic", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteNic",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.NetworkV1.DeleteNic(ctx, nicResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 40
	t.WithNewStep("Get deleted nic", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetNic",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nicName,
		}
		_, err = suite.client.NetworkV1.GetNic(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	// Step 41
	t.WithNewStep("Delete public ip", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeletePublicIP",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.NetworkV1.DeletePublicIp(ctx, publicIpResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 42
	t.WithNewStep("Get deleted public ip", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetPublicIP",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIPName,
		}
		_, err = suite.client.NetworkV1.GetPublicIp(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	// Step 43
	t.WithNewStep("Delete subnet", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteSubnet",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.NetworkV1.DeleteSubnet(ctx, subnetResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 44
	t.WithNewStep("Get deleted subnet", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetSubnet",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      subnetName,
		}
		_, err = suite.client.NetworkV1.GetSubnet(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	// Step 45
	t.WithNewStep("Delete route table", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteRouteTable",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.NetworkV1.DeleteRouteTable(ctx, routeResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 46
	t.WithNewStep("Get deleted route table", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetRouteTable",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      routeTableName,
		}
		_, err = suite.client.NetworkV1.GetRouteTable(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	// Step 47
	t.WithNewStep("Delete internet gateway", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteInternetGateway",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.NetworkV1.DeleteInternetGateway(ctx, gatewayResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 48
	t.WithNewStep("Get deleted internet gateway", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInternetGateway",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      internetGatewayName,
		}
		_, err = suite.client.NetworkV1.GetInternetGateway(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	// Step 49
	t.WithNewStep("Delete network", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteNetwork",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.NetworkV1.DeleteNetwork(ctx, networkResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 50
	t.WithNewStep("Get deleted network", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetNetwork",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

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

func verifyNetworkZonalMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *network.ZonalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Workspace:  *metadata.Workspace,
		Region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyNetworkRegionalMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *network.RegionalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Workspace:  *metadata.Workspace,
		Region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyNetworkSpecStep(ctx provider.StepCtx, expected *secalib.NetworkSpecV1, actual network.NetworkSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_cidr_ipv4", expected.Cidr.Ipv4,
			"actual_cidr_ipv4", actual.Cidr.Ipv4,

			"expected_cidr_ipv6", expected.Cidr.Ipv6,
			"actual_cidr_ipv6", actual.Cidr.Ipv6,

			"expected_sku_ref", expected.SkuRef,
			"actual_sku_ref", actual.SkuRef,

			"expected_route_table_ref", expected.RouteTableRef,
			"actual_route_table_ref", actual.RouteTableRef,
		)
		if actual.Cidr.Ipv4 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv4, *actual.Cidr.Ipv4, "Cidr.Ipv4 should match expected")
		}
		if actual.Cidr.Ipv6 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv6, *actual.Cidr.Ipv6, "Cidr.Ipv6 should match expected")
		}
		stepCtx.Require().Equal(expected.SkuRef, actual.SkuRef, "SkuRef should match expected")
		stepCtx.Require().Equal(expected.RouteTableRef, actual.RouteTableRef, "RouteTableRef should match expected")
	})
}

func verifyInternetGatewaySpecStep(ctx provider.StepCtx, expected *secalib.InternetGatewaySpecV1, actual network.InternetGatewaySpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_egress_only", expected.EgressOnly,
			"actual_egress_only", *actual.EgressOnly,
		)
		if actual.EgressOnly != nil {
			stepCtx.Require().Equal(expected.EgressOnly, *actual.EgressOnly, "EgressOnly should match expected")
		}
	})
}

func verifyRouteTableSpecStep(ctx provider.StepCtx, expected *secalib.RouteTableSpecV1, actual network.RouteTableSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_local_ref", expected.LocalRef,
			"actual_local_ref", actual.LocalRef,
		)
		stepCtx.Require().Equal(expected.LocalRef, actual.LocalRef, "LocalRef should match expected")

		stepCtx.Require().Equal(len(expected.Routes), len(actual.Routes), "Route list length should match expected")
		for i := 0; i < len(expected.Routes); i++ {
			expectedRoute := expected.Routes[i]
			actualRoute := actual.Routes[i]
			stepCtx.WithNewParameters(
				fmt.Sprintf("expected_route[%d]_destination_cidr_block", i), expectedRoute.DestinationCidrBlock,
				fmt.Sprintf("actual_route[%d]_destination_cidr_block", i), actualRoute.DestinationCidrBlock,

				fmt.Sprintf("expected_route[%d]_target_ref", i), expectedRoute.TargetRef,
				fmt.Sprintf("actual_route[%d]_target_ref", i), actualRoute.TargetRef,
			)
			stepCtx.Require().Equal(expectedRoute.DestinationCidrBlock, actualRoute.DestinationCidrBlock, fmt.Sprintf("Route [%d] DestinationCidrBlock should match expected", i))
			stepCtx.Require().Equal(expectedRoute.TargetRef, actualRoute.TargetRef, fmt.Sprintf("Route [%d] TargetRef should match expected", i))
		}
	})
}

func verifySubNetSpecStep(ctx provider.StepCtx, expected *secalib.SubnetSpecV1, actual network.SubnetSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_cidr_ipv4", expected.Cidr.Ipv4,
			"actual_cidr_ipv4", actual.Cidr.Ipv4,

			"expected_cidr_ipv6", expected.Cidr.Ipv6,
			"actual_cidr_ipv6", actual.Cidr.Ipv6,

			"expected_zone", expected.Zone,
			"actual_zone", actual.Zone,
		)
		if actual.Cidr.Ipv4 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv4, *actual.Cidr.Ipv4, "Cidr.Ipv4 should match expected")
		}
		if actual.Cidr.Ipv6 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv6, *actual.Cidr.Ipv6, "Cidr.Ipv6 should match expected")
		}
		stepCtx.Require().Equal(expected.Zone, actual.Zone, "Zone should match expected")
	})
}

func verifyPublicIpSpecStep(ctx provider.StepCtx, expected *secalib.PublicIpSpecV1, actual network.PublicIpSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_version", expected.Version,
			"actual_version", actual.Version,

			"expected_address", expected.Address,
			"actual_address", actual.Address,
		)
		stepCtx.Require().Equal(expected.Version, string(actual.Version), "Version should match expected")
		if actual.Address != nil {
			stepCtx.Require().Equal(expected.Address, *actual.Address, "Address should match expected")
		}
	})
}

func verifyNicSpecStep(ctx provider.StepCtx, expected *secalib.NICSpecV1, actual network.NicSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_addresses", expected.Addresses,
			"actual_addresses", actual.Addresses,

			"expected_public_ip_refs", expected.PublicIpRefs,
			"actual_public_ip_refs", actual.PublicIpRefs,

			"expected_subnet_ref", expected.SubnetRef,
			"actual_subnet_ref", actual.SubnetRef,
		)
		stepCtx.Require().Equal(expected.Addresses, actual.Addresses, "Addresses should match expected")
		if actual.PublicIpRefs != nil {
			stepCtx.Require().Equal(expected.PublicIpRefs, *actual.PublicIpRefs, "PublicIpRefs should match expected")
		}
		stepCtx.Require().Equal(expected.SubnetRef, actual.SubnetRef, "SubnetRef should match expected")
	})
}

func verifySecurityGroupSpecStep(ctx provider.StepCtx, expected *secalib.SecurityGroupSpecV1, actual network.SecurityGroupSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(len(expected.Rules), len(actual.Rules), "Rule list length should match expected")
		for i := 0; i < len(expected.Rules); i++ {
			expectedRule := expected.Rules[i]
			actualRule := actual.Rules[i]
			stepCtx.WithNewParameters(
				fmt.Sprintf("expected_rule[%d]_direction", i), expectedRule.Direction,
				fmt.Sprintf("actual_rule[%d]_direction", i), actualRule.Direction,
			)
			stepCtx.Require().Equal(expectedRule.Direction, string(actualRule.Direction), fmt.Sprintf("Rule [%d] Direction should match expected", i))
		}

	})
}
