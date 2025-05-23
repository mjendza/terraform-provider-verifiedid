---
subcategory: "Reference"
page_title: "applications/federatedIdentityCredentials - federated identity credentials associated with an application"
description: |-
  Manages a federated identity credentials associated with an application.
---

# applications/federatedIdentityCredentials - federated identity credentials associated with an application

This article demonstrates how to use `msgraph` provider to manage the federated identity credentials associated with an application resource in MSGraph.

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

resource "msgraph_resource" "application" {
  url = "applications"
  body = {
    displayName = "My Application"
  }
}

resource "msgraph_resource" "federatedIdentityCredential" {
  # url = "applications/{id}/federatedIdentityCredentials"
  url = "applications/${msgraph_resource.application.id}/federatedIdentityCredentials"
  body = {
    name        = "myFederatedIdentityCredentials"
    description = "My test federated identity credentials"
    audiences   = ["https://myapp.com"]
    issuer      = "https://sts.windows.net/00000000-0000-0000-0000-000000000000/"
    subject     = "00000000-0000-0000-0000-000000000000"
  }
}

```



## Arguments Reference

The following arguments are supported:

* `url` - (Required) The URL which is used to manage the resource. This should be set to `applications/{application-id}/federatedIdentityCredentials`.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/applications/federatedIdentityCredentials?pivots=deployment-language-terraform).

* `api_version` - (Optional) The API version used to manage the resource. The default value is `v1.0`. The allowed values are `v1.0` and `beta`.

For other arguments, please refer to the [msgraph_resource](https://registry.terraform.io/providers/Microsoft/msgraph/latest/docs/resources/resource) documentation.

### Read-Only

- `id` (String) The ID of the resource. Normally, it is in the format of UUID.

## Import

 ```shell
 # MSGraph resource can be imported using the resource id, e.g.
 terraform import msgraph_resource.example /applications/{application-id}/federatedIdentityCredentials/{federatedIdentityCredentials-id}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import msgraph_resource.example /applications/{application-id}/federatedIdentityCredentials/{federatedIdentityCredentials-id}?api-version=v1.0
 ```
