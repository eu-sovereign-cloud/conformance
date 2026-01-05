package usage

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestUsageV1Suite(t *testing.T) {
	TestSuite := &UsageV1TestSuite{
		MixedTestSuite: suites.MixedTestSuite{
			TestSuite: suites.TestSuite{
				Tenant:        conformance.Config.ClientTenant,
				AuthToken:     conformance.Config.ClientAuthToken,
				MockEnabled:   conformance.Config.MockEnabled,
				MockServerURL: &conformance.Config.MockServerURL,
				ScenarioName:  conformance.FoundationV1UsageSuiteName,
				BaseDelay:     conformance.Config.BaseDelay,
				BaseInterval:  conformance.Config.BaseInterval,
				MaxAttempts:   conformance.Config.MaxAttempts,
			},
			GlobalClient:   conformance.Clients.GlobalClient,
			RegionalClient: conformance.Clients.RegionalClient,
			Region:         conformance.Config.ClientRegion,
		},
		Users:          conformance.Config.ScenariosUsers,
		NetworkCidr:    conformance.Config.ScenariosCidr,
		PublicIpsRange: conformance.Config.ScenariosPublicIps,
		RegionZones:    conformance.Clients.RegionZones,
		InstanceSkus:   conformance.Clients.InstanceSkus,
		StorageSkus:    conformance.Clients.StorageSkus,
		NetworkSkus:    conformance.Clients.NetworkSkus,
	}

	if TestSuite.CanRun(conformance.Config.ScenariosRegexp) {
		suite.RunSuite(t, TestSuite)
	}
}
