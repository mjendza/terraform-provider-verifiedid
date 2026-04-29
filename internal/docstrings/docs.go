package docstrings

import "fmt"

func ApiVersion() string {
	return "The API version of the data source. The allowed values are `v1.0` and `beta`. Defaults to `v1.0`."
}

func Url(kind string) string {
	switch kind {
	case "data":
		return "The URL of the data source. It supports both collection URL which is used to list resources, for example `/verifiableCredentials/authorities/{authority-id}/contracts`, and item URL which is used to read an individual resource, for example `/verifiableCredentials/authorities/{authority-id}/contracts/{contract-id}`."
	case "resource":
		return `The URL which is used to manage the resource. It is the collection URL used to make a POST request to create a new resource, for example, "/verifiableCredentials/authorities/{authority-id}/contracts", and it must support the following operations:
  - Create a new resource: POST "/verifiableCredentials/authorities/{authority-id}/contracts"
  - Read an existing resource: GET "/verifiableCredentials/authorities/{authority-id}/contracts/{contract-id}"
  - Update an existing resource: PATCH "/verifiableCredentials/authorities/{authority-id}/contracts/{contract-id}"
  - Delete an existing resource: HTTP DELETE on the item URL.

  ~> **Soft delete for contracts** The Microsoft Entra Verified ID Admin API does not support HTTP DELETE on contracts. When the resource URL targets a contracts collection (` + "`verifiableCredentials/authorities/{authority-id}/contracts`" + `), this provider performs a soft delete on destroy by issuing PATCH ` + "`{\"status\": \"Disabled\"}`" + ` against the item URL. The contract remains in the tenant in a Disabled state. For all other resources (authorities, Microsoft Graph entities, ` + "`$ref`" + ` relationships) the provider issues a regular HTTP DELETE.

  More information about the Microsoft Entra Verified ID API can be found at [Microsoft Entra Verified ID API overview](https://learn.microsoft.com/en-us/entra/verified-id/admin-api).  
  And there are some [examples](https://github.com/mjendza/terraform-provider-verifiedid/tree/main/examples/quickstarts) to help you get started.
`
	case "update_resource":
		return `The item URL of the existing resource to update, for example "/verifiableCredentials/authorities/{authority-id}/contracts/{contract-id}".

	This resource performs PATCH requests against the item URL and expects the following operations to be supported by the API endpoint:
		- Read an existing resource: GET "/verifiableCredentials/authorities/{authority-id}/contracts/{contract-id}"
		- Update an existing resource: PATCH "/verifiableCredentials/authorities/{authority-id}/contracts/{contract-id}"

	More information about the Microsoft Entra Verified ID API can be found at [Microsoft Entra Verified ID API overview](https://learn.microsoft.com/en-us/entra/verified-id/admin-api).  
	There are also [examples](https://github.com/mjendza/terraform-provider-verifiedid/tree/main/examples/quickstarts) to help you get started.`
	default:
		return ""
	}
}

func Body() string {
	return "A dynamic attribute that contains the request body."
}

func Output() string {
	return fmt.Sprintf(`
The output HCL object containing the properties specified in %[1]sresponse_export_values%[1]s. Here are some examples to use the values.

	%[1]s%[1]s%[1]sterraform
	 output "app_id" {
	   // it will output the value of app_id
	   value = verifiedid_resource.application.output.app_id
	 }
	 
	 output "all" {
	   // it will output the whole response
	   value = verifiedid_resource.application.output.all
	 }
	%[1]s%[1]s%[1]s`, "`")
}

func ResponseExportValues() string {
	return fmt.Sprintf(`A map where the key is the name for the result and the value is a JMESPath query string to filter the response. Here's an example. If it sets to %[1]s{"all" = "@", "app_id" = "appId"}%[1]s, it will set the following HCL object to the computed property output.

	%[1]s%[1]s%[1]stext
	{
		"all" = {
			"appId" = "00000000-0000-0000-0000-000000000000"
			"displayName" = "example"
			"id" = "00000000-0000-0000-0000-000000000000"
			...
		}
		"app_id" = "00000000-0000-0000-0000-000000000000"
	}
	%[1]s%[1]s%[1]s

To learn more about JMESPath, visit [JMESPath](https://jmespath.org/).
`, "`")
}

func ResourceID() string {
	return "The ID of the resource. Normally, it is in the format of UUID."
}

func IgnoreMissingProperty() string {
	return "Whether ignore not returned properties like credentials in `body` to suppress plan-diff. Defaults to `true`. It's recommend to enable this option when some sensitive properties are not returned in response body, instead of setting them in `lifecycle.ignore_changes` because it will make the sensitive fields unable to update."
}
