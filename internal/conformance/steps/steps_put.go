package steps

import (
	"context"
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Params

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

// Steps

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
	requestResourceStep(sCtx, params.resource)
	resp, err := params.createOrUpdateFunc(t.Context(), params.resource)
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
}
