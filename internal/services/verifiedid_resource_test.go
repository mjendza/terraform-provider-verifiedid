package services_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance/check"
	"github.com/mjendza/terraform-provider-verifiedid/internal/clients"
	"github.com/mjendza/terraform-provider-verifiedid/internal/services"
	"github.com/mjendza/terraform-provider-verifiedid/internal/utils"
)

func defaultIgnores() []string {
	return []string{"body", "output", "retry"}
}

type VerifiedIDTestResource struct{}

func TestAcc_ResourceBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource", "test")

	r := VerifiedIDTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("id").IsUUID(),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(regexp.MustCompile(`^applications/[a-f0-9\-]+$`)),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAcc_ResourceUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource", "test")

	r := VerifiedIDTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("id").IsUUID(),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(regexp.MustCompile(`^applications/[a-f0-9\-]+$`)),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
		{
			Config: r.basicUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("id").IsUUID(),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(regexp.MustCompile(`^applications/[a-f0-9\-]+$`)),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAcc_ResourceIgnoreMissingProperty(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource", "test")

	r := VerifiedIDTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicWithIgnoreMissingProperty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("id").IsUUID(),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(regexp.MustCompile(`^applications/[a-f0-9\-]+$`)),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAcc_ResourceRetry(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource", "test")

	r := VerifiedIDTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withRetry(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("id").IsUUID(),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(regexp.MustCompile(`^applications/[a-f0-9\-]+$`)),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAcc_ResourceTimeouts_Create(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource", "test")

	r := VerifiedIDTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withCreateTimeout(),
			// Creating with 1ns should fail quickly with a deadline exceeded error
			ExpectError: regexp.MustCompile(`context deadline exceeded`),
		},
	})
}

func TestAcc_ResourceTimeouts_Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource", "test")

	r := VerifiedIDTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("id").IsUUID(),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(regexp.MustCompile(`^applications/[a-f0-9\-]+$`)),
			),
		},
		{
			Config:      r.withUpdateTimeout(),
			ExpectError: regexp.MustCompile(`context deadline exceeded`),
		},
	})
}

func (r VerifiedIDTestResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	apiVersion := state.Attributes["api_version"]
	url := state.Attributes["url"]

	checkUrl := fmt.Sprintf("%s/%s", url, state.ID)

	body, err := client.VerifiedIDClient.Read(ctx, checkUrl, apiVersion, clients.DefaultRequestOptions())
	if err == nil {
		// Verified ID contracts cannot be HTTP DELETEd; on destroy the provider
		// soft-deletes them by patching status=Disabled. Treat a disabled
		// contract as "no longer present" for CheckDestroy purposes.
		if services.IsContractURL(url) && contractStatusIsDisabled(body) {
			b := false
			return &b, nil
		}
		b := true
		return &b, nil
	}
	if utils.ResponseErrorWasNotFound(err) {
		b := false
		return &b, nil
	}
	return nil, fmt.Errorf("checking for presence of existing %s(api_version=%s) resource: %w", state.ID, apiVersion, err)
}

// contractStatusIsDisabled returns true when the JSON response body decoded
// from the Verified ID API contains a top-level "status" property equal to
// "Disabled" (case-insensitive).
func contractStatusIsDisabled(body interface{}) bool {
	m, ok := body.(map[string]interface{})
	if !ok {
		return false
	}
	status, ok := m["status"].(string)
	if !ok {
		return false
	}
	return strings.EqualFold(status, "Disabled")
}

func (r VerifiedIDTestResource) ImportIdFunc(tfState *terraform.State) (string, error) {
	state := tfState.RootModule().Resources["verifiedid_resource.test"].Primary
	url := state.Attributes["url"]
	return strings.TrimRight(url, "/") + "/" + state.ID, nil
}

func (r VerifiedIDTestResource) basic(data acceptance.TestData) string {
	return `
resource "verifiedid_resource" "test" {
  url = "applications"
  body = {
    displayName = "Demo App"
  }
}
`
}

func (r VerifiedIDTestResource) basicUpdate(data acceptance.TestData) string {
	return `
resource "verifiedid_resource" "test" {
  url = "applications"
  body = {
    displayName = "Demo App Updated"
  }
}
`
}

func (r VerifiedIDTestResource) basicWithIgnoreMissingProperty(data acceptance.TestData) string {
	return `
resource "verifiedid_resource" "test" {
  url = "applications"
  body = {
    displayName = "Demo App With Ignore Missing Property"
    passwordCredentials = [
      {
        displayName = "demo-credential"
      }
    ]
  }
  ignore_missing_property = true
}
`
}

func (r VerifiedIDTestResource) withRetry() string {
	return `
resource "verifiedid_resource" "test" {
  url = "applications"
  body = {
    displayName = "Demo App Retry"
  }
  retry = {
    error_message_regex = [
      "temporary error",
      ".*throttl.*",
    ]
  }
}`
}

func (r VerifiedIDTestResource) withCreateTimeout() string {
	return `
resource "verifiedid_resource" "test" {
  url = "applications"
  timeouts {
    create = "1ns"
  }
  body = {
    displayName = "Demo App Timeout Create"
  }
}
`
}

func (r VerifiedIDTestResource) withUpdateTimeout() string {
	return `
resource "verifiedid_resource" "test" {
  url = "applications"
  timeouts {
    update = "1ns"
  }
  body = {
    displayName = "Demo App Updated Timeout Update"
  }
}
`
}

