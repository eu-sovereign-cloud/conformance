package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/region"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestRegionV1Suites(t *testing.T) {
	globalTestSuite := createGlobalTestSuite(config.Parameters, config.Clients)

	// List Suite
	listTestSuite := &region.RegionListV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		RegionName:      config.Parameters.ClientRegion,
	}
	listTestSuite.ScenarioName = constants.RegionV1ListSuiteName
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}
}
