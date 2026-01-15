//nolint:dupl
package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/usage"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
)

func TestUsageV1Suites(t *testing.T) {
	mixedTestSuite := suites.CreateMixedTestSuite(config.Parameters, config.Clients)

	// Foundation Suite
	foundationTestSuite := &usage.FoundationUsageV1TestSuite{
		MixedTestSuite: mixedTestSuite,
		Users:          config.Parameters.ScenariosUsers,
		NetworkCidr:    config.Parameters.ScenariosCidr,
		PublicIpsRange: config.Parameters.ScenariosPublicIps,
		RegionZones:    config.Clients.RegionZones,
		InstanceSkus:   config.Clients.InstanceSkus,
		StorageSkus:    config.Clients.StorageSkus,
		NetworkSkus:    config.Clients.NetworkSkus,
	}
	foundationTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { foundationTestSuite.ScenarioName = constants.FoundationV1UsageSuiteName },
	)
}
