terraform {
  required_providers {
    st-byteplus = {
      source = "myklst/st-byteplus"
    }
  }
}

provider "st-byteplus" {
}

data "st-byteplus_cdn_domain" "def" {
  domain_name = "test.example.com"
}
