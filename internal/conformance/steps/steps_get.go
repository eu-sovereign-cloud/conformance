package steps

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/wrappers"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Params

type getGlobalResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType] struct {
	stepName           string
	stepParamsFunc     func(provider.StepCtx, constants.OperationName)
	operationName      constants.OperationName
	resourceName       string
	getFunc            func(context.Context, string) (wrappers.GlobalResourceWrapper[R, M, E], error)
	expectedMetadata   *M
	verifyMetadataFunc func(provider.StepCtx, *M, *M)
	expectedSpec       *E
	verifySpecFunc     func(provider.StepCtx, *E, *E)
}

type getTenantResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName)
	operationName          constants.OperationName
	tref                   secapi.TenantReference
	getValueFunc           func(context.Context, secapi.TenantReference, secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedResourceStatus S
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
}

type getWorkspaceResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID)
	operationName          constants.OperationName
	wref                   secapi.WorkspaceReference
	getValueFunc           func(context.Context, secapi.WorkspaceReference, secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedResourceStatus S
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
}

type getWorkspaceInstanceResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID)
	operationName          constants.OperationName
	wref                   secapi.WorkspaceReference
	getValueFunc           func(context.Context, secapi.WorkspaceReference, secapi.ResourceObserverUntilValueConfig[schema.InstanceStatusPowerState]) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedResourceStatus S
	expectedPowerState     schema.InstanceStatusPowerState
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
}

type getNetworkResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID, secapi.NetworkID)
	operationName          constants.OperationName
	nref                   secapi.NetworkReference
	getValueFunc           func(context.Context, secapi.NetworkReference, secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedResourceStatus S
	expectedLabels         schema.Labels
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
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
	expectedResourceStatus S
}

// Steps

func getGlobalResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType](
	t provider.T,
	params getGlobalResourceParams[R, M, E],
) *R {
	var err error
	var resp wrappers.GlobalResourceWrapper[R, M, E]
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		emptyRequestStep(sCtx)
		resp, err = params.getFunc(t.Context(), params.resourceName)

		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		resourceResponseStep(sCtx, resp.GetResource())

		// Metadata
		if resp.GetMetadata() != nil && params.expectedMetadata != nil {
			params.verifyMetadataFunc(sCtx, params.expectedMetadata, resp.GetMetadata())
		} else {
			slog.Error("Metadata verification failed: expected or actual metadata is nil")
			t.FailNow()
		}

		if params.expectedSpec != nil {
			params.verifySpecFunc(sCtx, params.expectedSpec, resp.GetSpec())
		}
	})
	return resp.GetResource()
}

func getTenantResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, params getTenantResourceParams[R, M, E, S],
) *R {
	var resp *R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		resp = getResourceUntilValueStep(t, suite, params.stepName, sCtx,
			getResourceUntilValueParams[R, M, E, S, secapi.TenantReference, schema.ResourceState]{
				reference:              params.tref,
				observerExpectedValues: []schema.ResourceState{types.GetStatusState(&params.expectedResourceStatus)},
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

func getWorkspaceResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, params getWorkspaceResourceParams[R, M, E, S],
) *R {
	var resp *R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.wref.Workspace)

		resp = getResourceUntilValueStep(t, suite, params.stepName, sCtx,
			getResourceUntilValueParams[R, M, E, S, secapi.WorkspaceReference, schema.ResourceState]{
				reference:              params.wref,
				observerExpectedValues: []schema.ResourceState{types.GetStatusState(&params.expectedResourceStatus)},
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

func getWorkspaceInstanceResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, params getWorkspaceInstanceResourceParams[R, M, E, S],
) *R {
	var resp *R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.wref.Workspace)

		resp = getResourceUntilValueStep(t, suite, params.stepName, sCtx,
			getResourceUntilValueParams[R, M, E, S, secapi.WorkspaceReference, schema.InstanceStatusPowerState]{
				reference:              params.wref,
				observerExpectedValues: []schema.InstanceStatusPowerState{types.GetStatusPowerState(&params.expectedResourceStatus)},
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

func getNetworkResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	t provider.T, suite *suites.TestSuite, params getNetworkResourceParams[R, M, E, S],
) *R {
	var resp *R
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.nref.Workspace, params.nref.Network)

		resp = getResourceUntilValueStep(t, suite, params.stepName, sCtx,
			getResourceUntilValueParams[R, M, E, S, secapi.NetworkReference, schema.ResourceState]{
				reference:              params.nref,
				observerExpectedValues: []schema.ResourceState{types.GetStatusState(&params.expectedResourceStatus)},
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

func getResourceUntilValueStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType, F secapi.Reference, V any](
	t provider.T, suite *suites.TestSuite, stepName string, sCtx provider.StepCtx, params getResourceUntilValueParams[R, M, E, S, F, V],
) *R {
	slog.Info(fmt.Sprintf("[%s] %s", suite.ScenarioName, stepName))

	config := secapi.ResourceObserverUntilValueConfig[V]{
		ExpectedValues: params.observerExpectedValues,
		Delay:          time.Duration(suite.BaseDelay) * time.Second,
		Interval:       time.Duration(suite.BaseInterval) * time.Second,
		MaxAttempts:    suite.MaxAttempts,
	}
	referenceRequestStep(sCtx, params.reference)

	resp, err := params.getValueFunc(t.Context(), params.reference, config)
	requireNoError(sCtx, err)
	requireNotNilResponse(sCtx, resp)

	resourceResponseStep(sCtx, resp.GetResource())

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
	expectedState := types.GetStatusState(&params.expectedResourceStatus)
	if expectedState != "" {
		suite.VerifyStatusStateStep(sCtx, expectedState, types.GetStatusState(resp.GetStatus()))
	} else {
		log.Fatalln("Status verification failed: expected or actual Status is nil")
	}

	// Conditions
	expectedConditions := types.GetStatusConditions(&params.expectedResourceStatus)
	if len(expectedConditions) > 0 {
		actualConditions := types.GetStatusConditions(resp.GetStatus())
		suite.VerifyStatusConditionsStep(sCtx, expectedConditions, actualConditions)
	}

	// Power state (Instance only)
	expectedPowerState := types.GetStatusPowerState(&params.expectedResourceStatus)
	if expectedPowerState != "" {
		suite.VerifyStatusPowerStateStep(sCtx, expectedPowerState, types.GetStatusPowerState(resp.GetStatus()))
	}

	return resp.GetResource()
}
