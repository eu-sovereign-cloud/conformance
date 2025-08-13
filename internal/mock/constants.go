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

	// Base URLs
	workspaceV1BaseURL = "/providers/seca.workspace/v1/tenants/%s/workspaces/%s"

	// Versions
	version1 = "v1"

	// Resource URLs
	workspaceResource = "tenants/%s/workspaces/%s"

	// Providers
	workspaceV1Provider = "seca.workspace/v1"

	// Kinds
	workspaceKind = "workspace"

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
