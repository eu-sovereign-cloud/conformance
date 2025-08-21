package secatest

const (
	// Step Parameters
	operationStepParameter = "operation"
	tenantStepParameter    = "tenant"
	workspaceStepParameter = "workspace"

	// Versions
	version1 = "v1"

	// Providers
	workspaceV1Provider = "seca.workspace/v1"
	storageV1Provider   = "seca.storage/v1"
	computeV1Provider   = "seca.compute/v1"

	// Labels
	architectureLabel = "architecture"
	providerLabel     = "provider"
	tierLabel         = "tier"

	// Resource URLs
	workspaceResource    = "tenants/%s/workspaces/%s"
	blockStorageResource = "tenants/%s/workspaces/%s/block-storages/%s"
	imageResource        = "tenants/%s/images/%s"
	instanceResource     = "tenants/%s/workspaces/%s/instances/%s"

	// Resource References
	skuRef           = "skus/%s"
	blockStoragesRef = "block-storages/%s"

	// Resource Kinds
	workspaceKind    = "workspace"
	storageSkuKind   = "storage-sku"
	blockStorageKind = "block-storage"
	imageKind        = "image"
	instanceSkuKind  = "instance-sku"
	instanceKind     = "instance"

	// Status States
	creatingStatusState  = "creating"
	activeStatusState    = "active"
	updatingStatusState  = "updating"
	suspendedStatusState = "suspended"

	// Resource Enums
	storageTypeLocalEphemeral = "local-ephemeral"
	cpuArchitectureAmd64      = "amd64"
	cpuArchitectureArm64      = "arm64"
)
