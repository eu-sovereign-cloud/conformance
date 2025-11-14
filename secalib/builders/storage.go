package builders

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func newBlockStorageResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.BlockStorageSpec) (*schema.BlockStorage, error) {
	medatata, err := NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.BlockStorageKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	return &schema.BlockStorage{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.BlockStorageStatus{},
	}, nil
}

func newImageResponse(name, provider, resource, apiVersion, tenant, region string, spec *schema.ImageSpec) (*schema.Image, error) {
	medatata, err := NewRegionalResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.ImageKind).
		Tenant(tenant).
		Region(region).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	return &schema.Image{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.ImageStatus{},
	}, nil
}
