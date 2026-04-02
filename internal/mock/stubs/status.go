package stubs

import (
	"time"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func addStatusCondition(conditions []schema.StatusCondition, state schema.ResourceState) []schema.StatusCondition {
	return append(conditions, schema.StatusCondition{
		LastTransitionAt: time.Now(),
		State:            state,
	})
}

func newResourceStatus(state schema.ResourceState) *schema.Status {
	return &schema.Status{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdateResourceStatus() *schema.Status {
	return &schema.Status{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setResourceState(status *schema.Status, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Worspace

func newWorkspaceStatus(state schema.ResourceState) *schema.WorkspaceStatus {
	return &schema.WorkspaceStatus{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdateWorkspaceStatus() *schema.WorkspaceStatus {
	return &schema.WorkspaceStatus{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setWorkspaceState(status *schema.WorkspaceStatus, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Storage

func newBlockStorageStatus(state schema.ResourceState) *schema.BlockStorageStatus {
	return &schema.BlockStorageStatus{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdateBlockStorageStatus() *schema.BlockStorageStatus {
	return &schema.BlockStorageStatus{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setBlockStorageState(status *schema.BlockStorageStatus, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newImageStatus(state schema.ResourceState) *schema.ImageStatus {
	return &schema.ImageStatus{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdateImageStatus() *schema.ImageStatus {
	return &schema.ImageStatus{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setImageState(status *schema.ImageStatus, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Instance

func newInstanceStatus(state schema.ResourceState) *schema.InstanceStatus {
	return &schema.InstanceStatus{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdateInstanceStatus() *schema.InstanceStatus {
	return &schema.InstanceStatus{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setInstanceState(status *schema.InstanceStatus, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Network

func newNetworkStatus(state schema.ResourceState) *schema.NetworkStatus {
	return &schema.NetworkStatus{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdateNetworkStatus() *schema.NetworkStatus {
	return &schema.NetworkStatus{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setNetworkState(status *schema.NetworkStatus, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newRouteTableStatus(state schema.ResourceState) *schema.RouteTableStatus {
	return &schema.RouteTableStatus{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdateRouteTableStatus() *schema.RouteTableStatus {
	return &schema.RouteTableStatus{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setRouteTableState(status *schema.RouteTableStatus, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newSubnetStatus(state schema.ResourceState) *schema.SubnetStatus {
	return &schema.SubnetStatus{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdateSubnetStatus() *schema.SubnetStatus {
	return &schema.SubnetStatus{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setSubnetState(status *schema.SubnetStatus, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newPublicIpStatus(state schema.ResourceState) *schema.PublicIpStatus {
	return &schema.PublicIpStatus{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdatePublicIpStatus() *schema.PublicIpStatus {
	return &schema.PublicIpStatus{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setPublicIpState(status *schema.PublicIpStatus, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newNicStatus(state schema.ResourceState) *schema.NicStatus {
	return &schema.NicStatus{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdateNicStatus() *schema.NicStatus {
	return &schema.NicStatus{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setNicState(status *schema.NicStatus, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newSecurityGroupRuleStatus(state schema.ResourceState) *schema.SecurityGroupRuleStatus {
	return &schema.SecurityGroupRuleStatus{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdateSecurityGroupRuleStatus() *schema.SecurityGroupRuleStatus {
	return &schema.SecurityGroupRuleStatus{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setSecurityGroupRuleState(status *schema.SecurityGroupRuleStatus, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newSecurityGroupStatus(state schema.ResourceState) *schema.SecurityGroupStatus {
	return &schema.SecurityGroupStatus{
		State: state,
		Conditions: []schema.StatusCondition{
			constants.PendingCondition,
		},
	}
}

func beforeUpdateSecurityGroupStatus() *schema.SecurityGroupStatus {
	return &schema.SecurityGroupStatus{
		State:      schema.ResourceStateActive,
		Conditions: constants.GetConditionBeforeUpdating,
	}
}

func setSecurityGroupState(status *schema.SecurityGroupStatus, state schema.ResourceState) {
	status.State = state
	status.Conditions = addStatusCondition(status.Conditions, state)
}
