package builders

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// workspace

type WorkspaceBuilder struct {
	*resourceBuilder
	labels   schema.Labels
	metadata *RegionalResourceMetadataBuilder
}

func NewWorkspaceBuilder() *WorkspaceBuilder {
	return &WorkspaceBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewRegionalResourceMetadataBuilder(),
	}
}
func (builder *WorkspaceBuilder) Labels(labels schema.Labels) *WorkspaceBuilder {
	builder.labels = labels
	return builder
}
func (builder *WorkspaceBuilder) Name(name string) *WorkspaceBuilder {
	builder.metadata.Name(name)
	return builder
}
func (builder *WorkspaceBuilder) Provider(provider string) *WorkspaceBuilder {
	builder.metadata.Provider(provider)
	return builder
}
func (builder *WorkspaceBuilder) Resource(resource string) *WorkspaceBuilder {
	builder.metadata.Resource(resource)
	return builder
}
func (builder *WorkspaceBuilder) ApiVersion(apiVersion string) *WorkspaceBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}
func (builder *WorkspaceBuilder) Tenant(tenant string) *WorkspaceBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}
func (builder *WorkspaceBuilder) Region(region string) *WorkspaceBuilder {
	builder.metadata.Region(region)
	return builder
}
func (builder *WorkspaceBuilder) BuildResponse() (*schema.Workspace, error) {

	medatata, err := builder.metadata.Kind(secalib.WorkspaceKind).BuildResponse()
	if err != nil {
		return nil, err
	}

	return &schema.Workspace{
		Metadata: medatata,
		Labels:   builder.labels,
		Spec:     schema.WorkspaceSpec{},
		Status:   &schema.WorkspaceStatus{},
	}, nil
}
