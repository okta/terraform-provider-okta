resource "okta_app_oauth" "test" {
  label                      = "testAcc_replace_with_uuid"
  type                       = "service"
  response_types             = ["token"]
  grant_types                = ["client_credentials"]
  token_endpoint_auth_method = "private_key_jwt"

  jwks {
    kty = "RSA"
    kid = "SIGNING_KEY_RSA"
    use = "sig"
    e   = "AQAB"
    n   = "owfoXNHcAlAVpIO41840ZU2tZraLGw3yEr3xZvAti7oEZPUKCytk88IDgH7440JOuz8GC_D6vtduWOqnEt0j0_faJnhKHgfj7DTWBOCxzSdjrM-Uyj6-e_XLFvZXzYsQvt52PnBJUV15G1W9QTjlghT_pFrW0xrTtbO1c281u1HJdPd5BeIyPb0pGbciySlx53OqGyxrAxPAt5P5h-n36HJkVsSQtNvgptLyOwWYkX50lgnh2szbJ0_O581bqkNBy9uqlnVeK1RZDQUl4mk8roWYhsx_JOgjpC3YyeXA6hHsT5xWZos_gNx98AHivNaAjzIzvyVItX2-hP0Aoscfff"
  }

  jwks {
    kty = "RSA"
    kid = "ENCRYPTION_KEY_RSA"
    use = "enc"
    e   = "AQAB"
    n   = "qCj2HvA3PMgdJJcFsF9MavDj3C0Q99pj1Tl2hxDUeYv5H8OGVrV4AoMFn8zBt126lRpH0w87CY8RCU1h7clDpXEhN8q6t1taJbTiIu4HwaxknTHQVV2pvuw6j1I-Ra5M1j_t0YBaCdF8_D80hT_EEZzefePQNOUKS4TTL5b60LyXEOQjXKiYZKWC-rZJJL4QzI3I17hVR39VVcVCtC0cQTSw5YJxJtavpfqObh9TsBJn81BeiT4OxFNkJIT5yb-jj6ejESkDCVdfbxdsM7AbNpzacX3VNh_Qs5Gam0wPUiY0oD8JaElzsaEoOIjwNJhgIPEYE-hxxaZ2fhRg0gfRnw"
  }
}

resource "okta_app_oauth" "test_ec" {
  label                      = "test_ecAcc_replace_with_uuid"
  type                       = "service"
  response_types             = ["token"]
  grant_types                = ["client_credentials"]
  token_endpoint_auth_method = "private_key_jwt"

  jwks {
    kty = "EC"
    kid = "SIGNING_KEY_EC"
    use = "sig"
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
