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

func responseResourceStep[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](ctx provider.StepCtx, response *stepFuncResponse[R, M, E, S]) {
	ctx.WithNewStep("Receive response", func(stepCtx provider.StepCtx) {
		if response == nil || response.resource == nil {
			return
		}

		if data, err := json.Marshal(response.resource); err != nil {
			slog.Error("Error marshaling resource to json", "error", err)
			stepCtx.FailNow()
		} else {
			stepCtx.WithNewParameters("resource", string(data))
		}
	})
}

func emptyResponseStep(ctx provider.StepCtx) {
	ctx.WithNewStep("Receive response", func(stepCtx provider.StepCtx) {
	})
}
