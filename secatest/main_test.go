package secatest

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	setupLogger()

	initConfig()

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
				if err := processConfig(); err != nil {
					return err
				}
				if err := initClients(cmd.Context()); err != nil {
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
			println(authorizationV1LifeCycleSuiteName)
			println(computeV1LifeCycleSuiteName)
			println(networkV1LifeCycleSuiteName)
			println(storageV1LifeCycleSuiteName)
			println(workspaceV1LifeCycleSuiteName)
			println(foundationV1UsageSuiteName)
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

	runCmd.Flags().StringVar(&config.providerRegionV1, "provider.region.v1", "", "Region V1 Provider Base URL")
	runCmd.Flags().StringVar(&config.providerAuthorizationV1, "provider.authorization.v1", "", "Authorization V1 Provider Base URL")

	runCmd.Flags().StringVar(&config.clientAuthToken, "client.auth.token", "", "Client Authentication Token")
	runCmd.Flags().StringVar(&config.clientRegion, "client.region", "", "Client Region Name")
	runCmd.Flags().StringVar(&config.clientTenant, "client.tenant", "", "Client Tenant Name")

	runCmd.Flags().StringVar(&config.scenariosFilter, "scenarios.filter", "", "Regular expression to filter scenarios to run")
	runCmd.Flags().StringSliceVar(&config.scenariosUsers, "scenarios.users", nil, "Scenario Available Users")
	runCmd.Flags().StringVar(&config.scenariosCidr, "scenarios.cidr", "", "Scenario Available Network CIDR")
	runCmd.Flags().StringVar(&config.scenariosPublicIps, "scenarios.public.ips", "", "Scenario Public IPs Range")

	runCmd.Flags().StringVar(&config.reportResultsPath, "report.results.path", "", "Report Results Path")

	runCmd.Flags().BoolVar(&config.mockEnabled, "mock.enabled", false, "Enable Mock Usage")
	runCmd.Flags().StringVar(&config.mockServerURL, "mock.server.url", "", "Mock Server URL")

	runCmd.MarkFlagsRequiredTogether("provider.region.v1", "provider.authorization.v1", "client.auth.token", "client.region", "client.tenant", "scenarios.users", "scenarios.cidr", "scenarios.public.ips")

	rootCmd.AddCommand(runCmd)

	reportCmd := newReportCmd()
	rootCmd.AddCommand(reportCmd)

	listCmd := newListCmd()
	rootCmd.AddCommand(listCmd)

	return rootCmd
}

func configureReports() {
	resultsPath := config.reportResultsPath

	outputPath := filepath.Dir(resultsPath)
	if err := os.Setenv("ALLURE_OUTPUT_PATH", outputPath); err != nil {
		slog.Error("Failed to configure reports", "error", err)
	}

	outputFolder := filepath.Base(resultsPath)
	if err := os.Setenv("ALLURE_OUTPUT_FOLDER", outputFolder); err != nil {
		slog.Error("Failed to configure reports", "error", err)
	}
}
