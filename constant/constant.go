package constant

const (
	EMPTY_STRING              = ""
	INVALID_TOKEN             = "This Token is Invalid"
	INVALID_CLAIMS            = "Token Cliams are Invalid"
	LEAVE_TEAM                = "Team Left Successfully."
	MEMBERS_ADDED_TO_TEAM     = "Members Added to Team."
	MEMBERS_REMOVED_FROM_TEAM = "Members Removed from Team."
	OTP_SENT                  = "OTP Sent to given Email ID Successfully."
	TOKEN_RESET_SUCCEED       = "Token Reset Done Successfully."
	TASK_CREATED              = "Task Created Successfully."
	TASK_UPDATED              = "Task Updated Successfully."
	TEAM_CREATED              = "Team Created Successfully."
	USER_REGISTRATION_SUCCEED = "User Registration Done Successfully."
	USER_LOGIN_SUCCEED        = "User Login Done Successfully."
	USER_PROFILE_UPDATED      = "User Profile Updated Successfully."
	USER_MAIL_QUEUE           = "user-mail-queue"
	OTP_VERIFICATION_SUCCEED  = "OTP Verification Done Successfully, You can proceed Further."
)

const (
	PG_Duplicate_Error_Code = "23505"
	PG_NO_ROWS = "no rows in result set"
)

type contextKey string

var (
	TokenKey        = contextKey("token")
	UserIdKey       = contextKey("userId")
	SocketServerKey = contextKey("socketServer")
)

var (
	TEAM_ID                 = "TeamID"
	TASK_ID                 = "TaskID"
	URL_PARAM_CONVERT_ERROR = "strconv.Atoi: parsing"
)
