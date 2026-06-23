package suites

import (
	"log/slog"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (suite *TestSuite) verifyAssertState(stepCtx provider.StepCtx) {
	if stepCtx.CurrentStep().Status != passed {
		slog.Error("Verification should have no assertion failures")
		stepCtx.FailNow()
	}
}
