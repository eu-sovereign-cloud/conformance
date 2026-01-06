package workspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/workspace"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type WorkspaceV1LifeCycleTestSuite struct {
	suites.RegionalTestSuite
}

func (suite *WorkspaceV1LifeCycleTestSuite) TestLifeCycleScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.WorkspaceProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	// Setup mock, if configured to use
	if suite.MockEnabled {
		mockParams := &mock.WorkspaceLifeCycleParamsV1{
			BaseParams: &mock.BaseParams{
				MockURL:   *suite.MockServerURL,
				AuthToken: suite.AuthToken,
				Tenant:    suite.Tenant,
				Region:    suite.Region,
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					constants.EnvLabel: constants.EnvDevelopmentLabel,
				},
				UpdatedLabels: schema.Labels{
					constants.EnvLabel: constants.EnvProductionLabel,
				},
			},
		}
		wm, err := workspace.ConfigureLifecycleScenarioV1(suite.ScenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.MockClient = wm
	}

	stepsBuilder := steps.NewBuilder(&suite.TestSuite, t)

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.Tenant,
			Name:   workspaceName,
		},
	}
	expectMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectLabels := schema.Labels{constants.EnvLabel: constants.EnvDevelopmentLabel}
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectLabels,
			Metadata:      expectMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	tref := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   workspaceName,
	}
	workspace = stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, *tref,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectLabels,
			Metadata:      expectMeta,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the workspace labels
	workspace.Labels = schema.Labels{
		constants.EnvLabel: constants.EnvProductionLabel,
	}
	expectLabels = workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Update the workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectLabels,
			Metadata:      expectMeta,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated workspace
	workspace = stepsBuilder.GetWorkspaceV1Step("Get the updated workspace", suite.Client.WorkspaceV1, *tref,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectLabels,
			Metadata:      expectMeta,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, *tref, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}
