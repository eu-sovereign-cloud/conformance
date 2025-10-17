package secalib

import (
	"fmt"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func AsReferenceURN(ref schema.Reference) (string, error) {
	urn, err := ref.AsReferenceURN()
	if err != nil {
		return "", fmt.Errorf("error extracting URN from reference: %w", err)
	}
	return string(urn), nil
}
