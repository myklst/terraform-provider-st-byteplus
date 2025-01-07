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

resource "st-byteplus_iam_policy" "name" {
  user_name         = "lq-user-1"
  attached_policies = ["VPCFullAccess", "TOSFullAccess", "lqtestpolicy"]
}

