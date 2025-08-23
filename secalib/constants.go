package secalib

const (
	// API Versions
	ApiVersion1 = "v1"

	// Providers
	WorkspaceV1Provider = "seca.workspace/v1"
	StorageV1Provider   = "seca.storage/v1"
	ComputeV1Provider   = "seca.compute/v1"

	// Labels
	ArchitectureLabel = "architecture"
	ProviderLabel     = "provider"
	TierLabel         = "tier"

	// Resource URLs
	SkuResource            = "tenants/%s/skus/%s"
	RoleResource           = "tenants/%s/roles/%s"
	RoleAssignmentResource = "tenants/%s/role-assignments/%s"
	WorkspaceResource      = "tenants/%s/workspaces/%s"
	BlockStorageResource   = "tenants/%s/workspaces/%s/block-storages/%s"
	ImageResource          = "tenants/%s/images/%s"
	InstanceResource       = "tenants/%s/workspaces/%s/instances/%s"

	// Resource References
	SkuRef          = "skus/%s"
	BlockStorageRef = "block-storages/%s"

	// Resource Kinds
	RoleKind           = "role"
	RoleAssignmentKind = "role-assignment"
	WorkspaceKind      = "workspace"
	StorageSkuKind     = "storage-sku"
	BlockStorageKind   = "block-storage"
	ImageKind          = "image"
	InstanceSkuKind    = "instance-sku"
	InstanceKind       = "instance"

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
