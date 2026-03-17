package slack

type (
	AuthDetails struct {
		// SigningSecret
		//
		// Used for HMAC-SHA256 verification of the request signature
		SigningSecret string

		// Timestamp
		//
		// Used for preventing replay attacks
		Timestamp int64

		// RequestSignature
		//
		// Used to confirm the veracity of the request
		RequestSignature string

		// RequestBody
		//
		// Raw request body from Slack, before deserialization
		RequestBody []byte
	}
)
