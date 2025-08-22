package secatest

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand"
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

func (suite *testSuite) generateSkuName() string {
	return fmt.Sprintf("sku-%d", rand.Intn(math.MaxInt32))
}

func (suite *testSuite) generateSkuRef(name string) string {
	return fmt.Sprintf(skuRef, name)
}

func (suite *testSuite) generateWorkspaceName() string {
	return fmt.Sprintf("workspace-%d", rand.Intn(math.MaxInt32))
}

func (suite *testSuite) generateWorkspaceResource(name string) string {
	return fmt.Sprintf(workspaceResource, suite.tenant, name)
}

func (suite *testSuite) generateBlockStorageName() string {
	return fmt.Sprintf("disk-%d", rand.Intn(math.MaxInt32))
}

func (suite *testSuite) generateBlockStorageResource(workspace string, blockStorage string) string {
	return fmt.Sprintf(blockStorageResource, suite.tenant, workspace, blockStorage)
}

func (suite *testSuite) generateBlockStorageRef(blockStorageName string) string {
	return fmt.Sprintf(blockStoragesRef, blockStorageName)
}

func (suite *testSuite) generateImageName() string {
	return fmt.Sprintf("image-%d", rand.Intn(math.MaxInt32))
}

func (suite *testSuite) generateImageResource(name string) string {
	return fmt.Sprintf(imageResource, suite.tenant, name)
}

func (suite *testSuite) generateInstanceName() string {
	return fmt.Sprintf("instance-%d", rand.Intn(math.MaxInt32))
}

func (suite *testSuite) generateInstanceResource(workspace string, instance string) string {
	return fmt.Sprintf(instanceResource, suite.tenant, workspace, instance)
}
