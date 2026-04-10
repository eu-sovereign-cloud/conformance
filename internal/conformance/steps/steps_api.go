package steps

import (
	"encoding/json"
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/types"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

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
