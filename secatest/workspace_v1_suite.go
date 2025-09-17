package secatest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type WorkspaceV1TestSuite struct {
	regionalTestSuite
}

func (suite *WorkspaceV1TestSuite) TestWorkspaceV1(t provider.T) {
	slog.Info("Starting Workspace Lifecycle Test")

	t.Title("Workspace Lifecycle Test")
	configureTags(t, secalib.WorkspaceProviderV1, secalib.WorkspaceKind)

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, workspaceName)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		wm, err := mock.CreateWorkspaceLifecycleScenarioV1("Workspace Lifecycle",
			mock.WorkspaceParamsV1{
				Params: &mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
					Region:    suite.region,
				},
				Workspace: &mock.ResourceParams[secalib.WorkspaceSpecV1]{
					Name: workspaceName,
					InitialSpec: &secalib.WorkspaceSpecV1{
						Labels: &[]secalib.Label{
							{
								Name:  secalib.EnvLabel,
								Value: secalib.EnvDevelopmentLabel,
							},
						},
					},
					UpdatedSpec: &secalib.WorkspaceSpecV1{
						Labels: &[]secalib.Label{
							{
								Name:  secalib.EnvLabel,
								Value: secalib.EnvProductionLabel,
							},
						},
					},
				},
			})
		if err != nil {
			t.Fatalf("Failed to create workspace scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var resp *workspace.Workspace
	var err error

	var expectedMeta *secalib.Metadata
	var expectedLabel *[]secalib.Label
	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		ws := &workspace.Workspace{
			Metadata: &workspace.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta = &secalib.Metadata{
			Name:       workspaceName,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		verifyWorkspaceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
	})

	t.WithNewStep("Get created workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		resp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		verifyWorkspaceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		expectedLabel = &[]secalib.Label{
			{
				Name:  secalib.EnvLabel,
				Value: secalib.EnvDevelopmentLabel,
			},
		}
		var actualLabels *[]secalib.Label
		if resp.Labels != nil {
			labels := make([]secalib.Label, 0, len(*resp.Labels))
			for k, v := range *resp.Labels {
				labels = append(labels, secalib.Label{Name: k, Value: v})
			}
			actualLabels = &labels
		}
		verifyLabelStep(sCtx, expectedLabel, actualLabels)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
	})

	t.WithNewStep("Update workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, resp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodPut
		verifyWorkspaceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
	})

	t.WithNewStep("Get updated workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		resp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		verifyWorkspaceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		expectedLabel = &[]secalib.Label{
			{
				Name:  secalib.EnvLabel,
				Value: secalib.EnvProductionLabel,
			},
		}
		var actualLabels *[]secalib.Label
		if resp.Labels != nil {
			labels := make([]secalib.Label, 0, len(*resp.Labels))
			for k, v := range *resp.Labels {
				labels = append(labels, secalib.Label{Name: k, Value: v})
			}
			actualLabels = &labels
		}
		verifyLabelStep(sCtx, expectedLabel, actualLabels)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
	})

	t.WithNewStep("Delete workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "DeleteWorkspace")

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		_, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	slog.Info("Finishing Workspace Lifecycle Test")
}

func (suite *WorkspaceV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}

func verifyWorkspaceMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *workspace.RegionalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Region:     &metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}
