package response

type Err struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload,omitempty"`
}

type Res struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload,omitempty"`
}

const (
	EntityTooLarge = "Entity too large, max size is %d MB"
)

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

// others
const (
	LimitExceeded      = "Too many requests. Please try again later."
	UserNotVerified    = "User is not verified"
	UserVerified       = "User is already verified"
	InvalidToken       = "Token is invalid or has expired"
	InvalidTokenFormat = "Token format is invalid"
	PhotoSizeLimit     = "Photo size is too large"
	PhotoOnly          = "Only photo is allowed"
	InvalidUUID        = "UUID is invalid"
)

// User Domain
const (
	ShowUserSuccess    = "User retrieved successfully"
	InvalidDateFormat  = "Date format is invalid"
	UpdateSuccess      = "User profile updated successfully"
	PhoneAlreadyExists = "Phone number already exists"
)

// Partner Domain
const (
	AlreadyRegisteredAsPartner = "Partner is already registered"
	NotPartner                 = "User is not a partner"
	PartnerNotVerified         = "Partner is not verified"
	PartnerRegisterSuccess     = "Partner registered successfully"
	PartnerVerificationSuccess = "Partner verified successfully"
	PartnerVerified            = "Partner is already verified"
	PartnerNotExists           = "Partner does not exist"
	GetPartnerSuccess          = "Partner retrieved successfully"
	UpdatePartnerPhotoSuccess  = "Partner photo updated successfully"
	InvalidBusinessType        = "Business type is invalid"
	GetAutoCompleteSuccess     = "Auto complete suggestions retrieved successfully"
)

// Product Domain
const (
	GetProductSuccess          = "Items retrieved successfully"
	CreateProductSuccess       = "Product created successfully"
	UpdateProductSuccess       = "Product updated successfully"
	ProfileNotFilledCompletely = "Profile is not filled completely"
	ProductNotExists           = "Product does not exist"
	RatingNotBelongToPartner   = "Product does not belong to partner"
	DeleteProductSuccess       = "Product deleted successfully"
	ProductAlreadyExists       = "Product already exists"
	InvalidDateTime            = "Date time format is invalid, format should be YYYY-MM-DD HH:MM:SS"
	ProductNotFound            = "The product with id %s does not exist"
)

// Rating Domain
const (
	GetRatingSuccess        = "Ratings retrieved successfully"
	RatingNotBelongToUser   = "Rating does not belong to user"
	CreateRatingSuccess     = "Rating created successfully"
	UpdateRatingSuccess     = "Rating updated successfully"
	UserAlreadyRated        = "User already rated"
	RatingDeleteSuccess     = "Rating deleted successfully"
	RatingNotFound          = "Rating does not exist"
	RatingNotExists         = "Rating does not exist"
	UserNotPurchasedProduct = "User has not purchased the product or the transaction is not completed yet"
)

// Business Type Domain
const (
	BusinessTypeEmpty      = "Business type is empty"
	GetBusinessTypeSuccess = "Business types retrieved successfully"
)

// Transaction Domain
const (
	GetTransactionSuccess    = "Transactions retrieved successfully"
	CreateTransactionSuccess = "Transaction created successfully"
	InsufficientStock        = "Insufficient stock for %s"
	InvalidQty               = "Quantity must be greater than 0"
	MissingProductID         = "Product ID is missing"
	MissingTransactionItems  = "Transaction items are missing"
	TransactionNotFound      = "Transaction not found"
)

// Transaction Detail Domain
const ()
