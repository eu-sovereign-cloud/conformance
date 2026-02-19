package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/region"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestRegionV1Suites(t *testing.T) {
	globalTestSuite := suites.CreateGlobalTestSuite(config.Parameters, config.Clients)

	// Provider Queries Suite
	providerQueriesSuite := region.CreateProviderQueriesV1TestSuite(globalTestSuite, config.Parameters.ClientRegion)
	if providerQueriesSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerQueriesSuite)
	}
}
