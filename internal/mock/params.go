package mock

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

// TODO Find a better name
type BaseParams struct {
	MockURL   string
	AuthToken string

	Tenant string
	Region string
}

func (p BaseParams) getBaseParams() *BaseParams {
	return &p
}

type ResourceParams[T any] struct {
	Name          string
	InitialLabels schema.Labels
	UpdatedLabels schema.Labels
	InitialSpec   *T
	UpdatedSpec   *T
}

type stubConfig struct {
	url          string
	httpMethod   string
	httpStatus   int
	params       *BaseParams
	pathParams   map[string]string
	headers      map[string]string
	responseBody any
	currentState string
	nextState    string
	priority     int
}
