package secatest

import (
	"context"

	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type WorkspaceV1TestSuite struct {
	suite.Suite
	client *secapi.RegionalClient
}

func (suite *WorkspaceV1TestSuite) TestWorkspaceV1(t provider.T) {
	t.Title("Workspace Lifecycle Test")
	t.Tags(
		"provider:workspace",
		"resource:workspace",
		"version:v1",
	)

	ctx := context.Background()
	var resp *workspace.Workspace
	var err error

	// Step 1
	t.WithNewStep("Create Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			"operation", "CreateOrUpdateWorkspace",
			"tenant", tenant1Name,
			"workspace", workspace1Name,
		)

		ws := &workspace.Workspace{
			Metadata: &workspace.RegionalResourceMetadata{
				Tenant: tenant1Name,
				Name:   workspace1Name,
			},
		}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)

		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireStatusEquals(sCtx, statusStateCreating, string(*resp.Status.State))

		sCtx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
			stepCtx.WithNewParameters(
				"expected_tenant", tenant1Name,
				"actual_tenant", resp.Metadata.Tenant,
				"expected_name", workspace1Name,
				"actual_name", resp.Metadata.Name,
			)
			stepCtx.Assert().Equal(tenant1Name, resp.Metadata.Tenant, "Tenant should match expected")
			stepCtx.Assert().Equal(workspace1Name, resp.Metadata.Name, "Name should match expected")
		})
	})

	// Step 2
	t.WithNewStep("Update Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			"operation", "CreateOrUpdateWorkspace",
			"tenant", tenant1Name,
			"workspace", workspace1Name,
		)

		ws := &workspace.Workspace{
			Metadata: &workspace.RegionalResourceMetadata{
				Tenant: tenant1Name,
				Name:   workspace1Name,
			},
		}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)

		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireStatusEquals(sCtx, statusStateUpdating, string(*resp.Status.State))

		sCtx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
			stepCtx.WithNewParameters(
				"expected_tenant", tenant1Name,
				"actual_tenant", resp.Metadata.Tenant,
				"expected_name", workspace1Name,
				"actual_name", resp.Metadata.Name,
			)
			stepCtx.Assert().Equal(tenant1Name, resp.Metadata.Tenant, "Tenant should match expected")
			stepCtx.Assert().Equal(workspace1Name, resp.Metadata.Name, "Name should match expected")
		})
	})

	// Step 3
	t.WithNewStep("Delete Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			"operation", "DeleteWorkspace",
			"tenant", tenant1Name,
			"workspace", workspace1Name,
		)

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
		requireNoError(sCtx, err)
	})

	// Step 4
	t.WithNewStep("Re-Delete Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			"operation", "DeleteWorkspace",
			"tenant", tenant1Name,
			"workspace", workspace1Name,
		)

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
		requireNoError(sCtx, err)
	})
}
