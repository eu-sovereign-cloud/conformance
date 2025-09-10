package secalib

// Authorization

type RoleSpecV1 struct {
	Permissions []*RoleSpecPermissionV1
}
type RoleSpecPermissionV1 struct {
	Provider  string
	Resources []string
	Verb      []string
}

type RoleAssignmentSpecV1 struct {
	Roles  []string
	Subs   []string
	Scopes []*RoleAssignmentSpecScopeV1
}
type RoleAssignmentSpecScopeV1 struct {
	Tenants    []string
	Regions    []string
	Workspaces []string
}

// Workspace

type WorkspaceSpecV1 struct {
	Labels *[]Label
}

// Compute

type InstanceSpecV1 struct {
	SkuRef        string
	Zone          string
	BootDeviceRef string
}

// Storage

type BlockStorageSpecV1 struct {
	SkuRef string `yaml:"skuRef"`
	SizeGB int    `yaml:"sizeGB"`
}

type ImageSpecV1 struct {
	BlockStorageRef string `yaml:"blockStorageRef"`
	CpuArchitecture string `yaml:"cpuArchitecture"`
}

type NetworkSpecV1 struct {
	Cidr            *NetworkSpecCIDRV1   `yaml:"cidr"`
	AdditionalCidrs []*NetworkSpecCIDRV1 `yaml:"additionalCidrs"`
	SkuRef          string               `yaml:"skuRef"`
	RouteTableRef   string               `yaml:"routeTableRef"`
}
type NetworkSpecCIDRV1 struct {
	Ipv4 string `yaml:"ipv4"`
	Ipv6 string `yaml:"ipv6"`
}

type InternetGatewaySpecV1 struct {
	EgressOnly bool `yaml:"egressOnly"`
}

type RouteTableSpecV1 struct {
	LocalRef string               `yaml:"localRef"`
	Routes   []*RouteTableRouteV1 `yaml:"routes"`
}
type RouteTableRouteV1 struct {
	DestinationCidrBlock string `yaml:"destinationCidrBlock"`
	TargetRef            string `yaml:"targetRef"`
}

type NICSpecV1 struct {
	Addresses    []string `yaml:"addresses"`
	PublicIpRefs []string `yaml:"publicIpRefs"`
	SubnetRef    string   `yaml:"subnetRef"`
}

type PublicIpSpecV1 struct {
	Version string `yaml:"version"`
	Address string `yaml:"address"`
}

type SecurityGroupSpecV1 struct {
	Rules []*SecurityGroupRuleV1 `yaml:"rules"`
}
type SecurityGroupRuleV1 struct {
	Direction string `yaml:"direction"`
}

type SubnetSpecV1 struct {
	Cidr *SubnetSpecCIDRV1 `yaml:"cidr"`
	Zone string            `yaml:"zone"`
}
type SubnetSpecCIDRV1 struct {
	Ipv4 string `yaml:"ipv4"`
	Ipv6 string `yaml:"ipv6"`
}
