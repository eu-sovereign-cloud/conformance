package mock

const (
	// Scenario States
	startedScenarioState    = "Started"
	creatingScenarioState   = "Creating"
	createdScenarioState    = "Created"
	updatingScenarioState   = "Updating"
	updatedScenarioState    = "Updated"
	deletingScenarioState   = "Deleting"
	redeletingScenarioState = "Re-Deleting"

	// Scenario Priorities
	defaultScenarioPriority = 1

	// Resource URLs
	workspaceResource      = "tenants/%s/workspaces/%s"
	instanceSkuResource    = "tenants/%s/workspaces/%s/skus/%s"
	instanceResource       = "tenants/%s/workspaces/%s/instances/%s"
	storageSkuResource     = "tenants/%s/workspaces/%s/skus/%s"
	blockStorageResource   = "tenants/%s/workspaces/%s/block-storages/%s"
	imageResource          = "tenants/%s/workspaces/%s/images/%s"
	rolesResource          = "tenants/%s/roles/%s"
	roleAssignmentResource = "tenants/%s/role-assignments/%s"

	// Resource Kinds
	workspaceKind      = "workspace"
	instanceSkuKind    = "instance-sku"
	instanceKind       = "instance"
	storageSkuKind     = "storage-sku"
	blockStorageKind   = "block-storage"
	imageKind          = "image"
	rolesKind          = "role"
	roleAssignmentKind = "role-assignment"

	// Status States
	creatingStatusState = "creating"
	activeStatusState   = "active"
	updatingStatusState = "updating"

	// Http Headers
	authorizationHttpHeaderKey         = "Authorization"
	authorizationHttpHeaderValuePrefix = "Bearer "

	contentTypeHttpHeaderKey   = "Content-Type"
	contentTypeHttpHeaderValue = "application/json"
)
