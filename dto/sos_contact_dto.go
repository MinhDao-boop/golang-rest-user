package dto

import (
	"golang-rest-user/enums"
	"time"
)

type CreateSOSContactRequest struct {
	Name     string `json:"name" binding:"required"`
	Position string `json:"position" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Note     string `json:"note"`
}

type UpdateSOSContactRequest struct {
	Name     *string `json:"name" patch:"true"`
	Position *string `json:"position" patch:"true"`
	Phone    *string `json:"phone" patch:"true"`
	Note     *string `json:"note" patch:"true"`
}

type ToggleSOSContactRequest struct {
	IsActive *enums.SOSContactStatus `json:"is_active" patch:"true"`
}

type DeleteSOSContactRequest struct {
	UUID []string `json:"uuid"`
}

type SOSContactResponse struct {
	ID        uint                   `json:"id"`
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Position  string                 `json:"position"`
	Phone     string                 `json:"phone"`
	Note      string                 `json:"note"`
	IsActive  enums.SOSContactStatus `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
