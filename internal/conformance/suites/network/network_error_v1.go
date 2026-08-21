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

// NetworkErrorV1TestSuite verifies that Network resources with invalid references
// are rejected by the API with 422 Unprocessable Entity, and that operations conflicting
// with the current state of a network are rejected with 409 Conflict.
//
// Scenarios tested:
//   - Create network with invalid region
//   - Create network with invalid SKU
//   - Create network with non-existent workspace
//   - Delete a network that still has an active route table — expect 409 conflict
type NetworkErrorV1TestSuite struct {
	suites.RegionalTestSuite

	config *NetworkErrorV1Config
	params *params.NetworkErrorV1Params
}

type NetworkErrorV1Config struct {
	NetworkCidr string
	NetworkSkus []string
}

func CreateNetworkErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *NetworkErrorV1Config) *NetworkErrorV1TestSuite {
	suite := &NetworkErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.NetworkErrorV1SuiteName.String()
	return suite
}

func (suite *NetworkErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.NetworkParentSuite)

	workspaceName := generators.GenerateWorkspaceName()
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]
	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName)
	invalidSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, "non-existent-sku")
	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace for network error scenarios testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	buildNetwork := func(name string, workspaceRef string, region string, skuRef schema.Reference) *schema.Network {
		n, err := builders.NewNetworkBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceRef).Region(region).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "Network for error scenario testing"}).
			Spec(&schema.NetworkSpec{
				Cidr:   schema.Cidr{Ipv4: suite.config.NetworkCidr},
				SkuRef: skuRef,
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build Network: %v", err)
		}
		return n
	}

	// Network kept alive by an active route table, for the delete-conflict scenario
	networkName := generators.GenerateNetworkName()
	internetGatewayName := generators.GenerateInternetGatewayName()
	routeTableName := generators.GenerateRouteTableName()
	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)

	network := buildNetwork(networkName, workspaceName, suite.Region, *networkSkuRefObj)

	internetGateway, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Internet gateway for network delete conflict scenario testing"}).
		Spec(&schema.InternetGatewaySpec{EgressOnly: false}).Build()
	if err != nil {
		t.Fatalf("Failed to build InternetGateway: %v", err)
	}

	routeTable, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Route table keeping the network active for delete conflict scenario testing"}).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RouteTable: %v", err)
	}

	p := &params.NetworkErrorV1Params{
		Workspace: workspace,

		// Invalid region — random string, valid workspace + valid SKU
		InvalidRegionNetwork: buildNetwork(
			generators.GenerateNetworkName(),
			workspaceName,
			"invalid-region",
			*networkSkuRefObj,
		),

		// Invalid SKU — valid workspace + valid region, SKU does not exist
		InvalidSkuNetwork: buildNetwork(
			generators.GenerateNetworkName(),
			workspaceName,
			suite.Region,
			*invalidSkuRefObj,
		),

		// Non-existent workspace — workspace was never created
		NonExistentWorkspaceNetwork: buildNetwork(
			generators.GenerateNetworkName(),
			"non-existent-workspace",
			suite.Region,
			*networkSkuRefObj,
		),

		Network:         network,
		InternetGateway: internetGateway,
		RouteTable:      routeTable,
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mocknetwork.ConfigureNetworkErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *NetworkErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork))
	suite.ConfigureDepends(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
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

	// Error scenarios — all must be rejected with 422
	stepsBuilder.CreateOrUpdateNetworkExpectViolationV1Step(
		"Create a network with invalid region — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidRegionNetwork,
	)

	stepsBuilder.CreateOrUpdateNetworkExpectViolationV1Step(
		"Create a network with invalid SKU — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidSkuNetwork,
	)

	stepsBuilder.CreateOrUpdateNetworkExpectViolationV1Step(
		"Create a network with non-existent workspace — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentWorkspaceNetwork,
	)

	// Delete conflict scenario — network with an active route table must not be deletable
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

	internetGateway := suite.params.InternetGateway
	expectInternetGatewayMeta := internetGateway.Metadata
	expectInternetGatewaySpec := &internetGateway.Spec
	expectInternetGatewayLabels := internetGateway.Labels
	expectInternetGatewayAnnotations := internetGateway.Annotations
	expectInternetGatewayExtensions := internetGateway.Extensions

	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create an internet gateway", t, suite.Client.NetworkV1, internetGateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:         expectInternetGatewayLabels,
			Annotations:    expectInternetGatewayAnnotations,
			Extensions:     expectInternetGatewayExtensions,
			Metadata:       expectInternetGatewayMeta,
			Spec:           expectInternetGatewaySpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	internetGatewayWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(internetGateway.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(internetGateway.Metadata.Workspace),
		Name:      internetGateway.Metadata.Name,
	}
	stepsBuilder.GetInternetGatewayV1Step("Get the created internet gateway", suite.Client.NetworkV1, internetGatewayWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.InternetGatewayStatus]{
			Metadata: expectInternetGatewayMeta,
			Spec:     expectInternetGatewaySpec,
			ResourceStatus: schema.InternetGatewayStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	routeTable := suite.params.RouteTable
	expectRouteTableMeta := routeTable.Metadata
	expectRouteTableSpec := &routeTable.Spec
	expectRouteTableLabels := routeTable.Labels
	expectRouteTableAnnotations := routeTable.Annotations
	expectRouteTableExtensions := routeTable.Extensions

	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", t, suite.Client.NetworkV1, routeTable,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Labels:         expectRouteTableLabels,
			Annotations:    expectRouteTableAnnotations,
			Extensions:     expectRouteTableExtensions,
			Metadata:       expectRouteTableMeta,
			Spec:           expectRouteTableSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	routeTableNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(routeTable.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(routeTable.Metadata.Workspace),
		Network:   secapi.NetworkID(routeTable.Metadata.Network),
		Name:      routeTable.Metadata.Name,
	}
	stepsBuilder.GetRouteTableV1Step("Get the created route table", suite.Client.NetworkV1, routeTableNRef,
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			Metadata: expectRouteTableMeta,
			Spec:     expectRouteTableSpec,
			ResourceStatus: schema.RouteTableStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Delete the network while the route table is still active — expect 409 conflict
	stepsBuilder.DeleteNetworkExpectConflictV1Step(
		"Delete a network with an active route table — expect 409 conflict",
		suite.Client.NetworkV1,
		network,
	)

	// Teardown — reverse dependency order
	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", t, suite.Client.NetworkV1, routeTable)
	stepsBuilder.WatchRouteTableUntilDeletedV1Step("Watch the route table deletion", t, suite.Client.NetworkV1, routeTableNRef)

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", t, suite.Client.NetworkV1, internetGateway)
	stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", t, suite.Client.NetworkV1, internetGatewayWRef)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", t, suite.Client.NetworkV1, network)
	stepsBuilder.WatchNetworkUntilDeletedV1Step("Watch the network deletion", t, suite.Client.NetworkV1, networkWRef)

	// Teardown
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *NetworkErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
