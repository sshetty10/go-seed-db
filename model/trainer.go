package model

type Trainer struct {
	ID           string `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v3()"`
	Name         string `json:"name"`
	Age          int    `json:"age"`
	City         string `json:"city"`
	LicenseID    string `json:"licenseID"`
	LicenseState string `json:"licenseState" gorm:"-"`
}
