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
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/google/uuid"

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

func (suite *NetworkV1TestSuite) generateLifecycleParams() (*secalib.NetworkLifeCycleParamsV1, error) {
	workspaceName := secalib.GenerateWorkspaceName()
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]
	blockStorageName := secalib.GenerateBlockStorageName()
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	instanceName := secalib.GenerateInstanceName()
	networkSkuName1 := suite.networkSkus[rand.Intn(len(suite.networkSkus))]
	networkSkuName2 := suite.networkSkus[rand.Intn(len(suite.networkSkus))]
	networkName := secalib.GenerateNetworkName()
	internetGatewayName := secalib.GenerateInternetGatewayName()
	routeTableName := secalib.GenerateRouteTableName()
	subnetName := secalib.GenerateSubnetName()
	nicName := secalib.GenerateNicName()
	publicIPName := secalib.GeneratePublicIPName()
	securityGroupName := secalib.GenerateSecurityGroupName()

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)
	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	instanceRef := secalib.GenerateInstanceRef(instanceName)
	networkSkuRef1 := secalib.GenerateSkuRef(networkSkuName1)
	networkSkuRef2 := secalib.GenerateSkuRef(networkSkuName2)
	networkRef := secalib.GenerateNetworkRef(networkName)
	internetGatewayRef := secalib.GenerateInternetGatewayRef(internetGatewayName)
	routeTableRef := secalib.GenerateRouteTableRef(routeTableName)
	subnetRef := secalib.GenerateSubnetRef(subnetName)
	publicIPRef := secalib.GeneratePublicIPRef(publicIPName)

	// Random data
	blockStorageSize := secalib.GenerateBlockStorageSize()

	// Select the zones
	zone1 := suite.regionZones[rand.Intn(len(suite.regionZones))]
	zone2 := suite.regionZones[rand.Intn(len(suite.regionZones))]

	subnetCidr, err := secalib.GenerateSubnetCidr(suite.networkCidr, 8, 1)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate subnet cidr: %v", err)
	}

	nicAddress1, err := secalib.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate nic address: %v", err)
	}
	nicAddress2, err := secalib.GenerateNicAddress(subnetCidr, 2)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate nic address: %v", err)
	}

	publicIpAddress1, err := secalib.GeneratePublicIp(suite.publicIpsRange, 1)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate public ip: %v", err)
	}
	publicIpAddress2, err := secalib.GeneratePublicIp(suite.publicIpsRange, 2)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate public ip: %v", err)
	}

	return &secalib.NetworkLifeCycleParamsV1{
		Workspace: &secalib.ResourceParams[secalib.WorkspaceSpecV1]{
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
		BlockStorage: &secalib.ResourceParams[secalib.BlockStorageSpecV1]{
			Name: blockStorageName,
			InitialSpec: &secalib.BlockStorageSpecV1{
				SkuRef: storageSkuRef,
				SizeGB: blockStorageSize,
			},
		},
		Instance: &secalib.ResourceParams[secalib.InstanceSpecV1]{
			Name: instanceName,
			InitialSpec: &secalib.InstanceSpecV1{
				SkuRef:        instanceSkuRef,
				Zone:          zone1,
				BootDeviceRef: blockStorageRef,
			},
		},
		Network: &secalib.ResourceParams[secalib.NetworkSpecV1]{
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
		InternetGateway: &secalib.ResourceParams[secalib.InternetGatewaySpecV1]{
			Name:        internetGatewayName,
			InitialSpec: &secalib.InternetGatewaySpecV1{EgressOnly: false},
			UpdatedSpec: &secalib.InternetGatewaySpecV1{EgressOnly: true},
		},
		RouteTable: &secalib.ResourceParams[secalib.RouteTableSpecV1]{
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
		Subnet: &secalib.ResourceParams[secalib.SubnetSpecV1]{
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
		Nic: &secalib.ResourceParams[secalib.NICSpecV1]{
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
		PublicIp: &secalib.ResourceParams[secalib.PublicIpSpecV1]{
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
		SecurityGroup: &secalib.ResourceParams[secalib.SecurityGroupSpecV1]{
			Name: securityGroupName,
			InitialSpec: &secalib.SecurityGroupSpecV1{
				Rules: []*secalib.SecurityGroupRuleV1{{Direction: secalib.SecurityRuleDirectionIngress}},
			},
			UpdatedSpec: &secalib.SecurityGroupSpecV1{
				Rules: []*secalib.SecurityGroupRuleV1{{Direction: secalib.SecurityRuleDirectionEgress}},
			},
		},
	}, nil
}

func (suite *NetworkV1TestSuite) TestNetworkLifeCycleV1(t provider.T) {
	slog.Info("Starting Network Lifecycle Test")

	t.Title("Network Lifecycle Test")
	configureTags(t, secalib.NetworkProviderV1, secalib.NetworkKind, secalib.InternetGatewayKind, secalib.NicKind, secalib.PublicIPKind, secalib.RouteTableKind,
		secalib.SubnetKind, secalib.SecurityGroupKind)

	params, err := suite.generateLifecycleParams()
	if err != nil {
		t.Fatalf("Failed to generate lifecycle params: %v", err)
	}

	// Resource URIs
	networkResource := secalib.GenerateNetworkResource(suite.tenant, params.Workspace.Name, params.Network.Name)
	internetGatewayResource := secalib.GenerateInternetGatewayResource(suite.tenant, params.Workspace.Name, params.InternetGateway.Name)
	routeTableResource := secalib.GenerateRouteTableResource(suite.tenant, params.Workspace.Name, params.RouteTable.Name)
	subnetResource := secalib.GenerateSubnetResource(suite.tenant, params.Workspace.Name, params.Subnet.Name)
	nicResource := secalib.GenerateNicResource(suite.tenant, params.Workspace.Name, params.Nic.Name)
	publicIPResource := secalib.GeneratePublicIPResource(suite.tenant, params.Workspace.Name, params.PublicIp.Name)
	securityGroupResource := secalib.GenerateSecurityGroupResource(suite.tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		scenarios := mock.NewNetworkScenariosV1(suite.authToken, suite.tenant, suite.region, suite.mockServerURL)

		id := uuid.New().String()
		wm, err := scenarios.ConfigureLifecycleScenario(id, params)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()

	// Workspace
	var workResp *workspace.Workspace

	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Workspace.Name,
		}
		ws := &workspace.Workspace{}
		workResp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, tref, ws, nil)
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
			Name:   params.Workspace.Name,
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
	var networkResp *network.Network
	var expectedNetworkMeta *secalib.Metadata

	t.WithNewStep("Create network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNetwork", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Network.Name,
		}
		net := &network.Network{
			Spec: network.NetworkSpec{
				Cidr:          network.Cidr{Ipv4: &suite.networkCidr},
				SkuRef:        params.Network.InitialSpec.SkuRef,
				RouteTableRef: params.Network.InitialSpec.RouteTableRef,
			},
		}
		networkResp, err = suite.client.NetworkV1.CreateOrUpdateNetwork(ctx, wref, net, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta = &secalib.Metadata{
			Name:       params.Network.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   networkResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NetworkKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		verifyNetworkSpecStep(sCtx, params.Network.InitialSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	t.WithNewStep("Get created network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Network.Name,
		}
		networkResp, err = suite.client.NetworkV1.GetNetwork(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		verifyNetworkSpecStep(sCtx, params.Network.InitialSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	t.WithNewStep("Update network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNetwork", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Network.Name,
		}
		networkResp, err = suite.client.NetworkV1.CreateOrUpdateNetwork(ctx, wref, networkResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta.Verb = http.MethodPut
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		verifyNetworkSpecStep(sCtx, params.Network.UpdatedSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Network.Name,
		}
		networkResp, err = suite.client.NetworkV1.GetNetwork(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		verifyNetworkSpecStep(sCtx, params.Network.UpdatedSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	// Internet gateway
	var gatewayResp *network.InternetGateway
	var expectedGatewayMeta *secalib.Metadata

	t.WithNewStep("Create internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateInternetGateway", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.InternetGateway.Name,
		}
		gtw := &network.InternetGateway{}
		gatewayResp, err = suite.client.NetworkV1.CreateOrUpdateInternetGateway(ctx, wref, gtw, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta = &secalib.Metadata{
			Name:       params.InternetGateway.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   internetGatewayResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InternetGatewayKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		verifyInternetGatewaySpecStep(sCtx, params.InternetGateway.InitialSpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	t.WithNewStep("Get created internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.InternetGateway.Name,
		}
		gatewayResp, err = suite.client.NetworkV1.GetInternetGateway(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		verifyInternetGatewaySpecStep(sCtx, params.InternetGateway.InitialSpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	t.WithNewStep("Update internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateInternetGateway", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.InternetGateway.Name,
		}
		gatewayResp, err = suite.client.NetworkV1.CreateOrUpdateInternetGateway(ctx, wref, gatewayResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta.Verb = http.MethodPut
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		verifyInternetGatewaySpecStep(sCtx, params.InternetGateway.UpdatedSpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.InternetGateway.Name,
		}
		gatewayResp, err = suite.client.NetworkV1.GetInternetGateway(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		verifyInternetGatewaySpecStep(sCtx, params.InternetGateway.UpdatedSpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	// Route table
	var routeResp *network.RouteTable
	var expectedRouteMeta *secalib.Metadata

	t.WithNewStep("Create route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateRouteTable", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.RouteTable.Name,
		}
		route := &network.RouteTable{}
		routeResp, err = suite.client.NetworkV1.CreateOrUpdateRouteTable(ctx, wref, route, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta = &secalib.Metadata{
			Name:       params.RouteTable.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   routeTableResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RouteTableKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		verifyRouteTableSpecStep(sCtx, params.RouteTable.InitialSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	t.WithNewStep("Get created route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.RouteTable.Name,
		}
		routeResp, err = suite.client.NetworkV1.GetRouteTable(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		verifyRouteTableSpecStep(sCtx, params.RouteTable.InitialSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	t.WithNewStep("Update route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateRouteTable", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.RouteTable.Name,
		}
		routeResp, err = suite.client.NetworkV1.CreateOrUpdateRouteTable(ctx, wref, routeResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta.Verb = http.MethodPut
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		verifyRouteTableSpecStep(sCtx, params.RouteTable.UpdatedSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.RouteTable.Name,
		}
		routeResp, err = suite.client.NetworkV1.GetRouteTable(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		verifyRouteTableSpecStep(sCtx, params.RouteTable.UpdatedSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	// Subnet
	var subnetResp *network.Subnet
	var expectedSubnetMeta *secalib.Metadata

	t.WithNewStep("Create subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSubnet", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Subnet.Name,
		}
		sub := &network.Subnet{
			Spec: network.SubnetSpec{
				Cidr: network.Cidr{Ipv4: &params.Subnet.InitialSpec.Cidr.Ipv4},
				Zone: params.Subnet.InitialSpec.Zone,
			},
		}
		subnetResp, err = suite.client.NetworkV1.CreateOrUpdateSubnet(ctx, wref, sub, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta = &secalib.Metadata{
			Name:       params.Subnet.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   subnetResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SubnetKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		verifySubNetSpecStep(sCtx, params.Subnet.InitialSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	t.WithNewStep("Get created subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Subnet.Name,
		}
		subnetResp, err = suite.client.NetworkV1.GetSubnet(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		verifySubNetSpecStep(sCtx, params.Subnet.InitialSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	t.WithNewStep("Update subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSubnet", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Subnet.Name,
		}
		subnetResp, err = suite.client.NetworkV1.CreateOrUpdateSubnet(ctx, wref, subnetResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta.Verb = http.MethodPut
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		verifySubNetSpecStep(sCtx, params.Subnet.UpdatedSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSubnet", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Subnet.Name,
		}
		subnetResp, err = suite.client.NetworkV1.GetSubnet(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		verifySubNetSpecStep(sCtx, params.Subnet.UpdatedSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	// Public ip
	var publicIpResp *network.PublicIp
	var expectedIpMeta *secalib.Metadata

	t.WithNewStep("Create public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdatePublicIp", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.PublicIp.Name,
		}
		ip := &network.PublicIp{
			Spec: network.PublicIpSpec{
				Address: &params.PublicIp.InitialSpec.Address,
				Version: network.IPVersion(params.PublicIp.InitialSpec.Version),
			},
		}
		publicIpResp, err = suite.client.NetworkV1.CreateOrUpdatePublicIp(ctx, wref, ip, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta = &secalib.Metadata{
			Name:       params.PublicIp.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   publicIPResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.PublicIPKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		verifyPublicIpSpecStep(sCtx, params.PublicIp.InitialSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	t.WithNewStep("Get created public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIP", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.PublicIp.Name,
		}
		publicIpResp, err = suite.client.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		verifyPublicIpSpecStep(sCtx, params.PublicIp.InitialSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	t.WithNewStep("Update public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdatePublicIp", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.PublicIp.Name,
		}
		publicIpResp, err = suite.client.NetworkV1.CreateOrUpdatePublicIp(ctx, wref, publicIpResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodPut
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		verifyPublicIpSpecStep(sCtx, params.PublicIp.UpdatedSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.PublicIp.Name,
		}
		publicIpResp, err = suite.client.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		verifyPublicIpSpecStep(sCtx, params.PublicIp.UpdatedSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	// Nic
	var nicResp *network.Nic
	var expectedNicMeta *secalib.Metadata

	t.WithNewStep("Create nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNic", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Nic.Name,
		}
		nic := &network.Nic{
			Spec: network.NicSpec{
				Addresses:    params.Nic.InitialSpec.Addresses,
				PublicIpRefs: &[]interface{}{params.Nic.InitialSpec.PublicIpRefs},
				SubnetRef:    params.Nic.InitialSpec.SubnetRef,
			},
		}
		nicResp, err = suite.client.NetworkV1.CreateOrUpdateNic(ctx, wref, nic, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta = &secalib.Metadata{
			Name:       params.Nic.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   nicResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NicKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		verifyNicSpecStep(sCtx, params.Nic.InitialSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	t.WithNewStep("Get created nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNic", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Nic.Name,
		}
		nicResp, err = suite.client.NetworkV1.GetNic(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		verifyNicSpecStep(sCtx, params.Nic.InitialSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	t.WithNewStep("Update nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNic", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Nic.Name,
		}
		nicResp, err = suite.client.NetworkV1.CreateOrUpdateNic(ctx, wref, nicResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta.Verb = http.MethodPut
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		verifyNicSpecStep(sCtx, params.Nic.UpdatedSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNic", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Nic.Name,
		}
		nicResp, err = suite.client.NetworkV1.GetNic(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		verifyNicSpecStep(sCtx, params.Nic.UpdatedSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	// Security Group
	var groupResp *network.SecurityGroup
	var expectedGroupMeta *secalib.Metadata

	t.WithNewStep("Create security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSecurityGroup", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.SecurityGroup.Name,
		}
		group := &network.SecurityGroup{
			Spec: network.SecurityGroupSpec{
				Rules: []network.SecurityGroupRuleSpec{
					{
						Direction: secalib.SecurityRuleDirectionIngress,
					},
				},
			},
		}
		groupResp, err = suite.client.NetworkV1.CreateOrUpdateSecurityGroup(ctx, wref, group, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta = &secalib.Metadata{
			Name:       params.SecurityGroup.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   securityGroupResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SecurityGroupKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		verifySecurityGroupSpecStep(sCtx, params.SecurityGroup.InitialSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
	})

	t.WithNewStep("Get created security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.SecurityGroup.Name,
		}
		groupResp, err = suite.client.NetworkV1.GetSecurityGroup(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		verifySecurityGroupSpecStep(sCtx, params.SecurityGroup.InitialSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
	})

	t.WithNewStep("Update security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSecurityGroup", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.SecurityGroup.Name,
		}
		groupResp, err = suite.client.NetworkV1.CreateOrUpdateSecurityGroup(ctx, wref, groupResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta.Verb = http.MethodPut
		verifyNetworkRegionalMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		verifySecurityGroupSpecStep(sCtx, params.SecurityGroup.UpdatedSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.SecurityGroup.Name,
		}
		groupResp, err = suite.client.NetworkV1.GetSecurityGroup(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		verifySecurityGroupSpecStep(sCtx, params.SecurityGroup.UpdatedSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
	})

	// Block storage
	var blockResp *storage.BlockStorage

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
		}
		block := &storage.BlockStorage{
			Spec: storage.BlockStorageSpec{
				SizeGB: params.BlockStorage.InitialSpec.SizeGB,
				SkuRef: params.BlockStorage.InitialSpec.SkuRef,
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

	t.WithNewStep("Get created block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
		}
		blockResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	// Instance
	var instanceResp *compute.Instance

	t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "CreateOrUpdateInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		inst := &compute.Instance{
			Spec: compute.InstanceSpec{
				SkuRef: params.Instance.InitialSpec.SkuRef,
				Zone:   params.Instance.InitialSpec.Zone,
			},
		}
		inst.Spec.BootVolume.DeviceRef = params.Instance.InitialSpec.BootDeviceRef

		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, wref, inst, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	t.WithNewStep("Get created instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
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
		suite.setComputeV1StepParams(sCtx, "DeleteInstance", params.Workspace.Name)

		err = suite.client.ComputeV1.DeleteInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		_, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteBlockStorage", params.Workspace.Name)

		err = suite.client.StorageV1.DeleteBlockStorage(ctx, blockResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
		}
		_, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSecurityGroup", params.Workspace.Name)

		err = suite.client.NetworkV1.DeleteSecurityGroup(ctx, groupResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.SecurityGroup.Name,
		}
		_, err = suite.client.NetworkV1.GetSecurityGroup(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNic", params.Workspace.Name)

		err = suite.client.NetworkV1.DeleteNic(ctx, nicResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNic", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Nic.Name,
		}
		_, err = suite.client.NetworkV1.GetNic(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeletePublicIp", params.Workspace.Name)

		err = suite.client.NetworkV1.DeletePublicIp(ctx, publicIpResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.PublicIp.Name,
		}
		_, err = suite.client.NetworkV1.GetPublicIp(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSubnet", params.Workspace.Name)

		err = suite.client.NetworkV1.DeleteSubnet(ctx, subnetResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSubnet", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Subnet.Name,
		}
		_, err = suite.client.NetworkV1.GetSubnet(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteRouteTable", params.Workspace.Name)

		err = suite.client.NetworkV1.DeleteRouteTable(ctx, routeResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.RouteTable.Name,
		}
		_, err = suite.client.NetworkV1.GetRouteTable(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteInternetGateway", params.Workspace.Name)

		err = suite.client.NetworkV1.DeleteInternetGateway(ctx, gatewayResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.InternetGateway.Name,
		}
		_, err = suite.client.NetworkV1.GetInternetGateway(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNetwork", params.Workspace.Name)

		err = suite.client.NetworkV1.DeleteNetwork(ctx, networkResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Network.Name,
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
		if actual.EgressOnly != nil {
			stepCtx.Require().Equal(expected.EgressOnly, *actual.EgressOnly, "EgressOnly should match expected")
		}
	})
}

func verifyRouteTableSpecStep(ctx provider.StepCtx, expected *secalib.RouteTableSpecV1, actual network.RouteTableSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.LocalRef, actual.LocalRef, "LocalRef should match expected")

		stepCtx.Require().Equal(len(expected.Routes), len(actual.Routes), "Route list length should match expected")
		for i := 0; i < len(expected.Routes); i++ {
			expectedRoute := expected.Routes[i]
			actualRoute := actual.Routes[i]
			stepCtx.Require().Equal(expectedRoute.DestinationCidrBlock, actualRoute.DestinationCidrBlock, fmt.Sprintf("Route [%d] DestinationCidrBlock should match expected", i))
			stepCtx.Require().Equal(expectedRoute.TargetRef, actualRoute.TargetRef, fmt.Sprintf("Route [%d] TargetRef should match expected", i))
		}
	})
}

func verifySubNetSpecStep(ctx provider.StepCtx, expected *secalib.SubnetSpecV1, actual network.SubnetSpec) {
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

func verifyPublicIpSpecStep(ctx provider.StepCtx, expected *secalib.PublicIpSpecV1, actual network.PublicIpSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Version, string(actual.Version), "Version should match expected")
		if actual.Address != nil {
			stepCtx.Require().Equal(expected.Address, *actual.Address, "Address should match expected")
		}
	})
}

func verifyNicSpecStep(ctx provider.StepCtx, expected *secalib.NICSpecV1, actual network.NicSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
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
			stepCtx.Require().Equal(expectedRule.Direction, string(actualRule.Direction), fmt.Sprintf("Rule [%d] Direction should match expected", i))
		}
	})
}
