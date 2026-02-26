//nolint:dupl
package builders

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// StorageSku

/// StorageSkuIteratorBuilder

type StorageSkuIteratorBuilder struct {
	*tenantResponseMetadataBuilder[StorageSkuIteratorBuilder]

	items []schema.StorageSku
}

func NewStorageSkuIteratorBuilder() *StorageSkuIteratorBuilder {
	builder := &StorageSkuIteratorBuilder{}
	builder.tenantResponseMetadataBuilder = newTenantResponseMetadataBuilder(builder)
	return builder
}

func (builder *StorageSkuIteratorBuilder) Items(items []schema.StorageSku) *StorageSkuIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *StorageSkuIteratorBuilder) Build() (*storage.SkuIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateSkuListResource(builder.tenant)

	return &storage.SkuIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}

// BlockStorage

/// BlockStorageMetadataBuilder

type BlockStorageMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[BlockStorageMetadataBuilder]
}

func NewBlockStorageMetadataBuilder() *BlockStorageMetadataBuilder {
	builder := &BlockStorageMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *BlockStorageMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage).
		Resource(generators.GenerateBlockStorageResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)).
		Ref(generators.GenerateBlockStorageRef(builder.metadata.Name)).
		build()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

/// BlockStorageBuilder

type BlockStorageBuilder struct {
	*regionalWorkspaceResourceBuilder[BlockStorageBuilder, schema.BlockStorageSpec]
	metadata *BlockStorageMetadataBuilder
	labels   schema.Labels
	spec     *schema.BlockStorageSpec
}

func NewBlockStorageBuilder() *BlockStorageBuilder {
	builder := &BlockStorageBuilder{
		metadata: NewBlockStorageMetadataBuilder(),
		spec:     &schema.BlockStorageSpec{},
	}

	builder.regionalWorkspaceResourceBuilder = newRegionalWorkspaceResourceBuilder(newRegionalWorkspaceResourceBuilderParams[BlockStorageBuilder, schema.BlockStorageSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[BlockStorageBuilder, schema.BlockStorageSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.BlockStorageSpec) { builder.spec = spec },
		},
		setTenant:    func(tenant string) { builder.metadata.Tenant(tenant) },
		setWorkspace: func(workspace string) { builder.metadata.Workspace(workspace) },
		setRegion:    func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *BlockStorageBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		field("spec", builder.spec),
		field("spec.SkuRef", builder.spec.SkuRef),
		field("spec.SizeGB", builder.spec.SizeGB),
	); err != nil {
		return err
	}

	return nil
}

func (builder *BlockStorageBuilder) Build() (*schema.BlockStorage, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.BlockStorage{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.BlockStorageStatus{},
	}, nil
}

/// BlockStorageIteratorBuilder

type BlockStorageIteratorBuilder struct {
	*workspaceResponseMetadataBuilder[BlockStorageIteratorBuilder]

	items []schema.BlockStorage
}

func NewBlockStorageIteratorBuilder() *BlockStorageIteratorBuilder {
	builder := &BlockStorageIteratorBuilder{}
	builder.workspaceResponseMetadataBuilder = newWorkspaceResponseMetadataBuilder(builder)
	return builder
}

func (builder *BlockStorageIteratorBuilder) Items(items []schema.BlockStorage) *BlockStorageIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *BlockStorageIteratorBuilder) Build() (*storage.BlockStorageIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateBlockStorageListResource(builder.tenant, builder.workspace)

	return &storage.BlockStorageIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}

// Image

/// ImageMetadataBuilder

type ImageMetadataBuilder struct {
	*regionalResourceMetadataBuilder[ImageMetadataBuilder]
}

func NewImageMetadataBuilder() *ImageMetadataBuilder {
	builder := &ImageMetadataBuilder{}
	builder.regionalResourceMetadataBuilder = newRegionalResourceMetadataBuilder(builder)
	return builder
}

func (builder *ImageMetadataBuilder) Build() (*schema.RegionalResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalResourceMetadataKindResourceKindImage).
		Resource(generators.GenerateImageResource(builder.metadata.Tenant, builder.metadata.Name)).
		Ref(generators.GenerateImageRef(builder.metadata.Name)).
		build()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

/// ImageBuilder

type ImageBuilder struct {
	*regionalResourceBuilder[ImageBuilder, schema.ImageSpec]
	metadata *ImageMetadataBuilder
	labels   schema.Labels
	spec     *schema.ImageSpec
}

func NewImageBuilder() *ImageBuilder {
	builder := &ImageBuilder{
		metadata: NewImageMetadataBuilder(),
		spec:     &schema.ImageSpec{},
	}

	builder.regionalResourceBuilder = newRegionalResourceBuilder(newRegionalResourceBuilderParams[ImageBuilder, schema.ImageSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[ImageBuilder, schema.ImageSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.ImageSpec) { builder.spec = spec },
		},
		setTenant: func(tenant string) { builder.metadata.Tenant(tenant) },
		setRegion: func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *ImageBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		field("spec", builder.spec),
		field("spec.BlockStorageRef", builder.spec.BlockStorageRef),
		field("spec.CpuArchitecture", builder.spec.CpuArchitecture),
	); err != nil {
		return err
	}

	return nil
}

func (builder *ImageBuilder) Build() (*schema.Image, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.Image{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.ImageStatus{},
	}, nil
}

/// ImageIteratorBuilder

type ImageIteratorBuilder struct {
	*tenantResponseMetadataBuilder[ImageIteratorBuilder]

	items []schema.Image
}

func NewImageIteratorBuilder() *ImageIteratorBuilder {
	builder := &ImageIteratorBuilder{}
	builder.tenantResponseMetadataBuilder = newTenantResponseMetadataBuilder(builder)
	return builder
}

func (builder *ImageIteratorBuilder) Items(items []schema.Image) *ImageIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *ImageIteratorBuilder) Build() (*storage.ImageIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateImageListResource(builder.tenant)

	return &storage.ImageIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}
