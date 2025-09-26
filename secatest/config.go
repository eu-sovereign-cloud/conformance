package secatest

import "sync"

type ConfigHolder struct {
	providerRegionV1        string
	providerAuthorizationV1 string

	clientAuthToken string
	clientRegion    string
	clientTenant    string

	scenarioUsers     []string
	scenarioCidr      string
	scenarioPublicIps string

	reportResultsPath string

	mockEnabled   bool
	mockServerURL string
}

var (
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
