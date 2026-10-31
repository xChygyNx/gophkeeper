package model

type LoginPassword struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Login       string `json:"Login"`
	Password    string `json:"Password"`
}
