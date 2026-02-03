package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/storage"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestStorageV1Suites(t *testing.T) {
	regionalTestSuite := suites.NewRegionalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := storage.NewLifeCycleV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if lifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, lifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := storage.NewListV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}

	// Create Block Storage  Suite
	createBlockStorageTestSuite := storage.NewCreateBlockStorageV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if createBlockStorageTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, createBlockStorageTestSuite)
	}

	// Update Block Storage Suite
	updateBlockStorageTestSuite := storage.NewUpdateBlockStorageV1TestSuite(regionalTestSuite, config.Clients.StorageSkus)
	if updateBlockStorageTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, updateBlockStorageTestSuite)
	}
}
