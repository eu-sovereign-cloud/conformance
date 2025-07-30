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

func requireNotNilResponse(sCtx provider.StepCtx, resp interface{}) {
	sCtx.WithNewStep("Verify nil response", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters("response", fmt.Sprintf("%v", resp))
		stepCtx.Require().NotNil(resp, "Should be not nil")
	})
}

func requireStatusEquals(sCtx provider.StepCtx, expectedState string, actualState string) {
	sCtx.WithNewStep("Verify status", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_state", expectedState,
			"actual_state", actualState,
		)
		stepCtx.Require().Equal(expectedState, actualState, fmt.Sprintf("Status.State should be '%s'", expectedState))
	})
}
