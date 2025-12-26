package mock

const (
	// Scenario Priorities
	startedScenarioState    = "Started"
	defaultScenarioPriority = 1

	// Http Headers
	authorizationHttpHeaderKey         = "Authorization"
	authorizationHttpHeaderValuePrefix = "Bearer "
	contentTypeHttpHeaderKey           = "Content-Type"
	contentTypeHttpHeaderValue         = "application/json"

	// Test Data

	// Providers
	authorizationProvider = "seca.authorization"
	regionProvider        = "seca.region"
	workspaceProvider     = "seca.workspace"
	storageProvider       = "seca.storage"
	computeProvider       = "seca.compute"
	networkProvider       = "seca.network"

	// Zones
	zoneA = "a"
	zoneB = "b"
)
