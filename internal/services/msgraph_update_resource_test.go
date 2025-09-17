package services_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance/check"
	"github.com/mjendza/terraform-provider-verifiedid/internal/clients"
	"github.com/mjendza/terraform-provider-verifiedid/internal/utils"
)

type VerifiedIDTestUpdateResource struct{}

func TestAcc_UpdateResourceBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_update_resource", "test")

	r := VerifiedIDTestUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic("Demo App Updated"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("id").IsUUID(),
			),
		},
	})
}

func TestAcc_UpdateResourceUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_update_resource", "test")

	r := VerifiedIDTestUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic("Demo App Updated"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
			),
		},
		{
			Config: r.basic("Demo App Updated Again"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
			),
		},
	})
}

func TestAcc_UpdateResourceTimeouts_Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_update_resource", "test")
	r := VerifiedIDTestUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withUpdateTimeout("Demo App"),
		},
		{
			Config:      r.withUpdateTimeout("Demo App Updated"),
			ExpectError: regexp.MustCompile(`context deadline exceeded`),
		},
	})
}

func TestAcc_UpdateResourceTimeouts_Create(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_update_resource", "test")
	r := VerifiedIDTestUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.withCreateTimeout("Demo App Updated"),
			ExpectError: regexp.MustCompile(`context deadline exceeded`),
		},
	})
}

func TestAcc_UpdateResourceTimeouts_Read(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_update_resource", "test")
	r := VerifiedIDTestUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.withReadTimeout("Demo App Updated"),
			ExpectError: regexp.MustCompile(`context deadline exceeded`),
		},
	})
}

func TestAcc_UpdateResourceRetry(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_update_resource", "test")

	r := VerifiedIDTestUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withRetry(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
			),
		},
	})
}

func TestAcc_UpdateResource_GroupOwnerBind_UpdateDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_update_resource", "test")
	r := VerifiedIDTestUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.groupWithOwnerUpdate("My Group Owners Bind 2"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
			),
		},
		{
			Config: r.groupWithOwnerUpdate("My Group Owners Bind 3"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
			),
		},
	})
}

func (r VerifiedIDTestUpdateResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	apiVersion := state.Attributes["api_version"]
	url := state.Attributes["url"]

	_, err := client.VerifiedIDClient.Read(ctx, url, apiVersion, clients.DefaultRequestOptions())
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

func (r VerifiedIDTestUpdateResource) basic(displayName string) string {
	return fmt.Sprintf(`
resource "verifiedid_resource" "application" {
  url = "applications"
  body = {
    displayName = "Demo App"
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}

resource "verifiedid_update_resource" "test" {
  url = "applications/${verifiedid_resource.application.id}"
  body = {
    displayName = "%s"
  }
}
`, displayName)
}

func (r VerifiedIDTestUpdateResource) withUpdateTimeout(displayName string) string {
	return fmt.Sprintf(`
%s

resource "verifiedid_update_resource" "test" {
  url = "applications/${verifiedid_resource.application.id}"
  body = {
    displayName = "%s"
  }
  timeouts {
    update = "1ns"
  }
}
`, VerifiedIDTestUpdateResource{}.applicationOnly(), displayName)
}

func (r VerifiedIDTestUpdateResource) withReadTimeout(displayName string) string {
	return fmt.Sprintf(`
%s

resource "verifiedid_update_resource" "test" {
  url = "applications/${verifiedid_resource.application.id}"
  body = {
    displayName = "%s"
  }
  timeouts {
    read = "1ns"
  }
}
`, VerifiedIDTestUpdateResource{}.applicationOnly(), displayName)
}

func (r VerifiedIDTestUpdateResource) withCreateTimeout(displayName string) string {
	return fmt.Sprintf(`
%s

resource "verifiedid_update_resource" "test" {
  url = "applications/${verifiedid_resource.application.id}"
  body = {
    displayName = "%s"
  }
  timeouts {
    create = "1ns"
  }
}
`, VerifiedIDTestUpdateResource{}.applicationOnly(), displayName)
}

// applicationOnly returns just the application resource to be used for composing
// different update resource configurations.
func (r VerifiedIDTestUpdateResource) applicationOnly() string {
	return `
resource "verifiedid_resource" "application" {
  url = "applications"
  body = {
    displayName = "Demo App"
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}
`
}

func (r VerifiedIDTestUpdateResource) withRetry() string {
	return `
resource "verifiedid_resource" "application" {
  url = "applications"
  body = {
    displayName = "Demo App"
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}

resource "verifiedid_update_resource" "test" {
  url = "applications/${verifiedid_resource.application.id}"
  body = {
    displayName = "Demo App Updated With Retry"
  }
  retry = {
    error_message_regex = [
      ".*throttl.*",
      "temporary error",
    ]
  }
}
`
}

func (r VerifiedIDTestUpdateResource) groupWithOwnerBase() string {
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
    displayName     = "My Group Owners Bind"
    mailEnabled     = false
    mailNickname    = "mygroup-owners-bind"
    securityEnabled = true
    "owners@odata.bind" = [
      "https://graph.microsoft.com/v1.0/directoryObjects/${verifiedid_resource.servicePrincipal_application.id}"
    ]
  }
  lifecycle {
    ignore_changes = [body.displayName]
  }
}
`
}

func (r VerifiedIDTestUpdateResource) groupWithOwnerUpdate(displayName string) string {
	return fmt.Sprintf(`
%s

resource "verifiedid_update_resource" "test" {
  url = "groups/${verifiedid_resource.group.id}"
  body = {
    displayName = "%s"
  }
}
`, r.groupWithOwnerBase(), displayName)
}
