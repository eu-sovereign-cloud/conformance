package secapi

import (
	"context"

	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	sdk "github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type WorkspaceTestSuite struct {
	suite.Suite
	client *sdk.RegionalClient
}

func (suite *WorkspaceTestSuite) SetupSuite() {
}

func (suite *WorkspaceTestSuite) TearDownSuite() {
}

func (suite *WorkspaceTestSuite) SetupTest() {
}

func (suite *WorkspaceTestSuite) TearDownTest() {
}

func (suite *WorkspaceTestSuite) TestCreateWorkspace(t provider.T) {
	ctx := context.Background()
	t.Title("Workspace Management")

	ws := &workspace.Workspace{
		Metadata: &workspace.RegionalResourceMetadata{
			Tenant: "my-tenant",
			Name:   "workspace-1",
		},
	}

	t.NewStep("Create a Workspace")
	resp, err := suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)
	t.Assert().NoError(err)
	t.Assert().NotNil(resp)
	t.Assert().Equal("creating", string(*resp.Status.State))

	t.NewStep("Update a Workspace")
	resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)
	t.Assert().NoError(err)
	t.Assert().NotNil(resp)
	t.Assert().Equal("updating", string(*resp.Status.State))

	t.NewStep("Delete a Workspace")
	err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
	t.Assert().NoError(err)

	t.NewStep("Redelete a Workspace")
	err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
	t.Assert().NoError(err)
}
