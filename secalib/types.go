package secalib

type GeneralParams struct {
	AuthToken string
	Tenant    string
	Region    string
}

type ResourceParams[T any] struct {
	Name        string
	InitialSpec *T
	UpdatedSpec *T
}

type Metadata struct {
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

type Status struct {
	State            string
	LastTransitionAt string
}

type Label struct {
	Name  string
	Value string
}
