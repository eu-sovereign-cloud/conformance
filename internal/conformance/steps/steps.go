package steps

import (
	"context"
	"fmt"
	"time"

	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Steps

/// Create Or Update

type createOrUpdateTenantResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string)
	operationName         string
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

type createOrUpdateWorkspaceResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string)
	operationName         string
	workspace             string
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

type createOrUpdateNetworkResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string, string)
	operationName         string
	workspace             string
	network               string
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

type createOrUpdateResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType] struct {
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

/// Get

type getTenantResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string)
	operationName         string
	tref                  secapi.TenantReference
	getFunc               func(context.Context, secapi.TenantReference, secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[R, M, E], error)
	expectedResourceState schema.ResourceState
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
}

type getWorkspaceResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string)
	operationName         string
	wref                  secapi.WorkspaceReference
	getFunc               func(context.Context, secapi.WorkspaceReference, secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[R, M, E], error)
	expectedResourceState schema.ResourceState
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
}

type getNetworkResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string, string)
	operationName         string
	nref                  secapi.NetworkReference
	getFunc               func(context.Context, secapi.NetworkReference, secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[R, M, E], error)
	expectedResourceState schema.ResourceState
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
}

type getResourceWithObserverParams[R types.ResourceType, M types.MetadataType, E types.SpecType, F secapi.Reference, V any] struct {
	reference             F
	observerExpectedValue V
	getFunc               func(context.Context, F, secapi.ResourceObserverConfig[V]) (*stepFuncResponse[R, M, E], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

/// Response

type stepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType] struct {
	resource *R
	labels   schema.Labels
	metadata *M
	spec     E
	state    *schema.ResourceState
}

func newStepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType](resource *R, labels schema.Labels, metadata *M, spec E, state *schema.ResourceState) *stepFuncResponse[R, M, E] {
	return &stepFuncResponse[R, M, E]{
		resource: resource,
		labels:   labels,
		metadata: metadata,
		spec:     spec,
		state:    state,
	}
}

type ResponseExpects[M types.MetadataType, E types.SpecType] struct {
	Labels        schema.Labels
	Metadata      *M
	Spec          *E
	ResourceState schema.ResourceState
}

func createOrUpdateTenantResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType](
	configurator *StepsConfigurator,
	params createOrUpdateTenantResourceParams[R, M, E],
) {
	configurator.withStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		createOrUpdateResourceStep(configurator, sCtx, createOrUpdateResourceParams[R, M, E]{
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

func createOrUpdateWorkspaceResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType](
	configurator *StepsConfigurator,
	params createOrUpdateWorkspaceResourceParams[R, M, E],
) {
	configurator.withStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace)

		createOrUpdateResourceStep(configurator, sCtx, createOrUpdateResourceParams[R, M, E]{
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

func createOrUpdateNetworkResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType](
	configurator *StepsConfigurator,
	params createOrUpdateNetworkResourceParams[R, M, E],
) {
	configurator.withStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace, params.network)

		createOrUpdateResourceStep(configurator, sCtx, createOrUpdateResourceParams[R, M, E]{
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

func createOrUpdateResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType](
	configurator *StepsConfigurator,
	sCtx provider.StepCtx,
	params createOrUpdateResourceParams[R, M, E],
) {
	if params.resource != nil {
		fmt.Printf("Request: %v\n", params.resource)
		configurator.suite.ReportRequestStep(sCtx, params.resource)
	}
	resp, err := params.createOrUpdateFunc(configurator.t.Context(), params.resource)

	if resp != nil {
		fmt.Printf("Response: %v\n", resp.resource)
		configurator.suite.ReportResponseStep(sCtx, resp.resource)
	}

	if err != nil {
		requireNoError(sCtx, err)
	}

	if params.expectedLabels != nil {
		configurator.suite.VerifyLabelsStep(sCtx, params.expectedLabels, resp.labels)
	}

	if params.expectedMetadata != nil {
		params.verifyMetadataFunc(sCtx, params.expectedMetadata, resp.metadata)
	}

	if params.expectedSpec != nil {
		params.verifySpecFunc(sCtx, params.expectedSpec, &resp.spec)
	}

	configurator.suite.VerifyStatusStep(sCtx, params.expectedResourceState, *resp.state)
}

func getTenantResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType](
	configurator *StepsConfigurator,
	params getTenantResourceParams[R, M, E],
) *R {
	var resp *stepFuncResponse[R, M, E]
	configurator.withStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		resp = getResourceWithObserver(configurator, sCtx,
			getResourceWithObserverParams[R, M, E, secapi.TenantReference, schema.ResourceState]{
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

func getWorkspaceResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType](
	configurator *StepsConfigurator,
	params getWorkspaceResourceParams[R, M, E],
) *R {
	var resp *stepFuncResponse[R, M, E]
	configurator.withStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.wref.Workspace))

		resp = getResourceWithObserver(configurator, sCtx,
			getResourceWithObserverParams[R, M, E, secapi.WorkspaceReference, schema.ResourceState]{
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

func getNetworkResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType](
	configurator *StepsConfigurator,
	params getNetworkResourceParams[R, M, E],
) *R {
	var resp *stepFuncResponse[R, M, E]
	configurator.withStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.nref.Workspace), string(params.nref.Network))

		resp = getResourceWithObserver(configurator, sCtx,
			getResourceWithObserverParams[R, M, E, secapi.NetworkReference, schema.ResourceState]{
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

func getResourceWithObserver[R types.ResourceType, M types.MetadataType, E types.SpecType, F secapi.Reference, V any](
	configurator *StepsConfigurator,
	sCtx provider.StepCtx,
	params getResourceWithObserverParams[R, M, E, F, V],
) *stepFuncResponse[R, M, E] {
	config := secapi.ResourceObserverConfig[V]{
		ExpectedValue: params.observerExpectedValue,
		Delay:         time.Duration(configurator.suite.BaseDelay) * time.Second,
		Interval:      time.Duration(configurator.suite.BaseInterval) * time.Second,
		MaxAttempts:   configurator.suite.MaxAttempts,
	}

	resp, err := params.getFunc(configurator.t.Context(), params.reference, config)
	requireNoError(sCtx, err)
	requireNotNilResponse(sCtx, resp)

	if params.expectedLabels != nil {
		configurator.suite.VerifyLabelsStep(sCtx, params.expectedLabels, resp.labels)
	}

	if params.expectedMetadata != nil {
		params.verifyMetadataFunc(sCtx, params.expectedMetadata, resp.metadata)
	}

	if params.expectedSpec != nil {
		params.verifySpecFunc(sCtx, params.expectedSpec, &resp.spec)
	}

	configurator.suite.VerifyStatusStep(sCtx, params.expectedResourceState, *resp.state)

	return resp
}
