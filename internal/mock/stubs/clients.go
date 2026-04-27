package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/wiremock"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Clients

func (configurator *Configurator) ConfigureClientsInitStub(response *schema.Region, url string, params wiremock.MockParams) error {
	setCreatedGlobalResourceMetadata(response.Metadata)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
