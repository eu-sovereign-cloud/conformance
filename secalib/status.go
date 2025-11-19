package secalib

import (
	"time"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"k8s.io/utils/ptr"
)

func SetResourceState(state string) *schema.ResourceState {
	return ptr.To(schema.ResourceState(state))
}

func addStatusCondition(conditions []schema.StatusCondition, state string) []schema.StatusCondition {
	return append(conditions, schema.StatusCondition{
		LastTransitionAt: time.Now(),
		State:            schema.ResourceState(state),
	})
}

func NewResourceStatus(state string) *schema.Status {
	return &schema.Status{
		State:      SetResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func SetStatusState(status *schema.Status, state string) {
	status.State = SetResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Worspace

func NewWorkspaceStatus(state string) *schema.WorkspaceStatus {
	return &schema.WorkspaceStatus{
		State:      SetResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func SetWorkspaceStatusState(status *schema.WorkspaceStatus, state string) {
	status.State = SetResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Storage

func NewBlockStorageStatus(state string) *schema.BlockStorageStatus {
	return &schema.BlockStorageStatus{
		State:      SetResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func SetBlockStorageStatusState(status *schema.BlockStorageStatus, state string) {
	status.State = SetResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewImageStatus(state string) *schema.ImageStatus {
	return &schema.ImageStatus{
		State:      SetResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func SetImageStatusState(status *schema.ImageStatus, state string) {
	status.State = SetResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Instance

func NewInstanceStatus(state string) *schema.InstanceStatus {
	return &schema.InstanceStatus{
		State:      SetResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func SetInstanceStatusState(status *schema.InstanceStatus, state string) {
	status.State = SetResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Network

func NewNetworkStatus(state string) *schema.NetworkStatus {
	return &schema.NetworkStatus{
		State:      SetResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func SetNetworkStatusState(status *schema.NetworkStatus, state string) {
	status.State = SetResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewRouteTableStatus(state string) *schema.RouteTableStatus {
	return &schema.RouteTableStatus{
		State:      SetResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func SetRouteTableStatusState(status *schema.RouteTableStatus, state string) {
	status.State = SetResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewSubnetStatus(state string) *schema.SubnetStatus {
	return &schema.SubnetStatus{
		State:      SetResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func SetSubnetStatusState(status *schema.SubnetStatus, state string) {
	status.State = SetResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewPublicIpStatus(state string) *schema.PublicIpStatus {
	return &schema.PublicIpStatus{
		State:      SetResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func SetPublicIpStatusState(status *schema.PublicIpStatus, state string) {
	status.State = SetResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewNicStatus(state string) *schema.NicStatus {
	return &schema.NicStatus{
		State:      SetResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func SetNicStatusState(status *schema.NicStatus, state string) {
	status.State = SetResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewSecurityGroupStatus(state string) *schema.SecurityGroupStatus {
	return &schema.SecurityGroupStatus{
		State:      SetResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func SetSecurityGroupStatusState(status *schema.SecurityGroupStatus, state string) {
	status.State = SetResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}
