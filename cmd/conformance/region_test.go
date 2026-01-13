package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/region"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestRegionV1Suites(t *testing.T) {
	globalTestSuite := suites.GlobalTestSuite{
		TestSuite: suites.TestSuite{
			AuthToken:     config.Parameters.ClientAuthToken,
			MockEnabled:   config.Parameters.MockEnabled,
			MockServerURL: &config.Parameters.MockServerURL,
			BaseDelay:     config.Parameters.BaseDelay,
			BaseInterval:  config.Parameters.BaseInterval,
			MaxAttempts:   config.Parameters.MaxAttempts,
		},
		Client: config.Clients.GlobalClient,
	}

	// List Suite
	testListSuite := &region.RegionV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		RegionName:      config.Parameters.ClientRegion,
	}
	testListSuite.ScenarioName = constants.RegionV1ListSuiteName
	if testListSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
