package mock

const (
	workspaceResponseTemplateV1 = `
	{
		"labels": {},
		"annotations": {},
		"extensions": {},
		"metadata": {
			"name": "[[.Metadata.Name]]",
			"provider": "[[.Metadata.Provider]]",
			"resource": "[[.Metadata.Resource]]",
			"verb": "[[.Metadata.Verb]]",
			"createdAt": "[[.Metadata.CreatedAt]]",
			"lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
			"resourceVersion": [[.Metadata.ResourceVersion]],
			"apiVersion": "[[.Metadata.ApiVersion]]",
			"kind": "[[.Metadata.Kind]]",
			"tenant": "[[.Metadata.Tenant]]",
			"region": "[[.Metadata.Region]]"
		},
		"spec": {},
		"status": {
			"state": "[[.Status.State]]",
			"conditions": [
				{
					"state": "[[.Status.State]]",
					"lastTransitionAt": "[[.Status.LastTransitionAt]]"
				}
			]
		}
	}`

	instanceSkuResponseTemplateV1 = `
	{
		"labels": {
			"architecture": "[[.Architecture]]",
        	"provider": "[[.Provider]]",
        	"tier": "[[.Tier]]"		
		},
		"annotations": {},
		"extensions": {},	
		"metadata": {
			"name": "[[.Metadata.Name]]",
			"provider": "[[.Metadata.Provider]]",
			"resource": "[[.Metadata.Resource]]",
			"verb": "[[.Metadata.Verb]]",
			"createdAt": "[[.Metadata.CreatedAt]]",
			"lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
			"resourceVersion": [[.Metadata.ResourceVersion]],
			"apiVersion": "[[.Metadata.ApiVersion]]",
			"kind": "[[.Metadata.Kind]]",
			"tenant": "[[.Metadata.Tenant]]",
			"region": "[[.Metadata.Region]]"
		},
		"spec": {
			"ram": [[.RAM]],
			"vCPU": [[.VCPU]]
		},
		"status": {
			"state": "[[.Status.State]]",
			"conditions": [
				{
					"state": "[[.Status.State]]",
					"lastTransitionAt": "[[.Status.LastTransitionAt]]"
				}
			]
		}
	}`

	instanceResponseTemplateV1 = `
	{
		"labels": {},
		"annotations": {},
		"extensions": {},	
		"metadata": {
			"name": "[[.Metadata.Name]]",
			"provider": "[[.Metadata.Provider]]",
			"resource": "[[.Metadata.Resource]]",
			"verb": "[[.Metadata.Verb]]",
			"createdAt": "[[.Metadata.CreatedAt]]",
			"lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
			"resourceVersion": [[.Metadata.ResourceVersion]],
			"apiVersion": "[[.Metadata.ApiVersion]]",
			"kind": "[[.Metadata.Kind]]",
			"tenant": "[[.Metadata.Tenant]]",
			"workspace": "[[.Metadata.Workspace]]",
			"region": "[[.Metadata.Region]]"
		},
		"spec": {
			"skuRef": "[[.SkuRef]]",
        	"zone": "[[.Zone]]",
        	"bootVolume": {
				"deviceRef": "[[.BootDeviceRef]]"
			}
		},
		"status": {
			"state": "[[.Status.State]]",
			"conditions": [
				{
					"state": "[[.Status.State]]",
					"lastTransitionAt": "[[.Status.LastTransitionAt]]"
				}
			]
		}
	}`

	storageSkuResponseTemplateV1 = `
	{
		"labels": {
			"provider": "[[.Provider]]",
        	"tier": "[[.Tier]]"
		},
		"annotations": {},
		"extensions": {},	
		"metadata": {
			"name": "[[.Metadata.Name]]",
			"provider": "[[.Metadata.Provider]]",
			"resource": "[[.Metadata.Resource]]",
			"verb": "[[.Metadata.Verb]]",
			"createdAt": "[[.Metadata.CreatedAt]]",
			"lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
			"resourceVersion": [[.Metadata.ResourceVersion]],
			"apiVersion": "[[.Metadata.ApiVersion]]",
			"kind": "[[.Metadata.Kind]]",
			"tenant": "[[.Metadata.Tenant]]",
			"region": "[[.Metadata.Region]]"
		},
		"spec": {
			"iops": [[.Iops]],
			"type": "[[.Type]]",
			"minVolumeSize": [[.MinVolumeSize]]
		},
		"status": {
			"state": "[[.Status.State]]",
			"conditions": [
				{
					"state": "[[.Status.State]]",
					"lastTransitionAt": "[[.Status.LastTransitionAt]]"
				}
			]
		}
	}`

	blockStorageResponseTemplateV1 = `
	{
		"labels": {},
		"annotations": {},
		"extensions": {},	
		"metadata": {
			"name": "[[.Metadata.Name]]",
			"provider": "[[.Metadata.Provider]]",
			"resource": "[[.Metadata.Resource]]",
			"verb": "[[.Metadata.Verb]]",
			"createdAt": "[[.Metadata.CreatedAt]]",
			"lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
			"resourceVersion": [[.Metadata.ResourceVersion]],
			"apiVersion": "[[.Metadata.ApiVersion]]",
			"kind": "[[.Metadata.Kind]]",
			"tenant": "[[.Metadata.Tenant]]",
			"workspace": "[[.Metadata.Workspace]]",
			"region": "[[.Metadata.Region]]"
		},
		"spec": {
			"skuRef": "[[.SkuRef]]",
			"sizeGB": [[.SizeGB]]
		},
		"status": {
			"state": "[[.Status.State]]",
			"conditions": [
				{
					"state": "[[.Status.State]]",
					"lastTransitionAt": "[[.Status.LastTransitionAt]]"
				}
			]
		}
	}`

	imageResponseTemplateV1 = `
	{
		"labels": {},
		"annotations": {},
		"extensions": {},	
		"metadata": {
			"name": "[[.Metadata.Name]]",
			"provider": "[[.Metadata.Provider]]",
			"resource": "[[.Metadata.Resource]]",
			"verb": "[[.Metadata.Verb]]",
			"createdAt": "[[.Metadata.CreatedAt]]",
			"lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
			"resourceVersion": [[.Metadata.ResourceVersion]],
			"apiVersion": "[[.Metadata.ApiVersion]]",
			"kind": "[[.Metadata.Kind]]",
			"tenant": "[[.Metadata.Tenant]]",
			"workspace": "[[.Metadata.Workspace]]",
			"region": "[[.Metadata.Region]]"
		},
		"spec": {
			"blockStorageRef": "[[.BlockStorageRef]]",
        	"cpuArchitecture": "[[.CpuArchitecture]]"
		},
		"status": {
			"state": "[[.Status.State]]",
			"conditions": [
				{
					"state": "[[.Status.State]]",
					"lastTransitionAt": "[[.Status.LastTransitionAt]]"
				}
			]
		}
	}`

	roleResponseTemplateV1 = `{
      "metadata": {
        "name": "[[.Metadata.Name]]",
        "createdAt": "[[.Metadata.CreatedAt]]",
        "lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
        "resourceVersion": [[.Metadata.ResourceVersion]],
        "tenant": "[[.Metadata.Tenant]]",
        "apiVersion": "[[.Metadata.ApiVersion]]",
        "kind": "[[.Metadata.Kind]]",
        "resource": "[[.Metadata.Resource]]",
        "verb": "[[.Metadata.Verb]]"
      },
      "spec": {
        "permissions": {
          "provider": "[[.Permissions.Provider]]",
          "resources": "[[.Permissions.Resources]]",
          "verbs": "[[.Permissions.Verbs]]"
        }
      },
      "status": {
        "state": "[[.Status.State]]",
        "conditions": [
          {
            "state": "[[.Status.Conditions.State]]",
            "lastTransitionAt": "[[.Status.Conditions.LastTransitionAt]]"
          }
        ]
      }
    }`

	roleAssignmentResponseTemplateV1 = `{
      "metadata": {
        "name": "[[.Metadata.Name]]",
        "createdAt": "[[.Metadata.CreatedAt]]",
        "lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
        "resourceVersion": [[.Metadata.ResourceVersion]],
        "tenant": "[[.Metadata.Tenant]]",
        "apiVersion": "[[.Metadata.ApiVersion]]",
        "kind": "[[.Metadata.Kind]]",
        "resource": "[[.Metadata.Resource]]",
        "verb": "[[.Metadata.Verb]]"
      },
      "spec": {
        "subs": "[[.Subs]]",
        "roles": "[[.Roles]]",
		"scopes": {
			"tenants": "[[.Scopes.Tenants]]",
			"regions": "[[.Scopes.Regions]]",
			"workspaces": "[[.Scopes.Workspaces]]"
		}
      },
      "status": {
        "state": "[[.Status.State]]",
        "conditions": [
          {
            "state": "[[.Status.Conditions.State]]",
            "lastTransitionAt": "[[.Status.Conditions.LastTransitionAt]]"
          }
        ]
      }
    }`
)
