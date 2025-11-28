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

	// URL Prefixes
	urlProvidersPrefix = "/providers/"

	// Resource URLs
	resourceTenantsPrefix    = "tenants/%s"
	resourceWorkspacesPrefix = resourceTenantsPrefix + "/workspaces/%s"

	regionResource              = "regions/%s"
	skuResource                 = resourceTenantsPrefix + "/skus/%s"
	skuListResource             = resourceTenantsPrefix + "/skus/"
	roleResource                = resourceTenantsPrefix + "/roles/%s"
	rolesResource               = resourceTenantsPrefix + "/roles"
	roleAssignmentResource      = resourceTenantsPrefix + "/role-assignments/%s"
	roleAssignmentsResource     = resourceTenantsPrefix + "/role-assignments"
	workspaceResource           = resourceTenantsPrefix + "/workspaces/%s"
	workspaceListResource       = resourceTenantsPrefix + "/workspaces"
	blockStorageResource        = resourceWorkspacesPrefix + "/block-storages/%s"
	blockStorageListResource    = resourceWorkspacesPrefix + "/block-storages"
	imageResource               = resourceTenantsPrefix + "/images/%s"
	imageListResource           = resourceTenantsPrefix + "/images"
	instanceResource            = resourceWorkspacesPrefix + "/instances/%s"
	instanceListResource        = resourceWorkspacesPrefix + "/instances"
	networkResource             = resourceWorkspacesPrefix + "/networks/%s"
	networkListResource         = resourceWorkspacesPrefix + "/networks"
	internetGatewayResource     = resourceWorkspacesPrefix + "/internet-gateways/%s"
	internetGatewayListResource = resourceWorkspacesPrefix + "/internet-gateways"
	nicResource                 = resourceWorkspacesPrefix + "/nics/%s"
	nicListResource             = resourceWorkspacesPrefix + "/nics"
	publicIpResource            = resourceWorkspacesPrefix + "/public-ips/%s"
	publicIpListResource        = resourceWorkspacesPrefix + "/public-ips"
	routeTableResource          = resourceWorkspacesPrefix + "/networks/%s/route-tables/%s"
	routeTableListResource      = resourceWorkspacesPrefix + "/networks/%s/route-tables"
	subnetResource              = resourceWorkspacesPrefix + "/networks/%s/subnets/%s"
	subnetListResource          = resourceWorkspacesPrefix + "/networks/%s/subnets"
	securityGroupResource       = resourceWorkspacesPrefix + "/security-groups/%s"
	securityGroupListResource   = resourceWorkspacesPrefix + "/security-groups"

	// Resource References
	skuRef             = "skus/%s"
	instanceRef        = "instances/%s"
	blockStorageRef    = "block-storages/%s"
	internetGatewayRef = "internet-gateways/%s"
	networkRef         = "networks/%s"
	routeTableRef      = "route-tables/%s"
	subnetRef          = "subnets/%s"
	publicIpRef        = "public-ips/%s"

	// Generators
	maxBlockStorageSize = 1000000 // GB

	// Labels
	EnvLabel            = "env"
	EnvDevelopmentLabel = "development"
	EnvConformanceLabel = "conformance"
	EnvProductionLabel  = "production"
)
