terraform {
  required_providers {
    msgraph = {
      source = "Microsoft/msgraph"
    }
  }
}

provider "msgraph" {
}

variable "application_id" {
  type    = string
  default = "00000000-0000-0000-0000-000000000000"
}

data "msgraph_resource" "application" {
  url = "applications/${var.application_id}"
  response_export_values = {
    all          = "@"
    display_name = "displayName"
  }
}

output "display_name" {
  // it will output "John Doe"
  value = data.msgraph_resource.application.output.display_name
}

output "all" {
  // it will output the whole response
  value = data.msgraph_resource.application.output.all
}
