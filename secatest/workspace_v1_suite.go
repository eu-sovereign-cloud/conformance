package secatest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/google/uuid"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type WorkspaceV1TestSuite struct {
	regionalTestSuite
}

func (suite *WorkspaceV1TestSuite) generateLifecycleParams() *secalib.WorkspaceLifeCycleParamsV1 {
	workspaceName := secalib.GenerateWorkspaceName()

	return &secalib.WorkspaceLifeCycleParamsV1{
		Workspace: &secalib.ResourceParams[secalib.WorkspaceSpecV1]{
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
	}
}

func (suite *WorkspaceV1TestSuite) TestWorkspaceLifeCycleV1(t provider.T) {
	slog.Info("Starting Workspace Lifecycle Test")

	t.Title("Workspace Lifecycle Test")
	configureTags(t, secalib.WorkspaceProviderV1, secalib.WorkspaceKind)

	params := suite.generateLifecycleParams()
	workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, params.Workspace.Name)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		scenarios := mock.NewWorkspaceScenariosV1(suite.authToken, suite.tenant, suite.region, suite.mockServerURL)

		id := uuid.New().String()
		wm, err := scenarios.ConfigureLifecycleScenario(id, params)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var resp *workspace.Workspace
	var err error

	var expectedMeta *secalib.Metadata
	var expectedLabels *[]secalib.Label

	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Workspace.Name,
		}
		ws := &workspace.Workspace{}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, tref, ws, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta = &secalib.Metadata{
			Name:       params.Workspace.Name,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
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
			Name:   params.Workspace.Name,
		}
		resp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		verifyWorkspaceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		expectedLabels = &[]secalib.Label{
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
		verifyLabelStep(sCtx, expectedLabels, actualLabels)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
	})

	t.WithNewStep("Update workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Workspace.Name,
		}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, tref, resp, nil)
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
			Name:   params.Workspace.Name,
		}
		resp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		verifyWorkspaceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		expectedLabels = &[]secalib.Label{
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
		verifyLabelStep(sCtx, expectedLabels, actualLabels)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
	})

	t.WithNewStep("Delete workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "DeleteWorkspace")

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Workspace.Name,
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
		Region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}
