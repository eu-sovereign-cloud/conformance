package wrappers

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

type RoleWrapper struct {
	*BaseResourceWrapper[schema.Role]
}

func NewRoleWrapper(resource *schema.Role) *RoleWrapper {
	return &RoleWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *RoleWrapper) GetResource() *schema.Role {
	return wrapper.resource
}

func (wrapper *RoleWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *RoleWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *RoleWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *RoleWrapper) GetMetadata() *schema.GlobalTenantResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *RoleWrapper) GetSpec() *schema.RoleSpec {
	return &wrapper.resource.Spec
}

func (wrapper *RoleWrapper) GetStatus() *schema.RoleStatus {
	return wrapper.resource.Status
}

type RoleAssignmentWrapper struct {
	*BaseResourceWrapper[schema.RoleAssignment]
}

func NewRoleAssignmentWrapper(resource *schema.RoleAssignment) *RoleAssignmentWrapper {
	return &RoleAssignmentWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *RoleAssignmentWrapper) GetResource() *schema.RoleAssignment {
	return wrapper.resource
}

func (wrapper *RoleAssignmentWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *RoleAssignmentWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *RoleAssignmentWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *RoleAssignmentWrapper) GetMetadata() *schema.GlobalTenantResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *RoleAssignmentWrapper) GetSpec() *schema.RoleAssignmentSpec {
	return &wrapper.resource.Spec
}

func (wrapper *RoleAssignmentWrapper) GetStatus() *schema.RoleAssignmentStatus {
	return wrapper.resource.Status
}
