package secalib

import (
	"time"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func addStatusCondition(conditions []schema.StatusCondition, state schema.ResourceState) []schema.StatusCondition {
	return append(conditions, schema.StatusCondition{
		LastTransitionAt: time.Now(),
		State:            state,
	})
}

func NewResourceStatus(state schema.ResourceState) *schema.Status {
	return &schema.Status{
		State:      &state,
		Conditions: []schema.StatusCondition{},
	}
}

func SetStatusState(status *schema.Status, state schema.ResourceState) {
	status.State = &state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Worspace

func NewWorkspaceStatus(state schema.ResourceState) *schema.WorkspaceStatus {
	return &schema.WorkspaceStatus{
		State:      &state,
		Conditions: []schema.StatusCondition{},
	}
}

func SetWorkspaceStatusState(status *schema.WorkspaceStatus, state schema.ResourceState) {
	status.State = &state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Storage

func NewBlockStorageStatus(state schema.ResourceState) *schema.BlockStorageStatus {
	return &schema.BlockStorageStatus{
		State:      &state,
		Conditions: []schema.StatusCondition{},
	}
}

func SetBlockStorageStatusState(status *schema.BlockStorageStatus, state schema.ResourceState) {
	status.State = &state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewImageStatus(state schema.ResourceState) *schema.ImageStatus {
	return &schema.ImageStatus{
		State:      &state,
		Conditions: []schema.StatusCondition{},
	}
}

func SetImageStatusState(status *schema.ImageStatus, state schema.ResourceState) {
	status.State = &state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Instance

func NewInstanceStatus(state schema.ResourceState) *schema.InstanceStatus {
	return &schema.InstanceStatus{
		State:      &state,
		Conditions: []schema.StatusCondition{},
	}
}

func SetInstanceStatusState(status *schema.InstanceStatus, state schema.ResourceState) {
	status.State = &state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Network

func NewNetworkStatus(state schema.ResourceState) *schema.NetworkStatus {
	return &schema.NetworkStatus{
		State:      &state,
		Conditions: []schema.StatusCondition{},
	}
}

func SetNetworkStatusState(status *schema.NetworkStatus, state schema.ResourceState) {
	status.State = &state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewRouteTableStatus(state schema.ResourceState) *schema.RouteTableStatus {
	return &schema.RouteTableStatus{
		State:      &state,
		Conditions: []schema.StatusCondition{},
	}
}

func SetRouteTableStatusState(status *schema.RouteTableStatus, state schema.ResourceState) {
	status.State = &state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewSubnetStatus(state schema.ResourceState) *schema.SubnetStatus {
	return &schema.SubnetStatus{
		State:      &state,
		Conditions: []schema.StatusCondition{},
	}
}

func SetSubnetStatusState(status *schema.SubnetStatus, state schema.ResourceState) {
	status.State = &state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewPublicIpStatus(state schema.ResourceState) *schema.PublicIpStatus {
	return &schema.PublicIpStatus{
		State:      &state,
		Conditions: []schema.StatusCondition{},
	}
}

func SetPublicIpStatusState(status *schema.PublicIpStatus, state schema.ResourceState) {
	status.State = &state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewNicStatus(state schema.ResourceState) *schema.NicStatus {
	return &schema.NicStatus{
		State:      &state,
		Conditions: []schema.StatusCondition{},
	}
}

func SetNicStatusState(status *schema.NicStatus, state schema.ResourceState) {
	status.State = &state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func NewSecurityGroupStatus(state schema.ResourceState) *schema.SecurityGroupStatus {
	return &schema.SecurityGroupStatus{
		State:      &state,
		Conditions: []schema.StatusCondition{},
	}
}

func SetSecurityGroupStatusState(status *schema.SecurityGroupStatus, state schema.ResourceState) {
	status.State = &state
	status.Conditions = addStatusCondition(status.Conditions, state)
}
