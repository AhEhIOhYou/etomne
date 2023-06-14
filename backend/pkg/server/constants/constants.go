package constants

const (

	//Server proccess
	ServerStartSuccess    = "server successfully started"
	ServerStartErr        = "server error :%v"
	ServerNotFoundEnvFile = ".env file not found"

	//Database proccess
	DatabaseConnectionStart   = "connecting to database..."
	DatabaseConnectionError   = "database connection error: %v"
	DatabaseConnectionSuccess = "connected to database"
	DatabaseError             = "database error: %v"

	//Handler error message
	Failed = "request failed: %v"

	//User errors
	UserIDInvalid                  = "user id invalid"
	UsernameCantBeEmpty            = "username can't be empty"
	UsernameInvalid                = "username is invalid"
	UsernameIsAlreadyTaken         = "username is already taken"
	UsernameLenError               = "username should be between 3-20 characters"
	EmailCantBeEmpty               = "email can't be empty"
	EmailWrongFormat               = "email worng format"
	PasswordCantBeEmpty            = "password can't be empty"
	PasswordIncorrect              = "password is incorrect"
	UsernameAndPasswordCantBeEmpty = "username and password can't be empty"
	Unauthorized                   = "unauthorized"
	AccountDoesNotExist            = "account does not exist"

	//Auth errors
	CannotGetUUID       = "cannot get uuid"
	RefreshTokenExpired = "refresh token expired"

	//Model errors
	ModelTitleCantBeEmpty = "title can't be empty"
	ModelNotAvaliable     = "model is not available to you"

	//File errors
	FileTitleCantBeEmpty = "title can't be empty"
	FileURLError         = "something wrong with file path"
	FileNotAvaliable     = "file not avaliable"

	//Security errors
	PasswordHashError   = "error hashing password: %v"
	PasswordVerifyError = "error verifying password: %v"
	NotEnoughRights = "not enough rights"

	UsernameIsAvailable    = "Username is available"
	RegistrationSuccessful = "Registration successful"
	LoginSuccessful        = "login successful"
	LogoutSuccessful       = "logout successful"
	YouAreLoggedIn         = "you are logged in"
	DeletedSuccessful  = "successfully deleted"
)
