package secalib

import (
	"fmt"
	"math"
	"math/rand"
)

// Names
func GenerateSkuName() string {
	return fmt.Sprintf("sku-%d", rand.Intn(math.MaxInt32))
}

func GenerateRoleName() string {
	return fmt.Sprintf("role-%d", rand.Intn(math.MaxInt32))
}

func GenerateRoleAssignmentName() string {
	return fmt.Sprintf("role-assignment-%d", rand.Intn(math.MaxInt32))
}

func GenerateWorkspaceName() string {
	return fmt.Sprintf("workspace-%d", rand.Intn(math.MaxInt32))
}

func GenerateBlockStorageName() string {
	return fmt.Sprintf("disk-%d", rand.Intn(math.MaxInt32))
}

func GenerateImageName() string {
	return fmt.Sprintf("image-%d", rand.Intn(math.MaxInt32))
}

func GenerateInstanceName() string {
	return fmt.Sprintf("instance-%d", rand.Intn(math.MaxInt32))
}

// Resources
func GenerateSkuResource(tenant string, sku string) string {
	return fmt.Sprintf(SkuResource, tenant, sku)
}

func GenerateRoleResource(tenant string, role string) string {
	return fmt.Sprintf(RoleResource, tenant, role)
}

func GenerateRoleAssignmentResource(tenant string, roleAssignment string) string {
	return fmt.Sprintf(RoleAssignmentResource, tenant, roleAssignment)
}

func GenerateWorkspaceResource(tenant string, workspace string) string {
	return fmt.Sprintf(WorkspaceResource, tenant, workspace)
}

func GenerateBlockStorageResource(tenant string, workspace string, blockStorage string) string {
	return fmt.Sprintf(BlockStorageResource, tenant, workspace, blockStorage)
}

func GenerateImageResource(tenant string, image string) string {
	return fmt.Sprintf(ImageResource, tenant, image)
}

func GenerateInstanceResource(tenant string, workspace string, instance string) string {
	return fmt.Sprintf(InstanceResource, tenant, workspace, instance)
}

// References
func GenerateSkuRef(name string) string {
	return fmt.Sprintf(SkuRef, name)
}

func GenerateBlockStorageRef(blockStorageName string) string {
	return fmt.Sprintf(BlockStorageRef, blockStorageName)
}

// URLs
func GenerateRoleURL(tenant string, role string) string {
	return fmt.Sprintf(RoleURLV1, tenant, role)
}

func GenerateRoleAssignmentURL(tenant string, roleAssignment string) string {
	return fmt.Sprintf(RoleAssignmentURLV1, tenant, roleAssignment)
}

func GenerateWorkspaceURL(tenant string, workspace string) string {
	return fmt.Sprintf(WorkspaceURLV1, tenant, workspace)
}

func GenerateStorageSkuURL(tenant string, sku string) string {
	return fmt.Sprintf(StorageSkuURLV1, tenant, sku)
}

func GenerateBlockStorageURL(tenant string, workspace string, blockStorage string) string {
	return fmt.Sprintf(BlockStorageURLV1, tenant, workspace, blockStorage)
}

func GenerateImageURL(tenant string, image string) string {
	return fmt.Sprintf(ImageURLV1, tenant, image)
}

func GenerateInstanceSkuURL(tenant string, sku string) string {
	return fmt.Sprintf(InstanceSkuURLV1, tenant, sku)
}

func GenerateInstanceURL(tenant string, workspace string, instance string) string {
	return fmt.Sprintf(InstanceURLV1, tenant, workspace, instance)
}

// Random

func GenerateStorageSkuIops() int {
	return rand.Intn(maxStorageSkuIops)
}
func GenerateStorageSkuMinVolumeSize(maxSize int) int {
	return rand.Intn(maxSize)
}

func GenerateBlockStorageSize() int {
	return rand.Intn(maxBlockStorageSize)
}
