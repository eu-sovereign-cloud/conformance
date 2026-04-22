package steps

import (
	"encoding/json"
	"log/slog"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func resourceRequestStep[R types.ResourceType](ctx provider.StepCtx, request *R) {
	ctx.WithNewStep("Send request", func(stepCtx provider.StepCtx) {
		if data, err := json.Marshal(request); err != nil {
			slog.Error("Error marshaling resource to json", "error", err)
			stepCtx.FailNow()
		} else {
			stepCtx.WithNewParameters("resource", string(data))
		}
	})
}

func referenceRequestStep[R secapi.ReferenceType](ctx provider.StepCtx, reference R) {
	ctx.WithNewStep("Send request", func(stepCtx provider.StepCtx) {
		if data, err := json.Marshal(reference); err != nil {
			slog.Error("Error marshaling reference to json", "error", err)
			stepCtx.FailNow()
		} else {
			stepCtx.WithNewParameters("reference", string(data))
		}
	})
}

func pathRequestStep[P secapi.PathType](ctx provider.StepCtx, path P, options *secapi.ListOptions) {
	ctx.WithNewStep("Send request", func(stepCtx provider.StepCtx) {
		if pathData, err := json.Marshal(path); err != nil {
			slog.Error("Error marshaling path to json", "error", err)
			stepCtx.FailNow()
		} else {
			stepCtx.WithNewParameters("path", string(pathData))
		}

		if options != nil {
			if optionsData, err := json.Marshal(options); err != nil {
				slog.Error("Error marshaling options to json", "error", err)
				stepCtx.FailNow()
			} else {
				stepCtx.WithNewParameters("options", string(optionsData))
			}
		}
	})
}

func emptyRequestStep(ctx provider.StepCtx) {
	ctx.WithNewStep("Send request", func(stepCtx provider.StepCtx) {})
}

func resourceResponseStep[R types.ResourceType](ctx provider.StepCtx, resource *R) {
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

func iteratorResponseStep[R types.ResourceType](ctx provider.StepCtx, resources []*R) {
	ctx.WithNewStep("Receive response", func(stepCtx provider.StepCtx) {
		if resources == nil {
			return
		}

		if data, err := json.Marshal(resources); err != nil {
			slog.Error("Error marshaling iterator to json", "error", err)
			stepCtx.FailNow()
		} else {
			stepCtx.WithNewParameters("iterator", string(data))
		}
	})
}

func emptyResponseStep(ctx provider.StepCtx) {
	ctx.WithNewStep("Receive response", func(stepCtx provider.StepCtx) {})
}
