package mock

// Params

type AuthorizationParamsV1 struct {
	Params
	Role           RoleParamsV1
	RoleAssignment RoleAssignmentParamsV1
}
type RoleParamsV1 struct {
	Name        string
	Permissions []RolePermissionParamsV1
}
type RolePermissionParamsV1 struct {
	Provider    string
	Resources   []string
	VerbInitial []string
	VerbUpdated []string
}

type RoleAssignmentParamsV1 struct {
	Name        string
	Roles       []string
	SubsInitial []string
	SubsUpdated []string
	Scopes      []RoleAssignmentScopeParamsV1
}
type RoleAssignmentScopeParamsV1 struct {
	Tenants    []string
	Regions    []string
	Workspaces []string
}

func (p AuthorizationParamsV1) getParams() Params { return p.Params }

type WorkspaceParamsV1 struct {
	Params
	Name string
}

func (p WorkspaceParamsV1) getParams() Params { return p.Params }

type ComputeParamsV1 struct {
	Params
	StorageSku   StorageSkuParamsV1
	BlockStorage BlockStorageParamsV1
	InstanceSku  InstanceSkuParamsV1
	Instance     InstanceParamsV1
}
type InstanceSkuParamsV1 struct {
	Name         string
	Architecture string
	Provider     string
	Tier         string
	RAM          int
	VCPU         int
}
type InstanceParamsV1 struct {
	Name          string
	SkuRef        string
	ZoneInitial   string
	ZoneUpdated   string
	BootDeviceRef string
}

func (p ComputeParamsV1) getParams() Params { return p.Params }

type StorageParamsV1 struct {
	Params
	Sku          StorageSkuParamsV1
	BlockStorage BlockStorageParamsV1
	Image        ImageParamsV1
}
type StorageSkuParamsV1 struct {
	Name          string
	Provider      string
	Tier          string
	Iops          int
	StorageType   string
	MinVolumeSize int
}
type BlockStorageParamsV1 struct {
	Name          string
	SkuRef        string
	SizeGBInitial int
	SizeGBUpdated int
}
type ImageParamsV1 struct {
	Name                   string
	BlockStorageRef        string
	CpuArchitectureInitial string
	CpuArchitectureUpdated string
}

func (p StorageParamsV1) getParams() Params { return p.Params }

// Responses

type roleResponseV1 struct {
	Metadata    metadataResponse
	Status      statusResponse
	Permissions []rolePermissionResponseV1
}
type rolePermissionResponseV1 struct {
	Provider  string
	Resources []string
	Verb      []string
}
type roleAssignmentResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Roles  []string
	Subs   []string
	Scopes []roleAssignmentScopeResponseV1
}
type roleAssignmentScopeResponseV1 struct {
	Tenants    []string
	Regions    []string
	Workspaces []string
}

type workspaceResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse
}

type instanceSkuResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Architecture string
	Provider     string
	Tier         string
	RAM          int
	VCPU         int
}
type instanceResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	SkuRef        string
	Zone          string
	BootDeviceRef string
}

type storageSkuResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Provider      string
	Tier          string
	Iops          int
	StorageType   string
	MinVolumeSize int
}
type blockStorageResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	SkuRef string
	SizeGB int
}
type imageResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	BlockStorageRef string
	CpuArchitecture string
}
