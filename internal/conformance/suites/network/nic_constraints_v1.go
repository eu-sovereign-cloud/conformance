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

// NicConstraintsValidationV1TestSuite verifies that Nic resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type NicConstraintsValidationV1TestSuite struct {
	suites.RegionalTestSuite
	config *NicLifeCycleV1Config

	params *params.NicConstraintsValidationV1Params
}

func CreateNicConstraintsValidationV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *NicLifeCycleV1Config) *NicConstraintsValidationV1TestSuite {
	suite := &NicConstraintsValidationV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.NicConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *NicConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	workspaceName := generators.GenerateWorkspaceName()
	networkName := generators.GenerateNetworkName()
	subnetName := generators.GenerateSubnetName()

	subnetCidr, err := generators.GenerateSubnetCidr(suite.config.NetworkCidr, 8, 1)
	if err != nil {
		t.Fatalf("Failed to generate subnet cidr: %v", err)
	}

	subnetRefObj := generators.GenerateSubnetRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, subnetName)

	// Generate the nic addresses
	nicAddress, err := generators.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Workspace for nic constraints violations testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	buildNic := func(name string, labels schema.Labels, annotations schema.Annotations) *schema.Nic {
		nic, err := builders.NewNicBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Labels(labels).Annotations(annotations).
			Spec(&schema.NicSpec{
				Addresses: []string{nicAddress},
				SubnetRef: *subnetRefObj,
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build Nic: %v", err)
		}
		return nic
	}

	p := &params.NicConstraintsValidationV1Params{
		Workspace: workspace,
		OverLengthNameNic: buildNic(
			strings.Repeat("a", 129),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Nic with over-length name"},
		),
		InvalidPatternNameNic: buildNic(
			"Invalid-Name-With-Uppercase",
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Nic with non-kebab-case name"},
		),
		OverLengthLabelValueNic: buildNic(
			generators.GenerateNicName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel, "constraint-test": strings.Repeat("x", 64)},
			schema.Annotations{"description": "Nic with over-length label value"},
		),
		OverLengthAnnotationNic: buildNic(
			generators.GenerateNicName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Nic with over-length annotation value", "long-annotation": strings.Repeat("y", 1025)},
		),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureNicConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *NicConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic))
	suite.ConfigureDepends(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
	)
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
	// Nic constraints violations
	stepsBuilder.CreateOrUpdateNicExpectViolationV1Step(
		"Create a nic with name exceeding maxLength:128 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthNameNic,
	)
	stepsBuilder.CreateOrUpdateNicExpectViolationV1Step(
		"Create a nic with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidPatternNameNic,
	)
	stepsBuilder.CreateOrUpdateNicExpectViolationV1Step(
		"Create a nic with label value exceeding maxLength:63 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthLabelValueNic,
	)
	stepsBuilder.CreateOrUpdateNicExpectViolationV1Step(
		"Create a nic with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.NetworkV1,
		suite.params.OverLengthAnnotationNic,
	)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *NicConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
