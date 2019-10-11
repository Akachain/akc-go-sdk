package model

// SuperAdminTable - Table name
const SuperAdminTable = "SuperAdmin"

// SuperAdmin - Super Admin
type SuperAdmin struct {
	SuperAdminID string `json:"SuperAdminID"` // keyhandle
	Name         string `json:"Name"`
	PublicKey    string `json:"PublicKey"` // format: pem
	Status       string `json:"Status"`
}
