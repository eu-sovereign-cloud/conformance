package steps

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StepResponseExpects[M types.MetadataType, E types.SpecType] struct {
	Labels         schema.Labels
	Metadata       *M
	Spec           *E
	ResourceStates []schema.ResourceState
}

// Watch

type watchTenantResourceUntilDeletedParams struct {
	watchResourceUntilDeletedParams[secapi.TenantReference]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName)
	operationName  constants.OperationName
}

func watchTenantResourceUntilDeletedStep(
	t provider.T, suite *suites.TestSuite, params watchTenantResourceUntilDeletedParams,
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)
		watchResourceUntilDeletedStep(t, suite, sCtx, params.watchResourceUntilDeletedParams)
	})
}

type watchWorkspaceResourceUntilDeletedParams struct {
	watchResourceUntilDeletedParams[secapi.WorkspaceReference]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID)
	operationName  constants.OperationName
}

func watchWorkspaceResourceUntilDeletedStep(
	t provider.T, suite *suites.TestSuite, params watchWorkspaceResourceUntilDeletedParams,
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.reference.Workspace)
		watchResourceUntilDeletedStep(t, suite, sCtx, params.watchResourceUntilDeletedParams)
	})
}

type watchNetworkResourceUntilDeletedParams struct {
	watchResourceUntilDeletedParams[secapi.NetworkReference]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID, secapi.NetworkID)
	operationName  constants.OperationName
}

func watchNetworkResourceUntilDeletedStep(
	t provider.T, suite *suites.TestSuite, params watchNetworkResourceUntilDeletedParams,
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.reference.Workspace, params.reference.Network)
		watchResourceUntilDeletedStep(t, suite, sCtx, params.watchResourceUntilDeletedParams)
	})
}

type watchResourceUntilDeletedParams[F secapi.Reference] struct {
	reference    F
	getErrorFunc func(context.Context, F, secapi.ResourceObserverConfig) error
}

func watchResourceUntilDeletedStep[F secapi.Reference](
	t provider.T, suite *suites.TestSuite, sCtx provider.StepCtx, params watchResourceUntilDeletedParams[F],
) {
	config := secapi.ResourceObserverConfig{
		Delay:       time.Duration(suite.BaseDelay) * time.Second,
		Interval:    time.Duration(suite.BaseInterval) * time.Second,
		MaxAttempts: suite.MaxAttempts,
	}

	err := params.getErrorFunc(t.Context(), params.reference, config)
	requireNoError(sCtx, err)
}

// API Request / Response

func requestResourceStep(ctx provider.StepCtx, request any) {
	ctx.WithNewStep("Send request", func(stepCtx provider.StepCtx) {
		if request == nil {
			return
		}

		if data, err := json.Marshal(request); err != nil {
			slog.Error("Error marshaling resource to json", "error", err)
			stepCtx.FailNow()
		} else {
			stepCtx.WithNewParameters("resource", string(data))
		}
	})
}

func emptyRequestStep(ctx provider.StepCtx) {
	ctx.WithNewStep("Send request", func(stepCtx provider.StepCtx) {})
}

func responseResourceStep[R types.ResourceType](ctx provider.StepCtx, resource *R) {
	ctx.WithNewStep("Receive response", func(stepCtx provider.StepCtx) {
		if resource == nil {
			return
		}

		if data, err := json.Marshal(resource); err != nil {
			slog.Error("Error marshaling resource to json", "error", err)
			stepCtx.FailNow()
		} else {
			stepCtx.WithNewParameters("resource", string(data))
		}
	})
}

func responseResourcesStep[R types.ResourceType](ctx provider.StepCtx, resources []*R) {
	ctx.WithNewStep("Receive response", func(stepCtx provider.StepCtx) {
		if resources == nil {
			return
		}

		if data, err := json.Marshal(resources); err != nil {
			slog.Error("Error marshaling resource to json", "error", err)
			stepCtx.FailNow()
		} else {
			stepCtx.WithNewParameters("resources", string(data))
		}
	})
}

func emptyResponseStep(ctx provider.StepCtx) {
	ctx.WithNewStep("Receive response", func(stepCtx provider.StepCtx) {})
}
