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
				labelsProvider: "seca",
				labelsTier:     "RD100",
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
				labelsProvider: "seca",
				labelsTier:     "RD500",
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
				labelsProvider: "seca",
				labelsTier:     "RD2k",
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
				labelsArchitecture: "amd64",
				labelsProvider:     "seca",
				labelsTier:         "D2XS",
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
				labelsArchitecture: "amd64",
				labelsProvider:     "seca",
				labelsTier:         "DXS",
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
				labelsArchitecture: "amd64",
				labelsProvider:     "seca",
				labelsTier:         "DS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  4,
				VCPU: 2,
			},
		},
	}
}

func generateNetworkSkusV1(tenant string) []schema.NetworkSku {
	return []schema.NetworkSku{
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "N1K",
				Provider: networkProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "N1K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				labelsProvider: "seca",
				labelsTier:     "N1K",
			},
			Spec: &schema.NetworkSkuSpec{
				Bandwidth: 1000,
				Packets:   10000,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "N5K",
				Provider: networkProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "N5K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				labelsProvider: "seca",
				labelsTier:     "N5K",
			},
			Spec: &schema.NetworkSkuSpec{
				Bandwidth: 5000,
				Packets:   40000,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "N10K",
				Provider: networkProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "N10K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				labelsProvider: "seca",
				labelsTier:     "N10K",
			},
			Spec: &schema.NetworkSkuSpec{
				Bandwidth: 10000,
				Packets:   80000,
			},
		},
	}
}
