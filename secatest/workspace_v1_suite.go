package secatest

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type WorkspaceV1TestSuite struct {
	suite.Suite
	client     *secapi.RegionalClient
	tenant     string
	mockParams *mock.MockParams
}

func (suite *WorkspaceV1TestSuite) TestWorkspaceV1(t provider.T) {
	t.Title("Workspace Lifecycle Test")
	t.Tags(
		"provider:workspace",
		"resource:workspace",
		"version:v1",
	)

	// Generate scenario data
	workspaceName := fmt.Sprintf("workspace-%d", rand.Intn(math.MaxInt32))

	// Setup mock, if configured to use
	if suite.mockParams != nil {
		if err := mock.CreateWorkspaceScenario(mock.MockParams{
			WireMockURL:   suite.mockParams.WireMockURL,
			TenantName:    suite.mockParams.TenantName,
			WorkspaceName: workspaceName,
			Region:        suite.mockParams.Region,
			Token:         suite.mockParams.Token,
		}); err != nil {
			slog.Error("Failed to create workspace scenario", "error", err)
			return
		}
	}

	ctx := context.Background()
	var resp *workspace.Workspace
	var err error

	// Step 1
	t.WithNewStep("Create Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			"operation", "CreateOrUpdateWorkspace",
			"tenant", suite.tenant,
			"workspace", workspaceName,
		)

		ws := &workspace.Workspace{
			Metadata: &workspace.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)

		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireStatusEquals(sCtx, statusStateCreating, string(*resp.Status.State))

		sCtx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
			stepCtx.WithNewParameters(
				"expected_tenant", suite.tenant,
				"actual_tenant", resp.Metadata.Tenant,
				"expected_name", workspaceName,
				"actual_name", resp.Metadata.Name,
			)
			stepCtx.Assert().Equal(suite.tenant, resp.Metadata.Tenant, "Tenant should match expected")
			stepCtx.Assert().Equal(workspaceName, resp.Metadata.Name, "Name should match expected")
		})
	})

	// Step 2
	t.WithNewStep("Update Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			"operation", "CreateOrUpdateWorkspace",
			"tenant", suite.tenant,
			"workspace", workspaceName,
		)

		ws := &workspace.Workspace{
			Metadata: &workspace.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)

		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireStatusEquals(sCtx, statusStateUpdating, string(*resp.Status.State))

		sCtx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
			stepCtx.WithNewParameters(
				"expected_tenant", suite.tenant,
				"actual_tenant", resp.Metadata.Tenant,
				"expected_name", workspaceName,
				"actual_name", resp.Metadata.Name,
			)
			stepCtx.Assert().Equal(suite.tenant, resp.Metadata.Tenant, "Tenant should match expected")
			stepCtx.Assert().Equal(workspaceName, resp.Metadata.Name, "Name should match expected")
		})
	})

	// Step 3
	t.WithNewStep("Delete Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			"operation", "DeleteWorkspace",
			"tenant", suite.tenant,
			"workspace", workspaceName,
		)

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
		requireNoError(sCtx, err)
	})

	// Step 4
	t.WithNewStep("Re-Delete Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			"operation", "DeleteWorkspace",
			"tenant", suite.tenant,
			"workspace", workspaceName,
		)

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
		requireNoError(sCtx, err)
	})
}
