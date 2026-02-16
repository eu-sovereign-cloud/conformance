package steps

import (
	"context"
	"log/slog"
	"time"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Params

type getTenantResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string)
	operationName         string
	tref                  secapi.TenantReference
	getFunc               func(context.Context, secapi.TenantReference, secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[R, M, E, S], error)
	expectedResourceState schema.ResourceState
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
}

type getTenantResourceWithErrorParams struct {
	stepName       string
	stepParamsFunc func(provider.StepCtx, string)
	operationName  string
	tref           secapi.TenantReference
	getFunc        func(context.Context, secapi.TenantReference) error
	expectedError  error
}

type getWorkspaceResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string)
	operationName         string
	wref                  secapi.WorkspaceReference
	getFunc               func(context.Context, secapi.WorkspaceReference, secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[R, M, E, S], error)
	expectedResourceState schema.ResourceState
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
}

type getWorkspaceResourceWithErrorParams struct {
	stepName       string
	stepParamsFunc func(provider.StepCtx, string, string)
	operationName  string
	wref           secapi.WorkspaceReference
	getFunc        func(context.Context, secapi.WorkspaceReference) error
	expectedError  error
}

type getNetworkResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string, string)
	operationName         string
	nref                  secapi.NetworkReference
	getFunc               func(context.Context, secapi.NetworkReference, secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[R, M, E, S], error)
	expectedResourceState schema.ResourceState
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
}

type getNetworkResourceWithErrorParams struct {
	stepName       string
	stepParamsFunc func(provider.StepCtx, string, string, string)
	operationName  string
	nref           secapi.NetworkReference
	getFunc        func(context.Context, secapi.NetworkReference) error
	expectedError  error
}

type getResourceWithObserverParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType, F secapi.Reference, V any] struct {
	reference             F
	observerExpectedValue V
	getFunc               func(context.Context, F, secapi.ResourceObserverConfig[V]) (*stepFuncResponse[R, M, E, S], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

type getResourceWithErrorParams[R secapi.Reference] struct {
	reference     R
	getFunc       func(context.Context, R) error
	expectedError error
}

// Steps

func getTenantResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T,
	suite *suites.TestSuite,
	params getTenantResourceParams[R, M, E, S],
) *R {
	var resp *stepFuncResponse[R, M, E, S]
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		resp = getResourceWithObserver(t, suite, sCtx,
			getResourceWithObserverParams[R, M, E, S, secapi.TenantReference, schema.ResourceState]{
				reference:             params.tref,
				observerExpectedValue: params.expectedResourceState,
				getFunc:               params.getFunc,
				expectedLabels:        params.expectedLabels,
				expectedMetadata:      params.expectedMetadata,
				verifyMetadataFunc:    params.verifyMetadataFunc,
				expectedSpec:          params.expectedSpec,
				verifySpecFunc:        params.verifySpecFunc,
				expectedResourceState: params.expectedResourceState,
			},
		)
		requireNotNilResponse(sCtx, resp)
	})
	return resp.resource
}

func getTenantResourceWithErrorStep(t provider.T, params getTenantResourceWithErrorParams) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		getResourceWithError(t, sCtx,
			getResourceWithErrorParams[secapi.TenantReference]{
				reference:     params.tref,
				getFunc:       params.getFunc,
				expectedError: params.expectedError,
			},
		)
	})
}

func getWorkspaceResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T,
	suite *suites.TestSuite,
	params getWorkspaceResourceParams[R, M, E, S],
) *R {
	var resp *stepFuncResponse[R, M, E, S]
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.wref.Workspace))

		resp = getResourceWithObserver(t, suite, sCtx,
			getResourceWithObserverParams[R, M, E, S, secapi.WorkspaceReference, schema.ResourceState]{
				reference:             params.wref,
				observerExpectedValue: params.expectedResourceState,
				getFunc:               params.getFunc,
				expectedLabels:        params.expectedLabels,
				expectedMetadata:      params.expectedMetadata,
				verifyMetadataFunc:    params.verifyMetadataFunc,
				expectedSpec:          params.expectedSpec,
				verifySpecFunc:        params.verifySpecFunc,
				expectedResourceState: params.expectedResourceState,
			},
		)
		requireNotNilResponse(sCtx, resp)
	})
	return resp.resource
}

func getWorkspaceResourceWithErrorStep(t provider.T, params getWorkspaceResourceWithErrorParams) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.wref.Workspace))

		getResourceWithError(t, sCtx,
			getResourceWithErrorParams[secapi.WorkspaceReference]{
				reference:     params.wref,
				getFunc:       params.getFunc,
				expectedError: params.expectedError,
			},
		)
	})
}

func getNetworkResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T,
	suite *suites.TestSuite,
	params getNetworkResourceParams[R, M, E, S],
) *R {
	var resp *stepFuncResponse[R, M, E, S]
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.nref.Workspace), string(params.nref.Network))

		resp = getResourceWithObserver(t, suite, sCtx,
			getResourceWithObserverParams[R, M, E, S, secapi.NetworkReference, schema.ResourceState]{
				reference:             params.nref,
				observerExpectedValue: params.expectedResourceState,
				getFunc:               params.getFunc,
				expectedLabels:        params.expectedLabels,
				expectedMetadata:      params.expectedMetadata,
				verifyMetadataFunc:    params.verifyMetadataFunc,
				expectedSpec:          params.expectedSpec,
				verifySpecFunc:        params.verifySpecFunc,
				expectedResourceState: params.expectedResourceState,
			},
		)
		requireNotNilResponse(sCtx, resp)
	})
	return resp.resource
}

func getNetworkResourceWithErrorStep(t provider.T, params getNetworkResourceWithErrorParams) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.nref.Workspace), string(params.nref.Network))

		getResourceWithError(t, sCtx,
			getResourceWithErrorParams[secapi.NetworkReference]{
				reference:     params.nref,
				getFunc:       params.getFunc,
				expectedError: params.expectedError,
			},
		)
	})
}

func getResourceWithObserver[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType, F secapi.Reference, V any](
	t provider.T,
	suite *suites.TestSuite,
	sCtx provider.StepCtx,
	params getResourceWithObserverParams[R, M, E, S, F, V],
) *stepFuncResponse[R, M, E, S] {
	config := secapi.ResourceObserverConfig[V]{
		ExpectedValue: params.observerExpectedValue,
		Delay:         time.Duration(suite.BaseDelay) * time.Second,
		Interval:      time.Duration(suite.BaseInterval) * time.Second,
		MaxAttempts:   suite.MaxAttempts,
	}

	requestResourceStep(sCtx, params.reference)
	resp, err := params.getFunc(t.Context(), params.reference, config)
	responseResourceStep(sCtx, resp)

	requireNoError(sCtx, err)
	requireNotNilResponse(sCtx, resp)

	// Label
	if params.expectedLabels != nil {
		suite.VerifyLabelsStep(sCtx, params.expectedLabels, resp.labels)
	}

	// Metadata
	if resp.metadata != nil && params.expectedMetadata != nil {
		params.verifyMetadataFunc(sCtx, params.expectedMetadata, resp.metadata)
	} else {
		slog.Error("Metadata verification failed: expected or actual metadata is nil")
		t.FailNow()
	}

	if params.expectedSpec != nil {
		params.verifySpecFunc(sCtx, params.expectedSpec, &resp.spec)
	}

	// Status
	if resp.status != nil && params.expectedResourceState != "" {
		suite.VerifyStatusStep(sCtx, params.expectedResourceState, types.GetStatusState(resp.status))
	} else {
		slog.Error("Status verification failed: expected or actual Status is nil")
		t.FailNow()
	}

	return resp
}

func getResourceWithError[F secapi.Reference](
	t provider.T, sCtx provider.StepCtx, params getResourceWithErrorParams[F],
) {
	requestResourceStep(sCtx, params.reference)
	err := params.getFunc(t.Context(), params.reference)
	emptyResponseStep(sCtx)

	requireError(sCtx, err, params.expectedError)
}
