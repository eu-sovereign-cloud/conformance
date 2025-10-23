package secalib

const (
	// API Versions
	ApiVersion1 = "v1"

	// Labels
	ArchitectureLabel   = "architecture"
	ProviderLabel       = "provider"
	TierLabel           = "tier"
	EnvLabel            = "env"
	EnvDevelopmentLabel = "development"
	EnvProductionLabel  = "production"

	// URL Prefixes
	urlProvidersPrefix = "/providers/"

	// Resource URLs
	resourceTenantsPrefix    = "tenants/%s"
	resourceWorkspacesPrefix = resourceTenantsPrefix + "/workspaces/%s"

	SkuResource             = resourceTenantsPrefix + "/skus/%s"
	RoleResource            = resourceTenantsPrefix + "/roles/%s"
	RoleAssignmentResource  = resourceTenantsPrefix + "/role-assignments/%s"
	WorkspaceResource       = resourceTenantsPrefix + "/workspaces/%s"
	BlockStorageResource    = resourceWorkspacesPrefix + "/block-storages/%s"
	ImageResource           = resourceTenantsPrefix + "/images/%s"
	InstanceResource        = resourceWorkspacesPrefix + "/instances/%s"
	NetworkResource         = resourceWorkspacesPrefix + "/networks/%s"
	InternetGatewayResource = resourceWorkspacesPrefix + "/internet-gateways/%s"
	NicResource             = resourceWorkspacesPrefix + "/nics/%s"
	PublicIpResource        = resourceWorkspacesPrefix + "/public-ips/%s"
	RouteTableResource      = resourceWorkspacesPrefix + "/networks/%s/route-tables/%s"
	SubnetResource          = resourceWorkspacesPrefix + "/networks/%s/subnets/%s"
	SecurityGroupResource   = resourceWorkspacesPrefix + "/security-groups/%s"

	// Resource References
	SkuRef             = "skus/%s"
	InstanceRef        = "instances/%s"
	BlockStorageRef    = "block-storages/%s"
	InternetGatewayRef = "internet-gateways/%s"
	NetworkRef         = "networks/%s"
	RouteTableRef      = "route-tables/%s"
	SubnetRef          = "subnets/%s"
	PublicIpRef        = "public-ips/%s"

	// Resource Kinds
	RoleKind            = "role"
	RoleAssignmentKind  = "role-assignment"
	WorkspaceKind       = "workspace"
	StorageSkuKind      = "storage-sku"
	BlockStorageKind    = "block-storage"
	ImageKind           = "image"
	InstanceSkuKind     = "instance-sku"
	InstanceKind        = "instance"
	NetworkKind         = "network"
	NetworkSkuKind      = "network-sku"
	InternetGatewayKind = "internet-gateway"
	NicKind             = "nic"
	PublicIpKind        = "public-ip"
	RouteTableKind      = "routing-table"
	SubnetKind          = "subnet"
	SecurityGroupKind   = "security-group"

	// Resource Enums
	CpuArchitectureAmd64 = "amd64"
	CpuArchitectureArm64 = "arm64"

	IpVersion4 = "ipv4"
	IpVersion6 = "ipv6"

	SecurityRuleDirectionIngress = "ingress"
	SecurityRuleDirectionEgress  = "egress"

	// Status States
	CreatingResourceState  = "creating"
	ActiveResourceState    = "active"
	UpdatingResourceState  = "updating"
	SuspendedResourceState = "suspended"

	// Generators
	maxBlockStorageSize = 1000000 // GB
)
