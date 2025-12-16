//nolint:dupl
package builders

import (
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/go-playground/validator/v10"
)

// ResponseMetadataBuilder

type responseMetadataBuilder struct {
	validator *validator.Validate
	metadata  *schema.ResponseMetadata
}

func NewResponseMetadataBuilder() *responseMetadataBuilder {
	return &responseMetadataBuilder{
		validator: validator.New(),
		metadata: &schema.ResponseMetadata{
			Verb: http.MethodGet,
		},
	}
}

func (builder *responseMetadataBuilder) Provider(provider string) *responseMetadataBuilder {
	builder.metadata.Provider = provider
	return builder
}

func (builder *responseMetadataBuilder) SkipToken(skipToken string) *responseMetadataBuilder {
	builder.metadata.SkipToken = &skipToken
	return builder
}

func (builder *responseMetadataBuilder) build() (*schema.ResponseMetadata, error) {
	if err := validateRequired(builder.validator,
		builder.metadata,
		builder.metadata.Provider,
		builder.metadata.Verb,
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

// GlobalResourceMetadataBuilder

type globalResourceMetadataBuilder[B any] struct {
	*resourceMetadataBuilder[B, schema.GlobalResourceMetadataKind]
	metadata *schema.GlobalResourceMetadata
}

func newGlobalResourceMetadataBuilder[B any](parent *B) *globalResourceMetadataBuilder[B] {
	builder := &globalResourceMetadataBuilder[B]{
		metadata: &schema.GlobalResourceMetadata{},
	}

	builder.resourceMetadataBuilder = newResourceMetadataBuilder(newResourceMetadataBuilderParams[B, schema.GlobalResourceMetadataKind]{
		parent:        parent,
		setName:       func(name string) { builder.metadata.Name = name },
		setProvider:   func(provider string) { builder.metadata.Provider = provider },
		setApiVersion: func(apiVersion string) { builder.metadata.ApiVersion = apiVersion },
		setKind:       func(kind schema.GlobalResourceMetadataKind) { builder.metadata.Kind = kind },
	})

	return builder
}

func (builder *globalResourceMetadataBuilder[B]) build() (*schema.GlobalResourceMetadata, error) {
	if err := validateRequired(builder.validator,
		builder.metadata,
		builder.metadata.Name,
		builder.metadata.Provider,
		builder.metadata.ApiVersion,
		builder.metadata.Kind,
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

// GlobalTenantResourceMetadataBuilder

type globalTenantResourceMetadataBuilder[B any] struct {
	*resourceMetadataBuilder[B, schema.GlobalTenantResourceMetadataKind]
	metadata *schema.GlobalTenantResourceMetadata
}

func newGlobalTenantResourceMetadataBuilder[B any](parent *B) *globalTenantResourceMetadataBuilder[B] {
	builder := &globalTenantResourceMetadataBuilder[B]{
		metadata: &schema.GlobalTenantResourceMetadata{},
	}

	builder.resourceMetadataBuilder = newResourceMetadataBuilder(newResourceMetadataBuilderParams[B, schema.GlobalTenantResourceMetadataKind]{
		parent:        parent,
		setName:       func(name string) { builder.metadata.Name = name },
		setProvider:   func(provider string) { builder.metadata.Provider = provider },
		setApiVersion: func(apiVersion string) { builder.metadata.ApiVersion = apiVersion },
		setKind:       func(kind schema.GlobalTenantResourceMetadataKind) { builder.metadata.Kind = kind },
	})

	return builder
}

func (builder *globalTenantResourceMetadataBuilder[B]) Tenant(tenant string) *B {
	builder.metadata.Tenant = tenant
	return builder.resourceMetadataBuilder.parent
}

func (builder *globalTenantResourceMetadataBuilder[B]) build() (*schema.GlobalTenantResourceMetadata, error) {
	if err := validateRequired(builder.validator,
		builder.metadata,
		builder.metadata.Name,
		builder.metadata.Provider,
		builder.metadata.ApiVersion,
		builder.metadata.Kind,
		builder.metadata.Tenant,
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

// RegionalResourceMetadata

type regionalResourceMetadataBuilder[B any] struct {
	*resourceMetadataBuilder[B, schema.RegionalResourceMetadataKind]
	metadata *schema.RegionalResourceMetadata
}

func newRegionalResourceMetadataBuilder[B any](parent *B) *regionalResourceMetadataBuilder[B] {
	builder := &regionalResourceMetadataBuilder[B]{
		metadata: &schema.RegionalResourceMetadata{},
	}

	builder.resourceMetadataBuilder = newResourceMetadataBuilder(newResourceMetadataBuilderParams[B, schema.RegionalResourceMetadataKind]{
		parent:        parent,
		setName:       func(name string) { builder.metadata.Name = name },
		setProvider:   func(provider string) { builder.metadata.Provider = provider },
		setApiVersion: func(apiVersion string) { builder.metadata.ApiVersion = apiVersion },
		setKind:       func(kind schema.RegionalResourceMetadataKind) { builder.metadata.Kind = kind },
	})

	return builder
}

func (builder *regionalResourceMetadataBuilder[B]) Tenant(tenant string) *B {
	builder.metadata.Tenant = tenant
	return builder.resourceMetadataBuilder.parent
}

func (builder *regionalResourceMetadataBuilder[B]) Region(region string) *B {
	builder.metadata.Region = region
	return builder.resourceMetadataBuilder.parent
}

func (builder *regionalResourceMetadataBuilder[B]) build() (*schema.RegionalResourceMetadata, error) {
	if err := validateRequired(builder.validator,
		builder.metadata,
		builder.metadata.Name,
		builder.metadata.Provider,
		builder.metadata.ApiVersion,
		builder.metadata.Kind,
		builder.metadata.Tenant,
		builder.metadata.Region,
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

// RegionalWorkspaceResourceMetadata

type regionalWorkspaceResourceMetadataBuilder[B any] struct {
	*resourceMetadataBuilder[B, schema.RegionalWorkspaceResourceMetadataKind]
	metadata *schema.RegionalWorkspaceResourceMetadata
}

func newRegionalWorkspaceResourceMetadataBuilder[B any](parent *B) *regionalWorkspaceResourceMetadataBuilder[B] {
	builder := &regionalWorkspaceResourceMetadataBuilder[B]{
		metadata: &schema.RegionalWorkspaceResourceMetadata{},
	}

	builder.resourceMetadataBuilder = newResourceMetadataBuilder(newResourceMetadataBuilderParams[B, schema.RegionalWorkspaceResourceMetadataKind]{
		parent:        parent,
		setName:       func(name string) { builder.metadata.Name = name },
		setProvider:   func(provider string) { builder.metadata.Provider = provider },
		setApiVersion: func(apiVersion string) { builder.metadata.ApiVersion = apiVersion },
		setKind:       func(kind schema.RegionalWorkspaceResourceMetadataKind) { builder.metadata.Kind = kind },
	})

	return builder
}

func (builder *regionalWorkspaceResourceMetadataBuilder[B]) Tenant(tenant string) *B {
	builder.metadata.Tenant = tenant
	return builder.resourceMetadataBuilder.parent
}

func (builder *regionalWorkspaceResourceMetadataBuilder[B]) Workspace(workspace string) *B {
	builder.metadata.Workspace = workspace
	return builder.resourceMetadataBuilder.parent
}

func (builder *regionalWorkspaceResourceMetadataBuilder[B]) Region(region string) *B {
	builder.metadata.Region = region
	return builder.resourceMetadataBuilder.parent
}

func (builder *regionalWorkspaceResourceMetadataBuilder[B]) build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	if err := validateRequired(builder.validator,
		builder.metadata,
		builder.metadata.Name,
		builder.metadata.Provider,
		builder.metadata.ApiVersion,
		builder.metadata.Kind,
		builder.metadata.Tenant,
		builder.metadata.Workspace,
		builder.metadata.Region,
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

// RegionalNetworkResourceMetadata

type regionalNetworkResourceMetadataBuilder[B any] struct {
	*resourceMetadataBuilder[B, schema.RegionalNetworkResourceMetadataKind]
	metadata *schema.RegionalNetworkResourceMetadata
}

func newRegionalNetworkResourceMetadataBuilder[B any](parent *B) *regionalNetworkResourceMetadataBuilder[B] {
	builder := &regionalNetworkResourceMetadataBuilder[B]{
		metadata: &schema.RegionalNetworkResourceMetadata{},
	}

	builder.resourceMetadataBuilder = newResourceMetadataBuilder(newResourceMetadataBuilderParams[B, schema.RegionalNetworkResourceMetadataKind]{
		parent:        parent,
		setName:       func(name string) { builder.metadata.Name = name },
		setProvider:   func(provider string) { builder.metadata.Provider = provider },
		setApiVersion: func(apiVersion string) { builder.metadata.ApiVersion = apiVersion },
		setKind:       func(kind schema.RegionalNetworkResourceMetadataKind) { builder.metadata.Kind = kind },
	})

	return builder
}

func (builder *regionalNetworkResourceMetadataBuilder[B]) Tenant(tenant string) *B {
	builder.metadata.Tenant = tenant
	return builder.resourceMetadataBuilder.parent
}

func (builder *regionalNetworkResourceMetadataBuilder[B]) Workspace(workspace string) *B {
	builder.metadata.Workspace = workspace
	return builder.resourceMetadataBuilder.parent
}

func (builder *regionalNetworkResourceMetadataBuilder[B]) Network(network string) *B {
	builder.metadata.Network = network
	return builder.resourceMetadataBuilder.parent
}

func (builder *regionalNetworkResourceMetadataBuilder[B]) Region(region string) *B {
	builder.metadata.Region = region
	return builder.resourceMetadataBuilder.parent
}

func (builder *regionalNetworkResourceMetadataBuilder[B]) build() (*schema.RegionalNetworkResourceMetadata, error) {
	if err := validateRequired(builder.validator,
		builder.metadata,
		builder.metadata.Name,
		builder.metadata.Provider,
		builder.metadata.ApiVersion,
		builder.metadata.Kind,
		builder.metadata.Tenant,
		builder.metadata.Workspace,
		builder.metadata.Network,
		builder.metadata.Region,
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}
