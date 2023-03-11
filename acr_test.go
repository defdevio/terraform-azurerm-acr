package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerregistry/mgmt/containerregistry"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/defdevio/terratest-helpers/pkg/helpers"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var (
	ctx            = context.Background()
	subscriptionID = os.Getenv("ARM_SUBSCRIPTION_ID")
	workDir, _     = os.Getwd()
)

func terraformVars() map[string]any {
	testVars := map[string]any{
		"environment":         "test",
		"location":            "westus",
		"name":                "devdevio",
		"resource_count":      1,
		"resource_group_name": "test",
	}

	return testVars
}

func TestCreateAzureContainerRegistry(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "",
		Vars:         terraformVars(),
	})

	testFiles := []string{
		"provider.tf",
		"terraform.tfstate",
		"terraform.tfstate.backup",
		".terraform.lock.hcl",
		".terraform",
	}

	// Defer cleaning up the test files created during the test
	defer helpers.CleanUpTestFiles(t, testFiles, workDir)

	// Create the provider file
	err := helpers.CreateAzureProviderFile(path.Join(workDir, "provider.tf"), t)
	if err != nil {
		t.Fatal(err)
	}

	// Use type assertions to ensure the interface values are the expected type for the given
	// terraform variable value
	environment, ok := terraformOptions.Vars["environment"].(string)
	if !ok {
		t.Fatal("A value type of 'string' was expected for 'environment'")
	}

	location, ok := terraformOptions.Vars["location"].(string)
	if !ok {
		t.Fatal("A value type of 'string' was expected for 'location'")
	}

	name, ok := terraformOptions.Vars["name"].(string)
	if !ok {
		t.Fatal("A value type of 'string' was expected for 'name'")
	}

	resourceGroup, ok := terraformOptions.Vars["resource_group_name"].(string)
	if !ok {
		t.Fatal("A value type of 'string' was expected for 'resourceGroup'")
	}

	// Create a resource group client
	resourceGroupClient, err := azure.CreateResourceGroupClientE(subscriptionID)
	if err != nil {
		t.Fatal(err)
	}

	// Defer the deletion of the resource group until all test functions have finished
	defer resourceGroupClient.Delete(ctx, resourceGroup)

	// Create the resource group using the resourceGroupClient
	resp, err := resourceGroupClient.CreateOrUpdate(ctx, resourceGroup, resources.Group{
		Location: &location,
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == 201 {
		log.Printf("Created resource group '%s'", *resp.Name)
	}

	// Defer destroying the terraform resources until the rest of the test functions finish
	defer terraform.Destroy(t, terraformOptions)

	// Init and apply the terraform module
	terraform.InitAndApply(t, terraformOptions)

	// Get the managed cluster resource the test created
	containerClient, err := azure.GetContainerRegistryClientE(subscriptionID)
	if err != nil {
		t.Fatal(err)
	}

	// Format the name in the same manner the module will create the name
	formattedName := fmt.Sprintf("%s%s%s", environment, location, name)

	registry, err := containerClient.Get(ctx, resourceGroup, formattedName)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the registry returns a succeeded provisioning state
	assert.Equal(t, containerregistry.Succeeded, registry.ProvisioningState)

	// Assert that the deployed registry resource has the same name as the desired resource
	assert.Equal(t, formattedName, *registry.Name)

	// Assert that the admin user account is false
	assert.Equal(t, false, *registry.AdminUserEnabled)
}
