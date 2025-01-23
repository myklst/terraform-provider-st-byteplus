resource "st-byteplus_iam_policy" "name" {
  user_name = "devopsuser01"
  attached_policies = ["VPCFullAccess", "TOSReadOnlyAccess", "VodReadOnlyAccess", "IAMFullAccess",]
}
