terraform {
  required_providers {
    st-byteplus = {
      source = "example.local/myklst/st-byteplus"
    }
  }
}

provider "st-byteplus" {
  access_key = "AKAPMWMzODdjNGVhZjhmNDYyN2FhYWIyM2RjNDdjMDBiODE"
  secret_key = ""
}

data "st-byteplus_cdn_domains" "example" {}
