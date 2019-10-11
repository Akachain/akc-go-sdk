package model

// ApprovalTable - Table name
const ApprovalTable = "Approval"

// Approval ...
type Approval struct {
	ApprovalID string `json:"ApprovalID"`
	ProposalID string `json:"ProposalID"`
	ApproverID string `json:"ApproverID"`
	Challenge  string `json:"Challenge"`
	Signature  string `json:"Signature"`
	Message    string `json:"Message"`
	Status     string `json:"Status"`
	CreatedAt  string `json:"CreatedAt"`
}
