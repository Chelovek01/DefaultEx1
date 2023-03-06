package models

type ClientInfo struct {
	Address string `json:"address,omitempty"`
	City    string `json:"city,omitempty"`
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Region  string `json:"region,omitempty"`
}
