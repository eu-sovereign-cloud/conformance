package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureRoleErrorV1 sets up mock stubs for the role error scenarios suite.
// All invalid role requests return 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create role with invalid provider in permission
//   - Create role with empty permissions list
//   - Create role with invalid HTTP verb in permission
func ConfigureRoleErrorV1(scenario *mockscenarios.Scenario, p params.RoleErrorV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	// Invalid provider permission role — expect 412
	invalidProviderURL := generators.GenerateRoleURL(
		sdkconsts.AuthorizationProviderV1Name,
		p.InvalidProviderPermissionRole.Metadata.Tenant,
		p.InvalidProviderPermissionRole.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidProviderURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty permissions role — expect 412
	emptyPermissionsURL := generators.GenerateRoleURL(
		sdkconsts.AuthorizationProviderV1Name,
		p.EmptyPermissionsRole.Metadata.Tenant,
		p.EmptyPermissionsRole.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyPermissionsURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid verb role — expect 412
	invalidVerbURL := generators.GenerateRoleURL(
		sdkconsts.AuthorizationProviderV1Name,
		p.InvalidVerbRole.Metadata.Tenant,
		p.InvalidVerbRole.Metadata.Name,
	)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidVerbURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
