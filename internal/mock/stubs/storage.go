//nolint:dupl
package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Block storage

func (configurator *scenarioConfigurator) ConfigureCreateBlockStorageStub(response *schema.BlockStorage, url string, params *mock.BaseParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newBlockStorageStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateBlockStorageStub(response *schema.BlockStorage, url string, params *mock.BaseParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setBlockStorageState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveBlockStorageStub(response *schema.BlockStorage, url string, params *mock.BaseParams) error {
	setBlockStorageState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListBlockStorageStub(response storage.BlockStorageIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Image

func (configurator *scenarioConfigurator) ConfigureCreateImageStub(response *schema.Image, url string, params *mock.BaseParams) error {
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = newImageStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateImageStub(response *schema.Image, url string, params *mock.BaseParams) error {
	setModifiedRegionalResourceMetadata(response.Metadata)
	setImageState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveImageStub(response *schema.Image, url string, params *mock.BaseParams) error {
	setImageState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListImageStub(response storage.ImageIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListStorageSkuStub(response storage.SkuIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
