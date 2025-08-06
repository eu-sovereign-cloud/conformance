package mock

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	WireMockURL         = "http://localhost:8080"
	testPutWorkspaceURL = "/providers/seca.workspace/v1/tenants/tenant1/workspaces/workspace1"
	TenantName          = "tenant1-"
	WorkspaceName       = "workspace-1"
	Region              = "region-1"
	Token               = "your_token"
	Version             = "v1"
)

/*
func TestProcessTemplate() {

}
*/
func TestWorkspaceScenario(t *testing.T) {
	WorkspaceMock := MockParams{
		WireMockURL:   WireMockURL,
		TenantName:    TenantName,
		WorkspaceName: WorkspaceName,
		Region:        Region,
		Token:         Token,
	}
	err := CreateWorkspaceScenario(WorkspaceMock)
	if err != nil {
		log.Printf("Error creating workspace scenario: %v\n", err)
		return
	}

	// Create Workspace
	url := WireMockURL + testPutWorkspaceURL
	responseUpdate, error := requestWorkspace("PUT", url, Token)
	if error != nil {
		log.Printf("Error updating workspace: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusCreated, responseUpdate.StatusCode, "Expected status code 201 OK")

	// Update Workspace
	responseUpdate, error = requestWorkspace("PUT", url, Token)
	if error != nil {
		log.Printf("Error updating workspace: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusOK, responseUpdate.StatusCode, "Expected status code 200 OK")

	// Delete Workspace
	responseUpdate, error = requestWorkspace("DELETE", url, Token)
	if error != nil {
		log.Printf("Error deleting workspace: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusAccepted, responseUpdate.StatusCode, "Expected status code 202 No Content")

	// Delete workspace 2 time
	responseUpdate, error = requestWorkspace("DELETE", url, Token)
	if error != nil {
		log.Printf("Error deleting workspace: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusAccepted, responseUpdate.StatusCode, "Expected status code 202 No Content")

	// Get workspace
	responseUpdate, error = requestWorkspace("GET", url, Token)
	if error != nil {
		log.Printf("Error getting workspace: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusAccepted, responseUpdate.StatusCode, "Expected status code 202 No Content")
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
