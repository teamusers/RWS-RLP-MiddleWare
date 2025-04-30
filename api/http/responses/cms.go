package responses

import (
	"time"
)

// CMS Member holds all the user fields
type GRProfilePayload struct {
	MemberNo                             string    `json:"member_no"`
	FirstName                            string    `json:"first_name"`
	LastName                             string    `json:"last_name"`
	DateOfBirth                          time.Time `json:"date_of_birth"`
	CardDisplayName                      string    `json:"card_display_name"`
	MemberClassID                        int64     `json:"member_class_id"`
	MemberClassCode                      string    `json:"member_class_code"`
	MemberClassLevel                     int       `json:"member_class_level"`
	PremiumMemberFlag                    bool      `json:"premium_member_flag"`
	PreferredLanguageID                  int64     `json:"preferred_language_id"`
	IdentificationTypeID                 int64     `json:"identification_type_id"`
	IdentificationNo                     string    `json:"identification_no"`
	IdentificationIssuanceCountryISOCode string    `json:"identification_issuance_country_iso_code"`
	NationalityCountryISOCode            string    `json:"nationality_country_iso_code"`
	IsSCPRFlag                           bool      `json:"is_scpr_flag"`
	ResidentialStatusID                  int64     `json:"residential_status_id"`
	CDDRequiredFlag                      bool      `json:"cdd_required_flag"`
	CDDDueDateTime                       time.Time `json:"cdd_due_date_time"`
	CDDRiskCategoryID                    int64     `json:"cdd_risk_category_id"`
	ConsentStatus                        bool      `json:"consent_status"`
	ContactOptionEmail                   bool      `json:"contact_option_email"`
	ContactOptionEmailStatus             int       `json:"contact_option_email_status"`
	ContactOptionPhone                   bool      `json:"contact_option_phone"`
	ContactOptionAddress                 bool      `json:"contact_option_address"`
	CasinoSiteID                         int64     `json:"casino_site_id"`
	EmailAddress                         string    `json:"email_address"`
	EmailModifiedDateTime                time.Time `json:"email_modified_date_time"`
	OccupationCategoryID                 int64     `json:"occupation_category_id"`
	OccupationDetails                    string    `json:"occupation_details"`
	BusinessNatureCategoryID             int64     `json:"business_nature_category_id"`
	BusinessNatureDetails                string    `json:"business_nature_details"`
	AnnualIncomeRangeCategoryID          int64     `json:"annual_income_range_category_id"`
}
