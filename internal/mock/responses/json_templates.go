package mock

const (
	WorkspacePutTemplateResponse = `
	{
		"metadata": {
			"name": "{{.Name}}",
			"createdAt": "{{.CreatedAt}}",
			"lastModifiedAt": "{{.LastModifiedAt}}",
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
					"lastTransitionAt": "{{.LastTransitionAt}}"
				}
			]
		}
	}`
)
