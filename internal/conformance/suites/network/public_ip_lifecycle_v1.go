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

type PublicIpLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	config *PublicIpLifeCycleV1Config
	params *params.PublicIpLifeCycleV1Params
}

type PublicIpLifeCycleV1Config struct {
	PublicIpsRange string
	NetworkSkus    []string
}

func CreatePublicIpLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *PublicIpLifeCycleV1Config) *PublicIpLifeCycleV1TestSuite {
	suite := &PublicIpLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.PublicIpLifeCycleV1SuiteName.String()
	return suite
}

func (suite *PublicIpLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Network")

	workspaceName := generators.GenerateWorkspaceName()
	publicIpName := generators.GeneratePublicIpName()

	// Generate the public ips
	publicIpAddress1, err := generators.GeneratePublicIp(suite.config.PublicIpsRange, 1)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}
	publicIpAddress2, err := generators.GeneratePublicIp(suite.config.PublicIpsRange, 2)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}

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

	publicIpInitial, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Public IP for conformance testing",
		}).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: publicIpAddress1,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	publicIpUpdated, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Public IP for conformance testing",
		}).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: publicIpAddress2,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	params := &params.PublicIpLifeCycleV1Params{
		Workspace:       workspace,
		PublicIpInitial: publicIpInitial,
		PublicIpUpdated: publicIpUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigurePublicIpLifecycleScenarioV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *PublicIpLifeCycleV1TestSuite) TestScenario(t provider.T) {
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
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			Labels:   expectWorkspaceLabels,
			Metadata: expectWorkspaceMeta,
			ResourceStatus: schema.WorkspaceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Public ip

	// Create a public ip
	publicIp := suite.params.PublicIpInitial
	expectPublicIpMeta := publicIp.Metadata
	expectPublicIpSpec := &publicIp.Spec
	expectPublicIpLabels := publicIp.Labels
	expectPublicIpAnnotations := publicIp.Annotations
	expectPublicIpExtensions := publicIp.Extensions
	stepsBuilder.CreateOrUpdatePublicIpV1Step("Create a public ip", suite.Client.NetworkV1, publicIp,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Labels:         expectPublicIpLabels,
			Annotations:    expectPublicIpAnnotations,
			Extensions:     expectPublicIpExtensions,
			Metadata:       expectPublicIpMeta,
			Spec:           expectPublicIpSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created public ip
	publicIpWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(publicIp.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(publicIp.Metadata.Workspace),
		Name:      publicIp.Metadata.Name,
	}
	stepsBuilder.GetPublicIpV1Step("Get the created public ip", suite.Client.NetworkV1, publicIpWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			Metadata: expectPublicIpMeta,
			Spec:     expectPublicIpSpec,
			ResourceStatus: schema.PublicIpStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Update the public ip
	publicIp = suite.params.PublicIpUpdated
	expectPublicIpSpec.Address = publicIp.Spec.Address
	expectPublicIpLabels = publicIp.Labels
	expectPublicIpAnnotations = publicIp.Annotations
	expectPublicIpExtensions = publicIp.Extensions
	stepsBuilder.CreateOrUpdatePublicIpV1Step("Update the public ip", suite.Client.NetworkV1, publicIp,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Labels:         expectPublicIpLabels,
			Annotations:    expectPublicIpAnnotations,
			Extensions:     expectPublicIpExtensions,
			Metadata:       expectPublicIpMeta,
			Spec:           expectPublicIpSpec,
			ResourceStates: suites.UpdatedResourceExpectedStates,
		},
	)

	// Get the updated public ip
	stepsBuilder.GetPublicIpV1Step("Get the updated public ip", suite.Client.NetworkV1, publicIpWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			Metadata: expectPublicIpMeta,
			Spec:     expectPublicIpSpec,
			ResourceStatus: schema.PublicIpStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)

	// Resources deletion
	stepsBuilder.DeletePublicIpV1Step("Delete the public ip", suite.Client.NetworkV1, publicIp)
	stepsBuilder.WatchPublicIpUntilDeletedV1Step("Watch the public ip deletion", suite.Client.NetworkV1, publicIpWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *PublicIpLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
