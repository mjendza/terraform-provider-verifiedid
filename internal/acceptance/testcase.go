package acceptance

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mjendza/terraform-provider-verifiedid/internal/clients"
	"github.com/mjendza/terraform-provider-verifiedid/internal/provider"
)

const (
	// charSetAlphaNum is the alphanumeric character set for use with randStringFromCharSet
	charSetAlphaNum = "abcdefghijklmnopqrstuvwxyz012346789"
)

type TestData struct {
	// RandomInteger is a random integer which is unique to this test case
	RandomInteger int

	// RandomString is a random 5 character string is unique to this test case
	RandomString string

	// ResourceName is the fully qualified resource name, comprising of the
	// resource type and then the resource label
	// e.g. `azurerm_resource_group.test`
	ResourceName string

	// ResourceType is the Terraform Resource Type - `azurerm_resource_group`
	ResourceType string

	// resourceLabel is the local used for the resource - generally "test""
	resourceLabel string
}

// BuildTestData generates some test data for the given resource
func BuildTestData(t *testing.T, resourceType string, resourceLabel string) TestData {
	return TestData{
		RandomInteger: RandTimeInt(),
		RandomString:  RandStringFromCharSet(5, charSetAlphaNum),
		ResourceName:  fmt.Sprintf("%s.%s", resourceType, resourceLabel),

		ResourceType:  resourceType,
		resourceLabel: resourceLabel,
	}
}

// RandomIntOfLength is a random 8 to 18 digit integer which is unique to this test case
func (td *TestData) RandomInt() int {
	return RandTimeInt()
}

func (td *TestData) RandomStringOfLength(len int) string {
	return RandStringFromCharSet(len, charSetAlphaNum)
}

// UpgradeTestDeployStep returns a test step used to deploy the configuration with previous version
func (td TestData) UpgradeTestDeployStep(step resource.TestStep, upgradeFrom string) resource.TestStep {
	if step.ExternalProviders == nil {
		step.ExternalProviders = td.externalProviders()
	}
	step.ExternalProviders["msgraph"] = resource.ExternalProvider{
		Source:            "registry.terraform.io/microsoft/msgraph",
		VersionConstraint: fmt.Sprintf("= %s", upgradeFrom),
	}
	step.ProtoV6ProviderFactories = nil
	return step
}

// UpgradeTestApplyStep returns a test step used to run terraform apply with the development version
func (td TestData) UpgradeTestApplyStep(applyStep resource.TestStep) resource.TestStep {
	if applyStep.ExternalProviders == nil {
		applyStep.ExternalProviders = td.externalProviders()
	}
	applyStep.ProtoV6ProviderFactories = td.providers()
	return applyStep
}

// UpgradeTestPlanStep returns a test step used to run terraform plan with the development version to check if there's any changes
func (td TestData) UpgradeTestPlanStep(planStep resource.TestStep) resource.TestStep {
	planStep.PlanOnly = true
	if planStep.ExternalProviders == nil {
		planStep.ExternalProviders = td.externalProviders()
	}
	planStep.ProtoV6ProviderFactories = td.providers()
	return planStep
}

func (td TestData) UpgradeTest(t *testing.T, testResource TestResource, steps []resource.TestStep) {
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := BuildTestClient()
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}
	resource.ParallelTest(t, testCase)
}

// lintignore:AT001
func (td TestData) DataSourceTest(t *testing.T, steps []resource.TestStep) {
	// DataSources don't need a check destroy - however since this is a wrapper function
	// and not matching the ignore pattern `XXX_data_source_test.go`, this needs to be explicitly opted out
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		Steps:    steps,
	}
	td.runAcceptanceTest(t, testCase)
}

func (td TestData) ResourceTest(t *testing.T, testResource TestResource, steps []resource.TestStep) {
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client, err := BuildTestClient()
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
			return CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}
	td.runAcceptanceTest(t, testCase)
}

func (td TestData) runAcceptanceTest(t *testing.T, testCase resource.TestCase) {
	testCase.ExternalProviders = td.externalProviders()
	// If any test steps require their own external providers, then we need to clear the global list
	providersInTestStep := false
	for i, step := range testCase.Steps {
		if step.ExternalProviders != nil {
			testCase.ExternalProviders = nil
			step.ProtoV6ProviderFactories = td.providers()
			testCase.Steps[i] = step
			providersInTestStep = true
		}
	}
	if !providersInTestStep {
		testCase.ProtoV6ProviderFactories = td.providers()
	}

	resource.ParallelTest(t, testCase)
}

func (td TestData) providers() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"verifiedid": providerserver.NewProtocol6WithError(&provider.VerifiedIDProvider{}),
	}
}

func (td TestData) externalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{}
}

func PreCheck(t *testing.T) {
	// `go test ./internal/services` (or any subpackage) sets the test process
	// CWD to that package directory, so a relative ARM_CLIENT_CERTIFICATE_PATH
	// like `./cert/cert.pfx` would resolve under e.g. internal/services/ and
	// fail. Anchor it to the repo root before the provider's Configure() reads
	// the env var. Production behavior (provider.go) is unchanged.
	if v := os.Getenv("ARM_CLIENT_CERTIFICATE_PATH"); v != "" && !filepath.IsAbs(v) {
		if root, err := repoRoot(); err == nil {
			abs := filepath.Join(root, v)
			if _, err := os.Stat(abs); err == nil {
				_ = os.Setenv("ARM_CLIENT_CERTIFICATE_PATH", abs)
			}
		}
	}

	if v := os.Getenv("TF_ACC"); v == "" {
		t.Fatalf(`TF_ACC must be set for acceptance tests!
For tests that authenticate with Azure by using a Service Principal, the following environment variables must be set:
- ARM_CLIENT_ID
- ARM_CLIENT_SECRET
- ARM_TENANT_ID

For tests that authenticate with Azure by OIDC in github action, the following environment variables must be set:
- ARM_CLIENT_ID
- ARM_TENANT_ID

For tests that authenticate with Azure by using a Service Principal with Certificate, the following environment variables must be set:
- ARM_CLIENT_ID
- ARM_CLIENT_CERTIFICATE_PATH
- ARM_TENANT_ID
`)
	}

	if !hasAcceptanceCredentials() {
		t.Fatalf(`acceptance tests need credentials for the Verified ID Admin API. Set ONE of:
  - ARM_CLIENT_ID + ARM_CLIENT_SECRET + ARM_TENANT_ID                     (service principal + secret)
  - ARM_CLIENT_ID + ARM_TENANT_ID + ARM_CLIENT_CERTIFICATE_PATH           (service principal + certificate)
  - ARM_CLIENT_ID + ARM_TENANT_ID + ARM_USE_OIDC=true                     (OIDC, e.g. GitHub Actions)
  - ARM_USE_CLI=true                                                      (Azure CLI, after 'az login')

ARM_USE_CLI defaults to true, so if you intended to use 'az login' make sure
ARM_USE_CLI is not explicitly set to "false". See docs/guides/ for details.`)
	}
}

// hasAcceptanceCredentials reports whether at least one credential mode that
// the provider's ChainedTokenCredential can satisfy is configured. The goal is
// to fail loudly on partial SP setups (e.g. ARM_CLIENT_ID set but
// ARM_CLIENT_SECRET missing) instead of silently falling through to Azure CLI
// and producing a confusing 401 mid-test.
func hasAcceptanceCredentials() bool {
	clientID := os.Getenv("ARM_CLIENT_ID")
	tenantID := os.Getenv("ARM_TENANT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	certPath := os.Getenv("ARM_CLIENT_CERTIFICATE_PATH")
	cert := os.Getenv("ARM_CLIENT_CERTIFICATE")
	useOIDC := os.Getenv("ARM_USE_OIDC") == "true"
	useCLI := os.Getenv("ARM_USE_CLI")

	if useOIDC {
		return clientID != "" && tenantID != ""
	}

	// If the user supplied any SP-style credential bits, validate the full
	// triple — partial SP setup is the most common silent-failure mode.
	spIntent := clientID != "" || clientSecret != "" || certPath != "" || cert != ""
	if spIntent {
		if clientID != "" && tenantID != "" && clientSecret != "" {
			return true
		}
		if clientID != "" && tenantID != "" && (certPath != "" || cert != "") {
			return true
		}
		return false
	}

	// No SP env at all -> rely on Azure CLI (provider default; see
	// testclient.go where ARM_USE_CLI defaults to true).
	return useCLI == "" || useCLI == "true"
}

// CheckDestroyedFunc returns a TestCheckFunc which validates the resource no longer exists
func CheckDestroyedFunc(client *clients.Client, testResource TestResource, resourceType, resourceName string) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		ctx := client.StopContext

		for label, resourceState := range state.RootModule().Resources {
			if resourceState.Type != resourceType {
				continue
			}
			if label != resourceName {
				continue
			}

			// Destroy is unconcerned with an error checking the status, since this is going to be "not found"
			result, err := testResource.Exists(ctx, client, resourceState.Primary)
			if result == nil && err == nil {
				return fmt.Errorf("should have either an error or a result when checking if %q has been destroyed", resourceName)
			}
			if result != nil && *result {
				return fmt.Errorf("%q still exists", resourceName)
			}
		}

		return nil
	}
}
