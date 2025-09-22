package secatest

import (
	"flag"
	"fmt"
	"net"
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
	mockServerURL *string
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

	providerRegionV1, err := readRequiredFlagOrEnv(providerRegionV1Flag, providerRegionV1Config)
	if err != nil {
		return nil, err
	}

	providerAuthorizationV1, err := readRequiredFlagOrEnv(providerAuthorizationV1Flag, providerAuthorizationV1Config)
	if err != nil {
		return nil, err
	}

	clientAuthToken, err := readRequiredFlagOrEnv(clientAuthTokenFlag, clientAuthTokenConfig)
	if err != nil {
		return nil, err
	}

	clientRegion, err := readRequiredFlagOrEnv(clientRegionFlag, clientRegionConfig)
	if err != nil {
		return nil, err
	}

	clientTenant, err := readRequiredFlagOrEnv(clientTenantFlag, clientTenantConfig)
	if err != nil {
		return nil, err
	}

	scenarioUsers, err := readRequiredFlagOrEnv(scenarioUsersFlag, scenarioUsersConfig)
	if err != nil {
		return nil, err
	}
	scenarioUsersList := strings.Split(scenarioUsers, ",")

	scenarioCidr, err := readRequiredFlagOrEnv(scenarioCidrFlag, scenarioCidrConfig)
	if err != nil {
		return nil, err
	}
	_, _, err = net.ParseCIDR(scenarioCidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR format for %s: %v", scenarioCidrConfig, err)
	}

	scenarioPublicIps, err := readRequiredFlagOrEnv(scenarioPublicIpsFlag, scenarioPublicIpsConfig)
	if err != nil {
		return nil, err
	}
	_, _, err = net.ParseCIDR(scenarioPublicIps)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR format for %s: %v", scenarioPublicIpsConfig, err)
	}

	reportResultsPath := readFlagOrEnvOrDefault(reportResultsPathFlag, reportResultsPathConfig, "./reports/results")

	mockEnabled := readFlagOrEnvOrDefault(mockEnabledFlag, mockEnabledConfig, "false")

	mockServerURL := readNoRequiredFlagOrEnv(mockServerURLFlag, mockServerURLConfig)

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

func convertFlagToEnvName(param string) string {
	return strings.ToUpper(strings.ReplaceAll(param, ".", "_"))
}

func readRequiredFlagOrEnv(flag *string, param string) (string, error) {
	value := *flag

	if value == "" {
		env := convertFlagToEnvName(param)

		value = os.Getenv(env)
		if value == "" {
			return "", fmt.Errorf("missing required configuration: %s", param)
		}
	}

	return value, nil
}

func readFlagOrEnvOrDefault(flag *string, param string, defaultValue string) string {
	value := *flag

	if value == "" {
		env := convertFlagToEnvName(param)

		value = os.Getenv(env)
		if value == "" {
			return defaultValue
		}
	}

	return value
}

func readNoRequiredFlagOrEnv(flag *string, param string) *string {
	value := *flag

	if value == "" {
		env := convertFlagToEnvName(param)

		value = os.Getenv(env)
		if value == "" {
			return nil
		}
	}

	return &value
}
