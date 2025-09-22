package secatest

import (
	"fmt"

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
