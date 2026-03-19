package steps

import (
	"context"
	"log"
	"time"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/conformance/pkg/wrappers"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ResponseExpects[M types.MetadataType, E types.SpecType] struct {
	Labels         schema.Labels
	Metadata       *M
	Spec           *E
	ResourceStates []schema.ResourceState
}

type ResponseExpectsWithCondition[M types.MetadataType, E types.SpecType] struct {
	Labels         schema.Labels
	Metadata       *M
	Spec           *E
	ResourceStatus schema.Status
}

// Steps

/// Create Or Update

type createOrUpdateTenantResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName)
	operationName          constants.OperationName
	resource               *R
	createOrUpdateFunc     func(context.Context, *R) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
	expectedResourceStates []schema.ResourceState
}

func createOrUpdateTenantResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, params createOrUpdateTenantResourceParams[R, M, E, S],
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		createOrUpdateResourceStep(t, suite, sCtx, createOrUpdateResourceParams[R, M, E, S]{
			resource:               params.resource,
			createOrUpdateFunc:     params.createOrUpdateFunc,
			expectedLabels:         params.expectedLabels,
			expectedMetadata:       params.expectedMetadata,
			verifyMetadataFunc:     params.verifyMetadataFunc,
			expectedSpec:           params.expectedSpec,
			verifySpecFunc:         params.verifySpecFunc,
			expectedResourceStates: params.expectedResourceStates,
		})
	})
}

type createOrUpdateWorkspaceResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName, string)
	operationName          constants.OperationName
	workspace              string
	resource               *R
	createOrUpdateFunc     func(context.Context, *R) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
	expectedResourceStates []schema.ResourceState
}

func createOrUpdateWorkspaceResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, params createOrUpdateWorkspaceResourceParams[R, M, E, S],
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace)

		createOrUpdateResourceStep(t, suite, sCtx, createOrUpdateResourceParams[R, M, E, S]{
			resource:               params.resource,
			createOrUpdateFunc:     params.createOrUpdateFunc,
			expectedLabels:         params.expectedLabels,
			expectedMetadata:       params.expectedMetadata,
			verifyMetadataFunc:     params.verifyMetadataFunc,
			expectedSpec:           params.expectedSpec,
			verifySpecFunc:         params.verifySpecFunc,
			expectedResourceStates: params.expectedResourceStates,
		})
	})
}

type createOrUpdateNetworkResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName, string, string)
	operationName          constants.OperationName
	workspace              string
	network                string
	resource               *R
	createOrUpdateFunc     func(context.Context, *R) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
	expectedResourceStates []schema.ResourceState
}

func createOrUpdateNetworkResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, params createOrUpdateNetworkResourceParams[R, M, E, S],
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace, params.network)

		createOrUpdateResourceStep(t, suite, sCtx, createOrUpdateResourceParams[R, M, E, S]{
			resource:               params.resource,
			createOrUpdateFunc:     params.createOrUpdateFunc,
			expectedLabels:         params.expectedLabels,
			expectedMetadata:       params.expectedMetadata,
			verifyMetadataFunc:     params.verifyMetadataFunc,
			expectedSpec:           params.expectedSpec,
			verifySpecFunc:         params.verifySpecFunc,
			expectedResourceStates: params.expectedResourceStates,
		})
	})
}

type createOrUpdateResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	resource               *R
	createOrUpdateFunc     func(context.Context, *R) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
	expectedResourceStates []schema.ResourceState
}

func createOrUpdateResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, sCtx provider.StepCtx, params createOrUpdateResourceParams[R, M, E, S],
) {
	resp, err := params.createOrUpdateFunc(t.Context(), params.resource)
	requireNoError(sCtx, err)

	// Label
	if params.expectedLabels != nil {
		suite.VerifyLabelsStep(sCtx, params.expectedLabels, resp.GetLabels())
	}

	// Metadata
	if resp.GetMetadata() != nil && params.expectedMetadata != nil {
		params.verifyMetadataFunc(sCtx, params.expectedMetadata, resp.GetMetadata())
	} else {
		log.Fatalln("Metadata verification failed: expected or actual metadata is nil")
	}

	if params.expectedSpec != nil {
		params.verifySpecFunc(sCtx, params.expectedSpec, resp.GetSpec())
	}

	// Status
	if resp.GetStatus() != nil && len(params.expectedResourceStates) > 0 {
		suite.VerifyStatusStatesStep(sCtx, params.expectedResourceStates, types.GetStatusState(resp.GetStatus()))
	} else {
		log.Fatalln("Status verification failed: expected or actual Status is nil")
	}
}

/// Get

type getTenantResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName)
	operationName          constants.OperationName
	tref                   secapi.TenantReference
	getValueFunc           func(context.Context, secapi.TenantReference, secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedResourceStates []schema.ResourceState
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
}
type getTenantResourceParamsWithConditions[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName)
	operationName          constants.OperationName
	tref                   secapi.TenantReference
	getValueFunc           func(context.Context, secapi.TenantReference, secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedResourceStatus schema.Status
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
}

type getTenantResourceWithConditionParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName)
	operationName          constants.OperationName
	tref                   secapi.TenantReference
	getValueFunc           func(context.Context, secapi.TenantReference, secapi.ResourceObserverUntilValueConfig[schema.Status]) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedResourceStatus []schema.Status
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
}

func getTenantResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, params getTenantResourceParams[R, M, E, S],
) *R {
	var resp *R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		resp = getResourceUntilValueStep(t, suite, sCtx,
			getResourceUntilValueParams[R, M, E, S, secapi.TenantReference, schema.ResourceState]{
				reference:              params.tref,
				observerExpectedValues: params.expectedResourceStates,
				getValueFunc:           params.getValueFunc,
				expectedLabels:         params.expectedLabels,
				expectedMetadata:       params.expectedMetadata,
				verifyMetadataFunc:     params.verifyMetadataFunc,
				expectedSpec:           params.expectedSpec,
				verifySpecFunc:         params.verifySpecFunc,
				expectedResourceStates: params.expectedResourceStates,
			},
		)
		requireNotNilResponse(sCtx, resp)
	})
	return resp
}

func getTenantResourceWithConditionStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, params getTenantResourceParamsWithConditions[R, M, E, S],
) *R {
	var resp *R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		resp = getResourceUntilValueWithConditionStep(t, suite, sCtx,
			getResourceUntilValueConditionParams[R, M, E, S, secapi.TenantReference, schema.ResourceState]{
				reference:              params.tref,
				observerExpectedValues: []schema.ResourceState{*params.expectedResourceStatus.State},
				getValueFunc:           params.getValueFunc,
				expectedLabels:         params.expectedLabels,
				expectedMetadata:       params.expectedMetadata,
				verifyMetadataFunc:     params.verifyMetadataFunc,
				expectedSpec:           params.expectedSpec,
				verifySpecFunc:         params.verifySpecFunc,
				expectedResourceStatus: params.expectedResourceStatus,
			},
		)
		requireNotNilResponse(sCtx, resp)
	})
	return resp
}

type getWorkspaceResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName, string)
	operationName          constants.OperationName
	wref                   secapi.WorkspaceReference
	getValueFunc           func(context.Context, secapi.WorkspaceReference, secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedResourceStates []schema.ResourceState
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
}

func getWorkspaceResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, params getWorkspaceResourceParams[R, M, E, S],
) *R {
	var resp *R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.wref.Workspace))

		resp = getResourceUntilValueStep(t, suite, sCtx,
			getResourceUntilValueParams[R, M, E, S, secapi.WorkspaceReference, schema.ResourceState]{
				reference:              params.wref,
				observerExpectedValues: params.expectedResourceStates,
				getValueFunc:           params.getValueFunc,
				expectedLabels:         params.expectedLabels,
				expectedMetadata:       params.expectedMetadata,
				verifyMetadataFunc:     params.verifyMetadataFunc,
				expectedSpec:           params.expectedSpec,
				verifySpecFunc:         params.verifySpecFunc,
				expectedResourceStates: params.expectedResourceStates,
			},
		)
		requireNotNilResponse(sCtx, resp)
	})
	return resp
}

type getNetworkResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName, string, string)
	operationName          constants.OperationName
	nref                   secapi.NetworkReference
	getValueFunc           func(context.Context, secapi.NetworkReference, secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedResourceStates []schema.ResourceState
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
}

func getNetworkResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, params getNetworkResourceParams[R, M, E, S],
) *R {
	var resp *R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.nref.Workspace), string(params.nref.Network))

		resp = getResourceUntilValueStep(t, suite, sCtx,
			getResourceUntilValueParams[R, M, E, S, secapi.NetworkReference, schema.ResourceState]{
				reference:              params.nref,
				observerExpectedValues: params.expectedResourceStates,
				getValueFunc:           params.getValueFunc,
				expectedLabels:         params.expectedLabels,
				expectedMetadata:       params.expectedMetadata,
				verifyMetadataFunc:     params.verifyMetadataFunc,
				expectedSpec:           params.expectedSpec,
				verifySpecFunc:         params.verifySpecFunc,
				expectedResourceStates: params.expectedResourceStates,
			},
		)
		requireNotNilResponse(sCtx, resp)
	})
	return resp
}

type getResourceUntilValueParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType, F secapi.Reference, V any] struct {
	reference              F
	observerExpectedValues []V
	getValueFunc           func(context.Context, F, secapi.ResourceObserverUntilValueConfig[V]) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
	expectedResourceStates []schema.ResourceState
}

type getResourceUntilValueConditionParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType, F secapi.Reference, V any] struct {
	reference              F
	observerExpectedValues []V
	getValueFunc           func(context.Context, F, secapi.ResourceObserverUntilValueConfig[V]) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
	expectedResourceStatus schema.Status
}

func getResourceUntilValueStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType, F secapi.Reference, V any](
	t provider.T, suite *suites.TestSuite, sCtx provider.StepCtx, params getResourceUntilValueParams[R, M, E, S, F, V],
) *R {
	config := secapi.ResourceObserverUntilValueConfig[V]{
		ExpectedValues: params.observerExpectedValues,
		Delay:          time.Duration(suite.BaseDelay) * time.Second,
		Interval:       time.Duration(suite.BaseInterval) * time.Second,
		MaxAttempts:    suite.MaxAttempts,
	}

	resp, err := params.getValueFunc(t.Context(), params.reference, config)
	requireNoError(sCtx, err)
	requireNotNilResponse(sCtx, resp)

	// Label
	if params.expectedLabels != nil {
		suite.VerifyLabelsStep(sCtx, params.expectedLabels, resp.GetLabels())
	}

	// Metadata
	if params.expectedMetadata != nil {
		params.verifyMetadataFunc(sCtx, params.expectedMetadata, resp.GetMetadata())
	} else {
		log.Fatalln("Metadata verification failed: expected or actual metadata is nil")
	}

	if params.expectedSpec != nil {
		params.verifySpecFunc(sCtx, params.expectedSpec, resp.GetSpec())
	}

	// Status
	if params.expectedResourceStates != nil {
		suite.VerifyStatusStatesStep(sCtx, params.expectedResourceStates, types.GetStatusState(resp.GetStatus()))
	} else {
		log.Fatalln("Status verification failed: expected or actual Status is nil")
	}

	return resp.GetResource()
}

func getResourceUntilValueWithConditionStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType, F secapi.Reference, V any](
	t provider.T, suite *suites.TestSuite, sCtx provider.StepCtx, params getResourceUntilValueConditionParams[R, M, E, S, F, V],
) *R {
	config := secapi.ResourceObserverUntilValueConfig[V]{
		ExpectedValues: params.observerExpectedValues,
		Delay:          time.Duration(suite.BaseDelay) * time.Second,
		Interval:       time.Duration(suite.BaseInterval) * time.Second,
		MaxAttempts:    suite.MaxAttempts,
	}

	resp, err := params.getValueFunc(t.Context(), params.reference, config)
	requireNoError(sCtx, err)
	requireNotNilResponse(sCtx, resp)

	// Label
	if params.expectedLabels != nil {
		suite.VerifyLabelsStep(sCtx, params.expectedLabels, resp.GetLabels())
	}

	// Metadata
	if params.expectedMetadata != nil {
		params.verifyMetadataFunc(sCtx, params.expectedMetadata, resp.GetMetadata())
	} else {
		log.Fatalln("Metadata verification failed: expected or actual metadata is nil")
	}

	if params.expectedSpec != nil {
		params.verifySpecFunc(sCtx, params.expectedSpec, resp.GetSpec())
	}

	// Status
	if params.expectedResourceStatus.State != nil {
		suite.VerifyStatusStateStep(sCtx, params.expectedResourceStatus.State, types.GetStatusState(resp.GetStatus()))
	} else {
		log.Fatalln("Status verification failed: expected or actual Status is nil")
	}

	// Conditions
	if len(params.expectedResourceStatus.Conditions) > 0 {
		actualConditions := types.GetStatusConditions(resp.GetStatus())
		suite.VerifyStatusConditionsStep(sCtx, params.expectedResourceStatus.Conditions, actualConditions)
	}

	return resp.GetResource()
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
	stepParamsFunc func(provider.StepCtx, constants.OperationName, string)
	operationName  constants.OperationName
}

func watchWorkspaceResourceUntilDeletedStep(
	t provider.T, suite *suites.TestSuite, params watchWorkspaceResourceUntilDeletedParams,
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.reference.Workspace))
		watchResourceUntilDeletedStep(t, suite, sCtx, params.watchResourceUntilDeletedParams)
	})
}

type watchNetworkResourceUntilDeletedParams struct {
	watchResourceUntilDeletedParams[secapi.NetworkReference]
	stepName       string
	stepParamsFunc func(provider.StepCtx, constants.OperationName, string, string)
	operationName  constants.OperationName
}

func watchNetworkResourceUntilDeletedStep(
	t provider.T, suite *suites.TestSuite, params watchNetworkResourceUntilDeletedParams,
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.reference.Workspace), string(params.reference.Network))
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
