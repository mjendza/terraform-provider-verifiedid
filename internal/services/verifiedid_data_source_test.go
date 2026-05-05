package services_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance/check"
)

type VerifiedIDTestDataSource struct{}

// TestAcc_DataSourceAuthorities reads the Verified ID authorities collection
// (verifiableCredentials/authorities) — the canonical list endpoint exposed by
// the Microsoft Entra Verified ID Admin API. The data source is verified by
// asserting that the dynamic `output` map is populated.
func TestAcc_DataSourceAuthorities(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.verifiedid_resource", "authorities")
	r := VerifiedIDTestDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.authorities(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.%").Exists(),
			),
		},
	})
}

func TestAcc_DataSourceAuthoritiesRetry(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.verifiedid_resource", "authorities")
	r := VerifiedIDTestDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.authoritiesWithRetry(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.%").Exists(),
			),
		},
	})
}

func TestAcc_DataSourceAuthoritiesTimeouts_Read(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.verifiedid_resource", "authorities")
	r := VerifiedIDTestDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config:      r.authoritiesWithReadTimeout(),
			ExpectError: regexp.MustCompile(`context deadline exceeded`),
		},
	})
}

func (r VerifiedIDTestDataSource) authorities() string {
	return `
data "verifiedid_resource" "authorities" {
  url = "verifiableCredentials/authorities"
  response_export_values = {
    all = "@"
  }
}
`
}

func (r VerifiedIDTestDataSource) authoritiesWithRetry() string {
	return `
data "verifiedid_resource" "authorities" {
  url = "verifiableCredentials/authorities"
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

func (r VerifiedIDTestDataSource) authoritiesWithReadTimeout() string {
	return `
data "verifiedid_resource" "authorities" {
  url = "verifiableCredentials/authorities"
  timeouts {
    read = "1ns"
  }
}
`
}
