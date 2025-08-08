package mock

const (
	WorkspaceTemplateResponse = `
	{
		"metadata": {
			"name": "[[.Name]]",
			"createdAt": "[[.CreatedAt]]",
			"lastModifiedAt": "[[.LastModifiedAt]]",
			"tenant": "[[.Tenant]]",
			"region": "[[.Region]]",
			"apiVersion": "[[.Version]]",
			"kind": "[[.Kind]]",
			"resource": "[[.Resource]]",
			"verb": "put"
		},
		"spec": {},
		"status": {
			"state": "[[.State]]",
			"conditions": [
				{
					"state": "[[.State]]",
					"lastTransitionAt": "[[.LastTransitionAt]]"
				}
			]
		}
	}`
	ComputePutTemplateResponse = `
	{
		"metadata": {
			"name": "[[.Name]]",
			"createdAt": "[[.CreatedAt]]",
			"lastModifiedAt": "[[.LastModifiedAt]]",
			"tenant": "[[.Tenant]]",
			"region": "[[.Region]]",
			"apiVersion": "[[.Version]]",
			"kind": "[[.Kind]]",
			"provider": "seca.compute",
			"resource": "[[.Resource]]",
			"workspace": "[[.Workspace]]",
			"verb": "put",
			"zone": "[[.Zone]]"
		},
		"spec": {
        "skuRef": \"{{jsonPath request.body '$.spec.skuRef'}}\",
        "zone": \"{{jsonPath request.body '$.spec.zone'}}\",
        "bootVolume": \"{{jsonPath request.body '$.spec.bootVolume'}}\"
      	},
		"status": {
			"state": "[[.State]]",
			"conditions": [
				{
					"state": "[[.State]]",
					"lastTransitionAt": "2025-07-21T15:18:49Z"
				}
			]
		}

	}`

	ComputeSkuTemplateResponse = `
	{
		"metadata": {
			"name": "[[.Name]]"
		},
		"spec": {
			"benchmarkPoints": 3000,
			"cpuType": "amd64",
			"provider": "SECA",
			"ramGB": 4,
			"tier": "[[.Name]]",
			"vCPU": 2
		}
	}`

	ComputeGetTemplateResponse = `
	{
      "metadata": {
        "name": "[[.Name]]",
        "createdAt": "[[.CreatedAt]]",
        "lastModifiedAt": "[[.LastModifiedAt]]",
        "resourceVersion": 1,
        "tenant": "[[.Tenant]]",
        "region": "[[.Region]]",
        "apiVersion": "v1",
        "kind": "instance",
        "resource": "[[.Resource]]",
        "verb": "get",
        "workspace": "[[.Workspace]]",
        "zone": "[[.Zone]]"
      },
      "spec": {
        "skuRef": "tenants/[[.Tenant]]/skus/seca.m",
        "zone": "[[.Zone]]",
        "bootVolume": {
          "deviceRef": "tenants/[[.Tenant]]/workspaces/[[.Workspace]]/block-storages/boot-volume"
        }
      },
      "status": {
        "state": "active",
        "conditions": [
          {
            "state": "active",
            "lastTransitionAt": "[[.LastTransitionAt]]"
          }
        ]
      }
    }`

	StoragePutImageTemplateResponse = `
	{
		"metadata": {
        "name": "[[.Storage.ImageName]]",
        "createdAt": "[[.Metadata.CreatedAt]]",
        "lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
        "resourceVersion": 1,
        "tenant": "[[.Storage.Tenant]]",
        "region": "[[.Storage.Region]]",
        "apiVersion": "[[.Metadata.Version]]",
        "kind": "[[.Metadata.Kind]]",
        "resource": "seca.storage/images/[[.Storage.ImageName]]",
        "verb": "put"
      },
      "spec": {
        "blockStorageRef": "{{jsonPath request.body '$.spec.blockStorageRef'}}",
        "cpuArchitecture": "{{jsonPath request.body '$.spec.cpuArchitecture'}}"
      },
      "status": {
        "state": "[[.Metadata.State]]",
        "conditions": [
          {
            "state": "[[.Metadata.State]]",
            "lastTransitionAt": "[[.Metadata.LastTransitionAt]]"
          }
        ]
      }
	}`

	StorageGetImageTemplateResponse = `
	{
		"metadata": {
			"name": "[[.Storage.ImageName]]",
			"createdAt": "[[.Metadata.CreatedAt]]",
			"lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
			"resourceVersion": 1,
			"tenant": "[[.Storage.Tenant]]",
			"region": "[[.Storage.Region]]",
			"apiVersion": "[[.Metadata.Version]]",
			"kind": "[[.Metadata.Kind]]",
			"resource": "seca.storage/images/[[.Storage.ImageName]]",
			"verb": "get"
		},
		"spec": {
			"blockStorageRef": "[[.Storage.BlockStorageRef]]",
			"cpuArchitecture": "[[.Storage.CpuArchitecture]]"
		},
		"status": {
			"state": "[[.Metadata.State]]",
			"conditions": [
			{
				"state": "[[.Metadata.State]]",
				"lastTransitionAt": "[[.Metadata.LastTransitionAt]]"
			}
			]
		}
    }`

	StoragePutBlockStorageTemplateResponse = `
	{
		"metadata": {
			"name": "[[.Storage.BlockStorageName]]",
			"createdAt": "[[.Metadata.CreatedAt]]",
			"lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
			"resourceVersion": 1,
			"tenant": "[[.Storage.Tenant]]",
			"region": "[[.Storage.Region]]",
			"apiVersion": "[[.Metadata.Version]]",
			"kind": "[[.Metadata.Kind]]",
			"resource": "seca.storage/images/[[.Storage.BlockStorageName]]",
			"verb": "put"
		},
		"spec": {
			"skuRef": "{{jsonPath request.body '$.spec.skuRef'}}",
			"sizeGB": "{{jsonPath request.body '$.spec.sizeGB'}}"
		},
		"status": {
			"state": "[[.Metadata.State]]",
			"conditions": [
			{
				"state": "[[.Metadata.State]]",
				"lastTransitionAt": "[[.Metadata.LastTransitionAt]]"
			}
			]
		}
	}`

	StorageGetBlockStorageTemplateResponse = `
	{
		"metadata": {
			"name": "[[.Storage.BlockStorageName]]",
			"createdAt": "[[.Metadata.CreatedAt]]",
			"lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
			"resourceVersion": 1,
			"tenant": "[[.Storage.Tenant]]",
			"region": "[[.Storage.Region]]",
			"apiVersion": "[[.Metadata.Version]]",
			"kind": "[[.Metadata.Kind]]",
			"resource": "seca.storage/images/[[.Storage.BlockStorageName]]",
			"verb": "put"
		},
		"spec": {
			"skuRef": "[[.Storage.SkuRef]]",
			"sizeGB": "[[.Storage.SizeGB]]"
		},
		"status": {
			"state": "[[.Metadata.State]]",
			"conditions": [
			{
				"state": "[[.Metadata.State]]",
				"lastTransitionAt": "[[.Metadata.LastTransitionAt]]"
			}
			]
		}
	}`

	StorageGetSkuTemplateResponse = `
	{
		"metadata": {
			"name": "[[.Storage.ImageName]]",
			"createdAt": "[[.Metadata.CreatedAt]]",
			"lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
			"resourceVersion": 1,
			"tenant": "[[.Storage.Tenant]]",
			"region": "[[.Storage.Region]]",
			"apiVersion": "[[.Metadata.Version]]",
			"kind": "[[.Metadata.Kind]]",
			"resource": "seca.storage/images/[[.Storage.ImageName]]",
			"verb": "put"
		},
		"spec": {
			"provider": "SECA",
			"type": "remote-durable",
			"iops": 100,
			"minVolumeSize": 1
		}
	}`
)
