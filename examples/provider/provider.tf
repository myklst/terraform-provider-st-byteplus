terraform {
  required_providers {
    st-byteplus = {
      source = "example.local/myklst/st-byteplus"
    }
  }
}

provider "st-byteplus" {
  region     = "ap-singapore-1"
  access_key = "NOT_USED"
  secret_key = "NOT_USED"
}

data "st-byteplus_cdn_domain" "example" {
  domain_name = "public1.sige-test4.com"

  client_config {
    access_key = "AKAPZjJjZWQwNzdkNTg1NGQwNzgyYTdhNzM4MjRiM2RmMzc"
    secret_key = ""
  }
}

output "cdn_domain_cname" {
  value = data.st-byteplus_cdn_domain.example
}
