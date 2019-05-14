package models

// ADMINTABLE - Table name in onchain
const ADMINTABLE = "Admin_"

// Proposal - struct
type Admin struct {
	AdminID   string `json:"AdminID"`
	Name      string `json:"Name"`
	PublicKey string `json:"PublicKey"`
}
