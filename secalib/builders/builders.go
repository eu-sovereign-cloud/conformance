package builders

import "github.com/eu-sovereign-cloud/conformance/secalib"

type metadataBuilder struct {
	validator *secalib.Validator
}

func newMetadataBuilder() *metadataBuilder {
	return &metadataBuilder{
		validator: secalib.NewValidator(),
	}
}

type resourceBuilder struct {
	validator *secalib.Validator
}

func newResourceBuilder() *resourceBuilder {
	return &resourceBuilder{
		validator: secalib.NewValidator(),
	}
}
