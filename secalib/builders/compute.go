package builders

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func newInstanceResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.InstanceSpec) (*schema.Instance, error) {
	medatata, err := NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.InstanceKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	return &schema.Instance{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.InstanceStatus{},
	}, nil
}
