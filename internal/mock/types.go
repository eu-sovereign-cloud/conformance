package mock

type MockParams struct {
	MockURL   string
	Tenant    string
	Workspace string
	Region    string
	AuthToken string
}

type scenarioConfig struct {
	name         string
	params       MockParams
	template     string
	response     interface{}
	currentState string
	nextState    string	
	httpStatus   int
	priority     int
}

type regionalMetadataResponseV1 struct {
	Name            string
	Provider        string
	Resource        string
	Verb            string
	CreatedAt       string
	LastModifiedAt  string
	ResourceVersion int
	ApiVersion      string
	Kind            string
	Tenant          string
	Region          string
}

type statusResponseV1 struct {
	State            string
	LastTransitionAt string
}

type workspaceResponseV1 struct {
	Metadata regionalMetadataResponseV1
	Status   statusResponseV1
}
