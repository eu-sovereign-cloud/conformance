package network

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mocknetwork "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/network"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// SubnetErrorV1TestSuite verifies that Subnet resources with invalid references
// are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create subnet with invalid region
//   - Create subnet with invalid zone
//   - Create subnet with non-existent workspace
//   - Create subnet with non-existent network
//   - Create subnet with non-existent routeTableRef
//   - Create subnet with CIDR outside network CIDR
type SubnetErrorV1TestSuite struct {
	suites.RegionalTestSuite

	config *SubnetLifeCycleV1Config
	params *params.SubnetErrorV1Params
}

func CreateSubnetErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *SubnetLifeCycleV1Config) *SubnetErrorV1TestSuite {
	suite := &SubnetErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.SubnetErrorV1SuiteName.String()
	return suite
}

func (suite *SubnetErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.NetworkParentSuite)

	workspaceName := generators.GenerateWorkspaceName()
	networkName := generators.GenerateNetworkName()
	internetGatewayName := generators.GenerateInternetGatewayName()
	routeTableName := generators.GenerateRouteTableName()
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]
	zone := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]

	subnetCidr, err := generators.GenerateSubnetCidr(suite.config.NetworkCidr, 8, 1)
	if err != nil {
		t.Fatalf("Failed to generate subnet cidr: %v", err)
	}

	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName)
	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)
	routeTableRefObj := generators.GenerateRouteTableRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, routeTableName)
	nonExistentRouteTableRefObj := generators.GenerateRouteTableRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, "non-existent-rt")
	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace for subnet error scenarios testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	network, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Network for subnet error scenarios testing"}).
		Spec(&schema.NetworkSpec{
			Cidr:   schema.Cidr{Ipv4: suite.config.NetworkCidr},
			SkuRef: *networkSkuRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	internetGateway, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "InternetGateway for subnet error scenarios testing"}).
		Spec(&schema.InternetGatewaySpec{EgressOnly: false}).Build()
	if err != nil {
		t.Fatalf("Failed to build InternetGateway: %v", err)
	}

	routeTable, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "RouteTable for subnet error scenarios testing"}).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RouteTable: %v", err)
	}

	buildSubnet := func(name string, workspaceRef string, networkRef string, region string, zone string, cidr string, routeTableRef schema.Reference) *schema.Subnet {
		s, err := builders.NewSubnetBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceRef).Region(region).Network(networkRef).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "Subnet for error scenario testing"}).
			Spec(&schema.SubnetSpec{
				Cidr:          schema.Cidr{Ipv4: cidr},
				RouteTableRef: routeTableRef,
				Zone:          zone,
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build Subnet: %v", err)
		}
		return s
	}

	p := &params.SubnetErrorV1Params{
		Workspace:       workspace,
		Network:         network,
		InternetGateway: internetGateway,
		RouteTable:      routeTable,

		// Invalid region — random string, everything else valid
		InvalidRegionSubnet: buildSubnet(
			generators.GenerateSubnetName(),
			workspaceName,
			networkName,
			"invalid-region",
			zone,
			subnetCidr,
			*routeTableRefObj,
		),

		// Invalid zone — random string not in provider zones
		InvalidZoneSubnet: buildSubnet(
			generators.GenerateSubnetName(),
			workspaceName,
			networkName,
			suite.Region,
			"invalid-zone",
			subnetCidr,
			*routeTableRefObj,
		),

		// Non-existent workspace — workspace + network were never created
		NonExistentWorkspaceSubnet: buildSubnet(
			generators.GenerateSubnetName(),
			"non-existent-workspace",
			"non-existent-network",
			suite.Region,
			zone,
			subnetCidr,
			*routeTableRefObj,
		),

		// Non-existent network — valid workspace, network does not exist
		NonExistentNetworkSubnet: buildSubnet(
			generators.GenerateSubnetName(),
			workspaceName,
			"non-existent-network",
			suite.Region,
			zone,
			subnetCidr,
			*routeTableRefObj,
		),

		// Non-existent routeTableRef — valid workspace + network, route table does not exist
		NonExistentRouteTableRefSubnet: buildSubnet(
			generators.GenerateSubnetName(),
			workspaceName,
			networkName,
			suite.Region,
			zone,
			subnetCidr,
			*nonExistentRouteTableRefObj,
		),

		// CIDR outside network CIDR — subnet CIDR not contained within network CIDR
		OutsideCidrSubnet: buildSubnet(
			generators.GenerateSubnetName(),
			workspaceName,
			networkName,
			suite.Region,
			zone,
			"192.168.99.0/24",
			*routeTableRefObj,
		),
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mocknetwork.ConfigureSubnetErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *SubnetErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSubnet))
	suite.ConfigureDepends(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace setup
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	expectWorkspaceAnnotations := workspace.Annotations
	expectWorkspaceExtensions := workspace.Extensions

	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Annotations:    expectWorkspaceAnnotations,
			Extensions:     expectWorkspaceExtensions,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			Labels:   expectWorkspaceLabels,
			Metadata: expectWorkspaceMeta,
			ResourceStatus: schema.WorkspaceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Network setup
	network := suite.params.Network
	expectNetworkMeta := network.Metadata
	expectNetworkSpec := &network.Spec
	expectNetworkLabels := network.Labels
	expectNetworkAnnotations := network.Annotations
	expectNetworkExtensions := network.Extensions

	stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", t, suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Labels:         expectNetworkLabels,
			Annotations:    expectNetworkAnnotations,
			Extensions:     expectNetworkExtensions,
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

	stepsBuilder.GetNetworkV1Step("Get the created network", suite.Client.NetworkV1, networkWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			Metadata: expectNetworkMeta,
			Spec:     expectNetworkSpec,
			ResourceStatus: schema.NetworkStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Internet gateway setup
	internetGat := suite.params.InternetGateway
	expectIgMeta := internetGat.Metadata
	expectIgSpec := &internetGat.Spec
	expectIgLabels := internetGat.Labels
	expectIgAnnotations := internetGat.Annotations
	expectIgExtensions := internetGat.Extensions

	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create an internet gateway", t, suite.Client.NetworkV1, internetGat,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:         expectIgLabels,
			Annotations:    expectIgAnnotations,
			Extensions:     expectIgExtensions,
			Metadata:       expectIgMeta,
			Spec:           expectIgSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	internetGatWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(internetGat.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(internetGat.Metadata.Workspace),
		Name:      internetGat.Metadata.Name,
	}
	stepsBuilder.GetInternetGatewayV1Step("Get the created internet gateway", suite.Client.NetworkV1, internetGatWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.InternetGatewayStatus]{
			Metadata: expectIgMeta,
			Spec:     expectIgSpec,
			ResourceStatus: schema.InternetGatewayStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Route table setup
	route := suite.params.RouteTable
	expectRouteMeta := route.Metadata
	expectRouteSpec := &route.Spec
	expectRouteLabels := route.Labels
	expectRouteAnnotations := route.Annotations
	expectRouteExtensions := route.Extensions

	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", t, suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Labels:         expectRouteLabels,
			Annotations:    expectRouteAnnotations,
			Extensions:     expectRouteExtensions,
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	routeNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(route.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(route.Metadata.Workspace),
		Network:   secapi.NetworkID(route.Metadata.Network),
		Name:      route.Metadata.Name,
	}
	stepsBuilder.GetRouteTableV1Step("Get the created route table", suite.Client.NetworkV1, routeNRef,
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			Metadata: expectRouteMeta,
			Spec:     expectRouteSpec,
			ResourceStatus: schema.RouteTableStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Error scenarios — all must be rejected with 422
	stepsBuilder.CreateOrUpdateSubnetExpectViolationV1Step(
		"Create a subnet with invalid region — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidRegionSubnet,
	)

	stepsBuilder.CreateOrUpdateSubnetExpectViolationV1Step(
		"Create a subnet with invalid zone — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidZoneSubnet,
	)

	stepsBuilder.CreateOrUpdateSubnetExpectViolationV1Step(
		"Create a subnet with non-existent workspace — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentWorkspaceSubnet,
	)

	stepsBuilder.CreateOrUpdateSubnetExpectViolationV1Step(
		"Create a subnet with non-existent network — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentNetworkSubnet,
	)

	stepsBuilder.CreateOrUpdateSubnetExpectViolationV1Step(
		"Create a subnet with non-existent routeTableRef — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentRouteTableRefSubnet,
	)

	stepsBuilder.CreateOrUpdateSubnetExpectViolationV1Step(
		"Create a subnet with CIDR outside network CIDR — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OutsideCidrSubnet,
	)

	// Teardown — reverse dependency order
	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", t, suite.Client.NetworkV1, route)
	stepsBuilder.WatchRouteTableUntilDeletedV1Step("Watch the route table deletion", t, suite.Client.NetworkV1, routeNRef)

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", t, suite.Client.NetworkV1, internetGat)
	stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", t, suite.Client.NetworkV1, internetGatWRef)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", t, suite.Client.NetworkV1, network)
	stepsBuilder.WatchNetworkUntilDeletedV1Step("Watch the network deletion", t, suite.Client.NetworkV1, networkWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *SubnetErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
