// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

// List of factors that are applicable to Okta Identity Engine (OIE)
var AuthenticatorProviders = []string{
	// NOTE: some authenticator types are available by feature flag on the org only
	DuoFactor,
	ExternalIdpFactor,
	GoogleOtpFactor,
	CustomOtpFactor,
	OktaEmailFactor,
	OktaPasswordFactor, // NOTE: Not configurable in OIE policies (Handle downstream as necessary)
	OktaVerifyFactor,
	OnPremMfaFactor,
	PhoneNumberFactor,
	RsaTokenFactor,
	SecurityQuestionFactor,
	WebauthnFactor,
	YubikeyTokenFactor,
	SmartCardIdpFactor,
}
