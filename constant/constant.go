package constant

const (
	TOKEN_REFRESH             = "Token Refresh Successfully."
	USER_REGISTRATION_SUCCEED = "User Registration Done Successfully."
	USER_PROFILE_UPDATED      = "User Profile Updated Successfully."
	ARTICLE_ADDED             = "Article Added Successfully."
	ARTICLE_UPDATED           = "Article Updated Successfully."
	ARTICLE_DELETED           = "Article Deleted Successfully."
	ARTICLE_VIEW_INCREASED    = "Article View Increased Successfully."
	ARTICLE_LIKE_ADDED        = "Like Added Successfully."
	ARTICLE_LIKE_REMOVED      = "Like Removed Successfully."
	FOLLOWING_NOW             = "You are Following the Author Now."
	NOT_FOLLOWING_NOW         = "You are Not Following the Author Now."
	TOPIC_ADDED               = "Topic Added Successfully."
	TOPIC_UPDATED             = "Topic Updated Successfully."
	TOPIC_DELETED             = "Topic Deleted Successfully."
	INVALID_TOKEN             = "This Token is Invalid."
	INVALID_CLAIMS            = "Token Cliams are Invalid."
	OTP_SENT                  = "OTP Sent to given Email ID Successfully."
	TEAM_CREATED              = "Team Created Successfully."
	MEMBERS_ADDED_TO_TEAM     = "Members Added to Team."
	MEMBERS_REMOVED_FROM_TEAM = "Members Removed from Team."
	LEFT_TEAM                 = "Team Left Successfully."
	TASK_CREATED              = "Task Created Successfully."
)

type contextKey string

var (
	TokenKey  = contextKey("token")
	UserIdKey = contextKey("userId")
)

// below I have declared keys for body, quer params and url params validations.
const (
	// for user model
	IdKey        = "id"
	FirstNameKey = "firstName"
	LastNameKey  = "lastName"
	BioKey       = "bio"
	EmailKey     = "email"
	PasswordKey  = "password"
	ProfileKey   = "profile"

	// for otp model
	OTPIdKey   = "id"
	OTPCodeKey = "otp"

	// for team model
	TeamNameKey      = "teamDetails.name"
	TeamMembersKey   = "teamMembers"
	TeamMembersIdKey = "teamMembers.memberID"

	// for team members model
	TeamIdKey       = "teamID"
	TeamMemberIdKey = "memberID"

	//for task model
	TitleKey              = "title"
	DescriptionKey        = "description"
	DeadlineKey           = "deadline"
	AssigneeIndividualKey = "assigneeIndividual"
	AssigneeTeamKey       = "assigneeTeam"
	StatusKey             = "status"
	PriorityKey           = "priority"

	// query params
	LimitKey  = "limit"
	OffsetKey = "offset"
)
