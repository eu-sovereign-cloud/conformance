//nolint:dupl
package builders

import (
	"fmt"

	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Region

type RegionMetadataBuilder struct {
	*globalResourceMetadataBuilder[RegionMetadataBuilder]
}

func NewRegionMetadataBuilder() *RegionMetadataBuilder {
	builder := &RegionMetadataBuilder{}
	builder.globalResourceMetadataBuilder = newGlobalResourceMetadataBuilder(builder)
	return builder
}

func (builder *RegionMetadataBuilder) Build() (*schema.GlobalResourceMetadata, error) {
	metadata, err := builder.kind(schema.GlobalResourceMetadataKindResourceKindRegion).
		Resource(generators.GenerateRegionResource(builder.metadata.Name)).
		Ref(generators.GenerateRegionRef(builder.metadata.Name)).
		build()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

type RegionBuilder struct {
	*globalResourceBuilder[RegionBuilder, schema.RegionSpec]
	metadata *RegionMetadataBuilder
	spec     *schema.RegionSpec
}

func NewRegionBuilder() *RegionBuilder {
	builder := &RegionBuilder{
		metadata: NewRegionMetadataBuilder(),
		spec:     &schema.RegionSpec{},
	}

	builder.globalResourceBuilder = newGlobalResourceBuilder(newGlobalResourceBuilderParams[RegionBuilder, schema.RegionSpec]{
		parent:        builder,
		setName:       func(name string) { builder.metadata.setName(name) },
		setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
		setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
		setSpec:       func(spec *schema.RegionSpec) { builder.spec = spec },
	})

	return builder
}

func (builder *RegionBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		field("spec", builder.spec),
		field("spec.AvailableZones", builder.spec.AvailableZones),
		field("spec.Providers", builder.spec.Providers),
	); err != nil {
		return err
	}

	// Validate each provider
	for i, provider := range builder.spec.Providers {
		if err := validateRequired(builder.validator,
			field(fmt.Sprintf("spec.Providers[%d].Name", i), provider.Name),
			field(fmt.Sprintf("spec.Providers[%d].Version", i), provider.Version),
			field(fmt.Sprintf("spec.Providers[%d].Url", i), provider.Url),
		); err != nil {
			return err
		}
	}

	return nil
}

func (builder *RegionBuilder) Build() (*schema.Region, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.Region{
		Metadata: metadata,
		Spec:     *builder.spec,
	}, nil
}
