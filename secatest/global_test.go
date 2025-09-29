package secatest

import (
	"log/slog"
	"os"
	"testing"

	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestGlobalSuites(t *testing.T) {
	// Initialize global client
	globalClient, err := secapi.NewGlobalClient(&secapi.GlobalConfig{
		AuthToken: config.clientAuthToken,
		Endpoints: secapi.GlobalEndpoints{
			RegionV1:        config.providerRegionV1,
			AuthorizationV1: config.providerAuthorizationV1,
		},
	})
	if err != nil {
		slog.Error("Failed to create global client", "error", err)
		os.Exit(1)
	}

	// Run test suites
	suite.RunNamedSuite(t, "Authorization V1", &AuthorizationV1TestSuite{
		globalTestSuite: globalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
			},
			client: globalClient,
		},
		users: config.scenarioUsers,
	})
}
