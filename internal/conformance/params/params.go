package params

import (
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
)

// Clients

type ClientsInitParams struct {
	*mock.MockParams
	Region string
}
