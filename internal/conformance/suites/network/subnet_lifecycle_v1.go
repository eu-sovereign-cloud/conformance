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
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"k8s.io/utils/ptr"
)

type SubnetLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	config *SubnetLifeCycleV1Config
	params *params.SubnetLifeCycleV1Params
}

type SubnetLifeCycleV1Config struct {
	NetworkCidr string
	RegionZones []string
	NetworkSkus []string
}

func CreateSubnetLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *SubnetLifeCycleV1Config) *SubnetLifeCycleV1TestSuite {
	suite := &SubnetLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.SubnetLifeCycleV1SuiteName.String()
	return suite
}

func (suite *SubnetLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Subnet")

	workspaceName := generators.GenerateWorkspaceName()
	subnetName := generators.GenerateSubnetName()
	networkName := generators.GenerateNetworkName()
	routeTableName := generators.GenerateRouteTableName()
	internetGatewayName := generators.GenerateInternetGatewayName()
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	// Generate the subnet cidr
	subnetCidr, err := generators.GenerateSubnetCidr(suite.config.NetworkCidr, 8, 1)
	if err != nil {
		t.Fatalf("Failed to generate subnet cidr: %v", err)
	}

	networkSkuRefObj, err := generators.GenerateSkuRefObject(networkSkuName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	routeTableRefObj, err := generators.GenerateRouteTableRefObject(routeTableName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	internetGatewayRefObj, err := generators.GenerateInternetGatewayRefObject(internetGatewayName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	// Select zones
	zone1 := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]
	zone2 := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	network, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: ptr.To(suite.config.NetworkCidr)},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	subnetInitial, err := builders.NewSubnetBuilder().
		Name(subnetName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone1,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	subnetUpdated, err := builders.NewSubnetBuilder().
		Name(subnetName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone2,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	routeTable, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Route Table: %v", err)
	}

	internetGat, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: ptr.To(false),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	params := &params.SubnetLifeCycleV1Params{
		Workspace:       workspace,
		Network:         network,
		RouteTable:      routeTable,
		InternetGateway: internetGat,
		SubnetInitial:   subnetInitial,
		SubnetUpdated:   subnetUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureSubnetLifecycleScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *SubnetLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.NetworkProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Network

	// Create a network
	network := suite.params.Network
	expectNetworkMeta := network.Metadata
	expectNetworkSpec := &network.Spec
	stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:      expectNetworkMeta,
			Spec:          expectNetworkSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created network
	networkWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(network.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
		Name:      network.Metadata.Name,
	}
	stepsBuilder.GetNetworkV1Step("Get the created network", suite.Client.NetworkV1, networkWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:      expectNetworkMeta,
			Spec:          expectNetworkSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Route table

	// Create a route table
	route := suite.params.RouteTable
	expectRouteMeta := route.Metadata
	expectRouteSpec := &route.Spec
	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:      expectRouteMeta,
			Spec:          expectRouteSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created route table
	routeNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(route.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(route.Metadata.Workspace),
		Network:   secapi.NetworkID(route.Metadata.Network),
		Name:      route.Metadata.Name,
	}
	stepsBuilder.GetRouteTableV1Step("Get the created route table", suite.Client.NetworkV1, routeNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:      expectRouteMeta,
			Spec:          expectRouteSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)
	// Internet gateway

	// Create an internet gateway
	internetGat := suite.params.InternetGateway
	expectInternetGatMeta := internetGat.Metadata
	expectInternetGatSpec := &internetGat.Spec
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create an internet gateway", suite.Client.NetworkV1, internetGat,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:      expectInternetGatMeta,
			Spec:          expectInternetGatSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)
	// Get the created internet gateway
	internetGatWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(internetGat.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(internetGat.Metadata.Workspace),
		Name:      internetGat.Metadata.Name,
	}
	stepsBuilder.GetInternetGatewayV1Step("Get the created internet gateway", suite.Client.NetworkV1, internetGatWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:      expectInternetGatMeta,
			Spec:          expectInternetGatSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Subnet

	// Create a subnet
	subnet := suite.params.SubnetInitial
	expectSubnetMeta := subnet.Metadata
	expectSubnetSpec := &subnet.Spec
	stepsBuilder.CreateOrUpdateSubnetV1Step("Create a subnet", suite.Client.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:      expectSubnetMeta,
			Spec:          expectSubnetSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created subnet
	subnetNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(subnet.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(subnet.Metadata.Workspace),
		Network:   secapi.NetworkID(subnet.Metadata.Network),
		Name:      subnet.Metadata.Name,
	}
	stepsBuilder.GetSubnetV1Step("Get the created subnet", suite.Client.NetworkV1, subnetNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:      expectSubnetMeta,
			Spec:          expectSubnetSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the subnet
	subnet.Spec = suite.params.SubnetUpdated.Spec
	expectSubnetSpec.Zone = subnet.Spec.Zone
	stepsBuilder.CreateOrUpdateSubnetV1Step("Update the subnet", suite.Client.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:      expectSubnetMeta,
			Spec:          expectSubnetSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated subnet
	stepsBuilder.GetSubnetV1Step("Get the updated subnet", suite.Client.NetworkV1, subnetNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:      expectSubnetMeta,
			Spec:          expectSubnetSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion
	stepsBuilder.DeleteSubnetV1Step("Delete the subnet", suite.Client.NetworkV1, subnet)
	stepsBuilder.GetSubnetWithErrorV1Step("Get deleted subnet", suite.Client.NetworkV1, subnetNRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, internetGat)
	stepsBuilder.GetInternetGatewayWithErrorV1Step("Get deleted internet gateway", suite.Client.NetworkV1, internetGatWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", suite.Client.NetworkV1, route)
	stepsBuilder.GetRouteTableWithErrorV1Step("Get deleted route table", suite.Client.NetworkV1, routeNRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", suite.Client.NetworkV1, network)
	stepsBuilder.GetNetworkWithErrorV1Step("Get deleted network", suite.Client.NetworkV1, networkWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *SubnetLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
