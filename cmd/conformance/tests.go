package main

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
)

func createTestSuite(params *config.ParametersHolder) *suites.TestSuite {
	return &suites.TestSuite{
		Tenant:        params.ClientTenant,
		AuthToken:     params.ClientAuthToken,
		MockEnabled:   params.MockEnabled,
		MockServerURL: &params.MockServerURL,
		BaseDelay:     params.BaseDelay,
		BaseInterval:  params.BaseInterval,
		MaxAttempts:   params.MaxAttempts,
	}
}

func createGlobalTestSuite(params *config.ParametersHolder, clients *config.ClientsHolder) suites.GlobalTestSuite {
	return suites.GlobalTestSuite{
		TestSuite: createTestSuite(params),
		Client:    clients.GlobalClient,
	}
}

func createRegionalTestSuite(params *config.ParametersHolder, clients *config.ClientsHolder) suites.RegionalTestSuite {
	return suites.RegionalTestSuite{
		TestSuite: createTestSuite(params),
		Region:    params.ClientRegion,
		Client:    clients.RegionalClient,
	}
}

func createMixedTestSuite(params *config.ParametersHolder, clients *config.ClientsHolder) suites.MixedTestSuite {
	return suites.MixedTestSuite{
		TestSuite:      createTestSuite(params),
		Region:         params.ClientRegion,
		GlobalClient:   clients.GlobalClient,
		RegionalClient: clients.RegionalClient,
	}
}
