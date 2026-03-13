package suites

import (
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

const (
	// Asserts
	passed = "passed"

	// Step Parameters
	providerStepParameter  = "provider"
	operationStepParameter = "operation"
	tenantStepParameter    = "tenant"
	workspaceStepParameter = "workspace"
	networkStepParameter   = "network"
)

// Expected States
var (
	CreatedResourceExpectedStates = constants.CreatedResourceExpectedStates
	UpdatedResourceExpectedStates = constants.UpdatedResourceExpectedStates
)

// Conditions
var (
	ActiveCondition   = constants.ActiveCondition
	CreatingCondition = constants.CreatingCondition
	PendingCondition  = constants.PendingCondition
	UpdatingCondition = constants.UpdatingCondition
	DeletingCondition = constants.DeletingCondition

	GetConditionAfterCreating = constants.GetConditionAfterCreating
	GetConditionAfterUpdating = constants.GetConditionAfterUpdating
)

// suppress unused import warning - schema is used transitively via constants
var _ schema.StatusCondition
