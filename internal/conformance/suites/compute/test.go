package compute

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestComputeV1Suite(t *testing.T) {
	regionalTestSuite := suites.RegionalTestSuite{
		TestSuite: suites.TestSuite{
			Tenant:        conformance.Config.ClientTenant,
			AuthToken:     conformance.Config.ClientAuthToken,
			MockEnabled:   conformance.Config.MockEnabled,
			MockServerURL: &conformance.Config.MockServerURL,
			ScenarioName:  conformance.ComputeV1LifeCycleSuiteName,
			BaseDelay:     conformance.Config.BaseDelay,
			BaseInterval:  conformance.Config.BaseInterval,
			MaxAttempts:   conformance.Config.MaxAttempts,
		},
		Region: conformance.Config.ClientRegion,
		Client: conformance.Clients.RegionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &ComputeV1LifeCycleTestSuite{
		RegionalTestSuite: regionalTestSuite,
		AvailableZones:    conformance.Clients.RegionZones,
		InstanceSkus:      conformance.Clients.InstanceSkus,
		StorageSkus:       conformance.Clients.StorageSkus,
	}

	if testLifeCycleSuite.CanRun(conformance.Config.ScenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &ComputeV1ListTestSuite{
		RegionalTestSuite: regionalTestSuite,
		AvailableZones:    conformance.Clients.RegionZones,
		InstanceSkus:      conformance.Clients.InstanceSkus,
		StorageSkus:       conformance.Clients.StorageSkus,
	}

	testListSuite.RegionalTestSuite.ScenarioName = conformance.ComputeV1ListSuiteName

	if testListSuite.CanRun(conformance.Config.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
