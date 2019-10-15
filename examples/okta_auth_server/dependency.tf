resource okta_auth_server test {
  name                      = "testAcc_replace_with_uuid"
  audiences                 = ["api://selfservice_client_1"]
}

resource okta_auth_server test1 {
  name                      = "testAcc_replace_with_uuid1"
  audiences                 = ["api://selfservice_client_2"]
  credentials_rotation_mode = "MANUAL"
}
