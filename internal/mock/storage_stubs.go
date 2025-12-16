//nolint:dupl
package mock

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Block storage

func (configurator *scenarioConfigurator) configureCreateBlockStorageStub(response *schema.BlockStorage, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newBlockStorageStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateBlockStorageStub(response *schema.BlockStorage, url string, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setBlockStorageState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveBlockStorageStub(response *schema.BlockStorage, url string, params HasParams) error {
	setBlockStorageState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Image

func (configurator *scenarioConfigurator) configureCreateImageStub(response *schema.Image, url string, params HasParams) error {
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = newImageStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateImageStub(response *schema.Image, url string, params HasParams) error {
	setModifiedRegionalResourceMetadata(response.Metadata)
	setImageState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveImageStub(response *schema.Image, url string, params HasParams) error {
	setImageState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
