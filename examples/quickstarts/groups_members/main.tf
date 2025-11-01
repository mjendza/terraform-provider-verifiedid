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

resource "verifiedid_resource" "member" {
  url = "groups/${verifiedid_resource.group.id}/members/$ref"
  body = {
    "@odata.id" = "https://graph.microsoft.com/v1.0/directoryObjects/${verifiedid_resource.servicePrincipal_application.id}"
  }
}
