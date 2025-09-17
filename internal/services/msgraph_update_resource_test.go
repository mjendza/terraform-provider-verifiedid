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

type MSGraphTestUpdateResource struct{}

func TestAcc_UpdateResourceBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "msgraph_update_resource", "test")

	r := MSGraphTestUpdateResource{}

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
	data := acceptance.BuildTestData(t, "msgraph_update_resource", "test")

	r := MSGraphTestUpdateResource{}

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
	data := acceptance.BuildTestData(t, "msgraph_update_resource", "test")
	r := MSGraphTestUpdateResource{}

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
	data := acceptance.BuildTestData(t, "msgraph_update_resource", "test")
	r := MSGraphTestUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.withCreateTimeout("Demo App Updated"),
			ExpectError: regexp.MustCompile(`context deadline exceeded`),
		},
	})
}

func TestAcc_UpdateResourceTimeouts_Read(t *testing.T) {
	data := acceptance.BuildTestData(t, "msgraph_update_resource", "test")
	r := MSGraphTestUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.withReadTimeout("Demo App Updated"),
			ExpectError: regexp.MustCompile(`context deadline exceeded`),
		},
	})
}

func TestAcc_UpdateResourceRetry(t *testing.T) {
	data := acceptance.BuildTestData(t, "msgraph_update_resource", "test")

	r := MSGraphTestUpdateResource{}

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
	data := acceptance.BuildTestData(t, "msgraph_update_resource", "test")
	r := MSGraphTestUpdateResource{}

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

func (r MSGraphTestUpdateResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	apiVersion := state.Attributes["api_version"]
	url := state.Attributes["url"]

	_, err := client.MSGraphClient.Read(ctx, url, apiVersion, clients.DefaultRequestOptions())
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

func (r MSGraphTestUpdateResource) basic(displayName string) string {
	return fmt.Sprintf(`
resource "msgraph_resource" "application" {
  url = "applications"
  body = {
    displayName = "Demo App"
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}

resource "msgraph_update_resource" "test" {
  url = "applications/${msgraph_resource.application.id}"
  body = {
    displayName = "%s"
  }
}
`, displayName)
}

func (r MSGraphTestUpdateResource) withUpdateTimeout(displayName string) string {
	return fmt.Sprintf(`
%s

resource "msgraph_update_resource" "test" {
  url = "applications/${msgraph_resource.application.id}"
  body = {
    displayName = "%s"
  }
  timeouts {
    update = "1ns"
  }
}
`, MSGraphTestUpdateResource{}.applicationOnly(), displayName)
}

func (r MSGraphTestUpdateResource) withReadTimeout(displayName string) string {
	return fmt.Sprintf(`
%s

resource "msgraph_update_resource" "test" {
  url = "applications/${msgraph_resource.application.id}"
  body = {
    displayName = "%s"
  }
  timeouts {
    read = "1ns"
  }
}
`, MSGraphTestUpdateResource{}.applicationOnly(), displayName)
}

func (r MSGraphTestUpdateResource) withCreateTimeout(displayName string) string {
	return fmt.Sprintf(`
%s

resource "msgraph_update_resource" "test" {
  url = "applications/${msgraph_resource.application.id}"
  body = {
    displayName = "%s"
  }
  timeouts {
    create = "1ns"
  }
}
`, MSGraphTestUpdateResource{}.applicationOnly(), displayName)
}

// applicationOnly returns just the application resource to be used for composing
// different update resource configurations.
func (r MSGraphTestUpdateResource) applicationOnly() string {
	return `
resource "msgraph_resource" "application" {
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

func (r MSGraphTestUpdateResource) withRetry() string {
	return `
resource "msgraph_resource" "application" {
  url = "applications"
  body = {
    displayName = "Demo App"
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}

resource "msgraph_update_resource" "test" {
  url = "applications/${msgraph_resource.application.id}"
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

func (r MSGraphTestUpdateResource) groupWithOwnerBase() string {
	return `
resource "msgraph_resource" "application" {
  url = "applications"
  body = {
    displayName = "My Application"
  }
  response_export_values = {
    appId = "appId"
  }
}

resource "msgraph_resource" "servicePrincipal_application" {
  url = "servicePrincipals"
  body = {
    appId = msgraph_resource.application.output.appId
  }
}

resource "msgraph_resource" "group" {
  url = "groups"
  body = {
    displayName     = "My Group Owners Bind"
    mailEnabled     = false
    mailNickname    = "mygroup-owners-bind"
    securityEnabled = true
    "owners@odata.bind" = [
      "https://graph.microsoft.com/v1.0/directoryObjects/${msgraph_resource.servicePrincipal_application.id}"
    ]
  }
  lifecycle {
    ignore_changes = [body.displayName]
  }
}
`
}

func (r MSGraphTestUpdateResource) groupWithOwnerUpdate(displayName string) string {
	return fmt.Sprintf(`
%s

resource "msgraph_update_resource" "test" {
  url = "groups/${msgraph_resource.group.id}"
  body = {
    displayName = "%s"
  }
}
`, r.groupWithOwnerBase(), displayName)
}
