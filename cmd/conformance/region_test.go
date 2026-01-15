//nolint:dupl
package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/region"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
)

func TestRegionV1Suites(t *testing.T) {
	globalTestSuite := suites.CreateGlobalTestSuite(config.Parameters, config.Clients)

	// List Suite
	listTestSuite := &region.RegionListV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		RegionName:      config.Parameters.ClientRegion,
	}
	listTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { listTestSuite.ScenarioName = constants.RegionV1ListSuiteName },
	)
}
