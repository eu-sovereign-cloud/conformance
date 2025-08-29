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

type WorkspaceSpecV1 struct{}

// Compute

type InstanceSpecV1 struct {
	SkuRef        string
	Zone          string
	BootDeviceRef string
}

// Storage

type BlockStorageSpecV1 struct {
	SkuRef string
	SizeGB int
}

type ImageSpecV1 struct {
	BlockStorageRef string
	CpuArchitecture string
}

type NetworkSpecV1 struct {
	Cidr            *NetworkSpecCIDRV1
	AdditionalCidrs []*NetworkSpecCIDRV1
	SkuRef          string
	RouteTableRef   string
}
type NetworkSpecCIDRV1 struct {
	Ipv4 string
	Ipv6 string
}

type InternetGatewaySpecV1 struct {
	EgressOnly bool
}

type RouteTableSpecV1 struct {
	LocalRef string
	Routes   []*RouteTableRouteV1
}
type RouteTableRouteV1 struct {
	DestinationCidrBlock string
	TargetRef            string
}

type NICSpecV1 struct {
	Addresses    []string
	PublicIpRefs []string
	SubnetRef    string
}

type PublicIPSpecV1 struct {
	Version string
	Address string
}

type SecurityGroupSpecV1 struct {
	Rules []*SecurityGroupRule
}
type SecurityGroupRule struct {
	Direction string
	Protocol  string
	Port      int
}

type SubnetSpecV1 struct {
	Cidr *SubnetSpecCIDRV1
	Zone string
}
type SubnetSpecCIDRV1 struct {
	Ipv4 string
	Ipv6 string
}
