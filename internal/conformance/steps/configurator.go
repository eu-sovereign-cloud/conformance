package steps

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StepsConfigurator struct {
	suite *suites.TestSuite
	t     provider.T
}

func NewStepsConfigurator(suite *suites.TestSuite, t provider.T) *StepsConfigurator {
	return &StepsConfigurator{
		suite: suite,
		t:     t,
	}
}
