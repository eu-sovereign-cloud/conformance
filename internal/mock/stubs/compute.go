package stubs

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	compute "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.compute.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Instance

func (configurator *scenarioConfigurator) ConfigureCreateInstanceStub(response *schema.Instance, url string, params *mock.BaseParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newInstanceStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateInstanceStub(response *schema.Instance, url string, params *mock.BaseParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setInstanceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureInstanceOperationStub(response *schema.Instance, url string, params *mock.BaseParams) error {
	if err := configurator.ConfigurePostStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveInstanceStub(response *schema.Instance, url string, params *mock.BaseParams) error {
	setInstanceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetSuspendedInstanceStub(response *schema.Instance, url string, params *mock.BaseParams) error {
	setInstanceState(response.Status, schema.ResourceStateSuspended)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListInstanceStub(response *compute.InstanceIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListSkuStub(response *compute.SkuIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
