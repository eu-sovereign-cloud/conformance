//nolint:dupl
package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Block storage

func (configurator *Configurator) ConfigureCreateBlockStorageStub(response *schema.BlockStorage, url string, params *mock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newBlockStorageStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateBlockStorageStub(response *schema.BlockStorage, url string, params *mock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setBlockStorageState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveBlockStorageStub(response *schema.BlockStorage, url string, params *mock.MockParams) error {
	setBlockStorageState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetListBlockStorageStub(response storage.BlockStorageIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Image

func (configurator *Configurator) ConfigureCreateImageStub(response *schema.Image, url string, params *mock.MockParams) error {
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = newImageStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateImageStub(response *schema.Image, url string, params *mock.MockParams) error {
	setModifiedRegionalResourceMetadata(response.Metadata)
	setImageState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveImageStub(response *schema.Image, url string, params *mock.MockParams) error {
	setImageState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetListImageStub(response storage.ImageIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetListStorageSkuStub(response storage.SkuIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
