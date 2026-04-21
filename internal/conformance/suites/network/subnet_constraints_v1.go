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

// SubnetConstraintsValidationV1TestSuite verifies that Subnet resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type SubnetConstraintsValidationV1TestSuite struct {
	suites.RegionalTestSuite

	config *SubnetLifeCycleV1Config
	params *params.SubnetConstraintsValidationV1Params
}

func CreateSubnetConstraintsValidationV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *SubnetLifeCycleV1Config) *SubnetConstraintsValidationV1TestSuite {
	suite := &SubnetConstraintsValidationV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.SubnetConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *SubnetConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

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

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Workspace for subnet constraints violations testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	network, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Network for subnet constraints testing"}).
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
		Annotations(schema.Annotations{"description": "InternetGateway for subnet constraints testing"}).
		Spec(&schema.InternetGatewaySpec{EgressOnly: false}).Build()
	if err != nil {
		t.Fatalf("Failed to build InternetGateway: %v", err)
	}

	routeTable, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "RouteTable for subnet constraints testing"}).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RouteTable: %v", err)
	}

	buildSubnet := func(name string, labels schema.Labels, annotations schema.Annotations) *schema.Subnet {
		s, err := builders.NewSubnetBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
			Labels(labels).Annotations(annotations).
			Spec(&schema.SubnetSpec{
				Cidr: schema.Cidr{Ipv4: subnetCidr},
				Zone: zone,
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build Subnet: %v", err)
		}
		return s
	}

	p := &params.SubnetConstraintsValidationV1Params{
		Workspace:       workspace,
		Network:         network,
		InternetGateway: internetGateway,
		RouteTable:      routeTable,
		OverLengthNameSubnet: buildSubnet(strings.Repeat("a", 129),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Subnet with over-length name"}),
		InvalidPatternNameSubnet: buildSubnet("Invalid-Name-With-Uppercase",
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Subnet with non-kebab-case name"}),
		OverLengthLabelValueSubnet: buildSubnet(generators.GenerateSubnetName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel, "constraint-test": strings.Repeat("x", 64)},
			schema.Annotations{"description": "Subnet with over-length label value"}),
		OverLengthAnnotationSubnet: buildSubnet(generators.GenerateSubnetName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Subnet with over-length annotation value", "long-annotation": strings.Repeat("y", 1025)}),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureSubnetConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *SubnetConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.NetworkProviderV1Name, string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	workspace := suite.params.Workspace
	network := suite.params.Network
	internetGateway := suite.params.InternetGateway
	routeTable := suite.params.RouteTable

	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   workspace.Metadata.Name,
	}
	networkWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(network.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
		Name:      network.Metadata.Name,
	}
	gatewayWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(internetGateway.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(internetGateway.Metadata.Workspace),
		Name:      internetGateway.Metadata.Name,
	}
	routeNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(routeTable.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(routeTable.Metadata.Workspace),
		Network:   secapi.NetworkID(routeTable.Metadata.Network),
		Name:      routeTable.Metadata.Name,
	}

	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create workspace for test environment", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         workspace.Labels,
			Annotations:    workspace.Annotations,
			Metadata:       workspace.Metadata,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			Labels:   workspace.Labels,
			Metadata: workspace.Metadata,
			ResourceStatus: schema.WorkspaceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	stepsBuilder.CreateOrUpdateNetworkV1Step("Create network for test environment", suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Labels:         network.Labels,
			Annotations:    network.Annotations,
			Metadata:       network.Metadata,
			Spec:           &network.Spec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create internet gateway for test environment", suite.Client.NetworkV1, internetGateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:         internetGateway.Labels,
			Annotations:    internetGateway.Annotations,
			Metadata:       internetGateway.Metadata,
			Spec:           &internetGateway.Spec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetInternetGatewayV1Step("Get the created internet gateway", suite.Client.NetworkV1, gatewayWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.InternetGatewayStatus]{
			Metadata: internetGateway.Metadata,
			Spec:     &internetGateway.Spec,
			ResourceStatus: schema.InternetGatewayStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create route table for test environment", suite.Client.NetworkV1, routeTable,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Labels:         routeTable.Labels,
			Annotations:    routeTable.Annotations,
			Metadata:       routeTable.Metadata,
			Spec:           &routeTable.Spec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetRouteTableV1Step("Get the created route table", suite.Client.NetworkV1, routeNRef,
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			Metadata: routeTable.Metadata,
			Spec:     &routeTable.Spec,
			ResourceStatus: schema.RouteTableStatus{
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

	stepsBuilder.CreateOrUpdateSubnetExpectViolationV1Step(
		"Create a subnet with name exceeding maxLength:128 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthNameSubnet,
	)
	stepsBuilder.CreateOrUpdateSubnetExpectViolationV1Step(
		"Create a subnet with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidPatternNameSubnet,
	)
	stepsBuilder.CreateOrUpdateSubnetExpectViolationV1Step(
		"Create a subnet with label value exceeding maxLength:63 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthLabelValueSubnet,
	)
	stepsBuilder.CreateOrUpdateSubnetExpectViolationV1Step(
		"Create a subnet with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthAnnotationSubnet,
	)

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, internetGateway)
	stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", suite.Client.NetworkV1, gatewayWRef)

	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", suite.Client.NetworkV1, routeTable)
	stepsBuilder.WatchRouteTableUntilDeletedV1Step("Watch the route table deletion", suite.Client.NetworkV1, routeNRef)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", suite.Client.NetworkV1, network)
	stepsBuilder.WatchNetworkUntilDeletedV1Step("Watch the network deletion", suite.Client.NetworkV1, networkWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *SubnetConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
