resource okta_social_idp facebook {
  type          = "FACEBOOK"
  protocol_type = "OAUTH2"
  name          = "facebook_replace_with_uuid"

  scopes = [
    "public_profile",
    "email",
  ]

  client_id         = "abcd123"
  client_secret     = "abcd123"
  username_template = "idpuser.email"
}

resource okta_social_idp google {
  type          = "GOOGLE"
  protocol_type = "OAUTH2"
  name          = "google_replace_with_uuid"

  scopes = [
    "profile",
    "email",
    "openid",
  ]

  client_id         = "abcd123"
  client_secret     = "abcd123"
  username_template = "idpuser.email"
}

resource okta_social_idp microsoft {
  type          = "MICROSOFT"
  protocol_type = "OIDC"
  name          = "microsoft_replace_with_uuid"

  scopes = [
    "openid",
    "email",
    "profile",
    "https://graph.microsoft.com/User.Read",
  ]

  client_id         = "abcd123"
  client_secret     = "abcd123"
  username_template = "idpuser.userPrincipalName"
}
