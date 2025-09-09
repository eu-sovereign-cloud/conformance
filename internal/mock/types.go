package mock

import "github.com/eu-sovereign-cloud/conformance/secalib"

type stubConfig struct {
	url          string
	params       secalib.GeneralParams
	response     any
	template     string
	currentState string
	nextState    string
	httpStatus   int
	priority     int
}

type resourceResponse[T any] struct {
	Metadata *secalib.Metadata
	Status   *secalib.Status
	Labels   *[]secalib.Label
	Spec     *T
}
