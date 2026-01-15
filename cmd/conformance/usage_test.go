package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/usage"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestUsageV1Suites(t *testing.T) {
	mixedTestSuite := suites.MixedTestSuite{
		TestSuite: suites.TestSuite{
			Tenant:        config.Parameters.ClientTenant,
			AuthToken:     config.Parameters.ClientAuthToken,
			MockEnabled:   config.Parameters.MockEnabled,
			MockServerURL: &config.Parameters.MockServerURL,
			BaseDelay:     config.Parameters.BaseDelay,
			BaseInterval:  config.Parameters.BaseInterval,
			MaxAttempts:   config.Parameters.MaxAttempts,
		},
		GlobalClient:   config.Clients.GlobalClient,
		RegionalClient: config.Clients.RegionalClient,
		Region:         config.Parameters.ClientRegion,
	}

	// Foundation Suite
	testFoundationSuite := &usage.FoundationUsageV1TestSuite{
		MixedTestSuite: mixedTestSuite,
		Users:          config.Parameters.ScenariosUsers,
		NetworkCidr:    config.Parameters.ScenariosCidr,
		PublicIpsRange: config.Parameters.ScenariosPublicIps,
		RegionZones:    config.Clients.RegionZones,
		InstanceSkus:   config.Clients.InstanceSkus,
		StorageSkus:    config.Clients.StorageSkus,
		NetworkSkus:    config.Clients.NetworkSkus,
	}
	testFoundationSuite.ScenarioName = constants.FoundationV1UsageSuiteName
	if testFoundationSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testFoundationSuite)
	}
}
