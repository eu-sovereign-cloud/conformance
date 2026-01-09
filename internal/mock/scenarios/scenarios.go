package scenarios

import "log/slog"

func LogScenarioMocking(scenarioName string) {
	slog.Info("Mocking of scenario " + scenarioName)
}
