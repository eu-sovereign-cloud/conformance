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
)

type RouteTableLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	config *RouteTableLifeCycleV1Config
	params *params.RouteTableLifeCycleV1Params
}

func CreateRouteTableLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *RouteTableLifeCycleV1Config) *RouteTableLifeCycleV1TestSuite {
	suite := &RouteTableLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.RouteTableLifeCycleV1SuiteName.String()
	return suite
}

type RouteTableLifeCycleV1Config struct {
	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	NetworkSkus    []string
}

func (suite *RouteTableLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Network")

	workspaceName := generators.GenerateWorkspaceName()
	networkName := generators.GenerateNetworkName()
	routeTableName := generators.GenerateRouteTableName()
	internetGatewayName := generators.GenerateInternetGatewayName()
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	instanceName := generators.GenerateInstanceName()

	instanceRefObj := generators.GenerateInstanceRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, workspaceName, instanceName)

	routeTableRefObj := generators.GenerateRouteTableRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, routeTableName)

	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName)

	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Workspace for conformance testing",
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	network, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Network for conformance testing",
		}).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: suite.config.NetworkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}
	routeTableInitial, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Route Table for conformance testing",
		}).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Route Table: %v", err)
	}

	routeTableUpdated, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Route Table for conformance testing",
		}).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *instanceRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Route Table: %v", err)
	}

	internetGat, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Internet Gateway for conformance testing",
		}).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: false,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	params := &params.RouteTableLifeCycleV1Params{
		Workspace:         workspace,
		Network:           network,
		InternetGateway:   internetGat,
		RouteTableInitial: routeTableInitial,
		RouteTableUpdated: routeTableUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureRouteTableLifecycleScenarioV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *RouteTableLifeCycleV1TestSuite) TestScenario(t provider.T) {
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
	expectWorkspaceAnnotations := workspace.Annotations
	expectWorkspaceExtensions := workspace.Extensions
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Annotations:    expectWorkspaceAnnotations,
			Extensions:     expectWorkspaceExtensions,
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
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:      expectWorkspaceLabels,
			Annotations: expectWorkspaceAnnotations,
			Extensions:  expectWorkspaceExtensions,
			Metadata:    expectWorkspaceMeta,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Network

	// Create a network
	network := suite.params.Network
	expectNetworkMeta := network.Metadata
	expectNetworkSpec := &network.Spec
	expectNetworkLabels := network.Labels
	expectNetworkAnnotations := network.Annotations
	expectNetworkExtensions := network.Extensions
	stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Labels:         expectNetworkLabels,
			Annotations:    expectNetworkAnnotations,
			Extensions:     expectNetworkExtensions,
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Internet gateway

	// Create an internet gateway
	internetGat := suite.params.InternetGateway
	expectInternetGatMeta := internetGat.Metadata
	expectInternetGatSpec := &internetGat.Spec
	expectInternetGatLabels := internetGat.Labels
	expectInternetGatAnnotations := internetGat.Annotations
	expectInternetGatExtensions := internetGat.Extensions
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create an internet gateway", suite.Client.NetworkV1, internetGat,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:         expectInternetGatLabels,
			Annotations:    expectInternetGatAnnotations,
			Extensions:     expectInternetGatExtensions,
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
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:      expectInternetGatLabels,
			Annotations: expectInternetGatAnnotations,
			Extensions:  expectInternetGatExtensions,
			Metadata:    expectInternetGatMeta,
			Spec:        expectInternetGatSpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Route table

	// Create a route table
	route := suite.params.RouteTableInitial
	expectRouteMeta := route.Metadata
	expectRouteSpec := &route.Spec
	expectRouteLabels := route.Labels
	expectRouteAnnotations := route.Annotations
	expectRouteExtensions := route.Extensions
	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Labels:         expectRouteLabels,
			Annotations:    expectRouteAnnotations,
			Extensions:     expectRouteExtensions,
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
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Labels:      expectRouteLabels,
			Annotations: expectRouteAnnotations,
			Extensions:  expectRouteExtensions,
			Metadata:    expectRouteMeta,
			Spec:        expectRouteSpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Get the created network
	networkWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(network.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
		Name:      network.Metadata.Name,
	}
	stepsBuilder.GetNetworkV1Step("Get the created network", suite.Client.NetworkV1, networkWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Labels:      expectNetworkLabels,
			Annotations: expectNetworkAnnotations,
			Extensions:  expectNetworkExtensions,
			Metadata:    expectNetworkMeta,
			Spec:        expectNetworkSpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Update the route table
	route = suite.params.RouteTableUpdated
	expectRouteSpec.Routes = route.Spec.Routes
	expectRouteLabels = route.Labels
	expectRouteAnnotations = route.Annotations
	expectRouteExtensions = route.Extensions
	stepsBuilder.CreateOrUpdateRouteTableV1Step("Update the route table", suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Labels:         expectRouteLabels,
			Annotations:    expectRouteAnnotations,
			Extensions:     expectRouteExtensions,
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: suites.UpdatedResourceExpectedStates,
		},
	)

	// Get the updated route table
	stepsBuilder.GetRouteTableV1Step("Get the updated route table", suite.Client.NetworkV1, routeNRef,
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Labels:      expectRouteLabels,
			Annotations: expectRouteAnnotations,
			Extensions:  expectRouteExtensions,
			Metadata:    expectRouteMeta,
			Spec:        expectRouteSpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)

	// Resources deletion

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, internetGat)
	stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", suite.Client.NetworkV1, internetGatWRef)

	// Resources deletion
	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", suite.Client.NetworkV1, route)
	stepsBuilder.WatchRouteTableUntilDeletedV1Step("Watch the route table deletion", suite.Client.NetworkV1, routeNRef)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", suite.Client.NetworkV1, network)
	stepsBuilder.WatchNetworkUntilDeletedV1Step("Watch the network deletion", suite.Client.NetworkV1, networkWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *RouteTableLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
