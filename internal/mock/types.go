package mock

import "github.com/eu-sovereign-cloud/conformance/secalib"

// TODO Find a better name
type Params struct {
	MockURL   string
	AuthToken string

	Tenant    string
	Workspace string
	Region    string
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
	template     string
	currentState string
	nextState    string
	httpStatus   int
	priority     int
}

// TODO Find a better name
type resourceResponse[T any] struct {
	Metadata *secalib.Metadata
	Status   *secalib.Status
	Spec     *T
}
