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

// Auth Domain
const (
	UsernameExist   = "username already exist"
	EmailExist      = "email already exist"
	RegisterSuccess = "user register successfully"

	OTPSent       = "OTP sent"
	InvalidOTP    = "invalid or expired OTP"
	VerifySuccess = "user verified successfully"

	IncorrectIdentifier = "email or password is incorrect"
	LoginSuccess        = "user login successfully"

	ResetPasswordSuccess = "password reset successfully"
)

// middleware
const (
	UserNotVerified = "user is not verified"
	UserVerified    = "user already verified"
	MissingToken    = "missing token"
	InvalidToken    = "invalid token"
)
