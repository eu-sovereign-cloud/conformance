//nolint:dupl
package mock

import (
	"net/http"

	authorization "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.authorization.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Role

func (configurator *scenarioConfigurator) configureCreateRoleStub(response *schema.Role, url string, params HasParams) error {
	setCreatedGlobalTenantResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateRoleStub(response *schema.Role, url string, params HasParams) error {
	setModifiedGlobalTenantResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveRoleStub(response *schema.Role, url string, params HasParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListRoleStub(response *authorization.RoleIterator, url string, params HasParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Role assignment

func (configurator *scenarioConfigurator) configureCreateRoleAssignmentStub(response *schema.RoleAssignment, url string, params HasParams) error {
	setCreatedGlobalTenantResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateRoleAssignmentStub(response *schema.RoleAssignment, url string, params HasParams) error {
	setModifiedGlobalTenantResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveRoleAssignmentStub(response *schema.RoleAssignment, url string, params HasParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListRoleAssignmentStub(response *authorization.RoleAssignmentIterator, url string, params HasParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
