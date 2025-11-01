# Verified ID Terraform Provider

This is a fork from https://github.com/microsoft/terraform-provider-msgraph

Alpha version of a Terraform provider for Microsoft Entra Verified ID.


## DEVELOPMENT
### Linux

BUILD
```
go build -o terraform-provider-verifiedid
```

COPY to local terraform plugin folder
```
cp terraform-provider-verifiedid ~/.terraform.d/plugins/local/custom/verifiedid/1.0.0/linux_amd64/
```

review http request to to the graph api
```
export TF_LOG=DEBUG
```

```
export TF_LOG=WARN
```