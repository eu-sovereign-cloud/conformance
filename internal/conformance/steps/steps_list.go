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

type listGlobalResourcesParams[R types.ResourceType, M types.MetadataType] struct {
	listOptions    *secapi.ListOptions
	listFunc       func(context.Context, *secapi.ListOptions) (*secapi.Iterator[R], error)
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName)
	operationName  constants.OperationName
}

type listTenantResourcesParams[R types.ResourceType, M types.MetadataType] struct {
	listResourcesParams[R, M, secapi.TenantPath]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName)
	operationName  constants.OperationName
}

type listWorkspaceResourcesParams[R types.ResourceType, M types.MetadataType] struct {
	listResourcesParams[R, M, secapi.WorkspacePath]
	stepName       string
	workspace      secapi.WorkspaceID
	stepParamsFunc func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID)
	operationName  constants.OperationName
}

type listNetworkResourcesParams[R types.ResourceType, M types.MetadataType] struct {
	listResourcesParams[R, M, secapi.NetworkPath]
	stepName       string
	workspace      secapi.WorkspaceID
	network        secapi.NetworkID
	stepParamsFunc func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID, secapi.NetworkID)
	operationName  constants.OperationName
}

type listResourcesParams[R types.ResourceType, M types.MetadataType, P types.PathType] struct {
	path        P
	listOptions *secapi.ListOptions
	listFunc    func(context.Context, P, *secapi.ListOptions) (*secapi.Iterator[R], error)
}

// Steps

func listGlobalResourcesStep[R types.ResourceType, M types.MetadataType](
	t provider.T, suite *suites.TestSuite, stepName string, params listGlobalResourcesParams[R, M],
) []*R {
	var items []*R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		slog.Info(fmt.Sprintf("[%s] %s", suite.ScenarioName, stepName))

		emptyRequestStep(sCtx)
		resp, err := params.listFunc(t.Context(), params.listOptions)

		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		// Items
		items, err = resp.All(t.Context())
		requireNoError(sCtx, err)

		responseResourcesStep(sCtx, items)

		requireNotNilResponse(sCtx, items)
		requireNotEmptyResponse(sCtx, items)
	})
	return items
}

func listTenantResourcesStep[R types.ResourceType, M types.MetadataType](
	t provider.T, suite *suites.TestSuite, params listTenantResourcesParams[R, M],
) []*R {
	var resp []*R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)
		resp = listResourcesStep(t, suite, params.stepName, sCtx, params.listResourcesParams)
	})
	return resp
}

func listWorkspaceResourcesStep[R types.ResourceType, M types.MetadataType](
	t provider.T, suite *suites.TestSuite, params listWorkspaceResourcesParams[R, M],
) []*R {
	var resp []*R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace)
		resp = listResourcesStep(t, suite, params.stepName, sCtx, params.listResourcesParams)
	})
	return resp
}

func listNetworkResourcesStep[R types.ResourceType, M types.MetadataType](
	t provider.T, suite *suites.TestSuite, params listNetworkResourcesParams[R, M],
) []*R {
	var resp []*R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace, params.network)
		resp = listResourcesStep(t, suite, params.stepName, sCtx, params.listResourcesParams)
	})
	return resp
}

func listResourcesStep[R types.ResourceType, M types.MetadataType, P types.PathType](
	t provider.T, suite *suites.TestSuite, stepName string, sCtx provider.StepCtx, params listResourcesParams[R, M, P],
) []*R {
	slog.Info(fmt.Sprintf("[%s] %s", suite.ScenarioName, stepName))

	requestResourceStep(sCtx, params.path)
	resp, err := params.listFunc(t.Context(), params.path, params.listOptions)

	requireNoError(sCtx, err)
	requireNotNilResponse(sCtx, resp)

	// Items
	items, err := resp.All(t.Context())
	requireNoError(sCtx, err)

	responseResourcesStep(sCtx, items)

	requireNotNilResponse(sCtx, items)
	requireNotEmptyResponse(sCtx, items)

	return items
}
