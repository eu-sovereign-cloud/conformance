package mock

// Params

type WorkspaceParamsV1 struct {
	Params
	Name string
}

func (p WorkspaceParamsV1) getParams() Params { return p.Params }

type ComputeParamsV1 struct {
	Params
	Sku      InstanceSkuParamsV1
	Instance InstanceParamsV1
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
	Zone          string
	BootDeviceRef string
}

func (p ComputeParamsV1) getParams() Params { return p.Params }

type StorageSkuParamsV1 struct {
	Params
	Sku          SkuParamsV1
	BlockStorage BlockStorageParamsV1
	Image        ImageParamsV1
}
type SkuParamsV1 struct {
	Provider      string
	Tier          string
	Iops          int
	StorageType   string
	MinVolumeSize int
}
type BlockStorageParamsV1 struct {
	SkuRef string
	SizeGB int
}
type ImageParamsV1 struct {
	BlockStorageRef string
	CpuArchitecture string
}

func (p StorageSkuParamsV1) getParams() Params { return p.Params }

type AuthorizationParamsV1 struct {
	Params
	roles          RolesParamsV1
	roleAssignment RoleAssignmentParamsV1
}

type RolesParamsV1 struct {
	Name        string
	Permissions []PermissionsParamsV1
}

type RoleAssignmentParamsV1 struct {
	Name   string
	roles  []string
	subs   []string
	scopes Scopes
}

func (p AuthorizationParamsV1) getParams() Params { return p.Params }

// Responses

type workspaceResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse
}

type instanceSkuResponseV1 struct {
	metadata metadataResponse
	status   statusResponse

	architecture string
	provider     string
	tier         string
	ram          int
	vCPU         int
}

type instanceResponseV1 struct {
	metadata metadataResponse
	status   statusResponse

	skuRef        string
	zone          string
	bootDeviceRef string
}

type storageSkuResponseV1 struct {
	metadata metadataResponse
	status   statusResponse

	provider      string
	tier          string
	iops          int
	storageType   string
	minVolumeSize int
}

type blockStorageResponseV1 struct {
	metadata metadataResponse
	status   statusResponse

	skuRef string
	sizeGB int
}

type imageResponseV1 struct {
	metadata metadataResponse
	status   statusResponse

	blockStorageRef string
	cpuArchitecture string
}

type rolesResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Permissions []PermissionsParamsV1
}
type PermissionsParamsV1 struct {
	Provider  string
	Resources []string
	Verbs     []string
}

type roleAssignmentResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Roles  []string
	Subs   []string
	Scopes []Scopes
}

type Scopes struct {
	Tenants    []string
	Regions    []string
	Workspaces []string
}
