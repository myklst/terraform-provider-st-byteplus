terraform {
  required_providers {
    st-byteplus = {
      source = "myklst/st-byteplus"
    }
  }
}

provider "st-byteplus" {
  region = "ap-singapore-1"
}
