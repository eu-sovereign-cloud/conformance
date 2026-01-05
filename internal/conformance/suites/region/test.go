package region

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestRegionV1Suite(t *testing.T) {
	TestSuite := &RegionV1TestSuite{
		GlobalTestSuite: suites.GlobalTestSuite{
			TestSuite: suites.TestSuite{
				AuthToken:     conformance.Config.ClientAuthToken,
				MockEnabled:   conformance.Config.MockEnabled,
				MockServerURL: &conformance.Config.MockServerURL,
				ScenarioName:  conformance.RegionV1LifeCycleSuiteName,
			},
			Client: conformance.Clients.GlobalClient,
		},
		RegionName: conformance.Config.ClientRegion,
	}

	if TestSuite.CanRun(conformance.Config.ScenariosRegexp) {
		suite.RunSuite(t, TestSuite)
	}
}
