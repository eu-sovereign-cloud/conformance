package secalib

const (
	// API Versions
	ApiVersion1 = "v1"

	// Providers
	AuthorizationProvider = "seca.authorization"
	WorkspaceProvider     = "seca.workspace"
	StorageProvider       = "seca.storage"
	ComputeProvider       = "seca.compute"
	NetworkProvider       = "seca.network"

	// Labels
	ArchitectureLabel   = "architecture"
	ProviderLabel       = "provider"
	TierLabel           = "tier"
	EnvLabel            = "env"
	EnvDevelopmentLabel = "development"
	EnvProductionLabel  = "production"

	// Resource URLs
	resourceTenantsPrefix    = "tenants/%s"
	resourceWorkspacesPrefix = resourceTenantsPrefix + "/workspaces/%s"

	SkuResource             = resourceTenantsPrefix + "/skus/%s"
	RoleResource            = resourceTenantsPrefix + "/roles/%s"
	RoleAssignmentResource  = resourceTenantsPrefix + "/role-assignments/%s"
	RegionResource          = "regions/%s"
	WorkspaceResource       = resourceTenantsPrefix + "/workspaces/%s"
	BlockStorageResource    = resourceWorkspacesPrefix + "/block-storages/%s"
	ImageResource           = resourceTenantsPrefix + "/images/%s"
	InstanceResource        = resourceWorkspacesPrefix + "/instances/%s"
	NetworkResource         = resourceWorkspacesPrefix + "/networks/%s"
	InternetGatewayResource = resourceWorkspacesPrefix + "/internet-gateways/%s"
	NicResource             = resourceWorkspacesPrefix + "/nics/%s"
	PublicIPResource        = resourceWorkspacesPrefix + "/public-ips/%s"
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
	PublicIPRef        = "public-ips/%s"

	// Resource Kinds
	RoleKind            = "role"
	RoleAssignmentKind  = "role-assignment"
	RegionKind          = "region"
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
	PublicIPKind        = "public-ip"
	RouteTableKind      = "routing-table"
	SubnetKind          = "subnet"
	SecurityGroupKind   = "security-group"

	// Resource Enums
	CpuArchitectureAmd64 = "amd64"
	CpuArchitectureArm64 = "arm64"

	IPVersion4 = "ipv4"
	IPVersion6 = "ipv6"

	SecurityRuleDirectionIngress = "ingress"
	SecurityRuleDirectionEgress  = "egress"

	// Status States
	CreatingStatusState  = "creating"
	ActiveStatusState    = "active"
	UpdatingStatusState  = "updating"
	SuspendedStatusState = "suspended"

	// Generators
	maxBlockStorageSize = 1000000 // GB
)
