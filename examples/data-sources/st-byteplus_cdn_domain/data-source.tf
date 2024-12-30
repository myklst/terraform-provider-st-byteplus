data "st-byteplus_cdn_domain" "example" {
  domain_name = "www.example.com"

  client_config {
    access_key = ""
    secret_key = ""
  }
}

output "cdn_domain_details" {
  value = data.st-byteplus_cdn_domain.example
}
