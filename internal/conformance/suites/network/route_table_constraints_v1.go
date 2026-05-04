package network

import (
	"math/rand"
	"strings"

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

// RouteTableConstraintsValidationV1TestSuite verifies that RouteTable resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type RouteTableConstraintsValidationV1TestSuite struct {
	suites.RegionalTestSuite

	config *RouteTableLifeCycleV1Config
	params *params.RouteTableConstraintsValidationV1Params
}

func CreateRouteTableConstraintsValidationV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *RouteTableLifeCycleV1Config) *RouteTableConstraintsValidationV1TestSuite {
	suite := &RouteTableConstraintsValidationV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.RouteTableConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *RouteTableConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	workspaceName := generators.GenerateWorkspaceName()
	networkName := generators.GenerateNetworkName()
	internetGatewayName := generators.GenerateInternetGatewayName()
	routeTableName := generators.GenerateRouteTableName()
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName)
	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)
	routeTableRefObj := generators.GenerateRouteTableRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, routeTableName)

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Workspace for route table constraints violations testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	network, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Network for route table constraints testing"}).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: suite.config.NetworkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	internetGateway, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "InternetGateway for route table constraints testing"}).
		Spec(&schema.InternetGatewaySpec{EgressOnly: false}).Build()
	if err != nil {
		t.Fatalf("Failed to build InternetGateway: %v", err)
	}

	buildRT := func(name string, labels schema.Labels, annotations schema.Annotations) *schema.RouteTable {
		rt, err := builders.NewRouteTableBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
			Labels(labels).Annotations(annotations).
			Spec(&schema.RouteTableSpec{
				Routes: []schema.RouteSpec{
					{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
				},
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build RouteTable: %v", err)
		}
		return rt
	}

	p := &params.RouteTableConstraintsValidationV1Params{
		Workspace:       workspace,
		Network:         network,
		InternetGateway: internetGateway,
		OverLengthNameRouteTable: buildRT(strings.Repeat("a", 129),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "RouteTable with over-length name"}),
		InvalidPatternNameRouteTable: buildRT("Invalid-Name-With-Uppercase",
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "RouteTable with non-kebab-case name"}),
		OverLengthLabelValueRouteTable: buildRT(generators.GenerateRouteTableName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel, "constraint-test": strings.Repeat("x", 64)},
			schema.Annotations{"description": "RouteTable with over-length label value"}),
		OverLengthAnnotationRouteTable: buildRT(generators.GenerateRouteTableName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "RouteTable with over-length annotation value", "long-annotation": strings.Repeat("y", 1025)}),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureRouteTableConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *RouteTableConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable))
	suite.ConfigureDepends(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalResourceMetadataKindResourceKindInternetGateway),
	)
	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
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

	// Get the created Workspace
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

	// Network

	// Create a network
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

	// Internet gateway

	// Create an internet gateway
	internetGat := suite.params.InternetGateway
	expectInternetGatMeta := internetGat.Metadata
	expectInternetGatSpec := &internetGat.Spec
	expectInternetGatLabels := internetGat.Labels
	expectInternetGatAnnotations := internetGat.Annotations
	expectInternetGatExtensions := internetGat.Extensions
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create an internet gateway", t, suite.Client.NetworkV1, internetGat,
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
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.InternetGatewayStatus]{
			Metadata: expectInternetGatMeta,
			Spec:     expectInternetGatSpec,
			ResourceStatus: schema.InternetGatewayStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	stepsBuilder.GetNetworkV1Step("Get the created network", suite.Client.NetworkV1, networkWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			Metadata: network.Metadata,
			Spec:     &network.Spec,
			ResourceStatus: schema.NetworkStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	stepsBuilder.CreateOrUpdateRouteTableExpectViolationV1Step(
		"Create a route table with name exceeding maxLength:128 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthNameRouteTable,
	)
	stepsBuilder.CreateOrUpdateRouteTableExpectViolationV1Step(
		"Create a route table with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidPatternNameRouteTable,
	)
	stepsBuilder.CreateOrUpdateRouteTableExpectViolationV1Step(
		"Create a route table with label value exceeding maxLength:63 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthLabelValueRouteTable,
	)
	stepsBuilder.CreateOrUpdateRouteTableExpectViolationV1Step(
		"Create a route table with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthAnnotationRouteTable,
	)

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", t, suite.Client.NetworkV1, internetGat)
	stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", t, suite.Client.NetworkV1, internetGatWRef)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", t, suite.Client.NetworkV1, network)
	stepsBuilder.WatchNetworkUntilDeletedV1Step("Watch the network deletion", t, suite.Client.NetworkV1, networkWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *RouteTableConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
