terraform {
  required_providers {
    st-byteplus = {
      source = "example.local/myklst/st-byteplus"
    }
  }
}

provider "st-byteplus" {

}

data "st-byteplus_cdn_domain" "example" {}
