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
        "permissions": [
		 [[- range $i, $p :=.Permissions]]
		 [[if $i]],[[end]]
			{
			"provider": "[[$p.Provider]]",
			"resources": "[[$p.Resources]]",
			"verbs": "[[$p.Verbs]]"
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
		"scopes": [
		 [[- range $i, $p :=.Scopes]]
		 [[if $i]],[[end]]
			{
			"tenants": "[[$p.Tenants]]",
			"regions": "[[$p.Regions]]",
			"workspaces": "[[$p.Workspaces]]"
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

	networkResponseTemplateV1 = `
	{
      "metadata": {
        "name": "[[.Metadata.Name]]",
        "createdAt": "[[.Metadata.CreatedAt]]",
        "lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
        "resourceVersion": [[.Metadata.ResourceVersion]],
        "tenant": "[[.Metadata.Tenant]]",
        "apiVersion": "[[.Metadata.ApiVersion]]",
        "kind": "[[.Metadata.Kind]]",
        "resource": "[[.Metadata.Resource]]",
        "verb": "[[.Metadata.Verb]]",
        "workspace": "[[.Metadata.Workspace]]"
      },
      "spec": {
        "cidr": "[[.Cidr]]",
        "skuRef": "[[.SkuRef]]",
        "routeTableRef": "[[.RouteTableRef]]"
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

	networkSkuResponseTemplateV1 = `
	{
      "metadata": {
        "name": "[[.Metadata.Name]]"
      },
      "spec": {
        "bandwidth": "[[.Bandwidth]]",
        "packets": "[[.Packets]]"
      }
    }`

	internetGatewayResponseTemplateV1 = `
	{
      "metadata": {
        "name": "[[.Metadata.Name]]",
        "createdAt": "[[.Metadata.CreatedAt]]",
        "lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
        "resourceVersion": [[.Metadata.ResourceVersion]],
        "tenant": "[[.Metadata.Tenant]]",
        "apiVersion": "[[.Metadata.ApiVersion]]",
        "kind": "[[.Metadata.Kind]]",
        "resource": "[[.Metadata.Resource]]",
        "verb": "[[.Metadata.Verb]]",
        "workspace": "[[.Metadata.Workspace]]"
      },
      "spec": {
        "egressOnly": "[[.EgressOnly]]"
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

	nicResponseTemplateV1 = `
	{
      "metadata": {
        "name": "[[.Metadata.Name]]",
        "createdAt": "[[.Metadata.CreatedAt]]",
        "lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
        "resourceVersion": [[.Metadata.ResourceVersion]],
        "tenant": "[[.Metadata.Tenant]]",
        "apiVersion": "[[.Metadata.ApiVersion]]",
        "kind": "[[.Metadata.Kind]]",
        "resource": "[[.Metadata.Resource]]",
        "verb": "[[.Metadata.Verb]]",
        "workspace": "[[.Metadata.Workspace]]"
      },
      "spec": {
        "addresses": "[[.Addresses]]",
        "subnetRef": "[[.SubnetRef]]"
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

	publicIPResponseTemplateV1 = `
	{
      "metadata": {
        "name": "[[.Metadata.Name]]",
        "createdAt": "[[.Metadata.CreatedAt]]",
        "lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
        "resourceVersion": [[.Metadata.ResourceVersion]],
        "tenant": "[[.Metadata.Tenant]]",
        "apiVersion": "[[.Metadata.ApiVersion]]",
        "kind": "[[.Metadata.Kind]]",
        "resource": "[[.Metadata.Resource]]",
        "verb": "[[.Metadata.Verb]]",
        "workspace": "[[.Metadata.Workspace]]"
      },
      "spec": {
        "version": "[[.Version]]",
        "address": "[[.Address]]"
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

	routeTableResponseTemplateV1 = `
	{
      "metadata": {
        "name": "[[.Metadata.Name]]",
        "createdAt": "[[.Metadata.CreatedAt]]",
        "lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
        "resourceVersion": [[.Metadata.ResourceVersion]],
        "tenant": "[[.Metadata.Tenant]]",
        "apiVersion": "[[.Metadata.ApiVersion]]",
        "kind": "[[.Metadata.Kind]]",
        "resource": "[[.Metadata.Resource]]",
        "verb": "[[.Metadata.Verb]]",
        "workspace": "[[.Metadata.Workspace]]"
      },
      "spec": {
	  	"localRef": "[[.LocalRef]]",
        "routes": [
		 [[- range $i, $r :=.Routes]]
		 [[if $i]],[[end]]
			{
			"destinationCidrBlock": "[[$r.DestinationCidrBlock]]",
			"targetRef": "[[$r.TargetRef]]"
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

	subnetResponseTemplateV1 = `
	{
      "metadata": {
        "name": "[[.Metadata.Name]]",
        "createdAt": "[[.Metadata.CreatedAt]]",
        "lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
        "resourceVersion": [[.Metadata.ResourceVersion]],
        "tenant": "[[.Metadata.Tenant]]",
        "apiVersion": "[[.Metadata.ApiVersion]]",
        "kind": "[[.Metadata.Kind]]",
        "resource": "[[.Metadata.Resource]]",
        "verb": "[[.Metadata.Verb]]",
        "workspace": "[[.Metadata.Workspace]]"
      },
      "spec": {
		"cidr": "[[.Cidr]]",
		"zone": "[[.Zone]]",
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

	securityGroupResponseTemplateV1 = `
	{
      "metadata": {
        "name": "[[.Metadata.Name]]",
        "createdAt": "[[.Metadata.CreatedAt]]",
        "lastModifiedAt": "[[.Metadata.LastModifiedAt]]",
        "resourceVersion": [[.Metadata.ResourceVersion]],
        "tenant": "[[.Metadata.Tenant]]",
        "apiVersion": "[[.Metadata.ApiVersion]]",
        "kind": "[[.Metadata.Kind]]",
        "resource": "[[.Metadata.Resource]]",
        "verb": "[[.Metadata.Verb]]",
        "workspace": "[[.Metadata.Workspace]]"
      },
      "spec": {
		"rules": [
			[[- range $i, $r :=.Rules]]
			[[if $i]],[[end]]
				{
				"direction": "[[$r.Direction]]"
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
)
