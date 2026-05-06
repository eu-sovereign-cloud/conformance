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
//   - spec.direction: enum [ingress, egress] (SecurityGroupRuleSpec)
//   - spec.version: enum [IPv4, IPv6] (SecurityGroupRuleSpec)
//   - spec.protocol: enum [tcp, udp, tcp+udp, icmp] (SecurityGroupRuleSpec)
//   - spec.ports.from: minimum 1, maximum 65535 (Port)
//   - spec.ports.to: minimum 1, maximum 65535 (Port)
//   - spec.ports.list[]: minimum 1, maximum 65535 (Port)
//   - spec.icmp.type: maximum 8 (IcmpConfig)
//   - spec.icmp.code: maximum 5 (IcmpConfig)
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

	buildSGRWithSpec := func(name string, spec *schema.SecurityGroupRuleSpec) *schema.SecurityGroupRule {
		sgr, err := builders.NewSecurityGroupRuleBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
			Annotations(schema.Annotations{"description": "SecurityGroupRule with invalid spec"}).
			Spec(spec).
			Build()
		if err != nil {
			t.Fatalf("Failed to build SecurityGroupRule: %v", err)
		}
		return sgr
	}

	overMaxPort := 65536
	underMinPort := 0
	overMaxIcmpType := 9
	overMaxIcmpCode := 6

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
		InvalidDirectionSecurityGroupRule: buildSGRWithSpec(
			generators.GenerateSecurityGroupRuleName(),
			&schema.SecurityGroupRuleSpec{Direction: "invalid-direction"},
		),
		InvalidVersionSecurityGroupRule: buildSGRWithSpec(
			generators.GenerateSecurityGroupRuleName(),
			&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress, Version: "IPv5"},
		),
		InvalidProtocolSecurityGroupRule: buildSGRWithSpec(
			generators.GenerateSecurityGroupRuleName(),
			&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress, Protocol: "ftp"},
		),
		OverMaxPortFromSecurityGroupRule: buildSGRWithSpec(
			generators.GenerateSecurityGroupRuleName(),
			&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{From: overMaxPort}},
		),
		UnderMinPortFromSecurityGroupRule: buildSGRWithSpec(
			generators.GenerateSecurityGroupRuleName(),
			&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{From: underMinPort}},
		),
		OverMaxPortToSecurityGroupRule: buildSGRWithSpec(
			generators.GenerateSecurityGroupRuleName(),
			&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{To: overMaxPort}},
		),
		UnderMinPortToSecurityGroupRule: buildSGRWithSpec(
			generators.GenerateSecurityGroupRuleName(),
			&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{To: underMinPort}},
		),
		OverMaxPortListSecurityGroupRule: buildSGRWithSpec(
			generators.GenerateSecurityGroupRuleName(),
			&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{List: []int{overMaxPort}}},
		),
		UnderMinPortListSecurityGroupRule: buildSGRWithSpec(
			generators.GenerateSecurityGroupRuleName(),
			&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress, Ports: &schema.Ports{List: []int{underMinPort}}},
		),
		OverMaxIcmpTypeSecurityGroupRule: buildSGRWithSpec(
			generators.GenerateSecurityGroupRuleName(),
			&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress, Icmp: &schema.IcmpConfig{Type: overMaxIcmpType, Code: 0}},
		),
		OverMaxIcmpCodeSecurityGroupRule: buildSGRWithSpec(
			generators.GenerateSecurityGroupRuleName(),
			&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress, Icmp: &schema.IcmpConfig{Type: 0, Code: overMaxIcmpCode}},
		),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureSecurityGroupRuleConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *SecurityGroupRuleConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalNetworkResourceMetadataKindResourceKindSecurityGroupRule))
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

	// Security Group Rule metadata violations
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

	// Security Group Rule spec violations
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with direction outside enum [ingress, egress] — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidDirectionSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with version outside enum [IPv4, IPv6] — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidVersionSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with protocol outside enum [tcp, udp, tcp+udp, icmp] — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidProtocolSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with ports.from exceeding maximum:65535 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverMaxPortFromSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with ports.from below minimum:1 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.UnderMinPortFromSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with ports.to exceeding maximum:65535 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverMaxPortToSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with ports.to below minimum:1 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.UnderMinPortToSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with ports.list[] exceeding maximum:65535 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverMaxPortListSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with ports.list[] below minimum:1 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.UnderMinPortListSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with icmp.type exceeding maximum:8 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverMaxIcmpTypeSecurityGroupRule,
	)
	stepsBuilder.CreateOrUpdateSecurityGroupRuleExpectViolationV1Step(
		"Create a security group rule with icmp.code exceeding maximum:5 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverMaxIcmpCodeSecurityGroupRule,
	)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *SecurityGroupRuleConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
