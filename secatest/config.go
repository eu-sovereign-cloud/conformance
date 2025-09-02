package secatest

import (
	"flag"
	"fmt"
	"os"
	"strings"
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

	mockEnabled   string
	mockServerURL string
}

const (
	providerRegionV1Config        = "seca.provider.region.v1"
	providerAuthorizationV1Config = "seca.provider.authorization.v1"

	clientAuthTokenConfig = "seca.client.authtoken"
	clientRegionConfig    = "seca.client.region"
	clientTenantConfig    = "seca.client.tenant"

	scenarioUsersConfig     = "seca.scenario.users"
	scenarioCidrConfig      = "seca.scenario.cidr"
	scenarioPublicIpsConfig = "seca.scenario.publicips"

	reportResultsPathConfig = "seca.report.resultspath"

	mockEnabledConfig   = "seca.mock.enabled"
	mockServerURLConfig = "seca.mock.serverurl"
)

func loadConfig() (*Config, error) {
	providerRegionV1Flag := flag.String(providerRegionV1Config, "", "Region V1 Provider Base URL")
	providerAuthorizationV1Flag := flag.String(providerAuthorizationV1Config, "", "Authorization V1 Provider Base URL")

	clientAuthTokenFlag := flag.String(clientAuthTokenConfig, "", "Client Authentication Token")
	clientRegionFlag := flag.String(clientRegionConfig, "", "Client Region Name")
	clientTenantFlag := flag.String(clientTenantConfig, "", "Client Tenant Name")

	scenarioUsersFlag := flag.String(scenarioUsersConfig, "", "Scenario Available Users")
	scenarioCidrFlag := flag.String(scenarioCidrConfig, "", "Scenario Available Network CIDR")
	scenarioPublicIpsFlag := flag.String(scenarioPublicIpsConfig, "", "Scenario Public IPs Range")

	reportResultsPathFlag := flag.String(reportResultsPathConfig, "", "Report Results Path")

	mockEnabledFlag := flag.String(mockEnabledConfig, "", "Enable Mock Usage")
	mockServerURLFlag := flag.String(mockServerURLConfig, "", "Mock Server URL")

	flag.Parse()

	providerRegionV1, err := readFlagOrEnv(providerRegionV1Flag, providerRegionV1Config)
	if err != nil {
		return nil, err
	}

	providerAuthorizationV1, err := readFlagOrEnv(providerAuthorizationV1Flag, providerAuthorizationV1Config)
	if err != nil {
		return nil, err
	}

	clientAuthToken, err := readFlagOrEnv(clientAuthTokenFlag, clientAuthTokenConfig)
	if err != nil {
		return nil, err
	}

	clientRegion, err := readFlagOrEnv(clientRegionFlag, clientRegionConfig)
	if err != nil {
		return nil, err
	}

	clientTenant, err := readFlagOrEnv(clientTenantFlag, clientTenantConfig)
	if err != nil {
		return nil, err
	}

	scenarioUsers, err := readFlagOrEnv(scenarioUsersFlag, scenarioUsersConfig)
	if err != nil {
		return nil, err
	}
	scenarioUsersList := strings.Split(scenarioUsers, ",")

	scenarioCidr, err := readFlagOrEnv(scenarioCidrFlag, scenarioCidrConfig)
	if err != nil {
		return nil, err
	}

	scenarioPublicIps, err := readFlagOrEnv(scenarioPublicIpsFlag, scenarioPublicIpsConfig)
	if err != nil {
		return nil, err
	}

	reportResultsPath, err := readFlagOrEnv(reportResultsPathFlag, reportResultsPathConfig)
	if err != nil {
		return nil, err
	}

	mockEnabled, err := readFlagOrEnv(mockEnabledFlag, mockEnabledConfig)
	if err != nil {
		return nil, err
	}

	mockServerURL, err := readFlagOrEnv(mockServerURLFlag, mockServerURLConfig)
	if err != nil {
		return nil, err
	}

	return &Config{
		providerRegionV1:        providerRegionV1,
		providerAuthorizationV1: providerAuthorizationV1,
		clientAuthToken:         clientAuthToken,
		clientRegion:            clientRegion,
		clientTenant:            clientTenant,
		scenarioUsers:           scenarioUsersList,
		scenarioCidr:            scenarioCidr,
		scenarioPublicIps:       scenarioPublicIps,
		reportResultsPath:       reportResultsPath,
		mockEnabled:             mockEnabled,
		mockServerURL:           mockServerURL,
	}, nil
}

func readFlagOrEnv(flag *string, param string) (string, error) {
	value := *flag

	if value == "" {
		// Convert flag to environment variable name format
		env := strings.ToUpper(strings.ReplaceAll(param, ".", "_"))

		value = os.Getenv(env)
		if value == "" {
			return "", fmt.Errorf("missing required configuration: %s", param)
		}
	}

	return value, nil
}
