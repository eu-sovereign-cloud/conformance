package mock

// Params

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
	CreatedZone   string
	UpdatedZone   string
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
	CreatedSizeGB int
	UpdatedSizeGB int
}
type ImageParamsV1 struct {
	Name                   string
	BlockStorageRef        string
	CreatedCpuArchitecture string
	UpdatedCpuArchitecture string
}

func (p StorageParamsV1) getParams() Params { return p.Params }

// Responses

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
