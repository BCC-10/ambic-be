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
	UsernameExist   = "username already exist"
	EmailExist      = "email already exist"
	RegisterSuccess = "user register successfully"

	OTPSent         = "OTP sent"
	InvalidOTP      = "invalid or expired OTP"
	UserVerified    = "OTP verified"
	VerifySuccess   = "user verified successfully"
	UserNotVerified = "user is not verified"

	IncorrectIdentifier = "email or password is incorrect"
	LoginSuccess        = "user login successfully"
)
