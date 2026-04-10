package steps

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/wrappers"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Params

type createOrUpdateTenantResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName)
	operationName          constants.OperationName
	resource               *R
	createOrUpdateFunc     func(context.Context, *R) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedLabels         schema.Labels
	expectedAnnotations    schema.Annotations
	expectedExtensions     schema.Extensions
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
	expectedResourceStates []schema.ResourceState
}

type createOrUpdateWorkspaceResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID)
	operationName          constants.OperationName
	workspace              secapi.WorkspaceID
	resource               *R
	createOrUpdateFunc     func(context.Context, *R) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedLabels         schema.Labels
	expectedAnnotations    schema.Annotations
	expectedExtensions     schema.Extensions
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
	expectedResourceStates []schema.ResourceState
}

type createOrUpdateNetworkResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	stepName               string
	stepParamsFunc         func(provider.StepCtx, constants.OperationName, secapi.WorkspaceID, secapi.NetworkID)
	operationName          constants.OperationName
	workspace              secapi.WorkspaceID
	network                secapi.NetworkID
	resource               *R
	createOrUpdateFunc     func(context.Context, *R) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedLabels         schema.Labels
	expectedAnnotations    schema.Annotations
	expectedExtensions     schema.Extensions
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
	expectedResourceStates []schema.ResourceState
}

type createOrUpdateResourceParams[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	resource               *R
	createOrUpdateFunc     func(context.Context, *R) (wrappers.ResourceWrapper[R, M, E, S], error)
	expectedLabels         schema.Labels
	expectedAnnotations    schema.Annotations
	expectedExtensions     schema.Extensions
	expectedMetadata       *M
	verifyMetadataFunc     func(provider.StepCtx, *M, *M)
	expectedSpec           *E
	verifySpecFunc         func(provider.StepCtx, *E, *E)
	expectedResourceStates []schema.ResourceState
}

// Steps

func createOrUpdateTenantResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	ctx context.Context, suite *suites.TestSuite, stepCreator StepCreator, params createOrUpdateTenantResourceParams[R, M, E, S],
) {
	stepCreator.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		createOrUpdateResourceStep(ctx, suite, params.stepName, sCtx, createOrUpdateResourceParams[R, M, E, S]{
			resource:               params.resource,
			createOrUpdateFunc:     params.createOrUpdateFunc,
			expectedLabels:         params.expectedLabels,
			expectedAnnotations:    params.expectedAnnotations,
			expectedExtensions:     params.expectedExtensions,
			expectedMetadata:       params.expectedMetadata,
			verifyMetadataFunc:     params.verifyMetadataFunc,
			expectedSpec:           params.expectedSpec,
			verifySpecFunc:         params.verifySpecFunc,
			expectedResourceStates: params.expectedResourceStates,
		})
	})
}

func createOrUpdateWorkspaceResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	ctx context.Context, suite *suites.TestSuite, stepCreator StepCreator, params createOrUpdateWorkspaceResourceParams[R, M, E, S],
) {
	stepCreator.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace)

		createOrUpdateResourceStep(ctx, suite, params.stepName, sCtx, createOrUpdateResourceParams[R, M, E, S]{
			resource:               params.resource,
			createOrUpdateFunc:     params.createOrUpdateFunc,
			expectedLabels:         params.expectedLabels,
			expectedAnnotations:    params.expectedAnnotations,
			expectedExtensions:     params.expectedExtensions,
			expectedMetadata:       params.expectedMetadata,
			verifyMetadataFunc:     params.verifyMetadataFunc,
			expectedSpec:           params.expectedSpec,
			verifySpecFunc:         params.verifySpecFunc,
			expectedResourceStates: params.expectedResourceStates,
		})
	})
}

func createOrUpdateNetworkResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	ctx context.Context, suite *suites.TestSuite, stepCreator StepCreator, params createOrUpdateNetworkResourceParams[R, M, E, S],
) {
	stepCreator.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace, params.network)

		createOrUpdateResourceStep(ctx, suite, params.stepName, sCtx, createOrUpdateResourceParams[R, M, E, S]{
			resource:               params.resource,
			createOrUpdateFunc:     params.createOrUpdateFunc,
			expectedLabels:         params.expectedLabels,
			expectedAnnotations:    params.expectedAnnotations,
			expectedExtensions:     params.expectedExtensions,
			expectedMetadata:       params.expectedMetadata,
			verifyMetadataFunc:     params.verifyMetadataFunc,
			expectedSpec:           params.expectedSpec,
			verifySpecFunc:         params.verifySpecFunc,
			expectedResourceStates: params.expectedResourceStates,
		})
	})
}

func createOrUpdateResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	ctx context.Context, suite *suites.TestSuite, stepName string, sCtx provider.StepCtx, params createOrUpdateResourceParams[R, M, E, S],
) {
	slog.Info(fmt.Sprintf("[%s] %s", suite.ScenarioName, stepName))

	requestResourceStep(sCtx, params.resource)

	resp, err := params.createOrUpdateFunc(ctx, params.resource)
	requireNoError(sCtx, err)
	requireNotNilResponse(sCtx, resp)

	responseResourceStep(sCtx, resp.GetResource())

	// Labels
	if params.expectedLabels != nil {
		suite.VerifyLabelsStep(sCtx, params.expectedLabels, resp.GetLabels())
	}

	// Extensions
	if params.expectedExtensions != nil {
		suite.VerifyExtensionsStep(sCtx, params.expectedExtensions, resp.GetExtensions())
	}

	// Annotations
	if params.expectedAnnotations != nil {
		suite.VerifyAnnotationsStep(sCtx, params.expectedAnnotations, resp.GetAnnotations())
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
