package mock

type MockParams struct {
	WireMockURL   string
	TenantName    string
	WorkspaceName string
	Region        string
	Token         string
}

type UsecaseMetadata struct {
	Name     string
	Tenant   string
	Region   string
	Version  string
	Kind     string
	Resource string
	State    string
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
