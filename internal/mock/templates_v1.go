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
        [[- range $i, $p := .Spec.Permissions]]
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
        [[- range $i, $s := .Spec.Subs]]
          [[if $i]],[[end]]
          "[[$s]]"
        [[- end]]
		    ],
        "roles": [
        [[- range $i, $r := .Spec.Roles]]
          [[if $i]],[[end]]
          "[[$r]]"
        [[- end]]
		    ],
        "scopes": [
        [[- range $i, $s := .Spec.Scopes]]
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

	regionListResponseTemplateV1 = `
  {
      "items": [
      [[- range $i, $r := .]]
        [[if $i]],[[end]]
        {
        "metadata": {
          "name": "[[$r.Metadata.Name]]",
          "provider": "[[$r.Metadata.Provider]]",
          "createdAt": "[[$r.Metadata.CreatedAt]]",
          "lastModifiedAt": "[[$r.Metadata.LastModifiedAt]]",
          "resourceVersion": [[$r.Metadata.ResourceVersion]],
          "tenant": "[[$r.Metadata.Tenant]]",
          "apiVersion": "[[$r.Metadata.ApiVersion]]",
          "kind": "[[$r.Metadata.Kind]]",
          "resource": "[[.Metadata.Resource]]",
          "verb": "[[$r.Metadata.Verb]]"
        },
          "spec": {
            "availableZones": [
              [[- range $i, $az := .Spec.AvailableZones]]
              [[if $i]],[[end]]"[[ $az ]]"
              [[- end]]
            ],
            "providers": [
              [[- range $i, $p := .Spec.Providers]]
                [[if $i]],[[end]]
                    {
                      "name": "[[$p.Name]]",
                      "version": "[[$p.Version]]",
                      "url": "[[$p.URL]]"
                    }[[- end]]
              ]
        }
        } [[- end]]
    ]
  }`
	regionResponseTemplateV1 = `
  {
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
            "availableZones": [
              [[- range $i, $az := .Spec.AvailableZones]]
              [[if $i]],[[end]]"[[ $az ]]"
              [[- end]]
            ],
            "providers": [
              [[- range $i, $p := .Spec.Providers]]
                [[if $i]],[[end]]
                    {
                      "name": "[[$p.Name]]",
                      "version": "[[$p.Version]]",
                      "url": "[[$p.URL]]"
                    }[[- end]]
              ]
        }
  }`

	workspaceResponseTemplateV1 = `
  {
		"labels": {
    [[- range $j, $w := $.Labels]]
      [[if $j]],[[end]]
      "[[$w.Name]]":  "[[$w.Value]]"
    [[- end]]
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
			"skuRef": "[[.Spec.SkuRef]]",
			"sizeGB": [[.Spec.SizeGB]]
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
			"blockStorageRef": "[[.Spec.BlockStorageRef]]",
        	"cpuArchitecture": "[[.Spec.CpuArchitecture]]"
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
			"skuRef": "[[.Spec.SkuRef]]",
        	"zone": "[[.Spec.Zone]]",
        	"bootVolume": {
				"deviceRef": "[[.Spec.BootDeviceRef]]"
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

	networkResponseTemplateV1 = `
	{
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
        "cidr": {
          "ipv4": "[[.Spec.Cidr.Ipv4]]",
          "ipv6": "[[.Spec.Cidr.Ipv6]]"
        },
        "skuRef": "[[.Spec.SkuRef]]",
        "routeTableRef": "[[.Spec.RouteTableRef]]"
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

	internetGatewayResponseTemplateV1 = `
	{
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
        "egressOnly": [[.Spec.EgressOnly]]
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
        "addresses": [
        [[- range $i, $a := .Spec.Addresses]]
          [[if $i]],[[end]]
          "[[$a]]"
        [[- end]]
		    ],
        "subnetRef": "[[.Spec.SubnetRef]]"
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
        "version": "[[.Spec.Version]]",
        "address": "[[.Spec.Address]]"
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
        "network": "[[.Metadata.Network]]",
        "region": "[[.Metadata.Region]]"
      },
      "spec": {
	  	  "localRef": "[[.Spec.LocalRef]]",
        "routes": [
		    [[- range $i, $r := .Spec.Routes]]
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
        "network": "[[.Metadata.Network]]",
        "region": "[[.Metadata.Region]]"
      },
      "spec": {
        "cidr": {
          "ipv4": "[[.Spec.Cidr.Ipv4]]"
        },
		    "zone": "[[.Spec.Zone]]"
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
		"rules": [
			[[- range $i, $r := .Spec.Rules]]
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
