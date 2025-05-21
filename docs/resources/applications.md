---
subcategory: "Reference"
page_title: "applications - application registration"
description: |-
  Manages a application registration.
---

# applications - application registration

This article demonstrates how to use `msgraph` provider to manage the application registration resource in MSGraph.

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

```



## Arguments Reference

The following arguments are supported:

* `url` - (Required) The URL which is used to manage the resource. This should be set to `applications`.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/applications?pivots=deployment-language-terraform).

* `api_version` - (Optional) The API version used to manage the resource. The default value is `v1.0`. The allowed values are `v1.0` and `beta`.

For other arguments, please refer to the [msgraph_resource](https://registry.terraform.io/providers/Microsoft/msgraph/latest/docs/resources/resource) documentation.

### Read-Only

- `id` (String) The ID of the resource. Normally, it is in the format of UUID.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import msgraph_resource.example /applications/{applications-id}
 ```
