package models

// COMMITTABLE - Table name in onchain
const COMMITTABLE = "Commit_"

// Commit - struct
type Commit struct {
	CommitID   string   `json:"CommitID"`
	ProposalID string   `json:"ProposalID"`
	QuorumID   []string `json:"QuorumID"`
	Status     string   `json:"Status"`
}
