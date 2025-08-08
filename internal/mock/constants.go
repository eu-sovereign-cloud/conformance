package mock

type MockParams struct {
	WireMockURL string
	Token       string
}

type UsecaseMetadata struct {
	Zone             string
	CreatedAt        string
	LastModifiedAt   string
	Version          string
	Kind             string
	Resource         string
	State            string
	LastTransitionAt string
}

type StubMetadata struct {
	Metadata           UsecaseMetadata
	Template           string
	RequestTemplate    string
	ScenarioState      string
	NextScenarioState  string
	ScenarioPriority   int
	ScenarioHttpStatus int

	MockConfig    MockParams
	Storage       StorageTemplateConfig
	Compute       ComputeTemplateConfig
	Network       NetworkTemplateConfig
	Authorization AuthorizationTemplateConfig
}

type StorageTemplateConfig struct {
	SkuName          string
	SkuRef           string
	ImageName        string
	BootVolume       int
	CpuArchitecture  string
	BlockStorageRef  string
	BlockStorageName string
	SizeGB           int

	WorkspaceName string
	Tenant        string
	Region        string
}
type WorkspaceTemplateConfig struct {
	SkuName         string
	ImageName       string
	BootVolume      string
	CpuArchitecture string
	BlockStorageRef string
	SizeGB          int
}
type ComputeTemplateConfig struct {
	SkuName         string
	ImageName       string
	BootVolume      string
	CpuArchitecture string
	BlockStorageRef string
	SizeGB          int
}

type NetworkTemplateConfig struct {
	SkuName         string
	ImageName       string
	BootVolume      string
	CpuArchitecture string
	BlockStorageRef string
	SizeGB          int
}

type AuthorizationTemplateConfig struct {
	SkuName         string
	ImageName       string
	BootVolume      string
	CpuArchitecture string
	BlockStorageRef string
	SizeGB          int
}
