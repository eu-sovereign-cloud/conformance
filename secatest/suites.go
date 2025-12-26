package secatest

import (
	"log/slog"
	"regexp"
	"strings"

	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/wiremock/go-wiremock"
)

type testSuite struct {
	suite.Suite
	tenant        string
	authToken     string
	mockEnabled   bool
	mockServerURL *string

	mockClient   *wiremock.Client
	scenarioName string

	baseDelay    int
	baseInterval int
	maxAttempts  int
}

func (suite *testSuite) canRun(regexp *regexp.Regexp) bool {
	if regexp == nil {
		return true
	} else {
		return regexp.MatchString(suite.scenarioName)
	}
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
	// Cleanup configured mock scenarios
	if suite.mockClient != nil {
		if err := suite.mockClient.ResetAllScenarios(); err != nil {
			slog.Error("Failed to reset scenarios", "error", err)
		}
	}
}

func (suite *testSuite) setAuthorizationV1StepParams(sctx provider.StepCtx, operation string) {
	sctx.WithNewParameters(
		providerStepParameter, authorizationProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
	)
}

func (suite *testSuite) setRegionV1StepParams(sctx provider.StepCtx, operation string) {
	sctx.WithNewParameters(
		providerStepParameter, regionProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
	)
}

func (suite *testSuite) setWorkspaceV1StepParams(sctx provider.StepCtx, operation string) {
	sctx.WithNewParameters(
		providerStepParameter, workspaceProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
	)
}

func (suite *testSuite) setStorageV1StepParams(sctx provider.StepCtx, operation string) {
	sctx.WithNewParameters(
		providerStepParameter, storageProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
	)
}

func (suite *testSuite) setStorageWorkspaceV1StepParams(sctx provider.StepCtx, operation string, workspace string) {
	sctx.WithNewParameters(
		providerStepParameter, storageProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
		workspaceStepParameter, workspace,
	)
}

func (suite *testSuite) setComputeV1StepParams(sctx provider.StepCtx, operation string, workspace string) {
	sctx.WithNewParameters(
		providerStepParameter, computeProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
		workspaceStepParameter, workspace,
	)
}

func (suite *testSuite) setNetworkV1StepParams(sctx provider.StepCtx, operation string, workspace string) {
	sctx.WithNewParameters(
		providerStepParameter, networkProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
		workspaceStepParameter, workspace,
	)
}

// TODO Find a better name for this function
func (suite *testSuite) setNetworkNetworkV1StepParams(sctx provider.StepCtx, operation string, workspace string, network string) {
	sctx.WithNewParameters(
		providerStepParameter, networkProviderV1,
		operationStepParameter, operation,
		tenantStepParameter, suite.tenant,
		workspaceStepParameter, workspace,
		networkStepParameter, network,
	)
}
