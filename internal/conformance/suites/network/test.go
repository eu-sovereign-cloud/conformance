package network

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestNetworkLifecycleV1Suite(t *testing.T) {
	regionalTest := suites.RegionalTestSuite{
		TestSuite: suites.TestSuite{
			Tenant:        conformance.Config.ClientTenant,
			AuthToken:     conformance.Config.ClientAuthToken,
			MockEnabled:   conformance.Config.MockEnabled,
			MockServerURL: &conformance.Config.MockServerURL,
			ScenarioName:  conformance.NetworkV1LifeCycleSuiteName,
			BaseDelay:     conformance.Config.BaseDelay,
			BaseInterval:  conformance.Config.BaseInterval,
			MaxAttempts:   conformance.Config.MaxAttempts,
		},
		Region: conformance.Config.ClientRegion,
		Client: conformance.Clients.RegionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &NetworkLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTest,
		NetworkCidr:       conformance.Config.ScenariosCidr,
		publicIpsRange:    conformance.Config.ScenariosPublicIps,
		RegionZones:       conformance.Clients.RegionZones,
		InstanceSkus:      conformance.Clients.InstanceSkus,
		StorageSkus:       conformance.Clients.StorageSkus,
		NetworkSkus:       conformance.Clients.NetworkSkus,
	}

	if testLifeCycleSuite.CanRun(conformance.Config.ScenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &NetworkListV1TestSuite{
		RegionalTestSuite: regionalTest,
		NetworkCidr:       conformance.Config.ScenariosCidr,
		PublicIpsRange:    conformance.Config.ScenariosPublicIps,
		RegionZones:       conformance.Clients.RegionZones,
		InstanceSkus:      conformance.Clients.InstanceSkus,
		StorageSkus:       conformance.Clients.StorageSkus,
		NetworkSkus:       conformance.Clients.NetworkSkus,
	}
	testListSuite.RegionalTestSuite.TestSuite.ScenarioName = conformance.NetworkV1ListSuiteName
	if testListSuite.CanRun(conformance.Config.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
