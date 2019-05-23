package models

// ADMINTABLE - Table name in onchain
const ADMINTABLE = "Admin_"

// Admin - struct
type Admin struct {
	AdminID   string `json:"AdminID"`
	Name      string `json:"Name"`
	PublicKey string `json:"PublicKey"`
	Status    string `json:"Status"`
}
