package wrappers

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

type BlockStorageWrapper struct {
	*BaseResourceWrapper[schema.BlockStorage]
}

func NewBlockStorageWrapper(resource *schema.BlockStorage) *BlockStorageWrapper {
	return &BlockStorageWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *BlockStorageWrapper) GetResource() *schema.BlockStorage {
	return wrapper.resource
}

func (wrapper *BlockStorageWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *BlockStorageWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *BlockStorageWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *BlockStorageWrapper) GetMetadata() *schema.RegionalWorkspaceResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *BlockStorageWrapper) GetSpec() *schema.BlockStorageSpec {
	return &wrapper.resource.Spec
}

func (wrapper *BlockStorageWrapper) GetStatus() *schema.BlockStorageStatus {
	return wrapper.resource.Status
}

type ImageWrapper struct {
	*BaseResourceWrapper[schema.Image]
}

func NewImageWrapper(resource *schema.Image) *ImageWrapper {
	return &ImageWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *ImageWrapper) GetResource() *schema.Image {
	return wrapper.resource
}

func (wrapper *ImageWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *ImageWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *ImageWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *ImageWrapper) GetMetadata() *schema.RegionalResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *ImageWrapper) GetSpec() *schema.ImageSpec {
	return &wrapper.resource.Spec
}

func (wrapper *ImageWrapper) GetStatus() *schema.ImageStatus {
	return wrapper.resource.Status
}
