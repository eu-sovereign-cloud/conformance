package secatest

import (
	"log/slog"
	"strings"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/wiremock/go-wiremock"
)

type testSuite struct {
	suite.Suite
	tenant        string
	authToken     string
	mockEnabled   string
	mockServerURL *string

	mockClient *wiremock.Client
}

func (suite *testSuite) isMockEnabled() bool {
	return suite.mockEnabled == "true"
}

type mixedTestSuite struct {
	testSuite

	globalClient *secapi.GlobalClient

	region         string
	regionalClient *secapi.RegionalClient
}

type globalTestSuite struct {
	testSuite
	client *secapi.GlobalClient
}

type regionalTestSuite struct {
	testSuite
	region string
	client *secapi.RegionalClient
}

func configureTags(t provider.T, provider string, kinds ...string) {
	t.Tags(
		"provider:"+provider,
		"resources:"+strings.Join(kinds, ", "),
	)
}

func (suite *testSuite) resetAllScenarios() {
	if suite.mockClient != nil {
		if err := suite.mockClient.ResetAllScenarios(); err != nil {
			slog.Error("Failed to reset scenarios", "error", err)
		}
	}
}

func (suite *testSuite) setAuthorizationV1StepParams(sctx provider.StepCtx, operation string) {
	sctx.WithNewParameters(
		providerStepParameter, secalib.AuthorizationProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
	)
}

func (suite *testSuite) setRegionV1StepParams(sctx provider.StepCtx, operation string) {
	sctx.WithNewParameters(
		providerStepParameter, secalib.RegionProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
	)
}

func (suite *testSuite) setWorkspaceV1StepParams(sctx provider.StepCtx, operation string) {
	sctx.WithNewParameters(
		providerStepParameter, secalib.WorkspaceProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
	)
}

func (suite *testSuite) setStorageV1StepParams(sctx provider.StepCtx, operation string, workspace string) {
	if workspace != "" {
		sctx.WithNewParameters(
			providerStepParameter, secalib.StorageProviderV1,
			operationStepParameter, operation,
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspace,
		)
	} else {
		sctx.WithNewParameters(
			providerStepParameter, secalib.StorageProviderV1,
			operationStepParameter, operation,
			tenantStepParameter, suite.tenant,
		)
	}
}

func (suite *testSuite) setComputeV1StepParams(sctx provider.StepCtx, operation string, workspace string) {
	sctx.WithNewParameters(
		providerStepParameter, secalib.ComputeProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
		workspaceStepParameter, workspace,
	)
}

func (suite *testSuite) setNetworkV1StepParams(sctx provider.StepCtx, operation string, workspace string) {
	sctx.WithNewParameters(
		providerStepParameter, secalib.NetworkProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
		workspaceStepParameter, workspace,
	)
}
