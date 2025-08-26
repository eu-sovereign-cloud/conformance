package mock

const (
	roleResponseTemplateV1 = `{
      "metadata": {
        "name": "[[.Metadata.Name]]",
		"provider": "[[.Metadata.Provider]]",
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
        "permissions": [
        [[- range $i, $p :=.Permissions]]
		  [[if $i]],[[end]]
          {
            "provider": "[[$p.Provider]]",
            "resources": [
            [[- range $j, $r := $p.Resources]]
              [[if $j]],[[end]]
              "[[$r]]"
            [[- end]]
            ],
            "verb": [
            [[- range $j, $v := $p.Verb]]
              [[if $j]],[[end]]
              "[[$v]]"
            [[- end]]
            ]
          }[[- end]]
		]
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

	roleAssignmentResponseTemplateV1 = `{
      "metadata": {
        "name": "[[.Metadata.Name]]",
		"provider": "[[.Metadata.Provider]]",
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
        "subs": [
        [[- range $i, $s := .Subs]]
          [[if $i]],[[end]]
          "[[$s]]"
        [[- end]]
		],
        "roles": [
        [[- range $i, $r := .Roles]]
          [[if $i]],[[end]]
          "[[$r]]"
        [[- end]]
		],
        "scopes": [
        [[- range $i, $s := .Scopes]]
		  [[if $i]],[[end]]
          {
            "tenants": [
            [[- range $j, $t := $s.Tenants]]
              [[if $j]],[[end]]
              "[[$t]]"
            [[- end]]
            ],
            "regions": [
            [[- range $j, $r := $s.Regions]]
              [[if $j]],[[end]]
              "[[$r]]"
            [[- end]]
            ],
            "workspaces": [
            [[- range $j, $w := $s.Workspaces]]
              [[if $j]],[[end]]
              "[[$w]]"
            [[- end]]
            ]
          }[[- end]]
        ]
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
			"type": "[[.StorageType]]",
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
)
