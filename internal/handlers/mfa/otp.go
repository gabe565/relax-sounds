package mfa

import (
	"errors"
	"net/http"

	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tools/security"
)

func RegisterRoutes(e *core.ServeEvent) {
	e.Router.POST("/api/mfa/request-otp", RequestOTP)
}

type requestOTPForm struct {
	MFAId string `form:"mfaId" json:"mfaId"`
}

type RequestOTPResponse struct {
	OTPId string `json:"otpId"`
	Email string `json:"email"`
}

func RequestOTP(e *core.RequestEvent) error {
	form := &requestOTPForm{}
	if err := e.BindBody(form); err != nil {
		return e.BadRequestError("An error occurred while loading the submitted data.", err)
	}
	if err := validation.ValidateStruct(form, validation.Field(&form.MFAId, validation.Required)); err != nil {
		return e.BadRequestError("An error occurred while validating the submitted data.", err)
	}

	mfa, err := e.App.FindMFAById(form.MFAId)
	if err != nil {
		return e.BadRequestError("Invalid or expired MFA session.", err)
	}

	collection, err := e.App.FindCachedCollectionByNameOrId(mfa.CollectionRef())
	if err != nil {
		return e.BadRequestError("Invalid MFA session.", err)
	}

	if !collection.MFA.Enabled || mfa.HasExpired(collection.MFA.DurationTime()) {
		return e.BadRequestError("Invalid or expired MFA session.", nil)
	}

	if !collection.OTP.Enabled {
		return e.ForbiddenError("The collection is not configured to allow OTP authentication.", nil)
	}

	if mfa.Method() == core.MFAMethodOTP {
		return e.BadRequestError("A different authentication method is required.", nil)
	}

	record, err := e.App.FindRecordById(collection, mfa.RecordRef())
	if err != nil {
		return e.BadRequestError("Invalid MFA session.", err)
	}

	if record.Email() == "" {
		return e.BadRequestError("The account has no email address.", nil)
	}

	if !e.App.IsDev() {
		otps, err := e.App.FindAllOTPsByRecord(record)
		if err != nil {
			return e.InternalServerError("Failed to fetch previous record OTPs.", err)
		}

		var recent int
		for _, existing := range otps {
			if !existing.HasExpired(collection.OTP.DurationTime()) {
				recent++
			}
		}
		if recent > 9 {
			return e.TooManyRequestsError("Too many code requests. Use the most recent code or try again later.", nil)
		}
	}

	password := security.RandomStringWithAlphabet(collection.OTP.Length, "1234567890")

	otp := core.NewOTP(e.App)
	otp.SetCollectionRef(collection.Id)
	otp.SetRecordRef(record.Id)
	otp.SetPassword(password)
	if err := e.App.Save(otp); err != nil {
		return e.InternalServerError("Failed to create OTP.", err)
	}

	if err := mails.SendRecordOTP(e.App, record, otp.Id, password); err != nil {
		return e.InternalServerError("Failed to send OTP email.", errors.Join(err, e.App.Delete(otp)))
	}

	return e.JSON(http.StatusOK, RequestOTPResponse{
		OTPId: otp.Id,
		Email: record.Email(),
	})
}
