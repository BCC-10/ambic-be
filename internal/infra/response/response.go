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
	UserNotExists   = "user not exists"

	OTPSent       = "OTP sent"
	InvalidOTP    = "invalid or expired OTP"
	VerifySuccess = "user verified successfully"

	IncorrectIdentifier = "email or password is incorrect"
	LoginSuccess        = "user login successfully"

	ForgotPasswordSuccess = "reset password link sent"
	MissingToken          = "missing token"
	ResetPasswordSuccess  = "password reset successfully"

	UpdateSuccess        = "user updated successfully"
	IncorrectOldPassword = "incorrect old password"
)

// middleware
const (
	LimitExceeded   = "too many requests"
	UserNotVerified = "user is not verified"
	UserVerified    = "user already verified"
	InvalidToken    = "invalid or expired token"
)
