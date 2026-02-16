package types

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

type StatusType interface {
	schema.Status |
		schema.WorkspaceStatus |
		schema.BlockStorageStatus |
		schema.ImageStatus |
		schema.InstanceStatus |
		schema.NetworkStatus |
		schema.SubnetStatus |
		schema.RouteTableStatus |
		schema.NicStatus |
		schema.PublicIpStatus |
		schema.SecurityGroupStatus
}

func GetStatusState[S StatusType](status *S) *schema.ResourceState {
	if status == nil {
		return nil
	}

	switch v := any(*status).(type) {
	case schema.Status:
		return v.State
	case schema.WorkspaceStatus:
		return v.State
	case schema.BlockStorageStatus:
		return v.State
	case schema.ImageStatus:
		return v.State
	case schema.InstanceStatus:
		return v.State
	case schema.NetworkStatus:
		return v.State
	case schema.SubnetStatus:
		return v.State
	case schema.RouteTableStatus:
		return v.State
	case schema.NicStatus:
		return v.State
	case schema.PublicIpStatus:
		return v.State
	case schema.SecurityGroupStatus:
		return v.State
	default:
		return nil
	}
}
