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
	Name             string
	Workspace        string
	Zone             string
	CreatedAt        string
	LastModifiedAt   string
	Tenant           string
	Region           string
	Version          string
	Kind             string
	Resource         string
	State            string
	LastTransitionAt string
}

type UsecaseStubMetadata struct {
	Params             MockParams
	Metadata           UsecaseMetadata
	Template           string
	RequestTemplate    string
	ScenarioState      string
	NextScenarioState  string
	ScenarioPriority   int
	ScenarioHttpStatus int
}
