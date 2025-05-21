---
subcategory: "Reference"
page_title: "servicePrincipals - service principal"
description: |-
  Manages a service principal.
---

# servicePrincipals - service principal

This article demonstrates how to use `msgraph` provider to manage the service principal resource in MSGraph.

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

locals {
  MicrosoftGraphAppId = "00000003-0000-0000-c000-000000000000"
}

data "msgraph_resource" "servicePrincipal_msgraph" {
  url = "servicePrincipals"
  query_parameters = {
    "$filter" = ["appId eq '${local.MicrosoftGraphAppId}'"]
  }
  response_export_values = {
    all = "@"
  }
}

resource "msgraph_resource" "application" {
  url = "applications"
  body = {
    displayName = "My Application"
  }
  response_export_values = {
    appId = "appId"
  }
}

resource "msgraph_resource" "servicePrincipal_application" {
  url = "servicePrincipals"
  body = {
    appId = msgraph_resource.application.output.appId
  }
}
```



## Arguments Reference

The following arguments are supported:

* `url` - (Required) The URL which is used to manage the resource. This should be set to `servicePrincipals`.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/servicePrincipals?pivots=deployment-language-terraform).

* `api_version` - (Optional) The API version used to manage the resource. The default value is `v1.0`. The allowed values are `v1.0` and `beta`.

For other arguments, please refer to the [msgraph_resource](https://registry.terraform.io/providers/Microsoft/msgraph/latest/docs/resources/resource) documentation.

### Read-Only

- `id` (String) The ID of the resource. Normally, it is in the format of UUID.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import msgraph_resource.example /servicePrincipals/{servicePrincipals-id}
 ```
