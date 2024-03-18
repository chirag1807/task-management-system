package response

import (
	"time"
)

type OTP struct {
	ID            int64     `json:"id"`
	OTP           int       `json:"otp"`
	OTPExpiryTime time.Time `json:"otpExpireTime,omitempty"`
}
