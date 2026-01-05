package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Clients

func (configurator *scenarioConfigurator) ConfigureClientsInitStub(response *schema.Region, url string, params *mock.BaseParams) error {
	// setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
