package steps

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type Builder struct {
	suite *suites.TestSuite
	t     provider.T
}

func NewBuilder(suite *suites.TestSuite, t provider.T) *Builder {
	return &Builder{
		suite: suite,
		t:     t,
	}
}
