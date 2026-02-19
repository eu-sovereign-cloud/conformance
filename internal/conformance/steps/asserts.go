package steps

import (
	"fmt"

	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func requireNoError(sCtx provider.StepCtx, err error) {
	sCtx.WithNewStep("Verify no error", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters("error", fmt.Sprintf("%v", err))
		stepCtx.Require().NoError(err, "Should not return an error")
	})
}

func requireError(sCtx provider.StepCtx, err error, expected error) {
	sCtx.WithNewStep("Verify error", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters("error", fmt.Sprintf("%v", err))
		stepCtx.Require().Error(err, "Should return an error")
		stepCtx.Require().EqualError(err, expected.Error(), "Should return the expected error")
	})
}

func requireNotNilResponse(sCtx provider.StepCtx, resp any) {
	sCtx.WithNewStep("Verify not nil response", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters("response", fmt.Sprintf("%v", resp))
		stepCtx.Require().NotNil(resp, "Should be not nil")
	})
}

func requireNotEmptyResponse[R types.ResourceType](sCtx provider.StepCtx, resp []*R) {
	sCtx.WithNewStep("Verify response length", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters("response", fmt.Sprintf("%v", resp))
		stepCtx.Require().NotNil(resp, "Should be not nil")
		stepCtx.Require().GreaterOrEqual(len(resp), 1, "Should have length greater than 1")
	})
}

func verifyIterListStep[R types.ResourceType](ctx provider.StepCtx, t provider.T, iter secapi.Iterator[R]) {
	ctx.WithNewStep("Verify Iter List", func(stepCtx provider.StepCtx) {
		// Iterate through all items
		resp, err := iter.All(t.Context())
		requireNoError(stepCtx, err)

		requireNotEmptyResponse(stepCtx, resp)
	})
}
