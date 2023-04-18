package entities

import (
	"html"
	"strings"
	"time"
)

type Model struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (m *Model) BeforeUpdate() {
	m.Title = html.EscapeString(strings.TrimSpace(m.Title))
	m.Description = html.EscapeString(strings.TrimSpace(m.Description))
	m.UpdatedAt = time.Now()
}

func (m *Model) Prepare() {
	m.Title = html.EscapeString(strings.TrimSpace(m.Title))
	m.Description = html.EscapeString(strings.TrimSpace(m.Description))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
}

func (m *Model) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if m.Title == "" || m.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
	default:
		if m.Title == "" || m.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
	}
	return errorMessages
}
