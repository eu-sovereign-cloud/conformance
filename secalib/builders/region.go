package builders

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Region

type RegionBuilder struct {
	*resourceBuilder
	metadata *GlobalResourceMetadataBuilder
	spec     *schema.RegionSpec
}

func NewRegionBuilder() *RegionBuilder {
	return &RegionBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewGlobalResourceMetadataBuilder(),
		spec:            &schema.RegionSpec{},
	}
}

func (builder *RegionBuilder) Name(name string) *RegionBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *RegionBuilder) Provider(provider string) *RegionBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *RegionBuilder) Resource(resource string) *RegionBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *RegionBuilder) ApiVersion(apiVersion string) *RegionBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *RegionBuilder) Spec(spec *schema.RegionSpec) *RegionBuilder {
	builder.spec = spec
	return builder
}

func (builder *RegionBuilder) BuildResponse() (*schema.Region, error) {
	medatata, err := builder.metadata.Kind(secalib.RegionKind).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequired(
		builder.spec,
		builder.spec.AvailableZones,
		builder.spec.Providers,
	); err != nil {
		return nil, err
	}
	// Validate each provider
	for _, provider := range builder.spec.Providers {
		if err := builder.validator.ValidateRequired(
			provider.Name,
			provider.Version,
			provider.Url,
		); err != nil {
			return nil, err
		}
	}

	return &schema.Region{
		Metadata: medatata,
		Spec:     *builder.spec,
	}, nil
}
