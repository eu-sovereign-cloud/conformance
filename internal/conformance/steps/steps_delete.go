package steps

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Params

type deleteTenantResourceParams[R types.ResourceType] struct {
	deleteResourceParams[R]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName)
	operationName  constants.OperationName
}

type deleteWorkspaceResourceParams[R types.ResourceType] struct {
	deleteResourceParams[R]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID)
	operationName  constants.OperationName
	workspace      secapi.WorkspaceID
}

type deleteNetworkResourceParams[R types.ResourceType] struct {
	deleteResourceParams[R]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID, secapi.NetworkID)
	operationName  constants.OperationName
	workspace      secapi.WorkspaceID
	network        secapi.NetworkID
}

type deleteResourceParams[R types.ResourceType] struct {
	resource   *R
	deleteFunc func(context.Context, *R) error
}

// Steps

func deleteTenantResourceStep[R types.ResourceType](t provider.T, suite *suites.TestSuite, params deleteTenantResourceParams[R]) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)
		deleteResourceStep(t, suite, params.stepName, sCtx, params.deleteResourceParams)
	})
}

func deleteWorkspaceResourceStep[R types.ResourceType](t provider.T, suite *suites.TestSuite, params deleteWorkspaceResourceParams[R]) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace)
		deleteResourceStep(t, suite, params.stepName, sCtx, params.deleteResourceParams)
	})
}

func deleteNetworkResourceStep[R types.ResourceType](t provider.T, suite *suites.TestSuite, params deleteNetworkResourceParams[R]) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace, params.network)
		deleteResourceStep(t, suite, params.stepName, sCtx, params.deleteResourceParams)
	})
}

func deleteResourceStep[R types.ResourceType](t provider.T, suite *suites.TestSuite, stepName string, sCtx provider.StepCtx, params deleteResourceParams[R]) {
	slog.Info(fmt.Sprintf("[%s] %s", suite.ScenarioName, stepName))

	requestResourceStep(sCtx, params.resource)
	err := params.deleteFunc(t.Context(), params.resource)
	emptyResponseStep(sCtx)

	requireNoError(sCtx, err)
}
