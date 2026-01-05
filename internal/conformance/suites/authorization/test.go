package authorization

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestAuthorizationV1Suite(t *testing.T) {
	globalTestSuite := suites.GlobalTestSuite{
		TestSuite: suites.TestSuite{
			Tenant:        conformance.Config.ClientTenant,
			AuthToken:     conformance.Config.ClientAuthToken,
			MockEnabled:   conformance.Config.MockEnabled,
			MockServerURL: &conformance.Config.MockServerURL,
			ScenarioName:  conformance.AuthorizationV1LifeCycleSuiteName,
			BaseDelay:     conformance.Config.BaseDelay,
			BaseInterval:  conformance.Config.BaseInterval,
			MaxAttempts:   conformance.Config.MaxAttempts,
		},
		Client: conformance.Clients.GlobalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &AuthorizationV1LifeCycleTestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           conformance.Config.ScenariosUsers,
	}

	if testLifeCycleSuite.CanRun(conformance.Config.ScenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &AuthorizationV1ListTestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           conformance.Config.ScenariosUsers,
	}

	testListSuite.ScenarioName = conformance.AuthorizationV1ListSuiteName

	if testListSuite.CanRun(conformance.Config.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
