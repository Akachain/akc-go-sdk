package models

// QUORUMTABLE - Table name in onchain
const QUORUMTABLE = "Quorum_"

// Proposal - struct
type Quorum struct {
	QuorumID   string `json:"QuorumID"`
	AdminID    string `json:"AdminID"`
	ProposalID string `json:"ProposalID"`
	Status     string `json:"Status"`
}