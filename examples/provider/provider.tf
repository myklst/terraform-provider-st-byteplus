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
  domain_name = "www.example.com"

    client_config {
      access_key = "AKAPZjJjZWQwNzdkNTg1NGQwNzgyYTdhNzM4MjRiM2RmMz"
      secret_key = ""
    }
}
