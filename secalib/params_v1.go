package secalib

type ParamsFileV1 struct {
	General                  *GeneralParams                  `yaml:"general"`
	WorkspaceLifeCycleV1     *WorkspaceLifeCycleParamsV1     `yaml:"scenarios:workspaceLifeCycleV1"`
	AuthorizationLifeCycleV1 *AuthorizationLifeCycleParamsV1 `yaml:"scenarios:authorizationLifeCycleV1"`
	ComputeLifeCycleV1       *ComputeLifeCycleParamsV1       `yaml:"scenarios:computeLifeCycleV1"`
	StorageLifeCycleV1       *StorageLifeCycleParamsV1       `yaml:"scenarios:storageLifeCycleV1"`
	NetworkLifeCycleV1       *NetworkLifeCycleParamsV1       `yaml:"scenarios:networkLifeCycleV1"`
	Usage                    *UsageParamsV1                  `yaml:"scenarios:usageV1"`
}

type AuthorizationLifeCycleParamsV1 struct {
	Role           *ResourceParams[RoleSpecV1]           `yaml:"role"`
	RoleAssignment *ResourceParams[RoleAssignmentSpecV1] `yaml:"roleAssignment"`
}

type WorkspaceLifeCycleParamsV1 struct {
	Workspace *ResourceParams[WorkspaceSpecV1] `yaml:"workspace"`
}

type ComputeLifeCycleParamsV1 struct {
	Workspace *ResourceParams[WorkspaceSpecV1] `yaml:"workspace"`

	BlockStorage *ResourceParams[BlockStorageSpecV1] `yaml:"blockStorage"`
	Instance     *ResourceParams[InstanceSpecV1]     `yaml:"instance"`
}

type StorageLifeCycleParamsV1 struct {
	Workspace *ResourceParams[WorkspaceSpecV1] `yaml:"workspace"`

	BlockStorage *ResourceParams[BlockStorageSpecV1] `yaml:"blockStorage"`
	Image        *ResourceParams[ImageSpecV1]        `yaml:"image"`
}

type NetworkLifeCycleParamsV1 struct {
	Workspace *ResourceParams[WorkspaceSpecV1] `yaml:"workspace"`

	BlockStorage *ResourceParams[BlockStorageSpecV1] `yaml:"blockStorage"`

	Instance *ResourceParams[InstanceSpecV1] `yaml:"instance"`

	Network         *ResourceParams[NetworkSpecV1]         `yaml:"network"`
	InternetGateway *ResourceParams[InternetGatewaySpecV1] `yaml:"internetGateway"`
	RouteTable      *ResourceParams[RouteTableSpecV1]      `yaml:"routeTable"`
	Subnet          *ResourceParams[SubnetSpecV1]          `yaml:"subnet"`
	Nic             *ResourceParams[NICSpecV1]             `yaml:"nic"`
	PublicIp        *ResourceParams[PublicIpSpecV1]        `yaml:"publicIP"`
	SecurityGroup   *ResourceParams[SecurityGroupSpecV1]   `yaml:"securityGroup"`
}

type UsageParamsV1 struct {
	Role           *ResourceParams[RoleSpecV1]           `yaml:"role"`
	RoleAssignment *ResourceParams[RoleAssignmentSpecV1] `yaml:"roleAssignment"`

	Workspace *ResourceParams[WorkspaceSpecV1] `yaml:"workspace"`

	BlockStorage *ResourceParams[BlockStorageSpecV1] `yaml:"blockStorage"`
	Image        *ResourceParams[ImageSpecV1]        `yaml:"image"`

	Network         *ResourceParams[NetworkSpecV1]         `yaml:"network"`
	InternetGateway *ResourceParams[InternetGatewaySpecV1] `yaml:"internetGateway"`
	RouteTable      *ResourceParams[RouteTableSpecV1]      `yaml:"routeTable"`
	Subnet          *ResourceParams[SubnetSpecV1]          `yaml:"subnet"`
	Nic             *ResourceParams[NICSpecV1]             `yaml:"nic"`
	PublicIp        *ResourceParams[PublicIpSpecV1]        `yaml:"publicIP"`
	SecurityGroup   *ResourceParams[SecurityGroupSpecV1]   `yaml:"securityGroup"`

	Instance *ResourceParams[InstanceSpecV1] `yaml:"instance"`
}
