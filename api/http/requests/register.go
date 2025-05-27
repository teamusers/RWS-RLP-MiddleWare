package requests

import (
	"errors"
	"lbe/codes"
	model "lbe/model"
)

// VerifyUseExistenceRequest is the payload to verify if an email is already registered.
// If not registered, an OTP will be sent to this email.
type VerifyUserExistence struct {
	// Email address to check for existing registration.
	Email string `json:"email" binding:"required" example:"user@example.com"`
}

type RegisterUser struct {
	User       model.User `json:"user"`
	SignUpType string     `json:"sign_up_type" example:"NEW"`
	RegId      string     `json:"reg_id" example:"123456"`
}

func (r *RegisterUser) Validate() error {
	signUpType := r.SignUpType

	if !codes.IsValidSignUpType(signUpType) {
		return errors.New("invalid sign_up_type provided")
	}

	if signUpType == codes.SignUpTypeTM {
		if r.User.UserProfile.EmployeeNumber == "" {
			return errors.New("user.user_profile.employee_number is required")
		}
	} else if signUpType == codes.SignUpTypeGRCMS {
		if r.RegId == "" {
			return errors.New("reg_id is required")
		}
	} else {
		if r.User.Email == "" {
			return errors.New("user.email is required")
		}
		if r.User.FirstName == "" {
			return errors.New("user.first_name is required")
		}
		if r.User.LastName == "" {
			return errors.New("user.last_name is required")
		}
		if r.User.DateOfBirth == nil {
			return errors.New("user.dob is required")
		}
		if r.User.PhoneNumbers == nil {
			return errors.New("user.phone_numbers is required")
		} else {
			if len(r.User.PhoneNumbers) == 0 || r.User.PhoneNumbers[0].PhoneNumber == "" {
				return errors.New("user.phone_numbers must be properly populated")
			}
		}
		if r.User.UserProfile.CountryCode == "" {
			return errors.New("user.user_profile.country_code is required")
		}
		if r.User.UserProfile.CountryName == "" {
			return errors.New("user.user_profile.country_name is required")
		}
		// marketing preference flags will be false by default

		if signUpType == codes.SignUpTypeGR {
			if r.User.GrProfile == nil {
				return errors.New("gr_profile is required")
			} else {
				if r.User.GrProfile.Class == "" {
					return errors.New("gr_profile.class is required")
				}
			}
		}
	}

	return nil
}

type VerifyGrUser struct {
	User model.User `json:"user" binding:"required"`
}

func (r *VerifyGrUser) Validate() error {
	if r.User.GrProfile == nil {
		return errors.New("gr_profile is required")
	} else {
		if r.User.GrProfile.Id == "" {
			return errors.New("gr_profile.id is required")
		}
		if r.User.GrProfile.Pin == "" {
			return errors.New("gr_profile.pin is required")
		}
	}

	return nil
}

type VerifyGrCmsUser struct {
	User model.User `json:"user" binding:"required"`
}

func (r *VerifyGrCmsUser) Validate() error {
	if r.User.GrProfile == nil {
		return errors.New("gr_profile is required")
	} else {
		if r.User.GrProfile.Id == "" {
			return errors.New("gr_profile.id is required")
		}
		if r.User.GrProfile.Class == "" {
			return errors.New("gr_profile.class is required")
		}
	}

	if r.User.Email == "" {
		return errors.New("user.email is required")
	}
	if r.User.FirstName == "" {
		return errors.New("user.first_name is required")
	}
	if r.User.LastName == "" {
		return errors.New("user.last_name is required")
	}
	if r.User.DateOfBirth == nil {
		return errors.New("user.dob is required")
	}
	if r.User.PhoneNumbers == nil {
		return errors.New("user.phone_numbers is required")
	} else {
		if len(r.User.PhoneNumbers) == 0 || r.User.PhoneNumbers[0].PhoneNumber == "" {
			return errors.New("user.phone_numbers must be properly populated")
		}
	}
	if r.User.UserProfile.CountryCode == "" {
		return errors.New("user.user_profile.country_code is required")
	}
	if r.User.UserProfile.CountryName == "" {
		return errors.New("user.user_profile.country_name is required")
	}

	return nil
}
