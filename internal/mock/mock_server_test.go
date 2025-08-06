package mock

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	WireMockURL            = "http://localhost:8080"
	testPutWorkspaceURL    = "/providers/seca.workspace/v1/tenants/tenant1/workspaces/workspace1"
	testComputeSkuURL      = "/providers/seca.compute/v1/tenants/tenant1/skus/"
	testComputeInstanceUrl = "/providers/seca.compute/v1/tenants/tenant1/workspaces/workspace1/instances/" + ComputeName
	TenantName             = "tenant1"
	WorkspaceName          = "workspace1"
	ComputeName            = "compute1"
	Region                 = "eu-central-1"
	Token                  = "your_token"
	Version                = "v1"
	ComputeSkuName         = "D2SX"
)

/*
func TestProcessTemplate() {

}
*/
func TestWorkspaceScenario(t *testing.T) {
	WorkspaceMock := MockParams{
		WireMockURL: WireMockURL,
		TenantName:  TenantName,
		Name:        WorkspaceName,
		Region:      Region,
		Token:       Token,
	}
	err := CreateWorkspaceScenario(WorkspaceMock)
	if err != nil {
		log.Printf("Error creating workspace scenario: %v\n", err)
		assert.Error(t, err, "Expected error when creating workspace scenario")
	}

	// Create Workspace
	url := WireMockURL + testPutWorkspaceURL
	responseUpdate, error := requestMethod("PUT", url, WorkspaceMock.Token)
	if error != nil {
		log.Printf("Error updating workspace: %v\n", error)
		assert.Error(t, err, "Expected error when updating workspace")
	}
	assert.Equal(t, http.StatusCreated, responseUpdate.StatusCode, "Expected status code 201 OK")

	// Update Workspace
	responseUpdate, error = requestMethod("PUT", url, WorkspaceMock.Token)
	if error != nil {
		log.Printf("Error updating workspace: %v\n", error)
		assert.Error(t, err, "Expected error when updating workspace")
	}
	assert.Equal(t, http.StatusOK, responseUpdate.StatusCode, "Expected status code 200 OK")

	// Delete Workspace
	responseUpdate, error = requestMethod("DELETE", url, WorkspaceMock.Token)
	if error != nil {
		log.Printf("Error deleting workspace: %v\n", error)
		assert.Error(t, err, "Expected error when deleting workspace")
	}
	assert.Equal(t, http.StatusAccepted, responseUpdate.StatusCode, "Expected status code 202 No Content")

	// Delete workspace 2 time
	responseUpdate, error = requestMethod("DELETE", url, WorkspaceMock.Token)
	if error != nil {
		log.Printf("Error deleting workspace: %v\n", error)
		assert.Error(t, err, "Expected error when deleting workspace")
	}
	assert.Equal(t, http.StatusAccepted, responseUpdate.StatusCode, "Expected status code 202 No Content")

	// Get workspace
	responseUpdate, error = requestMethod("GET", url, Token)
	if error != nil {
		log.Printf("Error getting workspace: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusAccepted, responseUpdate.StatusCode, "Expected status code 202 No Content")
}

func TestComputeScenario(t *testing.T) {
	ComputeMock := MockParams{
		WireMockURL:   WireMockURL,
		TenantName:    TenantName,
		WorkspaceName: WorkspaceName,
		Name:          ComputeName,
		Region:        Region,
		Token:         Token,
		SkuName:       ComputeSkuName,
	}
	err := CreateComputeScenario(ComputeMock)
	if err != nil {
		log.Printf("Error creating compute scenario: %v\n", err)
		assert.Error(t, err, "Expected error when creating compute scenario")
	}

	// Get Sku

	url := WireMockURL + testComputeSkuURL + ComputeMock.SkuName
	responseUpdate, error := requestMethod("GET", url, ComputeMock.Token)
	if error != nil {
		log.Printf("Error updating workspace: %v\n", error)
		assert.Error(t, err, "Expected error when updating workspace")
	}
	assert.Equal(t, http.StatusOK, responseUpdate.StatusCode, "Expected status code 200 OK")

	// Create Compute Instance
	url = WireMockURL + testComputeInstanceUrl
	responseUpdate, error = requestMethod("PUT", url, ComputeMock.Token)
	if error != nil {
		log.Printf("Error creating workspace: %v\n", error)
		assert.Error(t, err, "Expected error when creating workspace")
	}
	assert.Equal(t, http.StatusCreated, responseUpdate.StatusCode, "Expected status code 201 Created")

	// Get Compute Instance
	url = WireMockURL + testComputeInstanceUrl
	responseUpdate, error = requestMethod("GET", url, ComputeMock.Token)
	if error != nil {
		log.Printf("Error getting workspace: %v\n", error)
		assert.Error(t, err, "Expected error when getting workspace")
	}
	assert.Equal(t, http.StatusOK, responseUpdate.StatusCode, "Expected status code 200 OK")

	// Update Compute Instance
	url = WireMockURL + testComputeInstanceUrl
	responseUpdate, error = requestMethod("PUT", url, ComputeMock.Token)
	if error != nil {
		log.Printf("Error updating workspace: %v\n", error)
		assert.Error(t, err, "Expected error when updating workspace")
	}
	assert.Equal(t, http.StatusOK, responseUpdate.StatusCode, "Expected status code 200 OK")
}

func requestMethod(method string, url string, token string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %v\n", url, err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending PUT request to %s: %v\n", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}
