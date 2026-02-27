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
	t.AddParentSuite("SecurityGroupRule")

	workspaceName := generators.GenerateWorkspaceName()
	securityGroupRuleName := generators.GenerateSecurityGroupRuleName()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
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
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureSecurityGroupRuleLifecycleScenarioV1, params)
	if err != nil {
		slog.Error("Failed to setup mock", "error", err)
		t.FailNow()
	}
}

func (suite *SecurityGroupRuleLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.NetworkProviderV1Name,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
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
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Security Group Rule

	// Create a security group rule
	rule := suite.params.SecurityGroupRuleInitial
	expectRuleMeta := rule.Metadata
	expectRuleSpec := &rule.Spec
	stepsBuilder.CreateOrUpdateSecurityGroupRuleV1Step("Create a security group rule", suite.Client.NetworkV1, rule,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created security group
	ruleWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(rule.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(rule.Metadata.Workspace),
		Name:      rule.Metadata.Name,
	}
	stepsBuilder.GetSecurityGroupRuleV1Step("Get the created security group rule", suite.Client.NetworkV1, ruleWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the security group
	rule.Spec.Direction = schema.SecurityGroupRuleDirectionEgress
	expectRuleSpec.Direction = rule.Spec.Direction
	stepsBuilder.CreateOrUpdateSecurityGroupRuleV1Step("Update the security group rule", suite.Client.NetworkV1, rule,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: suites.UpdatedResourceExpectedStates,
		},
	)

	// Get the updated security group
	stepsBuilder.GetSecurityGroupRuleV1Step("Get the updated security group rule", suite.Client.NetworkV1, ruleWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Resources deletion
	stepsBuilder.DeleteSecurityGroupRuleV1Step("Delete the security group rule", suite.Client.NetworkV1, rule)
	stepsBuilder.GetSecurityGroupRuleWithErrorV1Step("Get deleted security group rule", suite.Client.NetworkV1, ruleWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)
}

func (suite *SecurityGroupRuleLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
