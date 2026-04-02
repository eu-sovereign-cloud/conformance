package constants

import (
	"time"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func buildConditionSequence(states ...schema.ResourceState) []schema.StatusCondition {
	base := time.Now()
	conditions := make([]schema.StatusCondition, len(states))
	for i, state := range states {
		conditions[i] = schema.StatusCondition{
			LastTransitionAt: base.Add(time.Duration(i) * time.Second),
			State:            state,
		}
	}
	return conditions
}

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

	GetConditionAfterCreating = buildConditionSequence(
		schema.ResourceStatePending,
		schema.ResourceStateCreating,
		schema.ResourceStateActive,
	)

	GetConditionAfterUpdating = buildConditionSequence(
		schema.ResourceStatePending,
		schema.ResourceStateCreating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
	)

	GetConditionAfterDeleting = buildConditionSequence(
		schema.ResourceStatePending,
		schema.ResourceStateCreating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateDeleting,
	)

	GetConditionAfterStopping = buildConditionSequence(
		schema.ResourceStatePending,
		schema.ResourceStateCreating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
	)

	GetConditionAfterStarting = buildConditionSequence(
		schema.ResourceStatePending,
		schema.ResourceStateCreating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
	)

	GetConditionAfterRestarting = buildConditionSequence(
		schema.ResourceStatePending,
		schema.ResourceStateCreating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
	)
)

func GetConditionBeforeUpdating() []schema.StatusCondition {
	return buildConditionSequence(
		schema.ResourceStatePending,
		schema.ResourceStateCreating,
		schema.ResourceStateActive,
	)
}
