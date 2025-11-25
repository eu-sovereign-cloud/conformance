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
	EnvConformance      = "conformance"

	// URL Prefixes
	UrlProvidersPrefix = "/providers/"

	// Resource URLs
	resourceTenantsPrefix    = "tenants/%s"
	resourceWorkspacesPrefix = resourceTenantsPrefix + "/workspaces/%s"

	RegionResource              = "regions/%s"
	SkuResource                 = resourceTenantsPrefix + "/skus/%s"
	SkuListResource             = resourceTenantsPrefix + "/skus/"
	RoleResource                = resourceTenantsPrefix + "/roles/%s"
	RolesResource               = resourceTenantsPrefix + "/roles"
	RoleAssignmentResource      = resourceTenantsPrefix + "/role-assignments/%s"
	RoleAssignmentsResource     = resourceTenantsPrefix + "/role-assignments"
	WorkspaceResource           = resourceTenantsPrefix + "/workspaces/%s"
	WorkspaceListResource       = resourceTenantsPrefix + "/workspaces"
	BlockStorageResource        = resourceWorkspacesPrefix + "/block-storages/%s"
	BlockStorageListResource    = resourceWorkspacesPrefix + "/block-storages"
	ImageResource               = resourceTenantsPrefix + "/images/%s"
	ImageListResource           = resourceTenantsPrefix + "/images"
	InstanceResource            = resourceWorkspacesPrefix + "/instances/%s"
	InstanceListResource        = resourceWorkspacesPrefix + "/instances"
	NetworkResource             = resourceWorkspacesPrefix + "/networks/%s"
	NetworkListResource         = resourceWorkspacesPrefix + "/networks"
	InternetGatewayResource     = resourceWorkspacesPrefix + "/internet-gateways/%s"
	InternetGatewayListResource = resourceWorkspacesPrefix + "/internet-gateways"
	NicResource                 = resourceWorkspacesPrefix + "/nics/%s"
	NicListResource             = resourceWorkspacesPrefix + "/nics"
	PublicIpResource            = resourceWorkspacesPrefix + "/public-ips/%s"
	PublicIpListResource        = resourceWorkspacesPrefix + "/public-ips"
	RouteTableResource          = resourceWorkspacesPrefix + "/networks/%s/route-tables/%s"
	RouteTableListResource      = resourceWorkspacesPrefix + "/networks/%s/route-tables"
	SubnetResource              = resourceWorkspacesPrefix + "/networks/%s/subnets/%s"
	SubnetListResource          = resourceWorkspacesPrefix + "/networks/%s/subnets"
	SecurityGroupResource       = resourceWorkspacesPrefix + "/security-groups/%s"
	SecurityGroupListResource   = resourceWorkspacesPrefix + "/security-groups"

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
