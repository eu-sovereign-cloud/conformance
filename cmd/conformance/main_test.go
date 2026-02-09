package main

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	setupLogger()

	config.InitParameters()

	rootCmd := initCommands(m)

	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		slog.Error("Error executing command", "error", err)
		os.Exit(1)
	}
}

func setupLogger() {
	// TODO Configure handler type and log level via env variables
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	slog.SetDefault(logger)
}

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "secatest",
		Short: "SECA Conformance Tests",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "run" {
				if err := config.ProcessParameters(); err != nil {
					return err
				}
				if err := config.InitClients(cmd.Context()); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newRunCmd(m *testing.M) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run Command",
		RunE: func(cmd *cobra.Command, args []string) error {
			configureReports()

			// Run the test suites
			code := m.Run()
			os.Exit(code)

			return nil
		},
	}
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List Command",
		RunE: func(cmd *cobra.Command, args []string) error {
			println("Available Test Scenarios:")
			println(constants.AuthorizationV1LifeCycleSuiteName)
			println(constants.AuthorizationV1ListSuiteName)
			println(constants.RegionV1ListSuiteName)
			println(constants.ComputeV1LifeCycleSuiteName)
			println(constants.ComputeV1ListSuiteName)
			println(constants.NetworkV1LifeCycleSuiteName)
			println(constants.NetworkV1ListSuiteName)
			println(constants.StorageV1LifeCycleSuiteName)
			println(constants.StorageV1ListSuiteName)
			println(constants.WorkspaceV1LifeCycleSuiteName)
			println(constants.WorkspaceV1ListSuiteName)
			println(constants.FoundationUsageV1SuiteName)

			return nil
		},
	}
}

func newReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "report",
		Short: "Report Command",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Extract the results path
			var path string
			if len(args) >= 1 {
				path = args[0]
			} else {
				path = "./reports/results"
			}

			// Run allure report
			cli := exec.Command("allure", "serve", path)
			if err := cli.Start(); err != nil {
				return err
			}

			return nil
		},
	}
}

func initCommands(m *testing.M) *cobra.Command {
	rootCmd := newRootCmd()

	runCmd := newRunCmd(m)

	runCmd.Flags().StringVar(&config.Parameters.ProviderRegionV1, "provider.region.v1", "", "Region V1 Provider Base URL")
	runCmd.Flags().StringVar(&config.Parameters.ProviderAuthorizationV1, "provider.authorization.v1", "", "Authorization V1 Provider Base URL")

	runCmd.Flags().StringVar(&config.Parameters.ClientAuthToken, "client.auth.token", "", "Client Authentication Token")
	runCmd.Flags().StringVar(&config.Parameters.ClientRegion, "client.region", "", "Client Region Name")
	runCmd.Flags().StringVar(&config.Parameters.ClientTenant, "client.tenant", "", "Client Tenant Name")

	runCmd.Flags().StringVar(&config.Parameters.ScenariosFilter, "scenarios.filter", "", "Regular expression to filter scenarios to run")
	runCmd.Flags().StringSliceVar(&config.Parameters.ScenariosUsers, "scenarios.users", nil, "Scenario Available Users")
	runCmd.Flags().StringVar(&config.Parameters.ScenariosCidr, "scenarios.cidr", "", "Scenario Available Network CIDR")
	runCmd.Flags().StringVar(&config.Parameters.ScenariosPublicIps, "scenarios.public.ips", "", "Scenario Public IPs Range")

	runCmd.Flags().StringVar(&config.Parameters.ReportResultsPath, "report.results.path", "", "Report Results Path")

	runCmd.Flags().BoolVar(&config.Parameters.MockEnabled, "mock.enabled", false, "Enable Mock Usage")
	runCmd.Flags().StringVar(&config.Parameters.MockServerURL, "mock.server.url", "", "Mock Server URL")

	runCmd.Flags().IntVar(&config.Parameters.BaseDelay, "retry.base.delay", 5, "Retry Base Delay in seconds")
	runCmd.Flags().IntVar(&config.Parameters.BaseInterval, "retry.base.interval", 30, "Retry Base Interval in seconds")
	runCmd.Flags().IntVar(&config.Parameters.MaxAttempts, "retry.max.attempts", 10, "Retry Max Attempts")

	runCmd.MarkFlagsRequiredTogether(
		"provider.region.v1",
		"client.auth.token",
		"client.region",
		"client.tenant",
	)

	rootCmd.AddCommand(runCmd)

	reportCmd := newReportCmd()
	rootCmd.AddCommand(reportCmd)

	listCmd := newListCmd()
	rootCmd.AddCommand(listCmd)

	return rootCmd
}

func configureReports() {
	resultsPath := config.Parameters.ReportResultsPath

	outputPath := filepath.Dir(resultsPath)
	if err := os.Setenv("ALLURE_OUTPUT_PATH", outputPath); err != nil {
		slog.Error("Failed to configure reports", "error", err)
	}

	outputFolder := filepath.Base(resultsPath)
	if err := os.Setenv("ALLURE_OUTPUT_FOLDER", outputFolder); err != nil {
		slog.Error("Failed to configure reports", "error", err)
	}
}
