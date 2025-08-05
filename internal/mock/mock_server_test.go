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
	testComputeSkuURL      = "/providers/seca.compute/v1/tenants/tenant1/skus/D2XS"
	testComputeInstanceUrl = "/providers/seca.compute/v1/tenants/tenant1/instances/workspace1/compute"
	TenantName             = "tenant1"
	WorkspaceName          = "workspace1"
	Region                 = "eu-central-1"
	Token                  = "your_token"
	Version                = "v1"
)

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
	responseUpdate, error := requestWorkspace("PUT", url, Token)
	if error != nil {
		log.Printf("Error updating workspace: %v\n", error)
		assert.Error(t, err, "Expected error when updating workspace")
	}
	assert.Equal(t, http.StatusCreated, responseUpdate.StatusCode, "Expected status code 201 OK")

	// Update Workspace
	responseUpdate, error = requestWorkspace("PUT", url, Token)
	if error != nil {
		log.Printf("Error updating workspace: %v\n", error)
		assert.Error(t, err, "Expected error when updating workspace")
	}
	assert.Equal(t, http.StatusOK, responseUpdate.StatusCode, "Expected status code 200 OK")

	// Delete Workspace
	responseUpdate, error = requestWorkspace("DELETE", url, Token)
	if error != nil {
		log.Printf("Error deleting workspace: %v\n", error)
		assert.Error(t, err, "Expected error when deleting workspace")
	}
	assert.Equal(t, http.StatusAccepted, responseUpdate.StatusCode, "Expected status code 202 No Content")

	// Delete workspace 2 time
	responseUpdate, error = requestWorkspace("DELETE", url, Token)
	if error != nil {
		log.Printf("Error deleting workspace: %v\n", error)
		assert.Error(t, err, "Expected error when deleting workspace")
	}
	assert.Equal(t, http.StatusAccepted, responseUpdate.StatusCode, "Expected status code 202 No Content")
}

func TestComputeScenario(t *testing.T) {
	ComputeMock := MockParams{
		WireMockURL:   WireMockURL,
		TenantName:    TenantName,
		WorkspaceName: WorkspaceName,
		Name:          "D2XS",
		Region:        Region,
		Token:         Token,
	}
	err := CreateComputeScenario(ComputeMock)
	if err != nil {
		log.Printf("Error creating compute scenario: %v\n", err)
		assert.Error(t, err, "Expected error when creating compute scenario")
	}

	// Get Sku
	url := WireMockURL + testPutWorkspaceURL
	responseUpdate, error := requestWorkspace("GET", url, Token)
	if error != nil {
		log.Printf("Error updating workspace: %v\n", error)
		assert.Error(t, err, "Expected error when updating workspace")
	}
	assert.Equal(t, http.StatusOK, responseUpdate.StatusCode, "Expected status code 200 OK")
}

func requestWorkspace(method string, url string, token string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %v\n", url, err)
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending PUT request to %s: %v\n", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}
