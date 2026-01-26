package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Workspace

func (configurator *stubConfigurator) ConfigureCreateWorkspaceStub(response *schema.Workspace, url string, params *mock.MockParams) error {
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = newWorkspaceStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureUpdateWorkspaceStubWithLabels(response *schema.Workspace, url string, params *mock.MockParams, labels schema.Labels) error {
	setModifiedRegionalResourceMetadata(response.Metadata)
	setWorkspaceState(response.Status, schema.ResourceStateUpdating)
	response.Labels = labels
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetActiveWorkspaceStub(response *schema.Workspace, url string, params *mock.MockParams) error {
	setWorkspaceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetListActiveWorkspaceStub(response *workspace.WorkspaceIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
