package secatest

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	RegionURL        string
	AuthorizationURL string
	RegionName       string
	AuthToken        string
	ResultsPath      string
}

const (
	regionUrlConfig        = "region_url"
	authorizationUrlConfig = "authorization_url"
	regionNameConfig       = "region_name"
	authTokenConfig        = "auth_token"
	resultsPathConfig      = "results_path"
)

func loadConfig() (*Config, error) {
	regionUrlFlag := flag.String(regionUrlConfig, "", "Region Provider Base URL")
	authorizationUrlFlag := flag.String(authorizationUrlConfig, "", "Authorization Provider Base URL")
	regionNameFlag := flag.String(regionNameConfig, "", "Regional Providers Region Name")
	authTokenFlag := flag.String(authTokenConfig, "", "JWT Authentication Token")
	resultsPathFlag := flag.String(resultsPathConfig, "../reports/results", "Report Results Path")
	flag.Parse()

	regionUrl, err := readFlagOrEnv(regionUrlFlag, regionUrlConfig)
	if err != nil {
		return nil, err
	}

	authorizationUrl, err := readFlagOrEnv(authorizationUrlFlag, authorizationUrlConfig)
	if err != nil {
		return nil, err
	}

	regionName, err := readFlagOrEnv(regionNameFlag, regionNameConfig)
	if err != nil {
		return nil, err
	}

	authToken, err := readFlagOrEnv(authTokenFlag, authTokenConfig)
	if err != nil {
		return nil, err
	}

	resultsPath := readFlagOrEnvOrDefault(resultsPathFlag, resultsPathConfig, "../reports/results")

	return &Config{
		RegionURL:        regionUrl,
		AuthorizationURL: authorizationUrl,
		RegionName:       regionName,
		AuthToken:        authToken,
		ResultsPath:      resultsPath,
	}, nil
}

func readFlagOrEnv(flag *string, param string) (string, error) {
	value := *flag

	if value == "" {
		value = os.Getenv(strings.ToUpper(param))

		if value == "" {
			return "", fmt.Errorf("missing required configuration: %s", param)
		}
	}

	return value, nil
}

func readFlagOrEnvOrDefault(flag *string, param string, defaultValue string) string {
	value := *flag

	if value == "" {
		value = os.Getenv(strings.ToUpper(param))

		if value == "" {
			return defaultValue
		}
	}

	return value
}
