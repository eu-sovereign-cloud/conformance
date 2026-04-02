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
	response.Status = newBlockStorageStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateBlockStorageStub(response *schema.BlockStorage, url string, params *mock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = beforeUpdateBlockStorageStatus()
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingBlockStorageStub(response *schema.BlockStorage, url string, params *mock.MockParams) error {
	setBlockStorageState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
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

func (configurator *Configurator) ConfigureGetUpdatingBlockStorageStub(response *schema.BlockStorage, url string, params *mock.MockParams) error {
	setBlockStorageState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetDeletingBlockStorageStub(response *schema.BlockStorage, url string, params *mock.MockParams) error {
	setBlockStorageState(response.Status, schema.ResourceStateDeleting)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListBlockStorageStub(response storage.BlockStorageIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Image

func (configurator *Configurator) ConfigureCreateImageStub(response *schema.Image, url string, params *mock.MockParams) error {
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = newImageStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateImageStub(response *schema.Image, url string, params *mock.MockParams) error {
	setModifiedRegionalResourceMetadata(response.Metadata)
	response.Status = beforeUpdateImageStatus()
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingImageStub(response *schema.Image, url string, params *mock.MockParams) error {
	setImageState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
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

func (configurator *Configurator) ConfigureGetUpdatingImageStub(response *schema.Image, url string, params *mock.MockParams) error {
	setImageState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetDeletingImageStub(response *schema.Image, url string, params *mock.MockParams) error {
	setImageState(response.Status, schema.ResourceStateDeleting)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListImageStub(response storage.ImageIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListStorageSkuStub(response storage.SkuIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
