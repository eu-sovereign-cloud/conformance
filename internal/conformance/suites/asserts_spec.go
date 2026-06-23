package suites

import (
	"fmt"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Authorization

func (suite *TestSuite) VerifyRoleSpecStep(ctx provider.StepCtx, expected *schema.RoleSpec, actual *schema.RoleSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(len(expected.Permissions), len(actual.Permissions),
			"Permissions list length should match expected")

		for i := 0; i < len(expected.Permissions); i++ {
			expectedPerm := expected.Permissions[i]
			actualPerm := actual.Permissions[i]
			stepCtx.Require().Equal(expectedPerm.Provider, actualPerm.Provider,
				fmt.Sprintf("Permission [%d] provider should match expected", i))
			stepCtx.Require().Equal(expectedPerm.Resources, actualPerm.Resources,
				fmt.Sprintf("Permission [%d] resources should match expected", i))
			stepCtx.Require().Equal(expectedPerm.Verb, actualPerm.Verb,
				fmt.Sprintf("Permission [%d] verb should match expected", i))
		}
	})
}

func (suite *TestSuite) VerifyRoleAssignmentSpecStep(ctx provider.StepCtx, expected *schema.RoleAssignmentSpec, actual *schema.RoleAssignmentSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Roles, actual.Roles, "Roles provider should match expected")
		stepCtx.Require().Equal(expected.Subs, actual.Subs, "Subs should match expected")
		stepCtx.Require().Equal(len(expected.Scopes), len(actual.Scopes), "Scope list length should match expected")

		for i := 0; i < len(expected.Scopes); i++ {
			expectedScope := expected.Scopes[i]
			actualScope := actual.Scopes[i]

			if len(actualScope.Tenants) > 0 {
				stepCtx.Require().Equal(expectedScope.Tenants, actualScope.Tenants, fmt.Sprintf("Scope [%d] tenants should match expected", i))
			}
			if len(actualScope.Regions) > 0 {
				stepCtx.Require().Equal(expectedScope.Regions, actualScope.Regions, fmt.Sprintf("Scope [%d] regions should match expected", i))
			}
			if len(actualScope.Workspaces) > 0 {
				stepCtx.Require().Equal(expectedScope.Workspaces, actualScope.Workspaces, fmt.Sprintf("Scope [%d] workspaces should match expected", i))
			}
		}
	})
}

// Storage

func (suite *TestSuite) VerifyBlockStorageSpecStep(ctx provider.StepCtx, expected *schema.BlockStorageSpec, actual *schema.BlockStorageSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.SizeGB, actual.SizeGB, "SizeGB should match expected")
		stepCtx.Require().Equal(expected.SkuRef, actual.SkuRef, "SkuRef should match expected")
	})
}

func (suite *TestSuite) VerifyImageSpecStep(ctx provider.StepCtx, expected *schema.ImageSpec, actual *schema.ImageSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.BlockStorageRef, actual.BlockStorageRef, "BlockStorageRef should match expected")
		stepCtx.Require().Equal(expected.CpuArchitecture, actual.CpuArchitecture, "CpuArchitecture should match expected")
	})
}

// Compute

func (suite *TestSuite) VerifyInstanceSpecStep(ctx provider.StepCtx, expected *schema.InstanceSpec, actual *schema.InstanceSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.SkuRef, actual.SkuRef, "SkuRef should match expected")
		stepCtx.Require().Equal(expected.Zone, actual.Zone, "Zone should match expected")
		stepCtx.Require().Equal(expected.BootVolume.DeviceRef, actual.BootVolume.DeviceRef, "BootVolume.DeviceRef should match expected")
	})
}

// Network

func (suite *TestSuite) VerifyNetworkSpecStep(ctx provider.StepCtx, expected *schema.NetworkSpec, actual *schema.NetworkSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		if actual.Cidr.Ipv4 != "" {
			stepCtx.Require().Equal(expected.Cidr.Ipv4, actual.Cidr.Ipv4, "Cidr.Ipv4 should match expected")
		}
		if actual.Cidr.Ipv6 != "" {
			stepCtx.Require().Equal(expected.Cidr.Ipv6, actual.Cidr.Ipv6, "Cidr.Ipv6 should match expected")
		}

		stepCtx.Require().Equal(expected.SkuRef, actual.SkuRef, "SkuRef should match expected")
		stepCtx.Require().Equal(expected.RouteTableRef, actual.RouteTableRef, "RouteTableRef should match expected")
	})
}

func (suite *TestSuite) VerifyInternetGatewaySpecStep(ctx provider.StepCtx, expected *schema.InternetGatewaySpec, actual *schema.InternetGatewaySpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.EgressOnly, actual.EgressOnly, "EgressOnly should match expected")
	})
}

func (suite *TestSuite) VerifyRouteTableSpecStep(ctx provider.StepCtx, expected *schema.RouteTableSpec, actual *schema.RouteTableSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(len(expected.Routes), len(actual.Routes), "Route list length should match expected")
		for i := 0; i < len(expected.Routes); i++ {
			expectedRoute := expected.Routes[i]
			actualRoute := actual.Routes[i]
			stepCtx.Require().Equal(expectedRoute.DestinationCidrBlock, actualRoute.DestinationCidrBlock, fmt.Sprintf("Route [%d] DestinationCidrBlock should match expected", i))
			stepCtx.Require().Equal(expectedRoute.TargetRef, actualRoute.TargetRef, fmt.Sprintf("Route [%d] TargetRef should match expected", i))
		}
	})
}

func (suite *TestSuite) VerifySubnetSpecStep(ctx provider.StepCtx, expected *schema.SubnetSpec, actual *schema.SubnetSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		if actual.Cidr.Ipv4 != "" {
			stepCtx.Require().Equal(expected.Cidr.Ipv4, actual.Cidr.Ipv4, "Cidr.Ipv4 should match expected")
		}
		if actual.Cidr.Ipv6 != "" {
			stepCtx.Require().Equal(expected.Cidr.Ipv6, actual.Cidr.Ipv6, "Cidr.Ipv6 should match expected")
		}
		stepCtx.Require().Equal(expected.Zone, actual.Zone, "Zone should match expected")
	})
}

func (suite *TestSuite) VerifyPublicIpSpecStep(ctx provider.StepCtx, expected *schema.PublicIpSpec, actual *schema.PublicIpSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Version, actual.Version, "Version should match expected")
		if actual.Address != "" {
			stepCtx.Require().Equal(expected.Address, actual.Address, "Address should match expected")
		}
	})
}

func (suite *TestSuite) VerifyNicSpecStep(ctx provider.StepCtx, expected *schema.NicSpec, actual *schema.NicSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Addresses, actual.Addresses, "Addresses should match expected")
		if actual.PublicIpRefs != nil {
			stepCtx.Require().Equal(expected.PublicIpRefs, actual.PublicIpRefs, "PublicIpRefs should match expected")
		}
		stepCtx.Require().Equal(expected.SubnetRef, actual.SubnetRef, "SubnetRef should match expected")
	})
}

func (suite *TestSuite) VerifySecurityGroupRuleSpecStep(ctx provider.StepCtx, expected *schema.SecurityGroupRuleSpec, actual *schema.SecurityGroupRuleSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Direction, actual.Direction, "Direction should match expected")
	})
}

func (suite *TestSuite) VerifySecurityGroupSpecStep(ctx provider.StepCtx, expected *schema.SecurityGroupSpec, actual *schema.SecurityGroupSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		if actual.Rules != nil && expected.Rules != nil {
			stepCtx.Require().Equal(len(expected.Rules), len(actual.Rules), "Rule list length should match expected")
			for i := 0; i < len(expected.Rules); i++ {
				expectedRule := expected.Rules[i]
				actualRule := actual.Rules[i]
				stepCtx.Require().Equal(expectedRule.Direction, actualRule.Direction, fmt.Sprintf("Rule [%d] Direction should match expected", i))
			}
		} else {
			stepCtx.Require().Equal(expected.RuleRefs, actual.RuleRefs, "RuleRefs should match expected")
		}
	})
}

// Region

func (suite *TestSuite) VerifyRegionSpecStep(ctx provider.StepCtx, _ *schema.RegionSpec, actual *schema.RegionSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().GreaterOrEqual(len(actual.AvailableZones), 1, "AvailableZones list length should match expected")

		stepCtx.Require().GreaterOrEqual(len(actual.Providers), 1, "Providers list length should greater then 1")
	})
}
