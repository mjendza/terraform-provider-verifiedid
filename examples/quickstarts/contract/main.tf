terraform {
  required_providers {
    verifiedid = {
      source = "mjendza/verifiedid"
    }
  }
}

provider "verifiedid" {
}

resource "verifiedid_resource" "credential_cert1" {
    patch_as_full_body = true
    #api_version = "beta"
    url  = "verifiableCredentials/authorities/____YOUR_GUID_HERE_FOR_AUTHORITY____/contracts"
    body = jsondecode(file("contract-demo.json"))
}
