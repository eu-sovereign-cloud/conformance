package mock

const (
	WorkspacePutTemplateResponse = `
	{
		"metadata": {
			"name": "{{.Name}}",
			"createdAt": "2025-07-21T15:18:49Z",
			"lastModifiedAt": "2025-07-21T15:18:49Z",
			"tenant": "{{.Tenant}}",
			"region": "{{.Region}}",
			"apiVersion": "{{.Version}}",
			"kind": "{{.Kind}}",
			"resource": "{{.Resource}}",
			"verb": "put"
		},
		"spec": {},
		"status": {
			"state": "{{.State}}",
			"conditions": [
				{
					"state": "{{.State}}",
					"lastTransitionAt": "2025-07-21T15:18:49Z"
				}
			]
		}
	}`
)
