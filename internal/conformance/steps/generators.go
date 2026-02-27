package steps

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Skus

// TODO Find a better package to it
func GenerateStorageSkusV1(tenant string) []schema.StorageSku {
	return []schema.StorageSku{
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "RD100",
				Provider: sdkconsts.StorageProviderV1Name,
				Resource: generators.GenerateSkuResource(tenant, "RD100"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				constants.ProviderLabel: constants.ProviderSecaLabel,
				constants.TierLabel:     "RD100",
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
				Provider: sdkconsts.StorageProviderV1Name,
				Resource: generators.GenerateSkuResource(tenant, "RD500"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				constants.ProviderLabel: constants.ProviderSecaLabel,
				constants.TierLabel:     "RD500",
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
				Provider: sdkconsts.StorageProviderV1Name,
				Resource: generators.GenerateSkuResource(tenant, "RD2K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				constants.ProviderLabel: constants.ProviderSecaLabel,
				constants.TierLabel:     "RD2k",
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
				Provider: sdkconsts.ComputeProviderV1Name,
				Resource: generators.GenerateSkuResource(tenant, "D2XS"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				constants.ArchitectureLabel: constants.ArchitectureAmd64Label,
				constants.ProviderLabel:     constants.ProviderSecaLabel,
				constants.TierLabel:         "D2XS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  1,
				VCPU: 1,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "DXS",
				Provider: sdkconsts.ComputeProviderV1Name,
				Resource: generators.GenerateSkuResource(tenant, "DXS"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				constants.ArchitectureLabel: constants.ArchitectureAmd64Label,
				constants.ProviderLabel:     constants.ProviderSecaLabel,
				constants.TierLabel:         "DXS",
			},
			Spec: &schema.InstanceSkuSpec{
				Ram:  2,
				VCPU: 1,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "DS",
				Provider: sdkconsts.ComputeProviderV1Name,
				Resource: generators.GenerateSkuResource(tenant, "DS"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				constants.ArchitectureLabel: constants.ArchitectureAmd64Label,
				constants.ProviderLabel:     constants.ProviderSecaLabel,
				constants.TierLabel:         "DS",
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
				Provider: sdkconsts.NetworkProviderV1Name,
				Resource: generators.GenerateSkuResource(tenant, "N1K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				constants.ProviderLabel: constants.ProviderSecaLabel,
				constants.TierLabel:     "N1K",
			},
			Spec: &schema.NetworkSkuSpec{
				Bandwidth: 1000,
				Packets:   10000,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "N5K",
				Provider: sdkconsts.NetworkProviderV1Name,
				Resource: generators.GenerateSkuResource(tenant, "N5K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				constants.ProviderLabel: constants.ProviderSecaLabel,
				constants.TierLabel:     "N5K",
			},
			Spec: &schema.NetworkSkuSpec{
				Bandwidth: 5000,
				Packets:   40000,
			},
		},
		{
			Metadata: &schema.SkuResourceMetadata{
				Name:     "N10K",
				Provider: sdkconsts.NetworkProviderV1Name,
				Resource: generators.GenerateSkuResource(tenant, "N10K"),
				Verb:     http.MethodGet,
				Tenant:   tenant,
			},
			Labels: schema.Labels{
				constants.ProviderLabel: constants.ProviderSecaLabel,
				constants.TierLabel:     "N10K",
			},
			Spec: &schema.NetworkSkuSpec{
				Bandwidth: 10000,
				Packets:   80000,
			},
		},
	}
}
