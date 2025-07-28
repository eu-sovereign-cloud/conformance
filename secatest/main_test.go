package secatest

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestMain(m *testing.M) {
	// Configure the report
	if os.Getenv("ALLURE_OUTPUT_PATH") == "" {
		os.Setenv("ALLURE_OUTPUT_PATH", "../reports")
	}
	if os.Getenv("ALLURE_OUTPUT_FOLDER") == "" {
		os.Setenv("ALLURE_OUTPUT_FOLDER", "results")
	}

	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	m.Run()
}

func TestSuites(t *testing.T) {
	ctx := context.Background()

	// Read Configurations
	// regionURL := os.Getenv("REGION_PROVIDER_URL")
	regionURL := "http://localhost:8080/providers/seca.region"
	if regionURL == "" {
		slog.Error("REGION_PROVIDER_URL is not set")
		os.Exit(1)
	}

	// authorizationURL := os.Getenv("AUTHORIZATION_PROVIDER_URL"
	authorizationURL := "http://localhost:8080/providers/seca.authorization"

	// regionName := os.Getenv("REGION_NAME")
	regionName := "eu-central-1"
	if regionName == "" {
		slog.Error("REGION_NAME is not set")
		os.Exit(1)
	}

	// authToken := os.Getenv("AUTH_TOKEN")
	authToken := "test-token"
	if authToken == "" {
		slog.Error("AUTH_TOKEN is not set")
		os.Exit(1)
	}

	// Setup clients
	config := &secapi.GlobalConfig{
		AuthToken: authToken,
		Endpoints: secapi.GlobalEndpoints{
			RegionV1:        regionURL,
			AuthorizationV1: authorizationURL,
		},
	}
	globalClient, err := secapi.NewGlobalClient(config)
	if err != nil {
		slog.Error("Failed to create global client", "error", err)
		os.Exit(1)
	}

	regionalClient, err := globalClient.NewRegionalClient(ctx, regionName)
	if err != nil {
		slog.Error("Failed to create regional client", "error", err)
		os.Exit(1)
	}

	suite.RunNamedSuite(t, "Workspace V1", &WorkspaceV1TestSuite{
		client: regionalClient,
	})
}
