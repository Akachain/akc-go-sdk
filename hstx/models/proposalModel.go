package models

// PROPOSALTABLE - Table name in onchain
const PROPOSALTABLE = "Proposal_"

// Proposal - struct
type Proposal struct {
	ProposalID string `json:"ProposalID"`
	Data       string `json:"Data"`
}
