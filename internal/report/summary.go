package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Summary struct {
	GeneratedAt time.Time        `json:"generated_at"`
	Totals      Totals           `json:"totals"`
	Scenarios   []ScenarioResult `json:"scenarios"`
}

type Totals struct {
	Total   int `json:"total"`
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Broken  int `json:"broken"`
	Skipped int `json:"skipped"`
}

type ScenarioResult struct {
	Name       string       `json:"name"`
	FullName   string       `json:"full_name"`
	Status     string       `json:"status"`
	DurationMs int64        `json:"duration_ms"`
	Error      *ErrorDetail `json:"error,omitempty"`
	Steps      []StepResult `json:"steps,omitempty"`
}

type StepResult struct {
	Name       string       `json:"name"`
	Status     string       `json:"status"`
	DurationMs int64        `json:"duration_ms"`
	Error      *ErrorDetail `json:"error,omitempty"`
	Steps      []StepResult `json:"steps,omitempty"`
}

type ErrorDetail struct {
	Message string `json:"message,omitempty"`
	Trace   string `json:"trace,omitempty"`
}

const (
	statusFailed  = "failed"
	statusBroken  = "broken"
	statusPassed  = "passed"
	statusSkipped = "skipped"
)

func statusPriority(status string) int {
	switch status {
	case statusFailed:
		return 0
	case statusBroken:
		return 1
	case statusPassed:
		return 2
	case statusSkipped:
		return 3
	default:
		return 4
	}
}

func convertStep(s allureStep) StepResult {
	sr := StepResult{
		Name:       s.Name,
		Status:     s.Status,
		DurationMs: s.Stop - s.Start,
	}
	if (s.Status == statusFailed || s.Status == statusBroken) &&
		(s.StatusDetails.Message != "" || s.StatusDetails.Trace != "") {
		sr.Error = &ErrorDetail{
			Message: s.StatusDetails.Message,
			Trace:   s.StatusDetails.Trace,
		}
	}
	for _, child := range s.Steps {
		sr.Steps = append(sr.Steps, convertStep(child))
	}
	return sr
}

func BuildSummary(resultsPath string) (*Summary, error) {
	entries, err := os.ReadDir(resultsPath)
	if err != nil {
		return nil, err
	}

	var scenarios []ScenarioResult
	totals := Totals{}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "-result.json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(resultsPath, entry.Name()))
		if err != nil {
			return nil, err
		}

		var ar allureResult
		if err := json.Unmarshal(data, &ar); err != nil {
			return nil, err
		}

		totals.Total++
		switch ar.Status {
		case statusPassed:
			totals.Passed++
		case statusFailed:
			totals.Failed++
		case statusBroken:
			totals.Broken++
		case statusSkipped:
			totals.Skipped++
		}

		sr := ScenarioResult{
			Name:       ar.Name,
			FullName:   ar.FullName,
			Status:     ar.Status,
			DurationMs: ar.Stop - ar.Start,
		}
		if (ar.Status == statusFailed || ar.Status == statusBroken) &&
			(ar.StatusDetails.Message != "" || ar.StatusDetails.Trace != "") {
			sr.Error = &ErrorDetail{
				Message: ar.StatusDetails.Message,
				Trace:   ar.StatusDetails.Trace,
			}
		}
		for _, step := range ar.Steps {
			sr.Steps = append(sr.Steps, convertStep(step))
		}

		scenarios = append(scenarios, sr)
	}

	sort.Slice(scenarios, func(i, j int) bool {
		pi, pj := statusPriority(scenarios[i].Status), statusPriority(scenarios[j].Status)
		if pi != pj {
			return pi < pj
		}
		return scenarios[i].Name < scenarios[j].Name
	})

	return &Summary{
		GeneratedAt: time.Now().UTC(),
		Totals:      totals,
		Scenarios:   scenarios,
	}, nil
}
