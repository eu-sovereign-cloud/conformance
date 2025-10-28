package secatest

import (
	"fmt"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type stepRetry struct {
	baseDelay    time.Duration
	baseInterval time.Duration
	maxAttempts  int

	actFunc    func() schema.ResourceState
	assertFunc func()
}

func newStepRetry(
	baseDelaySecs int,
	baseIntervalSecs int,
	maxAttempts int,
	actFunc func() schema.ResourceState,
	assertFunc func(),
) *stepRetry {
	return &stepRetry{
		baseDelay:    time.Duration(baseDelaySecs) * time.Second,
		baseInterval: time.Duration(baseIntervalSecs) * time.Second,
		maxAttempts:  maxAttempts,
		actFunc:      actFunc,
		assertFunc:   assertFunc,
	}
}

func (retry *stepRetry) run(ctx provider.StepCtx, operation string, expectedResourceState string) {
	time.Sleep(retry.baseDelay)

	for attempt := 1; attempt <= retry.maxAttempts; attempt++ {
		state := retry.actFunc()
		if state == *secalib.SetResourceState(expectedResourceState) {
			retry.assertFunc()
			return
		}

		if attempt >= retry.maxAttempts {
			ctx.WithNewStep("Max attempts reached", func(stepCtx provider.StepCtx) {
				stepCtx.Error(ctx,
					fmt.Sprintf("%s did not reach expected state '%s' after %d attempts ", operation, expectedResourceState, retry.maxAttempts),
				)
			})
		} else {
			time.Sleep(retry.baseInterval)
		}
	}
}
