package mock

import "github.com/eu-sovereign-cloud/conformance/secalib"

// TODO Find a better name
type Params struct {
	MockURL   string
	AuthToken string

	Tenant string
	Region string
}

type HasParams interface {
	getParams() *Params
}

// TODO Find a better name
type ResourceParams[T any] struct {
	Name        string
	InitialSpec *T
	UpdatedSpec *T
}

type stubConfig struct {
	url          string
	params       HasParams
	response     any
	currentState string
	nextState    string
	httpStatus   int
	priority     int
}

// TODO Find a better name
type resourceResponse[T any] struct {
	Metadata *secalib.Metadata
	Status   *secalib.Status
	Labels   *[]secalib.Label
	Spec     *T
}
