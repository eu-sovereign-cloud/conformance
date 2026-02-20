package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/storage"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestStorageV1Suites(t *testing.T) {
	regionalTestSuite := suites.CreateRegionalTestSuite(config.Parameters, config.Clients)

	// Provider LifeCycle Suite
	providerLifeCycleSuite := storage.CreateProviderLifeCycleV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if providerLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerLifeCycleSuite)
	}

	// Provider Queries Suite
	providerQueriesSuite := storage.CreateProviderQueriesV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if providerQueriesSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerQueriesSuite)
	}

	// Block Strage LifeCycle Suite
	blockStorageLifeCycleSuite := storage.CreateBlockStorageLifeCycleV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if blockStorageLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, blockStorageLifeCycleSuite)
	}

	// Image LifeCycle Suite
	imageLifeCycleSuite := storage.CreateImageLifeCycleV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if imageLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, imageLifeCycleSuite)
	}
}
