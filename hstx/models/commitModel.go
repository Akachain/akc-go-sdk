package models

// COMMITTABLE - Table name in onchain
const COMMITTABLE = "Commit_"

// Commit - struct
type Commit struct {
	CommitID   string   `json:"CommitID"`
	AdminID    string   `json:"AdminID"`
	ProposalID string   `json:"ProposalID"`
	QuorumList []string `json:"QuorumList"`
	Status     string   `json:"Status"`
}
