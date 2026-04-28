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

func (r VerifiedIDResourceActionDataSourceTestResource) basic() string {
	return `
resource "verifiedid_resource" "application" {
  url = "applications"
  body = {
    displayName = "Test Application"
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}

data "verifiedid_resource_action" "test" {
  resource_url = verifiedid_resource.application.resource_url
  action       = "owners"
  method       = "GET"
}
`
}

func (r VerifiedIDResourceActionDataSourceTestResource) withQueryParams() string {
	return `
resource "verifiedid_resource" "application" {
  url = "applications"
  body = {
    displayName = "Test Application"
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}

data "verifiedid_resource_action" "test" {
  resource_url = verifiedid_resource.application.resource_url
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
resource "verifiedid_resource" "application" {
  url = "applications"
  body = {
    displayName = "Test Application"
  }

  lifecycle {
    ignore_changes = [body.displayName]
  }
}

data "verifiedid_resource_action" "test" {
  resource_url = verifiedid_resource.application.resource_url
  action       = "owners"
  method       = "GET"

  headers = {
    "X-Custom-Header" = "test-value"
  }
}
`
}
