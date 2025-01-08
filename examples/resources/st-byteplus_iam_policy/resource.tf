resource "st-byteplus_iam_policy" "name" {
  user_name         = "lq-user-1"
  attached_policies = ["VPCFullAccess", "TOSFullAccess", "lqtestpolicy"]
}
