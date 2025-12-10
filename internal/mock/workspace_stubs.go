package mock

import (
	"net/http"

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
