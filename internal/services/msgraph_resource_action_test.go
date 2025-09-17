package services_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance/check"
	"github.com/mjendza/terraform-provider-verifiedid/internal/clients"
)

type VerifiedIDResourceActionTestResource struct{}

func TestAcc_ResourceActionBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource_action", "test")

	r := VerifiedIDResourceActionTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAcc_ResourceActionWithQueryParams(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource_action", "test")

	r := VerifiedIDResourceActionTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withQueryParams(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAcc_ResourceActionWithHeaders(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource_action", "test")

	r := VerifiedIDResourceActionTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withHeaders(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAcc_ResourceActionWithExportValues(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource_action", "test")

	r := VerifiedIDResourceActionTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withExportValues(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.%").HasValue("2"),
			),
		},
	})
}

func (r VerifiedIDResourceActionTestResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	exists := false
	return &exists, nil
}

func (r VerifiedIDResourceActionTestResource) basic() string {
	return `
provider "msgraph" {}

resource "verifiedid_resource" "group" {
  url = "groups"
  body = {
    displayName     = "Test Group"
    mailEnabled     = false
    mailNickname    = "mygroup"
    securityEnabled = true
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}

resource "verifiedid_resource_action" "test" {
  resource_url = verifiedid_resource.group.resource_url
  method       = "PATCH"

  body = {
    displayName = "Updated Group Name"
  }
}
`
}

func (r VerifiedIDResourceActionTestResource) withQueryParams() string {
	return `
provider "msgraph" {}

resource "verifiedid_resource" "group" {
  url = "groups"
  body = {
    displayName     = "Test Group"
    mailEnabled     = false
    mailNickname    = "mygroup"
    securityEnabled = true
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}

resource "verifiedid_resource_action" "test" {
  resource_url = verifiedid_resource.group.resource_url
  method       = "PATCH"

  query_parameters = {
    "$select" = ["id", "displayName"]
  }

  body = {
    displayName = "Updated Group Name with Query Params"
  }
}
`
}

func (r VerifiedIDResourceActionTestResource) withHeaders() string {
	return `
provider "msgraph" {}

resource "verifiedid_resource" "group" {
  url = "groups"
  body = {
    displayName     = "Test Group"
    mailEnabled     = false
    mailNickname    = "mygroup"
    securityEnabled = true
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}

resource "verifiedid_resource_action" "test" {
  resource_url = verifiedid_resource.group.resource_url
  method       = "PATCH"

  headers = {
    "X-Custom-Header" = "test-value"
    "X-Request-ID"    = "test-123"
  }

  body = {
    displayName = "Updated Group Name with Headers"
  }
}
`
}

func (r VerifiedIDResourceActionTestResource) withExportValues() string {
	return `
provider "msgraph" {}

resource "verifiedid_resource" "group" {
  url = "groups"
  body = {
    displayName     = "Test Group"
    mailEnabled     = false
    mailNickname    = "mygroup"
    securityEnabled = true
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}

resource "verifiedid_resource_action" "test" {
  resource_url = verifiedid_resource.group.resource_url
  method       = "PATCH"

  body = {
    displayName = "Updated Group Name with Export Values"
  }

  response_export_values = {
    group_id   = "id"
    group_name = "displayName"
  }
}
`
}
