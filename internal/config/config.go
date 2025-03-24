package config

type Config struct {
	DbURL string `json:"db_url"`
	CurrentUsername string `json:"current_user_name,omitempty"`
}