terraform {
  required_providers {
    verifiedid = {
      source = "mjendza/verifiedid"
    }
  }
}

provider "verifiedid" {
}

resource "verifiedid_resource" "application" {
  url = "applications"
  body = {
    displayName = "My Application"
  }
}
