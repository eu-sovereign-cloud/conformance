package suites

import (
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (suite *TestSuite) SetAuthorizationV1StepParams(sctx provider.StepCtx, operation string) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.AuthorizationProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
	)
}

func (suite *TestSuite) SetRegionV1StepParams(sctx provider.StepCtx, operation string) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.RegionProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
	)
}

func (suite *TestSuite) SetWorkspaceV1StepParams(sctx provider.StepCtx, operation string) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.WorkspaceProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
	)
}

func (suite *TestSuite) SetStorageV1StepParams(sctx provider.StepCtx, operation string) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.StorageProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
	)
}

func (suite *TestSuite) SetStorageWorkspaceV1StepParams(sctx provider.StepCtx, operation string, workspace string) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.StorageProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
		workspaceStepParameter, workspace,
	)
}

func (suite *TestSuite) SetComputeV1StepParams(sctx provider.StepCtx, operation string, workspace string) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.ComputeProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
		workspaceStepParameter, workspace,
	)
}

func (suite *TestSuite) SetNetworkV1StepParams(sctx provider.StepCtx, operation string, workspace string) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.NetworkProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
		workspaceStepParameter, workspace,
	)
}

// TODO Find a better name for this function
func (suite *TestSuite) SetNetworkNetworkV1StepParams(sctx provider.StepCtx, operation string, workspace string, network string) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.NetworkProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
		workspaceStepParameter, workspace,
		networkStepParameter, network,
	)
}
