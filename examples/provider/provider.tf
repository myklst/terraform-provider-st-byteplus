terraform {
  required_providers {
    st-byteplus = {
      source = "example.local/myklst/st-byteplus"
    }
  }
}

provider "st-byteplus" {
  region = "ap-singapore-1"
}

data "st-byteplus_cdn_domain" "example" {
  domain_name = "test.example.com"

  client_config {
    access_key = ""
    secret_key = ""
  }
}
