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

// SecurityGroupErrorV1TestSuite verifies that SecurityGroup resources with
// invalid references are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create security group with invalid region
//   - Create security group with non-existent workspace
//   - Create security group with non-existent rule ref
type SecurityGroupErrorV1TestSuite struct {
	suites.RegionalTestSuite
	params *params.SecurityGroupErrorV1Params
}

func CreateSecurityGroupErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *SecurityGroupErrorV1TestSuite {
	suite := &SecurityGroupErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.SecurityGroupErrorV1SuiteName.String()
	return suite
}

func (suite *SecurityGroupErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.NetworkParentSuite)

	workspaceName := generators.GenerateWorkspaceName()
	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace for security group error scenarios testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	// Non-existent rule ref — points to a rule that was never created
	nonExistentRuleRef := generators.GenerateSecurityGroupRuleRefObject(
		sdkconsts.NetworkProviderV1Name,
		suite.Tenant,
		workspaceName,
		"non-existent-rule",
	)

	buildSecurityGroup := func(name string, workspaceRef string, region string, ruleRefs []schema.Reference) *schema.SecurityGroup {
		sg, err := builders.NewSecurityGroupBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceRef).Region(region).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "SecurityGroup for error scenario testing"}).
			Spec(&schema.SecurityGroupSpec{
				RuleRefs: ruleRefs,
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build SecurityGroup: %v", err)
		}
		return sg
	}

	p := &params.SecurityGroupErrorV1Params{
		Workspace: workspace,

		// Invalid region — valid ruleRefs, invalid region
		InvalidRegionSecurityGroup: buildSecurityGroup(
			generators.GenerateSecurityGroupName(),
			workspaceName,
			"invalid-region",
			[]schema.Reference{*nonExistentRuleRef},
		),

		// Non-existent workspace — workspace was never created
		NonExistentWorkspaceSecurityGroup: buildSecurityGroup(
			generators.GenerateSecurityGroupName(),
			"non-existent-workspace",
			suite.Region,
			[]schema.Reference{*nonExistentRuleRef},
		),

		// Non-existent rule ref — valid workspace + region, rule does not exist
		NonExistentRuleRefSecurityGroup: buildSecurityGroup(
			generators.GenerateSecurityGroupName(),
			workspaceName,
			suite.Region,
			[]schema.Reference{*nonExistentRuleRef},
		),
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mocknetwork.ConfigureSecurityGroupErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *SecurityGroupErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup))
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
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with invalid region — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidRegionSecurityGroup,
	)

	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with non-existent workspace — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentWorkspaceSecurityGroup,
	)

	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with non-existent rule ref — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentRuleRefSecurityGroup,
	)

	// Teardown
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *SecurityGroupErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
