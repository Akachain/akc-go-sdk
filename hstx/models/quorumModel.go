package models

// QUORUMTABLE - Table name in onchain
const QUORUMTABLE = "Quorum_"

// Proposal - struct
type Quorum struct {
	AdminID    string `json:"AdminID"`
	QuorumID   string `json:"QuorumID"`
	ProposalID string `json:"ProposalID"`
	Status     string `json:"Status"`
}
