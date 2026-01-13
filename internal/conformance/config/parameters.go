package config

import (
	"fmt"
	"regexp"
	"sync"
)

type ParametersHolder struct {
	ProviderRegionV1        string
	ProviderAuthorizationV1 string

	ClientAuthToken string
	ClientRegion    string
	ClientTenant    string

	ScenariosFilter string
	ScenariosRegexp *regexp.Regexp

	ScenariosUsers     []string
	ScenariosCidr      string
	ScenariosPublicIps string

	ReportResultsPath string

	BaseDelay    int
	BaseInterval int
	MaxAttempts  int

	MockEnabled   bool
	MockServerURL string
}

var (
	Parameters     *ParametersHolder
	parametersLock sync.Mutex
)

func InitParameters() {
	parametersLock.Lock()
	defer parametersLock.Unlock()

	if Parameters != nil {
		return
	}

	Parameters = &ParametersHolder{}
}

func ProcessParameters() error {
	parametersLock.Lock()
	defer parametersLock.Unlock()

	if Parameters.ScenariosFilter != "" {
		expr, err := regexp.Compile(Parameters.ScenariosFilter)
		if err != nil {
			return fmt.Errorf("invalid scenarios.filter expression: %w", err)
		}
		Parameters.ScenariosRegexp = expr
	}

	return nil
}
