package mock

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

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
	Name          string
	InitialLabels schema.Labels
	UpdatedLabels schema.Labels
	InitialSpec   *T
	UpdatedSpec   *T
}

type stubConfig struct {
	url          string
	params       HasParams
	pathParams   map[string]string
	responseBody any
	currentState string
	nextState    string
	priority     int
}
