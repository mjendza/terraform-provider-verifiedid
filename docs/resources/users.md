---
subcategory: "Reference"
page_title: "users - user"
description: |-
  Manages a user.
---

# users - user

This article demonstrates how to use `msgraph` provider to manage the user resource in MSGraph.

## Example Usage

### default

```hcl
terraform {
  required_providers {
    msgraph = {
      source = "microsoft/msgraph"
    }
  }
}

provider "msgraph" {
}

data "msgraph_resource" "domains" {
  url = "domains"
  response_export_values = {
    all = "@"
  }
}

locals {
  domain = one([for domain in data.msgraph_resource.domains.output.all.value : domain.id if domain.isInitial])
}

resource "msgraph_resource" "user" {
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

```



## Arguments Reference

The following arguments are supported:

* `url` - (Required) The URL which is used to manage the resource. This should be set to `users`.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/users?pivots=deployment-language-terraform).

* `api_version` - (Optional) The API version used to manage the resource. The default value is `v1.0`. The allowed values are `v1.0` and `beta`.

For other arguments, please refer to the [msgraph_resource](https://registry.terraform.io/providers/Microsoft/msgraph/latest/docs/resources/resource) documentation.

### Read-Only

- `id` (String) The ID of the resource. Normally, it is in the format of UUID.

## Import

 ```shell
 # MSGraph resource can be imported using the resource id, e.g.
 terraform import msgraph_resource.example /users/{users-id}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import msgraph_resource.example /users/{users-id}?api-version=v1.0
 ```
