package models

import "time"

const (
	MandateStatusActive         = "active"
	MandateStatusElectedPending = "elected_pending"
	MandateStatusFormer         = "former"
)

// Politician is the public read model for an identity and all of its known federal mandates.
// It is populated from the internal classification API's Redis projection.
type Politician struct {
	ID       string                 `json:"id"`
	FullName string                 `json:"fullName"`
	Mandates []ParliamentaryMandate `json:"mandates"`
}

type ParliamentaryMandate struct {
	ID               string     `json:"id"`
	LegislativeHouse string     `json:"legislativeHouse"`
	Office           string     `json:"office"`
	State            string     `json:"state"`
	PartyAcronym     string     `json:"partyAcronym"`
	PartyDirection   string     `json:"partyDirection,omitempty"`
	PartyIdeology    string     `json:"partyIdeology,omitempty"`
	Status           string     `json:"status"`
	ElectionYear     *int       `json:"electionYear,omitempty"`
	StartsAt         *time.Time `json:"startsAt,omitempty"`
	EndsAt           *time.Time `json:"endsAt,omitempty"`
}
