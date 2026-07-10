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

// RouteTableErrorV1TestSuite verifies that RouteTable resources with invalid references
// are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create route table with invalid region
//   - Create route table with non-existent workspace
//   - Create route table with non-existent network
//   - Create route table with non-existent targetRef
type RouteTableErrorV1TestSuite struct {
	suites.RegionalTestSuite

	config *RouteTableErrorV1Config
	params *params.RouteTableErrorV1Params
}

type RouteTableErrorV1Config struct {
	NetworkCidr string
	NetworkSkus []string
}

func CreateRouteTableErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *RouteTableErrorV1Config) *RouteTableErrorV1TestSuite {
	suite := &RouteTableErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.RouteTableErrorV1SuiteName.String()
	return suite
}

func (suite *RouteTableErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.NetworkParentSuite)

	workspaceName := generators.GenerateWorkspaceName()
	networkName := generators.GenerateNetworkName()
	internetGatewayName := generators.GenerateInternetGatewayName()
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]
	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName)
	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)
	nonExistentIgRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, "non-existent-ig")
	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace for route table error scenarios testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	network, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Network for route table error scenarios testing"}).
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
		Annotations(schema.Annotations{"description": "Internet Gateway for route table error scenarios testing"}).
		Spec(&schema.InternetGatewaySpec{EgressOnly: false}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build InternetGateway: %v", err)
	}

	buildRouteTable := func(name string, workspaceRef string, networkRef string, region string, targetRef schema.Reference) *schema.RouteTable {
		rt, err := builders.NewRouteTableBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceRef).Region(region).Network(networkRef).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "RouteTable for error scenario testing"}).
			Spec(&schema.RouteTableSpec{
				Routes: []schema.RouteSpec{
					{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: targetRef},
				},
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build RouteTable: %v", err)
		}
		return rt
	}

	p := &params.RouteTableErrorV1Params{
		Workspace:       workspace,
		Network:         network,
		InternetGateway: internetGateway,

		// Invalid region — random string, valid workspace + network + targetRef
		InvalidRegionRouteTable: buildRouteTable(
			generators.GenerateRouteTableName(),
			workspaceName,
			networkName,
			"invalid-region",
			*internetGatewayRefObj,
		),

		// Non-existent workspace — workspace + network were never created
		NonExistentWorkspaceRouteTable: buildRouteTable(
			generators.GenerateRouteTableName(),
			"non-existent-workspace",
			"non-existent-network",
			suite.Region,
			*internetGatewayRefObj,
		),

		// Non-existent network — valid workspace, network does not exist
		NonExistentNetworkRouteTable: buildRouteTable(
			generators.GenerateRouteTableName(),
			workspaceName,
			"non-existent-network",
			suite.Region,
			*internetGatewayRefObj,
		),

		// Non-existent targetRef — valid workspace + network, IG does not exist
		NonExistentTargetRefRouteTable: buildRouteTable(
			generators.GenerateRouteTableName(),
			workspaceName,
			networkName,
			suite.Region,
			*nonExistentIgRefObj,
		),
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mocknetwork.ConfigureRouteTableErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *RouteTableErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable))
	suite.ConfigureDepends(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalResourceMetadataKindResourceKindInternetGateway),
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
	internetGateway := suite.params.InternetGateway
	expectIgMeta := internetGateway.Metadata
	expectIgSpec := &internetGateway.Spec
	expectIgLabels := internetGateway.Labels
	expectIgAnnotations := internetGateway.Annotations
	expectIgExtensions := internetGateway.Extensions

	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create an internet gateway", t, suite.Client.NetworkV1, internetGateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:         expectIgLabels,
			Annotations:    expectIgAnnotations,
			Extensions:     expectIgExtensions,
			Metadata:       expectIgMeta,
			Spec:           expectIgSpec,
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
			Metadata: expectIgMeta,
			Spec:     expectIgSpec,
			ResourceStatus: schema.InternetGatewayStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Error scenarios — all must be rejected with 422
	stepsBuilder.CreateOrUpdateRouteTableExpectViolationV1Step(
		"Create a route table with invalid region — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidRegionRouteTable,
	)

	stepsBuilder.CreateOrUpdateRouteTableExpectViolationV1Step(
		"Create a route table with non-existent workspace — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentWorkspaceRouteTable,
	)

	stepsBuilder.CreateOrUpdateRouteTableExpectViolationV1Step(
		"Create a route table with non-existent network — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentNetworkRouteTable,
	)

	stepsBuilder.CreateOrUpdateRouteTableExpectViolationV1Step(
		"Create a route table with non-existent targetRef — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentTargetRefRouteTable,
	)

	// Teardown — reverse dependency order
	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", t, suite.Client.NetworkV1, internetGateway)
	stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", t, suite.Client.NetworkV1, internetGatewayWRef)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", t, suite.Client.NetworkV1, network)
	stepsBuilder.WatchNetworkUntilDeletedV1Step("Watch the network deletion", t, suite.Client.NetworkV1, networkWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *RouteTableErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
