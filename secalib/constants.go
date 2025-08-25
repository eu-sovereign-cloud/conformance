package secalib

const (
	// API Versions
	ApiVersion1 = "v1"

	// Providers
	WorkspaceV1Provider = "seca.workspace/v1"
	StorageV1Provider   = "seca.storage/v1"
	ComputeV1Provider   = "seca.compute/v1"
	NetworkV1Provider   = "seca.network/v1"
	// Labels
	ArchitectureLabel = "architecture"
	ProviderLabel     = "provider"
	TierLabel         = "tier"

	// Resource URLs
	SkuResource             = "tenants/%s/skus/%s"
	RoleResource            = "tenants/%s/roles/%s"
	RoleAssignmentResource  = "tenants/%s/role-assignments/%s"
	WorkspaceResource       = "tenants/%s/workspaces/%s"
	BlockStorageResource    = "tenants/%s/workspaces/%s/block-storages/%s"
	ImageResource           = "tenants/%s/images/%s"
	InstanceResource        = "tenants/%s/workspaces/%s/instances/%s"
	NetworkResource         = "tenants/%s/workspaces/%s/networks/%s"
	NetworkSkuResource      = "tenants/%s/workspaces/%s/skus/%s"
	InternetGatewayResource = "tenants/%s/workspaces/%s/internet-gateways/%s"
	NicResource             = "tenants/%s/workspaces/%s/nics/%s"
	PublicIPResource        = "tenants/%s/workspaces/%s/public-ips/%s"
	RouteTableResource      = "tenants/%s/workspaces/%s/route-tables/%s"
	SubnetResource          = "tenants/%s/workspaces/%s/subnets/%s"
	SecurityGroupResource   = "tenants/%s/workspaces/%s/security-groups/%s"

	// Resource References
	SkuRef          = "skus/%s"
	BlockStorageRef = "block-storages/%s"

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
	PublicIPKind        = "public-ip"
	RouteTableKind      = "routing-table"
	SubnetKind          = "subnet"
	SecurityGroupKind   = "security-group"

	// Resource Enums
	StorageTypeLocalEphemeral = "local-ephemeral"

	CpuArchitectureAmd64 = "amd64"
	CpuArchitectureArm64 = "arm64"

	// Status States
	CreatingStatusState  = "creating"
	ActiveStatusState    = "active"
	UpdatingStatusState  = "updating"
	SuspendedStatusState = "suspended"
)
