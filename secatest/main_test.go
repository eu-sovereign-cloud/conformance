package secatest

import (
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

	if err := rootCmd.Execute(); err != nil {
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
				path = reportResultsPathDefault
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
	runCmd.Flags().StringVar(&config.providerRegionV1, providerRegionV1Config, "", "Region V1 Provider Base URL")
	runCmd.Flags().StringVar(&config.providerAuthorizationV1, providerAuthorizationV1Config, "", "Authorization V1 Provider Base URL")
	runCmd.Flags().StringVar(&config.clientAuthToken, clientAuthTokenConfig, "", "Client Authentication Token")
	runCmd.Flags().StringVar(&config.clientRegion, clientRegionConfig, "", "Client Region Name")
	runCmd.Flags().StringVar(&config.clientTenant, clientTenantConfig, "", "Client Tenant Name")
	runCmd.Flags().StringSliceVar(&config.scenarioUsers, scenarioUsersConfig, nil, "Scenario Available Users")
	runCmd.Flags().StringVar(&config.scenarioCidr, scenarioCidrConfig, "", "Scenario Available Network CIDR")
	runCmd.Flags().StringVar(&config.scenarioPublicIps, scenarioPublicIpsConfig, "", "Scenario Public IPs Range")
	runCmd.Flags().StringVar(&config.reportResultsPath, reportResultsPathConfig, "", "Report Results Path")
	runCmd.Flags().BoolVar(&config.mockEnabled, mockEnabledConfig, false, "Enable Mock Usage")
	runCmd.Flags().StringVar(&config.mockServerURL, mockServerURLConfig, "", "Mock Server URL")
	rootCmd.AddCommand(runCmd)

	reportCmd := newReportCmd()
	rootCmd.AddCommand(reportCmd)

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
