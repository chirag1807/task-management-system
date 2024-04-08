package request

import (
	"time"
)

// OTP model info
// @Description OTP information with id, 4-digit otp and expiry time of otp.
type OTP struct {
	ID            int64     `json:"id" example:"974751326021189896" validate:"required,number"`
	OTP           int       `json:"otp" example:"5896" validate:"required,number,min=1000,max=9999"`
	OTPExpiryTime time.Time `json:"otpExpireTime,omitempty" example:"2024-03-25T14:29:19.000Z"`
}
