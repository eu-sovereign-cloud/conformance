package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestRegionV1Suite(t *testing.T) {
	testSuite := &RegionV1TestSuite{
		globalTestSuite: globalTestSuite{
			testSuite: testSuite{
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
				scenarioName:  regionV1LifeCycleSuiteName,
			},
			client: clients.globalClient,
		},
	}

	if testSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testSuite)
	}
}
