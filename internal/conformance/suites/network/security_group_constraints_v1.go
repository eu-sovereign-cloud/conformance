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

// SecurityGroupConstraintsValidationV1TestSuite verifies that SecurityGroup resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
//   - spec.rules[].direction: enum [ingress, egress] (SecurityGroupRuleSpec)
//   - spec.rules[].version: enum [IPv4, IPv6] (SecurityGroupRuleSpec)
//   - spec.rules[].protocol: enum [tcp, udp, tcp+udp, icmp] (SecurityGroupRuleSpec)
//   - spec.rules[].ports.from: minimum 1, maximum 65535 (Port)
//   - spec.rules[].ports.to: minimum 1, maximum 65535 (Port)
//   - spec.rules[].ports.list[]: minimum 1, maximum 65535 (Port)
//   - spec.rules[].icmp.type: maximum 8 (IcmpConfig)
//   - spec.rules[].icmp.code: maximum 5 (IcmpConfig)
type SecurityGroupConstraintsValidationV1TestSuite struct {
	suites.RegionalTestSuite

	params *params.SecurityGroupConstraintsValidationV1Params
}

func CreateSecurityGroupConstraintsValidationV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *SecurityGroupConstraintsValidationV1TestSuite {
	suite := &SecurityGroupConstraintsValidationV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.SecurityGroupConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *SecurityGroupConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.NetworkParentSuite)

	workspaceName := generators.GenerateWorkspaceName()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Workspace for security group constraints violations testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	buildSG := func(name string, labels schema.Labels, annotations schema.Annotations) *schema.SecurityGroup {
		sg, err := builders.NewSecurityGroupBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Labels(labels).Annotations(annotations).
			Spec(&schema.SecurityGroupSpec{
				Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build SecurityGroup: %v", err)
		}
		return sg
	}

	buildSGWithRules := func(name string, rules []schema.SecurityGroupRuleSpec) *schema.SecurityGroup {
		sg, err := builders.NewSecurityGroupBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
			Annotations(schema.Annotations{"description": "SecurityGroup with invalid inline rule"}).
			Spec(&schema.SecurityGroupSpec{Rules: rules}).Build()
		if err != nil {
			t.Fatalf("Failed to build SecurityGroup: %v", err)
		}
		return sg
	}

	overMaxPort := 65536
	underMinPort := 0
	overMaxIcmpType := 9
	overMaxIcmpCode := 6

	p := &params.SecurityGroupConstraintsValidationV1Params{
		Workspace: workspace,
		OverLengthNameSecurityGroup: buildSG(
			strings.Repeat("a", 129),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "SecurityGroup with over-length name"},
		),
		InvalidPatternNameSecurityGroup: buildSG(
			"Invalid-Name-With-Uppercase",
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "SecurityGroup with non-kebab-case name"},
		),
		OverLengthLabelValueSecurityGroup: buildSG(
			generators.GenerateSecurityGroupName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel, "constraint-test": strings.Repeat("x", 64)},
			schema.Annotations{"description": "SecurityGroup with over-length label value"},
		),
		OverLengthAnnotationSecurityGroup: buildSG(
			generators.GenerateSecurityGroupName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "SecurityGroup with over-length annotation value", "long-annotation": strings.Repeat("y", 1025)},
		),
		InvalidDirectionSecurityGroup: buildSGWithRules(
			generators.GenerateSecurityGroupName(),
			[]schema.SecurityGroupRuleSpec{{Direction: "invalid-direction"}},
		),
		InvalidVersionSecurityGroup: buildSGWithRules(
			generators.GenerateSecurityGroupName(),
			[]schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress, Version: "IPv5"}},
		),
		InvalidProtocolSecurityGroup: buildSGWithRules(
			generators.GenerateSecurityGroupName(),
			[]schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress, Protocol: "ftp"}},
		),
		OverMaxPortFromSecurityGroup: buildSGWithRules(
			generators.GenerateSecurityGroupName(),
			[]schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{From: overMaxPort}}},
		),
		UnderMinPortFromSecurityGroup: buildSGWithRules(
			generators.GenerateSecurityGroupName(),
			[]schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{From: underMinPort}}},
		),
		OverMaxPortToSecurityGroup: buildSGWithRules(
			generators.GenerateSecurityGroupName(),
			[]schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{To: overMaxPort}}},
		),
		UnderMinPortToSecurityGroup: buildSGWithRules(
			generators.GenerateSecurityGroupName(),
			[]schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{To: underMinPort}}},
		),
		OverMaxPortListSecurityGroup: buildSGWithRules(
			generators.GenerateSecurityGroupName(),
			[]schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{List: []int{overMaxPort}}}},
		),
		UnderMinPortListSecurityGroup: buildSGWithRules(
			generators.GenerateSecurityGroupName(),
			[]schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{List: []int{underMinPort}}}},
		),
		OverMaxIcmpTypeSecurityGroup: buildSGWithRules(
			generators.GenerateSecurityGroupName(),
			[]schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress, Icmp: &schema.IcmpConfig{Type: overMaxIcmpType, Code: 0}}},
		),
		OverMaxIcmpCodeSecurityGroup: buildSGWithRules(
			generators.GenerateSecurityGroupName(),
			[]schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress, Icmp: &schema.IcmpConfig{Type: 0, Code: overMaxIcmpCode}}},
		),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureSecurityGroupConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *SecurityGroupConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalNetworkResourceMetadataKindResourceKindSecurityGroup))
	suite.ConfigureDepends(t, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace
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

	// Security Group metadata violations
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with name exceeding maxLength:128 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthNameSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidPatternNameSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with label value exceeding maxLength:63 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthLabelValueSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthAnnotationSecurityGroup,
	)

	// Security Group inline rule violations
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with rules[].direction outside enum [ingress, egress] — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidDirectionSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with rules[].version outside enum [IPv4, IPv6] — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidVersionSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with rules[].protocol outside enum [tcp, udp, tcp+udp, icmp] — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidProtocolSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with rules[].ports.from exceeding maximum:65535 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverMaxPortFromSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with rules[].ports.from below minimum:1 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.UnderMinPortFromSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with rules[].ports.to exceeding maximum:65535 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverMaxPortToSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with rules[].ports.to below minimum:1 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.UnderMinPortToSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with rules[].ports.list[] exceeding maximum:65535 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverMaxPortListSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with rules[].ports.list[] below minimum:1 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.UnderMinPortListSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with rules[].icmp.type exceeding maximum:8 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverMaxIcmpTypeSecurityGroup,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupExpectViolationV1Step(
		"Create a security group with rules[].icmp.code exceeding maximum:5 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverMaxIcmpCodeSecurityGroup,
	)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *SecurityGroupConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
