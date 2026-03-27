package report

import (
	"fmt"
	"io"
	"strings"
)

func WriteText(w io.Writer, s *Summary) error {
	if _, err := fmt.Fprintf(w, "Conformance Summary — %s\n", s.GeneratedAt.Format("2006-01-02T15:04:05Z")); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Total: %d  Passed: %d  Failed: %d  Broken: %d  Skipped: %d\n\n",
		s.Totals.Total, s.Totals.Passed, s.Totals.Failed, s.Totals.Broken, s.Totals.Skipped); err != nil {
		return err
	}

	for _, sc := range s.Scenarios {
		if _, err := fmt.Fprintf(w, "%-8s %s  (%d ms)\n", strings.ToUpper(sc.Status), sc.FullName, sc.DurationMs); err != nil {
			return err
		}
		if sc.Status != "passed" && sc.Status != "skipped" {
			if sc.Error != nil && sc.Error.Message != "" {
				if _, err := fmt.Fprintf(w, "  Error: %s\n", sc.Error.Message); err != nil {
					return err
				}
			}
			if len(sc.Steps) > 0 {
				if _, err := fmt.Fprintln(w, "  Steps:"); err != nil {
					return err
				}
				for _, step := range sc.Steps {
					if err := writeStep(w, step, 2); err != nil {
						return err
					}
				}
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}

	return nil
}

func writeStep(w io.Writer, step StepResult, depth int) error {
	indent := strings.Repeat("  ", depth)
	if _, err := fmt.Fprintf(w, "%s%-8s %s  (%d ms)\n", indent, strings.ToUpper(step.Status), step.Name, step.DurationMs); err != nil {
		return err
	}
	if step.Error != nil && step.Error.Message != "" {
		if _, err := fmt.Fprintf(w, "%s  Error: %s\n", indent, step.Error.Message); err != nil {
			return err
		}
	}
	for _, child := range step.Steps {
		if err := writeStep(w, child, depth+1); err != nil {
			return err
		}
	}
	return nil
}
