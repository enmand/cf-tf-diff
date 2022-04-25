# cf-tf-diff

cf-tf-diff is a utility for diffing the Cloudflare API against a given Terraform project to ensure
no configuration drift occurs between them. If drift exist, the command will exit with an error code
and print the differences.

## Installation

```bash
go get github.com/enmand/cf-tf-diff
```

## Usage

```bash
cf-tf-diff -p /path/to/terraform/project -c cloudflare-email -k cloudflare-key
```
