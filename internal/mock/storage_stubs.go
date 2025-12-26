//nolint:dupl
package mock

import (
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Block storage

func (configurator *scenarioConfigurator) configureCreateBlockStorageStub(response *schema.BlockStorage, url string, params *BaseParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newBlockStorageStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateBlockStorageStub(response *schema.BlockStorage, url string, params *BaseParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setBlockStorageState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveBlockStorageStub(response *schema.BlockStorage, url string, params *BaseParams) error {
	setBlockStorageState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListBlockStorageStub(response storage.BlockStorageIterator, url string, params *BaseParams, pathParams map[string]string) error {
	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Image

func (configurator *scenarioConfigurator) configureCreateImageStub(response *schema.Image, url string, params *BaseParams) error {
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = newImageStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateImageStub(response *schema.Image, url string, params *BaseParams) error {
	setModifiedRegionalResourceMetadata(response.Metadata)
	setImageState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveImageStub(response *schema.Image, url string, params *BaseParams) error {
	setImageState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListImageStub(response storage.ImageIterator, url string, params *BaseParams, pathParams map[string]string) error {
	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListStorageSkuStub(response storage.SkuIterator, url string, params *BaseParams, pathParams map[string]string) error {
	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
