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
  response_export_values = {
    all    = "@"
    app_id = "appId"
  }
}

output "app_id" {
  value = verifiedid_resource.application.output.app_id
}

output "all" {
  // it will output the whole response
  value = verifiedid_resource.application.output.all
}

output "resource_url" {
  // it will output something like "applications/12345678-1234-1234-1234-123456789abc"
  value = verifiedid_resource.application.resource_url
}
