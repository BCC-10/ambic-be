package response

type Err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Res struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload"`
}

const (
	EmailExist      = "email already exist"
	UsernameExist   = "username already exist"
	RegisterSuccess = "user register successfully"
	OTPSent         = "OTP sent"
	UserVerified    = "OTP verified"
	InvalidOTP      = "invalid or expired OTP"
	VerifySuccess   = "user verified successfully"
	LoginSuccess    = "user login successfully"
)
