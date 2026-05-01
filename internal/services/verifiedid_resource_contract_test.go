package services_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance"
	"github.com/mjendza/terraform-provider-verifiedid/internal/acceptance/check"
	"github.com/mjendza/terraform-provider-verifiedid/internal/clients"
	"github.com/mjendza/terraform-provider-verifiedid/internal/services"
	"github.com/mjendza/terraform-provider-verifiedid/internal/utils"
)

// envVerifiedIDAuthorityID is the optional environment variable that points
// the contract acceptance tests at an existing Verified ID authority.
// Creating an authority requires Key Vault setup, so for the contract
// soft-delete coverage we re-use a tenant-provided authority instead of
// provisioning one inline.
const envVerifiedIDAuthorityID = "ARM_VERIFIEDID_AUTHORITY_ID"

type VerifiedIDContractTestResource struct {
	authorityID string
}

// TestAcc_ContractUpdate pins the create -> update flow for a Verified ID
// contract to the JSON shape used by the PowerShell reference under
// reference/stage-101-vc/ (contracts.lifecycle.ps1 + contracts.create.sample.json
// + contracts.update.sample.json). It exercises the patch_as_full_body=true
// PATCH path: the Verified ID Admin API rejects partial PATCH bodies, so the
// update step must send the full body (minus immutable / server-owned fields:
// id, name, status, manifestUrl, issuerId, authorityId,
// issueNotificationAllowedToGroupOids).
func TestAcc_ContractUpdate(t *testing.T) {
	authorityID := os.Getenv(envVerifiedIDAuthorityID)
	if authorityID == "" {
		t.Skipf("skipping: %s is not set", envVerifiedIDAuthorityID)
	}

	data := acceptance.BuildTestData(t, "verifiedid_resource", "contract")
	r := VerifiedIDContractTestResource{authorityID: authorityID}

	pattern := fmt.Sprintf(
		`^verifiableCredentials/authorities/%s/contracts/[^/]+$`,
		regexp.QuoteMeta(authorityID),
	)
	urlRegex := regexp.MustCompile(pattern)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(urlRegex),
			),
		},
		{
			Config: r.basicUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(urlRegex),
			),
		},
	})
}

// TestAcc_ContractSoftDelete verifies that destroying a contract resource
// leaves the contract in place with status=Disabled rather than issuing an
// HTTP DELETE — the Microsoft Entra Verified ID Admin API does not support
// HTTP DELETE on contracts.
//
// See https://learn.microsoft.com/en-us/entra/verified-id/admin-api#contracts
func TestAcc_ContractSoftDelete(t *testing.T) {
	authorityID := os.Getenv(envVerifiedIDAuthorityID)
	if authorityID == "" {
		t.Skipf("skipping: %s is not set", envVerifiedIDAuthorityID)
	}

	data := acceptance.BuildTestData(t, "verifiedid_resource", "contract")
	r := VerifiedIDContractTestResource{authorityID: authorityID}

	pattern := fmt.Sprintf(
		`^verifiableCredentials/authorities/%s/contracts/[^/]+$`,
		regexp.QuoteMeta(authorityID),
	)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Exists(r),
				check.That(data.ResourceName).Key("resource_url").MatchesRegex(regexp.MustCompile(pattern)),
			),
		},
	})
}

// Exists treats a contract whose status has been flipped to Disabled as
// "no longer present". This is what soft delete means for contracts.
func (r VerifiedIDContractTestResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	apiVersion := state.Attributes["api_version"]
	url := state.Attributes["url"]
	checkUrl := fmt.Sprintf("%s/%s", url, state.ID)

	body, err := client.VerifiedIDClient.Read(ctx, checkUrl, apiVersion, clients.DefaultRequestOptions())
	if err == nil {
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

func (r VerifiedIDContractTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "verifiedid_resource" "contract" {
  url                = "verifiableCredentials/authorities/%[1]s/contracts"
  patch_as_full_body = true
  body = {
    name = "tfacc-%[2]s"
    rules = {
      attestations = {
        idTokenHints = [
          {
            mapping = [
              {
                outputClaim = "certNumber"
                required    = true
                inputClaim  = "certNumber"
                indexed     = false
              }
            ]
            required = true
          }
        ]
      }
      validityInterval = 2592000
      vc = {
        type = ["TfAccDemo"]
      }
    }
    displays = [
      {
        locale = "en-US"
        card = {
          backgroundColor = "#BDD0A7"
          description     = "tfacc demo"
          issuedBy        = "tfacc"
          textColor       = "#000000"
          title           = "tfacc title"
          logo = {
            description = "logo"
            uri         = "https://example.com/logo.png"
          }
        }
        consent = {
          instructions = "tfacc"
          title        = "tfacc"
        }
        claims = [
          {
            claim = "vc.credentialSubject.certNumber"
            label = "Certificate Number"
            type  = "String"
          }
        ]
      }
    ]
  }
}
`, r.authorityID, data.RandomString)
}

// basicUpdate mirrors reference/stage-101-vc/contracts.update.sample.json: the
// PATCH body OMITS the immutable "name" field (the Admin API rejects PATCH
// bodies that echo it back) and adds availableInVcDirectory +
// allowOverrideValidityIntervalOnIssuance. The display description visibly
// changes so the round-trip is observable.
func (r VerifiedIDContractTestResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "verifiedid_resource" "contract" {
  url                = "verifiableCredentials/authorities/%[1]s/contracts"
  patch_as_full_body = true
  body = {
    rules = {
      attestations = {
        idTokenHints = [
          {
            mapping = [
              {
                outputClaim = "certNumber"
                required    = true
                inputClaim  = "certNumber"
                indexed     = false
              }
            ]
            required = true
          }
        ]
      }
      validityInterval = 2592000
      vc = {
        type = ["TfAccDemo"]
      }
    }
    displays = [
      {
        locale = "en-US"
        card = {
          backgroundColor = "#BDD0A7"
          description     = "tfacc demo (updated)"
          issuedBy        = "tfacc"
          textColor       = "#000000"
          title           = "tfacc title"
          logo = {
            description = "logo"
            uri         = "https://example.com/logo.png"
          }
        }
        consent = {
          instructions = "tfacc"
          title        = "tfacc"
        }
        claims = [
          {
            claim = "vc.credentialSubject.certNumber"
            label = "Certificate Number"
            type  = "String"
          }
        ]
      }
    ]
    availableInVcDirectory                  = false
    allowOverrideValidityIntervalOnIssuance = true
  }
}
`, r.authorityID, data.RandomString)
}
