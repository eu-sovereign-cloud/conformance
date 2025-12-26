package mock

import (
	"net/http"

	compute "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.compute.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Instance

func (configurator *scenarioConfigurator) configureCreateInstanceStub(response *schema.Instance, url string, params *BaseParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newInstanceStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateInstanceStub(response *schema.Instance, url string, params *BaseParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setInstanceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureInstanceOperationStub(response *schema.Instance, url string, params *BaseParams) error {
	if err := configurator.configurePostStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveInstanceStub(response *schema.Instance, url string, params *BaseParams) error {
	setInstanceState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetSuspendedInstanceStub(response *schema.Instance, url string, params *BaseParams) error {
	setInstanceState(response.Status, schema.ResourceStateSuspended)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListInstanceStub(response *compute.InstanceIterator, url string, params *BaseParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListSkuStub(response *compute.SkuIterator, url string, params *BaseParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
