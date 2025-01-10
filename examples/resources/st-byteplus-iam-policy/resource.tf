terraform {
  required_providers {
    st-byteplus = {
      source = "example.local/myklst/st-byteplus"
    }
  }
}

provider "st-byteplus" {
  region = "cn-hongkong"
}

resource "st-byteplus_iam_policy" "name" {
  user_name = "lq-user-1" //temporary user for testing. Please either create said user or remove them from this list.
  attached_policies = [
    "VPCFullAccess", "TOSReadOnlyAccess", "VodReadOnlyAccess", "EIPReadOnlyAccess",
    "AccountProfileFullAccess", "BillingCenterReadOnlyAccess", "VODQualityFullAccess", "VPCReadOnlyAccess",
    "CBRReadOnlyAccess", "CRBReadOnlyAccess", "ECSReadOnlyAccess", "IAMFullAccess", "RTCReadOnlyAccess",
    "LiveSaaSReadOnlyAccess", "ImageXReadOnlyAccess", "LIVEReadOnlyAccess", "ESCloudReadOnlyAccess",
    "MSEReadOnlyAccess", "TLSReadOnlyAccess", "AutoScalingReadOnlyAccess", "CFSReadOnlyAccess",
  "LqTestPolicy"] //custom policies used for testing. Please either create said policies or remove them from this list.
}
