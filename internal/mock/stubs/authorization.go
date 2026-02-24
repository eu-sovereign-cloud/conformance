//nolint:dupl
package stubs

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	authorization "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.authorization.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Role

func (configurator *Configurator) ConfigureCreateRoleStub(response *schema.Role, url string, params *mock.MockParams) error {
	setCreatedGlobalTenantResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateRoleStub(response *schema.Role, url string, params *mock.MockParams) error {
	setModifiedGlobalTenantResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingRoleStub(response *schema.Role, url string, params *mock.MockParams) error {
	setResourceState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveRoleStub(response *schema.Role, url string, params *mock.MockParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetUpdatingRoleStub(response *schema.Role, url string, params *mock.MockParams) error {
	setResourceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetListRoleStub(response *authorization.RoleIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Role assignment

func (configurator *Configurator) ConfigureCreateRoleAssignmentStub(response *schema.RoleAssignment, url string, params *mock.MockParams) error {
	setCreatedGlobalTenantResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateRoleAssignmentStub(response *schema.RoleAssignment, url string, params *mock.MockParams) error {
	setModifiedGlobalTenantResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingRoleAssignmentStub(response *schema.RoleAssignment, url string, params *mock.MockParams) error {
	setResourceState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveRoleAssignmentStub(response *schema.RoleAssignment, url string, params *mock.MockParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetUpdatingRoleAssignmentStub(response *schema.RoleAssignment, url string, params *mock.MockParams) error {
	setResourceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetListRoleAssignmentStub(response *authorization.RoleAssignmentIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
