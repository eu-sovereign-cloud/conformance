package builders

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// BlockStorage

type BlockStorageBuilder struct {
	*resourceBuilder
	metadata *RegionalWorkspaceResourceMetadataBuilder
	spec     *schema.BlockStorageSpec
}

func NewBlockStorageBuilder() *BlockStorageBuilder {
	return &BlockStorageBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewRegionalWorkspaceResourceMetadataBuilder(),
		spec:            &schema.BlockStorageSpec{},
	}
}

func (builder *BlockStorageBuilder) Name(name string) *BlockStorageBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *BlockStorageBuilder) Provider(provider string) *BlockStorageBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *BlockStorageBuilder) Resource(resource string) *BlockStorageBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *BlockStorageBuilder) ApiVersion(apiVersion string) *BlockStorageBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *BlockStorageBuilder) Tenant(tenant string) *BlockStorageBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *BlockStorageBuilder) Workspace(workspace string) *BlockStorageBuilder {
	builder.metadata.Workspace(workspace)
	return builder
}

func (builder *BlockStorageBuilder) Region(region string) *BlockStorageBuilder {
	builder.metadata.Region(region)
	return builder
}

func (builder *BlockStorageBuilder) Spec(spec *schema.BlockStorageSpec) *BlockStorageBuilder {
	builder.spec = spec
	return builder
}

func (builder *BlockStorageBuilder) BuildResponse() (*schema.BlockStorage, error) {
	medatata, err := builder.metadata.Kind(secalib.BlockStorageKind).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequired(
		builder.spec,
		builder.spec.SkuRef,
		builder.spec.SizeGB,
	); err != nil {
		return nil, err
	}

	return &schema.BlockStorage{
		Metadata: medatata,
		Labels:   schema.Labels{},
		Spec:     *builder.spec,
		Status:   &schema.BlockStorageStatus{},
	}, nil
}

// Image

type ImageBuilder struct {
	*resourceBuilder
	metadata *RegionalResourceMetadataBuilder
	spec     *schema.ImageSpec
}

func NewImageBuilder() *ImageBuilder {
	return &ImageBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewRegionalResourceMetadataBuilder(),
		spec:            &schema.ImageSpec{},
	}
}

func (builder *ImageBuilder) Name(name string) *ImageBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *ImageBuilder) Provider(provider string) *ImageBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *ImageBuilder) Resource(resource string) *ImageBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *ImageBuilder) ApiVersion(apiVersion string) *ImageBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *ImageBuilder) Tenant(tenant string) *ImageBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *ImageBuilder) Region(region string) *ImageBuilder {
	builder.metadata.Region(region)
	return builder
}

func (builder *ImageBuilder) Spec(spec *schema.ImageSpec) *ImageBuilder {
	builder.spec = spec
	return builder
}

func (builder *ImageBuilder) BuildResponse() (*schema.Image, error) {
	medatata, err := builder.metadata.Kind(secalib.ImageKind).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequired(
		builder.spec,
		builder.spec.BlockStorageRef,
		builder.spec.CpuArchitecture,
	); err != nil {
		return nil, err
	}

	return &schema.Image{
		Metadata: medatata,
		Labels:   schema.Labels{},
		Spec:     *builder.spec,
		Status:   &schema.ImageStatus{},
	}, nil
}
