package secatest

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ProviderRegionV1        string
	ProviderAuthorizationV1 string
	ClientAuthToken         string
	ClientRegion            string
	ClientTenant            string
	ReportResultsPath       string
}

const (
	providerRegionV1Config        = "seca.provider.region.v1"
	providerAuthorizationV1Config = "seca.provider.authorization.v1"
	clientAuthTokenConfig         = "seca.client.authtoken"
	clientRegionConfig            = "seca.client.region"
	clientTenantConfig            = "seca.client.tenant"
	reportResultsPathConfig       = "seca.report.resultspath"
)

func loadConfig() (*Config, error) {
	providerRegionV1Flag := flag.String(providerRegionV1Config, "", "Region V1 Provider Base URL")
	providerAuthorizationV1Flag := flag.String(providerAuthorizationV1Config, "", "Authorization V1 Provider Base URL")
	clientAuthTokenFlag := flag.String(clientAuthTokenConfig, "", "Client Authentication Token")
	clientRegionFlag := flag.String(clientRegionConfig, "", "Client Region Name")
	clientTenantFlag := flag.String(clientTenantConfig, "", "Client Tenant Name")
	reportResultsPathFlag := flag.String(reportResultsPathConfig, "", "Report Results Path")
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

	reportResultsPath, err := readFlagOrEnv(reportResultsPathFlag, reportResultsPathConfig)
	if err != nil {
		return nil, err
	}

	return &Config{
		ProviderRegionV1:        providerRegionV1,
		ProviderAuthorizationV1: providerAuthorizationV1,
		ClientAuthToken:         clientAuthToken,
		ClientRegion:            clientRegion,
		ClientTenant:            clientTenant,
		ReportResultsPath:       reportResultsPath,
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
