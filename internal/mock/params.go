package mock

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

type Params struct {
	MockURL   string
	AuthToken string

	Tenant string
	Region string
}
type HasParams interface {
	getParams() *Params
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
	params       HasParams
	responseBody any
	currentState string
	nextState    string
	priority     int
}
