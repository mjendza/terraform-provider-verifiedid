terraform {
  required_providers {
    verifiedid = {
      source = "mjendza/verifiedid"
    }
  }
}

provider "verifiedid" {}

# This example creates an application first (so we have something to update),
# then uses verifiedid_update_resource to PATCH its displayName.

resource "verifiedid_resource" "application" {
  url = "applications"
  body = {
    displayName = "Demo App"
  }

  # We ignore the displayName change here because the update is handled by
  # the separate verifiedid_update_resource below.
  lifecycle {
    ignore_changes = [body.displayName]
  }
}

resource "verifiedid_update_resource" "application_update" {
  # Point directly at the item URL you want to PATCH
  url = "applications/${verifiedid_resource.application.id}"

  body = {
    displayName = "Demo App Updated"
  }
}