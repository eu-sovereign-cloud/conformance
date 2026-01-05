//nolint:dupl
package stubs

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	authorization "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.authorization.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Role

func (configurator *scenarioConfigurator) ConfigureCreateRoleStub(response *schema.Role, url string, params *mock.BaseParams) error {
	setCreatedGlobalTenantResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateRoleStub(response *schema.Role, url string, params *mock.BaseParams) error {
	setModifiedGlobalTenantResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveRoleStub(response *schema.Role, url string, params *mock.BaseParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListRoleStub(response *authorization.RoleIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Role assignment

func (configurator *scenarioConfigurator) ConfigureCreateRoleAssignmentStub(response *schema.RoleAssignment, url string, params *mock.BaseParams) error {
	setCreatedGlobalTenantResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateRoleAssignmentStub(response *schema.RoleAssignment, url string, params *mock.BaseParams) error {
	setModifiedGlobalTenantResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveRoleAssignmentStub(response *schema.RoleAssignment, url string, params *mock.BaseParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListRoleAssignmentStub(response *authorization.RoleAssignmentIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
