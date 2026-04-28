package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/wiremock"
	region "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.region.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Region
func (configurator *Configurator) ConfigureListRegionStub(response *region.RegionIterator, url string, params wiremock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetWithPathStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetRegionStub(response *schema.Region, url string, params wiremock.MockParams) error {
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
