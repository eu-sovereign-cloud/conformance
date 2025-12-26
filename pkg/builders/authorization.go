//nolint:dupl
package builders

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	authorization "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.authorization.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Role

/// RoleMetadataBuilder

type RoleMetadataBuilder struct {
	*globalTenantResourceMetadataBuilder[RoleMetadataBuilder]
}

func NewRoleMetadataBuilder() *RoleMetadataBuilder {
	builder := &RoleMetadataBuilder{}
	builder.globalTenantResourceMetadataBuilder = newGlobalTenantResourceMetadataBuilder(builder)
	return builder
}

func (builder *RoleMetadataBuilder) Build() (*schema.GlobalTenantResourceMetadata, error) {
	metadata, err := builder.kind(schema.GlobalTenantResourceMetadataKindResourceKindRole).build()
	if err != nil {
		return nil, err
	}

	resource := generators.GenerateRoleResource(builder.metadata.Tenant, builder.metadata.Name)
	metadata.Resource = resource

	return metadata, nil
}

/// RoleBuilder

type RoleBuilder struct {
	*globalTenantResourceBuilder[RoleBuilder, schema.RoleSpec]
	metadata *RoleMetadataBuilder
	labels   schema.Labels
	spec     *schema.RoleSpec
}

func NewRoleBuilder() *RoleBuilder {
	builder := &RoleBuilder{
		metadata: NewRoleMetadataBuilder(),
		spec:     &schema.RoleSpec{},
	}

	builder.globalTenantResourceBuilder = newGlobalTenantResourceBuilder(newGlobalTenantResourceBuilderParams[RoleBuilder, schema.RoleSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[RoleBuilder, schema.RoleSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.RoleSpec) { builder.spec = spec },
		},
		setTenant: func(tenant string) { builder.metadata.Tenant(tenant) },
	})

	return builder
}

func (builder *RoleBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		builder.spec,
		builder.spec.Permissions,
	); err != nil {
		return err
	}

	// Validate each permission
	for _, permission := range builder.spec.Permissions {
		if err := validateRequired(builder.validator,
			permission.Provider,
			permission.Resources,
			permission.Verb,
		); err != nil {
			return err
		}
	}

	return nil
}

func (builder *RoleBuilder) Build() (*schema.Role, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.Role{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.RoleStatus{},
	}, nil
}

/// RoleIteratorBuilder

type RoleIteratorBuilder struct {
	*tenantResponseMetadataBuilder[RoleIteratorBuilder]

	items []schema.Role
}

func NewRoleIteratorBuilder() *RoleIteratorBuilder {
	builder := &RoleIteratorBuilder{}
	builder.tenantResponseMetadataBuilder = newTenantResponseMetadataBuilder(builder)
	return builder
}

func (builder *RoleIteratorBuilder) Items(items []schema.Role) *RoleIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *RoleIteratorBuilder) Build() (*authorization.RoleIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateRoleListResource(builder.tenant)

	return &authorization.RoleIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}

// RoleAssignment

/// RoleAssignmentMetadataBuilder

type RoleAssignmentMetadataBuilder struct {
	*globalTenantResourceMetadataBuilder[RoleAssignmentMetadataBuilder]
}

func NewRoleAssignmentMetadataBuilder() *RoleAssignmentMetadataBuilder {
	builder := &RoleAssignmentMetadataBuilder{}
	builder.globalTenantResourceMetadataBuilder = newGlobalTenantResourceMetadataBuilder(builder)
	return builder
}

func (builder *RoleAssignmentMetadataBuilder) Build() (*schema.GlobalTenantResourceMetadata, error) {
	metadata, err := builder.kind(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment).build()
	if err != nil {
		return nil, err
	}

	resource := generators.GenerateRoleAssignmentResource(builder.metadata.Tenant, builder.metadata.Name)
	metadata.Resource = resource

	return metadata, nil
}

/// RoleAssignmentBuilder

type RoleAssignmentBuilder struct {
	*globalTenantResourceBuilder[RoleAssignmentBuilder, schema.RoleAssignmentSpec]
	metadata *RoleAssignmentMetadataBuilder
	labels   schema.Labels
	spec     *schema.RoleAssignmentSpec
}

func NewRoleAssignmentBuilder() *RoleAssignmentBuilder {
	builder := &RoleAssignmentBuilder{
		metadata: NewRoleAssignmentMetadataBuilder(),
		spec:     &schema.RoleAssignmentSpec{},
	}

	builder.globalTenantResourceBuilder = newGlobalTenantResourceBuilder(newGlobalTenantResourceBuilderParams[RoleAssignmentBuilder, schema.RoleAssignmentSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[RoleAssignmentBuilder, schema.RoleAssignmentSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.RoleAssignmentSpec) { builder.spec = spec },
		},
		setTenant: func(tenant string) { builder.metadata.Tenant(tenant) },
	})

	return builder
}

func (builder *RoleAssignmentBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		builder.spec,
		builder.spec.Subs,
		builder.spec.Scopes,
		builder.spec.Roles,
	); err != nil {
		return err
	}

	// Validate each scope
	for _, scope := range builder.spec.Scopes {
		if err := validateOneRequired(builder.validator,
			scope.Tenants,
			scope.Workspaces,
			scope.Regions,
		); err != nil {
			return err
		}
	}

	return nil
}

func (builder *RoleAssignmentBuilder) Build() (*schema.RoleAssignment, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.RoleAssignment{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.RoleAssignmentStatus{},
	}, nil
}

/// RoleAssignmentIteratorBuilder

type RoleAssignmentIteratorBuilder struct {
	*tenantResponseMetadataBuilder[RoleAssignmentIteratorBuilder]

	items []schema.RoleAssignment
}

func NewRoleAssignmentIteratorBuilder() *RoleAssignmentIteratorBuilder {
	builder := &RoleAssignmentIteratorBuilder{}
	builder.tenantResponseMetadataBuilder = newTenantResponseMetadataBuilder(builder)
	return builder
}

func (builder *RoleAssignmentIteratorBuilder) Items(items []schema.RoleAssignment) *RoleAssignmentIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *RoleAssignmentIteratorBuilder) Build() (*authorization.RoleAssignmentIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateRoleAssignmentListResource(builder.tenant)

	return &authorization.RoleAssignmentIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}
