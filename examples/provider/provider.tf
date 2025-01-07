terraform {
  required_providers {
    st-byteplus = {
      source = "example.local/myklst/st-byteplus"
    }
  }
}

provider "st-byteplus" {
  region     = "ap-singapore-1"
}

# resource "st-byteplus_iam_policy" "name" {}

# data "st-byteplus_cdn_domain" "example" {
#   domain_name = "www.example.com"
# }
