package steps

import (
	"context"
	"time"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Params

type watchTenantResourceUntilDeletedParams struct {
	watchResourceUntilDeletedParams[secapi.TenantReference]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName)
	operationName  constants.OperationName
}

type watchWorkspaceResourceUntilDeletedParams struct {
	watchResourceUntilDeletedParams[secapi.WorkspaceReference]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID)
	operationName  constants.OperationName
}

type watchNetworkResourceUntilDeletedParams struct {
	watchResourceUntilDeletedParams[secapi.NetworkReference]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID, secapi.NetworkID)
	operationName  constants.OperationName
}

type watchResourceUntilDeletedParams[F secapi.Reference] struct {
	reference    F
	getErrorFunc func(context.Context, F, secapi.ResourceObserverConfig) error
}

// Steps

func watchTenantResourceUntilDeletedStep(
	ctx context.Context, suite *suites.TestSuite, stepCreator StepCreator, params watchTenantResourceUntilDeletedParams,
) {
	stepCreator.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)
		watchResourceUntilDeletedStep(ctx, suite, sCtx, params.watchResourceUntilDeletedParams)
	})
}

func watchWorkspaceResourceUntilDeletedStep(
	ctx context.Context, suite *suites.TestSuite, stepCreator StepCreator, params watchWorkspaceResourceUntilDeletedParams,
) {
	stepCreator.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.reference.Workspace)
		watchResourceUntilDeletedStep(ctx, suite, sCtx, params.watchResourceUntilDeletedParams)
	})
}

func watchNetworkResourceUntilDeletedStep(
	ctx context.Context, suite *suites.TestSuite, stepCreator StepCreator, params watchNetworkResourceUntilDeletedParams,
) {
	stepCreator.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.reference.Workspace, params.reference.Network)
		watchResourceUntilDeletedStep(ctx, suite, sCtx, params.watchResourceUntilDeletedParams)
	})
}

func watchResourceUntilDeletedStep[F secapi.Reference](
	tctx context.Context, suite *suites.TestSuite, sCtx provider.StepCtx, params watchResourceUntilDeletedParams[F],
) {
	config := secapi.ResourceObserverConfig{
		Delay:       time.Duration(suite.BaseDelay) * time.Second,
		Interval:    time.Duration(suite.BaseInterval) * time.Second,
		MaxAttempts: suite.MaxAttempts,
	}
	referenceRequestStep(sCtx, params.reference)

	err := params.getErrorFunc(tctx, params.reference, config)
	requireNoError(sCtx, err)
}
