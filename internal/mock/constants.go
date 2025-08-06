package mock

type MockParams struct {
	WireMockURL   string
	TenantName    string
	WorkspaceName string
	Name          string
	SkuName       string
	Region        string
	Token         string
}

type UsecaseMetadata struct {
	Name      string
	Tenant    string
	Workspace string
	Region    string
	Version   string
	Kind      string
	Resource  string
	State     string
	Zone      string
}

type UsecaseStubMetadata struct {
	Params             MockParams
	Metadata           UsecaseMetadata
	Template           string
	ScenarioState      string
	NextScenarioState  string
	ScenarioPriority   int
	ScenarioHttpStatus int
}
