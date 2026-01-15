//nolint:dupl
package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/storage"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
)

func TestStorageV1Suites(t *testing.T) {
	regionalTestSuite := suites.CreateRegionalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := &storage.StorageLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       config.Clients.StorageSkus,
	}
	lifeCycleTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { lifeCycleTestSuite.ScenarioName = constants.StorageV1LifeCycleSuiteName },
	)

	// List Suite
	listTestSuite := &storage.StorageListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       config.Clients.StorageSkus,
	}
	listTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { listTestSuite.ScenarioName = constants.StorageV1ListSuiteName },
	)
}
