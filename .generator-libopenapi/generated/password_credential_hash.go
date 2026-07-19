// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordCredentialHash represents the PasswordCredentialHash schema
// Specifies a hashed password to import into Okta. This allows an existing password to be imported into Okta directly from some other store. Okta supports the BCRYPT, SHA-512, SHA-256, SHA-1, MD5, an...
type PasswordCredentialHash struct {
	// The number of iterations used when hashing passwords using PBKDF2. Must be >= 4096. Only required for PBKDF2 algorithm.
	IterationCount int `json:"iterationCount,omitempty"`
	// Size of the derived key in bytes. Only required for PBKDF2 algorithm.
	KeySize int `json:"keySize,omitempty"`
	// Only required for salted hashes. For BCRYPT, this specifies Radix-64 as the encoded salt used to generate the hash, which must be 22 characters long. For other salted hashes, this specifies the Bas...
	Salt string `json:"salt,omitempty"`
	// Specifies whether salt was pre- or postfixed to the password before hashing. Only required for salted algorithms.
	SaltOrder string `json:"saltOrder,omitempty"`
	// For SHA-512, SHA-256, SHA-1, MD5, and PBKDF2, this is the actual base64-encoded hash of the password (and salt, if used). This is the Base64-encoded `value` of the SHA-512/SHA-256/SHA-1/MD5/PBKDF2 ...
	Value string `json:"value,omitempty"`
	// Governs the strength of the hash and the time required to compute it. Only required for BCRYPT algorithm.
	WorkFactor int `json:"workFactor,omitempty"`
	Algorithm interface{} `json:"algorithm,omitempty"`
	DigestAlgorithm interface{} `json:"digestAlgorithm,omitempty"`
}
