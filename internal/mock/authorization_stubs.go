package mock

import (
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Role

func (configurator *scenarioConfigurator) configureCreateRoleStub(response *schema.Role, url string, params HasParams) error {
	setCreatedGlobalTenantResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateRoleStub(response *schema.Role, url string, params HasParams) error {
	setModifiedGlobalTenantResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateUpdating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveRoleStub(response *schema.Role, url string, params HasParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

// Role assignment

func (configurator *scenarioConfigurator) configureCreateRoleAssignmentStub(response *schema.RoleAssignment, url string, params HasParams) error {
	setCreatedGlobalTenantResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateRoleAssignmentStub(response *schema.RoleAssignment, url string, params HasParams) error {
	setModifiedGlobalTenantResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateUpdating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveRoleAssignmentStub(response *schema.RoleAssignment, url string, params HasParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}
