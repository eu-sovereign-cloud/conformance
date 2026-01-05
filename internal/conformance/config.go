package conformance

import (
	"fmt"
	"regexp"
	"sync"
)

type ConfigHolder struct {
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
	Config     *ConfigHolder
	configLock sync.Mutex
)

func InitConfig() {
	configLock.Lock()
	defer configLock.Unlock()

	if Config != nil {
		return
	}

	Config = &ConfigHolder{}
}

func ProcessConfig() error {
	configLock.Lock()
	defer configLock.Unlock()

	if Config.ScenariosFilter != "" {
		re, err := regexp.Compile(Config.ScenariosFilter)
		if err != nil {
			return fmt.Errorf("invalid scenarios.filter expression: %w", err)
		}
		Config.ScenariosRegexp = re
	}

	return nil
}
