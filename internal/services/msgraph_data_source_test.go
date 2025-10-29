package services_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance/check"
)

type VerifiedIDTestDataSource struct{}

func TestAcc_DataSourceBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.verifiedid_resource", "test")
	r := VerifiedIDTestDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.%").Exists(),
				check.That(data.ResourceName).Key("id").IsUUID(),
			),
		},
	})
}

func TestAcc_DataSourceQueryParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.verifiedid_resource", "test")
	r := VerifiedIDTestDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.query(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.%").Exists(),
			),
		},
	})
}

func TestAcc_DataSourceList(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.verifiedid_resource", "test")
	r := VerifiedIDTestDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.list(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.%").Exists(),
			),
		},
	})
}

func TestAcc_DataSourceRetry(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.verifiedid_resource", "test")
	r := VerifiedIDTestDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.withRetry(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.%").Exists(),
			),
		},
	})
}

func TestAcc_DataSourceTimeouts_Read(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.verifiedid_resource", "test")
	r := VerifiedIDTestDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config:      r.withReadTimeout(data),
			ExpectError: regexp.MustCompile(`context deadline exceeded`),
		},
	})
}

func (r VerifiedIDTestDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "verifiedid_resource" "test" {
  url = "applications/${verifiedid_resource.test.id}"
}
`, VerifiedIDTestResource{}.basic(data))
}

func (r VerifiedIDTestDataSource) query(data acceptance.TestData) string {
	return `
locals {
  MicrosoftGraphAppId = "00000003-0000-0000-c000-000000000000"
}

data "verifiedid_resource" "test" {
  url = "servicePrincipals"
  query_parameters = {
    "$filter" = ["appId eq '${local.MicrosoftGraphAppId}'"]
  }
  response_export_values = {
    all = "@"
  }
}`
}

func (r VerifiedIDTestDataSource) list(data acceptance.TestData) string {
	return `
data "verifiedid_resource" "test" {
  url = "groups"
  response_export_values = {
    all = "@"
  }
}`
}

func (r VerifiedIDTestDataSource) withRetry(data acceptance.TestData) string {
	return `
data "verifiedid_resource" "test" {
  url = "groups"
  retry = {
    error_message_regex = [
      "temporary error",
      ".*throttl.*",
    ]
  }
  response_export_values = {
    all = "@"
  }
}
`
}

func (r VerifiedIDTestDataSource) withReadTimeout(data acceptance.TestData) string {
	return VerifiedIDTestResource{}.basic(data) + "\n" + r.dataSourceWithReadTimeout()
}

func (r VerifiedIDTestDataSource) dataSourceWithReadTimeout() string {
	return `
data "verifiedid_resource" "test" {
  url = "applications/${verifiedid_resource.test.id}"
  timeouts {
    read = "1ns"
  }
}`
}
