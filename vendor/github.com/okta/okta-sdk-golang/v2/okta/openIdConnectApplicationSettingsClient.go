/*
* Copyright 2018 - Present Okta, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

// AUTO-GENERATED!  DO NOT EDIT FILE DIRECTLY

package okta

import ()

type OpenIdConnectApplicationSettingsClient struct {
	ApplicationType        string               `json:"application_type,omitempty"`
	ClientUri              string               `json:"client_uri,omitempty"`
	ConsentMethod          string               `json:"consent_method,omitempty"`
	GrantTypes             []*OAuthGrantType    `json:"grant_types,omitempty"`
	InitiateLoginUri       string               `json:"initiate_login_uri,omitempty"`
	IssuerMode             string               `json:"issuer_mode,omitempty"`
	LogoUri                string               `json:"logo_uri,omitempty"`
	PolicyUri              string               `json:"policy_uri,omitempty"`
	PostLogoutRedirectUris []string             `json:"post_logout_redirect_uris,omitempty"`
	RedirectUris           []string             `json:"redirect_uris,omitempty"`
	ResponseTypes          []*OAuthResponseType `json:"response_types,omitempty"`
	TosUri                 string               `json:"tos_uri,omitempty"`
}

func NewOpenIdConnectApplicationSettingsClient() *OpenIdConnectApplicationSettingsClient {
	return &OpenIdConnectApplicationSettingsClient{}
}

func (a *OpenIdConnectApplicationSettingsClient) IsApplicationInstance() bool {
	return true
}
