package steps

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StepsConfigurator struct {
	suite *suites.TestSuite
	t     provider.T
	ctx   provider.StepCtx
}

func NewStepsConfigurator(suite *suites.TestSuite, t provider.T) *StepsConfigurator {
	return &StepsConfigurator{
		suite: suite,
		t:     t,
	}
}

// NewStepsConfiguratorWithCtx creates a configurator that runs steps as sub-steps
// of the provided parent step context instead of directly on the test.
func NewStepsConfiguratorWithCtx(suite *suites.TestSuite, t provider.T, ctx provider.StepCtx) *StepsConfigurator {
	return &StepsConfigurator{
		suite: suite,
		t:     t,
		ctx:   ctx,
	}
}

// withStep runs a step either on the root test (t) or as a
// sub-step of an existing StepCtx when used inside a grouped step.
func (configurator *StepsConfigurator) withStep(stepName string, fn func(provider.StepCtx)) {
	if configurator.ctx != nil {
		configurator.ctx.WithNewStep(stepName, fn)
		return
	}

	configurator.t.WithNewStep(stepName, fn)
}
