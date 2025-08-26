package mock

// Params

type AuthorizationParamsV1 struct {
	Params
	role           RoleParamsV1
	roleAssignment RoleAssignmentParamsV1
}
type RoleAssignmentScopeParamsV1 struct {
	Tenants    []string
	Regions    []string
	Workspaces []string
}
type RoleParamsV1 struct {
	Name        string
	Permissions []RolePermissionParamsV1
}
type RolePermissionParamsV1 struct {
	Provider  string
	Resources []string
	Verbs     []string
}
type RoleAssignmentParamsV1 struct {
	Name   string
	roles  []string
	subs   []string
	scopes []RoleAssignmentScopeParamsV1
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
	BlockStorage    BlockStorageParamsV1
	StorageSku      StorageSkuParamsV1
}
type NetworkInstanceParamsV1 struct {
	Name            string
	Cidr            CIDR
	AdditionalCidrs []CIDR
	SkuRef          string
	RouteTableRef   string
}

type NetworkSkuParamsV1 struct {
	Name      string
	Bandwidth int
	Packets   int
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
	rules []Rules
}

type SubnetParamsV1 struct {
	Name string
	Cidr CIDR
	Zone string
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
	Direction string
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

type rolePermissionResponseV1 struct {
	Provider  string
	Resources []string
	Verbs     []string
}
type roleResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Permissions []rolePermissionResponseV1
}
type roleAssignmentScopeResponseV1 struct {
	Tenants    []string
	Regions    []string
	Workspaces []string
}
type roleAssignmentResponseV1 struct {
	Metadata metadataResponse
	Status   statusResponse

	Roles  []string
	Subs   []string
	Scopes []roleAssignmentScopeResponseV1
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

	Addresses []string
	SubnetRef string
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

	Cidr CIDR
	Zone string
}
