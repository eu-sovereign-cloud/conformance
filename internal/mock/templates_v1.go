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
)
