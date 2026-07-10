package mockworkspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureWorkspaceErrorV1 sets up mock stubs for the workspace error scenarios suite.
// All invalid workspace requests return 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create workspace with non-existent region
func ConfigureWorkspaceErrorV1(scenario *mockscenarios.Scenario, p params.WorkspaceErrorV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	// Non-existent region workspace — expect 412
	nonExistentRegionURL := generators.GenerateWorkspaceURL(
		sdkconsts.WorkspaceProviderV1Name,
		p.NonExistentRegionWorkspace.Metadata.Tenant,
		p.NonExistentRegionWorkspace.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(nonExistentRegionURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
