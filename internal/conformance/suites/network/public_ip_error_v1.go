package network

import (
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

// PublicIpErrorV1TestSuite verifies that PublicIp resources with
// invalid references are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create public ip with invalid region
//   - Create public ip with non-existent workspace
//   - Create public ip with invalid IP version
type PublicIpErrorV1TestSuite struct {
	suites.RegionalTestSuite

	config *PublicIpLifeCycleV1Config
	params *params.PublicIpErrorV1Params
}

func CreatePublicIpErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *PublicIpLifeCycleV1Config) *PublicIpErrorV1TestSuite {
	suite := &PublicIpErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}

	suite.ScenarioName = constants.PublicIpErrorV1SuiteName.String()
	return suite
}

func (suite *PublicIpErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.NetworkParentSuite)

	workspaceName := generators.GenerateWorkspaceName()
	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace for public ip error scenarios testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	buildPublicIp := func(name string, workspaceRef string, region string, ipVersion schema.IPVersion) *schema.PublicIp {
		pip, err := builders.NewPublicIpBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceRef).Region(region).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "PublicIp for error scenario testing"}).
			Spec(&schema.PublicIpSpec{
				Version: ipVersion,
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build PublicIp: %v", err)
		}
		return pip
	}

	p := &params.PublicIpErrorV1Params{
		Workspace: workspace,

		// Invalid region — random string, valid workspace
		InvalidRegionPublicIp: buildPublicIp(
			generators.GeneratePublicIpName(),
			workspaceName,
			"invalid-region",
			schema.IPVersionIPv4,
		),

		// Non-existent workspace — workspace was never created
		NonExistentWorkspacePublicIp: buildPublicIp(
			generators.GeneratePublicIpName(),
			"non-existent-workspace",
			suite.Region,
			schema.IPVersionIPv4,
		),

		// Invalid IP version — "IPv5" not in enum [IPv4, IPv6]
		InvalidVersionPublicIp: buildPublicIp(
			generators.GeneratePublicIpName(),
			workspaceName,
			suite.Region,
			"IPv5",
		),
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mocknetwork.ConfigurePublicIpErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *PublicIpErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP))
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
	stepsBuilder.CreateOrUpdatePublicIpExpectViolationV1Step(
		"Create a public ip with invalid region — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidRegionPublicIp,
	)

	stepsBuilder.CreateOrUpdatePublicIpExpectViolationV1Step(
		"Create a public ip with non-existent workspace — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentWorkspacePublicIp,
	)

	stepsBuilder.CreateOrUpdatePublicIpExpectViolationV1Step(
		"Create a public ip with invalid IP version (IPv5 not in [IPv4, IPv6]) — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidVersionPublicIp,
	)

	// Teardown
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *PublicIpErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
