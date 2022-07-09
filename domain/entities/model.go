package entities

import (
	"html"
	"strings"
	"time"
)

type Model struct {
	ID          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	UserID      uint64     `gorm:"size:100;not null;" json:"user_id"`
	Title       string     `gorm:"size:100;not null;unique" json:"title"`
	Description string     `gorm:"text;not null;" json:"description"`
	ModelFile   string     `gorm:"size:255;null;" json:"model_file"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (m *Model) BeforeSave() {
	m.Title = html.EscapeString(strings.TrimSpace(m.Title))
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
		if m.Description == "" || m.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
	default:
		if m.Title == "" || m.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
		if m.Description == "" || m.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
	}
	return errorMessages
}
