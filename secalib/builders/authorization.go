package builders

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Role

func newRoleResponse(name, provider, resource, apiVersion, tenant string, spec *schema.RoleSpec) (*schema.Role, error) {
	medatata, err := NewGlobalTenantResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.RoleKind).
		Tenant(tenant).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	return &schema.Role{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.RoleStatus{},
	}, nil
}

// RoleAssignment

func newRoleAssignmentResponse(name, provider, resource, apiVersion, tenant string, spec *schema.RoleAssignmentSpec) (*schema.RoleAssignment, error) {
	medatata, err := NewGlobalTenantResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.RoleAssignmentKind).
		Tenant(tenant).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	return &schema.RoleAssignment{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.RoleAssignmentStatus{},
	}, nil
}
