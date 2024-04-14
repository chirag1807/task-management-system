package response

import (
	"time"
)

// OTP model info
// @Description OTP information with id, 4-digit otp and expiry time of otp.
type OTP struct {
	ID            int64     `json:"id" example:"974751326021189896"`
	OTP           int       `json:"otp" example:"5896"`
	OTPExpiryTime time.Time `json:"otpExpireTime,omitempty" example:"2024-03-25T14:29:19.000Z"`
	Email         string    `json:"email" example:"chiragmakwana@gmail.com"`
	IsVerified    bool      `json:"isVerified" example:"true"`
}
