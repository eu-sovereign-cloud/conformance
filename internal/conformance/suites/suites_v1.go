package suites

import (
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (suite *TestSuite) SetAuthorizationV1StepParams(sctx provider.StepCtx, operation constants.OperationName) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.AuthorizationProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
	)
}

func (suite *TestSuite) SetRegionV1StepParams(sctx provider.StepCtx, operation constants.OperationName) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.RegionProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
	)
}

func (suite *TestSuite) SetWorkspaceV1StepParams(sctx provider.StepCtx, operation constants.OperationName) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.WorkspaceProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
	)
}

func (suite *TestSuite) SetStorageV1StepParams(sctx provider.StepCtx, operation constants.OperationName) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.StorageProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
	)
}

func (suite *TestSuite) SetStorageWorkspaceV1StepParams(sctx provider.StepCtx, operation constants.OperationName, workspace secapi.WorkspaceID) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.StorageProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
		workspaceStepParameter, workspace,
	)
}

func (suite *TestSuite) SetComputeV1StepParams(sctx provider.StepCtx, operation constants.OperationName, workspace secapi.WorkspaceID) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.ComputeProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
		workspaceStepParameter, workspace,
	)
}

func (suite *TestSuite) SetComputeSkuV1StepParams(sctx provider.StepCtx, operation constants.OperationName) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.ComputeProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
	)
}

func (suite *TestSuite) SetNetworkV1StepParams(sctx provider.StepCtx, operation constants.OperationName, workspace secapi.WorkspaceID) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.NetworkProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
		workspaceStepParameter, workspace,
	)
}

func (suite *TestSuite) SetNetworkSkuV1StepParams(sctx provider.StepCtx, operation constants.OperationName) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.NetworkProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
	)
}

// TODO Find a better name for this function
func (suite *TestSuite) SetNetworkNetworkV1StepParams(sctx provider.StepCtx, operation constants.OperationName, workspace secapi.WorkspaceID, network secapi.NetworkID) {
	sctx.WithNewParameters(
		providerStepParameter, sdkconsts.NetworkProviderV1Name,
		operationStepParameter, operation,
		tenantStepParameter, suite.Tenant,
		workspaceStepParameter, workspace,
		networkStepParameter, network,
	)
}
