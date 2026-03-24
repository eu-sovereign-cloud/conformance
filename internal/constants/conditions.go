package constants

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

// Expected States
var (
	CreatedResourceExpectedStates = []schema.ResourceState{schema.ResourceStatePending, schema.ResourceStateCreating, schema.ResourceStateActive}
	UpdatedResourceExpectedStates = []schema.ResourceState{schema.ResourceStateActive, schema.ResourceStateUpdating}
)

// Conditions
var (
	ActiveMessage = "Resource is active"
	ActiveReason  = "active"
	ActiveType    = "active"

	CreatingMessage = "Resource is being created"
	CreatingReason  = "creating"
	CreatingType    = "creating"

	PendingMessage = "Resource is being pending"
	PendingReason  = "pending"
	PendingType    = "pending"

	UpdatingMessage = "Resource is being updated"
	UpdatingReason  = "updating"
	UpdatingType    = "updating"

	DeletingMessage = "Resource is being deleted"
	DeletingReason  = "deleting"
	DeletingType    = "deleting"

	ActiveCondition = schema.StatusCondition{
		Message: ActiveMessage,
		Reason:  ActiveReason,
		State:   schema.ResourceStateActive,
		Type:    ActiveType,
	}

	CreatingCondition = schema.StatusCondition{
		Message: CreatingMessage,
		Reason:  CreatingReason,
		State:   schema.ResourceStateCreating,
		Type:    CreatingType,
	}

	PendingCondition = schema.StatusCondition{
		Message: PendingMessage,
		Reason:  PendingReason,
		State:   schema.ResourceStatePending,
		Type:    PendingType,
	}

	UpdatingCondition = schema.StatusCondition{
		Message: UpdatingMessage,
		Reason:  UpdatingReason,
		State:   schema.ResourceStateUpdating,
		Type:    UpdatingType,
	}

	DeletingCondition = schema.StatusCondition{
		Message: DeletingMessage,
		Reason:  DeletingReason,
		State:   schema.ResourceStateDeleting,
		Type:    DeletingType,
	}

	GetConditionAfterCreating = []schema.StatusCondition{
		PendingCondition,
		CreatingCondition,
		ActiveCondition,
	}

	GetConditionAfterUpdating = []schema.StatusCondition{
		PendingCondition,
		CreatingCondition,
		ActiveCondition,
		UpdatingCondition,
		ActiveCondition,
	}

	GetConditionAfterDeleting = []schema.StatusCondition{
		PendingCondition,
		CreatingCondition,
		ActiveCondition,
		UpdatingCondition,
		ActiveCondition,
		DeletingCondition,
	}

	GetConditionAfterStopping = []schema.StatusCondition{
		PendingCondition,
		CreatingCondition,
		ActiveCondition,
		UpdatingCondition,
		ActiveCondition,
		UpdatingCondition,
		ActiveCondition,
		UpdatingCondition,
		ActiveCondition,
		UpdatingCondition,
		ActiveCondition,
	}

	GetConditionAfterStarting = []schema.StatusCondition{
		PendingCondition,
		CreatingCondition,
		ActiveCondition,
		UpdatingCondition,
		ActiveCondition,
		UpdatingCondition,
		ActiveCondition,
	}

	GetConditionAfterRestarting = []schema.StatusCondition{
		PendingCondition,
		CreatingCondition,
		ActiveCondition,
		UpdatingCondition,
		ActiveCondition,
		UpdatingCondition,
		ActiveCondition,
		UpdatingCondition,
		ActiveCondition,
	}
)
