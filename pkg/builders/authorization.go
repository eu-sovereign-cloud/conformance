//nolint:dupl
package builders

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Role

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

// RoleAssignment

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
