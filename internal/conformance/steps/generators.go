package steps

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Skus

// TODO Find a better package to these methods

func GenerateStorageSkusV1(tenant string) []schema.StorageSku {
	return []schema.StorageSku{
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "RD100",
				Provider: conformance.StorageProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "RD100"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				conformance.ProviderLabel: "seca",
				conformance.TierLabel:     "RD100",
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
				Provider: conformance.StorageProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "RD500"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				conformance.ProviderLabel: "seca",
				conformance.TierLabel:     "RD500",
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
				Provider: conformance.StorageProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "RD2K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				conformance.ProviderLabel: "seca",
				conformance.TierLabel:     "RD2k",
			},
			Spec: &schema.StorageSkuSpec{
				Iops:          2000,
				MinVolumeSize: 50,
				Type:          schema.StorageSkuTypeRemoteDurable,
			},
		},
	}
}

func GenerateInstanceSkusV1(tenant string) []schema.InstanceSku {
	return []schema.InstanceSku{
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "D2XS",
				Provider: conformance.ComputeProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "D2XS"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				// TODO Create constants
				conformance.ArchitectureLabel: "amd64",
				conformance.ProviderLabel:     "seca",
				conformance.TierLabel:         "D2XS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  1,
				VCPU: 1,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "DXS",
				Provider: conformance.ComputeProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "DXS"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				conformance.ArchitectureLabel: "amd64",
				conformance.ProviderLabel:     "seca",
				conformance.TierLabel:         "DXS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  2,
				VCPU: 1,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "DS",
				Provider: conformance.ComputeProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "DS"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				conformance.ArchitectureLabel: "amd64",
				conformance.ProviderLabel:     "seca",
				conformance.TierLabel:         "DS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  4,
				VCPU: 2,
			},
		},
	}
}

func GenerateNetworkSkusV1(tenant string) []schema.NetworkSku {
	return []schema.NetworkSku{
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "N1K",
				Provider: conformance.NetworkProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "N1K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				conformance.ProviderLabel: "seca",
				conformance.TierLabel:     "N1K",
			},
			Spec: &schema.NetworkSkuSpec{
				Bandwidth: 1000,
				Packets:   10000,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "N5K",
				Provider: conformance.NetworkProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "N5K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				conformance.ProviderLabel: "seca",
				conformance.TierLabel:     "N5K",
			},
			Spec: &schema.NetworkSkuSpec{
				Bandwidth: 5000,
				Packets:   40000,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "N10K",
				Provider: conformance.NetworkProviderV1,
				Resource: generators.GenerateSkuResource(tenant, "N10K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				conformance.ProviderLabel: "seca",
				conformance.TierLabel:     "N10K",
			},
			Spec: &schema.NetworkSkuSpec{
				Bandwidth: 10000,
				Packets:   80000,
			},
		},
	}
}
