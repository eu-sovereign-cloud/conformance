package mock

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/stretchr/testify/assert"
)

const (
	WireMockURL = "http://localhost:8080"

	TenantName    = "tenant1"
	workspaceName = "workspace1"
	Region        = "eu-central-1"
	Token         = "your_token"
	Version       = "v1"
	Zone          = "zone1"
)

func Test_workspace(t *testing.T) {
	// TODO: implement test
	TenantName := "test-tenant"

	workspaceName := secalib.GenerateWorkspaceName()
	url := secalib.GenerateWorkspaceURL(TenantName, workspaceName)

	workspaceParams := WorkspaceParamsV1{
		Params: &Params{
			MockURL:   WireMockURL,
			AuthToken: Token,
			Tenant:    TenantName,
			Region:    Region,
		},
		Workspace: &ResourceParams[secalib.WorkspaceSpecV1]{
			Name: workspaceName,
			InitialSpec: &secalib.WorkspaceSpecV1{
				Label: &[]secalib.Label{
					{
						Name:  "env",
						Value: "test",
					},
				},
			},
			UpdatedSpec: &secalib.WorkspaceSpecV1{
				Label: &[]secalib.Label{
					{
						Name:  "env",
						Value: "prod",
					},
				},
			},
		},
	}

	wm, err := CreateWorkspaceLifecycleScenarioV1(fmt.Sprintf("workspace lifecycle %d", rand.Intn(100)), workspaceParams)
	if err != nil {
		t.Errorf("Failed to create workspace scenario: %v\n", err)
		t.FailNow()
		return
	}
	if wm == nil {
		t.Errorf("Failed to create workspace scenario: %v\n", err)
		t.FailNow()
		return
	}

	//Workspace
	// Create Workspace
	response, error := request("PUT", url, Token)
	if error != nil {
		t.Errorf("Error Put Workspace: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode, "Expected status code 201 for Put Workspace")
	// Get Workspace
	response, error = request("GET", url, Token)
	if error != nil {
		t.Errorf("Error Get Workspace: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Workspace")

	// Update Workspace
	response, error = request("PUT", url, Token)
	if error != nil {
		t.Errorf("Error Put Workspace: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Put Workspace")

	// Get Workspace
	response, error = request("GET", url, Token)
	if error != nil {
		t.Errorf("Error Get Workspace: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Workspace")

}
