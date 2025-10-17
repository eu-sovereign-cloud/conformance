package mock

import (
	"time"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"k8s.io/utils/ptr"
)

func setResourceState(state string) *schema.ResourceState {
	return ptr.To(schema.ResourceState(state))
}

func addStatusCondition(conditions []schema.StatusCondition, state string) []schema.StatusCondition {
	return append(conditions, schema.StatusCondition{
		LastTransitionAt: time.Now(),
		State:            schema.ResourceState(state),
	})
}

func newStatus(state string) *schema.Status {
	return &schema.Status{
		State:      setResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func setStatusState(status *schema.Status, state string) {
	status.State = setResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Worspace

func newWorkspaceStatus(state string) *schema.WorkspaceStatus {
	return &schema.WorkspaceStatus{
		State:      setResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func setWorkspaceStatusState(status *schema.WorkspaceStatus, state string) {
	status.State = setResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Storage

func newBlockStorageStatus(state string) *schema.BlockStorageStatus {
	return &schema.BlockStorageStatus{
		State:      setResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func setBlockStorageStatusState(status *schema.BlockStorageStatus, state string) {
	status.State = setResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newImageStatus(state string) *schema.ImageStatus {
	return &schema.ImageStatus{
		State:      setResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func setImageStatusState(status *schema.ImageStatus, state string) {
	status.State = setResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Instance

func newInstanceStatus(state string) *schema.InstanceStatus {
	return &schema.InstanceStatus{
		State:      setResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func setInstanceStatusState(status *schema.InstanceStatus, state string) {
	status.State = setResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

// Network

func newNetworkStatus(state string) *schema.NetworkStatus {
	return &schema.NetworkStatus{
		State:      setResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func setNetworkStatusState(status *schema.NetworkStatus, state string) {
	status.State = setResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newRouteTableStatus(state string) *schema.RouteTableStatus {
	return &schema.RouteTableStatus{
		State:      setResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func setRouteTableStatusState(status *schema.RouteTableStatus, state string) {
	status.State = setResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newSubnetStatus(state string) *schema.SubnetStatus {
	return &schema.SubnetStatus{
		State:      setResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func setSubnetStatusState(status *schema.SubnetStatus, state string) {
	status.State = setResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newPublicIpStatus(state string) *schema.PublicIpStatus {
	return &schema.PublicIpStatus{
		State:      setResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func setPublicIpStatusState(status *schema.PublicIpStatus, state string) {
	status.State = setResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newNicStatus(state string) *schema.NicStatus {
	return &schema.NicStatus{
		State:      setResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func setNicStatusState(status *schema.NicStatus, state string) {
	status.State = setResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}

func newSecurityGroupStatus(state string) *schema.SecurityGroupStatus {
	return &schema.SecurityGroupStatus{
		State:      setResourceState(state),
		Conditions: []schema.StatusCondition{},
	}
}

func setSecurityGroupStatusState(status *schema.SecurityGroupStatus, state string) {
	status.State = setResourceState(state)
	status.Conditions = addStatusCondition(status.Conditions, state)
}
