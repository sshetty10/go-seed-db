package model

type Trainer struct {
	ID   string `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
	Key  string `json:"key"`
}
