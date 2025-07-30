package mock

type Workspace struct {
	WireMockURL   string
	TenantName    string
	WorkspaceName string
	Region        string
	Token         string
}

type WorkspaceMetadata struct {
	Name     string
	Tenant   string
	Region   string
	Version  string
	Kind     string
	Resource string
	State    string
}

type WorkspaceStubMetadata struct {
	WorkspaceMock      Workspace
	WorkspaceMetadata  WorkspaceMetadata
	Template           string
	ScenarioState      string
	NextScenarioState  string
	ScenarioPriority   int
	ScenarioHttpStatus int
}
