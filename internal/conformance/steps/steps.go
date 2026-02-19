package steps

import (
	"encoding/json"
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StepResponseExpects[M types.MetadataType, E types.SpecType] struct {
	Labels        schema.Labels
	Metadata      *M
	Spec          *E
	ResourceState schema.ResourceState
}

// Step Responses

type createOrUpdateStepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	resource *R
	labels   schema.Labels
	metadata *M
	spec     E
	status   *S
}

func newCreateOrUpdateStepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	resource *R, labels schema.Labels, metadata *M, spec E, status *S,
) *createOrUpdateStepFuncResponse[R, M, E, S] {
	return &createOrUpdateStepFuncResponse[R, M, E, S]{
		resource: resource,
		labels:   labels,
		metadata: metadata,
		spec:     spec,
		status:   status,
	}
}

type getStepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	resource *R
	labels   schema.Labels
	metadata *M
	spec     E
	status   *S
}

func newGetStepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	resource *R, labels schema.Labels, metadata *M, spec E, status *S,
) *getStepFuncResponse[R, M, E, S] {
	return &getStepFuncResponse[R, M, E, S]{
		resource: resource,
		labels:   labels,
		metadata: metadata,
		spec:     spec,
		status:   status,
	}
}

type globalStepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType] struct {
	resource *R
	metadata *M
	spec     E
}

func newGlobalStepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType](resource *R, metadata *M, spec E) *globalStepFuncResponse[R, M, E] {
	return &globalStepFuncResponse[R, M, E]{
		resource: resource,
		metadata: metadata,
		spec:     spec,
	}
}

// API Request / Response

func requestResourceStep(ctx provider.StepCtx, request any) {
	ctx.WithNewStep("Send request", func(stepCtx provider.StepCtx) {
		if request == nil {
			return
		}

		if data, err := json.Marshal(request); err != nil {
			slog.Error("Error marshaling resource to json", "error", err)
			stepCtx.FailNow()
		} else {
			stepCtx.WithNewParameters("resource", string(data))
		}
	})
}

func emptyRequestStep(ctx provider.StepCtx) {
	ctx.WithNewStep("Send request", func(stepCtx provider.StepCtx) {})
}

func responseResourceStep[R types.ResourceType](ctx provider.StepCtx, resource *R) {
	ctx.WithNewStep("Receive response", func(stepCtx provider.StepCtx) {
		if resource == nil {
			return
		}

		if data, err := json.Marshal(resource); err != nil {
			slog.Error("Error marshaling resource to json", "error", err)
			stepCtx.FailNow()
		} else {
			stepCtx.WithNewParameters("resource", string(data))
		}
	})
}

func responseResourcesStep[R types.ResourceType](ctx provider.StepCtx, resources []*R) {
	ctx.WithNewStep("Receive response", func(stepCtx provider.StepCtx) {
		if resources == nil {
			return
		}

		if data, err := json.Marshal(resources); err != nil {
			slog.Error("Error marshaling resource to json", "error", err)
			stepCtx.FailNow()
		} else {
			stepCtx.WithNewParameters("resources", string(data))
		}
	})
}

func emptyResponseStep(ctx provider.StepCtx) {
	ctx.WithNewStep("Receive response", func(stepCtx provider.StepCtx) {})
}
