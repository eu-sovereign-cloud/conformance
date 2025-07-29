package secatest

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var config *Config

func TestMain(m *testing.M) {
	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	var err error
	config, err = loadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Configure allure
	configureReports(config)

	// Run tests
	os.Exit(m.Run())
}

func TestSuites(t *testing.T) {
	ctx := context.Background()

	// Initialize global client
	globalClient, err := secapi.NewGlobalClient(&secapi.GlobalConfig{
		AuthToken: config.AuthToken,
		Endpoints: secapi.GlobalEndpoints{
			RegionV1:        config.RegionURL,
			AuthorizationV1: config.AuthorizationURL,
		},
	})
	if err != nil {
		slog.Error("Failed to create global client", "error", err)
		os.Exit(1)
	}

	// Initialize regional client
	regionalClient, err := globalClient.NewRegionalClient(ctx, config.RegionName)
	if err != nil {
		slog.Error("Failed to create regional client", "error", err)
		os.Exit(1)
	}

	// Run test suites
	suite.RunNamedSuite(t, "Workspace V1", &WorkspaceV1TestSuite{
		client: regionalClient,
	})
}

func configureReports(config *Config) {
	resultsPath := config.ResultsPath

	outputPath := filepath.Dir(resultsPath)
	os.Setenv("ALLURE_OUTPUT_PATH", outputPath)

	outputFolder := filepath.Base(resultsPath)
	os.Setenv("ALLURE_OUTPUT_FOLDER", outputFolder)
}
