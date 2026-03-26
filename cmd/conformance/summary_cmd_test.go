package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/report"
	"github.com/spf13/cobra"
)

func newSummaryCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "summary <results-path>",
		Short: "Print a structured summary of Allure result files",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := report.BuildSummary(args[0])
			if err != nil {
				return fmt.Errorf("building summary: %w", err)
			}
			switch format {
			case "text":
				return report.WriteText(os.Stdout, s)
			default:
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(s)
			}
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "Output format: json or text")
	return cmd
}

func maybeWriteSummary() error {
	needFile := config.Parameters.SummaryOutputPath != ""
	needStdout := config.Parameters.SummaryFormat != ""
	if !needFile && !needStdout {
		return nil
	}

	s, err := report.BuildSummary(config.Parameters.ReportResultsPath)
	if err != nil {
		return fmt.Errorf("building summary: %w", err)
	}

	if needStdout {
		switch config.Parameters.SummaryFormat {
		case "text":
			if err := report.WriteText(os.Stdout, s); err != nil {
				return fmt.Errorf("writing text summary: %w", err)
			}
		default:
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(s); err != nil {
				return fmt.Errorf("encoding json summary: %w", err)
			}
		}
	}

	if needFile {
		data, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling summary: %w", err)
		}
		if err := os.WriteFile(config.Parameters.SummaryOutputPath, data, 0o600); err != nil {
			return fmt.Errorf("writing summary file: %w", err)
		}
	}

	return nil
}
