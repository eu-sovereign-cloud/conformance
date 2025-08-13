package mock

import (
	"bytes"
	"encoding/json"
	"text/template"
)

const (
	workspaceTemplateResponse = `
	{
		"metadata": {
			"name": "{{.Metadata.Name}}",
			"provider": "{{.Metadata.Provider}}",
			"resource": "{{.Metadata.Resource}}",
			"verb": "{{.Metadata.Verb}}",
			"createdAt": "{{.Metadata.CreatedAt}}",
			"lastModifiedAt": "{{.Metadata.LastModifiedAt}}",
			"resourceVersion": {{.Metadata.ResourceVersion}},
			"apiVersion": "{{.Metadata.ApiVersion}}",
			"kind": "{{.Metadata.Kind}}",
			"tenant": "{{.Metadata.Tenant}}",
			"region": "{{.Metadata.Region}}"
		},
		"spec": {},
		"status": {
			"state": "{{.Status.State}}",
			"conditions": [
				{
					"state": "{{.Status.State}}",
					"lastTransitionAt": "{{.Status.LastTransitionAt}}"
				}
			]
		}
	}`
)

func processTemplate(templ string, data any) (map[string]interface{}, error) {
	tmpl := template.Must(template.New("response").Parse(templ))

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return nil, err
	}
	var dataJsonMap map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &dataJsonMap)
	if err != nil {
		return nil, err
	}

	return dataJsonMap, nil
}
