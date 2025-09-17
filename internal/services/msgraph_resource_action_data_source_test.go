package services_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance"
	"github.com/mjendza/terraform-provider-verifiedid/internal/clients"
)

type VerifiedIDResourceActionDataSourceTestResource struct{}

func (VerifiedIDResourceActionDataSourceTestResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	exists := true
	return &exists, nil
}

func TestAcc_DataSourceResourceActionBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource_action", "test")

	r := VerifiedIDResourceActionDataSourceTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAcc_DataSourceResourceActionWithQueryParams(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource_action", "test")

	r := VerifiedIDResourceActionDataSourceTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withQueryParams(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAcc_DataSourceResourceActionWithHeaders(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource_action", "test")

	r := VerifiedIDResourceActionDataSourceTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withHeaders(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAcc_DataSourceResourceActionWithBody(t *testing.T) {
	data := acceptance.BuildTestData(t, "verifiedid_resource_action", "test")

	r := VerifiedIDResourceActionDataSourceTestResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withBody(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func (r VerifiedIDResourceActionDataSourceTestResource) basic() string {
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

data "verifiedid_resource_action" "test" {
  resource_url = verifiedid_resource.group.resource_url
  action       = "members"
  method       = "GET"
}
`
}

func (r VerifiedIDResourceActionDataSourceTestResource) withQueryParams() string {
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

data "verifiedid_resource_action" "test" {
  resource_url = verifiedid_resource.group.resource_url
  action       = "owners"
  method       = "GET"

  query_parameters = {
    "$select" = ["id", "displayName"]
    "$top"    = ["5"]
  }
}
`
}

func (r VerifiedIDResourceActionDataSourceTestResource) withHeaders() string {
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

data "verifiedid_resource_action" "test" {
  resource_url = verifiedid_resource.group.resource_url
  action       = "owners"
  method       = "GET"

  headers = {
    "X-Custom-Header" = "test-value"
  }
}
`
}

func (r VerifiedIDResourceActionDataSourceTestResource) withBody() string {
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

data "verifiedid_resource_action" "test" {
  resource_url = verifiedid_resource.group.resource_url
  action       = "checkMemberObjects"
  method       = "POST"

  body = {
    ids = ["00000000-0000-0000-0000-000000000000"]
  }

  response_export_values = {
    result = "value"
  }
}
`
}
