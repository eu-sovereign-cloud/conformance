package steps

import (
	"context"

	"github.com/eu-sovereign-cloud/conformance/pkg/types"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Params

type actionWorkspaceResourceParams[R types.ResourceType] struct {
	stepName       string
	stepParamsFunc func(provider.StepCtx, string, string)
	operationName  string
	workspace      string
	resource       *R
	actionFunc     func(context.Context, *R) error
}

type actionResourceParams[R types.ResourceType] struct {
	resource   *R
	actionFunc func(context.Context, *R) error
}

// Steps

func actionWorkspaceResourceStep[R types.ResourceType](t provider.T, params actionWorkspaceResourceParams[R]) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace)

		actionResourceStep(t, sCtx, actionResourceParams[R]{
			resource:   params.resource,
			actionFunc: params.actionFunc,
		})
	})
}

func actionResourceStep[R types.ResourceType](t provider.T, sCtx provider.StepCtx, params actionResourceParams[R]) {
	requestResourceStep(sCtx, params.resource)
	err := params.actionFunc(t.Context(), params.resource)
	emptyResponseStep(sCtx)

	requireNoError(sCtx, err)
}
