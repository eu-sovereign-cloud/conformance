package mock

const (
	// Scenario Priorities
	startedScenarioState    = "Started"
	defaultScenarioPriority = 1

	// Resource URLs
	skuResource          = "tenants/%s/skus/%s"
	workspaceResource    = "tenants/%s/workspaces/%s"
	instanceResource     = "tenants/%s/workspaces/%s/instances/%s"
	blockStorageResource = "tenants/%s/workspaces/%s/block-storages/%s"
	imageResource        = "tenants/%s/images/%s"

	// Resource Kinds
	workspaceKind    = "workspace"
	instanceSkuKind  = "instance-sku"
	instanceKind     = "instance"
	storageSkuKind   = "storage-sku"
	blockStorageKind = "block-storage"
	imageKind        = "image"

	// Status States
	creatingStatusState  = "creating"
	activeStatusState    = "active"
	updatingStatusState  = "updating"
	suspendedStatusState = "suspended"

	// Http Headers
	authorizationHttpHeaderKey         = "Authorization"
	authorizationHttpHeaderValuePrefix = "Bearer "
	contentTypeHttpHeaderKey           = "Content-Type"
	contentTypeHttpHeaderValue         = "application/json"
)
