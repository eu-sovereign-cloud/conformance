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

type NetworkParamsV1 struct {
	Params

	Network         NetworkInstanceParamsV1
	InternetGateway InternetGatewayParamsV1
	RouteTable      RouteTableParamsV1
	NIC             NICParamsV1
	PublicIP        PublicIPParamsV1
	SecurityGroup   SecurityGroupParamsV1
	Subnet          SubnetParamsV1
	NetworkSku      NetworkSkuParamsV1
	Instance        InstanceParamsV1
	InstanceSku     InstanceSkuParamsV1
}
type NetworkInstanceParamsV1 struct {
	Name            string
	Cidr            CIDR
	AdditionalCidrs []CIDR
	SkuRef          string
	RouteTableRef   string
}

type NetworkSkuParamsV1 struct {
	Name string
}

type InternetGatewayParamsV1 struct {
	Name       string
	EgressOnly bool
}

type RouteTableParamsV1 struct {
	Name     string
	LocalRef string
	Routes   []Routes
}

type NICParamsV1 struct {
	Name         string
	Addresses    []string
	PublicIpRefs []string
	SkuRef       string
	SubnetRef    string
}

type PublicIPParamsV1 struct {
	Name    string
	Version string
	Address string
}

type SecurityGroupParamsV1 struct {
	Name  string
	Rules []Rules
}

type SubnetParamsV1 struct {
	Name          string
	Cidr          CIDR
	Zone          string
	RouteTableRef string
	NetworkRef    string
	SkuRef        string
}

type CIDR struct {
	Ipv4 string
	Ipv6 string
}

type Routes struct {
	DestinationCidrBlock string
	TargetRef            string
}

type Rules struct {
	Direction  string
	Version    string
	Protocol   string
	Ports      Ports
	Icmp       Icmp
	sourceRefs []string
}
type Ports struct {
	From int
	To   int
	list []int
}

type Icmp struct {
	Type int
	Code int
}

func (p NetworkParamsV1) getParams() Params { return p.Params }

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

//network Responses

type networkResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Cidr            CIDR
	AdditionalCidrs []CIDR
	SkuRef          string
	RouteTableRef   string
}

type networkSkuResponseV1 struct {
	Metadata metadataResponse

	Bandwidth int
	Packets   int
}

type internetGatewayResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	EgressOnly bool
}

type routeTableResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	LocalRef string
	Routes   []Routes
}
type nicResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Addresses    []string
	PublicIpRefs []string
	SkuRef       string
	SubnetRef    string
}
type publicIPResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Version string
	Address string
}
type securityGroupResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Rules []Rules
}
type subnetResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Cidr          CIDR
	Zone          string
	RouteTableRef string
	NetworkRef    string
	SkuRef        string
}
