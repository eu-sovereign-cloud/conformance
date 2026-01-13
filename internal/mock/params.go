package mock

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

// TODO Find a better name
type BaseParams struct {
	MockURL   string
	AuthToken string

	Tenant string
	Region string
}

func (p BaseParams) GetBaseParams() *BaseParams {
	return &p
}

type ResourceParams[T any] struct {
	Name          string
	InitialLabels schema.Labels
	UpdatedLabels schema.Labels
	InitialSpec   *T
	UpdatedSpec   *T
}
