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

// NetworkConstraintsValidationV1TestSuite verifies that Network resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type NetworkConstraintsValidationV1TestSuite struct {
	suites.RegionalTestSuite

	config *NetworkLifeCycleV1Config
	params *params.NetworkConstraintsValidationV1Params
}

func CreateNetworkConstraintsValidationV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *NetworkLifeCycleV1Config) *NetworkConstraintsValidationV1TestSuite {
	suite := &NetworkConstraintsValidationV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.NetworkConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *NetworkConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	workspaceName := generators.GenerateWorkspaceName()
	internetGatewayName := generators.GenerateInternetGatewayName()
	routeTableName := generators.GenerateRouteTableName()
	networkName := generators.GenerateNetworkName()
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName)
	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)
	routeTableRefObj := generators.GenerateRouteTableRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, routeTableName)

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Workspace for network constraints violations testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	internetGateway, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "InternetGateway for network constraints testing"}).
		Spec(&schema.InternetGatewaySpec{EgressOnly: false}).Build()
	if err != nil {
		t.Fatalf("Failed to build InternetGateway: %v", err)
	}

	_ = internetGatewayRefObj

	buildNetwork := func(name string, labels schema.Labels, annotations schema.Annotations) *schema.Network {
		n, err := builders.NewNetworkBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Labels(labels).Annotations(annotations).
			Spec(&schema.NetworkSpec{
				Cidr:          schema.Cidr{Ipv4: suite.config.NetworkCidr},
				SkuRef:        *networkSkuRefObj,
				RouteTableRef: *routeTableRefObj,
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build Network: %v", err)
		}
		return n
	}

	p := &params.NetworkConstraintsValidationV1Params{
		Workspace:       workspace,
		InternetGateway: internetGateway,
		OverLengthNameNetwork: buildNetwork(strings.Repeat("a", 129),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Network with over-length name"}),
		InvalidPatternNameNetwork: buildNetwork("Invalid-Name-With-Uppercase",
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Network with non-kebab-case name"}),
		OverLengthLabelValueNetwork: buildNetwork(generators.GenerateNetworkName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel, "constraint-test": strings.Repeat("x", 64)},
			schema.Annotations{"description": "Network with over-length label value"}),
		OverLengthAnnotationNetwork: buildNetwork(generators.GenerateNetworkName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Network with over-length annotation value", "long-annotation": strings.Repeat("y", 1025)}),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureNetworkConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *NetworkConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.NetworkProviderV1Name, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	workspace := suite.params.Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   workspace.Metadata.Name,
	}
	internetGateway := suite.params.InternetGateway
	gatewayWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(internetGateway.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(internetGateway.Metadata.Workspace),
		Name:      internetGateway.Metadata.Name,
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

	stepsBuilder.CreateOrUpdateNetworkExpectViolationV1Step(
		"Create a network with name exceeding maxLength:128 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthNameNetwork,
	)
	stepsBuilder.CreateOrUpdateNetworkExpectViolationV1Step(
		"Create a network with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidPatternNameNetwork,
	)
	stepsBuilder.CreateOrUpdateNetworkExpectViolationV1Step(
		"Create a network with label value exceeding maxLength:63 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthLabelValueNetwork,
	)
	stepsBuilder.CreateOrUpdateNetworkExpectViolationV1Step(
		"Create a network with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthAnnotationNetwork,
	)

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, internetGateway)
	stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", suite.Client.NetworkV1, gatewayWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *NetworkConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
