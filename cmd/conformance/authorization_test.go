//nolint:dupl
package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/authorization"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
)

func TestAuthorizationV1Suites(t *testing.T) {
	globalTestSuite := suites.CreateGlobalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := &authorization.AuthorizationLifeCycleV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           config.Parameters.ScenariosUsers,
	}
	lifeCycleTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { lifeCycleTestSuite.ScenarioName = constants.AuthorizationV1LifeCycleSuiteName },
	)

	// List Suite
	listTestSuite := &authorization.AuthorizationListV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           config.Parameters.ScenariosUsers,
	}
	listTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { listTestSuite.ScenarioName = constants.AuthorizationV1ListSuiteName },
	)
}
