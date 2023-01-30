package model

type TechInfo struct {
	Company  string `json:"company"`
	Product  string `json:"product"`
	Lang     string `json:"lang"`
	Server   string `json:"server"`
	Category string `json:"category"`
}
type FingerPrint struct {
	Infos   string
	Matches []string
}
