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
