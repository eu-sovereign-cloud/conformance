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

// SecurityGroupRuleErrorV1TestSuite verifies that SecurityGroupRule resources with
// invalid references are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create security group rule with invalid region
//   - Create security group rule with non-existent workspace
type SecurityGroupRuleErrorV1TestSuite struct {
	suites.RegionalTestSuite
	params *params.SecurityGroupRuleErrorV1Params
}

func CreateSecurityGroupRuleErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *SecurityGroupRuleErrorV1TestSuite {
	suite := &SecurityGroupRuleErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.SecurityGroupRuleErrorV1SuiteName.String()
	return suite
}

func (suite *SecurityGroupRuleErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.NetworkParentSuite)

	workspaceName := generators.GenerateWorkspaceName()
	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace for security group rule error scenarios testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	buildSecurityGroupRule := func(name string, workspaceRef string, region string) *schema.SecurityGroupRule {
		sgr, err := builders.NewSecurityGroupRuleBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceRef).Region(region).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "SecurityGroupRule for error scenario testing"}).
			Spec(&schema.SecurityGroupRuleSpec{
				Direction: schema.SecurityGroupRuleDirectionIngress,
			}).
			Build()
		if err != nil {
			t.Fatalf("Failed to build SecurityGroupRule: %v", err)
		}
		return sgr
	}

	p := &params.SecurityGroupRuleErrorV1Params{
		Workspace: workspace,

		// Invalid region — random string, valid workspace
		InvalidRegionSecurityGroupRule: buildSecurityGroupRule(
			generators.GenerateSecurityGroupRuleName(),
			workspaceName,
			"invalid-region",
		),

		// Non-existent workspace — workspace was never created
		NonExistentWorkspaceSecurityGroupRule: buildSecurityGroupRule(
			generators.GenerateSecurityGroupRuleName(),
			"non-existent-workspace",
			suite.Region,
		),
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mocknetwork.ConfigureSecurityGroupRuleErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *SecurityGroupRuleErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroupRule))
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
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with invalid region — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidRegionSecurityGroupRule,
	)

	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with non-existent workspace — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentWorkspaceSecurityGroupRule,
	)

	// Teardown
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *SecurityGroupRuleErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
