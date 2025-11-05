package secatest

import (
	"fmt"
	"time"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type stepResourceStateRetry struct {
	baseDelay    time.Duration
	baseInterval time.Duration
	maxAttempts  int

	actFunc    func() (schema.ResourceState, error)
	assertFunc func()
}

func newStepResourceStateRetry(
	baseDelaySecs int,
	baseIntervalSecs int,
	maxAttempts int,

	actFunc func() (schema.ResourceState, error),
	assertFunc func(),
) *stepResourceStateRetry {
	return &stepResourceStateRetry{
		baseDelay:    time.Duration(baseDelaySecs) * time.Second,
		baseInterval: time.Duration(baseIntervalSecs) * time.Second,
		maxAttempts:  maxAttempts,

		actFunc:    actFunc,
		assertFunc: assertFunc,
	}
}

func (retry *stepResourceStateRetry) run(ctx provider.StepCtx, operation string, expectedResourceState string) {
	observer := secapi.NewResourceStateObserver(
		retry.baseDelay,
		retry.baseInterval,
		retry.maxAttempts,
		retry.actFunc,
	)

	err := observer.WaitUntil(schema.ResourceState(expectedResourceState))

	if err == secapi.ErrRetryMaxAttemptsReached {
		ctx.WithNewStep("Max attempts reached", func(stepCtx provider.StepCtx) {
			stepCtx.Error(ctx,
				fmt.Sprintf("%s did not reach expected state '%s' after %d attempts ", operation, expectedResourceState, retry.maxAttempts),
			)
		})
	}

	retry.assertFunc()
}
