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
		schema.ResourceStateActive,
		schema.ResourceStateCreating,
		schema.ResourceStatePending,
	)

	GetConditionAfterUpdating = buildConditionSequence(
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateCreating,
		schema.ResourceStatePending,
	)

	GetConditionAfterDeleting = buildConditionSequence(
		schema.ResourceStateDeleting,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateCreating,
		schema.ResourceStatePending,
	)

	GetConditionAfterStopping = buildConditionSequence(
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateCreating,
		schema.ResourceStatePending,
	)

	GetConditionAfterStarting = buildConditionSequence(
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateCreating,
		schema.ResourceStatePending,
	)
	GetConditionAfterStartingWithoutUpdate = buildConditionSequence(
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateCreating,
		schema.ResourceStatePending,
	)

	GetConditionAfterRestarting = buildConditionSequence(
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateUpdating,
		schema.ResourceStateActive,
		schema.ResourceStateCreating,
		schema.ResourceStatePending,
	)
)
