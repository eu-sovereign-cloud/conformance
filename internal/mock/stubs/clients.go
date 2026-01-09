package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Clients

func (configurator *stubConfigurator) ConfigureClientsInitStub(response *schema.Region, url string, params *mock.MockParams) error {
	setCreatedGlobalResourceMetadata(response.Metadata)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
