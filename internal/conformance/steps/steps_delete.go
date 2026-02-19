package steps

import (
	"context"

	"github.com/eu-sovereign-cloud/conformance/pkg/types"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Params

type deleteTenantResourceParams[R types.ResourceType] struct {
	deleteResourceParams[R]
	stepName       string
	stepParamsFunc func(provider.StepCtx, string)
	operationName  string
}

type deleteWorkspaceResourceParams[R types.ResourceType] struct {
	deleteResourceParams[R]
	stepName       string
	stepParamsFunc func(provider.StepCtx, string, string)
	operationName  string
	workspace      string
}

type deleteNetworkResourceParams[R types.ResourceType] struct {
	deleteResourceParams[R]
	stepName       string
	stepParamsFunc func(provider.StepCtx, string, string, string)
	operationName  string
	workspace      string
	network        string
}

type deleteResourceParams[R types.ResourceType] struct {
	resource   *R
	deleteFunc func(context.Context, *R) error
}

// Steps

func deleteTenantResourceStep[R types.ResourceType](t provider.T, params deleteTenantResourceParams[R]) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)
		deleteResourceStep(t, sCtx, params.deleteResourceParams)
	})
}

func deleteWorkspaceResourceStep[R types.ResourceType](t provider.T, params deleteWorkspaceResourceParams[R]) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace)
		deleteResourceStep(t, sCtx, params.deleteResourceParams)
	})
}

func deleteNetworkResourceStep[R types.ResourceType](t provider.T, params deleteNetworkResourceParams[R]) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace, params.network)
		deleteResourceStep(t, sCtx, params.deleteResourceParams)
	})
}

func deleteResourceStep[R types.ResourceType](t provider.T, sCtx provider.StepCtx, params deleteResourceParams[R]) {
	requestResourceStep(sCtx, params.resource)
	err := params.deleteFunc(t.Context(), params.resource)
	emptyResponseStep(sCtx)

	requireNoError(sCtx, err)
}
