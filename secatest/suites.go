package secatest

import (
	"log/slog"
	"strings"

	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/wiremock/go-wiremock"
)

type testSuite struct {
	suite.Suite
	tenant        string
	authToken     string
	mockEnabled   string
	mockServerURL string

	mockClient *wiremock.Client
}

func (suite *testSuite) isMockEnabled() bool {
	return suite.mockEnabled == "true"
}

type regionalTestSuite struct {
	testSuite
	region string
	client *secapi.RegionalClient
}

func configureTags(t provider.T, provider string, kinds ...string) {
	t.Tags(
		"provider:"+provider,
		"resources:"+strings.Join(kinds, ", "),
	)
}

func (suite *testSuite) resetAllScenarios() {
	if suite.mockClient != nil {
		if err := suite.mockClient.ResetAllScenarios(); err != nil {
			slog.Error("Failed to reset scenarios", "error", err)
		}
	}
}
