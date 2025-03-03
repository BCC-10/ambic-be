package response

type Err struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

type Res struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload"`
}

// Auth Domain
const (
	UsernameExist   = "Username already exists"
	EmailExist      = "Email already exists"
	RegisterSuccess = "User registered successfully"
	UserNotExists   = "User does not exist"

	VerificationLinkSent = "Verification link has been sent"
	VerifySuccess        = "User verified successfully"

	IncorrectIdentifier = "Email or password is incorrect"
	LoginSuccess        = "User logged in successfully"
	InvalidState        = "State is invalid"
	OAuthLoginSuccess   = "OAuth login URL generated successfully"

	ForgotPasswordSuccess = "Password reset link has been sent"
	MissingToken          = "Authentication token is missing"
	ResetPasswordSuccess  = "Password has been reset successfully"

	IncorrectOldPassword = "Old password is incorrect"
)

// middleware
const (
	LimitExceeded      = "Too many requests. Please try again later."
	UserNotVerified    = "User is not verified"
	UserVerified       = "User is already verified"
	InvalidToken       = "Token is invalid or has expired"
	InvalidTokenFormat = "Token format is invalid"
)

// User Domain
const (
	InvalidDateFormat = "Date format is invalid"
	UpdateSuccess     = "User profile updated successfully"
)

// Partner Domain
const (
	PartnerRegisterSuccess = "Partner registered successfully"
)

// Product Domain
const (
	ProductCreateSuccess       = "Product created successfully"
	ProfileNotFilledCompletely = "Profile is not filled completely"
)
