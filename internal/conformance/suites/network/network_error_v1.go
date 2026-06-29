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
// are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create network with invalid region
//   - Create network with invalid SKU
//   - Create network with non-existent workspace
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
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mocknetwork.ConfigureNetworkErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *NetworkErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork))
	suite.ConfigureDepends(t, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

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

	// Teardown
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *NetworkErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
