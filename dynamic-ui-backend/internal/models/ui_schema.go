package models

import "time"

type UISchemaResponse struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data,omitempty"`
	Message  string      `json:"message,omitempty"`
	Version  string      `json:"version"`
	CachedAt time.Time   `json:"cached_at"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

type VersionInfo struct {
	Success          bool      `json:"success"`
	AppVersion       string    `json:"app_version"`
	MinVersion       string    `json:"min_version"`
	ForceUpdate      bool      `json:"force_update"`
	AvailableScreens []string  `json:"available_screens"`
	UpdatedAt        time.Time `json:"updated_at"`
}
