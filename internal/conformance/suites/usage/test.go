package usage

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestUsageV1Suite(t *testing.T) {
	TestSuite := &FoundationUsageV1TestSuite{
		MixedTestSuite: suites.MixedTestSuite{
			TestSuite: suites.TestSuite{
				Tenant:        config.Parameters.ClientTenant,
				AuthToken:     config.Parameters.ClientAuthToken,
				MockEnabled:   config.Parameters.MockEnabled,
				MockServerURL: &config.Parameters.MockServerURL,
				ScenarioName:  constants.FoundationV1UsageSuiteName,
				BaseDelay:     config.Parameters.BaseDelay,
				BaseInterval:  config.Parameters.BaseInterval,
				MaxAttempts:   config.Parameters.MaxAttempts,
			},
			GlobalClient:   config.Clients.GlobalClient,
			RegionalClient: config.Clients.RegionalClient,
			Region:         config.Parameters.ClientRegion,
		},
		Users:          config.Parameters.ScenariosUsers,
		NetworkCidr:    config.Parameters.ScenariosCidr,
		PublicIpsRange: config.Parameters.ScenariosPublicIps,
		RegionZones:    config.Clients.RegionZones,
		InstanceSkus:   config.Clients.InstanceSkus,
		StorageSkus:    config.Clients.StorageSkus,
		NetworkSkus:    config.Clients.NetworkSkus,
	}

	if TestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, TestSuite)
	}
}
