---
subcategory: "Reference"
page_title: "{{.resource_type}} - {{.resource_type_friendly_name}}"
description: |-
  Manages a {{.resource_type_friendly_name}}.
---

# {{.resource_type}} - {{.resource_type_friendly_name}}

This article demonstrates how to use `msgraph` provider to manage the {{.resource_type_friendly_name}} resource in MSGraph.

## Example Usage

{{.example}}

## Arguments Reference

The following arguments are supported:

* `url` - (Required) The URL which is used to manage the resource. This should be set to `{{.url}}`.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation]({{.reference_link}}).

* `api_version` - (Optional) The API version used to manage the resource. The default value is `v1.0`. The allowed values are `v1.0` and `beta`.

For other arguments, please refer to the [msgraph_resource](https://registry.terraform.io/providers/Microsoft/msgraph/latest/docs/resources/resource) documentation.

### Read-Only

- `id` (String) The ID of the resource. Normally, it is in the format of UUID.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import msgraph_resource.example {{.resource_id}}
 ```
