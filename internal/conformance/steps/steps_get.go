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

type getGlobalResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType] struct {
	stepName           string
	stepParamsFunc     func(provider.StepCtx, string)
	operationName      string
	resourceName       string
	getFunc            func(context.Context, string) (*globalStepFuncResponse[R, M, E], error)
	expectedMetadata   *M
	verifyMetadataFunc func(provider.StepCtx, *M, *M)
	expectedSpec       *E
	verifySpecFunc     func(provider.StepCtx, *E, *E)
}

type getTenantResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	getResourceWithObserverParams[R, M, E, S, secapi.TenantReference, schema.ResourceState]
	stepName       string
	stepParamsFunc func(provider.StepCtx, string)
	operationName  string
}

type getTenantResourceWithErrorParams struct {
	getResourceWithErrorParams[secapi.TenantReference]
	stepName       string
	stepParamsFunc func(provider.StepCtx, string)
	operationName  string
}

type getWorkspaceResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	getResourceWithObserverParams[R, M, E, S, secapi.WorkspaceReference, schema.ResourceState]
	stepName       string
	stepParamsFunc func(provider.StepCtx, string, string)
	operationName  string
}

type getWorkspaceResourceWithErrorParams struct {
	getResourceWithErrorParams[secapi.WorkspaceReference]
	stepName       string
	stepParamsFunc func(provider.StepCtx, string, string)
	operationName  string
}

type getNetworkResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	getResourceWithObserverParams[R, M, E, S, secapi.NetworkReference, schema.ResourceState]
	stepName       string
	stepParamsFunc func(provider.StepCtx, string, string, string)
	operationName  string
}

type getNetworkResourceWithErrorParams struct {
	getResourceWithErrorParams[secapi.NetworkReference]
	stepName       string
	stepParamsFunc func(provider.StepCtx, string, string, string)
	operationName  string
}

type getResourceWithObserverParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType, F secapi.Reference, V any] struct {
	reference             F
	observerExpectedValue V
	getFunc               func(context.Context, F, secapi.ResourceObserverConfig[V]) (*getStepFuncResponse[R, M, E, S], error)
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

func getGlobalResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType](
	t provider.T,
	params getGlobalResourceParams[R, M, E],
) *R {
	var err error
	var resp *globalStepFuncResponse[R, M, E]
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		emptyRequestStep(sCtx)
		resp, err = params.getFunc(t.Context(), params.resourceName)

		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		responseResourceStep(sCtx, resp.resource)

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
	})
	return resp.resource
}

func getTenantResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T,
	suite *suites.TestSuite,
	params getTenantResourceParams[R, M, E, S],
) *R {
	var resp *getStepFuncResponse[R, M, E, S]
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)
		resp = getResourceWithObserver(t, suite, sCtx, params.getResourceWithObserverParams)
	})
	return resp.resource
}

func getTenantResourceWithErrorStep(t provider.T, params getTenantResourceWithErrorParams) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)
		getResourceWithError(t, sCtx, params.getResourceWithErrorParams)
	})
}

func getWorkspaceResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T,
	suite *suites.TestSuite,
	params getWorkspaceResourceParams[R, M, E, S],
) *R {
	var resp *getStepFuncResponse[R, M, E, S]
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.reference.Workspace))
		resp = getResourceWithObserver(t, suite, sCtx, params.getResourceWithObserverParams)
	})
	return resp.resource
}

func getWorkspaceResourceWithErrorStep(t provider.T, params getWorkspaceResourceWithErrorParams) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.reference.Workspace))
		getResourceWithError(t, sCtx, params.getResourceWithErrorParams)
	})
}

func getNetworkResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T,
	suite *suites.TestSuite,
	params getNetworkResourceParams[R, M, E, S],
) *R {
	var resp *getStepFuncResponse[R, M, E, S]
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.reference.Workspace), string(params.reference.Network))
		resp = getResourceWithObserver(t, suite, sCtx, params.getResourceWithObserverParams)
	})
	return resp.resource
}

func getNetworkResourceWithErrorStep(t provider.T, params getNetworkResourceWithErrorParams) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.reference.Workspace), string(params.reference.Network))
		getResourceWithError(t, sCtx, params.getResourceWithErrorParams)
	})
}

func getResourceWithObserver[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType, F secapi.Reference, V any](
	t provider.T,
	suite *suites.TestSuite,
	sCtx provider.StepCtx,
	params getResourceWithObserverParams[R, M, E, S, F, V],
) *getStepFuncResponse[R, M, E, S] {
	config := secapi.ResourceObserverConfig[V]{
		ExpectedValue: params.observerExpectedValue,
		Delay:         time.Duration(suite.BaseDelay) * time.Second,
		Interval:      time.Duration(suite.BaseInterval) * time.Second,
		MaxAttempts:   suite.MaxAttempts,
	}

	requestResourceStep(sCtx, params.reference)
	resp, err := params.getFunc(t.Context(), params.reference, config)

	requireNoError(sCtx, err)
	requireNotNilResponse(sCtx, resp)

	responseResourceStep(sCtx, resp.resource)

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
