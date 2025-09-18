package secalib

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
	Workspace       *string
	Network         *string
	Region          *string
}

type Status struct {
	State            string
	LastTransitionAt string
}

type Label struct {
	Name  string
	Value string
}
