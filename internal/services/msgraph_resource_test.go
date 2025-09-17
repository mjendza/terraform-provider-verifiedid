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

func TestAcc_ResourceGroupMember(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource", "test")

	r := VerifiedIDTestResource{}

	importStep := data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...)
	importStep.ImportStateVerify = false
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.groupMember(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("id").IsUUID(),
				check.That(data.ResourceName).Key("id").MatchesOtherKey(check.That("verifiedid_resource.servicePrincipal_application").Key("id")),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(regexp.MustCompile(`^groups/[a-f0-9\-]+/members/[a-f0-9\-]+$`)),
			),
		},
		importStep,
	})
}

func TestAcc_ResourceIgnoreMissingProperty(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource", "test")

	r := VerifiedIDTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.groupOwnerBind("My Group Owners Bind"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("id").IsUUID(),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(regexp.MustCompile(`^groups/[a-f0-9\-]+$`)),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
	})
}

func TestAcc_ResourceGroupOwnerBind_UpdateDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource", "test")

	r := VerifiedIDTestResource{}

	importStep := data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...)
	importStep.ImportStateVerify = false

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.groupOwnerBind("My Group Owners Bind"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("id").IsUUID(),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(regexp.MustCompile(`^groups/[a-f0-9\-]+$`)),
			),
		},
		data.ImportStepWithImportStateIdFunc(r.ImportIdFunc, defaultIgnores()...),
		{
			Config: r.groupOwnerBind("My Group Owners Bind Updated"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("id").IsUUID(),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(regexp.MustCompile(`^groups/[a-f0-9\-]+$`)),
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

	var checkUrl string
	if !strings.Contains(url, "/$ref") {
		checkUrl = fmt.Sprintf("%s/%s", url, state.ID)
	} else {
		checkUrl = url
	}

	_, err := client.VerifiedIDClient.Read(ctx, checkUrl, apiVersion, clients.DefaultRequestOptions())
	if err == nil {
		b := true
		return &b, nil
	}
	if utils.ResponseErrorWasNotFound(err) {
		b := false
		return &b, nil
	}
	return nil, fmt.Errorf("checking for presence of existing %s(api_version=%s) resource: %w", state.ID, apiVersion, err)
}

func (r VerifiedIDTestResource) ImportIdFunc(tfState *terraform.State) (string, error) {
	state := tfState.RootModule().Resources["verifiedid_resource.test"].Primary
	url := state.Attributes["url"]
	if !strings.Contains(url, "/$ref") {
		return fmt.Sprintf("%s/%s", url, state.ID), nil
	}
	return strings.ReplaceAll(url, "/$ref", fmt.Sprintf("/%s/$ref", state.ID)), nil
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

func (r VerifiedIDTestResource) groupMember() string {
	return `
resource "verifiedid_resource" "application" {
  url = "applications"
  body = {
    displayName = "My Application"
  }
  response_export_values = {
    appId = "appId"
  }
}

resource "verifiedid_resource" "servicePrincipal_application" {
  url = "servicePrincipals"
  body = {
    appId = verifiedid_resource.application.output.appId
  }
}

resource "verifiedid_resource" "group" {
  url = "groups"
  body = {
    displayName     = "My Group"
    mailEnabled     = false
    mailNickname    = "mygroup"
    securityEnabled = true
  }
}

resource "verifiedid_resource" "test" {
  url = "groups/${verifiedid_resource.group.id}/members/$ref"
  body = {
    "@odata.id" = "https://graph.microsoft.com/v1.0/directoryObjects/${verifiedid_resource.servicePrincipal_application.id}"
  }
}
`
}

func (r VerifiedIDTestResource) groupOwnerBind(displayName string) string {
	return fmt.Sprintf(`
resource "verifiedid_resource" "application" {
  url = "applications"
  body = {
    displayName = "My Application"
  }
  response_export_values = {
    appId = "appId"
  }
}

resource "verifiedid_resource" "servicePrincipal_application" {
  url = "servicePrincipals"
  body = {
    appId = verifiedid_resource.application.output.appId
  }
}

resource "verifiedid_resource" "test" {
  url = "groups"
  body = {
    displayName     = "%s"
    mailEnabled     = false
    mailNickname    = "mygroup-owners-bind"
    securityEnabled = true
    "owners@odata.bind" = [
      "https://graph.microsoft.com/v1.0/directoryObjects/${verifiedid_resource.servicePrincipal_application.id}"
    ]
  }
}
`, displayName)
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
