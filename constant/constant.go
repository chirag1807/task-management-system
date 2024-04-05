package constant

const (
	INVALID_TOKEN             = "This Token is Invalid"
	INVALID_CLAIMS            = "Token Cliams are Invalid"
	LEFT_TEAM                 = "Team Left Successfully."
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
)

type contextKey string

var (
	TokenKey        = contextKey("token")
	UserIdKey       = contextKey("userId")
	SocketServerKey = contextKey("socketServer")
)

// below I have declared keys for body, query params and url params validations.
const (
	// For User Model
	IdKey        = "id"
	FirstNameKey = "firstName"
	LastNameKey  = "lastName"
	BioKey       = "bio"
	EmailKey     = "email"
	PasswordKey  = "password"
	ProfileKey   = "profile"

	// For OTP Model
	OTPIdKey   = "id"
	OTPCodeKey = "otp"

	// For Team Model
	TeamNameKey      = "teamDetails.name"
	TeamProfileKey   = "teamDetails.teamProfile"
	TeamMembersKey   = "teamMembers"
	TeamMembersIdKey = "teamMembers.memberID"

	// For Team Members Model
	TeamIdKey       = "teamID"
	TeamMemberIdKey = "memberID"

	// For Task Model
	TaskIdKey             = "id"
	TitleKey              = "title"
	DescriptionKey        = "description"
	DeadlineKey           = "deadline"
	AssigneeIndividualKey = "assigneeIndividual"
	AssigneeTeamKey       = "assigneeTeam"
	StatusKey             = "status"
	PriorityKey           = "priority"

	// For Query Params
	LimitKey          = "limit"
	OffsetKey         = "offset"
	SearchKey         = "search"
	StatusFilterKey   = "status"
	SortByFilterKey   = "sortByFilter"
	SortByCreateAtKey = "sortByCreatedAt"
)
