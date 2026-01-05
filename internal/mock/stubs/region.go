package stubs

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	region "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.region.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Region
func (configurator *scenarioConfigurator) ConfigureGetListRegionStub(response *region.RegionIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetRegionStub(response *schema.Region, url string, params *mock.BaseParams) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
