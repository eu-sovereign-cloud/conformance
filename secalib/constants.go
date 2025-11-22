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

	// URL Prefixes
	UrlProvidersPrefix = "/providers/"

	// Resource URLs
	resourceTenantsPrefix    = "tenants/%s"
	resourceWorkspacesPrefix = resourceTenantsPrefix + "/workspaces/%s"

	RegionResource          = "regions/%s"
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

	// Generators
	maxBlockStorageSize = 1000000 // GB

	// Zones
	ZoneA = "a"
	ZoneB = "b"

	// Labels
	LabelKeyTier         = "tier"
	LabelEnvKey          = "env"
	LabelEnvValue        = "test"
	LabelMonitoringValue = "monitoring"
	LabelAlertLevelValue = "alert-level"
	LabelHightValue      = "high"
	LabelTierKey         = "tier"
	LabelTierValue       = "backend"
	LabelVersion         = "version"
	LabelUptime          = "uptime"
	LabelLoad            = "load"
)
