package request

import (
	"time"
)

type OTP struct {
	ID            int64     `json:"id"`
	OTP           int8      `json:"otp"`
	OTPExpiryTime time.Time `json:"otpExpireTime,omitempty"`
}
