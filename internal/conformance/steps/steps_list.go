package steps

import (
	"context"

	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Params

type listTenantResourcesParams[R types.ResourceType, M types.MetadataType] struct {
	listResourcesParams[R, M, secapi.TenantID]
	stepName       string
	stepParamsFunc func(provider.StepCtx, string)
	operationName  string
}

type listResourcesParams[R types.ResourceType, M types.MetadataType, F types.ReferenceType] struct {
	reference   F
	listOptions *secapi.ListOptions
	listFunc    func(context.Context, F, *secapi.ListOptions) (*secapi.Iterator[R], error)
}

// Steps

func listTenantResourcesStep[R types.ResourceType, M types.MetadataType](
	t provider.T,
	params listTenantResourcesParams[R, M],
) []*R {
	var resp []*R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)
		resp = listResources(t, sCtx, params.listResourcesParams)
	})
	return resp
}
func listResources[R types.ResourceType, M types.MetadataType, F types.ReferenceType](
	t provider.T,
	sCtx provider.StepCtx,
	params listResourcesParams[R, M, F],
) []*R {
	requestResourceStep(sCtx, params.reference)
	resp, err := params.listFunc(t.Context(), params.reference, params.listOptions)

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
