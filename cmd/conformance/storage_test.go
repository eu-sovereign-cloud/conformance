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

	// Block Storage Constraints Violations Suite
	blockStorageConstraintsSuite := storage.CreateBlockStorageConstraintsValidationV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if blockStorageConstraintsSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, blockStorageConstraintsSuite)
	}

	// Block Storage Error Suite
	blockStorageErrorSuite := storage.CreateBlockStorageErrorV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if blockStorageErrorSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, blockStorageErrorSuite)
	}

	// Image LifeCycle Suite
	imageLifeCycleSuite := storage.CreateImageLifeCycleV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if imageLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, imageLifeCycleSuite)
	}

	// Image Constraints Violations Suite
	imageConstraintsSuite := storage.CreateImageConstraintsValidationV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if imageConstraintsSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, imageConstraintsSuite)
	}
}
