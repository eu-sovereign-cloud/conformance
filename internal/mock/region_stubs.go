package mock

import (
	"net/http"

	region "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.region.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Region
func (configurator *scenarioConfigurator) configureGetListRegionStub(response *region.RegionIterator, url string, params *BaseParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetRegionStub(response *schema.Region, url string, params *BaseParams) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
