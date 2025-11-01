terraform {
  required_providers {
    verifiedid = {
      source = "mjendza/verifiedid"
    }
  }
}

provider "verifiedid" {
}

resource "verifiedid_resource" "group" {
  url = "groups"
  body = {
    displayName     = "My Group"
    mailEnabled     = false
    mailNickname    = "mygroup"
    securityEnabled = true
  }
}
