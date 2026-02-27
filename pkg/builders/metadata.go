//nolint:dupl
package builders

import (
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/go-playground/validator/v10"
)

// ResponseMetadata

/// responseMetadataBuilder

type responseMetadataBuilder[B any] struct {
	validator *validator.Validate
	parent    *B
	metadata  *schema.ResponseMetadata
}

func newResponseMetadataBuilder[B any](parent *B) *responseMetadataBuilder[B] {
	return &responseMetadataBuilder[B]{
		validator: validator.New(),
		parent:    parent,
		metadata: &schema.ResponseMetadata{
			Verb: http.MethodGet,
		},
	}
}

func (builder *responseMetadataBuilder[B]) Provider(provider string) *B {
	builder.metadata.Provider = provider
	return builder.parent
}

func (builder *responseMetadataBuilder[B]) SkipToken(skipToken string) *B {
	builder.metadata.SkipToken = &skipToken
	return builder.parent
}

func (builder *responseMetadataBuilder[B]) validate() error {
	if err := validateRequired(builder.validator,
		field("metadata", builder.metadata),
		field("metadata.Provider", builder.metadata.Provider),
		field("metadata.Verb", builder.metadata.Verb),
	); err != nil {
		return err
	}

	return nil
}

// tenantResponseMetadataBuilder

type tenantResponseMetadataBuilder[B any] struct {
	*responseMetadataBuilder[B]

	tenant string
}

func newTenantResponseMetadataBuilder[B any](parent *B) *tenantResponseMetadataBuilder[B] {
	return &tenantResponseMetadataBuilder[B]{
		responseMetadataBuilder: newResponseMetadataBuilder(parent),
	}
}

func (builder *tenantResponseMetadataBuilder[B]) Tenant(tenant string) *B {
	builder.tenant = tenant
	return builder.parent
}

func (builder *tenantResponseMetadataBuilder[B]) validate() error {
	// Validate response metadata
	if err := builder.responseMetadataBuilder.validate(); err != nil {
		return err
	}

	// Validate tenant response metadata
	if err := validateRequired(builder.validator,
		field("tenant", builder.tenant),
	); err != nil {
		return err
	}

	return nil
}

// workspaceResponseMetadataBuilder

type workspaceResponseMetadataBuilder[B any] struct {
	*responseMetadataBuilder[B]

	tenant    string
	workspace string
}

func newWorkspaceResponseMetadataBuilder[B any](parent *B) *workspaceResponseMetadataBuilder[B] {
	return &workspaceResponseMetadataBuilder[B]{
		responseMetadataBuilder: newResponseMetadataBuilder(parent),
	}
}

func (builder *workspaceResponseMetadataBuilder[B]) Tenant(tenant string) *B {
	builder.tenant = tenant
	return builder.parent
}

func (builder *workspaceResponseMetadataBuilder[B]) Workspace(workspace string) *B {
	builder.workspace = workspace
	return builder.parent
}

func (builder *workspaceResponseMetadataBuilder[B]) validate() error {
	// Validate response metadata
	if err := builder.responseMetadataBuilder.validate(); err != nil {
		return err
	}

	// Validate workspace response metadata
	if err := validateRequired(builder.validator,
		field("tenant", builder.tenant),
		field("workspace", builder.workspace),
	); err != nil {
		return err
	}

	return nil
}

// networkResponseMetadataBuilder

type networkResponseMetadataBuilder[B any] struct {
	*responseMetadataBuilder[B]

	tenant    string
	workspace string
	network   string
}

func newNetworkResponseMetadataBuilder[B any](parent *B) *networkResponseMetadataBuilder[B] {
	return &networkResponseMetadataBuilder[B]{
		responseMetadataBuilder: newResponseMetadataBuilder(parent),
	}
}

func (builder *networkResponseMetadataBuilder[B]) Tenant(tenant string) *B {
	builder.tenant = tenant
	return builder.parent
}

func (builder *networkResponseMetadataBuilder[B]) Workspace(workspace string) *B {
	builder.workspace = workspace
	return builder.parent
}

func (builder *networkResponseMetadataBuilder[B]) Network(network string) *B {
	builder.network = network
	return builder.parent
}

func (builder *networkResponseMetadataBuilder[B]) validate() error {
	// Validate response metadata
	if err := builder.responseMetadataBuilder.validate(); err != nil {
		return err
	}

	// Validate network response metadata
	if err := validateRequired(builder.validator,
		field("tenant", builder.tenant),
		field("workspace", builder.workspace),
		field("network", builder.network),
	); err != nil {
		return err
	}

	return nil
}

// ResourceMetadata

/// GlobalResourceMetadataBuilder

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
		setProvider:   func(resource string) { builder.metadata.Resource = resource },
		setResource:   func(provider string) { builder.metadata.Provider = provider },
		setApiVersion: func(apiVersion string) { builder.metadata.ApiVersion = apiVersion },
		setKind:       func(kind schema.GlobalResourceMetadataKind) { builder.metadata.Kind = kind },
		setRef:        func(ref string) { builder.metadata.Ref = ref },
	})

	return builder
}

func (builder *globalResourceMetadataBuilder[B]) build() (*schema.GlobalResourceMetadata, error) {
	if err := validateRequired(builder.validator,
		field("metadata", builder.metadata),
		field("metadata.Name", builder.metadata.Name),
		field("metadata.Provider", builder.metadata.Provider),
		field("metadata.Resource", builder.metadata.Resource),
		field("metadata.ApiVersion", builder.metadata.ApiVersion),
		field("metadata.Kind", builder.metadata.Kind),
		field("metadata.Ref", builder.metadata.Ref),
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

/// GlobalTenantResourceMetadataBuilder

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
		setResource:   func(provider string) { builder.metadata.Provider = provider },
		setApiVersion: func(apiVersion string) { builder.metadata.ApiVersion = apiVersion },
		setKind:       func(kind schema.GlobalTenantResourceMetadataKind) { builder.metadata.Kind = kind },
		setRef:        func(ref string) { builder.metadata.Ref = ref },
	})

	return builder
}

func (builder *globalTenantResourceMetadataBuilder[B]) Tenant(tenant string) *B {
	builder.metadata.Tenant = tenant
	return builder.resourceMetadataBuilder.parent
}

func (builder *globalTenantResourceMetadataBuilder[B]) build() (*schema.GlobalTenantResourceMetadata, error) {
	if err := validateRequired(builder.validator,
		field("metadata", builder.metadata),
		field("metadata.Name", builder.metadata.Name),
		field("metadata.Provider", builder.metadata.Provider),
		field("metadata.ApiVersion", builder.metadata.ApiVersion),
		field("metadata.Kind", builder.metadata.Kind),
		field("metadata.Ref", builder.metadata.Ref),
		field("metadata.Tenant", builder.metadata.Tenant),
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

/// RegionalResourceMetadata

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
		setResource:   func(provider string) { builder.metadata.Provider = provider },
		setApiVersion: func(apiVersion string) { builder.metadata.ApiVersion = apiVersion },
		setKind:       func(kind schema.RegionalResourceMetadataKind) { builder.metadata.Kind = kind },
		setRef:        func(ref string) { builder.metadata.Ref = ref },
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
		field("metadata", builder.metadata),
		field("metadata.Name", builder.metadata.Name),
		field("metadata.Provider", builder.metadata.Provider),
		field("metadata.ApiVersion", builder.metadata.ApiVersion),
		field("metadata.Kind", builder.metadata.Kind),
		field("metadata.Ref", builder.metadata.Ref),
		field("metadata.Tenant", builder.metadata.Tenant),
		field("metadata.Region", builder.metadata.Region),
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

/// RegionalWorkspaceResourceMetadata

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
		setResource:   func(provider string) { builder.metadata.Provider = provider },
		setApiVersion: func(apiVersion string) { builder.metadata.ApiVersion = apiVersion },
		setKind:       func(kind schema.RegionalWorkspaceResourceMetadataKind) { builder.metadata.Kind = kind },
		setRef:        func(ref string) { builder.metadata.Ref = ref },
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
		field("metadata", builder.metadata),
		field("metadata.Name", builder.metadata.Name),
		field("metadata.Provider", builder.metadata.Provider),
		field("metadata.ApiVersion", builder.metadata.ApiVersion),
		field("metadata.Kind", builder.metadata.Kind),
		field("metadata.Ref", builder.metadata.Ref),
		field("metadata.Tenant", builder.metadata.Tenant),
		field("metadata.Workspace", builder.metadata.Workspace),
		field("metadata.Region", builder.metadata.Region),
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

/// RegionalNetworkResourceMetadata

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
		setResource:   func(provider string) { builder.metadata.Provider = provider },
		setApiVersion: func(apiVersion string) { builder.metadata.ApiVersion = apiVersion },
		setKind:       func(kind schema.RegionalNetworkResourceMetadataKind) { builder.metadata.Kind = kind },
		setRef:        func(ref string) { builder.metadata.Ref = ref },
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
		field("metadata", builder.metadata),
		field("metadata.Name", builder.metadata.Name),
		field("metadata.Provider", builder.metadata.Provider),
		field("metadata.ApiVersion", builder.metadata.ApiVersion),
		field("metadata.Kind", builder.metadata.Kind),
		field("metadata.Ref", builder.metadata.Ref),
		field("metadata.Tenant", builder.metadata.Tenant),
		field("metadata.Workspace", builder.metadata.Workspace),
		field("metadata.Network", builder.metadata.Network),
		field("metadata.Region", builder.metadata.Region),
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}
