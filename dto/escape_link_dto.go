package dto

import (
	"time"

	"gorm.io/datatypes"
)

type ZoneEscapeLinksRequest struct {
	WebView datatypes.JSON `json:"webview" binding:"required"`
}

type ZoneEscapeLinksResponse struct {
	Id        uint           `json:"id"`
	Uuid      string         `json:"uuid"`
	WebView   datatypes.JSON `json:"webview"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
