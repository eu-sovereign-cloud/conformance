package steps

import (
	"context"
	"log"
	"time"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Steps

/// Create Or Update

type createOrUpdateTenantResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string)
	operationName         string
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E, S], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

type createOrUpdateWorkspaceResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string)
	operationName         string
	workspace             string
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E, S], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

type createOrUpdateNetworkResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string, string)
	operationName         string
	workspace             string
	network               string
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E, S], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

type createOrUpdateResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E, S], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

/// Get

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

/// Response

type stepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	resource *R
	labels   schema.Labels
	metadata *M
	spec     E
	status   *S
}

func newStepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	resource *R, labels schema.Labels, metadata *M, spec E, status *S,
) *stepFuncResponse[R, M, E, S] {
	return &stepFuncResponse[R, M, E, S]{
		resource: resource,
		labels:   labels,
		metadata: metadata,
		spec:     spec,
		status:   status,
	}
}

type ResponseExpects[M types.MetadataType, E types.SpecType] struct {
	Labels        schema.Labels
	Metadata      *M
	Spec          *E
	ResourceState schema.ResourceState
}

func createOrUpdateTenantResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T,
	suite *suites.TestSuite,
	params createOrUpdateTenantResourceParams[R, M, E, S],
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		createOrUpdateResourceStep(t, suite, sCtx, createOrUpdateResourceParams[R, M, E, S]{
			resource:              params.resource,
			createOrUpdateFunc:    params.createOrUpdateFunc,
			expectedLabels:        params.expectedLabels,
			expectedMetadata:      params.expectedMetadata,
			verifyMetadataFunc:    params.verifyMetadataFunc,
			expectedSpec:          params.expectedSpec,
			verifySpecFunc:        params.verifySpecFunc,
			expectedResourceState: params.expectedResourceState,
		})
	})
}

func createOrUpdateWorkspaceResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T,
	suite *suites.TestSuite,
	params createOrUpdateWorkspaceResourceParams[R, M, E, S],
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace)

		createOrUpdateResourceStep(t, suite, sCtx, createOrUpdateResourceParams[R, M, E, S]{
			resource:              params.resource,
			createOrUpdateFunc:    params.createOrUpdateFunc,
			expectedLabels:        params.expectedLabels,
			expectedMetadata:      params.expectedMetadata,
			verifyMetadataFunc:    params.verifyMetadataFunc,
			expectedSpec:          params.expectedSpec,
			verifySpecFunc:        params.verifySpecFunc,
			expectedResourceState: params.expectedResourceState,
		})
	})
}

func createOrUpdateNetworkResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T,
	suite *suites.TestSuite,
	params createOrUpdateNetworkResourceParams[R, M, E, S],
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace, params.network)

		createOrUpdateResourceStep(t, suite, sCtx, createOrUpdateResourceParams[R, M, E, S]{
			resource:              params.resource,
			createOrUpdateFunc:    params.createOrUpdateFunc,
			expectedLabels:        params.expectedLabels,
			expectedMetadata:      params.expectedMetadata,
			verifyMetadataFunc:    params.verifyMetadataFunc,
			expectedSpec:          params.expectedSpec,
			verifySpecFunc:        params.verifySpecFunc,
			expectedResourceState: params.expectedResourceState,
		})
	})
}

func createOrUpdateResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T,
	suite *suites.TestSuite,
	sCtx provider.StepCtx,
	params createOrUpdateResourceParams[R, M, E, S],
) {
	resp, err := params.createOrUpdateFunc(t.Context(), params.resource)
	requireNoError(sCtx, err)

	// Label
	if params.expectedLabels != nil {
		suite.VerifyLabelsStep(sCtx, params.expectedLabels, resp.labels)
	}

	// Metadata
	if resp.metadata != nil && params.expectedMetadata != nil {
		params.verifyMetadataFunc(sCtx, params.expectedMetadata, resp.metadata)
	} else {
		log.Fatalln("Metadata verification failed: expected or actual metadata is nil")
	}

	if params.expectedSpec != nil {
		params.verifySpecFunc(sCtx, params.expectedSpec, &resp.spec)
	}

	// Status
	if resp.status != nil && params.expectedResourceState != "" {
		suite.VerifyStatusStep(sCtx, params.expectedResourceState, types.GetStatusState(resp.status))
	} else {
		log.Fatalln("Status verification failed: expected or actual Status is nil")
	}
}

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

	resp, err := params.getFunc(t.Context(), params.reference, config)
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
		log.Fatalln("Metadata verification failed: expected or actual metadata is nil")
	}

	if params.expectedSpec != nil {
		params.verifySpecFunc(sCtx, params.expectedSpec, &resp.spec)
	}

	// Status
	if resp.status != nil && params.expectedResourceState != "" {
		suite.VerifyStatusStep(sCtx, params.expectedResourceState, types.GetStatusState(resp.status))
	} else {
		log.Fatalln("Status verification failed: expected or actual Status is nil")
	}

	return resp
}
