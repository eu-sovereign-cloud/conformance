package constants

import (
	"time"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Expected States
var (
	CreatedResourceExpectedStates = []schema.ResourceState{schema.ResourceStatePending, schema.ResourceStateCreating, schema.ResourceStateActive}
	UpdatedResourceExpectedStates = []schema.ResourceState{schema.ResourceStateActive, schema.ResourceStateUpdating}
)

// Conditions
var (
	ActiveMessage = "Resource is active and ready."
	ActiveReason  = "active"
	ActiveType    = "active"

	CreatingMessage = "Resource is being created."
	CreatingReason  = "creating"
	CreatingType    = "creating"

	PendingMessage = "Resource is pending initialization."
	PendingReason  = "pending"
	PendingType    = "pending"

	UpdatingMessage = "Resource is being updated."
	UpdatingReason  = "updating"
	UpdatingType    = "updating"

	DeletingMessage = "Resource is being deleting."
	DeletingReason  = "deleting"
	DeletingType    = "deleting"

	ActiveCondition = schema.StatusCondition{
		LastTransitionAt: time.Now(),
		State:            schema.ResourceStateActive,
	}

	CreatingCondition = schema.StatusCondition{
		LastTransitionAt: time.Now(),
		State:            schema.ResourceStateCreating,
	}

	PendingCondition = schema.StatusCondition{
		LastTransitionAt: time.Now(),
		State:            schema.ResourceStatePending,
	}

	UpdatingCondition = schema.StatusCondition{
		LastTransitionAt: time.Now(),
		State:            schema.ResourceStateUpdating,
	}

	DeletingCondition = schema.StatusCondition{
		LastTransitionAt: time.Now(),
		State:            schema.ResourceStateDeleting,
	}

	GetConditionAfterCreating = []schema.StatusCondition{
		PendingCondition,
		CreatingCondition,
		ActiveCondition,
	}

	GetConditionBeforeUpdating = []schema.StatusCondition{
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
