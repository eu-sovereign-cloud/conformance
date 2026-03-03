package wrappers

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

type WorkspaceWrapper struct {
	*BaseResourceWrapper[schema.Workspace]
}

func NewWorkspaceWrapper(resource *schema.Workspace) *WorkspaceWrapper {
	return &WorkspaceWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}
func (wrapper *WorkspaceWrapper) GetResource() *schema.Workspace {
	return wrapper.resource
}
func (wrapper *WorkspaceWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}
func (wrapper *WorkspaceWrapper) GetMetadata() *schema.RegionalResourceMetadata {
	return wrapper.resource.Metadata
}
func (wrapper *WorkspaceWrapper) GetSpec() *schema.WorkspaceSpec {
	return &wrapper.resource.Spec
}
func (wrapper *WorkspaceWrapper) GetStatus() *schema.WorkspaceStatus {
	return wrapper.resource.Status
}
