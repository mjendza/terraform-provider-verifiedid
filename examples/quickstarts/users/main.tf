terraform {
  required_providers {
    verifiedid = {
      source = "mjendza/verifiedid"
    }
  }
}

provider "verifiedid" {
}

data "verifiedid_resource" "domains" {
  url = "domains"
  response_export_values = {
    all = "@"
  }
}

locals {
  domain = one([for domain in data.verifiedid_resource.domains.output.all.value : domain.id if domain.isInitial])
}

resource "verifiedid_resource" "user" {
  url = "users"
  body = {
    accountEnabled    = false
    displayName       = "My User"
    mailNickname      = "myuser"
    userPrincipalName = "myuser@${local.domain}"
    passwordProfile = {
      forceChangePasswordNextSignIn = true
      password                      = "Str0ngP@ssword"
    }
  }
}
