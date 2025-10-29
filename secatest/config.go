package secatest

import (
	"fmt"
	"regexp"
	"sync"
)

type ConfigHolder struct {
	providerRegionV1        string
	providerAuthorizationV1 string

	clientAuthToken string
	clientRegion    string
	clientTenant    string

	scenariosFilter string
	scenariosRegexp *regexp.Regexp

	scenariosUsers     []string
	scenariosCidr      string
	scenariosPublicIps string

	reportResultsPath string

	baseDelay    int
	baseInterval int
	maxAttempts  int

	mockEnabled   bool
	mockServerURL string
}

var (
	// TODO Find a better way to share data between test suites
	config     *ConfigHolder
	configLock sync.Mutex
)

func initConfig() {
	configLock.Lock()
	defer configLock.Unlock()

	if config != nil {
		return
	}

	config = &ConfigHolder{}
}

func processConfig() error {
	configLock.Lock()
	defer configLock.Unlock()

	if config.scenariosFilter != "" {
		re, err := regexp.Compile(config.scenariosFilter)
		if err != nil {
			return fmt.Errorf("invalid scenarios.filter expression: %w", err)
		}
		config.scenariosRegexp = re
	}

	return nil
}
