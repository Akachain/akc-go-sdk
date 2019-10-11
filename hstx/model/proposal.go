package model

// ProposalTable - Table name
const ProposalTable = "Proposal"

// Proposal - struct
type Proposal struct {
	ProposalID string `json:"ProposalID"` // set
	Message    string `json:"Message"`    // args[0]
	CreatedBy  string `json:"CreatedBy"`  // args[1]: ID of Admin/SAdmin
	Status     string `json:"Status"`     // set
	CreatedAt  string `json:"CreatedAt"`  // args[2]
	UpdatedAt  string `json:"UpdatedAt"`  // args[3]
}
