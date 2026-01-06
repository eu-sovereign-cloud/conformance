package region

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestRegionV1Suite(t *testing.T) {
	TestSuite := &RegionV1TestSuite{
		GlobalTestSuite: suites.GlobalTestSuite{
			TestSuite: suites.TestSuite{
				AuthToken:     config.Parameters.ClientAuthToken,
				MockEnabled:   config.Parameters.MockEnabled,
				MockServerURL: &config.Parameters.MockServerURL,
				ScenarioName:  constants.RegionV1ListSuiteName,
			},
			Client: config.Clients.GlobalClient,
		},
		RegionName: config.Parameters.ClientRegion,
	}

	if TestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, TestSuite)
	}
}
