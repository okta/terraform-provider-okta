resource "okta_authenticator" "test1" {
key    = "custom_app"
name =  "VCRTestCustomAppAuth"
status = "ACTIVE"
agree_to_terms = "true"
legacy_ignore_name = false
settings = jsonencode({
    "userVerification": "REQUIRED",
    "appInstanceId": "0oaspABCDEF12345678"
    })
provider_json = jsonencode({
    "type": "PUSH",
    "configuration": {
        "fcm": {
            "id": "ppcrb12345678ABCDEF"
            }
        }      
    })
}
