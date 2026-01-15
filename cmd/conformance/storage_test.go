package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/storage"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestStorageV1Suites(t *testing.T) {
	regionalTestSuite := createRegionalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := &storage.StorageLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       config.Clients.StorageSkus,
	}
	lifeCycleTestSuite.ScenarioName = constants.StorageV1LifeCycleSuiteName
	if lifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, lifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := &storage.StorageListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       config.Clients.StorageSkus,
	}
	listTestSuite.ScenarioName = constants.StorageV1ListSuiteName
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}
}
