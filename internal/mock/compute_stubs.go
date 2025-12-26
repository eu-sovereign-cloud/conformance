package mock

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Instance

func (configurator *scenarioConfigurator) configureCreateInstanceStub(response *schema.Instance, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newInstanceStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateInstanceStub(response *schema.Instance, url string, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setInstanceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureInstanceOperationStub(response *schema.Instance, url string, params HasParams) error {
	if err := configurator.configurePostStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveInstanceStub(response *schema.Instance, url string, params HasParams) error {
	setInstanceState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetSuspendedInstanceStub(response *schema.Instance, url string, params HasParams) error {
	setInstanceState(response.Status, schema.ResourceStateSuspended)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
