terraform {
  required_providers {
    verifiedid = {
      source = "mjendza/verifiedid"
    }
  }
}

provider "verifiedid" {
}

locals {
  MicrosoftGraphAppId = "00000003-0000-0000-c000-000000000000"
}

data "verifiedid_resource" "servicePrincipal_msgraph" {
  url = "servicePrincipals"
  query_parameters = {
    "$filter" = ["appId eq '${local.MicrosoftGraphAppId}'"]
  }
  response_export_values = {
    all = "@"
  }
}

resource "verifiedid_resource" "application" {
  url = "applications"
  body = {
    displayName = "My Application"
  }
  response_export_values = {
    appId = "appId"
  }
}

resource "verifiedid_resource" "servicePrincipal_application" {
  url = "servicePrincipals"
  body = {
    appId = verifiedid_resource.application.output.appId
  }
}