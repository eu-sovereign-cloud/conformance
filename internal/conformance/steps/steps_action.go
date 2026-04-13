package steps

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Params

type actionWorkspaceResourceParams[R types.ResourceType] struct {
	actionResourceParams[R]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID)
	operationName  constants.OperationName
	workspace      secapi.WorkspaceID
}

type actionResourceParams[R types.ResourceType] struct {
	resource   *R
	actionFunc func(context.Context, *R) error
}

// Steps

func actionWorkspaceResourceStep[R types.ResourceType](t provider.T, suite *suites.TestSuite, params actionWorkspaceResourceParams[R]) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace)
		actionResourceStep(t, suite, params.stepName, sCtx, params.actionResourceParams)
	})
}

func actionResourceStep[R types.ResourceType](t provider.T, suite *suites.TestSuite, stepName string, sCtx provider.StepCtx, params actionResourceParams[R]) {
	slog.Info(fmt.Sprintf("[%s] %s", suite.ScenarioName, stepName))

	resourceRequestStep(sCtx, params.resource)
	err := params.actionFunc(t.Context(), params.resource)
	emptyResponseStep(sCtx)

	requireNoError(sCtx, err)
}
