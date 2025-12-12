package mock

import (
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Instance

func (configurator *scenarioConfigurator) configureCreateInstanceStub(response *schema.Instance, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newInstanceStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateInstanceStub(response *schema.Instance, url string, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setInstanceState(response.Status, schema.ResourceStateUpdating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureInstanceOperationStub(response *schema.Instance, url string, params HasParams) error {
	response.Metadata.Verb = http.MethodPost
	if err := configurator.configurePostStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveInstanceStub(response *schema.Instance, url string, params HasParams) error {
	setInstanceState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetSuspendedInstanceStub(response *schema.Instance, url string, params HasParams) error {
	setInstanceState(response.Status, schema.ResourceStateSuspended)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}
