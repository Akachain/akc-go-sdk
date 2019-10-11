package model

// AdminTable - Table name
const AdminTable = "Admin"

// Admin ...
type Admin struct {
	AdminID string `json:"AdminID"`
	Name    string `json:"Name"`
	Status  string `json:"Status"`
}
