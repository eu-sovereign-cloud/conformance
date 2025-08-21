package mock

type Params struct {
	MockURL   string
	AuthToken string

	Tenant    string
	Workspace string
	Region    string
}

type HasParams interface {
	getParams() Params
}

type scenarioConfig struct {
	url          string
	params       HasParams
	response     any
	template     string
	currentState string
	nextState    string
	httpStatus   int
	priority     int
}

type metadataResponse struct {
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
	Workspace       string
	Region          string
}

type statusResponse struct {
	State            string
	LastTransitionAt string
}
