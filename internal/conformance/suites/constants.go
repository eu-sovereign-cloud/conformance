package suites

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

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
	CreatedResourceExpectedStates = []schema.ResourceState{schema.ResourceStatePending, schema.ResourceStateCreating, schema.ResourceStateActive}
	UpdatedResourceExpectedStates = []schema.ResourceState{schema.ResourceStateActive, schema.ResourceStateUpdating}
)
