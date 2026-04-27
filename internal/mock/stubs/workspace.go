package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/wiremock"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Workspace

func (configurator *Configurator) ConfigureCreateWorkspaceStub(response *schema.Workspace, url string, params wiremock.MockParams) error {
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = newWorkspaceStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateWorkspaceStubWithLabels(response *schema.Workspace, url string, params wiremock.MockParams, labels schema.Labels) error {
	setModifiedRegionalResourceMetadata(response.Metadata)
	response.Labels = labels
	response.Status = beforeUpdateWorkspaceStatus()
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingWorkspaceStub(response *schema.Workspace, url string, params wiremock.MockParams) error {
	setWorkspaceState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveWorkspaceStub(response *schema.Workspace, url string, params wiremock.MockParams) error {
	setWorkspaceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetUpdatingWorkspaceStub(response *schema.Workspace, url string, params wiremock.MockParams) error {
	setWorkspaceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetDeletingWorkspaceStub(response *schema.Workspace, url string, params wiremock.MockParams) error {
	setWorkspaceState(response.Status, schema.ResourceStateDeleting)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListActiveWorkspaceStub(response *workspace.WorkspaceIterator, url string, params wiremock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetWithPathStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
