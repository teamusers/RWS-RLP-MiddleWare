package responses

import (
	"lbe/model"
	"time"
)

// CMS Member holds all the user fields
type GRProfilePayload struct {
	MemberNo                             string    `json:"MemberNo"`
	FirstName                            string    `json:"FirstName"`
	LastName                             string    `json:"LastName"`
	DateOfBirth                          time.Time `json:"DateOfBirth"`
	CardDisplayName                      string    `json:"CardDisplayName"`
	MemberClassID                        int64     `json:"MemberClassID"`
	MemberClassCode                      string    `json:"MemberClassCode"`
	MemberClassLevel                     int       `json:"MemberClassLevel"`
	PremiumMemberFlag                    bool      `json:"PremiumMemberFlag"`
	PreferredLanguageID                  int64     `json:"PreferredLanguageID"`
	IdentificationTypeID                 int64     `json:"IdentificationTypeID"`
	IdentificationNo                     string    `json:"IdentificationNo"`
	IdentificationIssuanceCountryISOCode string    `json:"IdentificationIssuanceCountryISOCode"`
	NationalityCountryISOCode            string    `json:"NationalityCountryISOCode"`
	IsSCPRFlag                           bool      `json:"IsSCPRFlag"`
	ResidentialStatusID                  int64     `json:"ResidentialStatusID"`
	CDDRequiredFlag                      bool      `json:"CDDRequiredFlag"`
	CDDDueDateTime                       time.Time `json:"CDDDueDateTime"`
	CDDRiskCategoryID                    int64     `json:"CDDRiskCategoryID"`
	ConsentStatus                        bool      `json:"ConsentStatus"`
	ContactOptionEmail                   bool      `json:"ContactOptionEmail"`
	ContactOptionEmailStatus             int       `json:"ContactOptionEmailStatus"`
	ContactOptionPhone                   bool      `json:"ContactOptionPhone"`
	ContactOptionAddress                 bool      `json:"ContactOptionAddress"`
	CasinoSiteID                         int64     `json:"CasinoSiteID"`
	EmailAddress                         string    `json:"EmailAddress"`
	EmailModifiedDateTime                time.Time `json:"EmailModifiedDateTime"`
	OccupationCategoryID                 int64     `json:"OccupationCategoryID"`
	OccupationDetails                    string    `json:"OccupationDetails"`
	BusinessNatureCategoryID             int64     `json:"BusinessNatureCategoryID"`
	BusinessNatureDetails                string    `json:"BusinessNatureDetails"`
	AnnualIncomeRangeCategoryID          int64     `json:"AnnualIncomeRangeCategoryID"`
}

func (g *GRProfilePayload) MapCmsToLbeGrProfile() model.GrProfile {
	return model.GrProfile{
		Id:          g.MemberNo,
		Class:       g.MemberClassCode,
		FirstName:   g.FirstName,
		LastName:    g.LastName,
		Email:       g.EmailAddress,
		DateOfBirth: model.Date(g.DateOfBirth),
		//TODO: Add mobile code and number
	}
}
