package network

import (
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

// SecurityGroupRuleConstraintsValidationV1TestSuite verifies that SecurityGroupRule resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type SecurityGroupRuleConstraintsValidationV1TestSuite struct {
	suites.RegionalTestSuite

	params *params.SecurityGroupRuleConstraintsValidationV1Params
}

func CreateSecurityGroupRuleConstraintsValidationV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *SecurityGroupRuleConstraintsValidationV1TestSuite {
	suite := &SecurityGroupRuleConstraintsValidationV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.SecurityGroupRuleConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *SecurityGroupRuleConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	workspaceName := generators.GenerateWorkspaceName()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Workspace for security group rule constraints violations testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	buildSGR := func(name string, labels schema.Labels, annotations schema.Annotations) *schema.SecurityGroupRule {
		sgr, err := builders.NewSecurityGroupRuleBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Labels(labels).Annotations(annotations).
			Spec(&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress}).
			Build()
		if err != nil {
			t.Fatalf("Failed to build SecurityGroupRule: %v", err)
		}
		return sgr
	}

	p := &params.SecurityGroupRuleConstraintsValidationV1Params{
		Workspace: workspace,
		OverLengthNameSecurityGroupRule: buildSGR(strings.Repeat("a", 129),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "SecurityGroupRule with over-length name"}),
		InvalidPatternNameSecurityGroupRule: buildSGR("Invalid-Name-With-Uppercase",
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "SecurityGroupRule with non-kebab-case name"}),
		OverLengthLabelValueSecurityGroupRule: buildSGR(generators.GenerateSecurityGroupRuleName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel, "constraint-test": strings.Repeat("x", 64)},
			schema.Annotations{"description": "SecurityGroupRule with over-length label value"}),
		OverLengthAnnotationSecurityGroupRule: buildSGR(generators.GenerateSecurityGroupRuleName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "SecurityGroupRule with over-length annotation value", "long-annotation": strings.Repeat("y", 1025)}),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureSecurityGroupRuleConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *SecurityGroupRuleConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.NetworkProviderV1Name, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroupRule))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	workspace := suite.params.Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   workspace.Metadata.Name,
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

	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with name exceeding maxLength:128 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthNameSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidPatternNameSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with label value exceeding maxLength:63 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthLabelValueSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthAnnotationSecurityGroupRule,
	)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *SecurityGroupRuleConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
