package mock

import (
	"net/http"

	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Workspace

func (configurator *scenarioConfigurator) configureCreateWorkspaceStub(response *schema.Workspace, url string, params HasParams) error {
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = newWorkspaceStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateWorkspaceStubWithLabels(response *schema.Workspace, url string, params HasParams, labels schema.Labels) error {
	setModifiedRegionalResourceMetadata(response.Metadata)
	setWorkspaceState(response.Status, schema.ResourceStateUpdating)
	response.Labels = labels
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveWorkspaceStub(response *schema.Workspace, url string, params HasParams) error {
	setWorkspaceState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListActiveWorkspaceStub(response *workspace.WorkspaceIterator, url string, params HasParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetListStub(url, params, pathParams, response, false); err != nil {
		return err
	}
	return nil
}
