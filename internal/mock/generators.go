package mock

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Skus

func generateStorageSkusV1(tenant string) []schema.StorageSku {
	return []schema.StorageSku{
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "RD100",
				Provider: storageProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "RD100"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				"provider": "seca",
				"tier":     "RD100",
			},
			Spec: &schema.StorageSkuSpec{
				Iops:          100,
				MinVolumeSize: 50,
				Type:          schema.StorageSkuTypeRemoteDurable,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "RD500",
				Provider: storageProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "RD500"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				"provider": "seca",
				"tier":     "RD500",
			},
			Spec: &schema.StorageSkuSpec{
				Iops:          500,
				MinVolumeSize: 50,
				Type:          schema.StorageSkuTypeRemoteDurable,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "RD2K",
				Provider: storageProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "RD2K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				"provider": "seca",
				"tier":     "RD2k",
			},
			Spec: &schema.StorageSkuSpec{
				Iops:          2000,
				MinVolumeSize: 50,
				Type:          schema.StorageSkuTypeRemoteDurable,
			},
		},
	}
}

func generateInstanceSkusV1(tenant string) []schema.InstanceSku {
	return []schema.InstanceSku{
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "D2XS",
				Provider: computeProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "D2XS"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				// TODO Create constants
				"architecture": "amd64",
				"provider":     "seca",
				"tier":         "D2XS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  1,
				VCPU: 1,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "DXS",
				Provider: computeProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "DXS"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				"architecture": "amd64",
				"provider":     "seca",
				"tier":         "DXS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  2,
				VCPU: 1,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "DS",
				Provider: computeProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "DS"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				"architecture": "amd64",
				"provider":     "seca",
				"tier":         "DS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  4,
				VCPU: 2,
			},
		},
	}
}
