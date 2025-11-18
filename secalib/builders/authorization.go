package builders

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Role

type RoleBuilder struct {
	*resourceBuilder
	metadata *GlobalTenantResourceMetadataBuilder
	spec     *schema.RoleSpec
}

func NewRoleBuilder() *RoleBuilder {
	return &RoleBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewGlobalTenantResourceMetadataBuilder(),
		spec:            &schema.RoleSpec{},
	}
}

func (builder *RoleBuilder) Name(name string) *RoleBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *RoleBuilder) Provider(provider string) *RoleBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *RoleBuilder) Resource(resource string) *RoleBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *RoleBuilder) ApiVersion(apiVersion string) *RoleBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *RoleBuilder) Tenant(tenant string) *RoleBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *RoleBuilder) Spec(spec *schema.RoleSpec) *RoleBuilder {
	builder.spec = spec
	return builder
}

func (builder *RoleBuilder) BuildResponse() (*schema.Role, error) {
	medatata, err := builder.metadata.Kind(secalib.RoleKind).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequireds(
		[]any{
			builder.spec,
			builder.spec.Permissions,
		},
	); err != nil {
		return nil, err
	}
	// Validate each permission
	for _, permission := range builder.spec.Permissions {
		if err := builder.validator.ValidateRequireds(
			[]any{
				permission.Provider,
				permission.Resources,
				permission.Verb,
			},
		); err != nil {
			return nil, err
		}
	}

	return &schema.Role{
		Metadata: medatata,
		Labels:   schema.Labels{},
		Spec:     *builder.spec,
		Status:   &schema.RoleStatus{},
	}, nil
}

// RoleAssignment

type RoleAssignmentBuilder struct {
	*resourceBuilder
	metadata *GlobalTenantResourceMetadataBuilder
	spec     *schema.RoleAssignmentSpec
}

func NewRoleAssignmentBuilder() *RoleAssignmentBuilder {
	return &RoleAssignmentBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewGlobalTenantResourceMetadataBuilder(),
		spec:            &schema.RoleAssignmentSpec{},
	}
}

func (builder *RoleAssignmentBuilder) Name(name string) *RoleAssignmentBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *RoleAssignmentBuilder) Provider(provider string) *RoleAssignmentBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *RoleAssignmentBuilder) Resource(resource string) *RoleAssignmentBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *RoleAssignmentBuilder) ApiVersion(apiVersion string) *RoleAssignmentBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *RoleAssignmentBuilder) Tenant(tenant string) *RoleAssignmentBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *RoleAssignmentBuilder) Spec(spec *schema.RoleAssignmentSpec) *RoleAssignmentBuilder {
	builder.spec = spec
	return builder
}

func (builder *RoleAssignmentBuilder) BuildResponse() (*schema.RoleAssignment, error) {
	medatata, err := builder.metadata.Kind(secalib.RoleAssignmentKind).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequireds(
		[]any{
			builder.spec,
			builder.spec.Subs,
			builder.spec.Scopes,
			builder.spec.Roles,
		},
	); err != nil {
		return nil, err
	}
	// Validate each scope
	for _, scope := range builder.spec.Scopes {
		if err := builder.validator.ValidateOneRequired(
			[]any{
				scope.Tenants,
				scope.Workspaces,
				scope.Regions,
			},
		); err != nil {
			return nil, err
		}
	}
	// TODO Validate each scope, if all fields are nil

	return &schema.RoleAssignment{
		Metadata: medatata,
		Labels:   schema.Labels{},
		Spec:     *builder.spec,
		Status:   &schema.RoleAssignmentStatus{},
	}, nil
}
