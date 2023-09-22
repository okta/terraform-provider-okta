resource "okta_app_oauth" "test" {
  label                      = "testAcc_replace_with_uuid"
  type                       = "service"
  response_types             = ["token"]
  grant_types                = ["client_credentials"]
  token_endpoint_auth_method = "private_key_jwt"

  jwks {
    kty = "RSA"
    kid = "SIGNING_KEY_RSA"
    e   = "AQAB"
    n   = "owfoXNHcAlAVpIO41840ZU2tZraLGw3yEr3xZvAti7oEZPUKCytk88IDgH7440JOuz8GC_D6vtduWOqnEt0j0_faJnhKHgfj7DTWBOCxzSdjrM-Uyj6-e_XLFvZXzYsQvt52PnBJUV15G1W9QTjlghT_pFrW0xrTtbO1c281u1HJdPd5BeIyPb0pGbciySlx53OqGyxrAxPAt5P5h-n36HJkVsSQtNvgptLyOwWYkX50lgnh2szbJ0_O581bqkNBy9uqlnVeK1RZDQUl4mk8roWYhsx_JOgjpC3YyeXA6hHsT5xWZos_gNx98AHivNaAjzIzvyVItX2-hP0Aoscfff"
  }
}

resource "okta_app_oauth" "test2" {
  label                      = "test2Acc_replace_with_uuid"
  type                       = "service"
  response_types             = ["token"]
  grant_types                = ["client_credentials"]
  token_endpoint_auth_method = "private_key_jwt"

  jwks {
    kty = "EC"
    kid = "SIGNING_KEY_EC"
    x   = "K37X78mXJHHldZYMzrwipjKR-YZUS2SMye0KindHp6I"
    y   = "8IfvsvXWzbFWOZoVOMwgF5p46mUj3kbOVf9Fk0vVVHo"
  }
}

# Test EC Key
# {
#     "kty": "EC",
#     "use": "sig",
#     "crv": "P-256",
#     "kid": "testing",
#     "x": "K37X78mXJHHldZYMzrwipjKR-YZUS2SMye0KindHp6I",
#     "y": "8IfvsvXWzbFWOZoVOMwgF5p46mUj3kbOVf9Fk0vVVHo",
#     "alg": "ES256"
# }
