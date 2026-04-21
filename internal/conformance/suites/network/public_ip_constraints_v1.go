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

// PublicIpConstraintsValidationV1TestSuite verifies that PublicIp resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type PublicIpConstraintsValidationV1TestSuite struct {
	suites.RegionalTestSuite

	config *PublicIpLifeCycleV1Config
	params *params.PublicIpConstraintsValidationV1Params
}

func CreatePublicIpConstraintsValidationV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *PublicIpLifeCycleV1Config) *PublicIpConstraintsValidationV1TestSuite {
	suite := &PublicIpConstraintsValidationV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.PublicIpConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *PublicIpConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	workspaceName := generators.GenerateWorkspaceName()

	publicIpAddress, err := generators.GeneratePublicIp(suite.config.PublicIpsRange, 1)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Workspace for public ip constraints violations testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	overLengthNamePublicIp, err := builders.NewPublicIpBuilder().
		Name(strings.Repeat("a", 129)).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "PublicIp with over-length name"}).
		Spec(&schema.PublicIpSpec{Version: schema.IPVersionIPv4, Address: publicIpAddress}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthNamePublicIp: %v", err)
	}

	invalidPatternNamePublicIp, err := builders.NewPublicIpBuilder().
		Name("Invalid-Name-With-Uppercase").
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "PublicIp with non-kebab-case name"}).
		Spec(&schema.PublicIpSpec{Version: schema.IPVersionIPv4, Address: publicIpAddress}).Build()
	if err != nil {
		t.Fatalf("Failed to build invalidPatternNamePublicIp: %v", err)
	}

	overLengthLabelPublicIp, err := builders.NewPublicIpBuilder().
		Name(generators.GeneratePublicIpName()).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
			"constraint-test":  strings.Repeat("x", 64),
		}).
		Annotations(schema.Annotations{"description": "PublicIp with over-length label value"}).
		Spec(&schema.PublicIpSpec{Version: schema.IPVersionIPv4, Address: publicIpAddress}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthLabelPublicIp: %v", err)
	}

	overLengthAnnotationPublicIp, err := builders.NewPublicIpBuilder().
		Name(generators.GeneratePublicIpName()).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{
			"description":     "PublicIp with over-length annotation value",
			"long-annotation": strings.Repeat("y", 1025),
		}).
		Spec(&schema.PublicIpSpec{Version: schema.IPVersionIPv4, Address: publicIpAddress}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthAnnotationPublicIp: %v", err)
	}

	p := &params.PublicIpConstraintsValidationV1Params{
		Workspace:                    workspace,
		OverLengthNamePublicIp:       overLengthNamePublicIp,
		InvalidPatternNamePublicIp:   invalidPatternNamePublicIp,
		OverLengthLabelValuePublicIp: overLengthLabelPublicIp,
		OverLengthAnnotationPublicIp: overLengthAnnotationPublicIp,
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigurePublicIpConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *PublicIpConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.NetworkProviderV1Name, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP))

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

	stepsBuilder.CreateOrUpdatePublicIpExpectViolationV1Step(
		"Create a public ip with name exceeding maxLength:128 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthNamePublicIp,
	)
	stepsBuilder.CreateOrUpdatePublicIpExpectViolationV1Step(
		"Create a public ip with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidPatternNamePublicIp,
	)
	stepsBuilder.CreateOrUpdatePublicIpExpectViolationV1Step(
		"Create a public ip with label value exceeding maxLength:63 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthLabelValuePublicIp,
	)
	stepsBuilder.CreateOrUpdatePublicIpExpectViolationV1Step(
		"Create a public ip with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthAnnotationPublicIp,
	)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *PublicIpConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
