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
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"k8s.io/utils/ptr"
)

type NetworkLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	config *NetworkLifeCycleV1Config
	params *params.NetworkLifeCycleV1Params
}

type NetworkLifeCycleV1Config struct {
	NetworkCidr string
	NetworkSkus []string
}

func CreateNetworkLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *NetworkLifeCycleV1Config) *NetworkLifeCycleV1TestSuite {
	suite := &NetworkLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.NetworkLifeCycleV1SuiteName.String()
	return suite
}

func (suite *NetworkLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Network")

	workspaceName := generators.GenerateWorkspaceName()
	networkName := generators.GenerateNetworkName()
	routeTableName := generators.GenerateRouteTableName()
	internetGatewayName := generators.GenerateInternetGatewayName()

	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]
	networkSkuName2 := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	routeTableRefObj := generators.GenerateRouteTableRefObject(routeTableName)

	networkSkuRefObj := generators.GenerateSkuRefObject(networkSkuName)
	networkSkuRef2Obj := generators.GenerateSkuRefObject(networkSkuName2)

	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(internetGatewayName)

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	networkInitial, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: ptr.To(suite.config.NetworkCidr)},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	networkUpdated, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: ptr.To(suite.config.NetworkCidr)},
			SkuRef:        *networkSkuRef2Obj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	routeTable, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: ptr.To(false),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	params := &params.NetworkLifeCycleV1Params{
		Workspace:       workspace,
		RouteTable:      routeTable,
		InternetGateway: internetGat,
		NetworkInitial:  networkInitial,
		NetworkUpdated:  networkUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureNetworkLifecycleScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *NetworkLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.NetworkProviderV1Name,
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
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Network

	// Create a network
	network := suite.params.NetworkInitial
	expectNetworkMeta := network.Metadata
	expectNetworkSpec := &network.Spec
	stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	networkWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(network.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
		Name:      network.Metadata.Name,
	}

	// Internet gateway

	// Create an internet gateway
	internetGat := suite.params.InternetGateway
	expectInternetGatMeta := internetGat.Metadata
	expectInternetGatSpec := &internetGat.Spec
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create an internet gateway", suite.Client.NetworkV1, internetGat,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:       expectInternetGatMeta,
			Spec:           expectInternetGatSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
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
			Metadata:       expectInternetGatMeta,
			Spec:           expectInternetGatSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Route table

	// Create a route table
	route := suite.params.RouteTable
	expectRouteMeta := route.Metadata
	expectRouteSpec := &route.Spec
	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
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
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the created network
	stepsBuilder.GetNetworkV1Step("Get the created network", suite.Client.NetworkV1, networkWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the network
	network.Spec = suite.params.NetworkUpdated.Spec
	expectNetworkSpec.SkuRef = network.Spec.SkuRef
	stepsBuilder.CreateOrUpdateNetworkV1Step("Update the network", suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated network
	stepsBuilder.GetNetworkV1Step("Get the updated network", suite.Client.NetworkV1, networkWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Resources deletion

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, internetGat)
	stepsBuilder.GetInternetGatewayWithErrorV1Step("Get deleted internet gateway", suite.Client.NetworkV1, internetGatWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", suite.Client.NetworkV1, route)
	stepsBuilder.GetRouteTableWithErrorV1Step("Get deleted route table", suite.Client.NetworkV1, routeNRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", suite.Client.NetworkV1, network)
	stepsBuilder.GetNetworkWithErrorV1Step("Get deleted network", suite.Client.NetworkV1, networkWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)
}

func (suite *NetworkLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
