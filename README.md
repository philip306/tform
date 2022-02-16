# TFORM

a quick and dirty wrapper that detects the version of terraform needed and automatically executes with that version.  Like tfswitch, but on windows.

## Pre-requisites

Tform expects that you manage your installed versions of terraform via chocolately and have installed the versions explicitly. e.g. `choco install terraform --version 1.0.5 -my`

Chocolately installs the terraform executables in `c:\programdata\chocolately\lib\terraform.x.x.x\tools\terraform.exe`.  Tform expects this to be the case.

## Usage

main.tf
```
terraform {
  required_version = "~> 1.0.0"
}

```

CMD from the main.tf dir

```
C:\DevOps>tform version
~~Executing with Terraform version: 1.0.5
Terraform v1.0.5
on windows_amd64
```

### Notes

- Will find the most suitable/latest version that matches the constraints from the required_versions parameter, including rightmost(~>)
- uses hashicorp/go-version, hashicorp/terraform-config-inspect, and re-uses some files from the the discovery plugin from terraform itself.
