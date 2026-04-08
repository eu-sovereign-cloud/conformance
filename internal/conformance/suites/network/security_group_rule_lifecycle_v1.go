package network

import (
	"log/slog"

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

type SecurityGroupRuleLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	params *params.SecurityGroupRuleLifeCycleV1Params
}

func CreateSecurityGroupRuleLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *SecurityGroupRuleLifeCycleV1TestSuite {
	suite := &SecurityGroupRuleLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.SecurityGroupRuleLifeCycleV1SuiteName.String()
	return suite
}

func (suite *SecurityGroupRuleLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.NetworkParentSuite)

	workspaceName := generators.GenerateWorkspaceName()
	securityGroupRuleName := generators.GenerateSecurityGroupRuleName()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Workspace for conformance testing",
		}).
		Build()
	if err != nil {
		slog.Error("Failed to build Workspace", "error", err)
		t.FailNow()
	}

	securityGroupRuleInitial, err := builders.NewSecurityGroupRuleBuilder().
		Name(securityGroupRuleName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Security Group Rule for conformance testing",
		}).
		Spec(&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress}).
		Build()
	if err != nil {
		slog.Error("Failed to build Security Group Rule", "error", err)
		t.FailNow()
	}

	securityGroupRuleUpdated, err := builders.NewSecurityGroupRuleBuilder().
		Name(securityGroupRuleName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Security Group Rule for conformance testing",
		}).
		Spec(&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionEgress}).
		Build()
	if err != nil {
		slog.Error("Failed to build Security Group Rule", "error", err)
		t.FailNow()
	}

	params := &params.SecurityGroupRuleLifeCycleV1Params{
		Workspace:                workspace,
		SecurityGroupRuleInitial: securityGroupRuleInitial,
		SecurityGroupRuleUpdated: securityGroupRuleUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureSecurityGroupRuleLifecycleScenarioV1, *params)
	if err != nil {
		slog.Error("Failed to setup mock", "error", err)
		t.FailNow()
	}
}

func (suite *SecurityGroupRuleLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalNetworkResourceMetadataKindResourceKindSecurityGroupRule))
	suite.ConfigureDepends(t, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

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

	// Security Group Rule

	// Create a security group rule
	rule := suite.params.SecurityGroupRuleInitial
	expectRuleMeta := rule.Metadata
	expectRuleSpec := &rule.Spec
	expectRuleLabels := rule.Labels
	expectRuleAnnotations := rule.Annotations
	expectRuleExtensions := rule.Extensions
	stepsBuilder.CreateOrUpdateSecurityGroupRuleV1Step("Create a security group rule", t, suite.Client.NetworkV1, rule,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Labels:         expectRuleLabels,
			Annotations:    expectRuleAnnotations,
			Extensions:     expectRuleExtensions,
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created security group rule
	ruleWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(rule.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(rule.Metadata.Workspace),
		Name:      rule.Metadata.Name,
	}
	stepsBuilder.GetSecurityGroupRuleV1Step("Get the created security group rule", suite.Client.NetworkV1, ruleWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec, schema.SecurityGroupRuleStatus]{
			Metadata: expectRuleMeta,
			Spec:     expectRuleSpec,
			ResourceStatus: schema.SecurityGroupRuleStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Update the security group rule
	rule.Spec.Direction = schema.SecurityGroupRuleDirectionEgress
	expectRuleSpec.Direction = rule.Spec.Direction
	expectRuleLabels = rule.Labels
	expectRuleAnnotations = rule.Annotations
	expectRuleExtensions = rule.Extensions
	stepsBuilder.CreateOrUpdateSecurityGroupRuleV1Step("Update the security group rule", t, suite.Client.NetworkV1, rule,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Labels:         expectRuleLabels,
			Annotations:    expectRuleAnnotations,
			Extensions:     expectRuleExtensions,
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: suites.UpdatedResourceExpectedStates,
		},
	)

	// Get the updated security group rule
	stepsBuilder.GetSecurityGroupRuleV1Step("Get the updated security group rule", suite.Client.NetworkV1, ruleWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec, schema.SecurityGroupRuleStatus]{
			Metadata: expectRuleMeta,
			Spec:     expectRuleSpec,
			ResourceStatus: schema.SecurityGroupRuleStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)

	// Resources deletion
	stepsBuilder.DeleteSecurityGroupRuleV1Step("Delete the security group rule", t, suite.Client.NetworkV1, rule)
	stepsBuilder.WatchSecurityGroupRuleUntilDeletedV1Step("Watch the security group rule deletion", t, suite.Client.NetworkV1, ruleWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *SecurityGroupRuleLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
