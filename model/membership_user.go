package model

import "time"

type MembershipUser struct {
	ID                    uint64    `json:"id"`
	ExternalID            string    `json:"external_id"`
	ExternalTYPE          string    `json:"external_id_type"`
	Email                 string    `json:"email"`
	BurnPin               uint64    `json:"burn_pin"`
	SessionToken          string    `json:"session_token"`
	SessionExpiry         int64     `json:"session_expiry"`
	GR_ID                 string    `json:"gr_id"`
	RLP_ID                string    `json:"rlp_id"`
	RWS_Membership_ID     string    `json:"rws_membership_id"`
	RWS_Membership_Number uint64    `json:"rws_membership_number"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}