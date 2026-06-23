package suites

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Region

func (suite *TestSuite) VerifyRegionItemsStep(ctx provider.StepCtx, items []*schema.Region) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "Region items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "Region item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsGlobalResourceMetadataFilled(stepCtx, item.Metadata, "Region", i)

		}
	})
}

// Compute

func (suite *TestSuite) VerifyInstanceSkuItemsStep(ctx provider.StepCtx, items []*schema.InstanceSku) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "InstanceSku items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "InstanceSku item[%d] metadata should not be nil", i)

			// Metadata
			// TODO

		}
	})
}

func (suite *TestSuite) VerifyInstanceItemsStep(ctx provider.StepCtx, items []*schema.Instance) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "Instance items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "Instance item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalWorkspaceResourceMetadataFilled(stepCtx, item.Metadata, "Instance", i)

		}
	})
}

// Network

func (suite *TestSuite) VerifyNetworkSkuItemsStep(ctx provider.StepCtx, items []*schema.NetworkSku) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "NetworkSku items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "NetworkSku item[%d] metadata should not be nil", i)

			// Metadata
			// TODO
		}
	})
}

func (suite *TestSuite) VerifyNetworkItemsStep(ctx provider.StepCtx, items []*schema.Network) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "Network items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "Network item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalWorkspaceResourceMetadataFilled(stepCtx, item.Metadata, "Network", i)

		}
	})
}

func (suite *TestSuite) VerifyInternetGatewayItemsStep(ctx provider.StepCtx, items []*schema.InternetGateway) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "InternetGateway items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "InternetGateway item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalWorkspaceResourceMetadataFilled(stepCtx, item.Metadata, "InternetGateway", i)

		}
	})
}

func (suite *TestSuite) VerifyRouteTableItemsStep(ctx provider.StepCtx, items []*schema.RouteTable) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "RouteTable items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "RouteTable item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalNetworkResourceMetadataFilled(stepCtx, item.Metadata, "RouteTable", i)

		}
	})
}

func (suite *TestSuite) VerifySubnetItemsStep(ctx provider.StepCtx, items []*schema.Subnet) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "Subnet items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "Subnet item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalNetworkResourceMetadataFilled(stepCtx, item.Metadata, "Subnet", i)

		}
	})
}

func (suite *TestSuite) VerifyPublicIpItemsStep(ctx provider.StepCtx, items []*schema.PublicIp) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "PublicIp items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "PublicIp item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalWorkspaceResourceMetadataFilled(stepCtx, item.Metadata, "PublicIp", i)

		}
	})
}

func (suite *TestSuite) VerifyNicItemsStep(ctx provider.StepCtx, items []*schema.Nic) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "NIC items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "NIC item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalWorkspaceResourceMetadataFilled(stepCtx, item.Metadata, "NIC", i)

		}
	})
}

func (suite *TestSuite) VerifySecurityGroupRuleItemsStep(ctx provider.StepCtx, items []*schema.SecurityGroupRule) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "SecurityGroupRule items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "SecurityGroupRule item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalWorkspaceResourceMetadataFilled(stepCtx, item.Metadata, "SecurityGroupRule", i)

		}
	})
}

func (suite *TestSuite) VerifySecurityGroupItemsStep(ctx provider.StepCtx, items []*schema.SecurityGroup) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "SecurityGroup items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "SecurityGroup item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalWorkspaceResourceMetadataFilled(stepCtx, item.Metadata, "SecurityGroup", i)

		}
	})
}

// Storage

func (suite *TestSuite) VerifyStorageSkuItemsStep(ctx provider.StepCtx, items []*schema.StorageSku) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "StorageSku items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "StorageSku item[%d] metadata should not be nil", i)

			// Metadata
			// TODO
		}
	})
}

func (suite *TestSuite) VerifyBlockStorageItemsStep(ctx provider.StepCtx, items []*schema.BlockStorage) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "BlockStorage items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "BlockStorage item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalWorkspaceResourceMetadataFilled(stepCtx, item.Metadata, "BlockStorage", i)

		}
	})
}

func (suite *TestSuite) VerifyImageItemsStep(ctx provider.StepCtx, items []*schema.Image) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotEmpty(item, "Image items should not be empty")
			stepCtx.Require().NotNil(item.Metadata, "Image item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalResourceMetadataFilled(stepCtx, item.Metadata, "Image", i)
		}
	})
}

// Workspace

func (suite *TestSuite) VerifyWorkspaceItemsStep(ctx provider.StepCtx, items []*schema.Workspace) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotNil(item, "Workspace item[%d] should not be nil", i)
			stepCtx.Require().NotNil(item.Metadata, "Workspace item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsRegionalResourceMetadataFilled(stepCtx, item.Metadata, "Workspace", i)
		}
	})
}

// Authorization

func (suite *TestSuite) VerifyRoleItemsStep(ctx provider.StepCtx, items []*schema.Role) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotNil(item, "Role item[%d] should not be nil", i)
			stepCtx.Require().NotNil(item.Metadata, "Role item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsGlobalTenantResourceMetadataFilled(stepCtx, item.Metadata, "Role", i)
		}
	})
}

func (suite *TestSuite) VerifyRoleAssignmentItemsStep(ctx provider.StepCtx, items []*schema.RoleAssignment) {
	ctx.WithNewStep("Verify items", func(stepCtx provider.StepCtx) {
		for i, item := range items {
			stepCtx.Require().NotNil(item, "RoleAssignment item[%d] should not be nil", i)
			stepCtx.Require().NotNil(item.Metadata, "RoleAssignment item[%d] metadata should not be nil", i)

			// Metadata
			suite.VerifyItemsGlobalTenantResourceMetadataFilled(stepCtx, item.Metadata, "RoleAssignment", i)
		}
	})
}
