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
	UsernameExist   = "Username already exist"
	EmailExist      = "Email already exist"
	RegisterSuccess = "User register successfully"
	UserNotExists   = "User not exists"

	OTPSent       = "OTP sent"
	InvalidOTP    = "Invalid or expired OTP"
	VerifySuccess = "User verified successfully"

	IncorrectIdentifier = "Email or password is incorrect"
	LoginSuccess        = "User login successfully"

	ForgotPasswordSuccess = "Reset password link sent"
	MissingToken          = "Missing token"
	ResetPasswordSuccess  = "Password reset successfully"

	UpdateSuccess        = "User updated successfully"
	IncorrectOldPassword = "Incorrect old password"
)

// middleware
const (
	LimitExceeded   = "Too many requests"
	UserNotVerified = "User is not verified"
	UserVerified    = "User already verified"
	InvalidToken    = "Invalid or expired token"
)
