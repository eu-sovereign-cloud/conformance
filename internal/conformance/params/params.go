package params

import (
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

type BaseParams struct {
	*mock.MockParams
	Tenant string
	Region string
}

type ResourceParams[T any] struct {
	Name          string
	InitialLabels schema.Labels
	UpdatedLabels schema.Labels
	InitialSpec   *T
	UpdatedSpec   *T
}

// Clients

type ClientsInitParams struct {
	*BaseParams
}
