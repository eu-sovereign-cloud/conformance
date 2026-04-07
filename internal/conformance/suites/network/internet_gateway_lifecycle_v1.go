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

type InternetGatewayLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	params *params.InternetGatewayLifeCycleV1Params
}

func CreateInternetGatewayLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *InternetGatewayLifeCycleV1TestSuite {
	suite := &InternetGatewayLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.InternetGatewayLifeCycleV1SuiteName.String()
	return suite
}

func (suite *InternetGatewayLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Network")

	workspaceName := generators.GenerateWorkspaceName()
	internetGatewayName := generators.GenerateInternetGatewayName()

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
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	internetGatInitial, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Internet Gateway for conformance testing",
		}).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: false,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	internetGatUpdated, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Internet Gateway for conformance testing",
		}).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: true,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	params := &params.InternetGatewayLifeCycleV1Params{
		Workspace:              workspace,
		InternetGatewayInitial: internetGatInitial,
		InternetGatewayUpdated: internetGatUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureInternetGatewayLifecycleScenarioV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *InternetGatewayLifeCycleV1TestSuite) TestScenario(t provider.T) {
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
	expectWorkspaceAnnotations := workspace.Annotations
	expectWorkspaceExtensions := workspace.Extensions
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
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
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:      expectWorkspaceLabels,
			Annotations: expectWorkspaceAnnotations,
			Extensions:  expectWorkspaceExtensions,
			Metadata:    expectWorkspaceMeta,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Internet gateway

	// Create an internet gateway
	gateway := suite.params.InternetGatewayInitial
	expectGatewayMeta := gateway.Metadata
	expectGatewaySpec := &gateway.Spec
	expectGatewayLabels := gateway.Labels
	expectGatewayAnnotations := gateway.Annotations
	expectGatewayExtensions := gateway.Extensions
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create a internet gateway", suite.Client.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:         expectGatewayLabels,
			Annotations:    expectGatewayAnnotations,
			Extensions:     expectGatewayExtensions,
			Metadata:       expectGatewayMeta,
			Spec:           expectGatewaySpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created internet gateway
	gatewayWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(gateway.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(gateway.Metadata.Workspace),
		Name:      gateway.Metadata.Name,
	}
	stepsBuilder.GetInternetGatewayV1Step("Get the created internet gateway", suite.Client.NetworkV1, gatewayWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:      expectGatewayLabels,
			Annotations: expectGatewayAnnotations,
			Extensions:  expectGatewayExtensions,
			Metadata:    expectGatewayMeta,
			Spec:        expectGatewaySpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Update the internet gateway
	gateway = suite.params.InternetGatewayUpdated
	expectGatewaySpec.EgressOnly = gateway.Spec.EgressOnly
	expectGatewayLabels = gateway.Labels
	expectGatewayAnnotations = gateway.Annotations
	expectGatewayExtensions = gateway.Extensions
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Update the internet gateway", suite.Client.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:         expectGatewayLabels,
			Annotations:    expectGatewayAnnotations,
			Extensions:     expectGatewayExtensions,
			Metadata:       expectGatewayMeta,
			Spec:           expectGatewaySpec,
			ResourceStates: suites.UpdatedResourceExpectedStates,
		},
	)

	// Get the updated internet gateway
	stepsBuilder.GetInternetGatewayV1Step("Get the updated internet gateway", suite.Client.NetworkV1, gatewayWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:      expectGatewayLabels,
			Annotations: expectGatewayAnnotations,
			Extensions:  expectGatewayExtensions,
			Metadata:    expectGatewayMeta,
			Spec:        expectGatewaySpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)

	// Resources deletion
	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, gateway)
	stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", suite.Client.NetworkV1, gatewayWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *InternetGatewayLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
