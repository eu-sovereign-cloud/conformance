package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Workspace

func (configurator *scenarioConfigurator) ConfigureCreateWorkspaceStub(response *schema.Workspace, url string, params *mock.BaseParams) error {
	setCreatedRegionalResourceMetadata(response.Metadata)
	response.Status = newWorkspaceStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateWorkspaceStubWithLabels(response *schema.Workspace, url string, params *mock.BaseParams, labels schema.Labels) error {
	setModifiedRegionalResourceMetadata(response.Metadata)
	setWorkspaceState(response.Status, schema.ResourceStateUpdating)
	response.Labels = labels
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveWorkspaceStub(response *schema.Workspace, url string, params *mock.BaseParams) error {
	setWorkspaceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListActiveWorkspaceStub(response *workspace.WorkspaceIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
