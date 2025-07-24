package secapi

import (
	"context"
	"log/slog"
	"os"
	"testing"

	sdk "github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestAll(t *testing.T) {
	ctx := context.Background()

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

	// Read Configurations
	regionURL := os.Getenv("REGION_PROVIDER_URL")
	if regionURL == "" {
		slog.Error("REGION_PROVIDER_URL is not set")
		os.Exit(1)
	}

	authorizationURL := os.Getenv("AUTHORIZATION_PROVIDER_URL")
	if authorizationURL == "" {
		slog.Error("AUTHORIZATION_PROVIDER_URL is not set")
		os.Exit(1)
	}

	regionName := os.Getenv("REGION_NAME")
	if regionName == "" {
		slog.Error("REGION_NAME is not set")
		os.Exit(1)
	}

	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		slog.Error("AUTH_TOKEN is not set")
		os.Exit(1)
	}

	// Setup clients
	globalClient, err := sdk.NewGlobalClient(&sdk.GlobalEndpoints{
		RegionV1:        regionURL,
		AuthorizationV1: authorizationURL,
	})
	if err != nil {
		slog.Error("Failed to create global client", "error", err)
		os.Exit(1)
	}

	regionalClient, err := globalClient.NewRegionalClient(ctx, regionName)
	if err != nil {
		slog.Error("Failed to create regional client", "error", err)
		os.Exit(1)
	}

	suite.RunNamedSuite(t, "Workspace Provider", &WorkspaceTestSuite{
		client: regionalClient,
	})
}
