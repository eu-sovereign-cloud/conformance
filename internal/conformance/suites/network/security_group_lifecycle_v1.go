package network

import (
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

type SecurityGroupLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	params *params.SecurityGroupLifeCycleV1Params
}

func CreateSecurityGroupLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *SecurityGroupLifeCycleV1TestSuite {
	suite := &SecurityGroupLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.SecurityGroupLifeCycleV1SuiteName.String()
	return suite
}

func (suite *SecurityGroupLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.NetworkParentSuite)

	workspaceName := generators.GenerateWorkspaceName()
	securityGroupName := generators.GenerateSecurityGroupName()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	securityGroupInitial, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	securityGroupUpdated, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionEgress}},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	params := &params.SecurityGroupLifeCycleV1Params{
		Workspace:            workspace,
		SecurityGroupInitial: securityGroupInitial,
		SecurityGroupUpdated: securityGroupUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureSecurityGroupLifecycleScenarioV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *SecurityGroupLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalNetworkResourceMetadataKindResourceKindSecurityGroup))
	suite.ConfigureDepends(t, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsConfigurator := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsConfigurator.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
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
	stepsConfigurator.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Security Group

	// Create a security group
	group := suite.params.SecurityGroupInitial
	expectGroupMeta := group.Metadata
	expectGroupSpec := &group.Spec
	stepsConfigurator.CreateOrUpdateSecurityGroupV1Step("Create a security group", t, suite.Client.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created security group
	groupWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(group.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(group.Metadata.Workspace),
		Name:      group.Metadata.Name,
	}
	stepsConfigurator.GetSecurityGroupV1Step("Get the created security group", suite.Client.NetworkV1, groupWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the security group
	groupRules := group.Spec.Rules
	groupRules[0] = schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionEgress}
	expectGroupSpec.Rules = group.Spec.Rules
	stepsConfigurator.CreateOrUpdateSecurityGroupV1Step("Update the security group", t, suite.Client.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: suites.UpdatedResourceExpectedStates,
		},
	)

	// Get the updated security group
	stepsConfigurator.GetSecurityGroupV1Step("Get the updated security group", suite.Client.NetworkV1, groupWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Resources deletion
	stepsConfigurator.DeleteSecurityGroupV1Step("Delete the security group", t, suite.Client.NetworkV1, group)
	stepsConfigurator.WatchSecurityGroupUntilDeletedV1Step("Watch the security group deletion", t, suite.Client.NetworkV1, groupWRef)

	stepsConfigurator.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsConfigurator.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *SecurityGroupLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
