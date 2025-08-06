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
		},
      	"labels": \"{{jsonPath request.body '$.labels' fallback null}}\",
      	"annotations": \"{{jsonPath request.body '$.annotations' fallback null}}\"

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
)
