resource "okta_authenticator" "test1" {
key    = "custom_app"
name =  "VCRTestCustomAppAuthRenamed"
status = "ACTIVE"
agree_to_terms = "true"
legacy_ignore_name = false
settings = jsonencode({
    "userVerification": "PREFERRED",
    "appInstanceId": "0oaspGHIJKL12345678"
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
