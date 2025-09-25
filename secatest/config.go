package secatest

const (
	providerRegionV1Config        = "provider.region.v1"
	providerAuthorizationV1Config = "provider.authorization.v1"
	clientAuthTokenConfig         = "client.authtoken"
	clientRegionConfig            = "client.region"
	clientTenantConfig            = "client.tenant"
	scenarioUsersConfig           = "scenario.users"
	scenarioCidrConfig            = "scenario.cidr"
	scenarioPublicIpsConfig       = "scenario.publicips"
	reportResultsPathConfig       = "report.resultspath"
	reportResultsPathDefault      = "./reports/results"
	mockEnabledConfig             = "mock.enabled"
	mockServerURLConfig           = "mock.serverurl"
)

type Config struct {
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

var config *Config

func initConfig() {
	config = &Config{}
}
