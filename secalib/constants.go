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

	regionResource          = "regions/%s"
	skuResource             = resourceTenantsPrefix + "/skus/%s"
	roleResource            = resourceTenantsPrefix + "/roles/%s"
	roleAssignmentResource  = resourceTenantsPrefix + "/role-assignments/%s"
	workspaceResource       = resourceTenantsPrefix + "/workspaces/%s"
	blockStorageResource    = resourceWorkspacesPrefix + "/block-storages/%s"
	imageResource           = resourceTenantsPrefix + "/images/%s"
	instanceResource        = resourceWorkspacesPrefix + "/instances/%s"
	networkResource         = resourceWorkspacesPrefix + "/networks/%s"
	internetGatewayResource = resourceWorkspacesPrefix + "/internet-gateways/%s"
	nicResource             = resourceWorkspacesPrefix + "/nics/%s"
	publicIpResource        = resourceWorkspacesPrefix + "/public-ips/%s"
	routeTableResource      = resourceWorkspacesPrefix + "/networks/%s/route-tables/%s"
	subnetResource          = resourceWorkspacesPrefix + "/networks/%s/subnets/%s"
	securityGroupResource   = resourceWorkspacesPrefix + "/security-groups/%s"

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
)
