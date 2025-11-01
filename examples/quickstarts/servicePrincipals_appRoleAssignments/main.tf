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

  # AppRoleAssignment
  userReadAllAppRoleId = one([for role in data.verifiedid_resource.servicePrincipal_msgraph.output.all.value[0].appRoles : role.id if role.value == "User.Read.All"])
  userReadWriteRoleId  = one([for role in data.verifiedid_resource.servicePrincipal_msgraph.output.all.value[0].oauth2PermissionScopes : role.id if role.value == "User.ReadWrite"])

  # ServicePrincipal
  MSGraphServicePrincipalId         = data.verifiedid_resource.servicePrincipal_msgraph.output.all.value[0].id
  TestApplicationServicePrincipalId = verifiedid_resource.servicePrincipal_application.output.all.id
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
    requiredResourceAccess = [
      {
        resourceAppId = local.MicrosoftGraphAppId
        resourceAccess = [
          {
            id   = local.userReadAllAppRoleId
            type = "Scope"
          },
          {
            id   = local.userReadWriteRoleId
            type = "Scope"
          }
        ]
      }
    ]
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
  response_export_values = {
    all = "@"
  }
}

resource "verifiedid_resource" "appRoleAssignment" {
  url = "servicePrincipals/${local.MSGraphServicePrincipalId}/appRoleAssignments"
  body = {
    appRoleId   = local.userReadAllAppRoleId
    principalId = local.TestApplicationServicePrincipalId
    resourceId  = local.MSGraphServicePrincipalId
  }
}
